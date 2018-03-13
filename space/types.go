package space

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (cfclient.Space, error)
	CreateSpaces(configDir, ldapBindPassword string) error
	UpdateSpaces(configDir string) (err error)
	DeleteSpaces(configFile string) (err error)
	ListSpaces(orgGUID string) ([]cfclient.Space, error)
}

type CFClient interface {
	GetSpaceByGuid(spaceGUID string) (cfclient.Space, error)
	UpdateSpace(spaceGUID string, req cfclient.SpaceRequest) (cfclient.Space, error)
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)
	CreateSpace(req cfclient.SpaceRequest) (cfclient.Space, error)
	DeleteSpace(guid string, recursive, async bool) error
}
