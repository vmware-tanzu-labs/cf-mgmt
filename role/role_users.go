package role

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func InitRoleUsers() *RoleUsers {
	return &RoleUsers{
		users:         make(map[string][]RoleUser),
		orphanedUsers: make(map[string]string),
	}
}

func NewRoleUsers(users []*resource.User, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()
	for _, user := range users {
		uaaUser := uaaUsers.GetByID(user.GUID)
		if uaaUser == nil {
			lo.G.Debugf("User with guid[%s] is not found in UAA", user.GUID)
			roleUsers.addOrphanedUser(user.GUID)
			continue
		}
		roleUser := RoleUser{
			UserName: uaaUser.Username,
			Origin:   uaaUser.Origin,
			GUID:     uaaUser.GUID,
		}

		if roleUser.UserName == "" {
			return nil, fmt.Errorf("username is blank for user with id %s", user.GUID)
		}
		roleUsers.AddUser(roleUser)
	}

	return roleUsers, nil
}

func (r *RoleUsers) OrphanedUsers() []string {
	var userList []string
	for _, userGUID := range r.orphanedUsers {
		userList = append(userList, userGUID)
	}
	return userList
}

// func (r *RoleUsers) HasUser(userName string) bool {
// 	_, ok := r.users[strings.ToLower(userName)]
// 	return ok
// }

func (r *RoleUsers) HasUserForGUID(userName, userGUID string) bool {
	userList := r.users[strings.ToLower(userName)]
	for _, user := range userList {
		if strings.EqualFold(user.GUID, userGUID) {
			return true
		}
	}
	return false
}

func (r *RoleUsers) HasUserForOrigin(userName, origin string) bool {
	userList := r.users[strings.ToLower(userName)]
	for _, user := range userList {
		if strings.EqualFold(user.Origin, origin) {
			return true
		}
	}
	return false
}

func (r *RoleUsers) RemoveUserForOrigin(userName, origin string) {
	var result []RoleUser
	userList := r.users[strings.ToLower(userName)]
	for _, user := range userList {
		if !strings.EqualFold(user.Origin, origin) {
			result = append(result, user)
		}
	}
	if len(result) == 0 {
		delete(r.users, strings.ToLower(userName))
	} else {
		r.users[strings.ToLower(userName)] = result
	}
}

func (r *RoleUsers) AddUsers(roleUsers []RoleUser) {
	for _, user := range roleUsers {
		r.AddUser(user)
	}
}

func (r *RoleUsers) AddOrphanedUsers(userGUIDs []string) {
	for _, userGUID := range userGUIDs {
		r.addOrphanedUser(userGUID)
	}
}
func (r *RoleUsers) addOrphanedUser(userGUID string) {
	r.orphanedUsers[strings.ToLower(userGUID)] = userGUID
}

func (r *RoleUsers) AddUser(roleUser RoleUser) {
	userList := r.users[strings.ToLower(roleUser.UserName)]
	userList = append(userList, roleUser)
	r.users[strings.ToLower(roleUser.UserName)] = userList
}

func (r *RoleUsers) Users() []RoleUser {
	var result []RoleUser
	if r.users != nil {
		for _, originUsers := range r.users {
			result = append(result, originUsers...)
		}
	}
	return result
}
