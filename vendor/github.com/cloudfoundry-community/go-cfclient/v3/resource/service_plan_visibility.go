package resource

import "fmt"

type ServicePlanVisibility struct {
	// Denotes the visibility of the plan; can be public, admin, organization, space
	Type string `json:"type"`

	// List of organizations whose members can access the plan; present if type is organization
	Organizations []ServicePlanVisibilityRelation `json:"organizations,omitempty"`

	// space whose members can access the plan; present if type is space
	Space *ServicePlanVisibilityRelation `json:"space,omitempty"`
}

type ServicePlanVisibilityRelation struct {
	// org or space GUID
	GUID string `json:"guid"`

	// org or space name, only used in responses
	Name *string `json:"name,omitempty"`
}

type ServicePlanVisibilityType int

const (
	ServicePlanVisibilityNone ServicePlanVisibilityType = iota
	ServicePlanVisibilityPublic
	ServicePlanVisibilityAdmin
	ServicePlanVisibilityOrganization
	ServicePlanVisibilitySpace
)

func (s ServicePlanVisibilityType) String() string {
	switch s {
	case ServicePlanVisibilityPublic:
		return "public"
	case ServicePlanVisibilityAdmin:
		return "admin"
	case ServicePlanVisibilityOrganization:
		return "organization"
	case ServicePlanVisibilitySpace:
		return "space"
	}
	return ""
}

func ParseServicePlanVisibilityType(visibilityType string) (ServicePlanVisibilityType, error) {
	switch visibilityType {
	case "public":
		return ServicePlanVisibilityPublic, nil
	case "admin":
		return ServicePlanVisibilityAdmin, nil
	case "organization":
		return ServicePlanVisibilityOrganization, nil
	case "space":
		return ServicePlanVisibilitySpace, nil
	}
	return ServicePlanVisibilityNone, fmt.Errorf("could not parse %s into a valid ServicePlanVisibilityType", visibilityType)
}

func NewServicePlanVisibilityUpdate(visibilityType ServicePlanVisibilityType) *ServicePlanVisibility {
	return &ServicePlanVisibility{
		Type: visibilityType.String(),
	}
}
