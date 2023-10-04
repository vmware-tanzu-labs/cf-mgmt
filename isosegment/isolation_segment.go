package isosegment

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(client CFClient, cfg config.Reader, orgReader organizationreader.Reader, spaceManager space.Manager, peek bool) (Manager, error) {
	globalCfg, err := cfg.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	return &Updater{
		Cfg:          cfg,
		Client:       client,
		OrgReader:    orgReader,
		SpaceManager: spaceManager,
		Peek:         peek,
		CleanUp:      globalCfg.EnableDeleteIsolationSegments,
	}, nil
}

// Updater performs the required updates to acheive the desired state wrt isolation segments.
// Updaters should always be created with NewUpdater.  It is save to modify Updater's
// exported fields after creation.
type Updater struct {
	Cfg          config.Reader
	Client       CFClient
	OrgReader    organizationreader.Reader
	SpaceManager space.Manager
	Peek         bool
	CleanUp      bool
}

func (u *Updater) Apply() error {
	lo.G.Debugf("Creating iso segments")
	if err := u.Create(); err != nil {
		return err
	}
	lo.G.Debugf("entitling iso segments")
	if err := u.Entitle(); err != nil {
		return err
	}
	lo.G.Debugf("update orgs")
	if err := u.UpdateOrgs(); err != nil {
		return err
	}
	lo.G.Debugf("update spaces")
	if err := u.UpdateSpaces(); err != nil {
		return err
	}
	lo.G.Debugf("unentitling iso segments")
	if err := u.Unentitle(); err != nil {
		return err
	}
	lo.G.Debugf("removing iso segments")
	if err := u.Remove(); err != nil {
		return err
	}
	return nil
}

