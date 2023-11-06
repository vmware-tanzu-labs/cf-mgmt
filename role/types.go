package role

import (
	"context"

	v3cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type RoleUsers struct {
	users         map[string][]RoleUser
	orphanedUsers map[string]string
}
type RoleUser struct {
	UserName string
	GUID     string
	Origin   string
}

type Manager interface {
	ClearRoles()
	DeleteUser(userGuid string) error
	ListSpaceUsersByRole(spaceGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error)
	ListOrgUsersByRole(orgGUID string) (*RoleUsers, *RoleUsers, *RoleUsers, *RoleUsers, error)
	AssociateOrgAuditor(orgGUID, orgName, entityGUID, userName, userGUID string) error
	AssociateOrgManager(orgGUID, orgName, entityGUID, userName, userGUID string) error
	AssociateOrgBillingManager(orgGUID, orgName, entityGUID, userName, userGUID string) error

	RemoveOrgManager(orgName, orgGUID, userName, userGUID string) error
	RemoveOrgBillingManager(orgName, orgGUID, userName, userGUID string) error
	RemoveOrgAuditor(orgName, orgGUID, userName, userGUID string) error
	RemoveOrgUser(orgName, orgGUID, userName, userGUID string) error

	AssociateSpaceAuditor(orgGUID, spaceName, spaceGUID, userName, userGUID string) error
	AssociateSpaceManager(orgGUID, spaceName, spaceGUID, userName, userGUID string) error
	AssociateSpaceDeveloper(orgGUID, spaceName, spaceGUID, userName, userGUID string) error
	AssociateSpaceSupporter(orgGUID, spaceName, spaceGUID, userName, userGUID string) error

	RemoveSpaceAuditor(spaceName, spaceGUID, userName, userGUID string) error
	RemoveSpaceDeveloper(spaceName, spaceGUID, userName, userGUID string) error
	RemoveSpaceManager(spaceName, spaceGUID, userName, userGUID string) error
	RemoveSpaceSupporter(spaceName, spaceGUID, userName, userGUID string) error
}

type CFRoleClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.RoleListOptions) ([]*resource.Role, error)
	CreateOrganizationRole(ctx context.Context, organizationGUID, userGUID string, roleType resource.OrganizationRoleType) (*resource.Role, error)
	CreateSpaceRole(ctx context.Context, spaceGUID, userGUID string, roleType resource.SpaceRoleType) (*resource.Role, error)
	Delete(ctx context.Context, guid string) (string, error)
}

type CFUserClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.UserListOptions) ([]*resource.User, error)
	Delete(ctx context.Context, guid string) (string, error)
}

type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *v3cfclient.PollingOptions) error
}
