package organization

import (
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/xchapter7x/lo"
)

// NewUserManager -
func NewUserManager(
	client CFClient,
	peek bool,
) UserMgr {
	return &UserManager{
		client: client,
		Peek:   peek,
	}
}

// UserManager -
type UserManager struct {
	client CFClient
	Peek   bool
}

func (m *UserManager) AddUserToOrg(userName, orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to orgGUID %s", userName, orgGUID)
		return nil
	}
	_, err := m.client.AssociateOrgUserByUsername(orgGUID, userName)
	return err
}

func (m *UserManager) RemoveOrgAuditorByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, orgGUID, "auditor")
		return nil
	}
	return m.client.RemoveOrgAuditorByUsername(orgGUID, userName)
}
func (m *UserManager) RemoveOrgBillingManagerByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, orgGUID, "billing manager")
		return nil
	}
	return m.client.RemoveOrgBillingManagerByUsername(orgGUID, userName)
}
func (m *UserManager) RemoveOrgManagerByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from GUID %s with role %s", userName, orgGUID, "manager")
		return nil
	}
	return m.client.RemoveOrgManagerByUsername(orgGUID, userName)
}
func (m *UserManager) ListOrgAuditors(orgGUID string) (map[string]string, error) {
	users, err := m.client.ListOrgAuditors(orgGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) ListOrgBillingManager(orgGUID string) (map[string]string, error) {
	users, err := m.client.ListOrgBillingManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) ListOrgManagers(orgGUID string) (map[string]string, error) {
	users, err := m.client.ListOrgManagers(orgGUID)
	if err != nil {
		return nil, err
	}
	return m.userListToMap(users), nil
}
func (m *UserManager) AssociateOrgAuditorByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org GUID %s", userName, "auditor", orgGUID)
		return nil
	}
	err := m.AddUserToOrg(userName, orgGUID)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateOrgAuditorByUsername(orgGUID, userName)
	return err
}
func (m *UserManager) AssociateOrgBillingManagerByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org GUID %s", userName, "billing manager", orgGUID)
		return nil
	}
	err := m.AddUserToOrg(userName, orgGUID)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateOrgBillingManagerByUsername(orgGUID, userName)
	return err
}

func (m *UserManager) AssociateOrgManagerByUsername(orgGUID, userName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org GUID %s", userName, "manager", orgGUID)
		return nil
	}
	err := m.AddUserToOrg(userName, orgGUID)
	if err != nil {
		return err
	}
	_, err = m.client.AssociateOrgManagerByUsername(orgGUID, userName)
	return err
}

func (m *UserManager) userListToMap(users []cfclient.User) map[string]string {
	userMap := make(map[string]string)
	for _, user := range users {
		userMap[strings.ToLower(user.Username)] = user.Guid
	}
	return userMap
}
