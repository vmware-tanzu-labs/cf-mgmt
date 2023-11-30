package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type OrganizationClient commonClient

type OrganizationListOptions struct {
	*ListOptions

	GUIDs Filter `qs:"guids"` // list of organization guids to filter by
	Names Filter `qs:"names"` // list of organization names to filter by
}

// NewOrganizationListOptions creates new options to pass to list
func NewOrganizationListOptions() *OrganizationListOptions {
	return &OrganizationListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o OrganizationListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// AssignDefaultIsolationSegment assigns a default iso segment to the specified organization
//
// Apps will not run in the new default isolation segment until they are restarted
// An empty isolationSegmentGUID will un-assign the default isolation segment
func (c *OrganizationClient) AssignDefaultIsolationSegment(ctx context.Context, guid, isolationSegmentGUID string) error {
	r := &resource.NullableToOneRelationship{
		Data: &resource.NullableRelationship{
			GUID: &isolationSegmentGUID,
		},
	}
	if isolationSegmentGUID == "" {
		r.Data.GUID = nil // set data to null to remove the relationship
	}
	_, err := c.client.patch(ctx, path.Format("/v3/organizations/%s/relationships/default_isolation_segment", guid), r, nil)
	return err
}

// Create an organization
func (c *OrganizationClient) Create(ctx context.Context, r *resource.OrganizationCreate) (*resource.Organization, error) {
	var org resource.Organization
	_, err := c.client.post(ctx, "/v3/organizations", r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// Delete the specified organization asynchronously and return a jobGUID
func (c *OrganizationClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/organizations/%s", guid))
}

// First returns the first organization matching the options or an error when less than 1 match
func (c *OrganizationClient) First(ctx context.Context, opts *OrganizationListOptions) (*resource.Organization, error) {
	return First[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForIsolationSegment returns the first organization matching the options and iso segment or an error when less than 1 match
func (c *OrganizationClient) FirstForIsolationSegment(ctx context.Context, isolationSegmentGUID string, opts *OrganizationListOptions) (*resource.Organization, error) {
	return First[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsolationSegment(ctx, isolationSegmentGUID, opts)
	})
}

// Get the specified organization
func (c *OrganizationClient) Get(ctx context.Context, guid string) (*resource.Organization, error) {
	var org resource.Organization
	err := c.client.get(ctx, path.Format("/v3/organizations/%s", guid), &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

// GetDefaultIsolationSegment gets the specified organization's default iso segment GUID if any
func (c *OrganizationClient) GetDefaultIsolationSegment(ctx context.Context, guid string) (string, error) {
	var relation resource.ToOneRelationship
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/relationships/default_isolation_segment", guid), &relation)
	if err != nil {
		return "", err
	}
	if relation.Data == nil {
		return "", nil
	}
	return relation.Data.GUID, nil
}

// GetDefaultDomain gets the specified organization's default domain if any
func (c *OrganizationClient) GetDefaultDomain(ctx context.Context, guid string) (*resource.Domain, error) {
	var domain resource.Domain
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/domains/default", guid), &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetUsageSummary gets the specified organization's usage summary
func (c *OrganizationClient) GetUsageSummary(ctx context.Context, guid string) (*resource.OrganizationUsageSummary, error) {
	var summary resource.OrganizationUsageSummary
	err := c.client.get(ctx, path.Format("/v3/organizations/%s/usage_summary", guid), &summary)
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

// List pages all organizations the user has access to
func (c *OrganizationClient) List(ctx context.Context, opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewOrganizationListOptions()
	}
	var res resource.OrganizationList
	err := c.client.list(ctx, "/v3/organizations", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all organizations the user has access to
func (c *OrganizationClient) ListAll(ctx context.Context, opts *OrganizationListOptions) ([]*resource.Organization, error) {
	if opts == nil {
		opts = NewOrganizationListOptions()
	}
	return AutoPage[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForIsolationSegment pages all organizations for the specified isolation segment
func (c *OrganizationClient) ListForIsolationSegment(ctx context.Context, isolationSegmentGUID string, opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewOrganizationListOptions()
	}
	var res resource.OrganizationList
	err := c.client.list(ctx, "/v3/isolation_segments/"+isolationSegmentGUID+"/organizations", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForIsolationSegmentAll retrieves all organizations for the specified isolation segment
func (c *OrganizationClient) ListForIsolationSegmentAll(ctx context.Context, isolationSegmentGUID string, opts *OrganizationListOptions) ([]*resource.Organization, error) {
	if opts == nil {
		opts = NewOrganizationListOptions()
	}
	return AutoPage[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsolationSegment(ctx, isolationSegmentGUID, opts)
	})
}

// ListUsers pages of all users that are members of the specified organization
func (c *OrganizationClient) ListUsers(ctx context.Context, guid string, opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	var res resource.UserList
	err := c.client.list(ctx, "/v3/organizations/"+guid+"/users", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListUsersAll retrieves all users that are members of the specified organization
func (c *OrganizationClient) ListUsersAll(ctx context.Context, guid string, opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.ListUsers(ctx, guid, opts)
	})
}

// Single returns a single organization matching the options or an error if not exactly 1 match
func (c *OrganizationClient) Single(ctx context.Context, opts *OrganizationListOptions) (*resource.Organization, error) {
	return Single[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForIsolationSegment returns a single organization matching the options and iso segment or an error if not exactly 1 match
func (c *OrganizationClient) SingleForIsolationSegment(ctx context.Context, isolationSegmentGUID string, opts *OrganizationListOptions) (*resource.Organization, error) {
	return Single[*OrganizationListOptions, *resource.Organization](opts, func(opts *OrganizationListOptions) ([]*resource.Organization, *Pager, error) {
		return c.ListForIsolationSegment(ctx, isolationSegmentGUID, opts)
	})
}

// Update the organization's specified attributes
func (c *OrganizationClient) Update(ctx context.Context, guid string, r *resource.OrganizationUpdate) (*resource.Organization, error) {
	var org resource.Organization
	_, err := c.client.patch(ctx, path.Format("/v3/organizations/%s", guid), r, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}
