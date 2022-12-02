package quota

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type CFSpaceQuotaClient interface {
	ListAll(ctx context.Context, opts *client.SpaceQuotaListOptions) ([]*resource.SpaceQuota, error)
	Update(ctx context.Context, guid string, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error)
	Apply(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error)
	Create(ctx context.Context, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error)
	Single(ctx context.Context, opts *client.SpaceQuotaListOptions) (*resource.SpaceQuota, error)
	Get(ctx context.Context, guid string) (*resource.SpaceQuota, error)
}

type CFOrganizationQuotaClient interface {
	ListAll(ctx context.Context, opts *client.OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, error)
	Update(ctx context.Context, guid string, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error)
	Create(ctx context.Context, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error)
	Apply(ctx context.Context, guid string, organizationGUIDs []string) ([]string, error)
	Single(ctx context.Context, opts *client.OrganizationQuotaListOptions) (*resource.OrganizationQuota, error)
	Get(ctx context.Context, guid string) (*resource.OrganizationQuota, error)
}
