package organization

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//Manager -
type Manager interface {
	CreateOrgs() error
	DeleteOrgs() error
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	RenameOrg(originalOrgName, newOrgName string) error
	UpdateOrgsMetadata() error
}

//Manager -
type Reader interface {
	ListOrgs() ([]cfclient.Org, error)
	FindOrg(orgName string) (cfclient.Org, error)
	FindOrgByGUID(orgGUID string) (cfclient.Org, error)
	GetOrgGUID(orgName string) (string, error)
	GetOrgByGUID(orgGUID string) (cfclient.Org, error)
	ClearOrgList()
	AddOrgToList(org cfclient.Org)
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
	SupportsMetadataAPI() (bool, error)
	UpdateOrgMetadata(orgGUID string, metadata cfclient.Metadata) error
	OrgMetadata(orgGUID string) (*cfclient.Metadata, error)
	RemoveOrgMetadata(orgGUID string) error
}
