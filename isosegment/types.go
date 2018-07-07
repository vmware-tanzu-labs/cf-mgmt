package isosegment

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Manager interface {
	Apply() error
	Ensure() error
	Entitle() error
	UpdateOrgs() error
	UpdateSpaces() error
	ListIsolationSegments() ([]cfclient.IsolationSegment, error)
}

type CFClient interface {
	ListIsolationSegments() ([]cfclient.IsolationSegment, error)
	ListIsolationSegmentsByQuery(query url.Values) ([]cfclient.IsolationSegment, error)
	CreateIsolationSegment(name string) (*cfclient.IsolationSegment, error)
	DeleteIsolationSegmentByGUID(guid string) error
	GetIsolationSegmentByGUID(guid string) (*cfclient.IsolationSegment, error)
	GetOrgByName(name string) (cfclient.Org, error)
	UpdateOrg(orgGUID string, orgRequest cfclient.OrgRequest) (cfclient.Org, error)
	UpdateSpace(spaceGUID string, req cfclient.SpaceRequest) (cfclient.Space, error)
	GetSpaceByName(spaceName string, orgGuid string) (cfclient.Space, error)
	AddIsolationSegmentToOrg(isolationSegmentGUID, orgGUID string) error
	RemoveIsolationSegmentFromOrg(isolationSegmentGUID, orgGUID string) error
	AddIsolationSegmentToSpace(isolationSegmentGUID, spaceGUID string) error
	RemoveIsolationSegmentFromSpace(isolationSegmentGUID, spaceGUID string) error
}
