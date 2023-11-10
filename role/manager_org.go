package role

import (
	"context"
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) AddUserToOrg(orgGUID string, userName, userGUID string) error {
	if m.Peek {
		return nil
	}
	orgUsers, _, _, _, err := m.ListOrgUsersByRole(orgGUID)
	if err != nil {
		return err
	}
	if !orgUsers.HasUserForGUID(userName, userGUID) {
		_, err := m.RoleClient.CreateOrganizationRole(context.Background(), orgGUID, userGUID, resource.OrganizationRoleUser)
		if err != nil {
			lo.G.Debugf("Error adding user [%s] to org with guid [%s] but should have succeeded missing from org roles %+v, message: [%s]", userName, userGUID, orgUsers, err.Error())
			return err
		}
		orgUsers.AddUser(RoleUser{UserName: userName, GUID: userGUID})
		m.UpdateOrgRoleUsers(orgGUID, orgUsers)
	}
	return nil
}

func (m *DefaultManager) AssociateOrgAuditor(orgGUID, orgName, entityGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "auditor", orgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "auditor", orgName)
	_, err = m.RoleClient.CreateOrganizationRole(context.Background(), orgGUID, userGUID, resource.OrganizationRoleAuditor)
	return err
}

func (m *DefaultManager) AssociateOrgManager(orgGUID, orgName, entityGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "manager", orgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "manager", orgName)
	_, err = m.RoleClient.CreateOrganizationRole(context.Background(), orgGUID, userGUID, resource.OrganizationRoleManager)
	return err
}

func (m *DefaultManager) AssociateOrgBillingManager(orgGUID, orgName, entityGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: Add User %s to role %s for org %s", userName, "billing manager", orgName)
		return nil
	}

	lo.G.Infof("Add User %s to role %s for org %s", userName, "billing manager", orgName)
	_, err = m.RoleClient.CreateOrganizationRole(context.Background(), orgGUID, userGUID, resource.OrganizationRoleBillingManager)
	return err
}

func (m *DefaultManager) RemoveOrgAuditor(orgName, orgGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, orgName, "auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, orgName, "auditor")
	roleGUID, err := m.GetOrgRoleGUID(orgGUID, userGUID, resource.OrganizationRoleAuditor.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}
func (m *DefaultManager) RemoveOrgBillingManager(orgName, orgGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, orgName, "billing manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, orgName, "billing manager")
	roleGUID, err := m.GetOrgRoleGUID(orgGUID, userGUID, resource.OrganizationRoleBillingManager.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}

func (m *DefaultManager) RemoveOrgManager(orgName, orgGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, orgName, "manager")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, orgName, "manager")
	roleGUID, err := m.GetOrgRoleGUID(orgGUID, userGUID, resource.OrganizationRoleManager.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}

func (m *DefaultManager) RemoveOrgUser(orgName, orgGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org %s with role %s", userName, orgName, "user")
		return nil
	}
	lo.G.Infof("removing user %s from org %s with role %s", userName, orgName, "user")
	roleGUID, err := m.GetOrgRoleGUID(orgGUID, userGUID, resource.OrganizationRoleUser.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}

func (m *DefaultManager) GetOrgRoleGUID(orgGUID, userGUID, role string) (string, error) {
	orgs, ok := m.OrgRolesUsers[orgGUID]
	if !ok {
		return "", fmt.Errorf("org with guid[%s] has no roles", orgGUID)
	}
	roles, ok := orgs[role]
	if !ok {
		return "", fmt.Errorf("org with guid[%s] has no roles of type [%s]", orgGUID, role)
	}
	roleGUID, ok := roles[userGUID]
	if !ok {
		return "", fmt.Errorf("org with guid[%s] has no role of type [%s] with user guid [%s]", orgGUID, role, userGUID)
	}
	return roleGUID, nil
}
