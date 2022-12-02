package client

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type JobClient commonClient

// Get the specified job
func (c *JobClient) Get(ctx context.Context, guid string) (*resource.Job, error) {
	var job resource.Job
	err := c.client.get(ctx, path.Format("/v3/jobs/%s", guid), &job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// PollComplete waits until the job completes, fails, or times out
func (c *JobClient) PollComplete(ctx context.Context, jobGUID string, opts *PollingOptions) error {
	err := PollForStateOrTimeout(func() (string, error) {
		job, err := c.Get(ctx, jobGUID)
		if job != nil {
			return string(job.State), err
		}
		return "", err
	}, string(resource.JobStateComplete), opts)

	// attempt to return the underlying saved job error
	if err == AsyncProcessFailedError {
		job, _ := c.Get(ctx, jobGUID)
		if job != nil && len(job.Errors) > 0 {
			return job.Errors[0]
		}
	}
	return err
}
