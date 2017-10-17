package organization

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	FindOrg(orgName string) (*cloudcontroller.Org, error)
	CreateOrgs() error
	CreatePrivateDomains() error
	DeleteOrgs(peekDeletion bool) error
	UpdateOrgUsers(configDir, ldapBindPassword string) error
	CreateQuotas() error
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
	Cfg             config.Reader
	CloudController cloudcontroller.Manager
	UAACMgr         uaac.Manager
	UtilsMgr        utils.Manager
	LdapMgr         ldap.Manager
	UserMgr         UserMgr
}
