package serviceaccess

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//go:generate counterfeiter -o fakes/fake_cf_client.go types.go CFClient
type CFClient interface {
	ListServicePlansByQuery(query url.Values) ([]cfclient.ServicePlan, error)
	MakeServicePlanPrivate(servicePlanGUID string) error
	ListServices() ([]cfclient.Service, error)
	ListServicePlanVisibilitiesByQuery(query url.Values) ([]cfclient.ServicePlanVisibility, error)
	CreateServicePlanVisibility(servicePlanGuid string, organizationGuid string) (cfclient.ServicePlanVisibility, error)
	DeleteServicePlanVisibilityByPlanAndOrg(servicePlanGuid string, organizationGuid string, async bool) error
}