// Create creates any isolation segments that do not yet exist,
func (u *Updater) Create() error {
	desired, err := u.allDesiredSegments()
	if err != nil {
		return err
	}
	current, err := u.Client.ListIsolationSegments()
	if err != nil {
		return err
	}

	c := classify(desired, current)
	for i := range c.missing {
		err := u.create(&c.missing[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// Create creates any isolation segments that do not yet exist,
func (u *Updater) Remove() error {
	desired, err := u.allDesiredSegments()
	if err != nil {
		return err
	}
	current, err := u.Client.ListIsolationSegments()
	if err != nil {
		return err
	}

	c := classify(desired, current)
	for i := range c.extra {
		err := u.delete(&c.extra[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Updater) Unentitle() error {
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
	sm := make(map[string][]*cfclient.IsolationSegment)
	for _, space := range spaces {
		org, err := u.OrgReader.FindOrg(space.Org)
		if err != nil {
			return errors.Wrap(err, "finding org for space configs")
		}
		if s := space.IsoSegment; s != "" {
			if isosegment, ok := isolationSegmentsMap[s]; ok {
				sm[org.GUID] = append(sm[org.GUID], &isosegment)
			} else {
				if !u.Peek {
					return fmt.Errorf("Isolation segment [%s] does not exist", s)
				}
			}
		} else {
			sm[org.GUID] = append(sm[org.GUID], nil)
		}
	}
	for _, orgConfig := range orgs {
		org, err := u.OrgReader.FindOrg(orgConfig.Org)
		if err != nil {
			return errors.Wrap(err, "finding org for org configs")
		}
		if s := orgConfig.DefaultIsoSegment; s != "" {
			if isosegment, ok := isolationSegmentsMap[s]; ok {
				sm[org.GUID] = append(sm[org.GUID], &isosegment)
			} else {
				if !u.Peek {
					return fmt.Errorf("Isolation segment [%s] does not exist", s)
				}
			}
		} else {
			sm[org.GUID] = append(sm[org.GUID], nil)
		}
	}

	for orgGUID, segments := range sm {
		orgIsolationSegments, err := u.Client.ListIsolationSegmentsByQuery(url.Values{
			"organization_guids": []string{orgGUID},
		})
		if err != nil {
			return err
		}

		var desiredSegments []cfclient.IsolationSegment
		for _, segment := range segments {
			if segment != nil {
				desiredSegments = append(desiredSegments, *segment)
			}
		}
		c := classify(desiredSegments, orgIsolationSegments)
		for i := range c.extra {
			err := u.revoke(&c.extra[i], orgGUID)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
				org, err := u.OrgReader.FindOrg(space.Org)
				if err != nil {
					return errors.Wrap(err, "finding org for space configs in entitle")
				}
				sm[org.GUID] = append(sm[org.GUID], isosegment)
			} else {
				if !u.Peek {
					return fmt.Errorf("Isolation segment [%s] does not exist", s)
				}
			}
		}
	}
	for _, orgConfig := range orgs {
		if s := orgConfig.DefaultIsoSegment; s != "" {
			org, err := u.OrgReader.FindOrg(orgConfig.Org)
			if err != nil {
				return errors.Wrap(err, "finding org for org configs in entitle")
			}
			if isosegment, ok := isolationSegmentsMap[s]; ok {
				sm[org.GUID] = append(sm[org.GUID], isosegment)
			} else {
				if !u.Peek {
					return fmt.Errorf("Isolation segment [%s] does not exist", s)
				}
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
		for i := range c.missing {
			err := u.entitle(&c.missing[i], orgGUID)
			if err != nil {
				return err
			}
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
		org, err := u.OrgReader.FindOrg(oc.Org)
		if err != nil {
			return errors.Wrap(err, "finding org for org configs in update orgs")
		}
		isolationSegmentMap, err := u.isolationSegmentMap()
		if err != nil {
			return err
		}

		isolationSegmentGUID, err := u.getIsolationSegmentGUID(oc.DefaultIsoSegment, isolationSegmentMap)
		if err != nil {
			return err
		}
		orgIsolationSegmentGUID, err := u.OrgReader.GetDefaultIsolationSegment(org)
		if err != nil {
			return err
		}
		if orgIsolationSegmentGUID != isolationSegmentGUID {
			if u.Peek {
				if isolationSegmentGUID != "" {
					lo.G.Infof("[dry-run]: set default isolation segment for org %s to %s", oc.Org, oc.DefaultIsoSegment)
				} else {
					lo.G.Infof("[dry-run]: reset default isolation segment for org %s", oc.Org)
				}
				continue
			}

			if isolationSegmentGUID != "" {
				lo.G.Infof("set default isolation segment for org %s to %s", oc.Org, oc.DefaultIsoSegment)
				err = u.Client.DefaultIsolationSegmentForOrg(org.GUID, isolationSegmentGUID)
				if err != nil {
					return err
				}
			} else {
				lo.G.Infof("reset default isolation segment for org %s", oc.Org)
				err = u.Client.ResetDefaultIsolationSegmentForOrg(org.GUID)
				if err != nil {
					return err
				}
			}

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
		if u.Peek {
			return fmt.Sprintf("%s-dry-run-isosegment-guid", isolationSegmentName), nil
		} else {
			return "", fmt.Errorf("Isolation Segment [%s] not found", isolationSegmentName)
		}
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
		space, err := u.SpaceManager.FindSpace(sc.Org, sc.Space)
		if err != nil {
			return err
		}
		isolationSegmentGUID, err := u.getIsolationSegmentGUID(sc.IsoSegment, isolationSegmentMap)
		if err != nil {
			return err
		}
		spaceIsoSegGUID, err := u.SpaceManager.GetSpaceIsolationSegmentGUID(space)
		if err != nil {
			return err
		}
		if spaceIsoSegGUID != isolationSegmentGUID {
			if u.Peek {
				if sc.IsoSegment != "" {
					lo.G.Infof("[dry-run]: set isolation segment for space %s to %s (org %s)", sc.Space, sc.IsoSegment, sc.Org)
				} else {
					lo.G.Infof("[dry-run]: reset isolation segment for space %s (org %s)", sc.Space, sc.Org)
				}
				continue
			}
			if sc.IsoSegment != "" {
				lo.G.Infof("set isolation segment for space %s to %s (org %s)", sc.Space, sc.IsoSegment, sc.Org)
				err = u.Client.IsolationSegmentForSpace(space.GUID, isolationSegmentGUID)
				if err != nil {
					return err
				}
			} else {
				lo.G.Infof("reset isolation segment for space %s (org %s)", sc.Space, sc.Org)
				err = u.Client.ResetIsolationSegmentForSpace(space.GUID)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}

func (u *Updater) create(s *cfclient.IsolationSegment) error {
	if u.Peek {
		lo.G.Info("[dry-run]: create segment", s.Name)
		return nil
	}

	lo.G.Info("create segment", s.Name)
	_, err := u.Client.CreateIsolationSegment(s.Name)
	return err
}

func (u *Updater) delete(s *cfclient.IsolationSegment) error {
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
	lo.G.Infof("entitle org %s to iso segment %s", orgGUID, s.Name)
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
	lo.G.Infof("revoke iso segment %s (%s) from org %s", s.Name, s.GUID, orgGUID)
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
