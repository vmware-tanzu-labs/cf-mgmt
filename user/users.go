package user

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

const ORG_USER string = "organization_user"
const ORG_AUDITOR string = "organization_auditor"
const ORG_MANAGER string = "organization_manager"
const ORG_BILLING_MANAGER string = "organization_billing_manager"
const SPACE_AUDITOR string = "space_auditor"
const SPACE_DEVELOPER string = "space_developer"
const SPACE_MANAGER string = "space_manager"
const SPACE_SUPPORTER string = "space_supporter"

// NewManager -
func NewManager(
	client CFClient,
	cfg config.Reader,
	spaceMgr space.Manager,
	orgReader organizationreader.Reader,
	uaaMgr uaa.Manager,
	peek bool) (Manager, error) {

	supports, err := client.SupportsSpaceSupporterRole()
	if err != nil {
		return nil, err
	}

	mgr := &DefaultManager{
		Peek:                   peek,
		SpaceMgr:               spaceMgr,
		OrgReader:              orgReader,
		UAAMgr:                 uaaMgr,
		Cfg:                    cfg,
		SupportsSpaceSupporter: supports,
	}
	_, logTimings := os.LookupEnv("LOG_TIMINGS")
	if logTimings {
		lo.G.Infof("Logging timings enabled")
		timer := NewCFClientTimer(client)
		mgr.Client = timer
	} else {
		mgr.Client = client
	}
	return mgr, nil
}

type DefaultManager struct {
	Client                 CFClient
	Cfg                    config.Reader
	SpaceMgr               space.Manager
	OrgReader              organizationreader.Reader
	UAAMgr                 uaa.Manager
	Peek                   bool
	LdapMgr                LdapManager
	LdapConfig             *config.LdapConfig
	SupportsSpaceSupporter bool
	UAAUsers               *uaa.Users
	CFUsers                map[string]cfclient.V3User
	OrgRoles               map[string]map[string]*RoleUsers
	SpaceRoles             map[string]map[string]*RoleUsers
}

func (m *DefaultManager) ClearRoles() {
	m.OrgRoles = nil
	m.SpaceRoles = nil
}

func (m *DefaultManager) LogResults() {
	v, ok := m.Client.(*CFClientTimer)
	if ok {
		v.LogResults()
	}
}

func (m *DefaultManager) GetCFUsers() (map[string]cfclient.V3User, error) {
	if m.CFUsers == nil {
		cfUserMap := make(map[string]cfclient.V3User)
		cfUsers, err := m.Client.ListV3UsersByQuery(url.Values{})
		if err != nil {
			return nil, err
		}
		lo.G.Debug("Begin CFUsers")
		for _, cfUser := range cfUsers {
			cfUserMap[cfUser.GUID] = cfUser
			lo.G.Debugf("CFUser with username [%s] and guid [%s]", cfUser.Username, cfUser.GUID)
		}
		lo.G.Debug("End CFUsers")
		m.CFUsers = cfUserMap
	}
	return m.CFUsers, nil
}

func (m *DefaultManager) GetUAAUsers() (*uaa.Users, error) {
	if m.UAAUsers == nil {
		uaaUsers, err := m.UAAMgr.ListUsers()
		if err != nil {
			return nil, err
		}
		m.UAAUsers = uaaUsers
	}
	return m.UAAUsers, nil
}

func (m *DefaultManager) AddUAAUser(user uaa.User) {
	m.UAAUsers.Add(user)
}

func (m *DefaultManager) dumpRolesUsers(entityType string, entityRoles map[string]map[string]*RoleUsers) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		for guid, roles := range entityRoles {
			for roleType, role := range roles {
				for _, user := range role.Users() {
					lo.G.Debugf("User [%s] with GUID[%s] and origin %s in role %s for entity type [%s/%s]", user.UserName, user.GUID, user.Origin, roleType, entityType, guid)
				}
			}
		}
	}
}

