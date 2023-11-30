package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type RevisionClient commonClient

// RevisionListOptions list filters
type RevisionListOptions struct {
	*ListOptions

	Versions Filter `qs:"versions"`
}

// NewRevisionListOptions creates new options to pass to list
func NewRevisionListOptions() *RevisionListOptions {
	return &RevisionListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o RevisionListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// FirstForApp returns the first revision matching the options and app or an error when less than 1 match
func (c *RevisionClient) FirstForApp(ctx context.Context, appGUID string, opts *RevisionListOptions) (*resource.Revision, error) {
	return First[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Get the specified revision
func (c *RevisionClient) Get(ctx context.Context, guid string) (*resource.Revision, error) {
	var res resource.Revision
	err := c.client.get(ctx, path.Format("/v3/revisions/%s", guid), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// GetEnvironmentVariables retrieves the specified revision's environment variables
func (c *RevisionClient) GetEnvironmentVariables(ctx context.Context, guid string) (map[string]*string, error) {
	var res resource.EnvVarResponse
	err := c.client.get(ctx, path.Format("/v3/revisions/%s/environment_variables", guid), &res)
	if err != nil {
		return nil, err
	}
	return res.Var, nil
}

// ListForApp pages revisions that are associated with the specified app
func (c *RevisionClient) ListForApp(ctx context.Context, appGUID string, opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	var res resource.RevisionList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/revisions", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all revisions that are associated with the specified app
func (c *RevisionClient) ListForAppAll(ctx context.Context, appGUID string, opts *RevisionListOptions) ([]*resource.Revision, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	return AutoPage[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// ListForAppDeployed pages deployed revisions that are associated with the specified app
func (c *RevisionClient) ListForAppDeployed(ctx context.Context, appGUID string, opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	var res resource.RevisionList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/revisions/deployed", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppDeployedAll pages deployed revisions that are associated with the specified app
func (c *RevisionClient) ListForAppDeployedAll(ctx context.Context, appGUID string, opts *RevisionListOptions) ([]*resource.Revision, error) {
	if opts == nil {
		opts = NewRevisionListOptions()
	}
	return AutoPage[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListForAppDeployed(ctx, appGUID, opts)
	})
}

// SingleForApp returns a single revision matching the options and app or an error if not exactly 1 match
func (c *RevisionClient) SingleForApp(ctx context.Context, appGUID string, opts *RevisionListOptions) (*resource.Revision, error) {
	return Single[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// SingleForAppDeployed returns a single deployed revision matching the options and app or an error if not exactly 1 match
func (c *RevisionClient) SingleForAppDeployed(ctx context.Context, appGUID string, opts *RevisionListOptions) (*resource.Revision, error) {
	return Single[*RevisionListOptions, *resource.Revision](opts, func(opts *RevisionListOptions) ([]*resource.Revision, *Pager, error) {
		return c.ListForAppDeployed(ctx, appGUID, opts)
	})
}

// Update the specified attributes of the deployment
func (c *RevisionClient) Update(ctx context.Context, guid string, r *resource.RevisionUpdate) (*resource.Revision, error) {
	var res resource.Revision
	_, err := c.client.patch(ctx, path.Format("/v3/revisions/%s", guid), r, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
