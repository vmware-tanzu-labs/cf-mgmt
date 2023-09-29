package config

import (
	"strings"

	"github.com/xchapter7x/lo"
)

// OrgConfig describes configuration for an org.
type OrgConfig struct {
	Org                        string              `yaml:"org"`
	OriginalOrg                string              `yaml:"original-org,omitempty"`
	BillingManagerGroup        string              `yaml:"org-billingmanager-group,omitempty"`
	ManagerGroup               string              `yaml:"org-manager-group,omitempty"`
	AuditorGroup               string              `yaml:"org-auditor-group,omitempty"`
	BillingManager             UserMgmt            `yaml:"org-billingmanager"`
	Manager                    UserMgmt            `yaml:"org-manager"`
	Auditor                    UserMgmt            `yaml:"org-auditor"`
	PrivateDomains             []string            `yaml:"private-domains"`
	RemovePrivateDomains       bool                `yaml:"enable-remove-private-domains"`
	SharedPrivateDomains       []string            `yaml:"shared-private-domains"`
	RemoveSharedPrivateDomains bool                `yaml:"enable-remove-shared-private-domains"`
	EnableOrgQuota             bool                `yaml:"enable-org-quota"`
	MemoryLimit                string              `yaml:"memory-limit,omitempty"`
	InstanceMemoryLimit        string              `yaml:"instance-memory-limit,omitempty"`
	TotalRoutes                string              `yaml:"total-routes,omitempty"`
	TotalServices              string              `yaml:"total-services,omitempty"`
	PaidServicePlansAllowed    bool                `yaml:"paid-service-plans-allowed"`
	RemoveUsers                bool                `yaml:"enable-remove-users"`
	TotalPrivateDomains        string              `yaml:"total_private_domains,omitempty"`
	TotalReservedRoutePorts    string              `yaml:"total_reserved_route_ports,omitempty"`
	TotalServiceKeys           string              `yaml:"total_service_keys,omitempty"`
	AppInstanceLimit           string              `yaml:"app_instance_limit,omitempty"`
	AppTaskLimit               string              `yaml:"app_task_limit,omitempty"`
	LogRateLimitBytesPerSecond string              `yaml:"log_rate_limit_bytes_per_second,omitempty"`
	DefaultIsoSegment          string              `yaml:"default_isolation_segment"`
	ServiceAccess              map[string][]string `yaml:"service-access,omitempty"`
	NamedQuota                 string              `yaml:"named_quota"`
	Metadata                   *Metadata           `yaml:"metadata"`
}

func (o *OrgConfig) GetQuota() OrgQuota {
	return OrgQuota{
		Name:                       o.Org,
		TotalPrivateDomains:        o.TotalPrivateDomains,
		TotalReservedRoutePorts:    o.TotalReservedRoutePorts,
		TotalServiceKeys:           o.TotalServiceKeys,
		AppInstanceLimit:           o.AppInstanceLimit,
		AppTaskLimit:               o.AppTaskLimit,
		MemoryLimit:                o.MemoryLimit,
		InstanceMemoryLimit:        o.InstanceMemoryLimit,
		TotalRoutes:                o.TotalRoutes,
		TotalServices:              o.TotalServices,
		PaidServicePlansAllowed:    o.PaidServicePlansAllowed,
		LogRateLimitBytesPerSecond: o.LogRateLimitBytesPerSecond,
	}
}

type OrgQuota struct {
	Name                       string `yaml:"-"`
	TotalPrivateDomains        string `yaml:"total_private_domains"`
	TotalReservedRoutePorts    string `yaml:"total_reserved_route_ports"`
	TotalServiceKeys           string `yaml:"total_service_keys"`
	AppInstanceLimit           string `yaml:"app_instance_limit"`
	AppTaskLimit               string `yaml:"app_task_limit"`
	MemoryLimit                string `yaml:"memory-limit"`
	InstanceMemoryLimit        string `yaml:"instance-memory-limit"`
	TotalRoutes                string `yaml:"total-routes"`
	TotalServices              string `yaml:"total-services"`
	PaidServicePlansAllowed    bool   `yaml:"paid-service-plans-allowed"`
	LogRateLimitBytesPerSecond string `yaml:"log_rate_limit_bytes_per_second"`
}

// Orgs contains cf-mgmt configuration for all orgs.
type Orgs struct {
	Orgs             []string `yaml:"orgs"`
	EnableDeleteOrgs bool     `yaml:"enable-delete-orgs"`
	ProtectedOrgs    []string `yaml:"protected_orgs"`
}

func (o *Orgs) ProtectedOrgList() []string {
	var allOrgNames []string
	uniqueNames := make(map[string]string)
	allOrgNames = append(o.ProtectedOrgs, DefaultProtectedOrgs...)
	for _, orgName := range allOrgNames {
		uniqueNames[orgName] = orgName
	}
	var returnList []string
	for _, name := range uniqueNames {
		returnList = append(returnList, name)
	}

	return returnList
}

func (o *Orgs) Replace(originalOrgName, newOrgName string) {
	lo.G.Debugf("Replacing %s with %s in org list", originalOrgName, newOrgName)
	var newList []string
	for _, orgName := range o.Orgs {
		if !strings.EqualFold(orgName, originalOrgName) {
			newList = append(newList, orgName)
		} else {
			lo.G.Debugf("Removing %s from org list", originalOrgName)
		}
	}
	o.Orgs = append(newList, newOrgName)
}

// Contains determines whether an org is present in a list of orgs.
func (o *Orgs) Contains(orgName string) bool {
	orgNameUpper := strings.ToUpper(orgName)
	for _, org := range o.Orgs {
		if strings.ToUpper(org) == orgNameUpper {
			return true
		}
	}
	return false
}

func (o *OrgConfig) GetBillingManagerGroups() []string {
	return o.BillingManager.groups(o.BillingManagerGroup)
}

func (o *OrgConfig) GetManagerGroups() []string {
	return o.Manager.groups(o.ManagerGroup)
}

func (o *OrgConfig) GetAuditorGroups() []string {
	return o.Auditor.groups(o.AuditorGroup)
}
