package resource

import "time"

type Process struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Process type; a unique identifier for processes belonging to an app: web, task
	Type string `json:"type"`

	// The command used to start the process; use null to revert to the buildpack-detected or procfile-provided start command
	Command *string `json:"command"`

	// The number of instances to run
	Instances int `json:"instances"`

	// The memory in MB allocated per instance
	MemoryInMB int `json:"memory_in_mb"`

	// The disk in MB allocated per instance
	DiskInMB int `json:"disk_in_mb"`

	// The log rate in bytes per second allocated per instance
	LogRateLimitInBytesPerSecond int `json:"log_rate_limit_in_bytes_per_second"`

	HealthCheck   ProcessHealthCheck   `json:"health_check"`
	Relationships ProcessRelationships `json:"relationships"`

	Links    map[string]Link `json:"links"`
	Metadata *Metadata       `json:"metadata"`
}

type ProcessList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Process `json:"resources"`
}

type ProcessUpdate struct {
	// The command used to start the process
	// Note - This doesn't currently support using null to revert to the buildpack-detected or
	// procfile-provided start command because of the omitempty
	Command *string `json:"command,omitempty"`

	HealthCheck *ProcessHealthCheck `json:"health_check,omitempty"`
	Metadata    *Metadata           `json:"metadata,omitempty"`
}

type ProcessStats struct {
	Stats []ProcessStat `json:"resources"`
}

type ProcessStat struct {
	Type                string           `json:"type"`  // Process type; a unique identifier for processes belonging to an app
	Index               int              `json:"index"` // The zero-based index of running instances
	State               string           `json:"state"` // The state of the instance; valid values are RUNNING, CRASHED, STARTING, DOWN
	Usage               Usage            `json:"usage"`
	Host                string           `json:"host"`
	InstancePorts       []map[string]int `json:"instance_ports"`
	Uptime              int              `json:"uptime"`
	MemoryQuota         int              `json:"mem_quota"`
	DiskQuota           int              `json:"disk_quota"`
	FileDescriptorQuota int              `json:"fds_quota"`

	// The current isolation segment that the instance is running on; the value is null when the
	// instance is not placed on a particular isolation segment
	IsolationSegment *string `json:"isolation_segment"`

	// Information about errors placing the instance; the value is null if there are no placement errors
	Details *string `json:"details"`
}

type ProcessScale struct {
	// The number of instances to run
	Instances *int `json:"instances,omitempty"`

	// The memory in MB allocated per instance
	MemoryInMB *int `json:"memory_in_mb,omitempty"`

	// The disk in MB allocated per instance
	DiskInMB *int `json:"disk_in_mb,omitempty"`

	// The log rate in bytes per second allocated per instance
	LogRateLimitInBytesPerSecond *int `json:"log_rate_limit_in_bytes_per_second,omitempty"`
}

type ProcessHealthCheck struct {
	// The type of health check to perform; valid values are http, port, and process; default is port
	Type string      `json:"type"`
	Data ProcessData `json:"data"`
}

type ProcessData struct {
	// The duration in seconds that health checks can fail before the process is restarted
	Timeout *int `json:"timeout"`

	// The timeout in seconds for individual health check requests for http and port health checks
	InvocationTimeout *int `json:"invocation_timeout,omitempty"`

	// The endpoint called to determine if the app is healthy; this key is only present for http health check
	Endpoint *string `json:"endpoint,omitempty"`
}

type ProcessRelationships struct {
	App      ToOneRelationship `json:"app"`      // The app the process belongs to
	Revision ToOneRelationship `json:"revision"` // The app revision the process is currently running
}

type Usage struct {
	Time   time.Time `json:"time"`
	CPU    float64   `json:"cpu"`
	Memory int       `json:"mem"`
	Disk   int       `json:"disk"`
}

func NewProcessScale() *ProcessScale {
	return &ProcessScale{}
}

func (p *ProcessScale) WithInstances(count int) *ProcessScale {
	p.Instances = &count
	return p
}

func (p *ProcessScale) WithMemoryInMB(mb int) *ProcessScale {
	p.MemoryInMB = &mb
	return p
}

func (p *ProcessScale) WithDiskInMB(mb int) *ProcessScale {
	p.DiskInMB = &mb
	return p
}

func (p *ProcessScale) WithLogRateLimitInBytesPerSecond(rate int) *ProcessScale {
	p.LogRateLimitInBytesPerSecond = &rate
	return p
}

func NewProcessUpdate() *ProcessUpdate {
	return &ProcessUpdate{}
}

func (p *ProcessUpdate) WithCommand(cmd string) *ProcessUpdate {
	p.Command = &cmd
	return p
}

func (p *ProcessUpdate) WithHealthCheckType(hcType string) *ProcessUpdate {
	if p.HealthCheck == nil {
		p.HealthCheck = &ProcessHealthCheck{}
	}
	p.HealthCheck.Type = hcType
	return p
}

func (p *ProcessUpdate) WithHealthCheckTimeout(timeout int) *ProcessUpdate {
	if p.HealthCheck == nil {
		p.HealthCheck = &ProcessHealthCheck{}
	}
	p.HealthCheck.Data.Timeout = &timeout
	return p
}

func (p *ProcessUpdate) WithHealthCheckInvocationTimeout(timeout int) *ProcessUpdate {
	if p.HealthCheck == nil {
		p.HealthCheck = &ProcessHealthCheck{}
	}
	p.HealthCheck.Data.InvocationTimeout = &timeout
	return p
}

func (p *ProcessUpdate) WithHealthCheckEndpoint(endpoint string) *ProcessUpdate {
	if p.HealthCheck == nil {
		p.HealthCheck = &ProcessHealthCheck{}
	}
	p.HealthCheck.Data.Endpoint = &endpoint
	return p
}
