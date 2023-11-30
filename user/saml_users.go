package user

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) SyncSamlUsers(roleUsers *role.RoleUsers, usersInput UsersInput) error {
	origin := m.LdapConfig.Origin
	uaaUsers, err := m.GetUAAUsers()
	if err != nil {
		return err
	}
	for _, userEmail := range usersInput.UniqueSamlUsers() {
		uaaUser := uaaUsers.GetByNameAndOrigin(userEmail, origin)
		if uaaUser == nil {
			lo.G.Debugf("user %s doesn't exist in cloud foundry with origin %s, so creating user", userEmail, origin)
			if userGUID, err := m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, origin); err != nil {
				return err
			} else {
				m.AddUAAUser(uaa.User{
					Username:   userEmail,
					Email:      userEmail,
					ExternalID: userEmail,
					Origin:     origin,
					GUID:       userGUID,
				})
			}
		}
		user := uaaUsers.GetByNameAndOrigin(userEmail, origin)
		if user == nil {
			return fmt.Errorf("unable to find user %s for origin %s", userEmail, origin)
		}
		if !roleUsers.HasUserForOrigin(userEmail, user.Origin) {
			m.dumpRoleUsers(fmt.Sprintf("Adding user [%s] with guid[%s] with origin [%s] as doesn't exist in users for %s/%s - Role %s", userEmail, user.GUID, origin, usersInput.OrgName, usersInput.SpaceName, usersInput.Role), roleUsers.Users())
			if err := usersInput.AddUser(usersInput.OrgGUID, usersInput.EntityName(), usersInput.EntityGUID(), user.Username, user.GUID); err != nil {
				return errors.Wrap(err, fmt.Sprintf("User %s with origin %s", user.Username, user.Origin))
			}
		} else {
			lo.G.Debugf("User %s already exists", userEmail)
			roleUsers.RemoveUserForOrigin(userEmail, user.Origin)
		}
	}
	return nil
}
