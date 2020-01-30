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

type CFClient interface {
	DeleteOrg(guid string, recursive, async bool) error
	CreateOrg(req cfclient.OrgRequest) (cfclient.Org, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	UnshareOrgPrivateDomain(orgGUID, privateDomainGUID string) error
	SupportsMetadataAPI() (bool, error)
	UpdateOrgMetadata(orgGUID string, metadata cfclient.Metadata) error
	OrgMetadata(orgGUID string) (*cfclient.Metadata, error)
	RemoveOrgMetadata(orgGUID string) error
}
