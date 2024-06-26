package role

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

type DefaultManager struct {
	RoleClient      CFRoleClient
	UserClient      CFUserClient
	JobClient       CFJobClient
	OrgRoles        map[string]map[string]*RoleUsers
	SpaceRoles      map[string]map[string]*RoleUsers
	UAAMgr          uaa.Manager
	Peek            bool
	OrgRolesUsers   map[string]map[string]map[string]string
	SpaceRolesUsers map[string]map[string]map[string]string
}

func New(roleClient CFRoleClient, userClient CFUserClient, jobClient CFJobClient, uaaMgr uaa.Manager, peek bool) Manager {
	return &DefaultManager{
		RoleClient: roleClient,
		UserClient: userClient,
		JobClient:  jobClient,
		UAAMgr:     uaaMgr,
		Peek:       peek,
	}
}

func (m *DefaultManager) ClearRoles() {
	m.OrgRoles = nil
	m.SpaceRoles = nil
	m.OrgRolesUsers = nil
	m.SpaceRolesUsers = nil
}

func (m *DefaultManager) ListOrgRoles() ([]*resource.Role, error) {

	var rolesToReturn []*resource.Role
	roles, err := m.listRolesForType(resource.OrganizationRoleAuditor.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.OrganizationRoleBillingManager.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.OrganizationRoleManager.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.OrganizationRoleUser.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)
	return rolesToReturn, err
}

func (m *DefaultManager) ListSpaceRoles() ([]*resource.Role, error) {
	var rolesToReturn []*resource.Role
	roles, err := m.listRolesForType(resource.SpaceRoleAuditor.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.SpaceRoleDeveloper.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.SpaceRoleManager.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)

	roles, err = m.listRolesForType(resource.SpaceRoleSupporter.String())
	if err != nil {
		return nil, err
	}
	rolesToReturn = append(rolesToReturn, roles...)
	return rolesToReturn, err
}

func (m *DefaultManager) listRolesForType(roleType string) ([]*resource.Role, error) {
	roles, err := m.RoleClient.ListAll(context.Background(), &client.RoleListOptions{
		Types: client.Filter{
			Values: []string{roleType},
		},
		ListOptions: &client.ListOptions{
			PerPage: 5000,
		},
	})
	if err != nil {
		return nil, err
	}
	lo.G.Debugf("Found %d roles from type %s", len(roles), roleType)
	err = m.checkResultsAllReturned(roles)
	if err != nil {
		return nil, err
	}
	m.dumpV3Roles(roles)
	return roles, err
}

func (m *DefaultManager) InitializeSpaceUserRolesMap() error {
	spaceV3UsersRolesMap := make(map[string]map[string][]*uaa.User)
	roles, err := m.ListSpaceRoles()
	if err != nil {
		return err
	}
	for _, role := range roles {
		spaceGUID := role.Relationships.Space.Data.GUID
		user, err := m.getUserForGUID(role.Relationships.User.Data.GUID)
		if err != nil {
			return err
		}
		spaceRoleMap, ok := spaceV3UsersRolesMap[spaceGUID]
		if !ok {
			spaceRoleMap = make(map[string][]*uaa.User)
			spaceV3UsersRolesMap[spaceGUID] = spaceRoleMap
		}
		spaceRoleMap[role.Type] = append(spaceRoleMap[role.Type], user)
	}
	spaceUsersRoleMap := make(map[string]map[string]*RoleUsers)
	for key, val := range spaceV3UsersRolesMap {
		for role, users := range val {
			uaaUsers, err := m.GetUAAUsers()
			if err != nil {
				return err
			}
			roleUsers, err := NewRoleUsers(users, uaaUsers)
			if err != nil {
				return err
			}
			roleMap, ok := spaceUsersRoleMap[key]
			if !ok {
				roleMap = make(map[string]*RoleUsers)
				spaceUsersRoleMap[key] = roleMap
			}
			roleMap[role] = roleUsers
		}
	}
	m.dumpRolesUsers("spaces", spaceUsersRoleMap)
	m.SpaceRoles = spaceUsersRoleMap
	m.SpaceRolesUsers = m.buildUserMap(
		func(role *resource.Role) string {
			return role.Relationships.Space.Data.GUID
		}, roles)
	return nil
}

func (m *DefaultManager) InitializeOrgUserRolesMap() error {
	orgV3UsersRolesMap := make(map[string]map[string][]*uaa.User)
	roles, err := m.ListOrgRoles()
	if err != nil {
		return err
	}
	for _, role := range roles {
		orgGUID := role.Relationships.Org.Data.GUID
		user, err := m.getUserForGUID(role.Relationships.User.Data.GUID)
		if err != nil {
			return err
		}
		orgRoleMap, ok := orgV3UsersRolesMap[orgGUID]
		if !ok {
			orgRoleMap = make(map[string][]*uaa.User)
			orgV3UsersRolesMap[orgGUID] = orgRoleMap
		}
		orgRoleMap[role.Type] = append(orgRoleMap[role.Type], user)
	}
	orgUsersRoleMap := make(map[string]map[string]*RoleUsers)
	for key, val := range orgV3UsersRolesMap {
		for role, users := range val {
			uaaUsers, err := m.GetUAAUsers()
			if err != nil {
				return err
			}
			roleUsers, err := NewRoleUsers(users, uaaUsers)
			if err != nil {
				return err
			}
			roleMap, ok := orgUsersRoleMap[key]
			if !ok {
				roleMap = make(map[string]*RoleUsers)
				orgUsersRoleMap[key] = roleMap
			}
			roleMap[role] = roleUsers
		}
	}
	m.dumpRolesUsers("organizations", orgUsersRoleMap)
	m.OrgRoles = orgUsersRoleMap
	m.OrgRolesUsers = m.buildUserMap(
		func(role *resource.Role) string {
			return role.Relationships.Org.Data.GUID
		}, roles)
	return nil
}

