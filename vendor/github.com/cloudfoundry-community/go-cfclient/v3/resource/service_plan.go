package resource

import (
	"encoding/json"
	"time"
)

type ServicePlan struct {
	GUID        string    `json:"guid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// the visibility of the plan; can be public, admin, organization
	VisibilityType string `json:"visibility_type"`

	// Whether the service plan is available
	Available bool `json:"available"`

	// Whether the service plan is free of charge
	Free bool `json:"free"`

	// The cost of the service plan as obtained from the service broker catalog
	Costs []ServicePlanCosts `json:"costs"`

	// Information about the version of this service plan
	MaintenanceInfo ServicePlanMaintenanceInfo `json:"maintenance_info"`

	// This object contains information obtained from the service broker catalog
	BrokerCatalog ServicePlanBrokerCatalog `json:"broker_catalog"`

	// Schema definitions for service instances and service bindings for the service plan
	Schemas ServicePlanSchemas `json:"schemas"`

	// The service offering that this service plan relates to
	Relationships ServicePlanRelationship `json:"relationships"`

	Metadata *Metadata       `json:"metadata"`
	Links    map[string]Link `json:"links"`
}

type ServicePlanList struct {
	Pagination Pagination           `json:"pagination"`
	Resources  []*ServicePlan       `json:"resources"`
	Included   *ServicePlanIncluded `json:"included"`
}

type ServicePlanWithIncluded struct {
	ServicePlan
	Included *ServicePlanIncluded `json:"included"`
}

type ServicePlanIncluded struct {
	Organizations    []*Organization    `json:"organizations"`
	Spaces           []*Space           `json:"spaces"`
	ServiceOfferings []*ServiceOffering `json:"service_offerings"`
}

type ServicePlanUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type ServicePlanCosts struct {
	// Currency code for the pricing amount, e.g. USD, GBP
	Currency string `json:"currency"`

	// Pricing amount
	Amount float64 `json:"amount"`

	// Display name for type of cost, e.g. Monthly, Hourly, Request, GB
	Unit string `json:"unit"`
}

type ServicePlanMaintenanceInfo struct {
	// The current semantic version of the service plan
	// comparing this version with that of a service instance can be used to determine
	// whether the service instance is up-to-date with this service plan
	Version string `json:"version"`

	// A textual explanation associated with this version
	Description string `json:"description"`
}

type ServicePlanBrokerCatalog struct {
	// The identifier that the service broker provided for this service plan
	ID string `json:"id"`

	// Additional information provided by the service broker as specified by OSBAPI
	// https://github.com/openservicebrokerapi/servicebroker/blob/master/profile.md#plan-metadata-fields
	Metadata *json.RawMessage `json:"metadata"`

	// The maximum number of seconds that Cloud Foundry will wait for an asynchronous service broker operation
	MaximumPollingDuration *int `json:"maximum_polling_duration"`

	// Features the service plan supports or not
	Features ServicePlanFeatures `json:"features"`
}

type ServicePlanFeatures struct {
	// Whether the service plan supports upgrade/downgrade for service plans
	// when the catalog does not specify a value, it is inherited from the service offering
	PlanUpdateable bool `json:"plan_updateable"`

	// Specifies whether service instances of the service can be bound to applications
	Bindable bool `json:"bindable"`
}

type ServicePlanSchemas struct {
	ServiceInstance ServicePlanServiceInstance `json:"service_instance"`
	ServiceBinding  ServicePlanServiceBinding  `json:"service_binding"`
}

type ServicePlanServiceInstance struct {
	// Schema definition for service instance creation
	Create ServicePlanSchemaCreateOrUpdate `json:"create"`

	// Schema definition for service instance update
	Update ServicePlanSchemaCreateOrUpdate `json:"update"`
}

type ServicePlanServiceBinding struct {
	// Schema definition for service Binding creation
	Create ServicePlanSchemaCreateOrUpdate `json:"create"`
}

type ServicePlanSchemaCreateOrUpdate struct {
	// The schema definition for the input parameters
	// each input parameter is expressed as a property within a JSON object
	Parameters *json.RawMessage `json:"parameters"`
}

type ServicePlanRelationship struct {
	ServiceOffering ToOneRelationship `json:"service_offering"`
}

// ServicePlanIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type ServicePlanIncludeType int

const (
	ServicePlanIncludeNone ServicePlanIncludeType = iota
	ServicePlanIncludeSpaceOrganization
	ServicePlanIncludeServiceOffering
)

func (a ServicePlanIncludeType) String() string {
	switch a {
	case ServicePlanIncludeSpaceOrganization:
		return "space.organization"
	case ServicePlanIncludeServiceOffering:
		return "service_offering"
	}
	return ""
}
