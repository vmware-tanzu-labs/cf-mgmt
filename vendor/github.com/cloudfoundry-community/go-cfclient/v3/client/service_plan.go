package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServicePlanClient commonClient

// ServicePlanListOptions list filters
type ServicePlanListOptions struct {
	*ListOptions

	Names                Filter `qs:"names"`
	BrokerCatalogIDs     Filter `qs:"broker_catalog_ids"`
	SpaceGUIDs           Filter `qs:"space_guids"`
	OrganizationGUIDs    Filter `qs:"organization_guids"`
	ServiceBrokerGUIDs   Filter `qs:"service_broker_guids"`
	ServiceBrokerNames   Filter `qs:"service_broker_names"`
	ServiceOfferingGUIDs Filter `qs:"service_offering_guids"`
	ServiceOfferingNames Filter `qs:"service_offering_names"`
	ServiceInstanceGUIDs Filter `qs:"service_instance_guids"`
	Available            *bool  `qs:"available"`

	Include resource.ServicePlanIncludeType `qs:"include"`
}

// NewServicePlanListOptions creates new options to pass to list
func NewServicePlanListOptions() *ServicePlanListOptions {
	return &ServicePlanListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServicePlanListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Delete the specified service plan
func (c *ServicePlanClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/service_plans/%s", guid))
	return err
}

// First returns the first service plan matching the options or an error when less than 1 match
func (c *ServicePlanClient) First(ctx context.Context, opts *ServicePlanListOptions) (*resource.ServicePlan, error) {
	return First[*ServicePlanListOptions, *resource.ServicePlan](opts, func(opts *ServicePlanListOptions) ([]*resource.ServicePlan, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service plan
func (c *ServicePlanClient) Get(ctx context.Context, guid string) (*resource.ServicePlan, error) {
	var ServicePlan resource.ServicePlan
	err := c.client.get(ctx, path.Format("/v3/service_plans/%s", guid), &ServicePlan)
	if err != nil {
		return nil, err
	}
	return &ServicePlan, nil
}

// GetIncludeServicePlan allows callers to fetch a service plan and include the associated service offering
func (c *ServicePlanClient) GetIncludeServicePlan(ctx context.Context, guid string) (*resource.ServicePlan, *resource.ServiceOffering, error) {
	var servicePlan resource.ServicePlanWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_plans/%s?include=%s", guid, resource.ServicePlanIncludeServiceOffering), &servicePlan)
	if err != nil {
		return nil, nil, err
	}
	return &servicePlan.ServicePlan, servicePlan.Included.ServiceOfferings[0], nil
}

// GetIncludeSpaceAndOrganization allows callers to fetch a service plan and include the parent space and organization
func (c *ServicePlanClient) GetIncludeSpaceAndOrganization(ctx context.Context, guid string) (*resource.ServicePlan, *resource.Space, *resource.Organization, error) {
	var servicePlan resource.ServicePlanWithIncluded
	err := c.client.get(ctx, path.Format("/v3/service_plans/%s?include=%s", guid, resource.ServicePlanIncludeSpaceOrganization), &servicePlan)
	if err != nil {
		return nil, nil, nil, err
	}
	return &servicePlan.ServicePlan, servicePlan.Included.Spaces[0], servicePlan.Included.Organizations[0], nil
}

// List pages service plans the user has access to
func (c *ServicePlanClient) List(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, *Pager, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}

	var res resource.ServicePlanList
	err := c.client.list(ctx, "/v3/service_plans", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service plans the user has access to
func (c *ServicePlanClient) ListAll(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}
	return AutoPage[*ServicePlanListOptions, *resource.ServicePlan](opts, func(opts *ServicePlanListOptions) ([]*resource.ServicePlan, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeServiceOffering page all service plans the user has access to and include the associated service offerings
func (c *ServicePlanClient) ListIncludeServiceOffering(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, []*resource.ServiceOffering, *Pager, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}
	opts.Include = resource.ServicePlanIncludeServiceOffering

	var res resource.ServicePlanList
	err := c.client.list(ctx, "/v3/service_plans", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.ServiceOfferings, pager, nil
}

// ListIncludeServiceOfferingAll retrieves all service plans the user has access to and include the associated service offerings
func (c *ServicePlanClient) ListIncludeServiceOfferingAll(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, []*resource.ServiceOffering, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}

	var all []*resource.ServicePlan
	var allServiceOfferings []*resource.ServiceOffering
	for {
		page, serviceOfferings, pager, err := c.ListIncludeServiceOffering(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allServiceOfferings = append(allServiceOfferings, serviceOfferings...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allServiceOfferings, nil
}

// ListIncludeSpacesAndOrganizations page all service plans the user has access to and include the associated spaces and organizations
func (c *ServicePlanClient) ListIncludeSpacesAndOrganizations(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, []*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}
	opts.Include = resource.ServicePlanIncludeSpaceOrganization

	var res resource.ServicePlanList
	err := c.client.list(ctx, "/v3/service_plans", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, res.Included.Organizations, pager, nil
}

// ListIncludeSpacesAndOrganizationsAll retrieves all service plans the user has access to and include the associated spaces and organizations
func (c *ServicePlanClient) ListIncludeSpacesAndOrganizationsAll(ctx context.Context, opts *ServicePlanListOptions) ([]*resource.ServicePlan, []*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewServicePlanListOptions()
	}

	var all []*resource.ServicePlan
	var allSpaces []*resource.Space
	var allOrgs []*resource.Organization
	for {
		page, spaces, orgs, pager, err := c.ListIncludeSpacesAndOrganizations(ctx, opts)
		if err != nil {
			return nil, nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		allOrgs = append(allOrgs, orgs...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, allOrgs, nil
}

// Single returns a single service plan matching the options or an error if not exactly 1 match
func (c *ServicePlanClient) Single(ctx context.Context, opts *ServicePlanListOptions) (*resource.ServicePlan, error) {
	return Single[*ServicePlanListOptions, *resource.ServicePlan](opts, func(opts *ServicePlanListOptions) ([]*resource.ServicePlan, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the service plan
func (c *ServicePlanClient) Update(ctx context.Context, guid string, r *resource.ServicePlanUpdate) (*resource.ServicePlan, error) {
	var res resource.ServicePlan
	_, err := c.client.patch(ctx, path.Format("/v3/service_plans/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
