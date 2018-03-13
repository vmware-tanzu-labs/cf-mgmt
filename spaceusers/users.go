package spaceusers

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
	ldapMgr ldap.Manager,
	uaaMgr uaa.Manager,
	peek bool) Manager {
	return &DefaultManager{
		client:   client,
		Peek:     peek,
		SpaceMgr: spaceMgr,
		LdapMgr:  ldapMgr,
		UAAMgr:   uaaMgr,
		Cfg:      cfg,
	}
}

type DefaultManager struct {
	client   CFClient
	Cfg      config.Reader
	SpaceMgr space.Manager
	LdapMgr  ldap.Manager
	UAAMgr   uaa.Manager
	Peek     bool
}

func (m *DefaultManager) RemoveSpaceAuditorByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Auditor")
		return nil
	}
	return m.client.RemoveSpaceAuditorByUsername(spaceGUID, userName)
}
func (m *DefaultManager) RemoveSpaceDeveloperByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Developer")
		return nil
	}
	return m.client.RemoveSpaceDeveloperByUsername(spaceGUID, userName)
}
func (m *DefaultManager) RemoveSpaceManagerByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Manager")
		return nil
	}
	return m.client.RemoveSpaceManagerByUsername(spaceGUID, userName)
}
func (m *DefaultManager) ListSpaceAuditors(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceAuditors(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) ListSpaceDevelopers(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceDevelopers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) ListSpaceManagers(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceManagers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *DefaultManager) associateOrgUserByUsername(orgGUID, userName string) error {
	_, err := m.client.AssociateOrgUserByUsername(orgGUID, userName)
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
	_, err = m.client.AssociateSpaceAuditorByUsername(spaceGUID, userName)
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
	_, err = m.client.AssociateSpaceDeveloperByUsername(spaceGUID, userName)
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
	_, err = m.client.AssociateSpaceManagerByUsername(spaceGUID, userName)
	return err
}

//UpdateSpaceUsers -
func (m *DefaultManager) UpdateSpaceUsers(configDir, ldapBindPassword string) error {
	config, err := m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

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
		if err := m.updateSpaceUsers(config, &input, uaaUsers); err != nil {
			return err
		}
	}

	return nil
}

func (m *DefaultManager) updateSpaceUsers(config *ldap.Config, input *config.SpaceConfig, uaaUsers map[string]string) error {
	space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
	if err != nil {
		return err
	}
	if err = m.syncSpaceUsers(config, uaaUsers, UpdateUsersInput{
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

	if err = m.syncSpaceUsers(config, uaaUsers,
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
	if err = m.syncSpaceUsers(config, uaaUsers,
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

//UpdateSpaceUsers Update space users
func (m *DefaultManager) syncSpaceUsers(config *ldap.Config, uaaUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	spaceUsers, err := updateUsersInput.ListUsers(updateUsersInput.SpaceGUID)
	if err != nil {
		return err
	}

	lo.G.Debugf("SpaceUsers before: %v", spaceUsers)
	if config.Enabled {
		var ldapUsers []ldap.User
		ldapUsers, err = m.LdapMgr.GetLdapUsers(config, updateUsersInput.LdapGroupNames, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %v", ldapUsers)
		for _, user := range ldapUsers {
			err = m.updateLdapUser(config, updateUsersInput, uaaUsers, user, spaceUsers)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		if _, userExists := uaaUsers[lowerUserID]; !userExists {
			return fmt.Errorf("user %s doesn't exist in cloud foundry, so must add internal user first", lowerUserID)
		}
		if _, ok := spaceUsers[lowerUserID]; !ok {
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, lowerUserID)
		}
	}

	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaaUsers[lowerUserEmail]; !userExists {
			lo.G.Debug("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err = m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				return err
			} else {
				uaaUsers[userEmail] = userEmail
			}
		}
		if _, ok := spaceUsers[lowerUserEmail]; !ok {
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userEmail); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, lowerUserEmail)
		}
	}
	if updateUsersInput.RemoveUsers {
		lo.G.Debugf("Deleting users for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
		for spaceUser, _ := range spaceUsers {
			err = updateUsersInput.RemoveUser(updateUsersInput.SpaceGUID, spaceUser)
			if err != nil {
				lo.G.Errorf("Cloud controller API error: %s", err)
				return err
			}
		}
	} else {
		lo.G.Debugf("Not removing users. Set enable-remove-users: true to spaceConfig for org/space: %s/%s", updateUsersInput.OrgName, updateUsersInput.SpaceName)
	}

	lo.G.Debugf("SpaceUsers after: %v", spaceUsers)
	return nil
}

func (m *DefaultManager) updateLdapUser(config *ldap.Config, updateUsersInput UpdateUsersInput,
	uaaUsers map[string]string,
	user ldap.User, spaceUsers map[string]string) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	} else {
		if user.Email == "" {
			user.Email = fmt.Sprintf("%s@user.from.ldap.cf", userID)
		}
	}
	userID = strings.ToLower(userID)

	if _, ok := spaceUsers[userID]; !ok {
		lo.G.Debugf("User[%s] not found in: %v", userID, spaceUsers)
		if _, userExists := uaaUsers[userID]; !userExists {
			lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userID)
				return nil
			} else {
				uaaUsers[userID] = userID
			}
		}
		if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, userID); err != nil {
			return err
		}
	} else {
		delete(spaceUsers, userID)
	}
	return nil
}
