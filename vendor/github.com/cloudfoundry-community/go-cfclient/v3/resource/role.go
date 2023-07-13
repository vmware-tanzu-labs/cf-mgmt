package resource

import "time"

// Role implements role object. Roles control access to resources in organizations and spaces. Roles are assigned to users.
type Role struct {
	GUID          string                                 `json:"guid"`
	CreatedAt     time.Time                              `json:"created_at"`
	UpdatedAt     time.Time                              `json:"updated_at"`
	Type          string                                 `json:"type,omitempty"`
	Relationships RoleSpaceUserOrganizationRelationships `json:"relationships,omitempty"`
	Links         map[string]Link                        `json:"links,omitempty"`
}

type RoleList struct {
	Pagination Pagination    `json:"pagination"`
	Resources  []*Role       `json:"resources"`
	Included   *RoleIncluded `json:"included"`
}

type RoleSpaceCreate struct {
	RoleType      string                     `json:"type"`
	Relationships RoleSpaceUserRelationships `json:"relationships"`
}

type RoleOrganizationCreate struct {
	RoleType      string                            `json:"type"`
	Relationships RoleOrganizationUserRelationships `json:"relationships"`
}

type RoleSpaceUserRelationships struct {
	Space ToOneRelationship `json:"space"`
	User  ToOneRelationship `json:"user"`
}

type RoleOrganizationUserRelationships struct {
	Org  ToOneRelationship `json:"organization"`
	User ToOneRelationship `json:"user"`
}

type RoleSpaceUserOrganizationRelationships struct {
	Space ToOneRelationship `json:"space"`
	User  ToOneRelationship `json:"user"`
	Org   ToOneRelationship `json:"organization"`
}

type RoleWithIncluded struct {
	Role
	Included *RoleIncluded `json:"included"`
}

type RoleIncluded struct {
	Users         []*User         `json:"users"`
	Organizations []*Organization `json:"organizations"`
	Spaces        []*Space        `json:"spaces"`
}

// SpaceRoleType https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#valid-role-types
type SpaceRoleType int

const (
	SpaceRoleNone SpaceRoleType = iota
	SpaceRoleAuditor
	SpaceRoleDeveloper
	SpaceRoleManager
	SpaceRoleSupporter
)

func (sr SpaceRoleType) String() string {
	switch sr {
	case SpaceRoleAuditor:
		return "space_auditor"
	case SpaceRoleDeveloper:
		return "space_developer"
	case SpaceRoleManager:
		return "space_manager"
	case SpaceRoleSupporter:
		return "space_supporter"
	}
	return ""
}

// OrganizationRoleType https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#valid-role-types
type OrganizationRoleType int

const (
	OrganizationRoleNone OrganizationRoleType = iota
	OrganizationRoleUser
	OrganizationRoleAuditor
	OrganizationRoleManager
	OrganizationRoleBillingManager
)

func (or OrganizationRoleType) String() string {
	switch or {
	case OrganizationRoleUser:
		return "organization_user"
	case OrganizationRoleAuditor:
		return "organization_auditor"
	case OrganizationRoleManager:
		return "organization_manager"
	case OrganizationRoleBillingManager:
		return "organization_billing_manager"
	}
	return ""
}

// RoleIncludeType https://v3-apidocs.cloudfoundry.org/version/3.126.0/index.html#include
type RoleIncludeType int

const (
	RoleIncludeNone RoleIncludeType = iota
	RoleIncludeUser
	RoleIncludeSpace
	RoleIncludeOrganization
)

func (r RoleIncludeType) String() string {
	switch r {
	case RoleIncludeUser:
		return "user"
	case RoleIncludeSpace:
		return "space"
	case RoleIncludeOrganization:
		return "organization"
	}
	return ""
}

func NewRoleSpaceCreate(spaceGUID, userGUID string, roleType SpaceRoleType) *RoleSpaceCreate {
	return &RoleSpaceCreate{
		RoleType: roleType.String(),
		Relationships: RoleSpaceUserRelationships{
			Space: ToOneRelationship{
				Data: &Relationship{
					GUID: spaceGUID,
				},
			},
			User: ToOneRelationship{
				Data: &Relationship{
					GUID: userGUID,
				},
			},
		},
	}
}

func NewRoleOrganizationCreate(orgGUID, userGUID string, roleType OrganizationRoleType) *RoleOrganizationCreate {
	return &RoleOrganizationCreate{
		RoleType: roleType.String(),
		Relationships: RoleOrganizationUserRelationships{
			Org: ToOneRelationship{
				Data: &Relationship{
					GUID: orgGUID,
				},
			},
			User: ToOneRelationship{
				Data: &Relationship{
					GUID: userGUID,
				},
			},
		},
	}
}
