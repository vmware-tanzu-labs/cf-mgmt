package user

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
)

func InitRoleUsers() *RoleUsers {
	return &RoleUsers{
		users:         make(map[string][]RoleUser),
		orphanedUsers: make(map[string]string),
	}
}

type RoleUsers struct {
	users         map[string][]RoleUser
	orphanedUsers map[string]string
}
type RoleUser struct {
	UserName string
	GUID     string
	Origin   string
}

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	InitializeLdap(ldapBindUser, ldapBindPassword, ldapServer string) error
	DeinitializeLdap() error
	UpdateSpaceUsers() error
	UpdateOrgUsers() error
	CleanupOrgUsers() error
	ListOrgUsersByRole(orgGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error)
	ListSpaceUsersByRole(spaceGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error)
}

type CFClient interface {
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)
	DeleteUser(userGuid string) error
	DeleteV3Role(roleGUID string) error
	ListV3SpaceRolesByGUIDAndType(spaceGUID string, roleType string) ([]cfclient.V3User, error)
	ListV3OrganizationRolesByGUIDAndType(orgGUID string, roleType string) ([]cfclient.V3User, error)
	CreateV3OrganizationRole(orgGUID, userGUID, roleType string) (*cfclient.V3Role, error)
	CreateV3SpaceRole(spaceGUID, userGUID, roleType string) (*cfclient.V3Role, error)
	SupportsSpaceSupporterRole() (bool, error)
	ListV3RolesByQuery(query url.Values) ([]cfclient.V3Role, error)
	ListV3UsersByQuery(query url.Values) ([]cfclient.V3User, error)
}

type LdapManager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByDN(userDN string) (*ldap.User, error)
	GetUserByID(userID string) (*ldap.User, error)
	Close()
}
