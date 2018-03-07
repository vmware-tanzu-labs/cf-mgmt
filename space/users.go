package space

import (
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/xchapter7x/lo"
)

// NewUserManager -
func NewUserManager(
	client CFClient,
	peek bool) UserMgr {
	return &UserManager{
		client: client,
		Peek:   peek,
	}
}

type UserManager struct {
	client CFClient
	Peek   bool
}

func (m *UserManager) RemoveSpaceAuditorByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Auditor")
		return nil
	}
	return m.client.RemoveSpaceAuditorByUsername(spaceGUID, userName)
}
func (m *UserManager) RemoveSpaceDeveloperByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Developer")
		return nil
	}
	return m.client.RemoveSpaceDeveloperByUsername(spaceGUID, userName)
}
func (m *UserManager) RemoveSpaceManagerByUsername(spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, spaceGUID, "Manager")
		return nil
	}
	return m.client.RemoveSpaceManagerByUsername(spaceGUID, userName)
}
func (m *UserManager) ListSpaceAuditors(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceAuditors(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) ListSpaceDevelopers(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceDevelopers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) ListSpaceManagers(spaceGUID string) (map[string]string, error) {
	users, err := m.client.ListSpaceManagers(spaceGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) associateOrgUserByUsername(orgGUID, userName string) error {
	_, err := m.client.AssociateOrgUserByUsername(orgGUID, userName)
	return err
}

func (m *UserManager) userListToMap(users []cfclient.User) map[string]string {
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[strings.ToLower(user.Username)] = user.Guid
	}
	return userMap
}

func (m *UserManager) AssociateSpaceAuditorByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "auditor", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateSpaceAuditorByUsername(spaceGUID, userName)
	return err
}
func (m *UserManager) AssociateSpaceDeveloperByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "developer", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateSpaceDeveloperByUsername(spaceGUID, userName)
	return err
}
func (m *UserManager) AssociateSpaceManagerByUsername(orgGUID, spaceGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for spaceGUID %s", userName, "manager", spaceGUID)
		return nil
	}
	err := m.associateOrgUserByUsername(orgGUID, userName)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateSpaceManagerByUsername(spaceGUID, userName)
	return err
}
