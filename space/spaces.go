package space

import (
	"fmt"
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(client CFClient, uaaMgr uaa.Manager,
	orgMgr organization.Manager,
	cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:    cfg,
		UAAMgr: uaaMgr,
		Client: client,
		OrgMgr: orgMgr,
		Peek:   peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg    config.Reader
	Client CFClient
	UAAMgr uaa.Manager
	OrgMgr organization.Manager
	Peek   bool
}

func (m *DefaultManager) UpdateSpaceSSH(sshAllowed bool, space cfclient.Space) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: setting sshAllowed to %v for space %s", sshAllowed, space.Name)
		return nil
	}

	_, err := m.Client.UpdateSpace(space.Guid, cfclient.SpaceRequest{
		Name:             space.Name,
		AllowSSH:         sshAllowed,
		OrganizationGuid: space.OrganizationGuid,
	})
	return err
}

//UpdateSpaces -
func (m *DefaultManager) UpdateSpaces() error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		space, err := m.FindSpace(input.Org, input.Space)
		if err != nil {
			continue
		}
		lo.G.Debug("Processing space", space.Name)
		if input.AllowSSH != space.AllowSSH {
			if err := m.UpdateSpaceSSH(input.AllowSSH, space); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) ListSpaces(orgGUID string) ([]cfclient.Space, error) {
	spaces, err := m.Client.ListSpacesByQuery(url.Values{
		"organization_guid": []string{orgGUID},
	})
	if err != nil {
		return nil, err
	}
	return spaces, err

}

//FindSpace -
func (m *DefaultManager) FindSpace(orgName, spaceName string) (cfclient.Space, error) {
	orgGUID, err := m.OrgMgr.GetOrgGUID(orgName)
	if err != nil {
		return cfclient.Space{}, err
	}
	spaces, err := m.ListSpaces(orgGUID)
	if err != nil {
		return cfclient.Space{}, err
	}
	for _, theSpace := range spaces {
		if theSpace.Name == spaceName {
			return theSpace, nil
		}
	}
	if m.Peek {
		return cfclient.Space{
			Name:             spaceName,
			Guid:             fmt.Sprintf("%s-dry-run-space-guid", spaceName),
			OrganizationGuid: fmt.Sprintf("%s-dry-run-org-guid", orgName),
		}, nil
	}
	return cfclient.Space{}, fmt.Errorf("space [%s] not found in org [%s]", spaceName, orgName)
}

func (m *DefaultManager) CreateSpace(spaceName, orgName, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: create space %s for org %s", spaceName, orgName)
		return nil
	}
	lo.G.Infof("create space %s for org %s", spaceName, orgName)
	_, err := m.Client.CreateSpace(cfclient.SpaceRequest{
		Name:             spaceName,
		OrganizationGuid: orgGUID,
	})
	return err
}

//CreateSpaces -
func (m *DefaultManager) CreateSpaces() error {
	configSpaceList, err := m.Cfg.Spaces()
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {
		if len(input.Spaces) == 0 {
			continue
		}
		orgGUID, err := m.OrgMgr.GetOrgGUID(input.Org)
		if err != nil {
			return err
		}
		spaces, err := m.ListSpaces(orgGUID)
		if err != nil {
			continue
		}
		for _, spaceName := range input.Spaces {
			if m.doesSpaceExist(spaces, spaceName) {
				lo.G.Debugf("[%s] space already exists", spaceName)
				continue
			}
			if err = m.CreateSpace(spaceName, input.Org, orgGUID); err != nil {
				lo.G.Error(err)
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) doesSpaceExist(spaces []cfclient.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Name == spaceName {
			return true
		}
	}
	return false
}

func (m *DefaultManager) DeleteSpaces() error {
	configSpaceList, err := m.Cfg.Spaces()
	if err != nil {
		return err
	}
	for _, input := range configSpaceList {

		if !input.EnableDeleteSpaces {
			lo.G.Debugf("Space deletion is not enabled for %s.  Set enable-delete-spaces: true in spaces.yml", input.Org)
			continue //Skip all orgs that have not opted-in
		}

		configuredSpaces := make(map[string]bool)
		for _, spaceName := range input.Spaces {
			configuredSpaces[spaceName] = true
		}

		org, err := m.OrgMgr.FindOrg(input.Org)
		if err != nil {
			return err
		}
		spaces, err := m.ListSpaces(org.Guid)
		if err != nil {
			return err
		}

		spacesToDelete := make([]cfclient.Space, 0)
		for _, space := range spaces {
			if _, exists := configuredSpaces[space.Name]; !exists {
				spacesToDelete = append(spacesToDelete, space)
			}
		}

		for _, space := range spacesToDelete {
			lo.G.Infof("Deleting [%s] space in org %s", space.Name, input.Org)
			if err := m.DeleteSpace(space.Guid); err != nil {
				return err
			}
		}

	}

	return nil
}

//DeleteSpace - deletes a space based on GUID
func (m *DefaultManager) DeleteSpace(spaceGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete space with GUID %s", spaceGUID)
		return nil
	}
	return m.Client.DeleteSpace(spaceGUID, true, true)
}
