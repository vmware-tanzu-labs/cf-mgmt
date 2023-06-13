package user

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
)

func NewRoleUsers(users []cfclient.V3User, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()
	for _, user := range users {
		uaaUser := uaaUsers.GetByID(user.GUID)
		if uaaUser == nil {
			roleUsers.addOrphanedUser(user.GUID)
			continue
		}
		roleUser := RoleUser{
			UserName: uaaUser.Username,
			Origin:   uaaUser.Origin,
			GUID:     uaaUser.GUID,
		}

		if roleUser.UserName == "" {
			return nil, fmt.Errorf("Username is blank for user with id %s", user.GUID)
		}
		roleUsers.addUser(roleUser)
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

func (r *RoleUsers) HasUser(userName string) bool {
	_, ok := r.users[strings.ToLower(userName)]
	return ok
}

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
		r.addUser(user)
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

func (r *RoleUsers) addUser(roleUser RoleUser) {
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

func (m *DefaultManager) ListOrgUsersByRole(orgGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	if m.OrgRoles == nil {
		err := m.initializeOrgUserRolesMap()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return m.getOrgRole(orgGUID, ORG_USER), m.getOrgRole(orgGUID, ORG_MANAGER), m.getOrgRole(orgGUID, ORG_BILLING_MANAGER), m.getOrgRole(orgGUID, ORG_AUDITOR), nil
}

func (m *DefaultManager) ListSpaceUsersByRole(spaceGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {

	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	if m.SpaceRoles == nil {
		err := m.initializeSpaceUserRolesMap()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return m.getSpaceRole(spaceGUID, SPACE_MANAGER), m.getSpaceRole(spaceGUID, SPACE_DEVELOPER), m.getSpaceRole(spaceGUID, SPACE_AUDITOR), m.getSpaceRole(spaceGUID, SPACE_SUPPORTER), nil
}

func (m *DefaultManager) getOrgRole(orgGUID, role string) *RoleUsers {
	orgRoles := m.OrgRoles[orgGUID]
	if orgRoles == nil {
		return InitRoleUsers()
	}
	roleUser := orgRoles[role]
	if roleUser == nil {
		return InitRoleUsers()
	}
	return roleUser
}

func (m *DefaultManager) getSpaceRole(spaceGUID, role string) *RoleUsers {
	spaceRoles := m.SpaceRoles[spaceGUID]
	if spaceRoles == nil {
		return InitRoleUsers()
	}
	roleUser := spaceRoles[role]
	if roleUser == nil {
		return InitRoleUsers()
	}
	return roleUser
}

func (m *DefaultManager) getUserForGUID(guid string) (*cfclient.V3User, error) {
	cfUsersMap, err := m.GetCFUsers()
	if err != nil {
		return nil, err
	}
	if user, ok := cfUsersMap[guid]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found for guid [%s]", guid)
}
