package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceRouteBindingClient commonClient

// ServiceRouteBindingListOptions list filters
type ServiceRouteBindingListOptions struct {
	*ListOptions

	GUIDs                Filter `qs:"guids"`
	RouteGUIDs           Filter `qs:"route_guids"`
	ServiceInstanceGUIDs Filter `qs:"service_instance_guids"`
	ServiceInstanceNames Filter `qs:"service_instance_names"`

	Include resource.ServiceRouteBindingIncludeType `qs:"include"`
}

// NewServiceRouteBindingListOptions creates new options to pass to list
func NewServiceRouteBindingListOptions() *ServiceRouteBindingListOptions {
	return &ServiceRouteBindingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceRouteBindingListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service route binding returning the jobGUID for managed service instances or the
// service route binding object for user provided service instances
func (c *ServiceRouteBindingClient) Create(ctx context.Context, r *resource.ServiceRouteBindingCreate) (string, *resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	jobGUID, err := c.client.post(ctx, "/v3/service_route_bindings", r, &srb)
	if err != nil {
		return "", nil, err
	}
	if jobGUID != "" {
		return jobGUID, nil, nil
	}
	return "", &srb, nil
}

// Delete the specified service route binding returning the jobGUID for managed service instances or empty string
// for user provided service instances
func (c *ServiceRouteBindingClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/service_route_bindings/%s", guid))
}

// First returns the first service route binding matching the options or an error when less than 1 match
func (c *ServiceRouteBindingClient) First(ctx context.Context, opts *ServiceRouteBindingListOptions) (*resource.ServiceRouteBinding, error) {
	return First[*ServiceRouteBindingListOptions, *resource.ServiceRouteBinding](opts, func(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service route binding
func (c *ServiceRouteBindingClient) Get(ctx context.Context, guid string) (*resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	err := c.client.get(ctx, path.Format("/v3/service_route_bindings/%s", guid), &srb)
	if err != nil {
		return nil, err
	}
	return &srb, nil
}

// GetIncludeRoute allows callers to fetch a service route binding and include the associated route
func (c *ServiceRouteBindingClient) GetIncludeRoute(ctx context.Context, guid string) (*resource.ServiceRouteBinding, *resource.Route, error) {
	var srb resource.ServiceRouteBindingWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_route_bindings/%s?include=%s", guid, resource.ServiceRouteBindingIncludeRoute), &srb)
	if err != nil {
		return nil, nil, err
	}
	return &srb.ServiceRouteBinding, srb.Included.Routes[0], nil
}

// GetIncludeServiceInstance allows callers to fetch a service route binding and include the associated service instance
func (c *ServiceRouteBindingClient) GetIncludeServiceInstance(ctx context.Context, guid string) (*resource.ServiceRouteBinding, *resource.ServiceInstance, error) {
	var srb resource.ServiceRouteBindingWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_route_bindings/%s?include=%s", guid, resource.ServiceRouteBindingIncludeServiceInstance), &srb)
	if err != nil {
		return nil, nil, err
	}
	return &srb.ServiceRouteBinding, srb.Included.ServiceInstances[0], nil
}

// GetParameters queries the Service Broker for the parameters associated with this service route binding
func (c *ServiceRouteBindingClient) GetParameters(ctx context.Context, guid string) (map[string]string, error) {
	var srbEnv map[string]string
	err := c.client.get(ctx, path.Format("/v3/service_route_bindings/%s/parameters", guid), &srbEnv)
	if err != nil {
		return nil, err
	}
	return srbEnv, nil
}

// List pages all the service route bindings the user has access to
func (c *ServiceRouteBindingClient) List(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.list(ctx, "/v3/service_route_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service route bindings the user has access to
func (c *ServiceRouteBindingClient) ListAll(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	return AutoPage[*ServiceRouteBindingListOptions, *resource.ServiceRouteBinding](opts, func(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeRoutes page all service route bindings the user has access to and include the associated routes
func (c *ServiceRouteBindingClient) ListIncludeRoutes(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.Route, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.list(ctx, "/v3/service_route_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Routes, pager, nil
}

// ListIncludeRoutesAll retrieves all service route bindings the user has access to and include the associated routes
func (c *ServiceRouteBindingClient) ListIncludeRoutesAll(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.Route, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}

	var all []*resource.ServiceRouteBinding
	var allRoutes []*resource.Route
	for {
		page, routes, pager, err := c.ListIncludeRoutes(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allRoutes = append(allRoutes, routes...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allRoutes, nil
}

// ListIncludeServiceInstances page all service route bindings the user has access to and include the
// associated service instances
func (c *ServiceRouteBindingClient) ListIncludeServiceInstances(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.ServiceInstance, *Pager, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}
	opts.Include = resource.ServiceRouteBindingIncludeNone

	var res resource.ServiceRouteBindingList
	err := c.client.list(ctx, "/v3/service_route_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.ServiceInstances, pager, nil
}

// ListIncludeServiceInstancesAll retrieves all service route bindings the user has access to and include the
// associated service instances
func (c *ServiceRouteBindingClient) ListIncludeServiceInstancesAll(ctx context.Context, opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, []*resource.ServiceInstance, error) {
	if opts == nil {
		opts = NewServiceRouteBindingListOptions()
	}

	var all []*resource.ServiceRouteBinding
	var allSIs []*resource.ServiceInstance
	for {
		page, sis, pager, err := c.ListIncludeServiceInstances(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSIs = append(allSIs, sis...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSIs, nil
}

// Single returns a single service route binding matching the options or an error if not exactly 1 match
func (c *ServiceRouteBindingClient) Single(ctx context.Context, opts *ServiceRouteBindingListOptions) (*resource.ServiceRouteBinding, error) {
	return Single[*ServiceRouteBindingListOptions, *resource.ServiceRouteBinding](opts, func(opts *ServiceRouteBindingListOptions) ([]*resource.ServiceRouteBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the service route binding
func (c *ServiceRouteBindingClient) Update(ctx context.Context, guid string, r *resource.ServiceRouteBindingUpdate) (*resource.ServiceRouteBinding, error) {
	var srb resource.ServiceRouteBinding
	_, err := c.client.patch(ctx, path.Format("/v3/service_route_bindings/%s", guid), r, &srb)
	if err != nil {
		return nil, err
	}
	return &srb, nil
}
