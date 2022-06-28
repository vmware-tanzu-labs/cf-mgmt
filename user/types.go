package user

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
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
	ListSpaceAuditors(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListSpaceDevelopers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListSpaceManagers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListSpaceSupporters(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgAuditors(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgBillingManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgUsers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
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
}

type LdapManager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByDN(userDN string) (*ldap.User, error)
	GetUserByID(userID string) (*ldap.User, error)
	Close()
}
