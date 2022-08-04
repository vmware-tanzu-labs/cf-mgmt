package user

import (
	"fmt"
	"net/url"
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

func (m *DefaultManager) ListOrgUsers(orgGUID string) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListV3OrganizationRolesByGUIDAndType(orgGUID, ORG_USER)
	if err != nil {
		return nil, err
	}
	return NewRoleUsers(users, m.UAAUsers)
}

func (m *DefaultManager) ListOrgUsersByRole(orgGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	managers := []cfclient.V3User{}
	billingManagers := []cfclient.V3User{}
	auditors := []cfclient.V3User{}
	orgUser := []cfclient.V3User{}
	query := url.Values{}
	query["organization_guids"] = []string{orgGUID}
	roles, err := m.Client.ListV3RolesByQuery(query)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, role := range roles {
		user, err := m.getUserForGUID(role.Relationships["user"].Data.GUID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if role.Type == ORG_MANAGER {
			managers = append(managers, *user)
		} else if role.Type == ORG_BILLING_MANAGER {
			billingManagers = append(billingManagers, *user)
		} else if role.Type == ORG_AUDITOR {
			auditors = append(auditors, *user)
		} else if role.Type == ORG_USER {
			orgUser = append(orgUser, *user)
		} else {
			return nil, nil, nil, nil, fmt.Errorf("type of %s is unknown", role.Type)
		}
	}
	orgUserRoles, err := NewRoleUsers(orgUser, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	managerRoleUsers, err := NewRoleUsers(managers, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	billingManagerRoleUsers, err := NewRoleUsers(billingManagers, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	auditorRoleUsers, err := NewRoleUsers(auditors, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return orgUserRoles, managerRoleUsers, billingManagerRoleUsers, auditorRoleUsers, nil
}

func (m *DefaultManager) ListSpaceUsersByRole(spaceGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	managers := []cfclient.V3User{}
	developers := []cfclient.V3User{}
	auditors := []cfclient.V3User{}
	supporters := []cfclient.V3User{}
	query := url.Values{}
	query["space_guids"] = []string{spaceGUID}
	roles, err := m.Client.ListV3RolesByQuery(query)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	for _, role := range roles {
		user, err := m.getUserForGUID(role.Relationships["user"].Data.GUID)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		if role.Type == SPACE_MANAGER {
			managers = append(managers, *user)
		} else if role.Type == SPACE_DEVELOPER {
			developers = append(developers, *user)
		} else if role.Type == SPACE_AUDITOR {
			auditors = append(auditors, *user)
		} else if role.Type == SPACE_SUPPORTER {
			supporters = append(supporters, *user)
		} else {
			return nil, nil, nil, nil, fmt.Errorf("type of %s is unknown", role.Type)
		}
	}
	managerRoleUsers, err := NewRoleUsers(managers, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	developerRoleUsers, err := NewRoleUsers(developers, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	auditorRoleUsers, err := NewRoleUsers(auditors, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	supporterRoleUsers, err := NewRoleUsers(supporters, m.UAAUsers)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return managerRoleUsers, developerRoleUsers, auditorRoleUsers, supporterRoleUsers, nil
}

func (m *DefaultManager) getUserForGUID(guid string) (*cfclient.V3User, error) {
	if user, ok := m.CFUsers[guid]; ok {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found for guid [%s]", guid)
}
