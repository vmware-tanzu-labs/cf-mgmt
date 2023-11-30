package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SidecarClient commonClient

// SidecarListOptions list filters
type SidecarListOptions struct {
	*ListOptions
}

// NewSidecarListOptions creates new options to pass to list
func NewSidecarListOptions() *SidecarListOptions {
	return &SidecarListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SidecarListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new app sidecar
func (c *SidecarClient) Create(ctx context.Context, appGUID string, r *resource.SidecarCreate) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	_, err := c.client.post(ctx, path.Format("/v3/apps/%s/sidecars", appGUID), r, &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

// Delete the specified sidecar
func (c *SidecarClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/sidecars/%s", guid))
	return err
}

// FirstForApp returns the first sidecar matching the options and app or an error when less than 1 match
func (c *SidecarClient) FirstForApp(ctx context.Context, appGUID string, opts *SidecarListOptions) (*resource.Sidecar, error) {
	return First[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// FirstForProcess returns the first sidecar matching the options and process or an error when less than 1 match
func (c *SidecarClient) FirstForProcess(ctx context.Context, processGUID string, opts *SidecarListOptions) (*resource.Sidecar, error) {
	return First[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForProcess(ctx, processGUID, opts)
	})
}

// Get the specified app
func (c *SidecarClient) Get(ctx context.Context, guid string) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	err := c.client.get(ctx, path.Format("/v3/sidecars/%s", guid), &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}

// ListForApp pages all sidecars associated with the specified app
func (c *SidecarClient) ListForApp(ctx context.Context, appGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	var res resource.SidecarList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/sidecars", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all sidecars associated with the specified app
func (c *SidecarClient) ListForAppAll(ctx context.Context, appGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	return AutoPage[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// ListForProcess pages all sidecars associated with the specified process
func (c *SidecarClient) ListForProcess(ctx context.Context, processGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	var res resource.SidecarList
	err := c.client.list(ctx, "/v3/processes/"+processGUID+"/sidecars", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForProcessAll retrieves all sidecars associated with the specified process
func (c *SidecarClient) ListForProcessAll(ctx context.Context, processGUID string, opts *SidecarListOptions) ([]*resource.Sidecar, error) {
	if opts == nil {
		opts = NewSidecarListOptions()
	}
	return AutoPage[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForProcess(ctx, processGUID, opts)
	})
}

// SingleForApp returns a single sidecar matching the options and app or an error if not exactly 1 match
func (c *SidecarClient) SingleForApp(ctx context.Context, appGUID string, opts *SidecarListOptions) (*resource.Sidecar, error) {
	return Single[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// SingleForProcess returns a single sidecar matching the options and process or an error if not exactly 1 match
func (c *SidecarClient) SingleForProcess(ctx context.Context, processGUID string, opts *SidecarListOptions) (*resource.Sidecar, error) {
	return Single[*SidecarListOptions, *resource.Sidecar](opts, func(opts *SidecarListOptions) ([]*resource.Sidecar, *Pager, error) {
		return c.ListForProcess(ctx, processGUID, opts)
	})
}

// Update the specified attributes of the app
func (c *SidecarClient) Update(ctx context.Context, guid string, r *resource.SidecarUpdate) (*resource.Sidecar, error) {
	var sc resource.Sidecar
	_, err := c.client.patch(ctx, path.Format("/v3/sidecars/%s", guid), r, &sc)
	if err != nil {
		return nil, err
	}
	return &sc, nil
}