func (m *DefaultManager) dumpV3Roles(entityType string, roles []cfclient.V3Role) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		for _, role := range roles {
			lo.G.Debugf("For entity [%s/%s] and role [%s] user guid [%s]", entityType, role.GUID, role.Type, role.Relationships["user"].Data.GUID)
		}
	}
}
func (m *DefaultManager) initializeSpaceUserRolesMap() error {
	spaceV3UsersRolesMap := make(map[string]map[string][]cfclient.V3User)
	query := url.Values{}
	query["per_page"] = []string{"5000"}
	query["order_by"] = []string{"created_at"}
	query["types"] = []string{SPACE_AUDITOR + "," + SPACE_DEVELOPER + "," + SPACE_MANAGER + "," + SPACE_SUPPORTER}
	roles, err := m.Client.ListV3RolesByQuery(query)
	if err != nil {
		return err
	}
	lo.G.Debugf("Found %d roles from %s API", len(roles), "space")
	err = m.checkResultsAllReturned("space", roles)
	if err != nil {
		return err
	}
	m.dumpV3Roles("space", roles)
	for _, role := range roles {
		spaceGUID := role.Relationships["space"].Data.GUID
		user, err := m.getUserForGUID(role.Relationships["user"].Data.GUID)
		if err != nil {
			return err
		}
		spaceRoleMap, ok := spaceV3UsersRolesMap[spaceGUID]
		if !ok {
			spaceRoleMap = make(map[string][]cfclient.V3User)
			spaceV3UsersRolesMap[spaceGUID] = spaceRoleMap
		}
		spaceRoleMap[role.Type] = append(spaceRoleMap[role.Type], *user)
	}
	spaceUsersRoleMap := make(map[string]map[string]*RoleUsers)
	for key, val := range spaceV3UsersRolesMap {
		for role, users := range val {
			uaaUsers, err := m.GetUAAUsers()
			if err != nil {
				return err
			}
			roleUsers, err := NewRoleUsers(users, uaaUsers)
			if err != nil {
				return err
			}
			roleMap, ok := spaceUsersRoleMap[key]
			if !ok {
				roleMap = make(map[string]*RoleUsers)
				spaceUsersRoleMap[key] = roleMap
			}
			roleMap[role] = roleUsers
		}
	}
	m.dumpRolesUsers("spaces", spaceUsersRoleMap)
	m.SpaceRoles = spaceUsersRoleMap
	return nil
}

func (m *DefaultManager) checkResultsAllReturned(entityType string, roles []cfclient.V3Role) error {
	tracker := make(map[string]string)
	for _, role := range roles {
		if _, ok := tracker[role.GUID]; !ok {
			tracker[role.GUID] = role.GUID
		} else {
			return fmt.Errorf("role for type %s with GUID[%s] is returned multiple times, pagination for v3 roles is not working", entityType, role.GUID)
		}
	}
	return nil
}

func (m *DefaultManager) initializeOrgUserRolesMap() error {
	orgV3UsersRolesMap := make(map[string]map[string][]cfclient.V3User)
	query := url.Values{}
	query["per_page"] = []string{"5000"}
	query["order_by"] = []string{"created_at"}
	query["types"] = []string{ORG_AUDITOR + "," + ORG_BILLING_MANAGER + "," + ORG_MANAGER + "," + ORG_USER}
	roles, err := m.Client.ListV3RolesByQuery(query)
	if err != nil {
		return err
	}
	lo.G.Debugf("Found %d roles from %s API", len(roles), "organization")
	err = m.checkResultsAllReturned("organization", roles)
	if err != nil {
		return err
	}
	m.dumpV3Roles("organization", roles)
	for _, role := range roles {
		orgGUID := role.Relationships["organization"].Data.GUID
		user, err := m.getUserForGUID(role.Relationships["user"].Data.GUID)
		if err != nil {
			return err
		}
		orgRoleMap, ok := orgV3UsersRolesMap[orgGUID]
		if !ok {
			orgRoleMap = make(map[string][]cfclient.V3User)
			orgV3UsersRolesMap[orgGUID] = orgRoleMap
		}
		orgRoleMap[role.Type] = append(orgRoleMap[role.Type], *user)
	}
	orgUsersRoleMap := make(map[string]map[string]*RoleUsers)
	for key, val := range orgV3UsersRolesMap {
		for role, users := range val {
			uaaUsers, err := m.GetUAAUsers()
			if err != nil {
				return err
			}
			roleUsers, err := NewRoleUsers(users, uaaUsers)
			if err != nil {
				return err
			}
			roleMap, ok := orgUsersRoleMap[key]
			if !ok {
				roleMap = make(map[string]*RoleUsers)
				orgUsersRoleMap[key] = roleMap
			}
			roleMap[role] = roleUsers
		}
	}
	m.dumpRolesUsers("organizations", orgUsersRoleMap)
	m.OrgRoles = orgUsersRoleMap

	return nil
}

