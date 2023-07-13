package resource

import "time"

type Organization struct {
	GUID          string            `json:"guid"`
	Name          string            `json:"name"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Suspended     *bool             `json:"suspended,omitempty"`
	Relationships QuotaRelationship `json:"relationships,omitempty"`
	Links         map[string]Link   `json:"links,omitempty"`
	Metadata      *Metadata         `json:"metadata,omitempty"`
}

type OrganizationCreate struct {
	Name      string    `json:"name"`
	Suspended *bool     `json:"suspended,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type OrganizationUpdate struct {
	Name      string    `json:"name,omitempty"`
	Suspended *bool     `json:"suspended,omitempty"`
	Metadata  *Metadata `json:"metadata,omitempty"`
}

type OrganizationUsageSummary struct {
	UsageSummary UsageSummary    `json:"usage_summary"`
	Links        map[string]Link `json:"links,omitempty"`
}

type OrganizationList struct {
	Pagination Pagination      `json:"pagination,omitempty"`
	Resources  []*Organization `json:"resources,omitempty"`
}

type UsageSummary struct {
	StartedInstances int `json:"started_instances"`
	MemoryInMb       int `json:"memory_in_mb"`
}

type QuotaRelationship struct {
	Quota ToOneRelationship `json:"quota"`
}

func NewOrganizationCreate(name string) *OrganizationCreate {
	return &OrganizationCreate{
		Name: name,
	}
}
