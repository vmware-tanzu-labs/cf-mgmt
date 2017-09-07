package config

// OrgConfig describes configuration for an org.
type OrgConfig struct {
	Org                     string   `yaml:"org"`
	BillingManagerGroup     string   `yaml:"org-billingmanager-group,omitempty"`
	ManagerGroup            string   `yaml:"org-manager-group,omitempty"`
	AuditorGroup            string   `yaml:"org-auditor-group,omitempty"`
	BillingManager          UserMgmt `yaml:"org-billingmanager"`
	Manager                 UserMgmt `yaml:"org-manager"`
	Auditor                 UserMgmt `yaml:"org-auditor"`
	PrivateDomains          []string `yaml:"private-domains"`
	RemovePrivateDomains    bool     `yaml:"enable-remove-private-domains"`
	EnableOrgQuota          bool     `yaml:"enable-org-quota"`
	MemoryLimit             int      `yaml:"memory-limit"`
	InstanceMemoryLimit     int      `yaml:"instance-memory-limit"`
	TotalRoutes             int      `yaml:"total-routes"`
	TotalServices           int      `yaml:"total-services"`
	PaidServicePlansAllowed bool     `yaml:"paid-service-plans-allowed"`
	RemoveUsers             bool     `yaml:"enable-remove-users"`
	TotalPrivateDomains     int      `yaml:"total_private_domains"`
	TotalReservedRoutePorts int      `yaml:"total_reserved_route_ports"`
	TotalServiceKeys        int      `yaml:"total_service_keys"`
	AppInstanceLimit        int      `yaml:"app_instance_limit"`
	DefaultIsoSegment       string   `yaml:"default_isolation_segment"`
}

// Orgs contains cf-mgmt configuration for all orgs.
type Orgs struct {
	Orgs             []string `yaml:"orgs"`
	EnableDeleteOrgs bool     `yaml:"enable-delete-orgs"`
	ProtectedOrgs    []string `yaml:"protected_orgs"`
}

// Contains determines whether an org is present in a list of orgs.
func (o *Orgs) Contains(orgName string) bool {
	for _, org := range o.Orgs {
		if org == orgName {
			return true
		}
	}
	return false
}

func (o *OrgConfig) GetBillingManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range o.BillingManager.LDAPGroups {
		groupMap[group] = group
	}
	if o.BillingManager.LDAPGroup != "" {
		groupMap[o.BillingManager.LDAPGroup] = o.BillingManager.LDAPGroup
	}
	if o.BillingManagerGroup != "" {
		groupMap[o.BillingManagerGroup] = o.BillingManagerGroup
	}
	return mapToKeys(groupMap)
}

func (o *OrgConfig) GetManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range o.Manager.LDAPGroups {
		groupMap[group] = group
	}
	if o.Manager.LDAPGroup != "" {
		groupMap[o.Manager.LDAPGroup] = o.Manager.LDAPGroup
	}
	if o.ManagerGroup != "" {
		groupMap[o.ManagerGroup] = o.ManagerGroup
	}
	return mapToKeys(groupMap)
}

func (o *OrgConfig) GetAuditorGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range o.Auditor.LDAPGroups {
		groupMap[group] = group
	}
	if o.Auditor.LDAPGroup != "" {
		groupMap[o.Auditor.LDAPGroup] = o.Auditor.LDAPGroup
	}
	if o.AuditorGroup != "" {
		groupMap[o.AuditorGroup] = o.AuditorGroup
	}
	return mapToKeys(groupMap)
}

func mapToKeys(aMap map[string]string) []string {
	var keys []string
	for k := range aMap {
		keys = append(keys, k)
	}
	return keys
}
