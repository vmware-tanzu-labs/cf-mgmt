package resource

import (
	"encoding/json"
	"time"
)

// ServiceCredentialBinding implements the service credential binding object
// a credential binding can be a binding between apps and a service instance or a service key
type ServiceCredentialBinding struct {
	GUID          string        `json:"guid"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	LastOperation LastOperation `json:"last_operation"`
	Metadata      *Metadata     `json:"metadata"`

	Relationships ServiceCredentialBindingRelationships `json:"relationships"`
	Links         map[string]Link                       `json:"links"`
}

type ServiceCredentialBindingDetails struct {
	Credentials    map[string]any `json:"credentials"`
	SyslogDrainURL string         `json:"syslog_drain_url"`
	VolumeMounts   []string       `json:"volume_mounts"`
}

type ServiceCredentialBindingCreate struct {
	Type          string                                `json:"type"`          // Type of the service credential binding. Valid values are key and app
	Relationships ServiceCredentialBindingRelationships `json:"relationships"` // The service instance to be bound

	Name       *string          `json:"name,omitempty"`       // Name of the service credential binding. name is optional when the type is app
	Parameters *json.RawMessage `json:"parameters,omitempty"` // A JSON object that is passed to the service broker
	Metadata   *Metadata        `json:"metadata,omitempty"`
}

type ServiceCredentialBindingUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type ServiceCredentialBindingList struct {
	Pagination Pagination                        `json:"pagination"`
	Resources  []*ServiceCredentialBinding       `json:"resources"`
	Included   *ServiceCredentialBindingIncluded `json:"included"`
}

type ServiceCredentialBindingRelationships struct {
	App             *ToOneRelationship `json:"app,omitempty"`
	ServiceInstance *ToOneRelationship `json:"service_instance,omitempty"`
}

type ServiceCredentialBindingWithIncluded struct {
	ServiceCredentialBinding
	Included *ServiceCredentialBindingIncluded `json:"included"`
}

type ServiceCredentialBindingIncluded struct {
	Apps             []*App             `json:"apps"`
	ServiceInstances []*ServiceInstance `json:"service_instances"`
}

// ServiceCredentialBindingIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type ServiceCredentialBindingIncludeType int

const (
	ServiceCredentialBindingIncludeNone ServiceCredentialBindingIncludeType = iota
	ServiceCredentialBindingIncludeApp
	ServiceCredentialBindingIncludeServiceInstance
)

func (a ServiceCredentialBindingIncludeType) String() string {
	switch a {
	case ServiceCredentialBindingIncludeApp:
		return "app"
	case ServiceCredentialBindingIncludeServiceInstance:
		return "service_instance"
	}
	return ""
}

func NewServiceCredentialBindingCreateApp(serviceInstanceGUID, appGUID string) *ServiceCredentialBindingCreate {
	return &ServiceCredentialBindingCreate{
		Type: "app",
		Relationships: ServiceCredentialBindingRelationships{
			App: &ToOneRelationship{
				Data: &Relationship{
					GUID: appGUID,
				},
			},
			ServiceInstance: &ToOneRelationship{
				Data: &Relationship{
					GUID: serviceInstanceGUID,
				},
			},
		},
	}
}

func NewServiceCredentialBindingCreateKey(serviceInstanceGUID, bindingName string) *ServiceCredentialBindingCreate {
	return &ServiceCredentialBindingCreate{
		Type: "key",
		Name: &bindingName,
		Relationships: ServiceCredentialBindingRelationships{
			ServiceInstance: &ToOneRelationship{
				Data: &Relationship{
					GUID: serviceInstanceGUID,
				},
			},
		},
	}
}

func (s *ServiceCredentialBindingCreate) WithName(name string) *ServiceCredentialBindingCreate {
	s.Name = &name
	return s
}

func (s *ServiceCredentialBindingCreate) WithJSONParameters(jsonParams string) *ServiceCredentialBindingCreate {
	j := json.RawMessage(jsonParams)
	s.Parameters = &j
	return s
}
