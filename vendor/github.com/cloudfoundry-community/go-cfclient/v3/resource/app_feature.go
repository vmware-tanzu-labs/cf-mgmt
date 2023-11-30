package resource

type AppFeature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type AppFeatureUpdate struct {
	Enabled bool `json:"enabled"`
}

type AppFeatureList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*AppFeature `json:"resources"`
}
