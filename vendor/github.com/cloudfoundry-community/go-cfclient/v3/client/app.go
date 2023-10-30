package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type AppClient commonClient

// AppListOptions list filters
type AppListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
	SpaceGUIDs        Filter `qs:"space_guids"`
	Stacks            Filter `qs:"stacks"`

	LifecycleType resource.LifecycleType  `qs:"lifecycle_type"`
	Include       resource.AppIncludeType `qs:"include"`
}

// NewAppListOptions creates new options to pass to list
func NewAppListOptions() *AppListOptions {
	return &AppListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o AppListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new app
func (c *AppClient) Create(ctx context.Context, r *resource.AppCreate) (*resource.App, error) {
	var app resource.App
	_, err := c.client.post(ctx, "/v3/apps", r, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Delete the specified app asynchronously and return a jobGUID.
func (c *AppClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/apps/%s", guid))
}

// First returns the first app matching the options or an error when less than 1 match
func (c *AppClient) First(ctx context.Context, opts *AppListOptions) (*resource.App, error) {
	return First[*AppListOptions, *resource.App](opts, func(opts *AppListOptions) ([]*resource.App, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified app
func (c *AppClient) Get(ctx context.Context, guid string) (*resource.App, error) {
	var app resource.App
	err := c.client.get(ctx, path.Format("/v3/apps/%s", guid), &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// GetIncludeSpace allows callers to fetch an app and include the parent space
func (c *AppClient) GetIncludeSpace(ctx context.Context, guid string) (*resource.App, *resource.Space, error) {
	var app resource.AppWithIncluded
	err := c.client.get(ctx, path.Format("/v3/apps/%s?include=%s", guid, resource.AppIncludeSpace), &app)
	if err != nil {
		return nil, nil, err
	}
	return &app.App, app.Included.Spaces[0], nil
}

// GetIncludeSpaceAndOrganization allows callers to fetch an app and include the parent space and organizations
func (c *AppClient) GetIncludeSpaceAndOrganization(ctx context.Context, guid string) (*resource.App, *resource.Space, *resource.Organization, error) {
	var app resource.AppWithIncluded
	err := c.client.get(ctx, path.Format("/v3/apps/%s?include=%s", guid, resource.AppIncludeSpaceOrganization), &app)
	if err != nil {
		return nil, nil, nil, err
	}
	return &app.App, app.Included.Spaces[0], app.Included.Organizations[0], nil
}

// GetEnvironment retrieves the environment variables that will be provided to an app at runtime.
// It will include environment variables for Environment Variable Groups and Service Bindings.
func (c *AppClient) GetEnvironment(ctx context.Context, guid string) (*resource.AppEnvironment, error) {
	var appEnv resource.AppEnvironment
	err := c.client.get(ctx, path.Format("/v3/apps/%s/env", guid), &appEnv)
	if err != nil {
		return nil, err
	}
	return &appEnv, nil
}

// GetEnvironmentVariables retrieves the environment variables that are associated with the given app
func (c *AppClient) GetEnvironmentVariables(ctx context.Context, guid string) (map[string]*string, error) {
	var appEnv resource.EnvVarResponse
	err := c.client.get(ctx, path.Format("/v3/apps/%s/environment_variables", guid), &appEnv)
	if err != nil {
		return nil, err
	}
	return appEnv.Var, nil
}

// List pages all the apps the user has access to
func (c *AppClient) List(ctx context.Context, opts *AppListOptions) ([]*resource.App, *Pager, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}
	opts.Include = resource.AppIncludeNone

	var res resource.AppList
	err := c.client.list(ctx, "/v3/apps", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all apps the user has access to
func (c *AppClient) ListAll(ctx context.Context, opts *AppListOptions) ([]*resource.App, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}
	return AutoPage[*AppListOptions, *resource.App](opts, func(opts *AppListOptions) ([]*resource.App, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListIncludeSpaces page all apps the user has access to and include the associated spaces
func (c *AppClient) ListIncludeSpaces(ctx context.Context, opts *AppListOptions) ([]*resource.App, []*resource.Space, *Pager, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}
	opts.Include = resource.AppIncludeSpace

	var res resource.AppList
	err := c.client.list(ctx, "/v3/apps", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, pager, nil
}

// ListIncludeSpacesAll retrieves all apps the user has access to and include the associated spaces
func (c *AppClient) ListIncludeSpacesAll(ctx context.Context, opts *AppListOptions) ([]*resource.App, []*resource.Space, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}

	var all []*resource.App
	var allSpaces []*resource.Space
	for {
		page, spaces, pager, err := c.ListIncludeSpaces(ctx, opts)
		if err != nil {
			return nil, nil, err
		}
		all = append(all, page...)
		allSpaces = append(allSpaces, spaces...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, allSpaces, nil
}

// ListIncludeSpacesAndOrganizations page all apps the user has access to and include the associated spaces and organizations
func (c *AppClient) ListIncludeSpacesAndOrganizations(ctx context.Context, opts *AppListOptions) ([]*resource.App, []*resource.Space, []*resource.Organization, *Pager, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}
	opts.Include = resource.AppIncludeSpaceOrganization

	var res resource.AppList
	err := c.client.list(ctx, "/v3/apps", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, res.Included.Spaces, res.Included.Organizations, pager, nil
}

// ListIncludeSpacesAndOrganizationsAll retrieves all apps the user has access to and include the associated spaces and organizations
func (c *AppClient) ListIncludeSpacesAndOrganizationsAll(ctx context.Context, opts *AppListOptions) ([]*resource.App, []*resource.Space, []*resource.Organization, error) {
	if opts == nil {
		opts = NewAppListOptions()
	}

	var all []*resource.App
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

// Permissions gets the current user’s permissions for the given app.
// If a user can see an app, then they can see its basic data.
// Only admin, read-only admins, and space developers can read sensitive data.
func (c *AppClient) Permissions(ctx context.Context, guid string) (*resource.AppPermissions, error) {
	var appPerms resource.AppPermissions
	err := c.client.get(ctx, path.Format("/v3/apps/%s/permissions", guid), &appPerms)
	if err != nil {
		return nil, err
	}
	return &appPerms, nil
}

// Restart will synchronously stop and start an application.
// Unlike the start and stop actions, this endpoint will error if the app is not successfully stopped in the runtime.
// For restarting applications without downtime, see the Deployments resource.
func (c *AppClient) Restart(ctx context.Context, guid string) (*resource.App, error) {
	var app resource.App
	_, err := c.client.post(ctx, path.Format("/v3/apps/%s/actions/restart", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// SetEnvironmentVariables updates the environment variables associated with the given app.
// The variables given in the request will be merged with the existing app environment variables.
// Any requested variables with a value of null will be removed from the app.
//
// Environment variable names may not start with VCAP_
// PORT is not a valid environment variable.
func (c *AppClient) SetEnvironmentVariables(ctx context.Context, guid string, envRequest map[string]*string) (map[string]*string, error) {
	req := &resource.EnvVar{
		Var: envRequest,
	}
	var res resource.EnvVarResponse
	_, err := c.client.patch(ctx, path.Format("/v3/apps/%s/environment_variables", guid), req, &res)
	if err != nil {
		return nil, err
	}
	return res.Var, nil
}

// Single returns a single app matching the options or an error if not exactly 1 match
func (c *AppClient) Single(ctx context.Context, opts *AppListOptions) (*resource.App, error) {
	return Single[*AppListOptions, *resource.App](opts, func(opts *AppListOptions) ([]*resource.App, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Start the app if not already started
func (c *AppClient) Start(ctx context.Context, guid string) (*resource.App, error) {
	var app resource.App
	_, err := c.client.post(ctx, path.Format("/v3/apps/%s/actions/start", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Stop the app if not already stopped
func (c *AppClient) Stop(ctx context.Context, guid string) (*resource.App, error) {
	var app resource.App
	_, err := c.client.post(ctx, path.Format("/v3/apps/%s/actions/stop", guid), nil, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Update the specified attributes of the app
func (c *AppClient) Update(ctx context.Context, guid string, r *resource.AppUpdate) (*resource.App, error) {
	var app resource.App
	_, err := c.client.patch(ctx, path.Format("/v3/apps/%s", guid), r, &app)
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// SSHEnabled returns if an application’s runtime environment will accept ssh connections.
// If ssh is disabled, the reason field will describe whether it is disabled globally,
// at the space level, or at the app level.
func (c *AppClient) SSHEnabled(ctx context.Context, guid string) (*resource.AppSSHEnabled, error) {
	var appSSH resource.AppSSHEnabled
	err := c.client.get(ctx, path.Format("/v3/apps/%s/ssh_enabled", guid), &appSSH)
	if err != nil {
		return nil, err
	}
	return &appSSH, nil
}
