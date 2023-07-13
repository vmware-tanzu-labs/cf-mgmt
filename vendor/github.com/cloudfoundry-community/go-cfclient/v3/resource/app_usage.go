package resource

import "time"

type AppUsage struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// the app that this event pertains to, if applicable
	App AppUsageGUIDName `json:"app"`

	// the process that this event pertains to, if applicable
	Process AppUsageGUIDType `json:"process"`

	// the space that this event pertains to, if applicable
	Space AppUsageGUIDName `json:"space"`

	// the organization that this event pertains to, if applicable
	Organization Relationship `json:"organization"`

	// the buildpack that this event pertains to, if applicable
	Buildpack AppUsageGUIDName `json:"buildpack"`

	// the task that this event pertains to, if applicable
	Task AppUsageGUIDName `json:"task"`

	// state of the app that this event pertains to, if applicable
	State AppUsageCurrentPreviousString `json:"state"`

	// memory in MB of the app that this event pertains to, if applicable
	MemoryInMbPerInstance AppUsageCurrentPreviousInt `json:"memory_in_mb_per_instance"`

	// instance count of the app that this event pertains to, if applicable
	InstanceCount AppUsageCurrentPreviousInt `json:"instance_count"`

	Links map[string]Link `json:"links"`
}

type AppUsageList struct {
	Pagination Pagination  `json:"pagination"`
	Resources  []*AppUsage `json:"resources"`
}

type AppUsageCurrentPreviousString struct {
	Current  string `json:"current"`
	Previous string `json:"previous"`
}

type AppUsageCurrentPreviousInt struct {
	Current  int `json:"current"`
	Previous int `json:"previous"`
}

type AppUsageGUIDName struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type AppUsageGUIDType struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
}
