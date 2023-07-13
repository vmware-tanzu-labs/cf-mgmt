package resource

import "time"

type Revision struct {
	GUID          string            `json:"guid"`
	Version       int               `json:"version"`
	Droplet       Relationship      `json:"droplet"`
	Processes     RevisionProcesses `json:"processes"`
	Sidecars      []RevisionSidecar `json:"sidecars"`
	Description   string            `json:"description"`
	Deployable    bool              `json:"deployable"`
	Relationships AppRelationship   `json:"relationships"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Metadata      *Metadata         `json:"metadata"`
	Links         map[string]Link   `json:"links"`
}

type RevisionList struct {
	Pagination Pagination  `json:"pagination"`
	Resources  []*Revision `json:"resources"`
}

type RevisionUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type RevisionProcess struct {
	Command string `json:"command"`
}

type RevisionProcesses struct {
	Web    *RevisionProcess `json:"web,omitempty"`
	Worker *RevisionProcess `json:"worker,omitempty"`
}

type RevisionSidecar struct {
	Name         string   `json:"name"`
	Command      string   `json:"command"`
	ProcessTypes []string `json:"process_types"`
	MemoryInMB   int      `json:"memory_in_mb"`
}
