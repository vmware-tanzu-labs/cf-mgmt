package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type StackClient commonClient

// StackListOptions list filters
type StackListOptions struct {
	*ListOptions

	Names Filter `qs:"names"` // list of stack names to filter by
}

// NewStackListOptions creates new options to pass to list
func NewStackListOptions() *StackListOptions {
	return &StackListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o StackListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new stack
func (c *StackClient) Create(ctx context.Context, r *resource.StackCreate) (*resource.Stack, error) {
	var stack resource.Stack
	_, err := c.client.post(ctx, "/v3/stacks", r, &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

// Delete the specified stack
func (c *StackClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/stacks/%s", guid))
	return err
}

// First returns the first stack matching the options or an error when less than 1 match
func (c *StackClient) First(ctx context.Context, opts *StackListOptions) (*resource.Stack, error) {
	return First[*StackListOptions, *resource.Stack](opts, func(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified stack
func (c *StackClient) Get(ctx context.Context, guid string) (*resource.Stack, error) {
	var stack resource.Stack
	err := c.client.get(ctx, path.Format("/v3/stacks/%s", guid), &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

// List pages all stacks the user has access to
func (c *StackClient) List(ctx context.Context, opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	var res resource.StackList
	err := c.client.list(ctx, "/v3/stacks", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all stacks the user has access to
func (c *StackClient) ListAll(ctx context.Context, opts *StackListOptions) ([]*resource.Stack, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	return AutoPage[*StackListOptions, *resource.Stack](opts, func(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListAppsOnStack pages all apps using a given stack
func (c *StackClient) ListAppsOnStack(ctx context.Context, guid string, opts *StackListOptions) ([]*resource.App, *Pager, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	var res resource.AppList
	err := c.client.list(ctx, "/v3/stacks/"+guid+"/apps", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAppsOnStackAll retrieves all apps using a given stack
func (c *StackClient) ListAppsOnStackAll(ctx context.Context, guid string, opts *StackListOptions) ([]*resource.App, error) {
	if opts == nil {
		opts = NewStackListOptions()
	}
	return AutoPage[*StackListOptions, *resource.App](opts, func(opts *StackListOptions) ([]*resource.App, *Pager, error) {
		return c.ListAppsOnStack(ctx, guid, opts)
	})
}

// Single returns a single stack matching the options or an error if not exactly 1 match
func (c *StackClient) Single(ctx context.Context, opts *StackListOptions) (*resource.Stack, error) {
	return Single[*StackListOptions, *resource.Stack](opts, func(opts *StackListOptions) ([]*resource.Stack, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of a stack
func (c *StackClient) Update(ctx context.Context, guid string, r *resource.StackUpdate) (*resource.Stack, error) {
	var stack resource.Stack
	_, err := c.client.patch(ctx, path.Format("/v3/stacks/%s", guid), r, &stack)
	if err != nil {
		return nil, err
	}
	return &stack, nil
}
