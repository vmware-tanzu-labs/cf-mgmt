package user

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/azureAD"
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
	uaaUsers, err := uaaMgr.ListUsers()
	if err != nil {
		return nil, err
	}
	cfUserMap := make(map[string]cfclient.V3User)
	cfUsers, err := client.ListV3UsersByQuery(url.Values{})
	if err != nil {
		return nil, err
	}
	for _, cfUser := range cfUsers {
		cfUserMap[cfUser.GUID] = cfUser
	}
	return &DefaultManager{
		Client:                 client,
		Peek:                   peek,
		SpaceMgr:               spaceMgr,
		OrgReader:              orgReader,
		UAAMgr:                 uaaMgr,
		Cfg:                    cfg,
		SupportsSpaceSupporter: supports,
		UAAUsers:               uaaUsers,
		CFUsers:                cfUserMap,
	}, nil
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
	AzureADMgr             AzureADManager
	AzureADConfig          *config.AzureADConfig
	SupportsSpaceSupporter bool
	UAAUsers               *uaa.Users
	CFUsers                map[string]cfclient.V3User
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
	r, _ := regexp.Compile("User '.+' already has '.+_auditor' role in space")
	if r.MatchString(err.Error()) {
		lo.G.Debug("User already exists with correct Role")
		return nil
	}
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
	r, _ := regexp.Compile("User '.+' already has '.+_developer' role in space")
	if r.MatchString(err.Error()) {
		lo.G.Debug("User already exists with correct Role")
		return nil
	}
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
	r, _ := regexp.Compile("User '.+' already has '.+_manager' role in space")
	if r.MatchString(err.Error()) {
		lo.G.Debug("User already exists with correct Role")
		return nil
	}
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
	r, _ := regexp.Compile("User '.+' already has '.+_supporter' role in space")
	if r.MatchString(err.Error()) {
		lo.G.Debug("User already exists with correct Role")
		return nil
	}
	return err
}

func (m *DefaultManager) AddUserToOrg(orgGUID string, userName, userGUID string) error {
	if m.Peek {
		return nil
	}
	orgUsers, err := m.ListOrgUsers(orgGUID)
	if err != nil {
		return err
	}
	if !orgUsers.HasUserForGUID(userName, userGUID) {
		_, err := m.Client.CreateV3OrganizationRole(orgGUID, userGUID, ORG_USER)
		if err != nil {
			return err
		}
		return err
	}
	return nil
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
func (m *DefaultManager) UpdateSpaceUsers() error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(&input); err != nil {
			return err
		}
	}

	return nil
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
		SpaceName:   space.Name,
		SpaceGUID:   space.Guid,
		OrgName:     input.Org,
		OrgGUID:     space.OrganizationGuid,
		GroupNames:  input.GetDeveloperGroups(),
		LdapUsers:   input.Developer.LDAPUsers,
		Users:       input.Developer.Users,
		SamlUsers:   input.Developer.SamlUsers,
		RemoveUsers: input.RemoveUsers,
		RoleUsers:   developers,
		RemoveUser:  m.RemoveSpaceDeveloper,
		AddUser:     m.AssociateSpaceDeveloper,
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "developer"))
	}

	if err = m.SyncUsers(
		UsersInput{
			SpaceName:   space.Name,
			SpaceGUID:   space.Guid,
			OrgGUID:     space.OrganizationGuid,
			OrgName:     input.Org,
			GroupNames:  input.GetManagerGroups(),
			LdapUsers:   input.Manager.LDAPUsers,
			Users:       input.Manager.Users,
			SamlUsers:   input.Manager.SamlUsers,
			RemoveUsers: input.RemoveUsers,
			RoleUsers:   managers,
			RemoveUser:  m.RemoveSpaceManager,
			AddUser:     m.AssociateSpaceManager,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "manager"))
	}
	if err = m.SyncUsers(
		UsersInput{
			SpaceName:   space.Name,
			SpaceGUID:   space.Guid,
			OrgGUID:     space.OrganizationGuid,
			OrgName:     input.Org,
			GroupNames:  input.GetAuditorGroups(),
			LdapUsers:   input.Auditor.LDAPUsers,
			Users:       input.Auditor.Users,
			SamlUsers:   input.Auditor.SamlUsers,
			RemoveUsers: input.RemoveUsers,
			RoleUsers:   auditors,
			RemoveUser:  m.RemoveSpaceAuditor,
			AddUser:     m.AssociateSpaceAuditor,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "auditor"))
	}

	if err = m.SyncUsers(UsersInput{
		SpaceName:   space.Name,
		SpaceGUID:   space.Guid,
		OrgName:     input.Org,
		OrgGUID:     space.OrganizationGuid,
		GroupNames:  input.GetSupporterGroups(),
		LdapUsers:   input.Supporter.LDAPUsers,
		Users:       input.Supporter.Users,
		SamlUsers:   input.Supporter.SamlUsers,
		RemoveUsers: input.RemoveUsers,
		RoleUsers:   supporters,
		RemoveUser:  m.RemoveSpaceSupporter,
		AddUser:     m.AssociateSpaceSupporter,
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
func (m *DefaultManager) UpdateOrgUsers() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(&input); err != nil {
			return err
		}

	}

	return nil
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
			OrgName:     org.Name,
			OrgGUID:     org.Guid,
			GroupNames:  input.GetBillingManagerGroups(),
			LdapUsers:   input.BillingManager.LDAPUsers,
			Users:       input.BillingManager.Users,
			SamlUsers:   input.BillingManager.SamlUsers,
			RemoveUsers: input.RemoveUsers,
			RoleUsers:   billingManagers,
			RemoveUser:  m.RemoveOrgBillingManager,
			AddUser:     m.AssociateOrgBillingManager,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "billing_managers"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:     org.Name,
		OrgGUID:     org.Guid,
		GroupNames:  input.GetAuditorGroups(),
		LdapUsers:   input.Auditor.LDAPUsers,
		Users:       input.Auditor.Users,
		SamlUsers:   input.Auditor.SamlUsers,
		RemoveUsers: input.RemoveUsers,
		RoleUsers:   auditors,
		RemoveUser:  m.RemoveOrgAuditor,
		AddUser:     m.AssociateOrgAuditor,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-auditors"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:     org.Name,
		OrgGUID:     org.Guid,
		GroupNames:  input.GetManagerGroups(),
		LdapUsers:   input.Manager.LDAPUsers,
		Users:       input.Manager.Users,
		SamlUsers:   input.Manager.SamlUsers,
		RemoveUsers: input.RemoveUsers,
		RoleUsers:   managers,
		RemoveUser:  m.RemoveOrgManager,
		AddUser:     m.AssociateOrgManager,
	})

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-manager"))
	}

	return nil
}

