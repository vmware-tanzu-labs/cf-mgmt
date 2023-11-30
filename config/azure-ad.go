package config

// Config -
type AzureADConfig struct {
	Enabled    bool   `yaml:"enabled,omitempty"`
	ClientId   string `yaml:"client-id,omitempty"`
	Secret     string `yaml:"client-secret,omitempty"`
	TenantID   string `yaml:"tenant-id,omitempty"`
	UserOrigin string `yaml:"origin,omitempty"`
	SPNOrigin  string `yaml:"spn-origin,omitempty"`
}
