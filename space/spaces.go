package space

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {
	return &DefaultSpaceManager{
		UAACMgr:         uaac.NewManager(sysDomain, uaacToken),
		CloudController: cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token),
		OrgMgr:          organization.NewManager(sysDomain, token, uaacToken),
		LdapMgr:         ldap.NewManager(),
		UtilsMgr:        utils.NewDefaultManager(),
	}
}

//CreateApplicationSecurityGroups -
func (m *DefaultSpaceManager) CreateApplicationSecurityGroups(configDir string) (err error) {
	var contents string
	var targetSGGUID string
	var sgs map[string]string
	var space cloudcontroller.Space
	files, _ := m.UtilsMgr.FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		input := &InputUpdateSpaces{}
		if err = m.UtilsMgr.LoadFile(f, input); err == nil {
			if input.EnableSecurityGroup {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					securityGroupFile := strings.Replace(f, "spaceConfig.yml", "security-group.json", -1)
					if contents, err = m.getSecurityFileContents(securityGroupFile); err == nil {
						sgName := fmt.Sprintf("%s-%s", input.Org, input.Space)
						if sgs, err = m.CloudController.ListSecurityGroups(); err == nil {
							if sgGUID, ok := sgs[sgName]; ok {
								lo.G.Info("Updating security group", sgName)
								if err = m.CloudController.UpdateSecurityGroup(sgGUID, sgName, contents); err == nil {
									lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
									m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, sgGUID)
								}
							} else {
								lo.G.Info("Creating security group", sgName)
								if targetSGGUID, err = m.CloudController.CreateSecurityGroup(sgName, contents); err == nil {
									lo.G.Info("Binding security group", sgName, "to space", space.Entity.Name)
									m.CloudController.AssignSecurityGroupToSpace(space.MetaData.GUID, targetSGGUID)
								}
							}
						}
					}
				}
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) getSecurityFileContents(securityGroupFile string) (contents string, err error) {
	var bytes []byte
	if bytes, err = ioutil.ReadFile(securityGroupFile); err == nil {
		contents = string(bytes)
	}
	return
}

//CreateQuotas -
func (m *DefaultSpaceManager) CreateQuotas(configDir string) (err error) {
	var quotas map[string]string
	var space cloudcontroller.Space
	var targetQuotaGUID string

	files, _ := m.UtilsMgr.FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		input := &InputUpdateSpaces{}
		if err = m.UtilsMgr.LoadFile(f, input); err == nil {
			lo.G.Info("Processing file", f)
			if input.EnableSpaceQuota {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					quotaName := space.Entity.Name
					if quotas, err = m.CloudController.ListSpaceQuotas(space.Entity.OrgGUID); err == nil {
						if quotaGUID, ok := quotas[quotaName]; ok {
							lo.G.Info("Updating quota", quotaName)
							if err = m.CloudController.UpdateSpaceQuota(space.Entity.OrgGUID, quotaGUID,
								quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								err = m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, quotaGUID)
							}
						} else {
							lo.G.Info("Creating quota", quotaName)
							if targetQuotaGUID, err = m.CloudController.CreateSpaceQuota(space.Entity.OrgGUID,
								quotaName, input.MemoryLimit, input.InstanceMemoryLimit, input.TotalRoutes, input.TotalServices, input.PaidServicePlansAllowed); err == nil {
								lo.G.Info("Assigning", quotaName, "to", space.Entity.Name)
								err = m.CloudController.AssignQuotaToSpace(space.MetaData.GUID, targetQuotaGUID)
							}
						}
					}
				}
			}
		}

	}

	if err != nil {
		lo.G.Error(err)
	}
	return
}

//UpdateSpaces -
func (m *DefaultSpaceManager) UpdateSpaces(configDir string) (err error) {
	var space cloudcontroller.Space
	files, _ := m.UtilsMgr.FindFiles(configDir, "spaceConfig.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputUpdateSpaces{}
		if err = m.UtilsMgr.LoadFile(f, input); err == nil {
			//if input, err = m.loadUpdateFile(f); err == nil {
			if space, err = m.FindSpace(input.Org, input.Space); err == nil {
				lo.G.Info("Processing space", space.Entity.Name)
				if input.AllowSSH != space.Entity.AllowSSH {
					if err = m.CloudController.UpdateSpaceSSH(input.AllowSSH, space.MetaData.GUID); err != nil {
						return
					}
				}
			}
		}
	}
	return
}

