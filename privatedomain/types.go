package privatedomain

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// Manager -
type Manager interface {
	CreatePrivateDomains() error
	SharePrivateDomains() error
	ListOrgSharedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error)
	ListOrgOwnedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error)
}

type CFClient interface {
	ListDomains() ([]cfclient.Domain, error)
	CreateDomain(name, orgGuid string) (*cfclient.Domain, error)
	ShareOrgPrivateDomain(orgGUID, privateDomainGUID string) (*cfclient.Domain, error)
	ListOrgPrivateDomains(orgGUID string) ([]cfclient.Domain, error)
	DeleteDomain(guid string) error
	UnshareOrgPrivateDomain(orgGUID, privateDomainGUID string) error
}
