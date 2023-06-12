package user

import (
	"fmt"
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func (m *DefaultManager) removeOrphanedUsers(orphanedUsers []string) error {
	for _, orphanedUser := range orphanedUsers {
		lo.G.Infof("Deleting orphaned CF user with guid %s", orphanedUser)
		err := m.Client.DeleteUser(orphanedUser)
		if err != nil {
			return err
		}
	}

	return nil
}

// CleanupOrgUsers -
func (m *DefaultManager) CleanupOrgUsers() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}
	uaaUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		return err
	}

	for _, input := range orgConfigs {
		if input.RemoveUsers {
			if err := m.cleanupOrgUsers(uaaUsers, &input); err != nil {
				return err
			}
		} else {
			lo.G.Infof("Not Removing Users from org %s", input.Org)
		}
	}
	return nil
}

func (m *DefaultManager) cleanupOrgUsers(uaaUsers *uaa.Users, input *config.OrgConfig) error {
	org, err := m.OrgReader.FindOrg(input.Org)
	if err != nil {
		return err
	}
	orgUsers, err := m.Client.ListV3OrganizationRolesByGUIDAndType(org.Guid, ORG_USER)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error listing org users for org %s", input.Org))
	}

	usersInRoles, err := m.usersInOrgRoles(org.Name, org.Guid)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error usersInOrgRoles for org %s", input.Org))
	}

	lo.G.Debugf("Users In Roles %+v", usersInRoles)

	cfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}
	for _, orgUser := range orgUsers {
		uaaUser := uaaUsers.GetByID(orgUser.GUID)
		var guid string
		if uaaUser == nil {
			lo.G.Infof("Unable to find user (%s) GUID from uaa, using org user guid instead", orgUser.Username)
			guid = orgUser.GUID
		} else {
			guid = uaaUser.GUID
		}
		if !util.Matches(orgUser.Username, cfg.ProtectedUsers) {
			if !usersInRoles.HasUser(orgUser.Username) {
				if m.Peek {
					lo.G.Infof("[dry-run]: Removing User %s from org %s", orgUser.Username, input.Org)
					continue
				}
				lo.G.Infof("Removing User %s from org %s", orgUser.Username, input.Org)
				role, err := m.GetOrgRoleGUID(org.Guid, guid, ORG_USER)
				if err != nil {
					return err
				}
				err = m.Client.DeleteV3Role(role)
				if err != nil {
					return errors.Wrap(err, fmt.Sprintf("Error removing user %s from org %s", orgUser.Username, input.Org))
				}
			}
		}
	}

	return m.removeOrphanedUsers(usersInRoles.OrphanedUsers())
}

func (m *DefaultManager) unassociatedOrphanedUser(input UsersInput, userGUIDs []string, unassign func(input UsersInput, userName string, userGUID string) error) error {
	for _, userGUID := range userGUIDs {
		err := unassign(input, "orphaned", userGUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) usersInOrgRoles(orgName, orgGUID string) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()

	userInput := UsersInput{
		OrgGUID: orgGUID,
		OrgName: orgName,
	}
	_, orgManagers, orgBillingManagers, orgAuditors, err := m.ListOrgUsersByRole(orgGUID)
	if err != nil {
		return nil, err
	}
	roleUsers.AddUsers(orgAuditors.Users())
	roleUsers.AddOrphanedUsers(orgAuditors.OrphanedUsers())
	err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveOrgAuditor)
	if err != nil {
		return nil, err
	}

	roleUsers.AddUsers(orgManagers.Users())
	roleUsers.AddOrphanedUsers(orgManagers.OrphanedUsers())

	err = m.unassociatedOrphanedUser(userInput, orgManagers.OrphanedUsers(), m.RemoveOrgManager)
	if err != nil {
		return nil, err
	}

	roleUsers.AddUsers(orgBillingManagers.Users())
	roleUsers.AddOrphanedUsers(orgBillingManagers.OrphanedUsers())

	err = m.unassociatedOrphanedUser(userInput, orgBillingManagers.OrphanedUsers(), m.RemoveOrgBillingManager)
	if err != nil {
		return nil, err
	}

	spaces, err := m.listSpaces(orgGUID)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing spaces for org %s", orgName))
	}
	for _, space := range spaces {
		userInput.SpaceGUID = space.Guid
		userInput.SpaceName = space.Name
		spaceManagers, spaceDevelopers, spaceAuditors, spaceSupporters, err := m.ListSpaceUsersByRole(space.Guid)
		if err != nil {
			return nil, err
		}
		roleUsers.AddUsers(spaceAuditors.Users())
		roleUsers.AddOrphanedUsers(spaceAuditors.OrphanedUsers())
		err = m.unassociatedOrphanedUser(userInput, spaceAuditors.OrphanedUsers(), m.RemoveSpaceAuditor)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceDevelopers.Users())
		roleUsers.AddOrphanedUsers(spaceDevelopers.OrphanedUsers())
		err = m.unassociatedOrphanedUser(userInput, spaceDevelopers.OrphanedUsers(), m.RemoveSpaceDeveloper)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceManagers.Users())
		roleUsers.AddOrphanedUsers(spaceManagers.OrphanedUsers())

		err = m.unassociatedOrphanedUser(userInput, spaceManagers.OrphanedUsers(), m.RemoveSpaceManager)
		if err != nil {
			return nil, err
		}

		roleUsers.AddUsers(spaceSupporters.Users())
		roleUsers.AddOrphanedUsers(spaceSupporters.OrphanedUsers())

		err = m.unassociatedOrphanedUser(userInput, spaceSupporters.OrphanedUsers(), m.RemoveSpaceSupporter)
		if err != nil {
			return nil, err
		}
	}

	return roleUsers, nil
}

func (m *DefaultManager) listSpaces(orgGUID string) ([]cfclient.Space, error) {
	spaces, err := m.Client.ListSpacesByQuery(url.Values{
		"q": []string{fmt.Sprintf("%s:%s", "organization_guid", orgGUID)},
	})
	if err != nil {
		return nil, err
	}
	return spaces, err

}