func (m *DefaultManager) RemoveSpaceAuditor(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Auditor")
	role, err := m.GetSpaceRoleGUID(input.SpaceGUID, userGUID, SPACE_AUDITOR)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}
func (m *DefaultManager) RemoveSpaceDeveloper(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Developer")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Developer")
	role, err := m.GetSpaceRoleGUID(input.SpaceGUID, userGUID, SPACE_DEVELOPER)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}
func (m *DefaultManager) RemoveSpaceManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Manager")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Manager")
	role, err := m.GetSpaceRoleGUID(input.SpaceGUID, userGUID, SPACE_MANAGER)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}

func (m *DefaultManager) RemoveSpaceSupporter(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "supporter")
		return nil
	}

	if !m.SupportsSpaceSupporter {
		lo.G.Infof("this instance of cloud foundry does not support space_supporter role")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "supporter")
	role, err := m.GetSpaceRoleGUID(input.SpaceGUID, userGUID, SPACE_SUPPORTER)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}

func (m *DefaultManager) AssociateSpaceAuditor(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "auditor", input.OrgName, input.SpaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "auditor", input.OrgName, input.SpaceName)
	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_AUDITOR)
	return err
}
func (m *DefaultManager) AssociateSpaceDeveloper(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "developer", input.OrgName, input.SpaceName)
		return nil
	}
	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "developer", input.OrgName, input.SpaceName)
	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_DEVELOPER)
	return err
}
func (m *DefaultManager) AssociateSpaceManager(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "manager", input.OrgName, input.SpaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "manager", input.OrgName, input.SpaceName)
	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_MANAGER)
	return err
}

func (m *DefaultManager) AssociateSpaceSupporter(input UsersInput, userName, userGUID string) error {

	if !m.SupportsSpaceSupporter {
		lo.G.Infof("this instance of cloud foundry does not support space_supporter role")
		return nil
	}
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "supporter", input.OrgName, input.SpaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "supporter", input.OrgName, input.SpaceName)
	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_SUPPORTER)
	return err
}

func (m *DefaultManager) AddUserToOrg(orgGUID string, userName, userGUID string) error {
	if m.Peek {
		return nil
	}
	orgUsers, _, _, _, err := m.ListOrgUsersByRole(orgGUID)
	if err != nil {
		return err
	}
	if !orgUsers.HasUserForGUID(userName, userGUID) {
		_, err := m.Client.CreateV3OrganizationRole(orgGUID, userGUID, ORG_USER)
		if err != nil {
			lo.G.Debugf("Error adding user [%s] to org with guid [%s] but should have succeeded missing from org roles %+v, message: [%s]", userName, userGUID, orgUsers, err.Error())
			return err
		}
		orgUsers.addUser(RoleUser{UserName: userName, GUID: userGUID})
		m.updateOrgRoleUsers(orgGUID, orgUsers)
	}
	return nil
}

