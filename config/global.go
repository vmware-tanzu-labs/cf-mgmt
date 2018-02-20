package config

// GlobalConfig configuration for global settings
type GlobalConfig struct {
	EnableDeleteIsolationSegments bool     `yaml:"enable-delete-isolation-segments"`
	EnableUnassignSecurityGroups  bool     `yaml:"enable-unassign-security-groups"`
	RunningSecurityGroups         []string `yaml:"running-security-groups"`
	StagingSecurityGroups         []string `yaml:"staging-security-groups"`
}
