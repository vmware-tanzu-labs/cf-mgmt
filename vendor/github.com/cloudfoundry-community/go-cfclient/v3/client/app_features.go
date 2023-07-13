package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type AppFeatureClient commonClient

// Get retrieves the named app feature
func (c *AppFeatureClient) Get(ctx context.Context, appGUID, featureName string) (*resource.AppFeature, error) {
	var a resource.AppFeature
	err := c.client.get(ctx, path.Format("/v3/apps/%s/features/%s", appGUID, featureName), &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetSSH retrieves the SSH app feature
func (c *AppFeatureClient) GetSSH(ctx context.Context, appGUID string) (*resource.AppFeature, error) {
	return c.Get(ctx, appGUID, "ssh")
}

// GetRevisions retrieves the revisions app feature
func (c *AppFeatureClient) GetRevisions(ctx context.Context, appGUID string) (*resource.AppFeature, error) {
	return c.Get(ctx, appGUID, "revisions")
}

// List pages all app features
func (c *AppFeatureClient) List(ctx context.Context, appGUID string) ([]*resource.AppFeature, *Pager, error) {
	var res resource.AppFeatureList
	err := c.client.get(ctx, path.Format("/v3/apps/%s/features", appGUID), &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// Update the enabled attribute of the named app feature
func (c *AppFeatureClient) Update(ctx context.Context, appGUID, featureName string, enabled bool) (*resource.AppFeature, error) {
	r := &resource.AppFeatureUpdate{
		Enabled: enabled,
	}
	var a resource.AppFeature
	_, err := c.client.patch(ctx, path.Format("/v3/apps/%s/features/%s", appGUID, featureName), r, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// UpdateSSH updated the enabled attribute of the SSH app feature
func (c *AppFeatureClient) UpdateSSH(ctx context.Context, appGUID string, enabled bool) (*resource.AppFeature, error) {
	return c.Update(ctx, appGUID, "ssh", enabled)
}

// UpdateRevisions updated the enabled attribute of the revisions app feature
func (c *AppFeatureClient) UpdateRevisions(ctx context.Context, appGUID string, enabled bool) (*resource.AppFeature, error) {
	return c.Update(ctx, appGUID, "revisions", enabled)
}
