package user

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/azureAD"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncAzureADUsers(roleUsers *role.RoleUsers, usersInput UsersInput) error {
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
			uaaUser := m.UAAUsers.GetByNameAndOrigin(userID, origin)
			lo.G.Debugf("SyncAzureADUsers: Processing user: %s", userID)
			if uaaUser == nil {
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
				if err := usersInput.AddUser(usersInput.OrgGUID, usersInput.EntityName(), usersInput.EntityGUID(), user.Username, user.GUID); err != nil {
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
	// a hack, at the moment, as the ldap group names are getting abused to also store the aad groups
	for _, groupName := range usersInput.UniqueLdapGroupNames() {

		userUPNList, err := m.AzureADMgr.GraphGetGroupMembers(m.AzureADMgr.GetADToken(), groupName)
		if err != nil {
			return nil, err
		}
		for _, userUPN := range userUPNList {
			lo.G.Debugf("AAD Adding userUPN %s", userUPN)
			// Check if user is also a SAML user, if so, it has already been added, if not, then add it here
			// NOTE: The SAML User List we get here is the list for the Role we are adding the user to (through the AddUser function, whcih maps to AssociateXXXXXXXXRole
			alreadyMember := false
			for _, u := range usersInput.SamlUsers {
				if strings.EqualFold(u, userUPN) {
					lo.G.Debugf("Group member %s is already defined as SAML user, skipping", userUPN)
					alreadyMember = true
					continue
				}
			}
			if !alreadyMember {
				azureADUsers = append(azureADUsers, azureAD.UserType{
					Upn: userUPN,
				})
			}

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
