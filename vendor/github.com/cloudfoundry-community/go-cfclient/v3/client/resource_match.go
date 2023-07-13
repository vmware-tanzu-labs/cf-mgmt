package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ResourceMatchClient commonClient

// Create a list of cached resources from the input list
func (c *ResourceMatchClient) Create(ctx context.Context, toMatch *resource.ResourceMatches) (*resource.ResourceMatches, error) {
	var matched resource.ResourceMatches
	_, err := c.client.post(ctx, "/v3/resource_matches", toMatch, &matched)
	if err != nil {
		return nil, err
	}
	return &matched, nil
}
