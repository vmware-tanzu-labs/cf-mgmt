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

const SPACE_AUDITOR string = "space_auditor"
const SPACE_DEVELOPER string = "space_developer"
const SPACE_MANAGER string = "space_manager"
const SPACE_SUPPORTER string = "space_supporter"

type DefaultManager struct {
	RoleClient      CFRoleClient
	UserClient      CFUserClient
	JobClient       CFJobClient
	OrgRoles        map[string]map[string]*RoleUsers
	SpaceRoles      map[string]map[string]*RoleUsers
	CFUsers         map[string]*resource.User
	UAAUsers        *uaa.Users
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
	roles, err := m.RoleClient.ListAll(context.Background(), &client.RoleListOptions{
		Types: client.Filter{
			Values: []string{resource.OrganizationRoleAuditor.String(),
				resource.OrganizationRoleBillingManager.String(),
				resource.OrganizationRoleManager.String(),
				resource.OrganizationRoleUser.String()},
		},
		ListOptions: &client.ListOptions{
			PerPage: 5000,
		},
	})
	if err != nil {
		return nil, err
	}
	lo.G.Debugf("Found %d roles from %s API", len(roles), "organization")
	err = m.checkResultsAllReturned("organization", roles)
	if err != nil {
		return nil, err
	}
	m.dumpV3Roles("organization", roles)
	return roles, err
}

func (m *DefaultManager) ListSpaceRoles() ([]*resource.Role, error) {
	roles, err := m.RoleClient.ListAll(context.Background(), &client.RoleListOptions{
		Types: client.Filter{
			Values: []string{SPACE_AUDITOR, SPACE_DEVELOPER, SPACE_MANAGER, SPACE_SUPPORTER},
		},
		ListOptions: &client.ListOptions{
			PerPage: 5000,
		},
	})
	if err != nil {
		return nil, err
	}
	lo.G.Debugf("Found %d roles from %s API", len(roles), "space")
	err = m.checkResultsAllReturned("space", roles)
	if err != nil {
		return nil, err
	}
	m.dumpV3Roles("space", roles)
	return roles, err
}

func (m *DefaultManager) InitializeSpaceUserRolesMap() error {
	spaceV3UsersRolesMap := make(map[string]map[string][]*resource.User)
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
			spaceRoleMap = make(map[string][]*resource.User)
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
	orgV3UsersRolesMap := make(map[string]map[string][]*resource.User)
	roles, err := m.ListOrgRoles()
	if err != nil {
		return err
	}
	lo.G.Debugf("Found %d roles from %s API", len(roles), "organization")
	err = m.checkResultsAllReturned("organization", roles)
	if err != nil {
		return err
	}
	m.dumpV3Roles("organization", roles)
	for _, role := range roles {
		orgGUID := role.Relationships.Org.Data.GUID
		user, err := m.getUserForGUID(role.Relationships.User.Data.GUID)
		if err != nil {
			return err
		}
		orgRoleMap, ok := orgV3UsersRolesMap[orgGUID]
		if !ok {
			orgRoleMap = make(map[string][]*resource.User)
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

func (m *DefaultManager) dumpV3Roles(entityType string, roles []*resource.Role) {
	level, logging := os.LookupEnv("LOG_LEVEL")
	if logging && strings.EqualFold(level, "DEBUG") {
		for _, role := range roles {
			lo.G.Debugf("For entity [%s/%s] and role [%s] user guid [%s]", entityType, role.GUID, role.Type, role.Relationships.User.Data.GUID)
		}
	}
}

func (m *DefaultManager) checkResultsAllReturned(entityType string, roles []*resource.Role) error {
	tracker := make(map[string]*resource.Role)
	for _, role := range roles {
		if priorRole, ok := tracker[role.GUID]; !ok {
			tracker[role.GUID] = role
		} else {
			return fmt.Errorf("role for type %s with GUID[%s] is returned multiple times, prior role [%s] and current role [%s]", entityType, role.GUID, asJson(priorRole), asJson(role))
		}
	}
	return nil
}

func asJson(role *resource.Role) string {
	bytes, _ := json.Marshal(role)
	return string(bytes)
}

func (m *DefaultManager) GetCFUsers() (map[string]*resource.User, error) {
	if m.CFUsers == nil {
		cfUserMap := make(map[string]*resource.User)
		cfUsers, err := m.UserClient.ListAll(context.Background(), &client.UserListOptions{
			ListOptions: &client.ListOptions{
				PerPage: 5000,
			},
		})
		if err != nil {
			return nil, err
		}
		lo.G.Debug("Begin CFUsers")
		for _, cfUser := range cfUsers {
			cfUserMap[cfUser.GUID] = cfUser
			lo.G.Debugf("CFUser with username [%s] and guid [%s]", cfUser.Username, cfUser.GUID)
		}
		lo.G.Debug("End CFUsers")
		m.CFUsers = cfUserMap
	}
	return m.CFUsers, nil
}

func (m *DefaultManager) GetUAAUsers() (*uaa.Users, error) {
	if m.UAAUsers == nil {
		uaaUsers, err := m.UAAMgr.ListUsers()
		if err != nil {
			return nil, err
		}
		m.UAAUsers = uaaUsers
	}
	return m.UAAUsers, nil
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

func (m *DefaultManager) getUserForGUID(guid string) (*resource.User, error) {
	cfUsersMap, err := m.GetCFUsers()
	if err != nil {
		return nil, err
	}
	if user, ok := cfUsersMap[guid]; ok {
		return user, nil
	}
	return nil, fmt.Errorf("user not found for guid [%s]", guid)
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
