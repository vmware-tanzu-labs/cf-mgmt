package resource

import (
	"encoding/json"
	"time"
)

// ServiceOffering represent the services offered by service brokers
type ServiceOffering struct {
	GUID        string    `json:"guid"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`        // Name of the service offering
	Description string    `json:"description"` // Description of the service offering
	Available   bool      `json:"available"`   // Whether the service offering is available
	Tags        []string  `json:"tags"`        // Descriptive tags for the service offering

	// A list of permissions that the user would have to give the service, if they provision it;
	// the only permissions currently supported are syslog_drain, route_forwarding and volume_mount
	Requires []string `json:"requires"`

	// Whether service Instances of this service offering can be shared across organizations and spaces
	Shareable bool `json:"shareable"`

	// URL that points to a documentation page for the service offering,
	// if provided by the service broker as part of the metadata field
	DocumentationURL string `json:"documentation_url"`

	// This object contains information obtained from the service broker catalog
	BrokerCatalog ServiceOfferingBrokerCatalog `json:"broker_catalog"`

	Relationships ServiceBrokerRelationship `json:"relationships"`
	Metadata      *Metadata                 `json:"metadata"`
	Links         map[string]Link           `json:"links,omitempty"`
}

type ServiceOfferingList struct {
	Pagination Pagination         `json:"pagination"`
	Resources  []*ServiceOffering `json:"resources"`
}

type ServiceOfferingUpdate struct {
	Metadata *Metadata `json:"metadata,omitempty"`
}

type ServiceOfferingBrokerCatalog struct {
	// The identifier that the service broker provided for this service offering
	ID string `json:"id"`

	// https://github.com/openservicebrokerapi/servicebroker/blob/master/profile.md#service-metadata-fields
	Metadata *json.RawMessage `json:"metadata"`

	Features ServiceOfferingFeatures `json:"features"`
}

type ServiceOfferingFeatures struct {
	// Whether the service offering supports upgrade/downgrade for service plans by default; service plans can override this field
	PlanUpdateable bool `json:"plan_updateable"`

	// Specifies whether service Instances of the service can be bound to applications
	Bindable bool `json:"bindable"`

	// Specifies whether the Fetching a service instance endpoint is supported for all service plans
	InstancesRetrievable bool `json:"instances_retrievable"`

	// Specifies whether the Fetching a service binding endpoint is supported for all service plans
	BindingsRetrievable bool `json:"bindings_retrievable"`

	// Specifies whether service instance updates relating only to context are propagated to the service broker
	AllowContextUpdates bool `json:"allow_context_updates"`
}

type ServiceBrokerRelationship struct {
	ServiceBroker ToOneRelationship `json:"service_broker"`
}
