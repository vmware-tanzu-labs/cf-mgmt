package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceBrokerClient commonClient

// ServiceBrokerListOptions list filters
type ServiceBrokerListOptions struct {
	*ListOptions

	SpaceGUIDs Filter `qs:"space_guids"`
	Names      Filter `qs:"names"`
}

// NewServiceBrokerListOptions creates new options to pass to list
func NewServiceBrokerListOptions() *ServiceBrokerListOptions {
	return &ServiceBrokerListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceBrokerListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service broker asynchronously and return a jobGUID
func (c *ServiceBrokerClient) Create(ctx context.Context, r *resource.ServiceBrokerCreate) (string, error) {
	return c.client.post(ctx, "/v3/service_brokers", r, nil)
}

// Delete the specified service broker asynchronously and return a jobGUID
func (c *ServiceBrokerClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/service_brokers/%s", guid))
}

// First returns the first service broker matching the options or an error when less than 1 match
func (c *ServiceBrokerClient) First(ctx context.Context, opts *ServiceBrokerListOptions) (*resource.ServiceBroker, error) {
	return First[*ServiceBrokerListOptions, *resource.ServiceBroker](opts, func(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service broker
func (c *ServiceBrokerClient) Get(ctx context.Context, guid string) (*resource.ServiceBroker, error) {
	var sb resource.ServiceBroker
	err := c.client.get(ctx, path.Format("/v3/service_brokers/%s", guid), &sb)
	if err != nil {
		return nil, err
	}
	return &sb, nil
}

// List pages all the service brokers the user has access to
func (c *ServiceBrokerClient) List(ctx context.Context, opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
	if opts == nil {
		opts = NewServiceBrokerListOptions()
	}

	var res resource.ServiceBrokerList
	err := c.client.list(ctx, "/v3/service_brokers", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service brokers the user has access to
func (c *ServiceBrokerClient) ListAll(ctx context.Context, opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, error) {
	if opts == nil {
		opts = NewServiceBrokerListOptions()
	}
	return AutoPage[*ServiceBrokerListOptions, *resource.ServiceBroker](opts, func(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single service broker matching the options or an error if not exactly 1 match
func (c *ServiceBrokerClient) Single(ctx context.Context, opts *ServiceBrokerListOptions) (*resource.ServiceBroker, error) {
	return Single[*ServiceBrokerListOptions, *resource.ServiceBroker](opts, func(opts *ServiceBrokerListOptions) ([]*resource.ServiceBroker, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the service broker returning either a jobGUID or a service broker instance.
// Only metadata updates synchronously and return a service broker instance, all other updates return a jobGUID
func (c *ServiceBrokerClient) Update(ctx context.Context, guid string, r *resource.ServiceBrokerUpdate) (string, *resource.ServiceBroker, error) {
	var sb resource.ServiceBroker
	jobGUID, err := c.client.patch(ctx, path.Format("/v3/service_brokers/%s", guid), r, &sb)
	if err != nil {
		return "", nil, err
	}
	if jobGUID != "" {
		return jobGUID, nil, nil
	}
	return "", &sb, nil
}