func (m *DefaultManager) updateOrgRoleUsers(orgGUID string, roleUser *RoleUsers) {
	orgRoles, ok := m.OrgRoles[orgGUID]
	if !ok {
		orgRoles = make(map[string]*RoleUsers)
	}
	orgRoles[ORG_USER] = roleUser
	m.OrgRoles[orgGUID] = orgRoles
}

func (m *DefaultManager) RemoveOrgAuditor(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "auditor")
	role, err := m.GetOrgRoleGUID(input.OrgGUID, userGUID, ORG_AUDITOR)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}
func (m *DefaultManager) RemoveOrgBillingManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "billing manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "billing manager")
	role, err := m.GetOrgRoleGUID(input.OrgGUID, userGUID, ORG_BILLING_MANAGER)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}

func (m *DefaultManager) RemoveOrgManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "manager")
	role, err := m.GetOrgRoleGUID(input.OrgGUID, userGUID, ORG_MANAGER)
	if err != nil {
		return err
	}
	return m.Client.DeleteV3Role(role)
}

func (m *DefaultManager) AssociateOrgAuditor(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "auditor", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "auditor", input.OrgName)
	_, err = m.Client.CreateV3OrganizationRole(input.OrgGUID, userGUID, ORG_AUDITOR)
	return err
}
func (m *DefaultManager) AssociateOrgBillingManager(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "billing manager", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "billing manager", input.OrgName)
	_, err = m.Client.CreateV3OrganizationRole(input.OrgGUID, userGUID, ORG_BILLING_MANAGER)
	return err
}

func (m *DefaultManager) AssociateOrgManager(input UsersInput, userName, userGUID string) error {
	err := m.AddUserToOrg(input.OrgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "manager", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "manager", input.OrgName)
	_, err = m.Client.CreateV3OrganizationRole(input.OrgGUID, userGUID, ORG_MANAGER)
	return err
}

// UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers() []error {
	errs := []error{}
	m.ClearRoles()
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return []error{err}
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(&input); err != nil {
			errs = append(errs, err)
		}
	}
	m.LogResults()
	return errs
}

func (m *DefaultManager) updateSpaceUsers(input *config.SpaceConfig) error {
	space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error finding space for org %s, space %s", input.Org, input.Space))
	}
	lo.G.Debug("")
	lo.G.Debug("")
	lo.G.Debugf("Processing Org(%s)/Space(%s)", input.Org, input.Space)
	lo.G.Debug("")
	lo.G.Debug("")

	managers, developers, auditors, supporters, err := m.ListSpaceUsersByRole(space.Guid)
	if err != nil {
		return err
	}

	if err = m.SyncUsers(UsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.Guid,
		OrgName:        input.Org,
		OrgGUID:        space.OrganizationGuid,
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      developers,
		RemoveUser:     m.RemoveSpaceDeveloper,
		AddUser:        m.AssociateSpaceDeveloper,
		Role:           SPACE_DEVELOPER,
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "developer"))
	}

	if err = m.SyncUsers(
		UsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.Guid,
			OrgGUID:        space.OrganizationGuid,
			OrgName:        input.Org,
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      managers,
			RemoveUser:     m.RemoveSpaceManager,
			AddUser:        m.AssociateSpaceManager,
			Role:           SPACE_MANAGER,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "manager"))
	}
	if err = m.SyncUsers(
		UsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.Guid,
			OrgGUID:        space.OrganizationGuid,
			OrgName:        input.Org,
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      auditors,
			RemoveUser:     m.RemoveSpaceAuditor,
			AddUser:        m.AssociateSpaceAuditor,
			Role:           SPACE_AUDITOR,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "auditor"))
	}

	if err = m.SyncUsers(UsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.Guid,
		OrgName:        input.Org,
		OrgGUID:        space.OrganizationGuid,
		LdapGroupNames: input.GetSupporterGroups(),
		LdapUsers:      input.Supporter.LDAPUsers,
		Users:          input.Supporter.Users,
		SamlUsers:      input.Supporter.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      supporters,
		RemoveUser:     m.RemoveSpaceSupporter,
		AddUser:        m.AssociateSpaceSupporter,
		Role:           SPACE_SUPPORTER,
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "developer"))
	}

	lo.G.Debug("")
	lo.G.Debug("")
	lo.G.Debugf("Done Processing Org(%s)/Space(%s)", input.Org, input.Space)
	lo.G.Debug("")
	lo.G.Debug("")
	return nil
}

