package resource

type ResourceMatches struct {
	Resources []ResourceMatch `json:"resources"`
}

type ResourceMatchChecksum struct {
	Value string `json:"value"`
}

type ResourceMatch struct {
	Checksum ResourceMatchChecksum `json:"checksum"`

	SizeInBytes int    `json:"size_in_bytes"`
	Path        string `json:"path"`
	Mode        string `json:"mode"`
}
