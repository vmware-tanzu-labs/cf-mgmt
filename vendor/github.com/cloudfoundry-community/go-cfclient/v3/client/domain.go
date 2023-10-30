package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type DomainClient commonClient

// DomainListOptions list filters
type DomainListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewDomainListOptions creates new options to pass to list
func NewDomainListOptions() *DomainListOptions {
	return &DomainListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DomainListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new domain
func (c *DomainClient) Create(ctx context.Context, r *resource.DomainCreate) (*resource.Domain, error) {
	var d resource.Domain
	_, err := c.client.post(ctx, "/v3/domains", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified domain asynchronously and return a jobGUID.
func (c *DomainClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/domains/%s", guid))
}

// First returns the first domain matching the options or an error when less than 1 match
func (c *DomainClient) First(ctx context.Context, opts *DomainListOptions) (*resource.Domain, error) {
	return First[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForOrganization returns the first domain matching the options and organization or an error when less than 1 match
func (c *DomainClient) FirstForOrganization(ctx context.Context, organizationGUID string, opts *DomainListOptions) (*resource.Domain, error) {
	return First[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.ListForOrganization(ctx, organizationGUID, opts)
	})
}

// Get the specified domain
func (c *DomainClient) Get(ctx context.Context, guid string) (*resource.Domain, error) {
	var d resource.Domain
	err := c.client.get(ctx, path.Format("/v3/domains/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages Domains the user has access to
func (c *DomainClient) List(ctx context.Context, opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
	var res resource.DomainList
	err := c.client.list(ctx, "/v3/domains", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all domains the user has access to
func (c *DomainClient) ListAll(ctx context.Context, opts *DomainListOptions) ([]*resource.Domain, error) {
	if opts == nil {
		opts = NewDomainListOptions()
	}
	return AutoPage[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForOrganization pages all domains for the specified org that the user has access to
func (c *DomainClient) ListForOrganization(ctx context.Context, organizationGUID string, opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
	if opts == nil {
		opts = NewDomainListOptions()
	}
	var res resource.DomainList
	err := c.client.list(ctx, "/v3/organizations/"+organizationGUID+"/domains", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForOrganizationAll retrieves all domains for the specified org that the user has access to
func (c *DomainClient) ListForOrganizationAll(ctx context.Context, organizationGUID string, opts *DomainListOptions) ([]*resource.Domain, error) {
	if opts == nil {
		opts = NewDomainListOptions()
	}
	return AutoPage[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.ListForOrganization(ctx, organizationGUID, opts)
	})
}

// Share an organization-scoped domain to the organization specified by the org guid
// This will allow the organization to use the organization-scoped domain
func (c *DomainClient) Share(ctx context.Context, domainGUID, organizationGUID string) (*resource.ToManyRelationships, error) {
	r := resource.NewDomainShare(organizationGUID)
	return c.ShareMany(ctx, domainGUID, r)
}

// ShareMany shares an organization-scoped domain to other organizations specified by a list of organization guids
// This will allow any of the other organizations to use the organization-scoped domain.
func (c *DomainClient) ShareMany(ctx context.Context, guid string, r *resource.ToManyRelationships) (*resource.ToManyRelationships, error) {
	var d resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/domains/%s/relationships/shared_organizations", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Single returns a single domain matching the options or an error if not exactly 1 match
func (c *DomainClient) Single(ctx context.Context, opts *DomainListOptions) (*resource.Domain, error) {
	return Single[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForOrganization returns a single domain matching the options and org or an error if not exactly 1 match
func (c *DomainClient) SingleForOrganization(ctx context.Context, organizationGUID string, opts *DomainListOptions) (*resource.Domain, error) {
	return Single[*DomainListOptions, *resource.Domain](opts, func(opts *DomainListOptions) ([]*resource.Domain, *Pager, error) {
		return c.ListForOrganization(ctx, organizationGUID, opts)
	})
}

// UnShare an organization-scoped domain to other organizations specified by a list of organization guids
// This will allow any of the other organizations to use the organization-scoped domain.
func (c *DomainClient) UnShare(ctx context.Context, domainGUID, organizationGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/domains/%s/relationships/shared_organizations/%s", domainGUID, organizationGUID))
	return err
}

// Update the specified attributes of the domain
func (c *DomainClient) Update(ctx context.Context, guid string, r *resource.DomainUpdate) (*resource.Domain, error) {
	var d resource.Domain
	_, err := c.client.patch(ctx, path.Format("/v3/domains/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
