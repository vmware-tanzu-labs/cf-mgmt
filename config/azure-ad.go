package config

// Config -
type AzureADConfig struct {
	Enabled    bool   `yaml:"enabled"`
	ClientId   string `yaml:"client-id"`
	Secret     string `yaml:"client-secret"`
	TenantID   string `yaml:"tenant-id"`
	UserOrigin string `yaml:"origin"`
}
