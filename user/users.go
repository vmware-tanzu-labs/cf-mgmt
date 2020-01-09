package user

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/util"
	"github.com/pkg/errors"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(
	client CFClient,
	cfg config.Reader,
	spaceMgr space.Manager,
	orgMgr organization.Manager,
	uaaMgr uaa.Manager,
	peek bool) Manager {
	return &DefaultManager{
		Client:   client,
		Peek:     peek,
		SpaceMgr: spaceMgr,
		OrgMgr:   orgMgr,
		UAAMgr:   uaaMgr,
		Cfg:      cfg,
	}
}

type DefaultManager struct {
	Client     CFClient
	Cfg        config.Reader
	SpaceMgr   space.Manager
	OrgMgr     organization.Manager
	UAAMgr     uaa.Manager
	Peek       bool
	LdapMgr    LdapManager
	LdapConfig *config.LdapConfig
}

func (m *DefaultManager) RemoveSpaceAuditor(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Auditor")
	return m.Client.RemoveSpaceAuditor(input.SpaceGUID, userGUID)
}
func (m *DefaultManager) RemoveSpaceDeveloper(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Developer")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Developer")
	return m.Client.RemoveSpaceDeveloper(input.SpaceGUID, userGUID)
}
func (m *DefaultManager) RemoveSpaceManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Manager")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s/%s with role %s", userName, input.OrgName, input.SpaceName, "Manager")
	return m.Client.RemoveSpaceManager(input.SpaceGUID, userGUID)
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
	_, err = m.Client.AssociateSpaceAuditor(input.SpaceGUID, userGUID)
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
	_, err = m.Client.AssociateSpaceDeveloper(input.SpaceGUID, userGUID)
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
	_, err = m.Client.AssociateSpaceManager(input.SpaceGUID, userGUID)
	return err
}

func (m *DefaultManager) AddUserToOrg(orgGUID string, userName, userGUID string) error {
	if m.Peek {
		return nil
	}
	_, err := m.Client.AssociateOrgUser(orgGUID, userGUID)
	return err
}

func (m *DefaultManager) RemoveOrgAuditor(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "auditor")
	return m.Client.RemoveOrgAuditor(input.OrgGUID, userGUID)
}
func (m *DefaultManager) RemoveOrgBillingManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "billing manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "billing manager")
	return m.Client.RemoveOrgBillingManager(input.OrgGUID, userGUID)
}

func (m *DefaultManager) RemoveOrgManager(input UsersInput, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, input.OrgName, "manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, input.OrgName, "manager")
	return m.Client.RemoveOrgManager(input.OrgGUID, userGUID)
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
	_, err = m.Client.AssociateOrgAuditor(input.OrgGUID, userGUID)
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
	_, err = m.Client.AssociateOrgBillingManager(input.OrgGUID, userGUID)
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
	_, err = m.Client.AssociateOrgManager(input.OrgGUID, userGUID)
	return err
}

//UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers() error {
	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		return err
	}

	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(&input, uaaUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultManager) updateSpaceUsers(input *config.SpaceConfig, uaaUsers *uaa.Users) error {
	space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error finding space for org %s, space %s", input.Org, input.Space))
	}
	lo.G.Debug("")
	lo.G.Debug("")
	lo.G.Debugf("Processing Org(%s)/Space(%s)", input.Org, input.Space)
	lo.G.Debug("")
	lo.G.Debug("")
	if err = m.SyncUsers(uaaUsers, UsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.Guid,
		OrgName:        input.Org,
		OrgGUID:        space.OrganizationGuid,
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		ListUsers:      m.listSpaceDevelopers,
		RemoveUser:     m.RemoveSpaceDeveloper,
		AddUser:        m.AssociateSpaceDeveloper,
	}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "developer"))
	}

	if err = m.SyncUsers(uaaUsers,
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
			ListUsers:      m.listSpaceManagers,
			RemoveUser:     m.RemoveSpaceManager,
			AddUser:        m.AssociateSpaceManager,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "manager"))
	}
	if err = m.SyncUsers(uaaUsers,
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
			ListUsers:      m.listSpaceAuditors,
			RemoveUser:     m.RemoveSpaceAuditor,
			AddUser:        m.AssociateSpaceAuditor,
		}); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s, space %s, role %s", input.Org, input.Space, "auditor"))
	}

	lo.G.Debug("")
	lo.G.Debug("")
	lo.G.Debugf("Done Processing Org(%s)/Space(%s)", input.Org, input.Space)
	lo.G.Debug("")
	lo.G.Debug("")
	return nil
}

//UpdateOrgUsers -
func (m *DefaultManager) UpdateOrgUsers() error {
	uaacUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		return err
	}

	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(&input, uaacUsers); err != nil {
			return err
		}

	}

	return nil
}

func (m *DefaultManager) updateOrgUsers(input *config.OrgConfig, uaaUsers *uaa.Users) error {
	org, err := m.OrgMgr.FindOrg(input.Org)
	if err != nil {
		return err
	}

	err = m.SyncUsers(
		uaaUsers, UsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetBillingManagerGroups(),
			LdapUsers:      input.BillingManager.LDAPUsers,
			Users:          input.BillingManager.Users,
			SamlUsers:      input.BillingManager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			ListUsers:      m.listOrgBillingManagers,
			RemoveUser:     m.RemoveOrgBillingManager,
			AddUser:        m.AssociateOrgBillingManager,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "billing_managers"))
	}

	err = m.SyncUsers(
		uaaUsers, UsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			ListUsers:      m.listOrgAuditors,
			RemoveUser:     m.RemoveOrgAuditor,
			AddUser:        m.AssociateOrgAuditor,
		})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-auditors"))
	}

	err = m.SyncUsers(
		uaaUsers, UsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
			ListUsers:      m.listOrgManagers,
			RemoveUser:     m.RemoveOrgManager,
			AddUser:        m.AssociateOrgManager,
		})

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error syncing users for org %s role %s", input.Org, "org-manager"))
	}

	return nil
}

//SyncUsers
func (m *DefaultManager) SyncUsers(uaaUsers *uaa.Users, usersInput UsersInput) error {
	roleUsers, err := usersInput.ListUsers(usersInput, uaaUsers)
	if err != nil {
		return err
	}
	lo.G.Debugf("Current Users In Role %+v", roleUsers.Users())

	if err := m.SyncLdapUsers(roleUsers, uaaUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding ldap users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after LDAP sync %+v", roleUsers.Users())
	}

	if err := m.SyncInternalUsers(roleUsers, uaaUsers, usersInput); err != nil {
		return errors.Wrap(err, "adding internal users")
	}
	if len(roleUsers.Users()) > 0 {
		lo.G.Debugf("Users after Internal sync %+v", roleUsers.Users())
	}

	if err := m.SyncSamlUsers(roleUsers, uaaUsers, usersInput); err != nil {
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

func (m *DefaultManager) SyncInternalUsers(roleUsers *RoleUsers, uaaUsers *uaa.Users, usersInput UsersInput) error {
	origin := "uaa"
	for _, userID := range usersInput.UniqueUsers() {
		lowerUserID := strings.ToLower(userID)
		uaaUserList := uaaUsers.GetByName(lowerUserID)
		if len(uaaUserList) == 0 || !strings.EqualFold(uaaUserList[0].Origin, origin) {
			return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
		}
		if !roleUsers.HasUser(lowerUserID) {
			lo.G.Debugf("Role Users %+v", roleUsers.users)
			user := uaaUsers.GetByNameAndOrigin(lowerUserID, origin)
			if user == nil {
				return fmt.Errorf("Unabled to find user %s for origin %s", lowerUserID, origin)
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
