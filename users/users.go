package users

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(
	client CFClient,
	cfg config.Reader,
	spaceMgr space.Manager,
	uaaMgr uaa.Manager,
	peek bool) Manager {
	return &DefaultManager{
		Client:   client,
		Peek:     peek,
		SpaceMgr: spaceMgr,
		UAAMgr:   uaaMgr,
		Cfg:      cfg,
	}
}

type DefaultManager struct {
	Client   CFClient
	Cfg      config.Reader
	SpaceMgr space.Manager
	LdapMgr  ldap.Manager
	UAAMgr   uaa.Manager
	Peek     bool
}

func (m *DefaultManager) InitializeLdap(ldapBindPassword string) error {
	ldapMgr, err := ldap.NewManager(m.Cfg, ldapBindPassword)
	if err != nil {
		return err
	}
	m.LdapMgr = ldapMgr
	return nil
}

func (m *DefaultManager) RemoveSpaceAuditorByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Auditor")
		return nil
	}
	return m.Client.RemoveSpaceAuditorByUsername(spaceGUID, userName)
}
func (m *DefaultManager) RemoveSpaceDeveloperByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Developer")
		return nil
	}
	return m.Client.RemoveSpaceDeveloperByUsername(spaceGUID, userName)
}
func (m *DefaultManager) RemoveSpaceManagerByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Manager")
		return nil
	}
	return m.Client.RemoveSpaceManagerByUsername(spaceGUID, userName)
}
func (m *DefaultManager) ListSpaceAuditors(spaceGUID string) (map[string]string, error) {
	users, err := m.Client.ListSpaceAuditors(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) ListSpaceDevelopers(spaceGUID string) (map[string]string, error) {
	users, err := m.Client.ListSpaceDevelopers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) ListSpaceManagers(spaceGUID string) (map[string]string, error) {
	users, err := m.Client.ListSpaceManagers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) associateOrgUserByUsername(orgGUID, userName string) error {
	_, err := m.Client.AssociateOrgUserByUsername(orgGUID, userName)
	return err
}

func (m *DefaultManager) userListToMap(users []cfclient.User) map[string]string {
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[strings.ToLower(user.Username)] = user.Guid
	}
	return userMap
}

func (m *DefaultManager) AssociateSpaceAuditorByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "auditor", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.Client.AssociateSpaceAuditorByUsername(spaceGUID, userName)
	return err
}
func (m *DefaultManager) AssociateSpaceDeveloperByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "developer", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.Client.AssociateSpaceDeveloperByUsername(spaceGUID, userName)
	return err
}
func (m *DefaultManager) AssociateSpaceManagerByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "manager", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.Client.AssociateSpaceManagerByUsername(spaceGUID, userName)
	return err
}

//UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers() error {
	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range spaceConfigs {
		if err := m.updateSpaceUsers(&input, uaaUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultManager) updateSpaceUsers(input *config.SpaceConfig, uaaUsers map[string]string) error {
	space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
	if err != nil {
		return err
	}
	if err = m.SyncSpaceUsers(uaaUsers, UpdateUsersInput{
		SpaceName:      space.Name,
		SpaceGUID:      space.Guid,
		OrgName:        input.Org,
		OrgGUID:        space.OrganizationGuid,
		LdapGroupNames: input.GetDeveloperGroups(),
		LdapUsers:      input.Developer.LDAPUsers,
		Users:          input.Developer.Users,
		SamlUsers:      input.Developer.SamlUsers,
		RemoveUsers:    input.RemoveUsers,
		ListUsers:      m.ListSpaceDevelopers,
		RemoveUser:     m.RemoveSpaceDeveloperByUsername,
		AddUser:        m.AssociateSpaceDeveloperByUsername,
	}); err != nil {
		return err
	}

	if err = m.SyncSpaceUsers(uaaUsers,
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
			ListUsers:      m.ListSpaceManagers,
			RemoveUser:     m.RemoveSpaceManagerByUsername,
			AddUser:        m.AssociateSpaceManagerByUsername,
		}); err != nil {
		return err
	}
	if err = m.SyncSpaceUsers(uaaUsers,
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
			ListUsers:      m.ListSpaceAuditors,
			RemoveUser:     m.RemoveSpaceAuditorByUsername,
			AddUser:        m.AssociateSpaceAuditorByUsername,
		}); err != nil {
		return err
	}
	return nil
}

