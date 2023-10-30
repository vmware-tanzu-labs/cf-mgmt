package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceOfferingClient commonClient

// ServiceOfferingListOptions list filters
type ServiceOfferingListOptions struct {
	*ListOptions

	Names              Filter `qs:"names"`
	ServiceBrokerGUIDs Filter `qs:"service_broker_guids"`
	ServiceBrokerNames Filter `qs:"service_broker_names"`
	SpaceGUIDs         Filter `qs:"space_guids"`
	OrganizationGUIDs  Filter `qs:"organization_guids"`
	Available          *bool  `qs:"available"`
}

// NewServiceOfferingListOptions creates new options to pass to list
func NewServiceOfferingListOptions() *ServiceOfferingListOptions {
	return &ServiceOfferingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceOfferingListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Delete the specified service offering
func (c *ServiceOfferingClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/service_offerings/%s", guid))
	return err
}

// First returns the first service offering matching the options or an error when less than 1 match
func (c *ServiceOfferingClient) First(ctx context.Context, opts *ServiceOfferingListOptions) (*resource.ServiceOffering, error) {
	return First[*ServiceOfferingListOptions, *resource.ServiceOffering](opts, func(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service offering
func (c *ServiceOfferingClient) Get(ctx context.Context, guid string) (*resource.ServiceOffering, error) {
	var ServiceOffering resource.ServiceOffering
	err := c.client.get(ctx, path.Format("/v3/service_offerings/%s", guid), &ServiceOffering)
	if err != nil {
		return nil, err
	}
	return &ServiceOffering, nil
}

// List pages service offerings the user has access to
func (c *ServiceOfferingClient) List(ctx context.Context, opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
	if opts == nil {
		opts = NewServiceOfferingListOptions()
	}

	var res resource.ServiceOfferingList
	err := c.client.list(ctx, "/v3/service_offerings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service offerings the user has access to
func (c *ServiceOfferingClient) ListAll(ctx context.Context, opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, error) {
	if opts == nil {
		opts = NewServiceOfferingListOptions()
	}
	return AutoPage[*ServiceOfferingListOptions, *resource.ServiceOffering](opts, func(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single service offering matching the options or an error if not exactly 1 match
func (c *ServiceOfferingClient) Single(ctx context.Context, opts *ServiceOfferingListOptions) (*resource.ServiceOffering, error) {
	return Single[*ServiceOfferingListOptions, *resource.ServiceOffering](opts, func(opts *ServiceOfferingListOptions) ([]*resource.ServiceOffering, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the service offering
func (c *ServiceOfferingClient) Update(ctx context.Context, guid string, r *resource.ServiceOfferingUpdate) (*resource.ServiceOffering, error) {
	var res resource.ServiceOffering
	_, err := c.client.patch(ctx, path.Format("/v3/service_offerings/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
