package serviceaccess

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_cf_client.go types.go CFClient
type CFClient interface {
	ListServicePlansByQuery(query url.Values) ([]cfclient.ServicePlan, error)
	MakeServicePlanPrivate(servicePlanGUID string) error
	MakeServicePlanPublic(servicePlanGUID string) error
	ListServiceBrokers() ([]cfclient.ServiceBroker, error)
	ListServicesByQuery(query url.Values) ([]cfclient.Service, error)
	ListServicePlanVisibilitiesByQuery(query url.Values) ([]cfclient.ServicePlanVisibility, error)
	CreateServicePlanVisibility(servicePlanGuid string, organizationGuid string) (cfclient.ServicePlanVisibility, error)
	DeleteServicePlanVisibilityByPlanAndOrg(servicePlanGuid string, organizationGuid string, async bool) error
	ListServices() ([]cfclient.Service, error)
}
