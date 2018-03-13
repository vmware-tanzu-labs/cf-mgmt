package spacequota

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Manager interface {
	CreateQuotas(configDir string) error
	SpaceQuotaByName(name string) (cfclient.SpaceQuota, error)
}

type CFClient interface {
	ListOrgSpaceQuotas(orgGUID string) ([]cfclient.SpaceQuota, error)
	UpdateSpaceQuota(spaceQuotaGUID string, spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	AssignSpaceQuota(quotaGUID, spaceGUID string) error
	CreateSpaceQuota(spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	GetSpaceQuotaByName(name string) (cfclient.SpaceQuota, error)
}