// UpdateOrgUsers -
func (m *DefaultManager) UpdateOrgUsers() []error {
	errs := []error{}
	m.ClearRoles()
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return []error{err}
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(&input); err != nil {
			errs = append(errs, err)
		}

	}
	m.LogResults()
	return errs
}

func (m *DefaultManager) updateOrgUsers(input *config.OrgConfig) error {
	org, err := m.OrgReader.FindOrg(input.Org)
	if err != nil {
		return err
	}

	_, managers, billingManagers, auditors, err := m.ListOrgUsersByRole(org.Guid)
	if err != nil {
		return err
	}
	err = m.SyncUsers(
		UsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetBillingManagerGroups(),
			LdapUsers:      input.BillingManager.LDAPUsers,
			Users:          input.BillingManager.Users,
			SamlUsers:      input.BillingManager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      billingManagers,
			RemoveUser:     m.RemoveOrgBillingManager,
			AddUser:        m.AssociateOrgBillingManager,
			Role:           ORG_BILLING_MANAGER,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "billing_managers"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:        org.Name,
		OrgGUID:        org.Guid,
		LdapGroupNames: input.GetAuditorGroups(),
		LdapUsers:      input.Auditor.LDAPUsers,
		Users:          input.Auditor.Users,
		SamlUsers:      input.Auditor.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      auditors,
		RemoveUser:     m.RemoveOrgAuditor,
		AddUser:        m.AssociateOrgAuditor,
		Role:           ORG_AUDITOR,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-auditors"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:        org.Name,
		OrgGUID:        org.Guid,
		LdapGroupNames: input.GetManagerGroups(),
		LdapUsers:      input.Manager.LDAPUsers,
		Users:          input.Manager.Users,
		SamlUsers:      input.Manager.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      managers,
		RemoveUser:     m.RemoveOrgManager,
		AddUser:        m.AssociateOrgManager,
		Role:           ORG_MANAGER,
	})

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-manager"))
	}

	return nil
}

func (m *DefaultManager) dumpRoleUsers(message string, users []RoleUser) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		lo.G.Debugf("Begin %s", message)
		for _, roleUser := range users {
			lo.G.Debugf("%+v", roleUser)
		}
		lo.G.Debugf("End %s", message)
	}
}