func (m DefaultManager) buildUserMap(keyFunction func(role *resource.Role) string, roles []*resource.Role) map[string]map[string]map[string]string {
	result := make(map[string]map[string]map[string]string)
	for _, role := range roles {
		guid := keyFunction(role)
		roleMap, ok := result[guid]
		if !ok {
			roleMap = make(map[string]map[string]string)
			result[guid] = roleMap
		}
		userMap, ok := roleMap[role.Type]
		if !ok {
			userMap = make(map[string]string)
			roleMap[role.Type] = userMap
		}
		userMap[role.Relationships.User.Data.GUID] = role.GUID
	}
	return result
}

func (m *DefaultManager) dumpRolesUsers(entityType string, entityRoles map[string]map[string]*RoleUsers) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		for guid, roles := range entityRoles {
			for roleType, role := range roles {
				for _, user := range role.Users() {
					lo.G.Debugf("User [%s] with GUID[%s] and origin %s in role %s for entity type [%s/%s]", user.UserName, user.GUID, user.Origin, roleType, entityType, guid)
				}
			}
		}
	}
}

func (m *DefaultManager) dumpV3Roles(roles []*resource.Role) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		for _, role := range roles {
			lo.G.Debugf("For role guid [%s] and role [%s] user guid [%s]", role.GUID, role.Type, role.Relationships.User.Data.GUID)
		}
	}
}

func (m *DefaultManager) checkResultsAllReturned(roles []*resource.Role) error {
	tracker := make(map[string]*resource.Role)
	for _, role := range roles {
		if priorRole, ok := tracker[role.GUID]; !ok {
			tracker[role.GUID] = role
		} else {
			return fmt.Errorf("role with GUID[%s] is returned multiple times, prior role [%s] and current role [%s]", role.GUID, asJson(priorRole), asJson(role))
		}
	}
	return nil
}

func asJson(role *resource.Role) string {
	bytes, _ := json.Marshal(role)
	return string(bytes)
}

func (m *DefaultManager) GetUAAUsers() (*uaa.Users, error) {
	return m.UAAMgr.ListUsers()
}

func (m *DefaultManager) UpdateOrgRoleUsers(orgGUID string, roleUser *RoleUsers) {
	orgRoles, ok := m.OrgRoles[orgGUID]
	if !ok {
		orgRoles = make(map[string]*RoleUsers)
	}
	orgRoles[resource.OrganizationRoleUser.String()] = roleUser
	m.OrgRoles[orgGUID] = orgRoles
}

func (m *DefaultManager) ListOrgUsersByRole(orgGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	if m.OrgRoles == nil {
		err := m.InitializeOrgUserRolesMap()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return m.getOrgRole(orgGUID, resource.OrganizationRoleUser.String()), m.getOrgRole(orgGUID, resource.OrganizationRoleManager.String()), m.getOrgRole(orgGUID, resource.OrganizationRoleBillingManager.String()), m.getOrgRole(orgGUID, resource.OrganizationRoleAuditor.String()), nil
}

func (m *DefaultManager) ListSpaceUsersByRole(spaceGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error) {

	if m.Peek && strings.Contains(spaceGUID, "dry-run-space-guid") {
		return InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), InitRoleUsers(), nil
	}
	if m.SpaceRoles == nil {
		err := m.InitializeSpaceUserRolesMap()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	}
	return m.getSpaceRole(spaceGUID, resource.SpaceRoleManager.String()), m.getSpaceRole(spaceGUID, resource.SpaceRoleDeveloper.String()), m.getSpaceRole(spaceGUID, resource.SpaceRoleAuditor.String()), m.getSpaceRole(spaceGUID, resource.SpaceRoleSupporter.String()), nil
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

func (m *DefaultManager) getUserForGUID(guid string) (*uaa.User, error) {
	uaaUsers, err := m.GetUAAUsers()
	if err != nil {
		return nil, err
	}
	user := uaaUsers.GetByID(guid)
	if user != nil {
		return user, nil
	} else {
		user, err := m.GetUser(guid)
		if err != nil {
			return nil, err
		}
		return &uaa.User{
			Username: user.Username,
			GUID:     user.GUID,
			Origin:   user.Origin,
		}, nil
	}
}

func (m *DefaultManager) GetUser(userGuid string) (*resource.User, error) {
	user, err := m.UserClient.Get(context.Background(), userGuid)
	return user, err
}

func (m *DefaultManager) DeleteUser(userGuid string) error {
	_, err := m.UserClient.Delete(context.Background(), userGuid)
	return err
}

func (m *DefaultManager) deleteRole(roleGUID string) error {
	jobGUID, err := m.RoleClient.Delete(context.Background(), roleGUID)
	if err != nil {
		return err
	}
	err = m.JobClient.PollComplete(context.Background(), jobGUID, &client.PollingOptions{
		FailedState:   "FAILED",
		Timeout:       time.Second * 30,
		CheckInterval: time.Second,
	})
	return err
}
