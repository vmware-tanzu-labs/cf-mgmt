package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type OrganizationQuotaClient commonClient

// OrganizationQuotaListOptions list filters
type OrganizationQuotaListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewOrganizationQuotaListOptions creates new options to pass to list
func NewOrganizationQuotaListOptions() *OrganizationQuotaListOptions {
	return &OrganizationQuotaListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o OrganizationQuotaListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Apply the specified organization quota to the organizations
func (c *OrganizationQuotaClient) Apply(ctx context.Context, guid string, organizationGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(organizationGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/organization_quotas/%s/relationships/organizations", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// Create a new organization quota
func (c *OrganizationQuotaClient) Create(ctx context.Context, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	_, err := c.client.post(ctx, "/v3/organization_quotas", r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// Delete the specified organization quota
func (c *OrganizationQuotaClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/organization_quotas/%s", guid))
	return err
}

// First returns the first organization quota matching the options or an error when less than 1 match
func (c *OrganizationQuotaClient) First(ctx context.Context, opts *OrganizationQuotaListOptions) (*resource.OrganizationQuota, error) {
	return First[*OrganizationQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified organization quota
func (c *OrganizationQuotaClient) Get(ctx context.Context, guid string) (*resource.OrganizationQuota, error) {
	var app resource.OrganizationQuota
	err := c.client.get(ctx, path.Format("/v3/organization_quotas/%s", guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// List pages all organization quotas the user has access to
func (c *OrganizationQuotaClient) List(ctx context.Context, opts *OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
	if opts == nil {
		opts = NewOrganizationQuotaListOptions()
	}

	var res resource.OrganizationQuotaList
	err := c.client.list(ctx, "/v3/organization_quotas", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all organization quotas the user has access to
func (c *OrganizationQuotaClient) ListAll(ctx context.Context, opts *OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, error) {
	if opts == nil {
		opts = NewOrganizationQuotaListOptions()
	}
	return AutoPage[*OrganizationQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single organization quota matching the options or an error if not exactly 1 match
func (c *OrganizationQuotaClient) Single(ctx context.Context, opts *OrganizationQuotaListOptions) (*resource.OrganizationQuota, error) {
	return Single[*OrganizationQuotaListOptions, *resource.OrganizationQuota](opts, func(opts *OrganizationQuotaListOptions) ([]*resource.OrganizationQuota, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the organization quota
func (c *OrganizationQuotaClient) Update(ctx context.Context, guid string, r *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	var q resource.OrganizationQuota
	_, err := c.client.patch(ctx, path.Format("/v3/organization_quotas/%s", guid), r, &q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
