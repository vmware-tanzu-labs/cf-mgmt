package user

import (
	"fmt"
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/util"
	"github.com/pkg/errors"
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

//CleanupOrgUsers -
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
	org, err := m.OrgMgr.FindOrg(input.Org)
	if err != nil {
		return err
	}
	orgUsers, err := m.Client.ListOrgUsers(org.Guid)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error listing org users for org %s", input.Org))
	}

	usersInRoles, err := m.usersInOrgRoles(org.Name, org.Guid, uaaUsers)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error usersInOrgRoles for org %s", input.Org))
	}

	lo.G.Debugf("Users In Roles %+v", usersInRoles)

	cfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}
	for _, orgUser := range orgUsers {
		uaaUser := uaaUsers.GetByID(orgUser.Guid)
		var guid string
		if uaaUser == nil {
			lo.G.Infof("Unable to find user (%s) GUID from uaa, using org user guid instead", orgUser.Username)
			guid = orgUser.Guid
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
				err := m.Client.RemoveOrgUser(org.Guid, guid)
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

func (m *DefaultManager) usersInOrgRoles(orgName, orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error) {
	roleUsers := InitRoleUsers()

	userInput := UsersInput{
		OrgGUID: orgGUID,
		OrgName: orgName,
	}
	orgAuditors, err := m.ListOrgAuditors(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org auditors for org %s", orgName))
	}
	roleUsers.AddUsers(orgAuditors.Users())
	roleUsers.AddOrphanedUsers(orgAuditors.OrphanedUsers())
	err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveOrgAuditor)
	if err != nil {
		return nil, err
	}

	orgManagers, err := m.ListOrgManagers(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org managers for org %s", orgName))
	}
	roleUsers.AddUsers(orgManagers.Users())
	roleUsers.AddOrphanedUsers(orgManagers.OrphanedUsers())

	err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveOrgManager)
	if err != nil {
		return nil, err
	}

	orgBillingManagers, err := m.ListOrgBillingManagers(orgGUID, uaaUsers)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error listing org billing managers for org %s", orgName))
	}
	roleUsers.AddUsers(orgBillingManagers.Users())
	roleUsers.AddOrphanedUsers(orgBillingManagers.OrphanedUsers())

	err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveOrgBillingManager)
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
		spaceAuditors, err := m.ListSpaceAuditors(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space auditors for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceAuditors.Users())
		roleUsers.AddOrphanedUsers(spaceAuditors.OrphanedUsers())
		err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveSpaceAuditor)
		if err != nil {
			return nil, err
		}

		spaceDevelopers, err := m.ListSpaceDevelopers(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space developers for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceDevelopers.Users())
		roleUsers.AddOrphanedUsers(spaceDevelopers.OrphanedUsers())
		err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveSpaceDeveloper)
		if err != nil {
			return nil, err
		}

		spaceManagers, err := m.ListSpaceManagers(space.Guid, uaaUsers)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error listing space managers for org/space %s/%s", orgName, space.Name))
		}
		roleUsers.AddUsers(spaceManagers.Users())
		roleUsers.AddOrphanedUsers(spaceManagers.OrphanedUsers())

		err = m.unassociatedOrphanedUser(userInput, orgAuditors.OrphanedUsers(), m.RemoveSpaceManager)
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
