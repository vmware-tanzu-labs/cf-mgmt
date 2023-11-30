package user

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/azureAD"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

const ORG_AUDITOR string = "organization_auditor"
const ORG_MANAGER string = "organization_manager"
const ORG_BILLING_MANAGER string = "organization_billing_manager"
const SPACE_AUDITOR string = "space_auditor"
const SPACE_DEVELOPER string = "space_developer"
const SPACE_MANAGER string = "space_manager"
const SPACE_SUPPORTER string = "space_supporter"

// NewManager -
func NewManager(
	cfg config.Reader,
	spaceMgr space.Manager,
	orgReader organizationreader.Reader,
	uaaMgr uaa.Manager, roleMgr role.Manager, ldapMgr *ldap.Manager,
	peek bool) (Manager, error) {

	ldapConfig, err := cfg.LdapConfig("", "", "")
	if err != nil {
		return nil, err
	}
	mgr := &DefaultManager{
		Peek:       peek,
		SpaceMgr:   spaceMgr,
		OrgReader:  orgReader,
		UAAMgr:     uaaMgr,
		RoleMgr:    roleMgr,
		LdapMgr:    ldapMgr,
		Cfg:        cfg,
		LdapConfig: ldapConfig,
	}
	return mgr, nil
}

type DefaultManager struct {
	Cfg           config.Reader
	SpaceMgr      space.Manager
	OrgReader     organizationreader.Reader
	UAAMgr        uaa.Manager
	RoleMgr       role.Manager
	Peek          bool
	LdapMgr       LdapManager
	LdapConfig    *config.LdapConfig
	AzureADMgr    AzureADManager
	AzureADConfig *config.AzureADConfig
	UAAUsers      *uaa.Users
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

// UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers() []error {
	errs := []error{}
	m.RoleMgr.ClearRoles()
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return []error{err}
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(&input); err != nil {
			errs = append(errs, err)
		}
	}
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

	managers, developers, auditors, supporters, err := m.RoleMgr.ListSpaceUsersByRole(space.GUID)
	if err != nil {
		return err
	}

	if err = m.SyncUsers(UsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.GUID,
		OrgName:        input.Org,
		OrgGUID:        space.Relationships.Organization.Data.GUID,
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SPNUsers:       input.Developer.SPNUsers,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      developers,
		RemoveUser:     m.RoleMgr.RemoveSpaceDeveloper,
		AddUser:        m.RoleMgr.AssociateSpaceDeveloper,
		Role:           SPACE_DEVELOPER,
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "developer"))
	}

	if err = m.SyncUsers(
		UsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.GUID,
			OrgGUID:        space.Relationships.Organization.Data.GUID,
			OrgName:        input.Org,
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SPNUsers:       input.Manager.SPNUsers,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      managers,
			RemoveUser:     m.RoleMgr.RemoveSpaceManager,
			AddUser:        m.RoleMgr.AssociateSpaceManager,
			Role:           SPACE_MANAGER,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "manager"))
	}
	if err = m.SyncUsers(
		UsersInput{
			SpaceName:      space.Name,
			SpaceGUID:      space.GUID,
			OrgGUID:        space.Relationships.Organization.Data.GUID,
			OrgName:        input.Org,
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			SPNUsers:       input.Auditor.SPNUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      auditors,
			RemoveUser:     m.RoleMgr.RemoveSpaceAuditor,
			AddUser:        m.RoleMgr.AssociateSpaceAuditor,
			Role:           SPACE_AUDITOR,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "auditor"))
	}

	if err = m.SyncUsers(UsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.GUID,
		OrgName:        input.Org,
		OrgGUID:        space.Relationships.Organization.Data.GUID,
		LdapGroupNames: input.GetSupporterGroups(),
		LdapUsers:      input.Supporter.LDAPUsers,
		Users:          input.Supporter.Users,
		SPNUsers:       input.Supporter.SPNUsers,
		SamlUsers:      input.Supporter.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      supporters,
		RemoveUser:     m.RoleMgr.RemoveSpaceSupporter,
		AddUser:        m.RoleMgr.AssociateSpaceSupporter,
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
	m.RoleMgr.ClearRoles()
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return []error{err}
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(&input); err != nil {
			errs = append(errs, err)
		}

	}
	return errs
}

func (m *DefaultManager) updateOrgUsers(input *config.OrgConfig) error {
	org, err := m.OrgReader.FindOrg(input.Org)
	if err != nil {
		return err
	}

	_, managers, billingManagers, auditors, err := m.RoleMgr.ListOrgUsersByRole(org.GUID)
	if err != nil {
		return err
	}
	err = m.SyncUsers(
		UsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.GUID,
			LdapGroupNames: input.GetBillingManagerGroups(),
			LdapUsers:      input.BillingManager.LDAPUsers,
			Users:          input.BillingManager.Users,
			SPNUsers:       input.BillingManager.SPNUsers,
			SamlUsers:      input.BillingManager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			RoleUsers:      billingManagers,
			RemoveUser:     m.RoleMgr.RemoveOrgBillingManager,
			AddUser:        m.RoleMgr.AssociateOrgBillingManager,
			Role:           ORG_BILLING_MANAGER,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "billing_managers"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:        org.Name,
		OrgGUID:        org.GUID,
		LdapGroupNames: input.GetAuditorGroups(),
		LdapUsers:      input.Auditor.LDAPUsers,
		Users:          input.Auditor.Users,
		SPNUsers:       input.Auditor.SPNUsers,
		SamlUsers:      input.Auditor.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      auditors,
		RemoveUser:     m.RoleMgr.RemoveOrgAuditor,
		AddUser:        m.RoleMgr.AssociateOrgAuditor,
		Role:           ORG_AUDITOR,
	})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-auditors"))
	}

	err = m.SyncUsers(UsersInput{
		OrgName:        org.Name,
		OrgGUID:        org.GUID,
		LdapGroupNames: input.GetManagerGroups(),
		LdapUsers:      input.Manager.LDAPUsers,
		Users:          input.Manager.Users,
		SPNUsers:       input.Manager.SPNUsers,
		SamlUsers:      input.Manager.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		RoleUsers:      managers,
		RemoveUser:     m.RoleMgr.RemoveOrgManager,
		AddUser:        m.RoleMgr.AssociateOrgManager,
		Role:           ORG_MANAGER,
	})

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-manager"))
	}

	return nil
}

