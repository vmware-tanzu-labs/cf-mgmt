package resource

import (
	"encoding/json"
	"time"
)

type ServiceRouteBinding struct {
	GUID          string        `json:"guid"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	LastOperation LastOperation `json:"last_operation"`

	// The URL for the route service
	RouteServiceURL string `json:"route_service_url"`

	// The route and service instance that the service route is bound to
	Relationships ServiceRouteBindingRelationships `json:"relationships"`

	Links    map[string]Link `json:"links"`
	Metadata *Metadata       `json:"metadata"`
}

type ServiceRouteBindingList struct {
	Pagination Pagination                   `json:"pagination"`
	Resources  []*ServiceRouteBinding       `json:"resources"`
	Included   *ServiceRouteBindingIncluded `json:"included"`
}

type ServiceRouteBindingCreate struct {
	Relationships ServiceRouteBindingRelationships `json:"relationships"`

	Metadata   *Metadata        `json:"metadata,omitempty"`
	Parameters *json.RawMessage `json:"parameters,omitempty"`
}

type ServiceRouteBindingUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type ServiceRouteBindingWithIncluded struct {
	ServiceRouteBinding
	Included *ServiceRouteBindingIncluded `json:"included"`
}

type ServiceRouteBindingIncluded struct {
	Routes           []*Route           `json:"routes"`
	ServiceInstances []*ServiceInstance `json:"service_instances"`
}

type ServiceRouteBindingRelationships struct {
	// The service instance that the route is bound to
	ServiceInstance ToOneRelationship `json:"service_instance"`

	// The route that the service instance is bound to
	Route ToOneRelationship `json:"route"`
}

// ServiceRouteBindingIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type ServiceRouteBindingIncludeType int

const (
	ServiceRouteBindingIncludeNone ServiceRouteBindingIncludeType = iota
	ServiceRouteBindingIncludeRoute
	ServiceRouteBindingIncludeServiceInstance
)

func (a ServiceRouteBindingIncludeType) String() string {
	switch a {
	case ServiceRouteBindingIncludeRoute:
		return "route"
	case ServiceRouteBindingIncludeServiceInstance:
		return "service_instance"
	}
	return ""
}

func NewServiceRouteBindingCreate(routeGUID, serviceInstanceGUID string) *ServiceRouteBindingCreate {
	return &ServiceRouteBindingCreate{
		Relationships: ServiceRouteBindingRelationships{
			ServiceInstance: ToOneRelationship{
				Data: &Relationship{
					GUID: serviceInstanceGUID,
				},
			},
			Route: ToOneRelationship{
				Data: &Relationship{
					GUID: routeGUID,
				},
			},
		},
	}
}
