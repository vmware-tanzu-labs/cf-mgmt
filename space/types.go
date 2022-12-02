package space

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (*resource.Space, error)
	CreateSpaces() error
	UpdateSpaces() (err error)
	DeleteSpaces() (err error)
	DeleteSpacesForOrg(orgGUID, orgName string) (err error)
	ListSpaces(orgGUID string) ([]*resource.Space, error)
	UpdateSpacesMetadata() error
	IsSSHEnabled(spaceGUID string) (bool, error)
}

type CFSpaceClient interface {
	Get(ctx context.Context, guid string) (*resource.Space, error)
	Update(ctx context.Context, guid string, r *resource.SpaceUpdate) (*resource.Space, error)
	Create(ctx context.Context, r *resource.SpaceCreate) (*resource.Space, error)
	Delete(ctx context.Context, guid string) (string, error)
	ListAll(ctx context.Context, opts *cfclient.SpaceListOptions) ([]*resource.Space, error)
}

type CFOrganizationClient interface {
	ListAll(ctx context.Context, opts *cfclient.OrganizationListOptions) ([]*resource.Organization, error)
}

type CFSpaceFeatureClient interface {
	EnableSSH(ctx context.Context, spaceGUID string, enable bool) error
	IsSSHEnabled(ctx context.Context, spaceGUID string) (bool, error)
}

type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *cfclient.PollingOptions) error
}
