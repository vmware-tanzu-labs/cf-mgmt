package space

import (
	"context"

	v3cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
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
	IsSSHEnabled(*resource.Space) (bool, error)
	UpdateSpacesMetadata() error
	GetSpaceIsolationSegmentGUID(*resource.Space) (string, error)
}

type CFSpaceClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.SpaceListOptions) ([]*resource.Space, error)
	Create(ctx context.Context, r *resource.SpaceCreate) (*resource.Space, error)
	Update(ctx context.Context, guid string, r *resource.SpaceUpdate) (*resource.Space, error)
	Delete(ctx context.Context, guid string) (string, error)
	GetAssignedIsolationSegment(ctx context.Context, guid string) (string, error)
}

type CFSpaceFeatureClient interface {
	IsSSHEnabled(ctx context.Context, spaceGUID string) (bool, error)
	EnableSSH(ctx context.Context, spaceGUID string, enable bool) error
}