func (m *DefaultManager) dumpRoleUsers(message string, users []role.RoleUser) {
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

	if err := m.SyncInternalUsers(roleUsers, usersInput, false); err != nil {
		return errors.Wrap(err, "adding internal users")
	}
	if len(roleUsers.Users()) > 0 {
		m.dumpRoleUsers(fmt.Sprintf("Users after Internal sync for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
	}

	if err := m.SyncInternalUsers(roleUsers, usersInput, true); err != nil {
		return errors.Wrap(err, "adding internal SPN users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after Internal SPN sync %+v", roleUsers.Users())
	}

	if err := m.SyncSamlUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding saml users")
	}
	if len(roleUsers.Users()) > 0 {
		m.dumpRoleUsers(fmt.Sprintf("Users after SAML sync for %s/%s - Role %s", usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
	}

	// Sync AAD users after SAML users. Do the uniqueness check in this function, so we don;t have to touch the SAML users function too much
	if err := m.SyncAzureADUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding Azure AD users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after AzureAD sync %+v", roleUsers.Users())
	}

	if err := m.RemoveUsers(roleUsers, usersInput); err != nil {
		return errors.Wrap(err, "removing users")
	}
	return nil
}

func (m *DefaultManager) SyncInternalUsers(roleUsers *role.RoleUsers, usersInput UsersInput, spnUsers bool) error {
	var userList []string
	var origin string
	if !spnUsers {
		origin = "uaa"
		userList = usersInput.UniqueUsers()
	} else {
		origin = m.AzureADConfig.SPNOrigin
		userList = usersInput.UniqueSPNUsers()
	}

	uaaUsers, err := m.GetUAAUsers()
	if err != nil {
		return err
	}

	for _, userID := range userList {
		lowerUserID := strings.ToLower(userID)
		uaaUser := uaaUsers.GetByNameAndOrigin(lowerUserID, origin)
		if uaaUser == nil {
			return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
		}
		if !roleUsers.HasUserForGUID(lowerUserID, uaaUser.GUID) {
			user := uaaUsers.GetByNameAndOrigin(lowerUserID, origin)
			if user == nil {
				return fmt.Errorf("unable to find user %s for origin %s", lowerUserID, origin)
			}
			m.dumpRoleUsers(fmt.Sprintf("Adding user [%s] with guid[%s] with origin [%s] as doesn't exist in users for %s/%s - Role %s", lowerUserID, user.GUID, origin, usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
			if err := usersInput.AddUser(usersInput.OrgGUID, usersInput.EntityName(), usersInput.EntityGUID(), user.Username, user.GUID); err != nil {
				return errors.Wrap(err, fmt.Sprintf("adding user %s for origin %s", user.Username, origin))
			}
		} else {
			roleUsers.RemoveUserForOrigin(lowerUserID, origin)
		}
	}
	return nil
}

//	func (m *DefaultManager) SyncInternalSPNUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
//		origin := m.AzureADConfig.SPNOrigin
//		for _, userID := range usersInput.UniqueSPNUsers() {
//			lowerUserID := strings.ToLower(userID)
//			uaaUserList := m.UAAUsers.GetByName(lowerUserID)
//			if len(uaaUserList) == 0 || !strings.EqualFold(uaaUserList[0].Origin, origin) {
//				return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
//			}
//			if !roleUsers.HasUser(lowerUserID) {
//				lo.G.Debugf("Role Users %+v", roleUsers.users)
//				user := m.UAAUsers.GetByNameAndOrigin(lowerUserID, origin)
//				if user == nil {
//					return fmt.Errorf("Unable to find user %s for origin %s", lowerUserID, origin)
//				}
//				if err := usersInput.AddUser(usersInput, user.Username, user.GUID); err != nil {
//					return errors.Wrap(err, fmt.Sprintf("adding user %s for origin %s", user.Username, origin))
//				}
//			} else {
//				roleUsers.RemoveUserForOrigin(lowerUserID, origin)
//			}
//		}
//		return nil
//	}
func (m *DefaultManager) RemoveUsers(roleUsers *role.RoleUsers, usersInput UsersInput) error {
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
				if err := usersInput.RemoveUser(usersInput.EntityName(), usersInput.EntityGUID(), roleUser.UserName, roleUser.GUID); err != nil {
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
		lo.G.Debugf("Azure AD is Enabled, with TenantId: %s, ClientID: %s, Origin: %s, SPNOrigin: %s", aadConfig.TenantID, aadConfig.ClientId, aadConfig.UserOrigin, aadConfig.SPNOrigin)

	}
	return nil
}
