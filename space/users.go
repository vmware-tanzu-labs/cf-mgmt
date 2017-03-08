package space

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

type UserMgr interface {
	UpdateSpaceUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error
}

//NewManager -
func NewUserManager(
	cloudController cloudcontroller.Manager,
	ldapMgr ldap.Manager,
	uaacMgr uaac.Manager) UserMgr {
	return &UserManager{
		cloudController: cloudController,
		LdapMgr:         ldapMgr,
		UAACMgr:         uaacMgr,
	}
}

type UserManager struct {
	cloudController cloudcontroller.Manager
	LdapMgr         ldap.Manager
	UAACMgr         uaac.Manager
}
type UpdateUsersInput struct {
	SpaceGUID        string
	OrgGUID          string
	Role             string
	LdapGroupName    string
	LdapUsers, Users []string
}

func (m *UserManager) UpdateSpaceUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error {
	spaceUsers, err := m.cloudController.GetSpaceUsers(updateUsersInput.SpaceGUID, updateUsersInput.Role)
	if err != nil {
		return err
	}
	if config.Enabled {
		var users []ldap.User
		users, err = m.getLdapUsers(config, updateUsersInput.LdapGroupName, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		for _, user := range users {
			if _, ok := spaceUsers[user.UserID]; !ok {
				err = m.updateLdapUser(config, updateUsersInput.SpaceGUID, updateUsersInput.OrgGUID, updateUsersInput.Role, uaacUsers, user)
				if err != nil {
					return err
				}
			} else {
				delete(spaceUsers, user.UserID)
			}
		}
	} else {
		lo.G.Info("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		if _, ok := spaceUsers[userID]; !ok {
			if err = m.addUserToOrgAndRole(userID, updateUsersInput.OrgGUID, updateUsersInput.SpaceGUID, updateUsersInput.Role); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(spaceUsers, userID)
		}
	}
	for spaceUser, spaceUserGUID := range spaceUsers {
		lo.G.Info(fmt.Sprintf("removing %s from %s", spaceUser, updateUsersInput.SpaceGUID))
		err = m.cloudController.RemoveSpaceUser(updateUsersInput.SpaceGUID, spaceUserGUID, "developers")
		if err != nil {
			lo.G.Error(fmt.Sprintf("Unable to remove user : %s from space %s role in spce : %s", spaceUser, updateUsersInput.Role, updateUsersInput.SpaceGUID))
			lo.G.Error(fmt.Errorf("Cloud controller API error : %s", err))
			return err
		}
	}
	return nil
}

func (m *UserManager) updateLdapUser(config *ldap.Config, spaceGUID, orgGUID string,
	role string, uaacUsers map[string]string,
	user ldap.User) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	}
	userID = strings.ToLower(userID)

	if _, userExists := uaacUsers[userID]; userExists {
		lo.G.Info("User : ", userID, "already exists in cloud foundry")
	} else {
		if userID != "" {
			lo.G.Info("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAACMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Error(err)
				return err
			}
			uaacUsers[userID] = userID
		}
	}
	if userID != "" {
		if err := m.addUserToOrgAndRole(userID, orgGUID, spaceGUID, role); err != nil {
			lo.G.Error(err)
			return err
		}
	}
	return nil
}

func (m *UserManager) getLdapUsers(config *ldap.Config, groupName string, userList []string) ([]ldap.User, error) {
	users := []ldap.User{}
	if groupName != "" {
		lo.G.Info("Finding LDAP user for group : ", groupName)
		if groupUsers, err := m.LdapMgr.GetUserIDs(config, groupName); err == nil {
			users = append(users, groupUsers...)
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	for _, user := range userList {
		if ldapUser, err := m.LdapMgr.GetUser(config, user); err == nil {
			if ldapUser != nil {
				users = append(users, *ldapUser)
			}
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	return users, nil
}

func (m *UserManager) addUserToOrgAndRole(userID, orgGUID, spaceGUID, role string) error {
	lo.G.Info(fmt.Sprintf("Adding user to org :  %s and space: %s ", orgGUID, spaceGUID))
	if err := m.cloudController.AddUserToOrg(userID, orgGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	lo.G.Info(fmt.Sprintf("Adding user to space: %s  with role: %s", spaceGUID, role))
	if err := m.cloudController.AddUserToSpaceRole(userID, role, spaceGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	return nil
}
