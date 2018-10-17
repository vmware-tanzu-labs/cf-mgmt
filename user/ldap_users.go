package user

import (
	"fmt"
	"strings"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncLdapUsers(roleUsers map[string]string, uaaUsers map[string]*uaaclient.User, updateUsersInput UpdateUsersInput) error {
	if m.LdapConfig.Enabled {
		ldapUsers, err := m.GetLDAPUsers(uaaUsers, updateUsersInput)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %+v", ldapUsers)
		for _, inputUser := range ldapUsers {
			userToUse := m.UpdateUserInfo(inputUser)
			userID := userToUse.UserID
			if _, ok := roleUsers[userID]; !ok {
				lo.G.Debugf("User[%s] not found in: %v", userID, roleUsers)
				if _, userExists := uaaUsers[userID]; !userExists {
					lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
					if err := m.UAAMgr.CreateExternalUser(userID, userToUse.Email, userToUse.UserDN, m.LdapConfig.Origin); err != nil {
						lo.G.Errorf("Unable to create user %s with error %s", userID, err.Error())
						continue
					} else {
						uaaUser := &uaaclient.User{
							Username:   userID,
							ExternalID: userToUse.UserDN,
							Origin:     m.LdapConfig.Origin,
							Emails:     []uaaclient.Email{uaaclient.Email{Value: userToUse.Email}},
						}
						uaaUsers[userID] = uaaUser
						uaaUsers[userToUse.UserDN] = uaaUser
					}
				}
				if err := updateUsersInput.AddUser(updateUsersInput, userID); err != nil {
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

func (m *DefaultManager) GetLDAPUsers(uaaUsers map[string]*uaaclient.User, updateUsersInput UpdateUsersInput) ([]ldap.User, error) {
	var ldapUsers []ldap.User
	for _, groupName := range updateUsersInput.LdapGroupNames {
		userDNList, err := m.LdapMgr.GetUserDNs(groupName)
		if err != nil {
			return nil, err
		}
		for _, userDN := range userDNList {
			if uaaUser, ok := uaaUsers[strings.ToLower(userDN)]; ok {
				lo.G.Debugf("UserDN [%s] found in UAA, skipping ldap lookup", userDN)
				ldapUsers = append(ldapUsers, ldap.User{
					UserID: uaaUser.Username,
					UserDN: userDN,
					Email:  Email(uaaUser),
				})
			} else {
				user, err := m.LdapMgr.GetUserByDN(userDN)
				if err != nil {
					return nil, err
				}
				if user != nil {
					ldapUsers = append(ldapUsers, *user)
				}
			}
		}
	}
	for _, userID := range updateUsersInput.LdapUsers {
		if uaaUser, ok := uaaUsers[strings.ToLower(userID)]; ok {
			lo.G.Debugf("UserID [%s] found in UAA, skipping ldap lookup", userID)
			ldapUsers = append(ldapUsers, ldap.User{
				UserID: userID,
				UserDN: uaaUser.ExternalID,
				Email:  Email(uaaUser),
			})
		} else {
			user, err := m.LdapMgr.GetUserByID(userID)
			if err != nil {
				return nil, err
			}
			if user != nil {
				ldapUsers = append(ldapUsers, *user)
			}
		}
	}
	return ldapUsers, nil
}

func Email(u *uaaclient.User) string {
	for _, email := range u.Emails {
		if *email.Primary {
			return email.Value
		}
	}
	if len(u.Emails) > 0 {
		return u.Emails[0].Value
	}
	return ""
}

func (m *DefaultManager) UpdateUserInfo(user ldap.User) ldap.User {
	userID := strings.ToLower(user.UserID)
	externalID := user.UserDN
	email := user.Email
	if m.LdapConfig.Origin != "ldap" {
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
