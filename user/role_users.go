package user

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func NewRoleUsers(users []cfclient.User, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()
	for _, user := range users {
		uaaUser := uaaUsers.GetByID(user.Guid)
		if uaaUser == nil {
			roleUsers.addOrphanedUser(user.Guid)
			continue
		}
		roleUser := RoleUser{
			UserName: uaaUser.Username,
			Origin:   uaaUser.Origin,
			GUID:     uaaUser.GUID,
		}

		if roleUser.UserName == "" {
			return nil, fmt.Errorf("Username is blank for user with id %s", user.Guid)
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
			for _, user := range originUsers {
				result = append(result, user)
			}
		}
	}
	return result
}

func (m *DefaultManager) ListSpaceAuditors(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListSpaceAuditors(spaceGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}
func (m *DefaultManager) ListSpaceDevelopers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListSpaceDevelopers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}
func (m *DefaultManager) ListSpaceManagers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListSpaceManagers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}

func (m *DefaultManager) listSpaceAuditors(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListSpaceAuditors(input.SpaceGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s, Space %s and role %s: %+v", input.OrgName, input.SpaceName, "space-auditor", roleUsers)
	}
	return roleUsers, err
}
func (m *DefaultManager) listSpaceDevelopers(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListSpaceDevelopers(input.SpaceGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s, Space %s and role %s: %+v", input.OrgName, input.SpaceName, "space-developer", roleUsers)
	}
	return roleUsers, err
}
func (m *DefaultManager) listSpaceManagers(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListSpaceManagers(input.SpaceGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s, Space %s and role %s: %+v", input.OrgName, input.SpaceName, "space-manager", roleUsers)
	}
	return roleUsers, err
}

func (m *DefaultManager) ListOrgAuditors(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgAuditors(orgGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}
func (m *DefaultManager) ListOrgBillingManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgBillingManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}

func (m *DefaultManager) ListOrgManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, uaaUsers)
}

func (m *DefaultManager) listOrgAuditors(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListOrgAuditors(input.OrgGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s and role %s: %+v", input.OrgName, "org-auditor", roleUsers)
	}
	return roleUsers, err
}
func (m *DefaultManager) listOrgBillingManagers(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListOrgBillingManagers(input.OrgGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s and role %s: %+v", input.OrgName, "org-billing-manager", roleUsers)
	}
	return roleUsers, err
}
func (m *DefaultManager) listOrgManagers(input UsersInput, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers, err := m.ListOrgManagers(input.OrgGUID, uaaUsers)
	if err == nil {
		lo.G.Debugf("RoleUsers for Org %s and role %s: %+v", input.OrgName, "org-manager", roleUsers)
	}
	return roleUsers, err
}
