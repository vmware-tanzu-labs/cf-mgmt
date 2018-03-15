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
	GetOrgGUID(orgName string) (string, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	GetOrgByGUID(orgGUID string) (cfclient.Org, error)
}

type CFClient interface {
	ListOrgs() ([]cfclient.Org, error)
	DeleteOrg(guid string, recursive, async bool) error
	CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error)
	GetOrgByGuid(guid string) (cfclient.Org, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	ListDomains() ([]cfclient.Domain, error)
	CreateDomain(name, orgGuid string) (*cfclient.Domain, error)
	ShareOrgPrivateDomain(orgGUID, privateDomainGUID string) (*cfclient.Domain, error)
	ListOrgPrivateDomains(orgGUID string) ([]cfclient.Domain, error)
	DeleteDomain(guid string) error
	UnshareOrgPrivateDomain(orgGUID, privateDomainGUID string) error
}
