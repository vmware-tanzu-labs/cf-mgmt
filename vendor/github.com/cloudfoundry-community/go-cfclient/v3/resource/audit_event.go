package resource

import (
	"encoding/json"
	"time"
)

type AuditEvent struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Type      string    `json:"type"`

	Actor  AuditEventRelatedObject `json:"actor"`
	Target AuditEventRelatedObject `json:"target"`

	Data         *json.RawMessage `json:"data"`
	Space        Relationship     `json:"space"`
	Organization Relationship     `json:"organization"`

	Links map[string]Link `json:"links"`
}

type AuditEventList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*AuditEvent `json:"resources"`
}

type AuditEventRelatedObject struct {
	GUID string `json:"guid"`
	Type string `json:"type"`
	Name string `json:"name"`
}
