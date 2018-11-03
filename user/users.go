package user

import (
	"fmt"
	"net/url"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
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
	LdapMgr    ldap.Manager
	LdapConfig *config.LdapConfig
}

func (m *DefaultManager) RemoveSpaceAuditor(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Auditor")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Auditor")
	return m.Client.RemoveSpaceAuditorByUsernameAndOrigin(input.SpaceGUID, userName, origin)
}
func (m *DefaultManager) RemoveSpaceDeveloper(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Developer")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Developer")
	return m.Client.RemoveSpaceDeveloperByUsernameAndOrigin(input.SpaceGUID, userName, origin)
}
func (m *DefaultManager) RemoveSpaceManager(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Manager")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org/space %s/%s with role %s", userName, origin, input.OrgName, input.SpaceName, "Manager")
	return m.Client.RemoveSpaceManagerByUsernameAndOrigin(input.SpaceGUID, userName, origin)
}

func (m *DefaultManager) AssociateSpaceAuditor(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "auditor", input.OrgName, input.SpaceName)
		return nil
	}

	lo.G.Infof("adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "auditor", input.OrgName, input.SpaceName)
	_, err = m.Client.AssociateSpaceAuditorByUsernameAndOrigin(input.SpaceGUID, userName, origin)
	return err
}
func (m *DefaultManager) AssociateSpaceDeveloper(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "developer", input.OrgName, input.SpaceName)
		return nil
	}
	lo.G.Infof("adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "developer", input.OrgName, input.SpaceName)
	_, err = m.Client.AssociateSpaceDeveloperByUsernameAndOrigin(input.SpaceGUID, userName, origin)
	return err
}
func (m *DefaultManager) AssociateSpaceManager(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "manager", input.OrgName, input.SpaceName)
		return nil
	}

	lo.G.Infof("adding %s with origin %s to role %s for org/space %s/%s", userName, origin, "manager", input.OrgName, input.SpaceName)
	_, err = m.Client.AssociateSpaceManagerByUsernameAndOrigin(input.SpaceGUID, userName, origin)
	return err
}

func (m *DefaultManager) AddUserToOrg(userName, origin string, input UpdateUsersInput) error {
	if m.Peek {
		return nil
	}
	_, err := m.Client.AssociateOrgUserByUsernameAndOrigin(input.OrgGUID, userName, origin)
	return err
}

func (m *DefaultManager) RemoveOrgAuditor(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "auditor")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "auditor")
	return m.Client.RemoveOrgAuditorByUsernameAndOrigin(input.OrgGUID, userName, origin)
}
func (m *DefaultManager) RemoveOrgBillingManager(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "billing manager")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "billing manager")
	return m.Client.RemoveOrgBillingManagerByUsernameAndOrigin(input.OrgGUID, userName, origin)
}

func (m *DefaultManager) RemoveOrgManager(input UpdateUsersInput, userName, origin string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "manager")
		return nil
	}
	lo.G.Infof("removing user %s with origin %s from org %s with role %s", userName, origin, input.OrgName, "manager")
	return m.Client.RemoveOrgManagerByUsernameAndOrigin(input.OrgGUID, userName, origin)
}

func (m *DefaultManager) AssociateOrgAuditor(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s with origin %s to role %s for org %s", userName, origin, "auditor", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s with origin %s to role %s for org %s", userName, origin, "auditor", input.OrgName)
	_, err = m.Client.AssociateOrgAuditorByUsernameAndOrigin(input.OrgGUID, userName, origin)
	return err
}
func (m *DefaultManager) AssociateOrgBillingManager(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s with origin %s to role %s for org %s", userName, origin, "billing manager", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s with origin %s to role %s for org %s", userName, origin, "billing manager", input.OrgName)
	_, err = m.Client.AssociateOrgBillingManagerByUsernameAndOrigin(input.OrgGUID, userName, origin)
	return err
}

func (m *DefaultManager) AssociateOrgManager(input UpdateUsersInput, userName, origin string) error {
	err := m.AddUserToOrg(userName, origin, input)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s with origin %s to role %s for org %s", userName, origin, "manager", input.OrgName)
		return nil
	}

	lo.G.Infof("Add User %s with origin %s to role %s for org %s", userName, origin, "manager", input.OrgName)
	_, err = m.Client.AssociateOrgManagerByUsernameAndOrigin(input.OrgGUID, userName, origin)
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

