package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	http2 "net/http"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type PackageClient commonClient

// PackageListOptions list filters
type PackageListOptions struct {
	*ListOptions

	GUIDs  Filter `qs:"guids"`  // list of package guids to filter by
	States Filter `qs:"states"` // list of package states to filter by
	Types  Filter `qs:"types"`  // list of package types to filter by, docker or bits
}

// NewPackageListOptions creates new options to pass to list
func NewPackageListOptions() *PackageListOptions {
	return &PackageListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o PackageListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Copy the bits of a source package to a target package
func (c *PackageClient) Copy(ctx context.Context, srcPackageGUID string, destAppGUID string) (*resource.Package, error) {
	var d resource.Package
	r := resource.NewPackageCopy(destAppGUID)
	_, err := c.client.post(ctx, path.Format("/v3/packages?source_guid=%s", srcPackageGUID), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Create a new package
func (c *PackageClient) Create(ctx context.Context, r *resource.PackageCreate) (*resource.Package, error) {
	var p resource.Package
	_, err := c.client.post(ctx, "/v3/packages", r, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Delete the specified package asynchronously and return a jobGUID
func (c *PackageClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/packages/%s", guid))
}

// Download the bits of an existing package
// It is the caller's responsibility to close the io.ReadCloser
func (c *PackageClient) Download(ctx context.Context, guid string) (io.ReadCloser, error) {
	// This is the initial request, which will redirect to the blobstore location.
	// The client will not automatically follow this redirect and uses a secondary
	// unauthenticated client to download the bits
	// https://v3-apidocs.cloudfoundry.org/version/3.128.0/index.html#download-package-bits
	p := path.Format("/v3/packages/%s/download", guid)
	req := http.NewRequest(ctx, http2.MethodGet, p).WithFollowRedirects(false)
	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error getting %s: %w", p, err)
	}
	if !http.IsResponseRedirect(resp.StatusCode) {
		return nil, fmt.Errorf("error downloading package %s bits, expected redirect to blobstore", guid)
	}

	// get the full URL to the blobstore via the Location header
	blobStoreLocation := resp.Header.Get("Location")
	if blobStoreLocation == "" {
		return nil, errors.New("response redirect Location header was empty")
	}

	// directly download the bits from blobstore using an unauthenticated client as
	// some blob stores will return a 400 if an Authorization header is sent
	req = http.NewRequest(ctx, http2.MethodGet, "")
	blobstoreHTTPExecutor := http.NewExecutor(
		c.client.unauthenticatedClientProvider, blobStoreLocation, c.client.config.UserAgent)

	resp, err = blobstoreHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error downloading package %s bits from blobstore", guid)
	}

	return resp.Body, nil
}

// First returns the first package matching the options or an error when less than 1 match
func (c *PackageClient) First(ctx context.Context, opts *PackageListOptions) (*resource.Package, error) {
	return First[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForApp returns the first package matching the options and app or an error when less than 1 match
func (c *PackageClient) FirstForApp(ctx context.Context, appGUID string, opts *PackageListOptions) (*resource.Package, error) {
	return First[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Get the specified build
func (c *PackageClient) Get(ctx context.Context, guid string) (*resource.Package, error) {
	var p resource.Package
	err := c.client.get(ctx, path.Format("/v3/packages/%s", guid), &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// List pages all the packages the user has access to
func (c *PackageClient) List(ctx context.Context, opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	var res resource.PackageList
	err := c.client.list(ctx, "/v3/packages", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all the packages the user has access to
func (c *PackageClient) ListAll(ctx context.Context, opts *PackageListOptions) ([]*resource.Package, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	return AutoPage[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForApp pages all the packages the user has access to
func (c *PackageClient) ListForApp(ctx context.Context, appGUID string, opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	var res resource.PackageList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/packages", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all the packages the user has access to
func (c *PackageClient) ListForAppAll(ctx context.Context, appGUID string, opts *PackageListOptions) ([]*resource.Package, error) {
	if opts == nil {
		opts = NewPackageListOptions()
	}
	return AutoPage[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// PollReady waits until the package is ready, fails, or times out
func (c *PackageClient) PollReady(ctx context.Context, guid string, opts *PollingOptions) error {
	return PollForStateOrTimeout(func() (string, error) {
		pkg, err := c.Get(ctx, guid)
		if pkg != nil {
			return string(pkg.State), err
		}
		return "", err
	}, string(resource.PackageStateReady), opts)
}

// Single returns a single package matching the options or an error if not exactly 1 match
func (c *PackageClient) Single(ctx context.Context, opts *PackageListOptions) (*resource.Package, error) {
	return Single[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForApp returns a single package matching the options for the app or an error if not exactly 1 match
func (c *PackageClient) SingleForApp(ctx context.Context, appGUID string, opts *PackageListOptions) (*resource.Package, error) {
	return Single[*PackageListOptions, *resource.Package](opts, func(opts *PackageListOptions) ([]*resource.Package, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Update the specified attributes of the package
func (c *PackageClient) Update(ctx context.Context, guid string, r *resource.PackageUpdate) (*resource.Package, error) {
	var p resource.Package
	_, err := c.client.patch(ctx, path.Format("/v3/packages/%s", guid), r, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// Upload an app's zip file contents
func (c *PackageClient) Upload(ctx context.Context, guid string, zipFile io.Reader) (*resource.Package, error) {
	p := path.Format("/v3/packages/%s/upload", guid)
	var pkg resource.Package
	_, err := c.client.postFileUpload(ctx, p, "bits", "package.zip", zipFile, &pkg)
	return &pkg, err
}
