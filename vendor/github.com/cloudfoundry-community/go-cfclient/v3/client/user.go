package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type UserClient commonClient

// UserListOptions list filters
type UserListOptions struct {
	*ListOptions

	// list of user guids to filter by
	GUIDs Filter `qs:"guids"`

	// list of usernames to filter by. Mutually exclusive with partial_usernames
	UserNames Filter `qs:"usernames"`

	// list of strings to search by. When using this query parameter, all the users that
	// contain the string provided in their username will be returned. Mutually exclusive with usernames
	PartialUsernames Filter `qs:"partial_usernames"`

	// list of user origins (user stores) to filter by, for example, users authenticated by
	// UAA have the origin “uaa”; users authenticated by an LDAP provider have the
	// origin ldap when filtering by origins, usernames must be included
	Origins Filter `qs:"origins"`
}

// NewUserListOptions creates new options to pass to list
func NewUserListOptions() *UserListOptions {
	return &UserListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o UserListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new user
func (c *UserClient) Create(ctx context.Context, r *resource.UserCreate) (*resource.User, error) {
	var user resource.User
	_, err := c.client.post(ctx, "/v3/users", r, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Delete the specified user
func (c *UserClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/users/%s", guid))
}

// First returns the first user matching the options or an error when less than 1 match
func (c *UserClient) First(ctx context.Context, opts *UserListOptions) (*resource.User, error) {
	return First[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified user
func (c *UserClient) Get(ctx context.Context, guid string) (*resource.User, error) {
	var user resource.User
	err := c.client.get(ctx, path.Format("/v3/users/%s", guid), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List pages all users the user has access to
func (c *UserClient) List(ctx context.Context, opts *UserListOptions) ([]*resource.User, *Pager, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	var res resource.UserList
	err := c.client.list(ctx, "/v3/users", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all users the user has access to
func (c *UserClient) ListAll(ctx context.Context, opts *UserListOptions) ([]*resource.User, error) {
	if opts == nil {
		opts = NewUserListOptions()
	}
	return AutoPage[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single user matching the options or an error if not exactly 1 match
func (c *UserClient) Single(ctx context.Context, opts *UserListOptions) (*resource.User, error) {
	return Single[*UserListOptions, *resource.User](opts, func(opts *UserListOptions) ([]*resource.User, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of a user
func (c *UserClient) Update(ctx context.Context, guid string, r *resource.UserUpdate) (*resource.User, error) {
	var user resource.User
	_, err := c.client.patch(ctx, path.Format("/v3/users/%s", guid), r, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
