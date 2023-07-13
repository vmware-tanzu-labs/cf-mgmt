package resource

import (
	"time"
)

type Meta struct {
	GUID      string    `json:"guid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
}

// Pagination is used by the apis to page list results
type Pagination struct {
	TotalResults int  `json:"total_results"`
	TotalPages   int  `json:"total_pages"`
	First        Link `json:"first"`
	Last         Link `json:"last"`
	Next         Link `json:"next"`
	Previous     Link `json:"previous"`
}

// Link is a HATEOAS-style link for apis
type Link struct {
	Href   string `json:"href"`
	Method string `json:"method,omitempty"`
}

type SpaceRelationship struct {
	Space ToOneRelationship `json:"space"`
}

type AppRelationship struct {
	App ToOneRelationship `json:"app"`
}

// ToOneRelationship is a relationship to a single object
type ToOneRelationship struct {
	Data *Relationship `json:"data"`
}

// ToManyRelationships is a relationship to multiple objects
type ToManyRelationships struct {
	Data []Relationship `json:"data"`
}

type Relationship struct {
	GUID string `json:"guid,omitempty"`
}

type NullableToOneRelationship struct {
	Data *NullableRelationship `json:"data"`
}

type NullableRelationship struct {
	GUID *string `json:"guid"`
}

type LastOperation struct {
	Type        string    `json:"type"`
	State       string    `json:"state"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewToManyRelationships(guids []string) *ToManyRelationships {
	r := &ToManyRelationships{
		Data: make([]Relationship, len(guids)),
	}
	for i, g := range guids {
		r.Data[i] = Relationship{
			GUID: g,
		}
	}
	return r
}
