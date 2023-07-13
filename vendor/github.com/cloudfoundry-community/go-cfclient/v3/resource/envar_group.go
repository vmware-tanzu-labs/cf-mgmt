package resource

import "time"

type EnvVarGroup struct {
	UpdatedAt time.Time         `json:"updated_at"`
	Name      string            `json:"name"` // The name of the group; can only be running or staging
	Var       map[string]string `json:"var"`
	Links     map[string]Link   `json:"links"`
}

type EnvVarGroupUpdate struct {
	Var map[string]string `json:"var"`
}
