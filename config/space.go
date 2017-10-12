package config

import "fmt"

// Spaces describes cf-mgmt config for all spaces.
type Spaces struct {
	Org                string   `yaml:"org"`
	Spaces             []string `yaml:"spaces"`
	EnableDeleteSpaces bool     `yaml:"enable-delete-spaces"`
}

// SpaceConfig describes attributes for a space.
type SpaceConfig struct {
	Org                     string   `yaml:"org"`
	Space                   string   `yaml:"space"`
	Developer               UserMgmt `yaml:"space-developer"`
	Manager                 UserMgmt `yaml:"space-manager"`
	Auditor                 UserMgmt `yaml:"space-auditor"`
	DeveloperGroup          string   `yaml:"space-developer-group,omitempty"`
	ManagerGroup            string   `yaml:"space-manager-group,omitempty"`
	AuditorGroup            string   `yaml:"space-auditor-group,omitempty"`
	AllowSSH                bool     `yaml:"allow-ssh"`
	EnableSpaceQuota        bool     `yaml:"enable-space-quota"`
	MemoryLimit             int      `yaml:"memory-limit"`
	InstanceMemoryLimit     int      `yaml:"instance-memory-limit"`
	TotalRoutes             int      `yaml:"total-routes"`
	TotalServices           int      `yaml:"total-services"`
	PaidServicePlansAllowed bool     `yaml:"paid-service-plans-allowed"`
	EnableSecurityGroup     bool     `yaml:"enable-security-group"`
	SecurityGroupContents   string   `yaml:"security-group-contents,omitempty"`
	RemoveUsers             bool     `yaml:"enable-remove-users"`
	TotalPrivateDomains     int      `yaml:"total_private_domains"`
	TotalReservedRoutePorts int      `yaml:"total_reserved_route_ports"`
	TotalServiceKeys        int      `yaml:"total_service_keys"`
	AppInstanceLimit        int      `yaml:"app_instance_limit"`
	IsoSegment              string   `yaml:"isolation_segment"`
}

// Contains determines whether a space is present in a list of spaces.
func (s *Spaces) Contains(spaceName string) bool {
	for _, v := range s.Spaces {
		if v == spaceName {
			return true
		}
	}
	return false
}

func (i *SpaceConfig) GetSpaceConfigFilenameAndPath(configDir, orgName, spaceName string) string {
	return fmt.Sprintf("%s/%s/%s/spaceConfig.yml", configDir, orgName, spaceName)
}

func (i *SpaceConfig) GetSpaceConfigFilePath(configDir, orgName, spaceName string) string {
	return fmt.Sprintf("%s/%s/%s", configDir, orgName, spaceName)
}

func (i *SpaceConfig) GetDeveloperGroups() []string {
	return i.Developer.groups(i.DeveloperGroup)
}

func (i *SpaceConfig) GetManagerGroups() []string {
	return i.Manager.groups(i.ManagerGroup)
}

func (i *SpaceConfig) GetAuditorGroups() []string {
	return i.Auditor.groups(i.AuditorGroup)
}
