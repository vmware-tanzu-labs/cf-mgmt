package user

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncSamlUsers(roleUsers *RoleUsers, usersInput UsersInput) error {
	origin := m.LdapConfig.Origin
	for _, userEmail := range usersInput.UniqueSamlUsers() {
		userList := m.UAAUsers.GetByName(userEmail)
		if len(userList) == 0 {
			lo.G.Debug("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if userGUID, err := m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				continue
			} else {
				m.UAAUsers.Add(uaa.User{
					Username:   userEmail,
					Email:      userEmail,
					ExternalID: userEmail,
					Origin:     origin,
					GUID:       userGUID,
				})
			}
		}
		user := m.UAAUsers.GetByNameAndOrigin(userEmail, origin)
		if user == nil {
			return fmt.Errorf("Unable to find user %s for origin %s", userEmail, origin)
		}
		if !roleUsers.HasUserForOrigin(userEmail, user.Origin) {
			if err := usersInput.AddUser(usersInput, user.Username, user.GUID); err != nil {
				return errors.Wrap(err, fmt.Sprintf("User %s with origin %s", user.Username, user.Origin))
			}
		} else {
			roleUsers.RemoveUserForOrigin(userEmail, user.Origin)
		}
	}
	return nil
}
