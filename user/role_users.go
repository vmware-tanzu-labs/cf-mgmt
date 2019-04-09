package user

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pkg/errors"
	"github.com/xchapter7x/lo"
)

func NewRoleUsers(users []cfclient.User, uaaUsers *uaa.Users) (*RoleUsers, []string, error) {
	roleUsers := InitRoleUsers()
	orphanedUsers := []string{}
	for _, user := range users {
		uaaUser := uaaUsers.GetByID(user.Guid)
		if uaaUser == nil {
			orphanedUsers = append(orphanedUsers, user.Guid)
			continue
		}
		roleUser := RoleUser{
			UserName: uaaUser.Username,
			Origin:   uaaUser.Origin,
			GUID:     uaaUser.GUID,
		}

		if roleUser.UserName == "" {
			return nil, nil, fmt.Errorf("Username is blank for user with id %s", user.Guid)
		}
		roleUsers.addUser(roleUser)
	}
	return roleUsers, orphanedUsers, nil
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
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
}
func (m *DefaultManager) ListSpaceDevelopers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListSpaceDevelopers(spaceGUID)
	if err != nil {
		return nil, err
	}
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
}
func (m *DefaultManager) ListSpaceManagers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListSpaceManagers(spaceGUID)
	if err != nil {
		return nil, err
	}
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
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

func (m *DefaultManager) usersInOrgRoles(orgName, orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()

	orgAuditors, err := m.ListOrgAuditors(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org auditors for org %s", orgName))
	}
	roleUsers.AddUsers(orgAuditors.Users())

	orgManagers, err := m.ListOrgManagers(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org managers for org %s", orgName))
	}
	roleUsers.AddUsers(orgManagers.Users())

	orgBillingManagers, err := m.ListOrgBillingManagers(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org billing managers for org %s", orgName))
	}
	roleUsers.AddUsers(orgBillingManagers.Users())

	spaces, err := m.listSpaces(orgGUID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing spaces for org %s", orgName))
	}
	for _, space := range spaces {
		spaceAuditors, err := m.ListSpaceAuditors(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space auditors for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceAuditors.Users())

		spaceDevelopers, err := m.ListSpaceDevelopers(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space developers for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceDevelopers.Users())

		spaceManagers, err := m.ListSpaceManagers(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space managers for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceManagers.Users())
	}

	return roleUsers, nil
}

func (m *DefaultManager) ListOrgAuditors(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgAuditors(orgGUID)
	if err != nil {
		return nil, err
	}
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
}
func (m *DefaultManager) ListOrgBillingManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgBillingManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
}
func (m *DefaultManager) ListOrgManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), nil
	}
	users, err := m.Client.ListOrgManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	roleUsers, orphanedUsers, err := NewRoleUsers(users, uaaUsers)
	if err != nil {
		return nil, err
	}
	err = m.removeOrphanedUsers(orphanedUsers)
	if err != nil {
		return nil, err
	}
	return roleUsers, nil
}

func (m *DefaultManager) removeOrphanedUsers(orphanedUsers []string) error {
	for _, orphanedUser := range orphanedUsers {
		lo.G.Infof("Deleting orphaned CF user with guid %s", orphanedUser)
		err := m.Client.DeleteUser(orphanedUser)
		if err != nil {
			return err
		}
	}

	return nil
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
