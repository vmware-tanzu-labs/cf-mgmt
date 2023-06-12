package space

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (cfclient.Space, error)
	CreateSpaces() error
	UpdateSpaces() (err error)
	DeleteSpaces() (err error)
	DeleteSpacesForOrg(orgGUID, orgName string) (err error)
	ListSpaces(orgGUID string) ([]cfclient.Space, error)
	UpdateSpacesMetadata() error
}

type CFClient interface {
	GetSpaceByGuid(spaceGUID string) (cfclient.Space, error)
	UpdateSpace(spaceGUID string, req cfclient.SpaceRequest) (cfclient.Space, error)
	CreateSpace(req cfclient.SpaceRequest) (cfclient.Space, error)
	DeleteSpace(guid string, recursive, async bool) error
	ListSpaces() ([]cfclient.Space, error)
	SupportsMetadataAPI() (bool, error)
	UpdateSpaceMetadata(spaceGUID string, metadata cfclient.Metadata) error
	SpaceMetadata(spaceGUID string) (*cfclient.Metadata, error)
	RemoveSpaceMetadata(spaceGUID string) error
	ListOrgs() ([]cfclient.Org, error)
}
