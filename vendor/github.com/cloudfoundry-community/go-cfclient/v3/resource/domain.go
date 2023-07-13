package resource

import (
	"time"
)

type Domain struct {
	GUID               string              `json:"guid"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	Name               string              `json:"name"`
	Internal           bool                `json:"internal"`
	RouterGroup        *Relationship       `json:"router_group"`
	SupportedProtocols []string            `json:"supported_protocols"`
	Relationships      DomainRelationships `json:"relationships"`
	Metadata           *Metadata           `json:"metadata"`
	Links              map[string]Link     `json:"links"`
}

type DomainCreate struct {
	Name string `json:"name"`

	Internal            *bool                `json:"internal,omitempty"`
	RouterGroup         *Relationship        `json:"router_group,omitempty"`
	Organization        *ToOneRelationship   `json:"organization,omitempty"`
	SharedOrganizations *ToManyRelationships `json:"shared_organizations,omitempty"`
	Metadata            *Metadata            `json:"metadata,omitempty"`
}

type DomainUpdate struct {
	Metadata *Metadata `json:"metadata"`
}

type DomainList struct {
	Pagination Pagination `json:"pagination"`
	Resources  []*Domain  `json:"resources"`
}

type DomainRelationships struct {
	Organization        ToOneRelationship   `json:"organization"`
	SharedOrganizations ToManyRelationships `json:"shared_organizations"`
}

func NewDomainCreate(name string) *DomainCreate {
	return &DomainCreate{
		Name: name,
	}
}

func NewDomainShare(orgGUID string) *ToManyRelationships {
	return &ToManyRelationships{
		Data: []Relationship{
			{
				GUID: orgGUID,
			},
		},
	}
}
