package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceCredentialBindingClient commonClient

// ServiceCredentialBindingListOptions list filters
type ServiceCredentialBindingListOptions struct {
	*ListOptions

	Names                Filter `qs:"names"`                  // list of service credential binding names to filter by
	ServiceInstanceGUIDs Filter `qs:"service_instance_guids"` // list of SI guids to filter by
	ServiceInstanceNames Filter `qs:"service_instance_names"` // list of SI names to filter by
	AppGUIDs             Filter `qs:"app_guids"`              // list of app guids to filter by
	AppNames             Filter `qs:"app_names"`              // list of app names to filter by
	ServicePlanGUIDs     Filter `qs:"service_plan_guids"`     // list of service plan guids to filter by
	ServicePlanNames     Filter `qs:"service_plan_names"`     // list of service plan names to filter by
	ServiceOfferingGUIDs Filter `qs:"service_offering_guids"` // list of service offering guids to filter by
	ServiceOfferingNames Filter `qs:"service_offering_names"` // list of service offering names to filter by
	Type                 Filter `qs:"type"`                   // list of service credential binding types to filter by, app or key
	GUIDs                Filter `qs:"guids"`                  // list of service route binding guids to filter by

	Include resource.ServiceCredentialBindingIncludeType `qs:"include"`
}

// NewServiceCredentialBindingListOptions creates new options to pass to list
func NewServiceCredentialBindingListOptions() *ServiceCredentialBindingListOptions {
	return &ServiceCredentialBindingListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceCredentialBindingListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new service credential binding
func (c *ServiceCredentialBindingClient) Create(ctx context.Context, r *resource.ServiceCredentialBindingCreate) (string, *resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	jobGUID, err := c.client.post(ctx, "/v3/service_credential_bindings", r, &d)
	if err != nil {
		return "", nil, err
	}
	if jobGUID != "" {
		return jobGUID, nil, nil
	}
	return "", &d, nil
}

// Delete the specified service credential binding
func (c *ServiceCredentialBindingClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/service_credential_bindings/%s", guid))
	return err
}

// First returns the first service credential binding matching the options or an error when less than 1 match
func (c *ServiceCredentialBindingClient) First(ctx context.Context, opts *ServiceCredentialBindingListOptions) (*resource.ServiceCredentialBinding, error) {
	return First[*ServiceCredentialBindingListOptions, *resource.ServiceCredentialBinding](opts, func(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service credential binding
func (c *ServiceCredentialBindingClient) Get(ctx context.Context, guid string) (*resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	err := c.client.get(ctx, path.Format("/v3/service_credential_bindings/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetDetails the specified service credential binding details
func (c *ServiceCredentialBindingClient) GetDetails(ctx context.Context, guid string) (*resource.ServiceCredentialBindingDetails, error) {
	var d resource.ServiceCredentialBindingDetails
	err := c.client.get(ctx, path.Format("/v3/service_credential_bindings/%s/details", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// GetParameters the specified service credential binding details
func (c *ServiceCredentialBindingClient) GetParameters(ctx context.Context, guid string) (map[string]string, error) {
	var p map[string]string
	err := c.client.get(ctx, path.Format("/v3/service_credential_bindings/%s/parameters", guid), &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// GetIncludeApp allows callers to fetch a service credential binding and include the associated app
func (c *ServiceCredentialBindingClient) GetIncludeApp(ctx context.Context, guid string) (*resource.ServiceCredentialBinding, *resource.App, error) {
	var r resource.ServiceCredentialBindingWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_credential_bindings/%s?include=%s", guid, resource.ServiceCredentialBindingIncludeApp), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.ServiceCredentialBinding, r.Included.Apps[0], nil
}

// GetIncludeServiceInstance allows callers to fetch a service credential binding and include the associated service instance
func (c *ServiceCredentialBindingClient) GetIncludeServiceInstance(ctx context.Context, guid string) (*resource.ServiceCredentialBinding, *resource.ServiceInstance, error) {
	var r resource.ServiceCredentialBindingWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_credential_bindings/%s?include=%s", guid, resource.ServiceCredentialBindingIncludeServiceInstance), &r)
	if err != nil {
		return nil, nil, err
	}
	return &r.ServiceCredentialBinding, r.Included.ServiceInstances[0], nil
}

// List pages ServiceCredentialBindings the user has access to
func (c *ServiceCredentialBindingClient) List(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
	var res resource.ServiceCredentialBindingList
	err := c.client.list(ctx, "/v3/service_credential_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all ServiceCredentialBindings the user has access to
func (c *ServiceCredentialBindingClient) ListAll(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}
	return AutoPage[*ServiceCredentialBindingListOptions, *resource.ServiceCredentialBinding](opts, func(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeApps pages all service credential bindings the user has access to and include the associated apps
func (c *ServiceCredentialBindingClient) ListIncludeApps(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, []*resource.App, *Pager, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}
	opts.Include = resource.ServiceCredentialBindingIncludeApp

	var res resource.ServiceCredentialBindingList
	err := c.client.list(ctx, "/v3/service_credential_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Apps, pager, nil
}

// ListIncludeAppsAll retrieves all service credential bindings the user has access to and include the associated apps
func (c *ServiceCredentialBindingClient) ListIncludeAppsAll(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, []*resource.App, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}

	var all []*resource.ServiceCredentialBinding
	var allApps []*resource.App
	for {
		page, apps, pager, err := c.ListIncludeApps(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allApps = append(allApps, apps...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allApps, nil
}

// ListIncludeServiceInstances pages all service credential bindings the user has access to and include the associated SIs
func (c *ServiceCredentialBindingClient) ListIncludeServiceInstances(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, []*resource.ServiceInstance, *Pager, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}
	opts.Include = resource.ServiceCredentialBindingIncludeServiceInstance

	var res resource.ServiceCredentialBindingList
	err := c.client.list(ctx, "/v3/service_credential_bindings", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.ServiceInstances, pager, nil
}

// ListIncludeServiceInstancesAll retrieves all service credential bindings the user has access to and include the associated SIs
func (c *ServiceCredentialBindingClient) ListIncludeServiceInstancesAll(ctx context.Context, opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, []*resource.ServiceInstance, error) {
	if opts == nil {
		opts = NewServiceCredentialBindingListOptions()
	}

	var all []*resource.ServiceCredentialBinding
	var allServiceInstances []*resource.ServiceInstance
	for {
		page, serviceInstances, pager, err := c.ListIncludeServiceInstances(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allServiceInstances = append(allServiceInstances, serviceInstances...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allServiceInstances, nil
}

// Single returns a single service credential binding matching the options or an error if not exactly 1 match
func (c *ServiceCredentialBindingClient) Single(ctx context.Context, opts *ServiceCredentialBindingListOptions) (*resource.ServiceCredentialBinding, error) {
	return Single[*ServiceCredentialBindingListOptions, *resource.ServiceCredentialBinding](opts, func(opts *ServiceCredentialBindingListOptions) ([]*resource.ServiceCredentialBinding, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the app
func (c *ServiceCredentialBindingClient) Update(ctx context.Context, guid string, r *resource.ServiceCredentialBindingUpdate) (*resource.ServiceCredentialBinding, error) {
	var d resource.ServiceCredentialBinding
	_, err := c.client.patch(ctx, path.Format("/v3/service_credential_bindings/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
