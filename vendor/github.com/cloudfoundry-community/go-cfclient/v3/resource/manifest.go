package resource

type ManifestDiff struct {
	Diff []ManifestDiffItem `json:"diff"`
}

type ManifestDiffItem struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Was   string `json:"was,omitempty"`
	Value string `json:"value,omitempty"`
}
