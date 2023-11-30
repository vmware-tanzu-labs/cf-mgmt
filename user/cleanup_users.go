package user

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) removeOrphanedUsers(orphanedUsers []string) error {
	for _, orphanedUser := range orphanedUsers {
		lo.G.Infof("Deleting orphaned CF user with guid %s", orphanedUser)
		err := m.RoleMgr.DeleteUser(orphanedUser)
		if err != nil {
			return err
		}
	}

	return nil
}

// CleanupOrgUsers -
func (m *DefaultManager) CleanupOrgUsers() []error {
	errs := []error{}
	m.RoleMgr.ClearRoles()
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return []error{err}
	}
	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		return []error{err}
	}

	for _, input := range orgConfigs {
		if input.RemoveUsers {
			if err := m.cleanupOrgUsers(uaaUsers, &input); err != nil {
				errs = append(errs, err)
			}
		} else {
			lo.G.Infof("Not Removing Users from org %s", input.Org)
		}
	}
	return errs
}

func (m *DefaultManager) cleanupOrgUsers(uaaUsers *uaa.Users, input *config.OrgConfig) error {
	org, err := m.OrgReader.FindOrg(input.Org)
	if err != nil {
		return err
	}
	orgUsers, _, _, _, err := m.RoleMgr.ListOrgUsersByRole(org.GUID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error listing org users for org %s", input.Org))
	}

	usersInRoles, err := m.usersInOrgRoles(org.Name, org.GUID)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error usersInOrgRoles for org %s", input.Org))
	}

	lo.G.Debugf("Users In Roles %+v", usersInRoles)

	cfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}
	for _, orgUser := range orgUsers.Users() {
		uaaUser := uaaUsers.GetByID(orgUser.GUID)
		var guid string
		if uaaUser == nil {
			lo.G.Infof("Unable to find user (%s) GUID from uaa, using org user guid instead", orgUser.UserName)
			guid = orgUser.GUID
		} else {
			guid = uaaUser.GUID
		}
		if !util.Matches(orgUser.UserName, cfg.ProtectedUsers) {
			if !usersInRoles.HasUserForGUID(orgUser.UserName, guid) {
				if m.Peek {
					lo.G.Infof("[dry-run]: Removing User %s from org %s", orgUser.UserName, input.Org)
					continue
				}
				lo.G.Infof("Removing User %s from org %s", orgUser.UserName, input.Org)
				err = m.RoleMgr.RemoveOrgUser(org.Name, org.GUID, orgUser.UserName, guid)
				if err != nil {
					return err
				}
			}
		}
	}

	return m.removeOrphanedUsers(usersInRoles.OrphanedUsers())
}

func (m *DefaultManager) unassociatedOrphanedSpaceUser(input UsersInput, userGUIDs []string, unassign func(entityName string, entityGUID string, userName string, userGUID string) error) error {
	for _, userGUID := range userGUIDs {
		err := unassign(input.SpaceName, input.SpaceGUID, "orphaned", userGUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) unassociatedOrphanedOrgUser(input UsersInput, userGUIDs []string, unassign func(entityName string, entityGUID string, userName string, userGUID string) error) error {
	for _, userGUID := range userGUIDs {
		err := unassign(input.OrgName, input.OrgGUID, "orphaned", userGUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) usersInOrgRoles(orgName, orgGUID string) (*role.RoleUsers, error) {
	roleUsers := role.InitRoleUsers()

	userInput := UsersInput{
		OrgGUID: orgGUID,
		OrgName: orgName,
	}
	_, orgManagers, orgBillingManagers, orgAuditors, err := m.RoleMgr.ListOrgUsersByRole(orgGUID)
	if err != nil {
		return nil, err
	}
	roleUsers.AddUsers(orgAuditors.Users())
	roleUsers.AddOrphanedUsers(orgAuditors.OrphanedUsers())
	err = m.unassociatedOrphanedOrgUser(userInput, orgAuditors.OrphanedUsers(), m.RoleMgr.RemoveOrgAuditor)
	if err != nil {
		return nil, err
	}

	roleUsers.AddUsers(orgManagers.Users())
	roleUsers.AddOrphanedUsers(orgManagers.OrphanedUsers())

	err = m.unassociatedOrphanedOrgUser(userInput, orgManagers.OrphanedUsers(), m.RoleMgr.RemoveOrgManager)
	if err != nil {
		return nil, err
	}

	roleUsers.AddUsers(orgBillingManagers.Users())
	roleUsers.AddOrphanedUsers(orgBillingManagers.OrphanedUsers())

	err = m.unassociatedOrphanedOrgUser(userInput, orgBillingManagers.OrphanedUsers(), m.RoleMgr.RemoveOrgBillingManager)
	if err != nil {
		return nil, err
	}
	spaces, err := m.SpaceMgr.ListSpaces(orgGUID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing spaces for org %s", orgName))
	}
	for _, space := range spaces {
		userInput.SpaceGUID = space.GUID
		userInput.SpaceName = space.Name
		spaceManagers, spaceDevelopers, spaceAuditors, spaceSupporters, err := m.RoleMgr.ListSpaceUsersByRole(space.GUID)
		if err != nil {
			return nil, err
		}
		roleUsers.AddUsers(spaceAuditors.Users())
		roleUsers.AddOrphanedUsers(spaceAuditors.OrphanedUsers())
		err = m.unassociatedOrphanedSpaceUser(userInput, spaceAuditors.OrphanedUsers(), m.RoleMgr.RemoveSpaceAuditor)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceDevelopers.Users())
		roleUsers.AddOrphanedUsers(spaceDevelopers.OrphanedUsers())
		err = m.unassociatedOrphanedSpaceUser(userInput, spaceDevelopers.OrphanedUsers(), m.RoleMgr.RemoveSpaceDeveloper)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceManagers.Users())
		roleUsers.AddOrphanedUsers(spaceManagers.OrphanedUsers())

		err = m.unassociatedOrphanedSpaceUser(userInput, spaceManagers.OrphanedUsers(), m.RoleMgr.RemoveSpaceManager)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceSupporters.Users())
		roleUsers.AddOrphanedUsers(spaceSupporters.OrphanedUsers())

		err = m.unassociatedOrphanedSpaceUser(userInput, spaceSupporters.OrphanedUsers(), m.RoleMgr.RemoveSpaceSupporter)
		if err != nil {
			return nil, err
		}
	}

	return roleUsers, nil
}
