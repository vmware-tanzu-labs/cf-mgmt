package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type FeatureFlagClient commonClient

// FeatureFlagListOptions list filters
type FeatureFlagListOptions struct {
	*ListOptions
}

// NewFeatureFlagListOptions creates new options to pass to list
func NewFeatureFlagListOptions() *FeatureFlagListOptions {
	return &FeatureFlagListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o FeatureFlagListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Get the specified feature flag
func (c *FeatureFlagClient) Get(ctx context.Context, featureFlag resource.FeatureFlagType) (*resource.FeatureFlag, error) {
	var ff resource.FeatureFlag
	err := c.client.get(ctx, path.Format("/v3/feature_flags/%s", featureFlag), &ff)
	if err != nil {
		return nil, err
	}
	return &ff, nil
}

// List pages feature flags
func (c *FeatureFlagClient) List(ctx context.Context, opts *FeatureFlagListOptions) ([]*resource.FeatureFlag, *Pager, error) {
	if opts == nil {
		opts = NewFeatureFlagListOptions()
	}
	var res resource.FeatureFlagList
	err := c.client.list(ctx, "/v3/feature_flags", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all feature flags
func (c *FeatureFlagClient) ListAll(ctx context.Context, opts *FeatureFlagListOptions) ([]*resource.FeatureFlag, error) {
	if opts == nil {
		opts = NewFeatureFlagListOptions()
	}
	return AutoPage[*FeatureFlagListOptions, *resource.FeatureFlag](opts, func(opts *FeatureFlagListOptions) ([]*resource.FeatureFlag, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the feature flag
func (c *FeatureFlagClient) Update(ctx context.Context, featureFlag resource.FeatureFlagType, r *resource.FeatureFlagUpdate) (*resource.FeatureFlag, error) {
	var d resource.FeatureFlag
	_, err := c.client.patch(ctx, path.Format("/v3/feature_flags/%s", featureFlag), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
