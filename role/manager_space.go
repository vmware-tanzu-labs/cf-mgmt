package role

import (
	"context"
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) RemoveSpaceAuditor(spaceName, spaceGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s with role %s", userName, spaceName, "Auditor")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s with role %s", userName, spaceName, "Auditor")
	roleGUID, err := m.GetSpaceRoleGUID(spaceGUID, userGUID, resource.SpaceRoleAuditor.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}
func (m *DefaultManager) RemoveSpaceDeveloper(spaceName, spaceGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s with role %s", userName, spaceName, "Developer")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s with role %s", userName, spaceName, "Developer")
	roleGUID, err := m.GetSpaceRoleGUID(spaceGUID, userGUID, resource.SpaceRoleDeveloper.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}
func (m *DefaultManager) RemoveSpaceManager(spaceName, spaceGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s with role %s", userName, spaceName, "Manager")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s with role %s", userName, spaceName, "Manager")
	roleGUID, err := m.GetSpaceRoleGUID(spaceGUID, userGUID, resource.SpaceRoleManager.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}
func (m *DefaultManager) RemoveSpaceSupporter(spaceName, spaceGUID, userName, userGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: removing user %s from org/space %s with role %s", userName, spaceName, "Supporter")
		return nil
	}
	lo.G.Infof("removing user %s from org/space %s with role %s", userName, spaceName, "Supporter")
	roleGUID, err := m.GetSpaceRoleGUID(spaceGUID, userGUID, resource.SpaceRoleSupporter.String())
	if err != nil {
		return err
	}
	return m.deleteRole(roleGUID)
}

func (m *DefaultManager) AssociateSpaceAuditor(orgGUID, spaceName, spaceGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s", userName, "auditor", spaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s", userName, "auditor", spaceName)
	_, err = m.RoleClient.CreateSpaceRole(context.Background(), spaceGUID, userGUID, resource.SpaceRoleAuditor)
	return err
}
func (m *DefaultManager) AssociateSpaceManager(orgGUID, spaceName, spaceGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s", userName, "manager", spaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s", userName, "manager", spaceName)
	_, err = m.RoleClient.CreateSpaceRole(context.Background(), spaceGUID, userGUID, resource.SpaceRoleManager)
	return err
}
func (m *DefaultManager) AssociateSpaceDeveloper(orgGUID, spaceName, spaceGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s", userName, "developer", spaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s", userName, "developer", spaceName)
	_, err = m.RoleClient.CreateSpaceRole(context.Background(), spaceGUID, userGUID, resource.SpaceRoleDeveloper)
	return err
}
func (m *DefaultManager) AssociateSpaceSupporter(orgGUID, spaceName, spaceGUID, userName, userGUID string) error {
	err := m.AddUserToOrg(orgGUID, userName, userGUID)
	if err != nil {
		return err
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s", userName, "supporter", spaceName)
		return nil
	}

	lo.G.Infof("adding %s to role %s for org/space %s", userName, "supporter", spaceName)
	_, err = m.RoleClient.CreateSpaceRole(context.Background(), spaceGUID, userGUID, resource.SpaceRoleSupporter)
	return err
}

func (m *DefaultManager) GetSpaceRoleGUID(spaceGUID, userGUID, role string) (string, error) {
	spaces, ok := m.SpaceRolesUsers[spaceGUID]
	if !ok {
		return "", fmt.Errorf("space with guid[%s] has no roles", spaceGUID)
	}
	roles, ok := spaces[role]
	if !ok {
		return "", fmt.Errorf("space with guid[%s] has no roles of type [%s]", spaceGUID, role)
	}
	roleGUID, ok := roles[userGUID]
	if !ok {
		return "", fmt.Errorf("space with guid[%s] has no role of type [%s] with user guid [%s]", spaceGUID, role, userGUID)
	}
	return roleGUID, nil
}

// func (m *DefaultManager) AssociateSpaceAuditor(input UsersInput, userName, userGUID string) error {
// 	err := m.RoleMgr.AddUserToOrg(input.OrgGUID, userName, userGUID)
// 	if err != nil {
// 		return err
// 	}
// 	if m.Peek {
// 		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "auditor", input.OrgName, input.SpaceName)
// 		return nil
// 	}

// 	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "auditor", input.OrgName, input.SpaceName)
// 	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_AUDITOR)
// 	return err
// }
// func (m *DefaultManager) AssociateSpaceDeveloper(input UsersInput, userName, userGUID string) error {
// 	err := m.RoleMgr.AddUserToOrg(input.OrgGUID, userName, userGUID)
// 	if err != nil {
// 		return err
// 	}
// 	if m.Peek {
// 		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "developer", input.OrgName, input.SpaceName)
// 		return nil
// 	}
// 	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "developer", input.OrgName, input.SpaceName)
// 	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_DEVELOPER)
// 	return err
// }
// func (m *DefaultManager) AssociateSpaceManager(input UsersInput, userName, userGUID string) error {
// 	err := m.RoleMgr.AddUserToOrg(input.OrgGUID, userName, userGUID)
// 	if err != nil {
// 		return err
// 	}
// 	if m.Peek {
// 		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "manager", input.OrgName, input.SpaceName)
// 		return nil
// 	}

// 	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "manager", input.OrgName, input.SpaceName)
// 	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_MANAGER)
// 	return err
// }

// func (m *DefaultManager) AssociateSpaceSupporter(input UsersInput, userName, userGUID string) error {

// 	if !m.SupportsSpaceSupporter {
// 		lo.G.Infof("this instance of cloud foundry does not support space_supporter role")
// 		return nil
// 	}
// 	err := m.RoleMgr.AddUserToOrg(input.OrgGUID, userName, userGUID)
// 	if err != nil {
// 		return err
// 	}
// 	if m.Peek {
// 		lo.G.Infof("[dry-run]: adding %s to role %s for org/space %s/%s", userName, "supporter", input.OrgName, input.SpaceName)
// 		return nil
// 	}

// 	lo.G.Infof("adding %s to role %s for org/space %s/%s", userName, "supporter", input.OrgName, input.SpaceName)
// 	_, err = m.Client.CreateV3SpaceRole(input.SpaceGUID, userGUID, SPACE_SUPPORTER)
// 	return err
// }
