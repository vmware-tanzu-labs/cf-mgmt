package isosegment

import (
	"fmt"
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(client CFClient, cfg config.Reader, peek bool) (Manager, error) {
	globalCfg, err := cfg.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	return &Updater{
		Cfg:     cfg,
		Client:  client,
		Peek:    peek,
		CleanUp: globalCfg.EnableDeleteIsolationSegments,
	}, nil
}

// Updater performs the required updates to acheive the desired state wrt isolation segments.
// Updaters should always be created with NewUpdater.  It is save to modify Updater's
// exported fields after creation.
type Updater struct {
	Cfg     config.Reader
	Client  CFClient
	Peek    bool
	CleanUp bool
}

// Ensure creates any isolation segments that do not yet exist,
// and optionally removes unneeded isolation segments.
func (u *Updater) Ensure() error {
	desired, err := u.allDesiredSegments()
	if err != nil {
		return err
	}
	current, err := u.Client.ListIsolationSegments()
	if err != nil {
		return err
	}

	c := classify(desired, current)
	return c.update("", u.create, u.delete)
}

// Entitle ensures that each org is entitled to the isolation segments it needs to use.
func (u *Updater) Entitle() error {
	spaces, err := u.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	orgs, err := u.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	isolationSegmentsMap, err := u.isolationSegmentMap()
	if err != nil {
		return err
	}
	// build up a list of segments required by each org (grouped by org name)
	// this includes segments used by all of the orgs spaces, as well as the
	// org's default segment
	sm := make(map[string][]cfclient.IsolationSegment)
	for _, space := range spaces {
		if s := space.IsoSegment; s != "" {
			if isosegment, ok := isolationSegmentsMap[s]; ok {
				org, err := u.Client.GetOrgByName(space.Org)
				if err != nil {
					return err
				}
				sm[org.Guid] = append(sm[org.Guid], isosegment)
			} else {
				return fmt.Errorf("Isolation segment [%s] does not exist", s)
			}
		}
	}
	for _, orgConfig := range orgs {
		if s := orgConfig.DefaultIsoSegment; s != "" {
			org, err := u.Client.GetOrgByName(orgConfig.Org)
			if err != nil {
				return err
			}
			if isosegment, ok := isolationSegmentsMap[s]; ok {
				sm[org.Guid] = append(sm[org.Guid], isosegment)
			} else {
				return fmt.Errorf("Isolation segment [%s] does not exist", s)
			}
		}
	}

	for orgGUID, desiredSegments := range sm {
		orgIsolationSegments, err := u.Client.ListIsolationSegmentsByQuery(url.Values{
			"organization_guids": []string{orgGUID},
		})
		if err != nil {
			return err
		}

		c := classify(desiredSegments, orgIsolationSegments)
		err = c.update(orgGUID, u.entitle, u.revoke)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateOrgs sets the default isolation segment for each org,
// as specified in the cf-mgmt config.
func (u *Updater) UpdateOrgs() error {
	ocs, err := u.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}
	for _, oc := range ocs {
		if u.Peek {
			if oc.DefaultIsoSegment != "" {
				lo.G.Infof("[dry-run]: set default isolation segment for org %s to %s", oc.Org, oc.DefaultIsoSegment)
			} else {
				lo.G.Infof("[dry-run]: reset default isolation segment for org %s", oc.Org)
			}
			continue
		}
		org, err := u.Client.GetOrgByName(oc.Org)
		if err != nil {
			return err
		}
		isolationSegmentMap, err := u.isolationSegmentMap()
		if err != nil {
			return err
		}

		isolationSegmentGUID, err := u.getIsolationSegmentGUID(oc.DefaultIsoSegment, isolationSegmentMap)
		if err != nil {
			return err
		}
		orgRequest := cfclient.OrgRequest{
			Name: org.Name,
			DefaultIsolationSegmentGuid: isolationSegmentGUID,
		}
		_, err = u.Client.UpdateOrg(org.Guid, orgRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Updater) getIsolationSegmentGUID(isolationSegmentName string, isolationSegmentMap map[string]cfclient.IsolationSegment) (string, error) {
	if isolationSegmentName == "" {
		return "", nil
	}
	if isosegment, ok := isolationSegmentMap[isolationSegmentName]; ok {
		return isosegment.GUID, nil
	} else {
		return "", fmt.Errorf("Isolation Segment [%s] not found", isolationSegmentName)
	}
}

// UpdateSpaces sets the isolation segment for each space,
// as specified in the cf-mgmt config.
func (u *Updater) UpdateSpaces() error {
	scs, err := u.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	isolationSegmentMap, err := u.isolationSegmentMap()
	if err != nil {
		return err
	}
	for _, sc := range scs {
		if u.Peek {
			if sc.IsoSegment != "" {
				lo.G.Infof("[dry-run]: set isolation segment for space %s to %s (org %s)", sc.Space, sc.IsoSegment, sc.Org)
			} else {
				lo.G.Infof("[dry-run]: reset isolation segment for space %s (org %s)", sc.Space, sc.Org)
			}
			continue
		}
		org, err := u.Client.GetOrgByName(sc.Org)
		if err != nil {
			return err
		}
		space, err := u.Client.GetSpaceByName(sc.Space, org.Guid)
		if err != nil {
			return err
		}
		isolationSegmentGUID, err := u.getIsolationSegmentGUID(sc.IsoSegment, isolationSegmentMap)
		if err != nil {
			return err
		}

		spaceRequest := cfclient.SpaceRequest{
			Name:                 space.Name,
			OrganizationGuid:     org.Guid,
			IsolationSegmentGuid: isolationSegmentGUID,
		}
		_, err = u.Client.UpdateSpace(space.Guid, spaceRequest)
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Updater) create(s *cfclient.IsolationSegment, _ string) error {
	if u.Peek {
		lo.G.Info("[dry-run]: create segment", s.Name)
		return nil
	}
	_, err := u.Client.CreateIsolationSegment(s.Name)
	return err
}

func (u *Updater) delete(s *cfclient.IsolationSegment, _ string) error {
	if !u.CleanUp {
		return nil
	}
	if s.Name == "shared" {
		return nil
	}
	if u.Peek {
		lo.G.Infof("[dry-run]: delete segment %s (%s)", s.Name, s.GUID)
		return nil
	}
	lo.G.Infof("delete segment %s (%s)", s.Name, s.GUID)
	return u.Client.DeleteIsolationSegmentByGUID(s.GUID)
}

func (u *Updater) entitle(s *cfclient.IsolationSegment, orgGUID string) error {
	if u.Peek {
		lo.G.Infof("[dry-run]: entitle org %s to iso segment %s", orgGUID, s.Name)
		return nil
	}
	return u.Client.AddIsolationSegmentToOrg(s.GUID, orgGUID)
}

func (u *Updater) revoke(s *cfclient.IsolationSegment, orgGUID string) error {
	if !u.CleanUp {
		return nil
	}
	if u.Peek {
		lo.G.Infof("[dry-run]: revoke iso segment %s from org %s", s.Name, orgGUID)
		return nil
	}
	return u.Client.RemoveIsolationSegmentFromOrg(s.GUID, orgGUID)
}

// allDesiredSegments iterates through the cf-mgmt configuration for all
// orgs and spaces and builds the complete set of isolation segments that
// should exist
func (u *Updater) allDesiredSegments() ([]cfclient.IsolationSegment, error) {
	orgs, err := u.Cfg.GetOrgConfigs()
	if err != nil {
		return nil, err
	}
	spaces, err := u.Cfg.GetSpaceConfigs()
	if err != nil {
		return nil, err
	}

	segments := make(map[string]struct{})
	for _, org := range orgs {
		if org.DefaultIsoSegment != "" {
			segments[org.DefaultIsoSegment] = struct{}{}
		}
	}
	for _, space := range spaces {
		if space.IsoSegment != "" {
			segments[space.IsoSegment] = struct{}{}
		}
	}

	result := make([]cfclient.IsolationSegment, 0, len(segments))
	for k := range segments {
		result = append(result, cfclient.IsolationSegment{Name: k})
	}
	return result, nil
}

func (u *Updater) isolationSegmentMap() (map[string]cfclient.IsolationSegment, error) {
	isolationSegments, err := u.ListIsolationSegments()
	if err != nil {
		return nil, err
	}

	isolationSegmentsMap := make(map[string]cfclient.IsolationSegment)
	for _, isosegment := range isolationSegments {
		isolationSegmentsMap[isosegment.Name] = isosegment
	}
	return isolationSegmentsMap, nil
}

func (m *Updater) ListIsolationSegments() ([]cfclient.IsolationSegment, error) {
	isolationSegments, err := m.Client.ListIsolationSegments()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total isolation segments returned :", len(isolationSegments))
	return isolationSegments, nil
}
