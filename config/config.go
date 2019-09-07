// Package config provides utilities for reading and writing cf-mgmt's configuration.
package config

// DefaultProtectedOrgs lists the organizations that are considered protected
// and should never be deleted by cf-mgmt.
var DefaultProtectedOrgs = []string{
	"system",
	"p-spring-cloud-services",
	"splunk-nozzle-org",
	"redis-test-ORG*",
	"appdynamics-org",
	"credhub-service-broker-org",
}

// Manager can read and write the cf-mgmt configuration.
type Manager interface {
	Updater
	Reader
}

// Updater is used to update the cf-mgmt configuration.
type Updater interface {
	AddOrgToConfig(orgConfig *OrgConfig) error
	AddSpaceToConfig(spaceConfig *SpaceConfig) error
	AddSecurityGroupToSpace(orgName, spaceName string, securityGroupDefinition []byte) error
	AddSecurityGroup(securityGroupName string, securityGroupDefinition []byte) error
	AddOrgQuota(orgQuota OrgQuota) error
	AddSpaceQuota(spaceQuota SpaceQuota) error
	AddDefaultSecurityGroup(securityGroupName string, securityGroupDefinition []byte) error
	CreateConfigIfNotExists(uaaOrigin string) error
	DeleteConfigIfExists() error

	SaveOrgSpaces(spaces *Spaces) error
	SaveSpaceConfig(spaceConfig *SpaceConfig) error
	SaveOrgConfig(orgConfig *OrgConfig) error
	RenameOrgConfig(orgConfig *OrgConfig) error
	RenameSpaceConfig(spaceConfig *SpaceConfig) error

	DeleteOrgConfig(orgName string) error
	DeleteSpaceConfig(orgName, spaceName string) error

	SaveOrgs(*Orgs) error
	SaveGlobalConfig(*GlobalConfig) error
	SaveOrgQuota(*OrgQuota) error
	SaveSpaceQuota(*SpaceQuota) error
}

// Reader is used to read the cf-mgmt configuration.
type Reader interface {
	Orgs() (*Orgs, error)
	OrgSpaces(orgName string) (*Spaces, error)
	Spaces() ([]Spaces, error)
	GetOrgConfigs() ([]OrgConfig, error)
	GetSpaceConfigs() ([]SpaceConfig, error)
	GetASGConfigs() ([]ASGConfig, error)
	GetDefaultASGConfigs() ([]ASGConfig, error)
	GetGlobalConfig() (*GlobalConfig, error)
	GetSpaceDefaults() (*SpaceConfig, error)
	GetOrgConfig(orgName string) (*OrgConfig, error)
	GetSpaceConfig(orgName, spaceName string) (*SpaceConfig, error)
	LdapConfig(bindUser, bindPassword, ldapServer string) (*LdapConfig, error)
	GetOrgQuotas() ([]OrgQuota, error)
	GetSpaceQuotas(org string) ([]SpaceQuota, error)
	GetOrgQuota(name string) (*OrgQuota, error)
	GetSpaceQuota(name, org string) (*SpaceQuota, error)
}

// NewManager creates a Manager that is backed by a set of YAML
// files in the specified configuration directory.
func NewManager(configDir string) Manager {
	return &yamlManager{
		ConfigDir: configDir,
	}
}
