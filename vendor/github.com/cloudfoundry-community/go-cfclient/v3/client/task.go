package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type TaskClient commonClient

// TaskListOptions list filters
type TaskListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`
	Names             Filter `qs:"names"`
	States            Filter `qs:"states"`
	SpaceGUIDs        Filter `qs:"space_guids"`
	OrganizationGUIDs Filter `qs:"organization_guids"`
}

// NewTaskListOptions creates new options to pass to list
func NewTaskListOptions() *TaskListOptions {
	return &TaskListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o TaskListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// Cancel the specified task
//
// Canceled tasks will initially be in state CANCELING and will move to state FAILED once the cancel request
// has been processed. Cancel requests are idempotent and will be processed according to the state of the
// task when the request is executed. Canceling a task that is in SUCCEEDED or FAILED state will return an error.
func (c *TaskClient) Cancel(ctx context.Context, guid string) (*resource.Task, error) {
	var task resource.Task
	_, err := c.client.post(ctx, path.Format("/v3/tasks/%s/actions/cancel", guid), nil, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// Create a new task for the specified app
func (c *TaskClient) Create(ctx context.Context, appGUID string, r *resource.TaskCreate) (*resource.Task, error) {
	var task resource.Task
	_, err := c.client.post(ctx, path.Format("/v3/apps/%s/tasks", appGUID), r, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// First returns the first task matching the options or an error when less than 1 match
func (c *TaskClient) First(ctx context.Context, opts *TaskListOptions) (*resource.Task, error) {
	return First[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// FirstForApp returns the first task matching the options and app or an error when less than 1 match
func (c *TaskClient) FirstForApp(ctx context.Context, appGUID string, opts *TaskListOptions) (*resource.Task, error) {
	return First[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Get the specified task
func (c *TaskClient) Get(ctx context.Context, guid string) (*resource.Task, error) {
	var task resource.Task
	err := c.client.get(ctx, path.Format("/v3/tasks/%s", guid), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List pages all the tasks the user has access to. The command field is excluded in the response.
func (c *TaskClient) List(ctx context.Context, opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
	if opts == nil {
		opts = NewTaskListOptions()
	}

	var res resource.TaskList
	err := c.client.list(ctx, "/v3/tasks", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all tasks the user has access to. The command field is excluded in the response.
func (c *TaskClient) ListAll(ctx context.Context, opts *TaskListOptions) ([]*resource.Task, error) {
	if opts == nil {
		opts = NewTaskListOptions()
	}
	return AutoPage[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListForApp pages all the tasks for the specified app that the user has access to. The command field
// may be excluded in the response based on the user’s role.
func (c *TaskClient) ListForApp(ctx context.Context, appGUID string, opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
	if opts == nil {
		opts = NewTaskListOptions()
	}

	var res resource.TaskList
	err := c.client.list(ctx, "/v3/apps/"+appGUID+"/tasks", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListForAppAll retrieves all the tasks for the specified app that the user has access to. The command field
// may be excluded in the response based on the user’s role.
func (c *TaskClient) ListForAppAll(ctx context.Context, appGUID string, opts *TaskListOptions) ([]*resource.Task, error) {
	if opts == nil {
		opts = NewTaskListOptions()
	}
	return AutoPage[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Single returns a single task matching the options or an error if not exactly 1 match
func (c *TaskClient) Single(ctx context.Context, opts *TaskListOptions) (*resource.Task, error) {
	return Single[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// SingleForApp returns a single task matching the options or an error if not exactly 1 match
func (c *TaskClient) SingleForApp(ctx context.Context, appGUID string, opts *TaskListOptions) (*resource.Task, error) {
	return Single[*TaskListOptions, *resource.Task](opts, func(opts *TaskListOptions) ([]*resource.Task, *Pager, error) {
		return c.ListForApp(ctx, appGUID, opts)
	})
}

// Update the specified attributes of the task
func (c *TaskClient) Update(ctx context.Context, guid string, r *resource.TaskUpdate) (*resource.Task, error) {
	var task resource.Task
	_, err := c.client.patch(ctx, path.Format("/v3/tasks/%s", guid), r, &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}
