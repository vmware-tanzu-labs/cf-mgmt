package isosegment

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Manager interface {
	Apply() error
	Create() error
	Remove() error
	Entitle() error
	Unentitle() error
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
	GetSpaceByName(spaceName string, orgGuid string) (cfclient.Space, error)
	AddIsolationSegmentToOrg(isolationSegmentGUID, orgGUID string) error
	RemoveIsolationSegmentFromOrg(isolationSegmentGUID, orgGUID string) error
	AddIsolationSegmentToSpace(isolationSegmentGUID, spaceGUID string) error
	RemoveIsolationSegmentFromSpace(isolationSegmentGUID, spaceGUID string) error
	DefaultIsolationSegmentForOrg(orgGUID, isolationSegmentGUID string) error
	ResetDefaultIsolationSegmentForOrg(orgGUID string) error
	IsolationSegmentForSpace(spaceGUID, isolationSegmentGUID string) error
	ResetIsolationSegmentForSpace(spaceGUID string) error
}
