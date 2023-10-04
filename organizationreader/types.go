package organizationreader

import (
	"context"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	v3cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Reader -
type Reader interface {
	ListOrgs() ([]*resource.Organization, error)
	FindOrg(orgName string) (*resource.Organization, error)
	FindOrgByGUID(orgGUID string) (*resource.Organization, error)
	GetOrgGUID(orgName string) (string, error)
	ClearOrgList()
	AddOrgToList(org *resource.Organization)
	GetDefaultIsolationSegment(org *resource.Organization) (string, error)
}

type CFClient interface {
	ListOrgs() ([]cfclient.Org, error)
	DeleteOrg(guid string, recursive, async bool) error
	GetOrgByGuid(guid string) (cfclient.Org, error)
}

type CFOrgClient interface {
	ListAll(ctx context.Context, opts *v3cfclient.OrganizationListOptions) ([]*resource.Organization, error)
	GetDefaultIsolationSegment(ctx context.Context, guid string) (string, error)
}