//UpdateSpaceUsers -
func (m *DefaultSpaceManager) UpdateSpaceUsers(configDir, ldapBindPassword string) (err error) {
	var space cloudcontroller.Space
	var config *ldap.Config

	if config, err = m.LdapMgr.GetConfig(configDir, ldapBindPassword); err != nil {
		return
	}

	if config.Enabled {
		files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaceConfig.yml")
		for _, f := range files {
			lo.G.Info("Processing space file", f)
			input := &InputUpdateSpaces{}
			if err = m.UtilsMgr.LoadFile(f, input); err == nil {
				if space, err = m.FindSpace(input.Org, input.Space); err == nil {
					lo.G.Info("User sync for space", space.Entity.Name)
					if err = m.updateUsers(config, space, "developers", input.DeveloperGroup); err != nil {
						return
					}
					if err = m.updateUsers(config, space, "managers", input.ManagerGroup); err != nil {
						return
					}
					if err = m.updateUsers(config, space, "auditors", input.AuditorGroup); err != nil {
						return
					}
				}
			}
		}
	}

	return
}
func (m *DefaultSpaceManager) updateUsers(config *ldap.Config, space cloudcontroller.Space, role, groupName string) (err error) {
	var groupUsers []ldap.User
	var uaacUsers map[string]string
	if groupName != "" {
		lo.G.Info("Getting users for group", groupName)
		if groupUsers, err = m.LdapMgr.GetUserIDs(config, groupName); err == nil {
			if uaacUsers, err = m.UAACMgr.ListUsers(); err == nil {
				for _, groupUser := range groupUsers {
					if _, userExists := uaacUsers[strings.ToLower(groupUser.UserID)]; userExists {
						lo.G.Info("User", groupUser.UserID, "already exists")
					} else {
						lo.G.Info("User", groupUser.UserID, "doesn't exist so creating in UAA")
						if err = m.UAACMgr.CreateLdapUser(groupUser.UserID, groupUser.Email, groupUser.UserDN); err != nil {
							return
						}
					}
					lo.G.Info("Adding user to groups")
					if err = m.addRole(groupUser.UserID, role, space.Entity.OrgGUID, space.MetaData.GUID); err != nil {
						lo.G.Error(err)
						return
					}
				}
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) addRole(userName, role, orgGUID, spaceGUID string) error {
	if err := m.CloudController.AddUserToOrg(userName, orgGUID); err != nil {
		return err
	}
	return m.CloudController.AddUserToSpaceRole(userName, role, spaceGUID)
}

//CreateSpace -
func (m *DefaultSpaceManager) CreateSpace(orgName, spaceName string) (space cloudcontroller.Space, err error) {
	var orgGUID string
	if orgGUID, err = m.OrgMgr.GetOrgGUID(orgName); err == nil {
		if err = m.CloudController.CreateSpace(spaceName, orgGUID); err == nil {
			space, err = m.FindSpace(orgGUID, spaceName)
		}
	}
	return
}

//FindSpace -
func (m *DefaultSpaceManager) FindSpace(orgName, spaceName string) (space cloudcontroller.Space, err error) {
	var spaces []cloudcontroller.Space
	if spaces, err = m.fetchSpaces(orgName); err == nil {
		for _, theSpace := range spaces {
			if theSpace.Entity.Name == spaceName {
				space = theSpace
				return
			}
		}
	}
	return
}

//CreateSpaces -
func (m *DefaultSpaceManager) CreateSpaces(configDir string) (err error) {
	files, _ := utils.NewDefaultManager().FindFiles(configDir, "spaces.yml")
	for _, f := range files {
		lo.G.Info("Processing space file", f)
		input := &InputCreateSpaces{}
		if err = utils.NewDefaultManager().LoadFile(f, input); err == nil {
			if len(input.Spaces) == 0 {
				lo.G.Info("No spaces in config file", f)
			}
			var spaces []cloudcontroller.Space
			if spaces, err = m.fetchSpaces(input.Org); err == nil {
				for _, spaceName := range input.Spaces {
					if m.doesSpaceExist(spaces, spaceName) {
						lo.G.Info(fmt.Sprintf("[%s] space already exists", spaceName))
					} else {
						lo.G.Info(fmt.Sprintf("Creating [%s] space in [%s] org", spaceName, input.Org))
						m.CreateSpace(input.Org, spaceName)
					}
				}
			} else {
				return
			}
		}
	}
	return
}

func (m *DefaultSpaceManager) doesSpaceExist(spaces []cloudcontroller.Space, spaceName string) (result bool) {
	result = false
	for _, space := range spaces {
		if space.Entity.Name == spaceName {
			result = true
			return
		}
	}
	return

}

func (m *DefaultSpaceManager) fetchSpaces(orgName string) (spaces []cloudcontroller.Space, err error) {
	var orgGUID string
	if orgGUID, err = m.OrgMgr.GetOrgGUID(orgName); err == nil {
		spaces, err = m.CloudController.ListSpaces(orgGUID)
	}
	return
}
