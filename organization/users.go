package organization

import (
	"fmt"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaa"
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
	uaaMgr uaa.Manager) UserMgr {
	return &UserManager{
		cloudController: cloudController,
		LdapMgr:         ldapMgr,
		UAAMgr:          uaaMgr,
	}
}

// UserManager -
type UserManager struct {
	cloudController cloudcontroller.Manager
	LdapMgr         ldap.Manager
	UAAMgr          uaa.Manager
}

// UpdateUsersInput -
type UpdateUsersInput struct {
	OrgName                                     string
	OrgGUID                                     string
	Role                                        string
	LdapUsers, Users, LdapGroupNames, SamlUsers []string
	RemoveUsers                                 bool
}

//UpdateOrgUsers -
func (m *UserManager) UpdateOrgUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error {

	orgUsers, err := m.cloudController.GetCFUsers(updateUsersInput.OrgGUID, ORGS, updateUsersInput.Role)

	if err != nil {
		return err
	}
	if config.Enabled {
		var ldapUsers []ldap.User
		ldapUsers, err = m.getLdapUsers(config, updateUsersInput.LdapGroupNames, updateUsersInput.LdapUsers)
		if err != nil {
			return err
		}
		for _, user := range ldapUsers {
			err = m.updateLdapUser(config, updateUsersInput.OrgGUID, updateUsersInput.Role, updateUsersInput.OrgName, uaacUsers, user, orgUsers)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		if _, ok := orgUsers[lowerUserID]; !ok {
			if _, userExists := uaacUsers[lowerUserID]; !userExists {
				return fmt.Errorf("User %s doesn't exist in cloud foundry, so must add internal user first", userID)
			}
			if err = m.addUserToOrgAndRole(userID, updateUsersInput.OrgGUID, updateUsersInput.Role, updateUsersInput.OrgName); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(orgUsers, lowerUserID)
		}
	}

	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaacUsers[lowerUserEmail]; !userExists {
			lo.G.Info("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err = m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, config.Origin); err != nil {
				lo.G.Errorf("Unable to create user %s due to error %s", userEmail, err.Error())
				return err
			} else {
				uaacUsers[userEmail] = userEmail
			}
		}
		if _, ok := orgUsers[lowerUserEmail]; !ok {
			if err = m.addUserToOrgAndRole(userEmail, updateUsersInput.OrgGUID, updateUsersInput.Role, updateUsersInput.OrgName); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(orgUsers, lowerUserEmail)
		}
	}

	if updateUsersInput.RemoveUsers {
		lo.G.Debugf("Deleting users for org: %s", updateUsersInput.OrgName)
		for orgUser, orgUserGUID := range orgUsers {
			lo.G.Infof("removing user: %s from org: %s and role: %s", orgUser, updateUsersInput.OrgName, updateUsersInput.Role)
			err = m.cloudController.RemoveCFUser(updateUsersInput.OrgGUID, ORGS, orgUserGUID, updateUsersInput.Role)
			if err != nil {
				lo.G.Errorf("Unable to remove user : %s from org %s with role %s", orgUser, updateUsersInput.OrgGUID, updateUsersInput.Role)
				lo.G.Errorf("Cloud controller API error : %s", err)
				return err
			}
		}
	} else {
		lo.G.Debugf("Not removing users. Set enable-remove-users: true to orgConfig for org: %s", updateUsersInput.OrgName)
	}
	return nil
}

func (m *UserManager) updateLdapUser(config *ldap.Config, orgGUID string,
	role string, orgName string, uaacUsers map[string]string,
	user ldap.User, orgUsers map[string]string) error {

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

	if _, ok := orgUsers[userID]; !ok {
		if _, userExists := uaacUsers[userID]; !userExists {
			lo.G.Info("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Errorf("Unable to create user %s due to error %s", userID, err.Error())
			} else {
				uaacUsers[userID] = userID
				if err := m.addUserToOrgAndRole(userID, orgGUID, role, orgName); err != nil {
					lo.G.Error(err)
					return err
				}
			}
		} else {
			if err := m.addUserToOrgAndRole(userID, orgGUID, role, orgName); err != nil {
				lo.G.Error(err)
				return err
			}
		}
	} else {
		delete(orgUsers, userID)
	}
	return nil
}

func (m *UserManager) getLdapUsers(config *ldap.Config, groupNames []string, userList []string) ([]ldap.User, error) {
	users := []ldap.User{}
	for _, groupName := range groupNames {
		if groupName != "" {
			lo.G.Debug("Finding LDAP user for group:", groupName)
			if groupUsers, err := m.LdapMgr.GetUserIDs(config, groupName); err == nil {
				users = append(users, groupUsers...)
			} else {
				lo.G.Warning(err)
			}
		}
	}
	for _, user := range userList {
		if ldapUser, err := m.LdapMgr.GetUser(config, user); err == nil {
			if ldapUser != nil {
				users = append(users, *ldapUser)
			}
		} else {
			lo.G.Warning(err)
		}
	}
	return users, nil
}

func (m *UserManager) addUserToOrgAndRole(userID, orgGUID, role, orgName string) error {
	if err := m.cloudController.AddUserToOrg(userID, orgGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	lo.G.Infof("Adding user: %s to org: %s with role: %s", userID, orgName, role)
	if err := m.cloudController.AddUserToOrgRole(userID, role, orgGUID); err != nil {
		lo.G.Error(err)
		return err
	}
	return nil
}
