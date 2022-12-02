package organizationreader

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Reader -
type Reader interface {
	ListOrgs() ([]*resource.Organization, error)
	FindOrg(orgName string) (*resource.Organization, error)
	FindOrgByGUID(orgGUID string) (*resource.Organization, error)
	GetOrgGUID(orgName string) (string, error)
	GetOrgByGUID(orgGUID string) (*resource.Organization, error)
	ClearOrgList()
	AddOrgToList(org *resource.Organization)
}

type CFOrganizationClient interface {
	ListAll(ctx context.Context, opts *cfclient.OrganizationListOptions) ([]*resource.Organization, error)
	Delete(ctx context.Context, guid string) (string, error)
	Get(ctx context.Context, guid string) (*resource.Organization, error)
}
