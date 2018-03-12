package organization

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//Manager -
type Manager interface {
	ListOrgs() ([]cfclient.Org, error)
	ListOrgSharedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error)
	ListOrgOwnedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error)
	FindOrg(orgName string) (cfclient.Org, error)
	CreateOrgs() error
	CreatePrivateDomains() error
	SharePrivateDomains() error
	DeleteOrgs() error
	UpdateOrgUsers(configDir, ldapBindPassword string) error
	CreateQuotas() error
	GetOrgGUID(orgName string) (string, error)
	ListOrgAuditors(orgGUID string) (map[string]string, error)
	ListOrgBillingManager(orgGUID string) (map[string]string, error)
	ListOrgManagers(orgGUID string) (map[string]string, error)
	OrgQuotaByName(name string) (cfclient.OrgQuota, error)
}

// UserMgr - interface type encapsulating Update space users behavior
type UserMgr interface {
	RemoveOrgAuditorByUsername(orgGUID, userName string) error
	RemoveOrgBillingManagerByUsername(orgGUID, userName string) error
	RemoveOrgManagerByUsername(orgGUID, userName string) error
	ListOrgAuditors(orgGUID string) (map[string]string, error)
	ListOrgBillingManager(orgGUID string) (map[string]string, error)
	ListOrgManagers(orgGUID string) (map[string]string, error)
	AssociateOrgAuditorByUsername(orgGUID, userName string) error
	AssociateOrgBillingManagerByUsername(orgGUID, userName string) error
	AssociateOrgManagerByUsername(orgGUID, userName string) error
}

type CFClient interface {
	RemoveOrgUserByUsername(orgGUID, name string) error
	RemoveOrgAuditorByUsername(orgGUID, name string) error
	RemoveOrgBillingManagerByUsername(orgGUID, name string) error
	RemoveOrgManagerByUsername(orgGUID, name string) error
	ListOrgAuditors(orgGUID string) ([]cfclient.User, error)
	ListOrgManagers(orgGUID string) ([]cfclient.User, error)
	ListOrgBillingManagers(orgGUID string) ([]cfclient.User, error)
	AssociateOrgAuditorByUsername(orgGUID, name string) (cfclient.Org, error)
	AssociateOrgManagerByUsername(orgGUID, name string) (cfclient.Org, error)
	AssociateOrgBillingManagerByUsername(orgGUID, name string) (cfclient.Org, error)
	ListOrgs() ([]cfclient.Org, error)
	DeleteOrg(guid string, recursive, async bool) error
	CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error)
	GetOrgByGuid(guid string) (cfclient.Org, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	AssociateOrgUserByUsername(orgGUID, userName string) (cfclient.Org, error)
	ListDomains() ([]cfclient.Domain, error)
	CreateDomain(name, orgGuid string) (*cfclient.Domain, error)
	ShareOrgPrivateDomain(orgGUID, privateDomainGUID string) (*cfclient.Domain, error)
	ListOrgPrivateDomains(orgGUID string) ([]cfclient.Domain, error)
	DeleteDomain(guid string) error
	UnshareOrgPrivateDomain(orgGUID, privateDomainGUID string) error
	ListOrgQuotas() ([]cfclient.OrgQuota, error)
	CreateOrgQuota(orgQuote cfclient.OrgQuotaRequest) (*cfclient.OrgQuota, error)
	UpdateOrgQuota(orgQuotaGUID string, orgQuota cfclient.OrgQuotaRequest) (*cfclient.OrgQuota, error)
	GetOrgQuotaByName(name string) (cfclient.OrgQuota, error)
}

// UpdateUsersInput -
type UpdateUsersInput struct {
	OrgName                                     string
	OrgGUID                                     string
	LdapUsers, Users, LdapGroupNames, SamlUsers []string
	RemoveUsers                                 bool
	ListUsers                                   func(orgGUID string) (map[string]string, error)
	AddUser                                     func(orgGUID, userName string) error
	RemoveUser                                  func(orgGUID, userName string) error
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

//Entity -
type Entity struct {
	Name string `json:"name"`
}

//Org -
type Org struct {
	AccessToken string `json:"access_token"`
}
