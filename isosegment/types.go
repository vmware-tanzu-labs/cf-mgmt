package isosegment

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type Manager interface {
	GetIsoSegmentForOrg(orgGUID string) (*resource.IsolationSegment, error)
	GetIsoSegmentForSpace(spaceGUID string) (*resource.IsolationSegment, error)
	Apply() error
	Create() error
	Remove() error
	Entitle() error
	Unentitle() error
	UpdateOrgs() error
	UpdateSpaces() error
	ListIsolationSegments() ([]*resource.IsolationSegment, error)
}

type CFIsolationSegmentClient interface {
	ListAll(ctx context.Context, opts *cfclient.IsolationSegmentListOptions) ([]*resource.IsolationSegment, error)
	Create(ctx context.Context, r *resource.IsolationSegmentCreate) (*resource.IsolationSegment, error)
	Delete(ctx context.Context, guid string) error
	Get(ctx context.Context, guid string) (*resource.IsolationSegment, error)
	EntitleOrganization(ctx context.Context, guid string, organizationGUID string) (*resource.IsolationSegmentRelationship, error)
	RevokeOrganization(ctx context.Context, guid string, organizationGUID string) error
}

type CFOrganizationClient interface {
	AssignDefaultIsolationSegment(ctx context.Context, guid, isoSegmentGUID string) error
	GetDefaultIsolationSegment(ctx context.Context, guid string) (string, error)
}

type CFSpaceClient interface {
	AssignIsolationSegment(ctx context.Context, guid, isolationSegmentGUID string) error
	GetAssignedIsolationSegment(ctx context.Context, guid string) (string, error)
}
