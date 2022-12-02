package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SpaceFeatureClient commonClient

// EnableSSH toggles the SSH feature for a space
func (c *SpaceFeatureClient) EnableSSH(ctx context.Context, spaceGUID string, enable bool) error {
	r := resource.SpaceFeatureUpdate{
		Enabled: enable,
	}
	_, err := c.client.patch(ctx, path.Format("/v3/spaces/%s/features/ssh", spaceGUID), r, nil)
	return err
}

// IsSSHEnabled returns true if SSH is enabled for the specified space
func (c *SpaceFeatureClient) IsSSHEnabled(ctx context.Context, spaceGUID string) (bool, error) {
	var sf resource.SpaceFeature
	err := c.client.get(ctx, path.Format("/v3/spaces/%s/features/ssh", spaceGUID), &sf)
	if err != nil {
		return false, err
	}
	return sf.Enabled, nil
}
