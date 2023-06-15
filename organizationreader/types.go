package organizationreader

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// Reader -
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
	GetOrgByGuid(guid string) (cfclient.Org, error)
}
