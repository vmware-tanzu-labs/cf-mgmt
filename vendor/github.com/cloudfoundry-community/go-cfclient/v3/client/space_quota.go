package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SpaceQuotaClient commonClient

// SpaceQuotaListOptions list filters
type SpaceQuotaListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
	SpaceGUIDs        Filter `qs:"space_guids"`
}

// NewSpaceQuotaListOptions creates new options to pass to list
func NewSpaceQuotaListOptions() *SpaceQuotaListOptions {
	return &SpaceQuotaListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SpaceQuotaListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Apply the quota to the specified spaces
func (c *SpaceQuotaClient) Apply(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/space_quotas/%s/relationships/spaces", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// Create a new space quota
func (c *SpaceQuotaClient) Create(ctx context.Context, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	_, err := c.client.post(ctx, "/v3/space_quotas", r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Delete the specified space quota asynchronously and return a jobGUID
func (c *SpaceQuotaClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/space_quotas/%s", guid))
}

// First returns the first space quota matching the options or an error when less than 1 match
func (c *SpaceQuotaClient) First(ctx context.Context, opts *SpaceQuotaListOptions) (*resource.SpaceQuota, error) {
	return First[*SpaceQuotaListOptions, *resource.SpaceQuota](opts, func(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified space quota
func (c *SpaceQuotaClient) Get(ctx context.Context, guid string) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	err := c.client.get(ctx, path.Format("/v3/space_quotas/%s", guid), &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// List pages all space quotas the user has access to
func (c *SpaceQuotaClient) List(ctx context.Context, opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
	if opts == nil {
		opts = NewSpaceQuotaListOptions()
	}

	var res resource.SpaceQuotaList
	err := c.client.list(ctx, "/v3/space_quotas", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all space quotas the user has access to
func (c *SpaceQuotaClient) ListAll(ctx context.Context, opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, error) {
	if opts == nil {
		opts = NewSpaceQuotaListOptions()
	}
	return AutoPage[*SpaceQuotaListOptions, *resource.SpaceQuota](opts, func(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Remove the space quota from the specified space
func (c *SpaceQuotaClient) Remove(ctx context.Context, guid, spaceGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/space_quotas/%s/relationships/spaces/%s", guid, spaceGUID))
	return err
}

// Single returns a single space quota matching the options or an error if not exactly 1 match
func (c *SpaceQuotaClient) Single(ctx context.Context, opts *SpaceQuotaListOptions) (*resource.SpaceQuota, error) {
	return Single[*SpaceQuotaListOptions, *resource.SpaceQuota](opts, func(opts *SpaceQuotaListOptions) ([]*resource.SpaceQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the organization quota
func (c *SpaceQuotaClient) Update(ctx context.Context, guid string, r *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	var q resource.SpaceQuota
	_, err := c.client.patch(ctx, path.Format("/v3/space_quotas/%s", guid), r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
