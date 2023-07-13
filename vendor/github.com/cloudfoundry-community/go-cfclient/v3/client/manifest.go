package client

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/http"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"io"
	http2 "net/http"
	"strings"
)

type ManifestClient commonClient

// Generate the specified app manifest as a yaml text string
func (c *ManifestClient) Generate(ctx context.Context, appGUID string) (string, error) {
	p := path.Format("/v3/apps/%s/manifest", appGUID)
	req := http.NewRequest(ctx, http2.MethodGet, p)

	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error getting %s: %w", p, err)
	}
	defer func(b io.ReadCloser) {
		_ = b.Close()
	}(resp.Body)

	if resp.StatusCode != http2.StatusOK {
		return "", c.client.decodeError(resp)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response %s: %w", p, err)
	}
	return buf.String(), nil
}

// ApplyManifest applies the changes specified in a manifest to the named apps and their underlying processes
// asynchronously and returns a jobGUID.
//
// The apps must reside in the space. These changes are additive and will not modify any unspecified
// properties or remove any existing environment variables, routes, or services.
func (c *ManifestClient) ApplyManifest(ctx context.Context, spaceGUID string, manifest string) (string, error) {
	reader := strings.NewReader(manifest)
	req := http.NewRequest(ctx, http2.MethodPost, path.Format("/v3/spaces/%s/actions/apply_manifest", spaceGUID)).
		WithContentType("application/x-yaml").
		WithBody(reader)

	resp, err := c.client.authenticatedHTTPExecutor.ExecuteRequest(req)
	if err != nil {
		return "", fmt.Errorf("error uploading manifest %s bits: %w", spaceGUID, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http2.StatusAccepted {
		return "", c.client.decodeError(resp)
	}

	jobGUID, err := c.client.decodeJobIDOrBody(resp, nil)
	if err != nil {
		return "", fmt.Errorf("error reading jobGUID: %w", err)
	}
	return jobGUID, nil
}
