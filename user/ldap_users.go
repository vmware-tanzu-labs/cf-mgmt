package user

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncLdapUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	origin := m.LdapConfig.Origin
	if m.LdapConfig.Enabled {
		uaaUsers, err := m.GetUAAUsers()
		if err != nil {
			return err
		}
		ldapUsers, err := m.GetLDAPUsers(usersInput)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %+v", ldapUsers)
		for _, inputUser := range ldapUsers {
			userToUse := m.UpdateUserInfo(inputUser)
			userID := userToUse.UserID
			userList := uaaUsers.GetByName(userID)
			if len(userList) == 0 {
				lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
				if userGUID, err := m.UAAMgr.CreateExternalUser(userID, userToUse.Email, userToUse.UserDN, m.LdapConfig.Origin); err != nil {
					lo.G.Errorf("Unable to create user %s with error %s", userID, err.Error())
					continue
				} else {
					m.AddUAAUser(uaa.User{
						Username:   userID,
						ExternalID: userToUse.UserDN,
						Origin:     m.LdapConfig.Origin,
						Email:      userToUse.Email,
						GUID:       userGUID,
					})
				}
			}
			if !roleUsers.HasUserForOrigin(userID, origin) {
				lo.G.Debugf("Role Users %+v", roleUsers.users)
				user := uaaUsers.GetByNameAndOrigin(userID, origin)
				if user == nil {
					return fmt.Errorf("Unable to find user %s for origin %s", userID, origin)
				}
				if err := usersInput.AddUser(usersInput, user.Username, user.GUID); err != nil {
					return errors.Wrap(err, fmt.Sprintf("User %s with origin %s", user.Username, user.Origin))
				}
			} else {
				lo.G.Debugf("User[%s] found in role", userID)
				roleUsers.RemoveUserForOrigin(userID, origin)
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	return nil
}

func (m *DefaultManager) GetLDAPUsers(usersInput UsersInput) ([]ldap.User, error) {
	var ldapUsers []ldap.User
	uaaUsers, err := m.GetUAAUsers()
	if err != nil {
		return nil, err
	}
	for _, groupName := range usersInput.UniqueLdapGroupNames() {
		userDNList, err := m.LdapMgr.GetUserDNs(groupName)
		if err != nil {
			return nil, err
		}
		for _, userDN := range userDNList {
			lo.G.Debugf("Checking for userDN %s", userDN)
			uaaUser := uaaUsers.GetByExternalID(userDN)
			if uaaUser != nil {
				lo.G.Debugf("UserDN [%s] found in UAA as [%s], skipping ldap lookup", userDN, uaaUser.Username)
				ldapUsers = append(ldapUsers, ldap.User{
					UserID: uaaUser.Username,
					UserDN: userDN,
					Email:  uaaUser.Email,
				})
			} else {
				lo.G.Debugf("UserDN [%s] not found in UAA, executing ldap lookup", userDN)
				user, err := m.LdapMgr.GetUserByDN(userDN)
				if err != nil {
					return nil, err
				}
				if user != nil {
					ldapUsers = append(ldapUsers, *user)
				} else {
					lo.G.Infof("UserDN %s not found in ldap", userDN)
				}
			}
		}
	}
	for _, userID := range usersInput.LdapUsers {
		userList := uaaUsers.GetByName(userID)
		if len(userList) > 0 {
			lo.G.Debugf("UserID [%s] found in UAA, skipping ldap lookup", userID)
			for _, uaaUser := range userList {
				lo.G.Debugf("Checking if userID [%s] with origin [%s] and externalID [%s] matches ldap origin", uaaUser.Username, uaaUser.Origin, uaaUser.ExternalID)
				if strings.EqualFold(uaaUser.Origin, m.LdapConfig.Origin) {
					ldapUsers = append(ldapUsers, ldap.User{
						UserID: userID,
						UserDN: uaaUser.ExternalID,
						Email:  uaaUser.Email,
					})
				}
			}
		} else {
			lo.G.Debugf("User [%s] not found in UAA, executing ldap lookup", userID)
			user, err := m.LdapMgr.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			if user != nil {
				ldapUsers = append(ldapUsers, *user)
			} else {
				lo.G.Infof("User %s not found in ldap", userID)
			}
		}
	}
	lo.G.Debugf("LdapUsers before unique check: %+v", ldapUsers)
	ldapUsersToReturn := []ldap.User{}
	uniqueLDAPUsers := make(map[string]ldap.User)
	for _, ldapUser := range ldapUsers {
		if len(strings.TrimSpace(ldapUser.UserDN)) == 0 {
			lo.G.Debugf("User [%s] has a blank externalID", ldapUser.UserID)
			ldapUsersToReturn = append(ldapUsersToReturn, ldapUser)
		} else {
			uniqueLDAPUsers[strings.ToUpper(ldapUser.UserDN)] = ldapUser
		}
	}
	for _, uniqueLDAPUser := range uniqueLDAPUsers {
		ldapUsersToReturn = append(ldapUsersToReturn, uniqueLDAPUser)
	}
	return ldapUsersToReturn, nil
}

func (m *DefaultManager) UpdateUserInfo(user ldap.User) ldap.User {
	userID := strings.ToLower(user.UserID)
	externalID := user.UserDN
	email := user.Email
	if m.LdapConfig.Origin != "ldap" {
		if m.LdapConfig.UseIDForSAMLUser {
			userID = strings.ToLower(user.UserID)
			externalID = user.UserID
		} else {
			userID = strings.ToLower(user.Email)
			externalID = user.Email
		}
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
