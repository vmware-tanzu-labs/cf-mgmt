package resource

import "time"

type Buildpack struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`     // The name of the buildpack; to be used by app buildpack field (only alphanumeric characters)
	State     string    `json:"state"`    // The state of the buildpack Valid value is: AWAITING_UPLOAD
	Filename  *string   `json:"filename"` // The filename of the buildpack, if any
	Stack     string    `json:"stack"`    // The name of the stack that the buildpack will use
	Position  int       `json:"position"` // The order in which the buildpacks are checked during buildpack auto-detection
	Enabled   bool      `json:"enabled"`  // Whether the buildpack can be used for staging
	Locked    bool      `json:"locked"`   // Whether the buildpack is locked to prevent updating the bits

	Metadata *Metadata       `json:"metadata"`
	Links    map[string]Link `json:"links"`
}

type BuildpackCreateOrUpdate struct {
	Name     *string   `json:"name,omitempty"`     // The name of the buildpack; to be used by app buildpack field (only alphanumeric characters)
	Position *int      `json:"position,omitempty"` // The order in which the buildpacks are checked during buildpack auto-detection
	Enabled  *bool     `json:"enabled,omitempty"`  // Whether the buildpack can be used for staging
	Locked   *bool     `json:"locked,omitempty"`   // Whether the buildpack is locked to prevent updating the bits
	Stack    *string   `json:"stack,omitempty"`    // The name of the stack that the buildpack will use
	Metadata *Metadata `json:"metadata,omitempty"`
}

type BuildpackList struct {
	Pagination Pagination   `json:"pagination"`
	Resources  []*Buildpack `json:"resources"`
}

func NewBuildpackCreate(name string) *BuildpackCreateOrUpdate {
	return &BuildpackCreateOrUpdate{
		Name: &name,
	}
}

func (bp *BuildpackCreateOrUpdate) WithName(name string) *BuildpackCreateOrUpdate {
	bp.Name = &name
	return bp
}

func (bp *BuildpackCreateOrUpdate) WithPosition(position int) *BuildpackCreateOrUpdate {
	bp.Position = &position
	return bp
}

func (bp *BuildpackCreateOrUpdate) WithStack(stack string) *BuildpackCreateOrUpdate {
	bp.Stack = &stack
	return bp
}

func (bp *BuildpackCreateOrUpdate) WithEnabled(enabled bool) *BuildpackCreateOrUpdate {
	bp.Enabled = &enabled
	return bp
}

func (bp *BuildpackCreateOrUpdate) WithLocked(locked bool) *BuildpackCreateOrUpdate {
	bp.Locked = &locked
	return bp
}

func NewBuildpackUpdate() *BuildpackCreateOrUpdate {
	return &BuildpackCreateOrUpdate{}
}
