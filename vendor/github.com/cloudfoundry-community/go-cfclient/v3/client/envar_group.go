package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type EnvVarGroupClient commonClient

// Get retrieves the specified envvar group
func (c *EnvVarGroupClient) Get(ctx context.Context, name string) (*resource.EnvVarGroup, error) {
	var e resource.EnvVarGroup
	err := c.client.get(ctx, path.Format("/v3/environment_variable_groups/%s", name), &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// GetRunning retrieves the running envvar group
func (c *EnvVarGroupClient) GetRunning(ctx context.Context) (*resource.EnvVarGroup, error) {
	return c.Get(ctx, "running")
}

// GetStaging retrieves the running envvar group
func (c *EnvVarGroupClient) GetStaging(ctx context.Context) (*resource.EnvVarGroup, error) {
	return c.Get(ctx, "staging")
}

// Update the specified attributes of the envar group
func (c *EnvVarGroupClient) Update(ctx context.Context, name string, r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	var e resource.EnvVarGroup
	_, err := c.client.patch(ctx, path.Format("/v3/environment_variable_groups/%s", name), r, &e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

// UpdateRunning updates the specified attributes of the running envar group
func (c *EnvVarGroupClient) UpdateRunning(ctx context.Context, r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	return c.Update(ctx, "running", r)
}

// UpdateStaging updates the specified attributes of the staging envar group
func (c *EnvVarGroupClient) UpdateStaging(ctx context.Context, r *resource.EnvVarGroupUpdate) (*resource.EnvVarGroup, error) {
	return c.Update(ctx, "staging", r)
}