//SyncSpaceUsers
func (m *DefaultManager) SyncSpaceUsers(uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	roleUsers, err := updateUsersInput.ListUsers(updateUsersInput.SpaceGUID)
	if err != nil {
		return err
	}

	lo.G.Debugf("RoleUsers before: %v", roleUsers)
	if err := m.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return err
	}
	if err := m.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return err
	}
	if err := m.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput); err != nil {
		return err
	}
	if err := m.RemoveUsers(roleUsers, updateUsersInput); err != nil {
		return err
	}
	lo.G.Debugf("RoleUsers after: %v", roleUsers)
	return nil
}

func (m *DefaultManager) SyncLdapUsers(roleUsers, uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	config := m.LdapMgr.LdapConfig()
	if config.Enabled {
		ldapUsers, err := m.LdapMgr.GetLdapUsers(updateUsersInput.LdapGroupNames, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %v", ldapUsers)
		for _, inputUser := range ldapUsers {
			userToUse := m.UpdateUserInfo(inputUser)
			config := m.LdapMgr.LdapConfig()
			userID := userToUse.UserID
			if _, ok := roleUsers[userID]; !ok {
				lo.G.Debugf("User[%s] not found in: %v", userID, roleUsers)
				if _, userExists := uaaUsers[userID]; !userExists {
					lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
					if err := m.UAAMgr.CreateExternalUser(userID, userToUse.Email, userToUse.UserDN, config.Origin); err != nil {
						lo.G.Error("Unable to create user", userID)
						continue
					} else {
						uaaUsers[userID] = userID
					}
				}
				if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
					return err
				}
			} else {
				delete(roleUsers, userID)
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	return nil
}

func (m *DefaultManager) SyncInternalUsers(roleUsers, uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		if _, userExists := uaaUsers[lowerUserID]; !userExists {
			return fmt.Errorf("user %s doesn't exist in cloud foundry, so must add internal user first", lowerUserID)
		}
		if _, ok := roleUsers[lowerUserID]; !ok {
			if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
				return err
			}
		} else {
			delete(roleUsers, lowerUserID)
		}
	}
	return nil
}

func (m *DefaultManager) SyncSamlUsers(roleUsers, uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	config := m.LdapMgr.LdapConfig()
	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaaUsers[lowerUserEmail]; !userExists {
			lo.G.Debug("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				continue
			} else {
				uaaUsers[userEmail] = userEmail
			}
		}
		if _, ok := roleUsers[lowerUserEmail]; !ok {
			if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userEmail); err != nil {
				return err
			}
		} else {
			delete(roleUsers, lowerUserEmail)
		}
	}
	return nil
}

func (m *DefaultManager) RemoveUsers(roleUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	if updateUsersInput.RemoveUsers {
		lo.G.Debugf("Deleting users for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
		for roleUser, _ := range roleUsers {
			if err := updateUsersInput.RemoveUser(updateUsersInput.SpaceGUID, roleUser); err != nil {
				return err
			}
		}
	} else {
		lo.G.Debugf("Not removing users. Set enable-remove-users: true to spaceConfig for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
	}
	return nil
}

func (m *DefaultManager) UpdateUserInfo(user ldap.User) ldap.User {
	config := m.LdapMgr.LdapConfig()
	userID := strings.ToLower(user.UserID)
	externalID := user.UserDN
	email := user.Email
	if config.Origin != "ldap" {
		userID = strings.ToLower(user.Email)
		externalID = user.Email
	} else {
		if email == "" {
			email = fmt.Sprintf("%s@user.from.ldap.cf", userID)
		}
	}

	return ldap.User{
		UserID: userID,
		UserDN: externalID,
		Email:  email,
	}
}
