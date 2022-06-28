package cfclient

// V3User implements the user object
type V3User struct {
	GUID             string          `json:"guid,omitempty"`
	CreatedAt        string          `json:"created_at,omitempty"`
	UpdatedAt        string          `json:"updated_at,omitempty"`
	Username         string          `json:"username,omitempty"`
	PresentationName string          `json:"presentation_name,omitempty"`
	Origin           string          `json:"origin,omitempty"`
	Links            map[string]Link `json:"links,omitempty"`
	Metadata         V3Metadata      `json:"metadata,omitempty"`
}
