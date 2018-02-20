package export

import (
	"fmt"

	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

//NewExportManager Creates a new instance of the ImportConfig manager
func NewExportManager(
	configDir string,
	uaaMgr uaa.Manager,
	cloudController cc.Manager) Manager {
	return &DefaultImportManager{
		ConfigDir:       configDir,
		UAAMgr:          uaaMgr,
		CloudController: cloudController,
	}
}

//ExportConfig Imports org and space configuration from an existing CF instance
//Entries part of excludedOrgs and excludedSpaces are not included in the import
func (im *DefaultImportManager) ExportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error {
	//Get all the users from the foundation
	userIDToUserMap, err := im.UAAMgr.UsersByID()
	if err != nil {
		lo.G.Error("Unable to retrieve users")
		return err
	}
	lo.G.Debugf("uaa user id map %v", userIDToUserMap)
	//Get all the orgs
	orgs, err := im.CloudController.ListOrgs()
	if err != nil {
		lo.G.Errorf("Unable to retrieve orgs. Error : %s", err)
		return err
	}

	securityGroups, err := im.CloudController.ListNonDefaultSecurityGroups()
	if err != nil {
		lo.G.Errorf("Unable to retrieve security groups. Error : %s", err)
		return err
	}

	defaultSecurityGroups, err := im.CloudController.ListDefaultSecurityGroups()
	if err != nil {
		lo.G.Errorf("Unable to retrieve security groups. Error : %s", err)
		return err
	}

	configMgr := config.NewManager(im.ConfigDir)
	lo.G.Info("Trying to delete existing config directory")
	//Delete existing config directory
	err = configMgr.DeleteConfigIfExists()
	if err != nil {
		return err
	}
	//Create a brand new directory
	lo.G.Info("Trying to create new config folder")

	var uaaUserOrigin string
	for _, usr := range userIDToUserMap {
		if usr.Origin != "" {
			uaaUserOrigin = usr.Origin
			break
		}
	}
	lo.G.Infof("Using UAA user origin: %s", uaaUserOrigin)
	err = configMgr.CreateConfigIfNotExists(uaaUserOrigin)
	if err != nil {
		return err
	}

	globalConfig, err := configMgr.GetGlobalConfig()
	if err != nil {
		return err
	}

	lo.G.Debugf("Orgs to process: %s", orgs)

	for _, org := range orgs {
		orgName := org.Entity.Name
		if _, ok := excludedOrgs[orgName]; ok {
			lo.G.Infof("Skipping org: %s as it is ignored from import", orgName)
			continue
		}

		lo.G.Infof("Processing org: %s ", orgName)
		orgConfig := &config.OrgConfig{Org: orgName}
		//Add users
		addOrgUsers(orgConfig, im.CloudController, userIDToUserMap, org.MetaData.GUID)
		//Add Quota definition if applicable
		if org.Entity.QuotaDefinitionGUID != "" {
			quota := quotaDefinition(im.CloudController, org.Entity.QuotaDefinitionGUID, organization.ORGS)
			orgConfig.EnableOrgQuota = quota.IsQuotaEnabled()
			orgConfig.MemoryLimit = quota.GetMemoryLimit()
			orgConfig.InstanceMemoryLimit = quota.GetInstanceMemoryLimit()
			orgConfig.TotalRoutes = quota.GetTotalRoutes()
			orgConfig.TotalServices = quota.GetTotalServices()
			orgConfig.PaidServicePlansAllowed = quota.IsPaidServicesAllowed()
			orgConfig.TotalPrivateDomains = quota.TotalPrivateDomains
			orgConfig.TotalReservedRoutePorts = quota.TotalReservedRoutePorts
			orgConfig.TotalServiceKeys = quota.TotalServiceKeys
			orgConfig.AppInstanceLimit = quota.AppInstanceLimit
		}
		configMgr.AddOrgToConfig(orgConfig)

		lo.G.Infof("Done creating org %s", orgConfig.Org)
		lo.G.Infof("Listing spaces for org %s", orgConfig.Org)
		spaces, _ := im.CloudController.ListSpaces(org.MetaData.GUID)
		lo.G.Infof("Found %d Spaces for org %s", len(spaces), orgConfig.Org)
		for _, orgSpace := range spaces {
			spaceName := orgSpace.Entity.Name
			if _, ok := excludedSpaces[spaceName]; ok {
				lo.G.Infof("Skipping space: %s as it is ignored from import", spaceName)
				continue
			}
			lo.G.Infof("Processing space: %s", spaceName)

			spaceConfig := &config.SpaceConfig{Org: org.Entity.Name, Space: spaceName}
			//Add users
			addSpaceUsers(spaceConfig, im.CloudController, userIDToUserMap, orgSpace.MetaData.GUID)
			//Add Quota definition if applicable
			if orgSpace.Entity.QuotaDefinitionGUID != "" {
				quota := quotaDefinition(im.CloudController, orgSpace.Entity.QuotaDefinitionGUID, space.SPACES)
				spaceConfig.EnableSpaceQuota = quota.IsQuotaEnabled()
				spaceConfig.MemoryLimit = quota.GetMemoryLimit()
				spaceConfig.InstanceMemoryLimit = quota.GetInstanceMemoryLimit()
				spaceConfig.TotalRoutes = quota.GetTotalRoutes()
				spaceConfig.TotalServices = quota.GetTotalServices()
				spaceConfig.PaidServicePlansAllowed = quota.IsPaidServicesAllowed()
				spaceConfig.TotalPrivateDomains = quota.TotalPrivateDomains
				spaceConfig.TotalReservedRoutePorts = quota.TotalReservedRoutePorts
				spaceConfig.TotalServiceKeys = quota.TotalServiceKeys
				spaceConfig.AppInstanceLimit = quota.AppInstanceLimit
			}
			if orgSpace.Entity.AllowSSH {
				spaceConfig.AllowSSH = true
			}

			spaceSGName := fmt.Sprintf("%s-%s", orgName, spaceName)
			if spaceSGNames, err := im.CloudController.ListSpaceSecurityGroups(orgSpace.MetaData.GUID); err == nil {
				for securityGroupName, _ := range spaceSGNames {
					lo.G.Infof("Adding named security group [%s] to space [%s]", securityGroupName, spaceName)
					if securityGroupName != spaceSGName {
						spaceConfig.ASGs = append(spaceConfig.ASGs, securityGroupName)
					}
				}
			}

			configMgr.AddSpaceToConfig(spaceConfig)

			if sgInfo, ok := securityGroups[spaceSGName]; ok {
				delete(securityGroups, spaceSGName)
				if rules, err := im.CloudController.GetSecurityGroupRules(sgInfo.GUID); err == nil {
					configMgr.AddSecurityGroupToSpace(orgName, spaceName, rules)
				}
			}

		}
	}

	for sgName, sgInfo := range securityGroups {
		lo.G.Infof("Adding security group %s", sgName)
		if rules, err := im.CloudController.GetSecurityGroupRules(sgInfo.GUID); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			configMgr.AddSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	for sgName, sgInfo := range defaultSecurityGroups {
		lo.G.Infof("Adding default security group %s", sgName)
		if sgInfo.DefaultRunning {
			globalConfig.RunningSecurityGroups = append(globalConfig.RunningSecurityGroups, sgName)
		}
		if sgInfo.DefaultStaging {
			globalConfig.StagingSecurityGroups = append(globalConfig.StagingSecurityGroups, sgName)
		}
		if rules, err := im.CloudController.GetSecurityGroupRules(sgInfo.GUID); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			configMgr.AddDefaultSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	return configMgr.SaveGlobalConfig(globalConfig)
}

func quotaDefinition(controller cc.Manager, quotaDefinitionGUID, entityType string) cc.QuotaEntity {
	quotaDef, _ := controller.QuotaDef(quotaDefinitionGUID, entityType)
	if quotaDef.Entity.Name != "default" {
		return quotaDef.Entity
	}
	return cc.QuotaEntity{}
}

func addOrgUsers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, orgGUID string) {
	addOrgManagers(orgConfig, controller, userIDToUserMap, orgGUID)
	addBillingManagers(orgConfig, controller, userIDToUserMap, orgGUID)
	addOrgAuditors(orgConfig, controller, userIDToUserMap, orgGUID)
}

func addSpaceUsers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	addSpaceDevelopers(spaceConfig, controller, userIDToUserMap, spaceGUID)
	addSpaceManagers(spaceConfig, controller, userIDToUserMap, spaceGUID)
	addSpaceAuditors(spaceConfig, controller, userIDToUserMap, spaceGUID)
}

func addOrgManagers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgMgrs, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_MANAGERS)
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs), orgConfig.Org)
	doAddUsers(orgMgrs, &orgConfig.Manager.Users, &orgConfig.Manager.LDAPUsers, &orgConfig.Manager.SamlUsers, userIDToUserMap)
}

