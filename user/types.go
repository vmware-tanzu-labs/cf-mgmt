package user

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
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

type CFRoleClient interface {
	Delete(ctx context.Context, guid string) (string, error)
	CreateOrganizationRole(ctx context.Context, organizationGUID, userGUID string, roleType resource.OrganizationRoleType) (*resource.Role, error)
	CreateSpaceRole(ctx context.Context, spaceGUID, userGUID string, roleType resource.SpaceRoleType) (*resource.Role, error)
	ListAll(ctx context.Context, opts *client.RoleListOptions) ([]*resource.Role, error)
	ListIncludeUsersAll(ctx context.Context, opts *client.RoleListOptions) ([]*resource.Role, []*resource.User, error)
}

type CFUserClient interface {
	Delete(ctx context.Context, guid string) (string, error)
	ListAll(ctx context.Context, opts *client.UserListOptions) ([]*resource.User, error)
}

type CFSpaceClient interface {
	ListAll(ctx context.Context, opts *client.SpaceListOptions) ([]*resource.Space, error)
}

type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *client.PollingOptions) error
}

type LdapManager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByDN(userDN string) (*ldap.User, error)
	GetUserByID(userID string) (*ldap.User, error)
	Close()
}
