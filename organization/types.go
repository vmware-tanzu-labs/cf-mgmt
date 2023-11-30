package organization

import (
	"context"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Manager -
type Manager interface {
	CreateOrgs() error
	DeleteOrgs() error
	RenameOrg(originalOrgName, newOrgName string) error
	UpdateOrgsMetadata() error
}

type CFOrgClient interface {
	Create(ctx context.Context, r *resource.OrganizationCreate) (*resource.Organization, error)
	Update(ctx context.Context, guid string, r *resource.OrganizationUpdate) (*resource.Organization, error)
	Delete(ctx context.Context, guid string) (string, error)
}
