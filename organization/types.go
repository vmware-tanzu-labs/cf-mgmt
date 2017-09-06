package organization

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	FindOrg(orgName string) (*cloudcontroller.Org, error)
	CreateOrgs(configFile string) error
	CreatePrivateDomains(configFile string) error
	DeleteOrgs(configFile string, peekDeletion bool) error
	UpdateOrgUsers(configDir, ldapBindPassword string) error
	CreateQuotas(configDir string) error
	GetOrgGUID(orgName string) (string, error)
}

// ORGS represents orgs constant
const ORGS = "organizations"
const ROLE_ORG_BILLING_MANAGERS = "billing_managers"
const ROLE_ORG_MANAGERS = "managers"
const ROLE_ORG_AUDITORS = "auditors"

//Resources -
type Resources struct {
	Resource []*Resource `json:"resources"`
}

//Resource -
type Resource struct {
	MetaData MetaData `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

//MetaData -
type MetaData struct {
	GUID string `json:"guid"`
}

//InputOrgs -
type InputOrgs struct {
	Orgs             []string `yaml:"orgs"`
	EnableDeleteOrgs bool     `yaml:"enable-delete-orgs"`
	ProtectedOrgs    []string `yaml:"protected_orgs"`
}

//Contains -
func (s *InputOrgs) Contains(orgName string) bool {
	for _, org := range s.Orgs {
		if org == orgName {
			return true
		}
	}
	return false
}

//InputUpdateOrgs -
type InputUpdateOrgs struct {
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
	IsoSegments             []string `yaml:"isolation_segments"`
	DefaultIsoSegment       string   `yaml:"default_isolation_segment"`
}

func (i *InputUpdateOrgs) GetBillingManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.BillingManager.LdapGroups {
		groupMap[group] = group
	}
	if i.BillingManager.LdapGroup != "" {
		groupMap[i.BillingManager.LdapGroup] = i.BillingManager.LdapGroup
	}
	if i.BillingManagerGroup != "" {
		groupMap[i.BillingManagerGroup] = i.BillingManagerGroup
	}
	return mapToKeys(groupMap)
}

func (i *InputUpdateOrgs) GetManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Manager.LdapGroups {
		groupMap[group] = group
	}
	if i.Manager.LdapGroup != "" {
		groupMap[i.Manager.LdapGroup] = i.Manager.LdapGroup
	}
	if i.ManagerGroup != "" {
		groupMap[i.ManagerGroup] = i.ManagerGroup
	}
	return mapToKeys(groupMap)
}

func (i *InputUpdateOrgs) GetAuditorGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Auditor.LdapGroups {
		groupMap[group] = group
	}
	if i.Auditor.LdapGroup != "" {
		groupMap[i.Auditor.LdapGroup] = i.Auditor.LdapGroup
	}
	if i.AuditorGroup != "" {
		groupMap[i.AuditorGroup] = i.AuditorGroup
	}
	return mapToKeys(groupMap)
}

func mapToKeys(aMap map[string]string) []string {
	keys := make([]string, 0, len(aMap))
	for k := range aMap {
		keys = append(keys, k)
	}
	return keys
}

type UserMgmt struct {
	LdapUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	LdapGroup  string   `yaml:"ldap_group"`
	LdapGroups []string `yaml:"ldap_groups"`
}

//Entity -
type Entity struct {
	Name string `json:"name"`
}

//Org -
type Org struct {
	AccessToken string `json:"access_token"`
}

//DefaultOrgManager -
type DefaultOrgManager struct {
	CloudController cloudcontroller.Manager
	UAACMgr         uaac.Manager
	UtilsMgr        utils.Manager
	LdapMgr         ldap.Manager
	UserMgr         UserMgr
}
