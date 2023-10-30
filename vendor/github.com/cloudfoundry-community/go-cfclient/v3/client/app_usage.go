package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type AppUsageClient commonClient

// AppUsageListOptions list filters
type AppUsageListOptions struct {
	*ListOptions
}

// NewAppUsageOptions creates new options to pass to list
func NewAppUsageOptions() *AppUsageListOptions {
	return &AppUsageListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o AppUsageListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Get retrieves the specified app event
func (c *AppUsageClient) Get(ctx context.Context, guid string) (*resource.AppUsage, error) {
	var a resource.AppUsage
	err := c.client.get(ctx, path.Format("/v3/app_usage_events/%s", guid), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// List pages all app usage events
func (c *AppUsageClient) List(ctx context.Context, opts *AppUsageListOptions) ([]*resource.AppUsage, *Pager, error) {
	if opts == nil {
		opts = NewAppUsageOptions()
	}
	var res resource.AppUsageList
	err := c.client.list(ctx, "/v3/app_usage_events", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all app usage events
func (c *AppUsageClient) ListAll(ctx context.Context, opts *AppUsageListOptions) ([]*resource.AppUsage, error) {
	if opts == nil {
		opts = NewAppUsageOptions()
	}
	return AutoPage[*AppUsageListOptions, *resource.AppUsage](opts, func(opts *AppUsageListOptions) ([]*resource.AppUsage, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Purge destroys all existing events. Populates new usage events, one for each started app.
// All populated events will have a created_at value of current time.
//
// There is the potential race condition if apps are currently being started, stopped, or scaled.
// The seeded usage events will have the same guid as the app.
func (c *AppUsageClient) Purge(ctx context.Context) error {
	_, err := c.client.post(ctx, "/v3/app_usage_events/actions/destructively_purge_all_and_reseed", nil, nil)
	return err
}
