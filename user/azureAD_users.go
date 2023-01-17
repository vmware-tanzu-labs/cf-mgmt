package user

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/azureAD"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncAzureADUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	origin := m.AzureADConfig.UserOrigin
	if m.AzureADConfig.Enabled {
		azureADUsers, err := m.GetAzureADUsers(usersInput)
		if err != nil {
			return err
		}
		lo.G.Debugf("azureADUsers: %+v", azureADUsers)
		for _, inputUser := range azureADUsers {
			userToUse := m.UpdateADUserInfo(inputUser)
			userID := userToUse.Upn
			userList := m.UAAUsers.GetByName(userID)
			if len(userList) == 0 {
				lo.G.Debugf("AAD User %s doesn't exist in cloud foundry, so creating user", userToUse.Upn)
				if userGUID, err := m.UAAMgr.CreateExternalUser(userToUse.Upn, userToUse.Upn, userToUse.Upn, m.AzureADConfig.UserOrigin); err != nil {
					lo.G.Errorf("Unable to create AAD user %s with error %s", userToUse.Upn, err.Error())
					continue
				} else {
					m.UAAUsers.Add(uaa.User{
						Username:   userToUse.Upn,
						ExternalID: userToUse.Upn,
						Origin:     m.AzureADConfig.UserOrigin,
						Email:      userToUse.Upn,
						GUID:       userGUID,
					})
				}
			}
			if !roleUsers.HasUserForOrigin(userID, origin) {
				user := m.UAAUsers.GetByNameAndOrigin(userID, origin)
				if user == nil {
					return fmt.Errorf("Unable to find user %s for origin %s", userID, origin)
				}
				if err := usersInput.AddUser(usersInput, user.Username, user.GUID); err != nil {
					return errors.Wrap(err, fmt.Sprintf("User %s with origin %s", user.Username, user.Origin))
				}
			} else {
				lo.G.Debugf("AAD User[%s] found in role", userID)
				roleUsers.RemoveUserForOrigin(userID, origin)
			}
		}
	} else {
		lo.G.Debug("Skipping Azure AD sync as it is disabled (enable by updating config/azureAD.yml)")
	}
	return nil
}

func (m *DefaultManager) GetAzureADUsers(usersInput UsersInput) ([]azureAD.UserType, error) {
	var azureADUsers []azureAD.UserType
	// TODO Base of LDAP groups for now, later change to separate group entry for AAD
	for _, groupName := range usersInput.UniqueLdapGroupNames() {

		userUPNList, err := m.AzureADMgr.GraphGetGroupMembers(m.AzureADMgr.GetADToken(),groupName) 
		if err != nil {
			return nil, err
		}
		for _, userUPN := range userUPNList {
			lo.G.Debugf("AAD Checking for userDN %s", userUPN)
			azureADUsers = append(azureADUsers, azureAD.UserType{
				Upn: userUPN,
			})
			// uaaUser := m.UAAUsers.GetByExternalID(userUPN)
			// if uaaUser != nil {
			// 	lo.G.Debugf("AAD userUPN [%s] found in UAA as [%s]", userUPN, uaaUser.Username)
			// 	azureADUsers = append(azureADUsers, azureAD.UserType{
			// 		Upn: uaaUser.Username,
			// 	})
			// } else {
			// 	lo.G.Debugf("userUPN [%s] not found in UAA, executing azure AD lookup", userUPN)
			// 	user, err := m.LdapMgr.GetUserByDN(userUPN) // TODO implement azure AD version
			// 	if err != nil {
			// 		return nil, err
			// 	}
			// 	if user != nil {
			// 		azureADUsers = append(azureADUsers, *user) // TODO: Solves itself after other TODO's
			// 	} else {
			// 		lo.G.Infof("user %s not found in Azure AD", userUPN)
			// 	}
			// }
		}
	}
	lo.G.Debugf("Azure AD Users before uniqueness check: %+v", azureADUsers)

	ADUsersToReturn := []azureAD.UserType{}
	uniqueADUsers := make(map[string]azureAD.UserType)
	for _, azureADUser := range azureADUsers {
		uniqueADUsers[strings.ToUpper(azureADUser.Upn)] = azureADUser
	}	
	for _, uniqueADUser := range uniqueADUsers {
		ADUsersToReturn = append(ADUsersToReturn, uniqueADUser)
	}
	return ADUsersToReturn, nil
}


func (m *DefaultManager) UpdateADUserInfo(user azureAD.UserType) azureAD.UserType {
	upnToLower := strings.ToLower(user.Upn)

	return azureAD.UserType{
		Upn: upnToLower,
	}
}
