package organization

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Manager -
type Manager interface {
	CreateOrgs() error
	DeleteOrgs() error
	UpdateOrg(orgGUID string, orgRequest *resource.OrganizationUpdate) (*resource.Organization, error)
	RenameOrg(originalOrgName, newOrgName string) error
	UpdateOrgsMetadata() error
}

type CFOrganizationClient interface {
	Delete(ctx context.Context, guid string) (string, error)
	Create(ctx context.Context, r *resource.OrganizationCreate) (*resource.Organization, error)
	Update(ctx context.Context, guid string, r *resource.OrganizationUpdate) (*resource.Organization, error)
}

type CFDomainClient interface {
	UnShare(ctx context.Context, domainGUID, organizationGUID string) error
}

type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *cfclient.PollingOptions) error
}
