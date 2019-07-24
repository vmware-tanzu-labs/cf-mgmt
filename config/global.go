package config

// GlobalConfig configuration for global settings
type GlobalConfig struct {
	EnableDeleteIsolationSegments bool                    `yaml:"enable-delete-isolation-segments"`
	EnableUnassignSecurityGroups  bool                    `yaml:"enable-unassign-security-groups"`
	EnableServiceAccess           bool                    `yaml:"enable-service-access"`
	RunningSecurityGroups         []string                `yaml:"running-security-groups"`
	StagingSecurityGroups         []string                `yaml:"staging-security-groups"`
	SharedDomains                 map[string]SharedDomain `yaml:"shared-domains"`
	EnableDeleteSharedDomains     bool                    `yaml:"enable-remove-shared-domains"`
	MetadataPrefix                string                  `yaml:"metadata-prefix"`
	ServiceAccess                 []ServiceVisibility     `yaml:"service-access"`
}

type SharedDomain struct {
	Internal    bool   `yaml:"internal"`
	RouterGroup string `yaml:"router-group,omitempty"`
}

type ServiceVisibility struct {
	Service string   `yaml:"service"`
	Plan    string   `yaml:"plan,omitempty"`
	Orgs    []string `yaml:"orgs,omitempty"`
	Disable bool     `yaml:"disable,omitempty"`
}
