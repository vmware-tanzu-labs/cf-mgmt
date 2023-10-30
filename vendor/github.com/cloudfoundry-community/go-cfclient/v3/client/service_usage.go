package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceUsageClient commonClient

// ServiceUsageListOptions list filters
type ServiceUsageListOptions struct {
	*ListOptions

	AfterGUID            string `qs:"after_guid"`
	GUIDs                Filter `qs:"guids"`
	ServiceInstanceTypes Filter `qs:"service_instance_types"`
	ServiceOfferingGUIDs Filter `qs:"service_offering_guids"`
}

// NewServiceUsageOptions creates new options to pass to list
func NewServiceUsageOptions() *ServiceUsageListOptions {
	return &ServiceUsageListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceUsageListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// First returns the first space matching the options or an error when less than 1 match
func (c *ServiceUsageClient) First(ctx context.Context, opts *ServiceUsageListOptions) (*resource.ServiceUsage, error) {
	return First[*ServiceUsageListOptions, *resource.ServiceUsage](opts, func(opts *ServiceUsageListOptions) ([]*resource.ServiceUsage, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get retrieves the specified service event
func (c *ServiceUsageClient) Get(ctx context.Context, guid string) (*resource.ServiceUsage, error) {
	var a resource.ServiceUsage
	err := c.client.get(ctx, path.Format("/v3/service_usage_events/%s", guid), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// List pages all service usage events
func (c *ServiceUsageClient) List(ctx context.Context, opts *ServiceUsageListOptions) ([]*resource.ServiceUsage, *Pager, error) {
	if opts == nil {
		opts = NewServiceUsageOptions()
	}
	var res resource.ServiceUsageList
	err := c.client.list(ctx, "/v3/service_usage_events", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service usage events
func (c *ServiceUsageClient) ListAll(ctx context.Context, opts *ServiceUsageListOptions) ([]*resource.ServiceUsage, error) {
	if opts == nil {
		opts = NewServiceUsageOptions()
	}
	return AutoPage[*ServiceUsageListOptions, *resource.ServiceUsage](opts, func(opts *ServiceUsageListOptions) ([]*resource.ServiceUsage, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Purge destroys all existing events. Populates new usage events, one for each existing service instance.
// All populated events will have a created_at value of current time.
//
// There is the potential race condition if service instances are currently being created or deleted.
// The seeded usage events will have the same guid as the service instance.
func (c *ServiceUsageClient) Purge(ctx context.Context) error {
	_, err := c.client.post(ctx, "/v3/service_usage_events/actions/destructively_purge_all_and_reseed", nil, nil)
	return err
}

// Single returns a single service usage matching the options or an error if not exactly 1 match
func (c *ServiceUsageClient) Single(ctx context.Context, opts *ServiceUsageListOptions) (*resource.ServiceUsage, error) {
	return Single[*ServiceUsageListOptions, *resource.ServiceUsage](opts, func(opts *ServiceUsageListOptions) ([]*resource.ServiceUsage, *Pager, error) {
		return c.List(ctx, opts)
	})
}
