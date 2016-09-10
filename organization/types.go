package organization

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	FindOrg(orgName string) (org *cloudcontroller.Org, err error)
	CreateOrgs(configFile string) (err error)
	UpdateOrgUsers(configDir, ldapBindPassword string) (err error)
	CreateQuotas(configDir string) (err error)
	GetOrgGUID(orgName string) (orgGUID string, err error)
}

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
	Orgs []string `yaml:"orgs"`
}

//Contains -
func (s *InputOrgs) Contains(orgName string) bool {
	set := make(map[string]bool)
	for _, v := range s.Orgs {
		set[v] = true
	}
	_, ok := set[orgName]
	return ok
}

//InputUpdateOrgs -
type InputUpdateOrgs struct {
	Org                     string `yaml:"org"`
	BillingManagerGroup     string `yaml:"org-billingmanager-group"`
	ManagerGroup            string `yaml:"org-manager-group"`
	AuditorGroup            string `yaml:"org-auditor-group"`
	EnableOrgQuota          bool   `yaml:"enable-org-quota"`
	MemoryLimit             int    `yaml:"memory-limit"`
	InstanceMemoryLimit     int    `yaml:"instance-memory-limit"`
	TotalRoutes             int    `yaml:"total-routes"`
	TotalServices           int    `yaml:"total-services"`
	PaidServicePlansAllowed bool   `yaml:"paid-service-plans-allowed"`
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
}
