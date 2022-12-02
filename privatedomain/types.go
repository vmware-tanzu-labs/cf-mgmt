package privatedomain

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// Manager -
type Manager interface {
	CreatePrivateDomains() error
	SharePrivateDomains() error
	ListOrgSharedPrivateDomains(orgGUID string) (map[string]*resource.Domain, error)
	ListOrgOwnedPrivateDomains(orgGUID string) (map[string]*resource.Domain, error)
}

type CFDomainClient interface {
	ListAll(ctx context.Context, opts *cfclient.DomainListOptions) ([]*resource.Domain, error)
	Create(ctx context.Context, r *resource.DomainCreate) (*resource.Domain, error)
	Share(ctx context.Context, domainGUID, organizationGUID string) (*resource.ToManyRelationships, error)
	ListForOrganizationAll(ctx context.Context, organizationGUID string, opts *cfclient.DomainListOptions) ([]*resource.Domain, error)
	Delete(ctx context.Context, guid string) (string, error)
	UnShare(ctx context.Context, domainGUID, organizationGUID string) error
}

type CFJobClient interface {
	PollComplete(ctx context.Context, jobGUID string, opts *cfclient.PollingOptions) error
}
