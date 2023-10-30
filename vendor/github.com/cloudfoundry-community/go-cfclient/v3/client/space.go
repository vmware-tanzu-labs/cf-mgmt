package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SpaceClient commonClient

// SpaceListOptions list filters
type SpaceListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`              // list of space guids to filter by
	Names             Filter `qs:"names"`              // list of space names to filter by
	OrganizationGUIDs Filter `qs:"organization_guids"` // list of organization guids to filter by

	Include resource.SpaceIncludeType `qs:"include"` // include parent objects if any
}

// NewSpaceListOptions creates new options to pass to list
func NewSpaceListOptions() *SpaceListOptions {
	return &SpaceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SpaceListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// AssignIsolationSegment assigns an isolation segment to the space
//
// Apps will not run in the isolation segment until they are restarted
// An empty isolationSegmentGUID will un-assign the isolation segment
func (c *SpaceClient) AssignIsolationSegment(ctx context.Context, guid, isolationSegmentGUID string) error {
	r := &resource.NullableToOneRelationship{
		Data: &resource.NullableRelationship{
			GUID: &isolationSegmentGUID,
		},
	}
	if isolationSegmentGUID == "" {
		r.Data.GUID = nil // set data to null to remove the relationship
	}
	_, err := c.client.patch(ctx, path.Format("/v3/spaces/%s/relationships/isolation_segment", guid), r, nil)
	return err
}

// Create a new space
func (c *SpaceClient) Create(ctx context.Context, r *resource.SpaceCreate) (*resource.Space, error) {
	var space resource.Space
	_, err := c.client.post(ctx, "/v3/spaces", r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// Delete the specified space asynchronously and return a jobGUID
func (c *SpaceClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/spaces/%s", guid))
}

// First returns the first space matching the options or an error when less than 1 match
func (c *SpaceClient) First(ctx context.Context, opts *SpaceListOptions) (*resource.Space, error) {
	return First[*SpaceListOptions, *resource.Space](opts, func(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified space
func (c *SpaceClient) Get(ctx context.Context, guid string) (*resource.Space, error) {
	var space resource.Space
	err := c.client.get(ctx, path.Format("/v3/spaces/%s", guid), &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}

// GetAssignedIsolationSegment gets the space's assigned isolation segment, if any
func (c *SpaceClient) GetAssignedIsolationSegment(ctx context.Context, guid string) (string, error) {
	var relation resource.ToOneRelationship
	err := c.client.get(ctx, path.Format("/v3/spaces/%s/relationships/isolation_segment", guid), &relation)
	if err != nil {
		return "", err
	}
	if relation.Data == nil {
		return "", nil
	}
	return relation.Data.GUID, nil
}

// GetIncludeOrganization allows callers to fetch a space and include the parent organization
func (c *SpaceClient) GetIncludeOrganization(ctx context.Context, guid string) (*resource.Space, *resource.Organization, error) {
	var space resource.SpaceWithIncluded
	err := c.client.get(ctx, path.Format("/v3/spaces/%s?include=%s", guid, resource.SpaceIncludeOrganization), &space)
	if err != nil {
		return nil, nil, err
	}
	return &space.Space, space.Included.Organizations[0], nil
}

// List pages all spaces the user has access to
func (c *SpaceClient) List(ctx context.Context, opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	opts.Include = resource.SpaceIncludeNone

	var res resource.SpaceList
	err := c.client.list(ctx, "/v3/spaces", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all spaces the user has access to
func (c *SpaceClient) ListAll(ctx context.Context, opts *SpaceListOptions) ([]*resource.Space, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	return AutoPage[*SpaceListOptions, *resource.Space](opts, func(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeOrganizations page all spaces the user has access to and include the parent organizations
func (c *SpaceClient) ListIncludeOrganizations(ctx context.Context, opts *SpaceListOptions) ([]*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}
	opts.Include = resource.SpaceIncludeOrganization

	var res resource.SpaceList
	err := c.client.list(ctx, "/v3/spaces", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Organizations, pager, nil
}

// ListIncludeOrganizationsAll retrieves all spaces the user has access to and include the parent organizations
func (c *SpaceClient) ListIncludeOrganizationsAll(ctx context.Context, opts *SpaceListOptions) ([]*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewSpaceListOptions()
	}

	var all []*resource.Space
	var allOrgs []*resource.Organization
	for {
		page, orgs, pager, err := c.ListIncludeOrganizations(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allOrgs = append(allOrgs, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allOrgs, nil
}

// ListUsers pages users by space GUID
func (c *SpaceClient) ListUsers(ctx context.Context, spaceGUID string, opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	var res resource.UserList
	err := c.client.list(ctx, "/v3/spaces/"+spaceGUID+"/users", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListUsersAll retrieves all users by space GUID
func (c *SpaceClient) ListUsersAll(ctx context.Context, spaceGUID string, opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.ListUsers(ctx, spaceGUID, opts)
	})
}

// Single returns a single space matching the options or an error if not exactly 1 match
func (c *SpaceClient) Single(ctx context.Context, opts *SpaceListOptions) (*resource.Space, error) {
	return Single[*SpaceListOptions, *resource.Space](opts, func(opts *SpaceListOptions) ([]*resource.Space, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of a space
func (c *SpaceClient) Update(ctx context.Context, guid string, r *resource.SpaceUpdate) (*resource.Space, error) {
	var space resource.Space
	_, err := c.client.patch(ctx, path.Format("/v3/spaces/%s", guid), r, &space)
	if err != nil {
		return nil, err
	}
	return &space, nil
}
