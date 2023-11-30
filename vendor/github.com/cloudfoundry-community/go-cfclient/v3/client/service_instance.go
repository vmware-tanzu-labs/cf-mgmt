package client

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ServiceInstanceClient commonClient

// ServiceInstanceListOptions list filters
type ServiceInstanceListOptions struct {
	*ListOptions

	Names             Filter `qs:"names"` // list of service instance names to filter by
	GUIDs             Filter `qs:"guids"` // list of service instance guids to filter by
	Type              string `qs:"type"`  // Filter by type; valid values are managed and user-provided
	SpaceGUIDs        Filter `qs:"space_guids"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
	ServicePlanGUIDs  Filter `qs:"service_plan_guids"`
	ServicePlanNames  Filter `qs:"service_plan_names"`
}

// NewServiceInstanceListOptions creates new options to pass to list
func NewServiceInstanceListOptions() *ServiceInstanceListOptions {
	return &ServiceInstanceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ServiceInstanceListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// CreateManaged requests a new service instance asynchronously from a broker. The result
// of this call is an error or the jobGUID.
func (c *ServiceInstanceClient) CreateManaged(ctx context.Context, r *resource.ServiceInstanceCreate) (string, error) {
	var si resource.ServiceInstance
	jobGUID, err := c.client.post(ctx, "/v3/service_instances", r, &si)
	if err != nil {
		return "", err
	}
	return jobGUID, nil
}

// CreateUserProvided creates a new user provided service instance. User provided service instances
// do not require interactions with service brokers.
func (c *ServiceInstanceClient) CreateUserProvided(ctx context.Context, r *resource.ServiceInstanceCreate) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	_, err := c.client.post(ctx, "/v3/service_instances", r, &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

// Delete the specified service instance returning the async deletion jobGUID
func (c *ServiceInstanceClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/service_instances/%s", guid))
}

// First returns the first service instance matching the options or an error when less than 1 match
func (c *ServiceInstanceClient) First(ctx context.Context, opts *ServiceInstanceListOptions) (*resource.ServiceInstance, error) {
	return First[*ServiceInstanceListOptions, *resource.ServiceInstance](opts, func(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified service instance
func (c *ServiceInstanceClient) Get(ctx context.Context, guid string) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s", guid), &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}

// GetUserPermissions retrieves the current user’s permissions for the given service instance
//
// If a user can get a service instance then they can ‘read’ it. Users who can update a service instance can ‘manage’ it.
//
// This endpoint’s primary purpose is to enable third-party service dashboards to determine the permissions of a
// given Cloud Foundry user that has authenticated with the dashboard via single sign-on (SSO). For more information,
// see the Cloud Foundry documentation on Dashboard Single Sign-On.
func (c *ServiceInstanceClient) GetUserPermissions(ctx context.Context, guid string) (*resource.ServiceInstanceUserPermissions, error) {
	var permissions resource.ServiceInstanceUserPermissions
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s/permissions", guid), &permissions)
	if err != nil {
		return nil, err
	}
	return &permissions, nil
}

// GetManagedParameters queries the service broker for the parameters associated with this managed service instance
//
// The broker catalog must have enabled the instances_retrievable feature for the Service Offering.
// Check the Service Offering object for the value of this feature flag.
func (c *ServiceInstanceClient) GetManagedParameters(ctx context.Context, guid string) (*json.RawMessage, error) {
	var parameters json.RawMessage
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s/parameters", guid), &parameters)
	if err != nil {
		return nil, err
	}
	return &parameters, nil
}

// GetUserProvidedCredentials the specified user provided service instance credentials
func (c *ServiceInstanceClient) GetUserProvidedCredentials(ctx context.Context, guid string) (*json.RawMessage, error) {
	var credentials json.RawMessage
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s/credentials", guid), &credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

// GetSharedSpaceRelationships lists the spaces that the service instance has been shared to
func (c *ServiceInstanceClient) GetSharedSpaceRelationships(ctx context.Context, guid string) (*resource.ServiceInstanceSharedSpaceRelationships, error) {
	var relations resource.ServiceInstanceSharedSpaceRelationships
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s/relationships/shared_spaces", guid), &relations)
	if err != nil {
		return nil, err
	}
	return &relations, nil
}

// GetSharedSpaceUsageSummary retrieves the number of bound apps in spaces where the service instance has been shared to
func (c *ServiceInstanceClient) GetSharedSpaceUsageSummary(ctx context.Context, guid string) (*resource.ServiceInstanceUsageSummary, error) {
	var usage resource.ServiceInstanceUsageSummary
	err := c.client.get(ctx, path.Format("/v3/service_instances/%s/relationships/shared_spaces/usage_summary", guid), &usage)
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

// List pages all service instances the user has access to
func (c *ServiceInstanceClient) List(ctx context.Context, opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
	if opts == nil {
		opts = NewServiceInstanceListOptions()
	}
	var res resource.ServiceInstanceList
	err := c.client.list(ctx, "/v3/service_instances", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all service instances the user has access to
func (c *ServiceInstanceClient) ListAll(ctx context.Context, opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, error) {
	if opts == nil {
		opts = NewServiceInstanceListOptions()
	}
	return AutoPage[*ServiceInstanceListOptions, *resource.ServiceInstance](opts, func(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ShareWithSpace shares the service instance with the specified space
//
// In order to share into a space the requesting user must be a space developer in the target space
func (c *ServiceInstanceClient) ShareWithSpace(ctx context.Context, guid string, spaceGUID string) (*resource.ServiceInstanceSharedSpaceRelationships, error) {
	return c.ShareWithSpaces(ctx, guid, []string{spaceGUID})
}

// ShareWithSpaces shares the service instance with the specified spaces
//
// In order to share into a space the requesting user must be a space developer in the target space
func (c *ServiceInstanceClient) ShareWithSpaces(ctx context.Context, guid string, spaceGUIDs []string) (*resource.ServiceInstanceSharedSpaceRelationships, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relationships resource.ServiceInstanceSharedSpaceRelationships
	_, err := c.client.post(ctx, path.Format("/v3/service_instances/%s/relationships/shared_spaces", guid), req, &relationships)
	if err != nil {
		return nil, err
	}
	return &relationships, nil
}

// Single returns a single service instance matching the options or an error if not exactly 1 match
func (c *ServiceInstanceClient) Single(ctx context.Context, opts *ServiceInstanceListOptions) (*resource.ServiceInstance, error) {
	return Single[*ServiceInstanceListOptions, *resource.ServiceInstance](opts, func(opts *ServiceInstanceListOptions) ([]*resource.ServiceInstance, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// UnShareWithSpace un=shares the service instance with the specified space
//
// This will automatically unbind any applications bound to this service instance in the specified space
// Un-sharing a service instance from a space will not delete any service keys
func (c *ServiceInstanceClient) UnShareWithSpace(ctx context.Context, guid string, spaceGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/service_instances/%s/relationships/shared_spaces/%s", guid, spaceGUID))
	return err
}

// UnShareWithSpaces un-shares the service instance with the specified spaces
//
// This will automatically unbind any applications bound to this service instance in the specified space
// Un-sharing a service instance from a space will not delete any service keys
func (c *ServiceInstanceClient) UnShareWithSpaces(ctx context.Context, guid string, spaceGUIDs []string) error {
	for _, s := range spaceGUIDs {
		err := c.UnShareWithSpace(ctx, guid, s)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateManaged updates the specified attributes of the managed service instance returning either a jobGUID or a
// service instance object
//
// Only metadata, tags, and name (when allow_context_updates feature disabled) updates synchronously and return a service
// instance object, all other updates return a jobGUID
func (c *ServiceInstanceClient) UpdateManaged(ctx context.Context, guid string, r *resource.ServiceInstanceManagedUpdate) (string, *resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	jobGUID, err := c.client.patch(ctx, path.Format("/v3/service_instances/%s", guid), r, &si)
	if err != nil {
		return "", nil, err
	}
	if jobGUID != "" {
		return jobGUID, nil, nil
	}
	return "", &si, nil
}

// UpdateUserProvided updates the specified attributes of the user-provided service instance returning a
// service instance object
func (c *ServiceInstanceClient) UpdateUserProvided(ctx context.Context, guid string, r *resource.ServiceInstanceUserProvidedUpdate) (*resource.ServiceInstance, error) {
	var si resource.ServiceInstance
	_, err := c.client.patch(ctx, path.Format("/v3/service_instances/%s", guid), r, &si)
	if err != nil {
		return nil, err
	}
	return &si, nil
}