func (m *DefaultManager) updateSpaceUsers(input *config.SpaceConfig, uaaUsers map[string]uaa.User) error {
	space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error finding space for org %s, space %s", input.Org, input.Space))
	}

	if err = m.SyncUsers(uaaUsers, UpdateUsersInput{
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
		UpdateUsersInput{
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
		UpdateUsersInput{
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

//CleanupOrgUsers -
func (m *DefaultManager) CleanupOrgUsers() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}
	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		return err
	}

	for _, input := range orgConfigs {
		if err := m.cleanupOrgUsers(uaaUsers, &input); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) cleanupOrgUsers(uaaUsers map[string]uaa.User, input *config.OrgConfig) error {
	org, err := m.OrgMgr.FindOrg(input.Org)
	if err != nil {
		return err
	}
	orgUsers, err := m.Client.ListOrgUsers(org.Guid)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error listing org users for org %s", input.Org))
	}

	usersInRoles, err := m.usersInOrgRoles(org.Name, org.Guid, uaaUsers)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error usersInOrgRoles for org %s", input.Org))
	}

	lo.G.Debugf("Users In Roles %+v", usersInRoles)

	for _, orgUser := range orgUsers {
		if !usersInRoles.HasUser(strings.ToLower(orgUser.Username)) {
			uaaUser, ok := uaaUsers[orgUser.Guid]
			if !ok {
				return fmt.Errorf("Unable to find user with id %s and userName %s", orgUser.Guid, orgUser.Username)
			}
			if strings.EqualFold(uaaUser.Origin, "uaa") {
				lo.G.Infof("Skipping removal of user %s with origin %s from org %s", orgUser.Username, uaaUser.Origin, input.Org)
				continue
			}
			if m.Peek {
				lo.G.Infof("[dry-run]: Removing User %s with origin %s from org %s", orgUser.Username, uaaUser.Origin, input.Org)
				continue
			}
			lo.G.Infof("Removing User %s with origin %s from org %s", orgUser.Username, uaaUser.Origin, input.Org)
			err := m.Client.RemoveOrgUserByUsernameAndOrigin(org.Guid, orgUser.Username, uaaUser.Origin)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("Error removing user %s with origin %s from org %s", orgUser.Username, uaaUser.Origin, input.Org))
			}
		}
	}
	return nil
}

func (m *DefaultManager) listSpaces(orgGUID string) ([]cfclient.Space, error) {
	spaces, err := m.Client.ListSpacesByQuery(url.Values{
		"q": []string{fmt.Sprintf("%s:%s", "organization_guid", orgGUID)},
	})
	if err != nil {
		return nil, err
	}
	return spaces, err

}

func (m *DefaultManager) updateOrgUsers(input *config.OrgConfig, uaaUsers map[string]uaa.User) error {
	org, err := m.OrgMgr.FindOrg(input.Org)
	if err != nil {
		return err
	}

	err = m.SyncUsers(
		uaaUsers, UpdateUsersInput{
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
		uaaUsers, UpdateUsersInput{
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
		uaaUsers, UpdateUsersInput{
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
func (m *DefaultManager) SyncUsers(uaaUsers map[string]uaa.User, updateUsersInput UpdateUsersInput) error {
	roleUsers, err := updateUsersInput.ListUsers(updateUsersInput, uaaUsers)
	if err != nil {
		return err
	}

	if err := m.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return errors.Wrap(err, "adding ldap users")
	}
	if err := m.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return errors.Wrap(err, "adding internal users")
	}
	if err := m.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return errors.Wrap(err, "adding saml users")
	}
	if err := m.RemoveUsers(roleUsers, updateUsersInput); err != nil {
		return errors.Wrap(err, "removing users")
	}
	return nil
}

func (m *DefaultManager) SyncInternalUsers(roleUsers *RoleUsers, uaaUsers map[string]uaa.User, updateUsersInput UpdateUsersInput) error {
	origin := "uaa"
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		uaaUser, userExists := uaaUsers[lowerUserID]
		if !userExists || !strings.EqualFold(uaaUser.Origin, origin) {
			return fmt.Errorf("user %s doesn't exist in origin %s, so must add internal user first", lowerUserID, origin)
		}
		if !roleUsers.HasUserForOrigin(lowerUserID, origin) {
			if err := updateUsersInput.AddUser(updateUsersInput, userID, origin); err != nil {
				return errors.Wrap(err, fmt.Sprintf("adding user %s for origin %s", userID, origin))
			}
		} else {
			roleUsers.RemoveUserForOrigin(lowerUserID, origin)
		}
	}
	return nil
}

func (m *DefaultManager) SyncSamlUsers(roleUsers *RoleUsers, uaaUsers map[string]uaa.User, updateUsersInput UpdateUsersInput) error {
	origin := m.LdapConfig.Origin
	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaaUsers[lowerUserEmail]; !userExists {
			lo.G.Debug("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				continue
			} else {
				uaaUsers[userEmail] = uaa.User{
					Username:   userEmail,
					Email:      userEmail,
					ExternalID: userEmail,
					Origin:     origin,
				}
			}
		}
		if !roleUsers.HasUserForOrigin(lowerUserEmail, origin) {
			if err := updateUsersInput.AddUser(updateUsersInput, userEmail, origin); err != nil {
				return err
			}
		} else {
			roleUsers.RemoveUserForOrigin(lowerUserEmail, origin)
		}
	}
	return nil
}

func (m *DefaultManager) RemoveUsers(roleUsers *RoleUsers, updateUsersInput UpdateUsersInput) error {
	if updateUsersInput.RemoveUsers {
		for _, roleUser := range roleUsers.Users() {
			if err := updateUsersInput.RemoveUser(updateUsersInput, roleUser.UserName, roleUser.Origin); err != nil {
				return errors.Wrap(err, fmt.Sprintf("error removing user %s with origin %s", roleUser.UserName, roleUser.Origin))
			}
		}
	} else {
		if updateUsersInput.SpaceName == "" {
			lo.G.Debugf("Not removing users. Set enable-remove-users: true to orgConfig for org: %s", updateUsersInput.OrgName)
		} else {
			lo.G.Debugf("Not removing users. Set enable-remove-users: true to spaceConfig for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
		}
	}
	return nil
}

func (m *DefaultManager) InitializeLdap(ldapBindPassword string) error {
	ldapConfig, err := m.Cfg.LdapConfig(ldapBindPassword)
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
