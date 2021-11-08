package user

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncLdapUsers(roleUsers *RoleUsers, uaaUsers *uaa.Users, usersInput UsersInput) error {

	if m.LdapConfig.Enabled {
		ldapUsers, err := m.GetLDAPUsers(uaaUsers, usersInput)
		if err != nil {
			return err
		}
		lo.G.Debugf("LdapUsers: %+v", ldapUsers)
		for _, inputUser := range ldapUsers {
			userToUse := m.UpdateUserInfo(inputUser)
			userID := userToUse.UserID
			userList := uaaUsers.GetByName(userID)
			origin := userToUse.Origin
			if origin == "" {
				return fmt.Errorf("Unable to find user %s for origin %s", userID, origin)
			}
			if len(userList) == 0 {
				lo.G.Debug("User", userID, "doesn't exist in cloud foundry, so creating user")
				if userGUID, err := m.UAAMgr.CreateExternalUser(userID, userToUse.Email, userToUse.UserDN, origin); err != nil {
					lo.G.Errorf("Unable to create user %s with error %s", userID, err.Error())
					continue
				} else {
					uaaUsers.Add(uaa.User{
						Username:   userID,
						ExternalID: userToUse.UserDN,
						Origin:     userToUse.Origin,
						Email:      userToUse.Email,
						GUID:       userGUID,
					})
				}
			}
			if !roleUsers.HasUserForOrigin(userID, origin) {
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

func (m *DefaultManager) FilterUsers(userFilter string, userFilterMode string, dnList []string, groupName string) []string {
	filteredDNList := []string{}
	filter, _ := regexp.Compile(userFilter)
	filterMode := userFilterMode
	for _, userDN := range dnList {
		if filter.MatchString(userDN) {
			if filterMode == "include" {
				filteredDNList = append(filteredDNList, userDN)
			} else if filterMode == "exclude" {
				lo.G.Debugf("Removed user %s from group %s because of filter %s and mode %s", userDN, groupName, filter, filterMode)
			}
		} else {
			if filterMode == "include" {
				lo.G.Debugf("Removed user %s from group %s because of filter %s and mode %s", userDN, groupName, filter, filterMode)
			} else if filterMode == "exclude" {
				filteredDNList = append(filteredDNList, userDN)
			}
		}
	}
	return filteredDNList
}

func (m *DefaultManager) GetLDAPUsers(uaaUsers *uaa.Users, usersInput UsersInput) ([]ldap.User, error) {
	var ldapUsers []ldap.User

	// if saml groups are present, saml group members get the saml origin (LdapConfig.origin)
	// ldap_group members and ldap_users get the origin ldap
	// ldap_group members than are interpreted with LdapConfig.LdapOrigin
	// if LdapOrigin is not present it defaults to origin, in this case the saml origin
	if m.LdapConfig.LdapOrigin == "" {
		m.LdapConfig.LdapOrigin = m.LdapConfig.Origin
	}

	if m.LdapConfig.LdapUserFilterMode == "" {
		m.LdapConfig.LdapUserFilterMode = "include"
	}
	if m.LdapConfig.SamlUserFilterMode == "" {
		m.LdapConfig.SamlUserFilterMode = "include"
	}

	// only when ldapOriging is set AND SamlGroups are present, ldapOrigin is used
	// for migration only - remove when migrated
	length := len(usersInput.UniqueSamlGroupNames())
	originForLdapGroups := ""
	if length > 0 {
		originForLdapGroups = m.LdapConfig.LdapOrigin
		lo.G.Debugf("SAML Groups found %s", usersInput.UniqueSamlGroupNames())
	} else {
		originForLdapGroups = m.LdapConfig.Origin
	}
	// remove when migrated

	for _, groupName := range usersInput.UniqueLdapGroupNames() {
		lo.G.Debugf("LDAP GROUP found %s", groupName)
		userDNList, err := m.LdapMgr.GetUserDNs(groupName)
		if err != nil {
			return nil, err
		}
		lo.G.Debugf("LDAP group users pre filter %s", userDNList)

		// filter ldap users
		lo.G.Debugf("LDAP filter %s", m.LdapConfig.LdapUserFilter)
		if m.LdapConfig.LdapUserFilter != "" {
			filteredDNList := m.FilterUsers(m.LdapConfig.LdapUserFilter, m.LdapConfig.LdapUserFilterMode, userDNList, groupName)
			userDNList = filteredDNList
		}

		lo.G.Debugf("LDAP group users post filter %s", userDNList)

		for _, userDN := range userDNList {
			lo.G.Debugf("Checking for userDN %s", userDN)
			uaaUser := uaaUsers.GetByExternalID(userDN)
			if uaaUser != nil {
				lo.G.Debugf("UserDN [%s] found in UAA as [%s], skipping ldap lookup", userDN, uaaUser.Username)
				ldapUsers = append(ldapUsers, ldap.User{
					UserID: uaaUser.Username,
					UserDN: userDN,
					Email:  uaaUser.Email,
					// remove when migrated
					Origin: originForLdapGroups,
				})
			} else {
				lo.G.Debugf("UserDN [%s] not found in UAA, executing ldap lookup", userDN)
				user, err := m.LdapMgr.GetUserByDN(userDN)
				if err != nil {
					return nil, err
				}
				if user != nil {
					// remove when migrated
					user.Origin = originForLdapGroups
					ldapUsers = append(ldapUsers, *user)
				} else {
					lo.G.Debugf("UserDN %s not found in ldap", userDN)
				}
			}
		}
	}

	for _, groupName := range usersInput.UniqueSamlGroupNames() {
		lo.G.Debugf("SAML Group %v", groupName)
		userDNList, err := m.LdapMgr.GetUserDNs(groupName)
		if err != nil {
			return nil, err
		}

		// filter saml users
		lo.G.Debugf("user list before %+v", userDNList)

		if m.LdapConfig.SamlUserFilter != "" {
			lo.G.Debugf("SAML filter %s", m.LdapConfig.SamlUserFilter)
			filteredDNList := m.FilterUsers(m.LdapConfig.SamlUserFilter, m.LdapConfig.SamlUserFilterMode, userDNList, groupName)
			userDNList = filteredDNList
		}

		lo.G.Debugf("user list after %v", userDNList)

		for _, userDN := range userDNList {
			lo.G.Debugf("Checking for userDN %s", userDN)
			uaaUser := uaaUsers.GetByExternalID(userDN)
			if uaaUser != nil {
				lo.G.Debugf("UserDN [%s] found in UAA as [%s], skipping ldap lookup", userDN, uaaUser.Username)
				ldapUsers = append(ldapUsers, ldap.User{
					UserID: uaaUser.Username,
					UserDN: userDN,
					Email:  uaaUser.Email,
					Origin: m.LdapConfig.Origin,
				})
			} else {
				lo.G.Debugf("UserDN [%s] not found in UAA, executing ldap lookup", userDN)
				user, err := m.LdapMgr.GetUserByDN(userDN)
				if err != nil {
					return nil, err
				}
				if user != nil {
					user.Origin = m.LdapConfig.Origin
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
				if strings.EqualFold(uaaUser.Origin, m.LdapConfig.LdapOrigin) {
					ldapUsers = append(ldapUsers, ldap.User{
						UserID: userID,
						UserDN: uaaUser.ExternalID,
						Email:  uaaUser.Email,
						Origin: m.LdapConfig.LdapOrigin,
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
				user.Origin = m.LdapConfig.LdapOrigin
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
	lo.G.Debugf("LdapUsers to return: %+v", ldapUsersToReturn)
	return ldapUsersToReturn, nil
}

func (m *DefaultManager) UpdateUserInfo(user ldap.User) ldap.User {
	userID := strings.ToLower(user.UserID)
	externalID := user.UserDN
	email := user.Email
	origin := user.Origin
	if origin != "ldap" {
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
		Origin: origin,
	}
}
