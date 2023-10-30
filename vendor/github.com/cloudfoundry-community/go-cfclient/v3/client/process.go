package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type ProcessClient commonClient

// ProcessListOptions list filters
type ProcessListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Types             Filter `qs:"types"`
	Names             Filter `qs:"names"`
	AppGUIDs          Filter `qs:"app_guids"`
	SpaceGUIDs        Filter `qs:"space_guids"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewProcessOptions creates new options to pass to list
func NewProcessOptions() *ProcessListOptions {
	return &ProcessListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o ProcessListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// First returns the first process matching the options or an error when less than 1 match
func (c *ProcessClient) First(ctx context.Context, opts *ProcessListOptions) (*resource.Process, error) {
	return First[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForApp returns the first process matching the options and app or an error when less than 1 match
func (c *ProcessClient) FirstForApp(ctx context.Context, appGUID string, opts *ProcessListOptions) (*resource.Process, error) {
	return First[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Get the specified process
func (c *ProcessClient) Get(ctx context.Context, guid string) (*resource.Process, error) {
	var iso resource.Process
	err := c.client.get(ctx, path.Format("/v3/processes/%s", guid), &iso)
	if err != nil {
		return nil, err
	}
	return &iso, nil
}

// GetStats for the specified process
func (c *ProcessClient) GetStats(ctx context.Context, guid string) (*resource.ProcessStats, error) {
	var stats resource.ProcessStats
	err := c.client.get(ctx, path.Format("/v3/processes/%s/stats", guid), &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetStatsForApp for the specified app
func (c *ProcessClient) GetStatsForApp(ctx context.Context, appGUID, processType string) (*resource.ProcessStats, error) {
	var stats resource.ProcessStats
	err := c.client.get(ctx, path.Format("/v3/apps/%s/processes/%s/stats", appGUID, processType), &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// List pages all processes
func (c *ProcessClient) List(ctx context.Context, opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}

	var isos resource.ProcessList
	err := c.client.list(ctx, "/v3/processes", opts.ToQueryString, &isos)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(isos.Pagination)
	return isos.Resources, pager, nil
}

// ListAll retrieves all processes
func (c *ProcessClient) ListAll(ctx context.Context, opts *ProcessListOptions) ([]*resource.Process, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}
	return AutoPage[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForApp pages all processes for the specified app
func (c *ProcessClient) ListForApp(ctx context.Context, appGUID string, opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}

	var processes resource.ProcessList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/processes", opts.ToQueryString, &processes)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(processes.Pagination)
	return processes.Resources, pager, nil
}

// ListForAppAll retrieves all processes for the specified app
func (c *ProcessClient) ListForAppAll(ctx context.Context, appGUID string, opts *ProcessListOptions) ([]*resource.Process, error) {
	if opts == nil {
		opts = NewProcessOptions()
	}
	return AutoPage[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Scale the process using the specified scaling requirements
func (c *ProcessClient) Scale(ctx context.Context, guid string, scale *resource.ProcessScale) (*resource.Process, error) {
	var process resource.Process
	_, err := c.client.post(ctx, path.Format("/v3/processes/%s/actions/scale", guid), scale, &process)
	if err != nil {
		return nil, err
	}
	return &process, nil
}

// Single returns a single package matching the options or an error if not exactly 1 match
func (c *ProcessClient) Single(ctx context.Context, opts *ProcessListOptions) (*resource.Process, error) {
	return Single[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForApp returns a single package matching the options for the app or an error if not exactly 1 match
func (c *ProcessClient) SingleForApp(ctx context.Context, appGUID string, opts *ProcessListOptions) (*resource.Process, error) {
	return Single[*ProcessListOptions, *resource.Process](opts, func(opts *ProcessListOptions) ([]*resource.Process, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Update the specified attributes of the process
func (c *ProcessClient) Update(ctx context.Context, guid string, r *resource.ProcessUpdate) (*resource.Process, error) {
	var process resource.Process
	_, err := c.client.patch(ctx, path.Format("/v3/processes/%s", guid), r, &process)
	if err != nil {
		return nil, err
	}
	return &process, nil
}

// Terminate an instance of a specific process. Health management will eventually restart the instance.
func (c *ProcessClient) Terminate(ctx context.Context, guid string, index int) error {
	_, err := c.client.delete(ctx, path.Format("/v3/processes/%s/instances/%d", guid, index))
	return err
}
