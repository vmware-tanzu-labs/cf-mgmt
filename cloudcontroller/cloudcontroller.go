package cloudcontroller

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/xchapter7x/lo"
)

func NewManager(client *cfclient.Client, peek bool) Manager {
	return &DefaultManager{
		Client: *client,
		Peek:   peek,
	}
}

//ListIsolationSegments : Returns all isolation segments
func (m *DefaultManager) ListIsolationSegments() ([]cfclient.IsolationSegment, error) {
	isolationSegments, err := m.Client.ListIsolationSegments()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total isolation segments returned :", len(isolationSegments))
	return isolationSegments, nil
}

func (m *DefaultManager) OrgQuotaByName(name string) (cfclient.OrgQuota, error) {
	return m.Client.GetOrgQuotaByName(name)
}
func (m *DefaultManager) SpaceQuotaByName(name string) (cfclient.SpaceQuota, error) {
	return m.Client.GetSpaceQuotaByName(name)
}
