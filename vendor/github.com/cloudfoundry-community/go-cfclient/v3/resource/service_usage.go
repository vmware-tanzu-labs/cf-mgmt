package resource

import "time"

type ServiceUsage struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Current state of the service that this event pertains to, if applicable
	State string `json:"state"`

	// space that this event pertains to, if applicable
	Space ServiceUsageGUIDName `json:"space"`

	// organization that this event pertains to, if applicable
	Organization Relationship `json:"organization"`

	// service instance that this event pertains to, if applicable
	ServiceInstance ServiceUsageGUIDNameType `json:"service_instance"`

	// service plan that this event pertains to, if applicable
	ServicePlan ServiceUsageGUIDName `json:"service_plan"`

	// service offering that this event pertains to, if applicable
	ServiceOffering ServiceUsageGUIDName `json:"service_offering"`

	// service broker that this event pertains to, if applicable
	ServiceBroker ServiceUsageGUIDName `json:"service_broker"`

	Links map[string]Link `json:"links"`
}

type ServiceUsageList struct {
	Pagination Pagination      `json:"pagination"`
	Resources  []*ServiceUsage `json:"resources"`
}

type ServiceUsageGUIDName struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
}

type ServiceUsageGUIDNameType struct {
	GUID string `json:"guid"`
	Name string `json:"name"`
	Type string `json:"type"`
}