// SyncUsers
func (m *DefaultManager) SyncUsers(usersInput UsersInput) error {
	roleUsers := usersInput.RoleUsers
	m.dumpRoleUsers(fmt.Sprintf("Current Users for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())

	if err := m.SyncLdapUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding ldap users")
	}
	if len(roleUsers.Users()) > 0 {
		m.dumpRoleUsers(fmt.Sprintf("Users after LDAP sync for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
	}

	if err := m.SyncInternalUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding internal users")
	}
	if len(roleUsers.Users()) > 0 {
		m.dumpRoleUsers(fmt.Sprintf("Users after Internal sync for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
	}

	if err := m.SyncSamlUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding saml users")
	}
	if len(roleUsers.Users()) > 0 {
		m.dumpRoleUsers(fmt.Sprintf("Users after SAML sync for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
	}

	if err := m.RemoveUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "removing users")
	}
	return nil
}

func (m *DefaultManager) SyncInternalUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	origin := "uaa"
	uaaUsers, err := m.GetUAAUsers()
	if err != nil {
		return err
	}
	for _, userID := range usersInput.UniqueUsers() {
		lowerUserID := strings.ToLower(userID)
		uaaUserList := uaaUsers.GetByName(lowerUserID)
		if len(uaaUserList) == 0 || !strings.EqualFold(uaaUserList[0].Origin, origin) {
			return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
		}
		if !roleUsers.HasUser(lowerUserID) {
			user := uaaUsers.GetByNameAndOrigin(lowerUserID, origin)
			if user == nil {
				return fmt.Errorf("unable to find user %s for origin %s", lowerUserID, origin)
			}
			m.dumpRoleUsers(fmt.Sprintf("Adding user [%s] with guid[%s] with origin [%s] as doesn't exist in users for %s/%s - Role %s", lowerUserID, user.GUID, origin, usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
			if err := usersInput.AddUser(usersInput, user.Username, user.GUID); err != nil {
				return errors.Wrap(err, fmt.Sprintf("adding user %s for origin %s", user.Username, origin))
			}
		} else {
			roleUsers.RemoveUserForOrigin(lowerUserID, origin)
		}
	}
	return nil
}

func (m *DefaultManager) RemoveUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	if usersInput.RemoveUsers {
		cfg, err := m.Cfg.GetGlobalConfig()
		if err != nil {
			return err
		}
		protectedUsers := cfg.ProtectedUsers

		if len(roleUsers.Users()) > 0 {
			m.dumpRoleUsers(fmt.Sprintf("The following users are being removed for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
		}
		for _, roleUser := range roleUsers.Users() {
			if !util.Matches(roleUser.UserName, protectedUsers) {
				if err := usersInput.RemoveUser(usersInput, roleUser.UserName, roleUser.GUID); err != nil {
					return errors.Wrap(err, fmt.Sprintf("error removing user %s", roleUser.UserName))
				}
			}
		}
	} else {
		if usersInput.SpaceName == "" {
			lo.G.Debugf("Not removing users. Set enable-remove-users: true to orgConfig for org: %s", usersInput.OrgName)
		} else {
			lo.G.Debugf("Not removing users. Set enable-remove-users: true to spaceConfig for org/space: %s/%s", usersInput.OrgName, usersInput.SpaceName)
		}
	}
	return nil
}

func (m *DefaultManager) InitializeLdap(ldapBindUser, ldapBindPassword, ldapServer string) error {
	ldapConfig, err := m.Cfg.LdapConfig(ldapBindUser, ldapBindPassword, ldapServer)
	if err != nil {
		return err
	}
	m.LdapConfig = ldapConfig
	if m.LdapConfig.Enabled {
		ldapMgr, err := ldap.NewManager(ldapConfig)
		if err != nil {
			return err
		}
		m.LdapMgr = ldapMgr
	}
	return nil
}

func (m *DefaultManager) DeinitializeLdap() error {
	if m.LdapMgr != nil {
		m.LdapMgr.Close()
	}
	return nil
}

func (m *DefaultManager) GetOrgRoleGUID(orgGUID, userGUID, role string) (string, error) {
	roles, err := m.Client.ListV3RolesByQuery(url.Values{
		"organization_guids": []string{orgGUID},
		"user_guids":         []string{userGUID},
		"types":              []string{role},
	})
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", fmt.Errorf("no role found for orgGUID: %s, userGUID: %s and types: %s", orgGUID, userGUID, role)
	}
	if len(roles) > 1 {
		return "", fmt.Errorf("more than 1 role found for orgGUID: %s, userGUID: %s and types: %s", orgGUID, userGUID, role)
	}
	return roles[0].GUID, nil
}

func (m *DefaultManager) GetSpaceRoleGUID(spaceGUID, userGUID, role string) (string, error) {
	roles, err := m.Client.ListV3RolesByQuery(url.Values{
		"space_guids": []string{spaceGUID},
		"user_guids":  []string{userGUID},
		"types":       []string{role},
	})
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", fmt.Errorf("no role found for spaceGUID: %s, userGUID: %s and types: %s", spaceGUID, userGUID, role)
	}
	if len(roles) > 1 {
		return "", fmt.Errorf("more than 1 role found for spaceGUID: %s, userGUID: %s and types: %s", spaceGUID, userGUID, role)
	}
	return roles[0].GUID, nil
}
