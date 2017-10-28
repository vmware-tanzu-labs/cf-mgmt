package isosegment

import (
	"fmt"
	"net/url"

	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	ccwrap "code.cloudfoundry.org/cli/api/cloudcontroller/wrapper"
	"code.cloudfoundry.org/cli/api/uaa"
	uaawrap "code.cloudfoundry.org/cli/api/uaa/wrapper"
	uaautil "code.cloudfoundry.org/cli/api/uaa/wrapper/util"
	"github.com/xchapter7x/lo"
)

type manager interface {
	GetIsolationSegments() ([]Segment, error)
	EntitledIsolationSegments(org string) ([]Segment, error)

	CreateIsolationSegment(name string) error
	DeleteIsolationSegment(segmentName string) error

	EnableOrgIsolation(orgName, segmentName string) error
	RevokeOrgIsolation(orgName, segmentName string) error

	SetOrgIsolationSegment(orgName string, s Segment) error
	SetSpaceIsolationSegment(orgName, spaceName string, s Segment) error
}

type ccv3Manager struct {
	cc *ccv3.Client

	// a cache mapping segment/org names to their GUIDs
	segments map[string]string
	orgs     map[string]string
}

func (c *ccv3Manager) EnableOrgIsolation(orgName, segmentName string) error {
	orgGUID, err := c.orgGUID(orgName)
	if err != nil {
		return err
	}
	segmentGUID, err := c.segmentGUID(segmentName)
	if err != nil {
		return err
	}

	_, _, err = c.cc.EntitleIsolationSegmentToOrganizations(segmentGUID, []string{orgGUID})
	return err
}

func (c *ccv3Manager) RevokeOrgIsolation(orgName, segmentName string) error {
	orgGUID, err := c.orgGUID(orgName)
	if err != nil {
		return err
	}
	segmentGUID, err := c.segmentGUID(segmentName)
	if err != nil {
		return err
	}

	_, err = c.cc.RevokeIsolationSegmentFromOrganization(segmentGUID, orgGUID)
	return err
}

// EntitledIsolationSegments gets the isolations segments that an org is entitled to.
func (c *ccv3Manager) EntitledIsolationSegments(org string) ([]Segment, error) {
	orgGUID, err := c.orgGUID(org)
	if err != nil {
		return nil, err
	}
	is, warnings, err := c.cc.GetIsolationSegments(url.Values{
		"organization_guids": []string{orgGUID},
	})
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		lo.G.Info(warnings)
	}

	result := make([]Segment, len(is))
	for i := range is {
		result[i] = Segment{
			Name: is[i].Name,
			GUID: is[i].GUID,
		}
	}
	return result, nil
}

func (c *ccv3Manager) GetIsolationSegments() ([]Segment, error) {
	is, warnings, err := c.cc.GetIsolationSegments(nil)
	if err != nil {
		return nil, err
	}

	if len(warnings) > 0 {
		lo.G.Info(warnings)
	}

	result := make([]Segment, len(is))
	for i := range is {
		result[i] = Segment{
			Name: is[i].Name,
			GUID: is[i].GUID,
		}
	}
	return result, nil
}

func (c *ccv3Manager) CreateIsolationSegment(name string) error {
	segment, _, err := c.cc.CreateIsolationSegment(ccv3.IsolationSegment{
		Name: name,
	})
	if err != nil {
		return err
	}
	c.segments[name] = segment.GUID
	return nil
}

func (c *ccv3Manager) DeleteIsolationSegment(segmentName string) error {
	guid, err := c.segmentGUID(segmentName)
	if err != nil {
		return err
	}
	_, err = c.cc.DeleteIsolationSegment(guid)
	if err != nil {
		return err
	}
	delete(c.segments, segmentName)
	return nil
}

// SetOrgIsolationSegment sets the default isolation segment for the specified org.
// If the segment name is empty, it resets the default isolation segment.
func (c *ccv3Manager) SetOrgIsolationSegment(orgName string, s Segment) error {
	orgGUID, err := c.orgGUID(orgName)
	if err != nil {
		return err
	}
	var segmentGUID string
	if s.GUID == "" && s.Name != "" {
		segmentGUID, err = c.segmentGUID(s.Name)
		if err != nil {
			return err
		}
	}
	_, err = c.cc.PatchOrganizationDefaultIsolationSegment(orgGUID, segmentGUID)
	return err
}

func (c *ccv3Manager) SetSpaceIsolationSegment(orgName, spaceName string, s Segment) error {
	orgGUID, err := c.orgGUID(orgName)
	if err != nil {
		return err
	}

	spaces, _, err := c.cc.GetSpaces(url.Values{
		"organization_guids": []string{orgGUID},
		"names":              []string{spaceName},
	})
	if err != nil {
		return err
	}
	if l := len(spaces); l != 1 {
		return fmt.Errorf("found %d spaces with name %s in org %s", l, spaceName, orgName)
	}

	var segmentGUID string
	if s.GUID == "" && s.Name != "" {
		segmentGUID, err = c.segmentGUID(s.Name)
		if err != nil {
			return err
		}
	}

	_, _, err = c.cc.AssignSpaceToIsolationSegment(spaces[0].GUID, segmentGUID)
	return err
}

func (c *ccv3Manager) orgGUID(name string) (string, error) {
	if guid, ok := c.orgs[name]; ok {
		return guid, nil
	}

	orgs, _, err := c.cc.GetOrganizations(url.Values{"names": []string{name}})
	if err != nil {
		return "", err
	}
	if l := len(orgs); l != 1 {
		return "", fmt.Errorf("found %d orgs with name %s", l, name)
	}
	return orgs[0].GUID, nil
}

func (c *ccv3Manager) segmentGUID(name string) (string, error) {
	if guid, ok := c.segments[name]; ok {
		return guid, nil
	}

	ss, _, err := c.cc.GetIsolationSegments(url.Values{"names": []string{name}})
	if err != nil {
		return "", err
	}
	if l := len(ss); l != 1 {
		return "", fmt.Errorf("found %d iso segments with name %s", l, name)
	}
	return ss[0].GUID, nil
}

// ccv3Client creates a client for the V3 Cloud Controller API.
func ccv3Client(cfmgmtVersion, cfURL, uaaToken string) (*ccv3.Client, error) {
	tokenCache := uaautil.NewInMemoryTokenCache()
	tokenCache.SetAccessToken("bearer " + uaaToken)

	uaaClient := uaa.NewClient(uaa.Config{
		AppName:           appName,
		AppVersion:        cfmgmtVersion,
		SkipSSLValidation: true, // TODO

		//ClientID:     "cf-mgmt",
		//ClientSecret: "cf-mgmt-secret",
	})
	uaaClient.WrapConnection(uaawrap.NewRetryRequest(2))

	ccClient := ccv3.NewClient(ccv3.Config{
		AppName:    appName,
		AppVersion: cfmgmtVersion,

		Wrappers: []ccv3.ConnectionWrapper{
			ccwrap.NewUAAAuthentication(uaaClient, tokenCache),
			ccwrap.NewRetryRequest(2),
		},
	})

	warnings, err := ccClient.TargetCF(ccv3.TargetSettings{
		SkipSSLValidation: true, // TODO
		URL:               cfURL,
	})
	if len(warnings) > 0 {
		lo.G.Warning(warnings)
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't target CF: %v", err)
	}

	return ccClient, nil
}
