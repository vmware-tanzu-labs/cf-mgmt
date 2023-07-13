package resource

import "time"

type IsolationSegment struct {
	GUID      string          `json:"guid"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Links     map[string]Link `json:"links"`
	Metadata  *Metadata       `json:"metadata"`
}

type IsolationSegmentCreate struct {
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type IsolationSegmentUpdate struct {
	Name     *string   `json:"name,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

type IsolationSegmentRelationship struct {
	Data  []Relationship  `json:"data"`
	Links map[string]Link `json:"links"`
}

type IsolationSegmentList struct {
	Pagination Pagination          `json:"pagination"`
	Resources  []*IsolationSegment `json:"resources"`
}

func NewIsolationSegmentCreate(name string) *IsolationSegmentCreate {
	return &IsolationSegmentCreate{
		Name: name,
	}
}
