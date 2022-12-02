package resource

import "time"

type Task struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`

	// User-facing id of the task; this number is unique for every task associated with a given app
	SequenceID int `json:"sequence_id"`

	// Command that will be executed; this field may be excluded based on a user’s role
	Command string `json:"command"`

	// State of the task Possible states are PENDING, RUNNING, SUCCEEDED, CANCELING, and FAILED
	State string `json:"state"`

	// Amount of memory to allocate for the task in MB
	MemoryInMB int `json:"memory_in_mb"`

	// Amount of disk to allocate for the task in MB
	DiskInMB int `json:"disk_in_mb"`

	// Amount of log rate to allocate for the task in bytes
	LogRateLimitInBytesPerSecond int `json:"log_rate_limit_in_bytes_per_second"`

	// Results from the task
	Result TaskResult `json:"result"`

	// The guid of the droplet that will be used to run the command
	DropletGUID string `json:"droplet_guid"`

	// The app the task belongs to
	Relationships AppRelationship `json:"relationships"`

	Links    map[string]Link `json:"links"`
	Metadata *Metadata       `json:"metadata"`
}

type TaskCreate struct {
	// Command that will be executed; NOTE: optional if a template.process.guid is provided
	Command *string `json:"command,omitempty"`

	// Name of the task, otherwise auto-generated
	Name *string `json:"name,omitempty"`

	// Amount of memory to allocate for the task in MB
	MemoryInMB *int `json:"memory_in_mb,omitempty"`

	// Amount of disk to allocate for the task in MB
	DiskInMB *int `json:"disk_in_mb,omitempty"`

	// Amount of log rate to allocate for the task in bytes
	LogRateLimitInBytesPerSecond *int `json:"log_rate_limit_in_bytes_per_second,omitempty"`

	// The guid of the droplet that will be used to run the command
	DropletGUID *string `json:"droplet_guid,omitempty"`

	// The process that will be used as a template
	Template *TaskTemplate `json:"template,omitempty"`

	Metadata *Metadata `json:"metadata,omitempty"`
}

type TaskUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type TaskList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Task    `json:"resources"`
}

type TaskResult struct {
	// nil if the task succeeds, contains the error message if it fails
	FailureReason *string `json:"failure_reason"`
}

type TaskTemplate struct {
	Process TaskProcess `json:"process"`
}

type TaskProcess struct {
	// The guid of the process that will be used as a template
	GUID string `json:"guid"`
}

func NewTaskCreateWithProcessTemplate(processGUID string) *TaskCreate {
	return &TaskCreate{
		Template: &TaskTemplate{
			Process: TaskProcess{
				GUID: processGUID,
			},
		},
	}
}

func NewTaskCreateWithCommand(command string) *TaskCreate {
	return &TaskCreate{
		Command: &command,
	}
}

func (t *TaskCreate) WithName(name string) *TaskCreate {
	t.Name = &name
	return t
}

func (t *TaskCreate) WithMemoryInMB(mb int) *TaskCreate {
	t.MemoryInMB = &mb
	return t
}

func (t *TaskCreate) WithDiskInMB(mb int) *TaskCreate {
	t.DiskInMB = &mb
	return t
}

func (t *TaskCreate) WithLogRateLimitInBytesPerSecond(mbps int) *TaskCreate {
	t.LogRateLimitInBytesPerSecond = &mbps
	return t
}

func (t *TaskCreate) WithDropletGUID(dropletGUID string) *TaskCreate {
	t.DropletGUID = &dropletGUID
	return t
}
