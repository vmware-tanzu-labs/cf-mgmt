package resource

type SpaceFeature struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Description string `json:"description"`
}

type SpaceFeatureUpdate struct {
	Enabled bool `json:"enabled"`
}