// SyncUsers
func (m *DefaultManager) SyncUsers(usersInput UsersInput) error {
	// roleUsers, err := usersInput.ListUsers(usersInput, uaaUsers)
	// if err != nil {
	// 	return err
	// }
	roleUsers := usersInput.RoleUsers
	lo.G.Debugf("Current Users In Role %+v", roleUsers.Users())

	if err := m.SyncLdapUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding ldap users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after LDAP sync %+v", roleUsers.Users())
	}

	if err := m.SyncAzureADUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding Azure AD users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after AzureAD sync %+v", roleUsers.Users())
	}

	if err := m.SyncInternalUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding internal users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after Internal sync %+v", roleUsers.Users())
	}

	if err := m.SyncSamlUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding saml users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after SAML sync %+v", roleUsers.Users())
	}

	if err := m.RemoveUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "removing users")
	}
	return nil
}

func (m *DefaultManager) SyncInternalUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	origin := "uaa"
	for _, userID := range usersInput.UniqueUsers() {
		lowerUserID := strings.ToLower(userID)
		uaaUserList := m.UAAUsers.GetByName(lowerUserID)
		if len(uaaUserList) == 0 || !strings.EqualFold(uaaUserList[0].Origin, origin) {
			return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
		}
		if !roleUsers.HasUser(lowerUserID) {
			lo.G.Debugf("Role Users %+v", roleUsers.users)
			user := m.UAAUsers.GetByNameAndOrigin(lowerUserID, origin)
			if user == nil {
				return fmt.Errorf("Unable to find user %s for origin %s", lowerUserID, origin)
			}
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
			lo.G.Debugf("The following users are being removed %+v", roleUsers.Users())
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
		tmpAadConfig, err := m.Cfg.AzureADConfig("", "", "", "") // Just to see if both ldap and AAD are configured, which is not supported (yet)
		if tmpAadConfig.Enabled {
			return errors.New("Both LDAP and Azure AD groups are enabled. This is not supported yet")
		}
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

func (m *DefaultManager) InitializeAzureAD(tenantId, clientId, secret, origin string) error {
	aadConfig, err := m.Cfg.AzureADConfig(tenantId, clientId, secret, origin)
	if err != nil {
		return err
	}
	m.AzureADConfig = aadConfig
	if m.AzureADConfig.Enabled {
		azureAdMgr, err := azureAD.NewManager(aadConfig)
		if err != nil {
			return err
		}
		m.AzureADMgr = azureAdMgr
		lo.G.Debugf("Azure AD is Enabled, with TenantId: %s, ClientID: %s, Origin: %s", aadConfig.TenantID, aadConfig.ClientId, aadConfig.UserOrigin)

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