func addBillingManagers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgBillingMgrs, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_BILLING_MANAGERS)
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs), orgConfig.Org)
	doAddUsers(orgBillingMgrs, &orgConfig.BillingManager.Users, &orgConfig.BillingManager.LDAPUsers, &orgConfig.BillingManager.SamlUsers, userIDToUserMap)
}

func addOrgAuditors(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgAuditors, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_AUDITORS)
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors), orgConfig.Org)
	doAddUsers(orgAuditors, &orgConfig.Auditor.Users, &orgConfig.Auditor.LDAPUsers, &orgConfig.Auditor.SamlUsers, userIDToUserMap)
}

func addSpaceManagers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceMgrs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_MANAGERS)
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceMgrs, &spaceConfig.Manager.Users, &spaceConfig.Manager.LDAPUsers, &spaceConfig.Manager.SamlUsers, userIDToUserMap)
}

func addSpaceDevelopers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceDevs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_DEVELOPERS)
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceDevs, &spaceConfig.Developer.Users, &spaceConfig.Developer.LDAPUsers, &spaceConfig.Developer.SamlUsers, userIDToUserMap)
}

func addSpaceAuditors(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceAuditors, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_AUDITORS)
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceAuditors, &spaceConfig.Auditor.Users, &spaceConfig.Auditor.LDAPUsers, &spaceConfig.Auditor.SamlUsers, userIDToUserMap)
}

func doAddUsers(cfUsers map[string]string, uaaUsers *[]string, ldapUsers *[]string, samlUsers *[]string, userIDToUserMap map[string]uaa.User) {
	for cfUser := range cfUsers {
		if usr, ok := userIDToUserMap[cfUser]; ok {
			if usr.Origin == "uaa" {
				*uaaUsers = append(*uaaUsers, usr.UserName)
			} else if usr.Origin == "ldap" {
				*ldapUsers = append(*ldapUsers, usr.UserName)
			} else {
				*samlUsers = append(*samlUsers, usr.UserName)
			}
		} else {
			lo.G.Infof("CFUser [%s] not found in uaa user list", cfUser)
		}
	}
}

func getCFUsers(cc cc.Manager, entityGUID, entityType, role string) (map[string]string, error) {
	return cc.GetCFUsers(entityGUID, entityType, role)
}
