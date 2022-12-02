package serviceaccess

import (
	"context"
	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

//go:generate counterfeiter -o fakes/fake_svc_plan_client.go types.go CFServicePlanClient
type CFServicePlanClient interface {
	ListAll(ctx context.Context, opts *cfclient.ServicePlanListOptions) ([]*resource.ServicePlan, error)
}

//go:generate counterfeiter -o fakes/fake_svc_plan_visibility_client.go types.go CFServicePlanVisibilityClient
type CFServicePlanVisibilityClient interface {
	Apply(ctx context.Context, servicePlanGUID string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error)
	Delete(ctx context.Context, servicePlanGUID, organizationGUID string) error
	Get(ctx context.Context, servicePlanGUID string) (*resource.ServicePlanVisibility, error)
	Update(ctx context.Context, servicePlanGUID string, r *resource.ServicePlanVisibility) (*resource.ServicePlanVisibility, error)
}

//go:generate counterfeiter -o fakes/fake_svc_offering_client.go types.go CFServiceOfferingClient
type CFServiceOfferingClient interface {
	ListAll(ctx context.Context, opts *cfclient.ServiceOfferingListOptions) ([]*resource.ServiceOffering, error)
}

//go:generate counterfeiter -o fakes/fake_svc_broker_client.go types.go CFServiceBrokerClient
type CFServiceBrokerClient interface {
	ListAll(ctx context.Context, opts *cfclient.ServiceBrokerListOptions) ([]*resource.ServiceBroker, error)
}
