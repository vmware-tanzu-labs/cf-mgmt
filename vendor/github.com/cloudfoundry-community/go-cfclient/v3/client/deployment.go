package client

import (
	"context"
	"errors"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type DeploymentClient commonClient

// DeploymentListOptions list filters
type DeploymentListOptions struct {
	*ListOptions

	AppGUIDs      Filter `qs:"app_guids"`
	States        Filter `qs:"states"`
	StatusReasons Filter `qs:"status_reasons"`
	StatusValues  Filter `qs:"status_values"`
}

// NewDeploymentListOptions creates new options to pass to list
func NewDeploymentListOptions() *DeploymentListOptions {
	return &DeploymentListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o DeploymentListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Cancel the ongoing deployment
func (c *DeploymentClient) Cancel(ctx context.Context, guid string) error {
	_, err := c.client.post(ctx, path.Format("/v3/deployments/%s/actions/cancel", guid), nil, nil)
	return err
}

// Create a new deployment
func (c *DeploymentClient) Create(ctx context.Context, r *resource.DeploymentCreate) (*resource.Deployment, error) {
	// validate the params
	if r.Droplet != nil && r.Revision != nil {
		return nil, errors.New("droplet and revision cannot both be set")
	}

	var d resource.Deployment
	_, err := c.client.post(ctx, "/v3/deployments", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// First returns the first deployment matching the options or an error when less than 1 match
func (c *DeploymentClient) First(ctx context.Context, opts *DeploymentListOptions) (*resource.Deployment, error) {
	return First[*DeploymentListOptions, *resource.Deployment](opts, func(opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified deployment
func (c *DeploymentClient) Get(ctx context.Context, guid string) (*resource.Deployment, error) {
	var d resource.Deployment
	err := c.client.get(ctx, path.Format("/v3/deployments/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages deployments the user has access to
func (c *DeploymentClient) List(ctx context.Context, opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
	if opts == nil {
		opts = NewDeploymentListOptions()
	}
	var res resource.DeploymentList
	err := c.client.list(ctx, "/v3/deployments", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all deployments the user has access to
func (c *DeploymentClient) ListAll(ctx context.Context, opts *DeploymentListOptions) ([]*resource.Deployment, error) {
	if opts == nil {
		opts = NewDeploymentListOptions()
	}
	return AutoPage[*DeploymentListOptions, *resource.Deployment](opts, func(opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single deployment matching the options or an error if not exactly 1 match
func (c *DeploymentClient) Single(ctx context.Context, opts *DeploymentListOptions) (*resource.Deployment, error) {
	return Single[*DeploymentListOptions, *resource.Deployment](opts, func(opts *DeploymentListOptions) ([]*resource.Deployment, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the deployment
func (c *DeploymentClient) Update(ctx context.Context, guid string, r *resource.DeploymentUpdate) (*resource.Deployment, error) {
	var d resource.Deployment
	_, err := c.client.patch(ctx, path.Format("/v3/deployments/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
