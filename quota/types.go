package quota

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type CFClient interface {
	ListOrgSpaceQuotas(orgGUID string) ([]cfclient.SpaceQuota, error)
	UpdateSpaceQuota(spaceQuotaGUID string, spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	AssignSpaceQuota(quotaGUID, spaceGUID string) error
	CreateSpaceQuota(spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	GetSpaceQuotaByName(name string) (cfclient.SpaceQuota, error)
	ListOrgQuotas() ([]cfclient.OrgQuota, error)
	CreateOrgQuota(orgQuote cfclient.OrgQuotaRequest) (*cfclient.OrgQuota, error)
	UpdateOrgQuota(orgQuotaGUID string, orgQuota cfclient.OrgQuotaRequest) (*cfclient.OrgQuota, error)
	GetOrgQuotaByName(name string) (cfclient.OrgQuota, error)
}
