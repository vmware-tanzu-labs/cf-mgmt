package config

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
}

func (s *Spaces) Contains(spaceName string) bool {
	for _, v := range s.Spaces {
		if v == spaceName {
			return true
		}
	}
	return false
}

func (i *SpaceConfig) GetDeveloperGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Developer.LDAPGroups {
		groupMap[group] = group
	}
	if i.Developer.LDAPGroup != "" {
		groupMap[i.Developer.LDAPGroup] = i.Developer.LDAPGroup
	}
	if i.DeveloperGroup != "" {
		groupMap[i.DeveloperGroup] = i.DeveloperGroup
	}
	return mapToKeys(groupMap)
}

func (i *SpaceConfig) GetManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Manager.LDAPGroups {
		groupMap[group] = group
	}
	if i.Manager.LDAPGroup != "" {
		groupMap[i.Manager.LDAPGroup] = i.Manager.LDAPGroup
	}
	if i.ManagerGroup != "" {
		groupMap[i.ManagerGroup] = i.ManagerGroup
	}
	return mapToKeys(groupMap)
}

func (i *SpaceConfig) GetAuditorGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Auditor.LDAPGroups {
		groupMap[group] = group
	}
	if i.Auditor.LDAPGroup != "" {
		groupMap[i.Auditor.LDAPGroup] = i.Auditor.LDAPGroup
	}
	if i.AuditorGroup != "" {
		groupMap[i.AuditorGroup] = i.AuditorGroup
	}
	return mapToKeys(groupMap)
}
