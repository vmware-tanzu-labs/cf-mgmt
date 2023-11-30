package client

import (
	"context"
	"io"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type BuildpackClient commonClient

// BuildpackListOptions list filters
type BuildpackListOptions struct {
	*ListOptions

	Names  Filter `qs:"names"`  // list of buildpack names to filter by
	Stacks Filter `qs:"stacks"` // list of stack names to filter by
}

// NewBuildpackListOptions creates new options to pass to list
func NewBuildpackListOptions() *BuildpackListOptions {
	return &BuildpackListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o BuildpackListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Create a new buildpack
func (c *BuildpackClient) Create(ctx context.Context, r *resource.BuildpackCreateOrUpdate) (*resource.Buildpack, error) {
	var bp resource.Buildpack
	_, err := c.client.post(ctx, "/v3/buildpacks", r, &bp)
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

// Delete the specified buildpack
func (c *BuildpackClient) Delete(ctx context.Context, guid string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/buildpacks/%s", guid))
	return err
}

// First returns the first buildpack matching the options or an error when less than 1 match
func (c *BuildpackClient) First(ctx context.Context, opts *BuildpackListOptions) (*resource.Buildpack, error) {
	return First[*BuildpackListOptions, *resource.Buildpack](opts, func(opts *BuildpackListOptions) ([]*resource.Buildpack, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get retrieves the specified buildpack
func (c *BuildpackClient) Get(ctx context.Context, guid string) (*resource.Buildpack, error) {
	var bp resource.Buildpack
	err := c.client.get(ctx, path.Format("/v3/buildpacks/%s", guid), &bp)
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

// List pages all buildpacks the user has access to
func (c *BuildpackClient) List(ctx context.Context, opts *BuildpackListOptions) ([]*resource.Buildpack, *Pager, error) {
	if opts == nil {
		opts = NewBuildpackListOptions()
	}
	var res resource.BuildpackList
	err := c.client.list(ctx, "/v3/buildpacks", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all buildpacks the user has access to
func (c *BuildpackClient) ListAll(ctx context.Context, opts *BuildpackListOptions) ([]*resource.Buildpack, error) {
	if opts == nil {
		opts = NewBuildpackListOptions()
	}

	var all []*resource.Buildpack
	for {
		page, pager, err := c.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}
	return all, nil
}

// Single returns a single buildpack matching the options or an error if not exactly 1 match
func (c *BuildpackClient) Single(ctx context.Context, opts *BuildpackListOptions) (*resource.Buildpack, error) {
	return Single[*BuildpackListOptions, *resource.Buildpack](opts, func(opts *BuildpackListOptions) ([]*resource.Buildpack, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Update the specified attributes of the buildpack
func (c *BuildpackClient) Update(ctx context.Context, guid string, r *resource.BuildpackCreateOrUpdate) (*resource.Buildpack, error) {
	var bp resource.Buildpack
	_, err := c.client.patch(ctx, path.Format("/v3/buildpacks/%s", guid), r, &bp)
	if err != nil {
		return nil, err
	}
	return &bp, nil
}

// Upload a gzip compressed (zip) file containing a Cloud Foundry compatible buildpack
func (c *BuildpackClient) Upload(ctx context.Context, guid string, zipFile io.Reader) (string, *resource.Buildpack, error) {
	p := path.Format("/v3/buildpacks/%s/upload", guid)
	var b resource.Buildpack
	jobGUID, err := c.client.postFileUpload(ctx, p, "bits", "buildpack.zip", zipFile, &b)
	if err != nil {
		return "", nil, err
	}
	return jobGUID, &b, nil
}
