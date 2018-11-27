package config

import (
	"strings"

	"github.com/xchapter7x/lo"
)

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
	OriginalSpace           string   `yaml:"original-space,omitempty"`
	Developer               UserMgmt `yaml:"space-developer"`
	Manager                 UserMgmt `yaml:"space-manager"`
	Auditor                 UserMgmt `yaml:"space-auditor"`
	DeveloperGroup          string   `yaml:"space-developer-group,omitempty"`
	ManagerGroup            string   `yaml:"space-manager-group,omitempty"`
	AuditorGroup            string   `yaml:"space-auditor-group,omitempty"`
	AllowSSH                bool     `yaml:"allow-ssh"`
	EnableSpaceQuota        bool     `yaml:"enable-space-quota"`
	MemoryLimit             string   `yaml:"memory-limit"`
	InstanceMemoryLimit     string   `yaml:"instance-memory-limit"`
	TotalRoutes             string   `yaml:"total-routes"`
	TotalServices           string   `yaml:"total-services"`
	PaidServicePlansAllowed bool     `yaml:"paid-service-plans-allowed"`
	EnableSecurityGroup     bool     `yaml:"enable-security-group"`
	SecurityGroupContents   string   `yaml:"security-group-contents,omitempty"`
	RemoveUsers             bool     `yaml:"enable-remove-users"`
	TotalPrivateDomains     string   `yaml:"total_private_domains"`
	TotalReservedRoutePorts string   `yaml:"total_reserved_route_ports"`
	TotalServiceKeys        string   `yaml:"total_service_keys"`
	AppInstanceLimit        string   `yaml:"app_instance_limit"`
	AppTaskLimit            string   `yaml:"app_task_limit"`
	IsoSegment              string   `yaml:"isolation_segment"`
	ASGs                    []string `yaml:"named-security-groups"`
}

func (s *SpaceConfig) InstanceMemoryLimitAsInt() (int, error) {
	return ToMegabytes(s.InstanceMemoryLimit)
}

func (s *SpaceConfig) MemoryLimitAsInt() (int, error) {
	return ToMegabytes(s.MemoryLimit)
}

// Contains determines whether a space is present in a list of spaces.
func (s *Spaces) Contains(spaceName string) bool {
	spaceNameToUpper := strings.ToUpper(spaceName)
	for _, v := range s.Spaces {
		if strings.ToUpper(v) == spaceNameToUpper {
			return true
		}
	}
	return false
}

func (s *Spaces) Replace(originalSpaceName, newSpaceName string) {
	lo.G.Debugf("Replacing %s with %s in space list", originalSpaceName, newSpaceName)
	var newList []string
	for _, spaceName := range s.Spaces {
		if !strings.EqualFold(spaceName, originalSpaceName) {
			newList = append(newList, spaceName)
		} else {
			lo.G.Debugf("Removing %s from space list", originalSpaceName)
		}
	}
	s.Spaces = append(newList, newSpaceName)
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
