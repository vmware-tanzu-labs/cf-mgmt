package organization

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

// UserMgr -
type UserMgr interface {
	UpdateOrgUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error
}

// NewUserManager -
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

// UserManager -
type UserManager struct {
	cloudController cloudcontroller.Manager
	LdapMgr         ldap.Manager
	UAACMgr         uaac.Manager
}

// UpdateUsersInput -
type UpdateUsersInput struct {
	OrgGUID          string
	Role             string
	LdapGroupName    string
	LdapUsers, Users []string
}

//UpdateOrgUsers -
func (m *UserManager) UpdateOrgUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error {

	orgUsers, err := m.cloudController.GetCFUsers(updateUsersInput.OrgGUID, ORGS, updateUsersInput.Role)

	if err != nil {
		return err
	}
	if config.Enabled {
		var ldapUsers []ldap.User
		ldapUsers, err = m.getLdapUsers(config, updateUsersInput.LdapGroupName, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		for _, user := range ldapUsers {
			err = m.updateLdapUser(config, updateUsersInput.OrgGUID, updateUsersInput.Role, uaacUsers, user)
			if err != nil {
				return err
			}
			if _, ok := orgUsers[user.UserID]; !ok {
				delete(orgUsers, user.UserID)
			}
		}
	} else {
		lo.G.Info("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		if _, ok := orgUsers[userID]; !ok {
			if err = m.addUserToOrgAndRole(userID, updateUsersInput.OrgGUID, updateUsersInput.Role); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(orgUsers, userID)
		}
	}
	for orgUser, orgUserGUID := range orgUsers {
		lo.G.Info(fmt.Sprintf("removing %s from %s", orgUser, updateUsersInput.OrgGUID))
		err = m.cloudController.RemoveCFUser(updateUsersInput.OrgGUID, ORGS, orgUserGUID, updateUsersInput.Role)
		if err != nil {
			lo.G.Error(fmt.Sprintf("Unable to remove user : %s from org %s with role %s", orgUser, updateUsersInput.OrgGUID, updateUsersInput.Role))
			lo.G.Error(fmt.Errorf("Cloud controller API error : %s", err))
			return err
		}
	}
	return nil
}

func (m *UserManager) updateLdapUser(config *ldap.Config, orgGUID string,
	role string, uaacUsers map[string]string,
	user ldap.User) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	}
	userID = strings.ToLower(userID)

	if _, userExists := uaacUsers[userID]; !userExists {
		lo.G.Info("User", userID, "doesn't exist in cloud foundry, so creating user")
		if err := m.UAACMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
			lo.G.Error(err)
			return err
		}
		uaacUsers[userID] = userID
		if err := m.addUserToOrgAndRole(userID, orgGUID, role); err != nil {
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

func (m *UserManager) addUserToOrgAndRole(userID, orgGUID, role string) error {
	lo.G.Info(fmt.Sprintf("Adding user to org :  %s ", orgGUID))
	if err := m.cloudController.AddUserToOrg(userID, orgGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	return nil
}
