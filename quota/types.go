package quota

import (
	"context"

	v3cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type CFSpaceQuotaClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.SpaceQuotaListOptions) ([]*resource.SpaceQuota, error)
	Update(ctx context.Context, guid string, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error)
	Create(ctx context.Context, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error)
	Apply(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error)
	Get(ctx context.Context, guid string) (*resource.SpaceQuota, error)
}

type CFOrgQuotaClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, error)
	Update(ctx context.Context, guid string, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error)
	Create(ctx context.Context, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error)
	Get(ctx context.Context, guid string) (*resource.OrganizationQuota, error)
	Apply(ctx context.Context, guid string, organizationGUIDs []string) ([]string, error)
}
