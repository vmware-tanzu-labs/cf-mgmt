package client

import "context"

type AdminClient commonClient

// ClearBuildpackCache will delete all the existing buildpack caches in the blobstore. Success returns a JobID.
//
// The buildpack cache is used during staging by buildpacks as a way to cache certain resources, e.g. downloaded
// Ruby gems. An admin who wants to decrease the size of their blobstore could use this endpoint to delete
// unnecessary blobs.
func (c *AdminClient) ClearBuildpackCache(ctx context.Context) (string, error) {
	jobGUID, err := c.client.post(ctx, "/v3/admin/actions/clear_buildpack_cache", nil, nil)
	if err != nil {
		return "", err
	}
	return jobGUID, nil
}
