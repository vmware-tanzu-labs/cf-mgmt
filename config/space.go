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
	Org                         string    `yaml:"org"`
	Space                       string    `yaml:"space"`
	OriginalSpace               string    `yaml:"original-space,omitempty"`
	Developer                   UserMgmt  `yaml:"space-developer"`
	Manager                     UserMgmt  `yaml:"space-manager"`
	Auditor                     UserMgmt  `yaml:"space-auditor"`
	Supporter                   UserMgmt  `yaml:"space-supporter"`
	DeveloperGroup              string    `yaml:"space-developer-group,omitempty"`
	ManagerGroup                string    `yaml:"space-manager-group,omitempty"`
	AuditorGroup                string    `yaml:"space-auditor-group,omitempty"`
	SupporterGroup              string    `yaml:"space-supporter-group,omitempty"`
	AllowSSH                    bool      `yaml:"allow-ssh"`
	AllowSSHUntil               string    `yaml:"allow-ssh-until,omitempty"`
	EnableSpaceQuota            bool      `yaml:"enable-space-quota"`
	EnableSecurityGroup         bool      `yaml:"enable-security-group"`
	EnableUnassignSecurityGroup bool      `yaml:"enable-unassign-security-group"`
	SecurityGroupContents       string    `yaml:"security-group-contents,omitempty"`
	RemoveUsers                 bool      `yaml:"enable-remove-users"`
	IsoSegment                  string    `yaml:"isolation_segment"`
	ASGs                        []string  `yaml:"named-security-groups"`
	MemoryLimit                 string    `yaml:"memory-limit,omitempty"`
	InstanceMemoryLimit         string    `yaml:"instance-memory-limit,omitempty"`
	TotalRoutes                 string    `yaml:"total-routes,omitempty"`
	TotalServices               string    `yaml:"total-services,omitempty"`
	PaidServicePlansAllowed     bool      `yaml:"paid-service-plans-allowed"`
	TotalReservedRoutePorts     string    `yaml:"total_reserved_route_ports,omitempty"`
	TotalServiceKeys            string    `yaml:"total_service_keys,omitempty"`
	AppInstanceLimit            string    `yaml:"app_instance_limit,omitempty"`
	AppTaskLimit                string    `yaml:"app_task_limit,omitempty"`
	LogRateLimitBytesPerSecond  string    `yaml:"log_rate_limit_bytes_per_second,omitempty"`
	NamedQuota                  string    `yaml:"named_quota"`
	Metadata                    *Metadata `yaml:"metadata"`
}

func (s *SpaceConfig) GetSecurityGroupContents() string {
	if len(s.SecurityGroupContents) > 0 {
		return s.SecurityGroupContents
	} else {
		return "[]"
	}
}
func (s *SpaceConfig) GetQuota() SpaceQuota {
	return SpaceQuota{
		Name:                       s.Space,
		Org:                        s.Org,
		MemoryLimit:                s.MemoryLimit,
		InstanceMemoryLimit:        s.InstanceMemoryLimit,
		TotalRoutes:                s.TotalRoutes,
		TotalServices:              s.TotalServices,
		PaidServicePlansAllowed:    s.PaidServicePlansAllowed,
		TotalReservedRoutePorts:    s.TotalReservedRoutePorts,
		TotalServiceKeys:           s.TotalServiceKeys,
		AppInstanceLimit:           s.AppInstanceLimit,
		AppTaskLimit:               s.AppTaskLimit,
		LogRateLimitBytesPerSecond: s.LogRateLimitBytesPerSecond,
	}
}

type SpaceQuota struct {
	Name                       string `yaml:"-"`
	Org                        string `yaml:"-"`
	MemoryLimit                string `yaml:"memory-limit"`
	InstanceMemoryLimit        string `yaml:"instance-memory-limit"`
	TotalRoutes                string `yaml:"total-routes"`
	TotalServices              string `yaml:"total-services"`
	PaidServicePlansAllowed    bool   `yaml:"paid-service-plans-allowed"`
	TotalReservedRoutePorts    string `yaml:"total_reserved_route_ports"`
	TotalServiceKeys           string `yaml:"total_service_keys"`
	AppInstanceLimit           string `yaml:"app_instance_limit"`
	AppTaskLimit               string `yaml:"app_task_limit"`
	LogRateLimitBytesPerSecond string `yaml:"log_rate_limit_bytes_per_second"`
}

func (s *SpaceQuota) IsUnlimitedMemory() bool {
	return strings.EqualFold(s.MemoryLimit, "-1") || strings.EqualFold(s.MemoryLimit, "unlimited")
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

func (i *SpaceConfig) GetSupporterGroups() []string {
	return i.Supporter.groups(i.SupporterGroup)
}
