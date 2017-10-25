package export

import (
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

//NewExportManager Creates a new instance of the ImportConfig manager
func NewExportManager(
	configDir string,
	uaacMgr uaac.Manager,
	cloudController cc.Manager) Manager {
	return &DefaultImportManager{
		ConfigDir:       configDir,
		UAACMgr:         uaacMgr,
		CloudController: cloudController,
	}
}

//ExportConfig Imports org and space configuration from an existing CF instance
//Entries part of excludedOrgs and excludedSpaces are not included in the import
func (im *DefaultImportManager) ExportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error {
	//Get all the users from the foundation
	userIDToUserMap, err := im.UAACMgr.UsersByID()
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
	lo.G.Debugf("Orgs to process: %s", orgs)

	for _, org := range orgs {
		if _, ok := excludedOrgs[org.Entity.Name]; ok {
			lo.G.Infof("Skipping org: %s as it is ignored from import", org.Entity.Name)
			continue
		}
		lo.G.Infof("Processing org: %s ", org.Entity.Name)
		orgConfig := &config.OrgConfig{Org: org.Entity.Name}
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
			if _, ok := excludedSpaces[orgSpace.Entity.Name]; ok {
				lo.G.Infof("Skipping space: %s as it is ignored from import", orgSpace.Entity.Name)
				continue
			}
			lo.G.Infof("Processing space: %s", orgSpace.Entity.Name)

			spaceConfig := &config.SpaceConfig{Org: org.Entity.Name, Space: orgSpace.Entity.Name}
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
			configMgr.AddSpaceToConfig(spaceConfig)
		}
	}
	return nil
}

func quotaDefinition(controller cc.Manager, quotaDefinitionGUID, entityType string) cc.QuotaEntity {
	quotaDef, _ := controller.QuotaDef(quotaDefinitionGUID, entityType)
	if quotaDef.Entity.Name != "default" {
		return quotaDef.Entity
	}
	return cc.QuotaEntity{}
}

func addOrgUsers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	addOrgManagers(orgConfig, controller, userIDToUserMap, orgGUID)
	addBillingManagers(orgConfig, controller, userIDToUserMap, orgGUID)
	addOrgAuditors(orgConfig, controller, userIDToUserMap, orgGUID)
}

func addSpaceUsers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	addSpaceDevelopers(spaceConfig, controller, userIDToUserMap, spaceGUID)
	addSpaceManagers(spaceConfig, controller, userIDToUserMap, spaceGUID)
	addSpaceAuditors(spaceConfig, controller, userIDToUserMap, spaceGUID)
}

func addOrgManagers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	orgMgrs, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_MANAGERS)
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs), orgConfig.Org)
	doAddUsers(orgMgrs, &orgConfig.Manager.Users, &orgConfig.Manager.LDAPUsers, userIDToUserMap)
}

func addBillingManagers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	orgBillingMgrs, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_BILLING_MANAGERS)
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs), orgConfig.Org)
	doAddUsers(orgBillingMgrs, &orgConfig.BillingManager.Users, &orgConfig.BillingManager.LDAPUsers, userIDToUserMap)
}

func addOrgAuditors(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	orgAuditors, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_AUDITORS)
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors), orgConfig.Org)
	doAddUsers(orgAuditors, &orgConfig.Auditor.Users, &orgConfig.Auditor.LDAPUsers, userIDToUserMap)
}

func addSpaceManagers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceMgrs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_MANAGERS)
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceMgrs, &spaceConfig.Manager.Users, &spaceConfig.Manager.LDAPUsers, userIDToUserMap)
}

func addSpaceDevelopers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceDevs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_DEVELOPERS)
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceDevs, &spaceConfig.Developer.Users, &spaceConfig.Developer.LDAPUsers, userIDToUserMap)
}

func addSpaceAuditors(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceAuditors, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_AUDITORS)
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceAuditors, &spaceConfig.Auditor.Users, &spaceConfig.Auditor.LDAPUsers, userIDToUserMap)
}

func doAddUsers(cfUsers map[string]string, uaaUsers *[]string, ldapUsers *[]string, userIDToUserMap map[string]uaac.User) {
	for cfUser := range cfUsers {
		if usr, ok := userIDToUserMap[cfUser]; ok {
			if usr.Origin == "uaa" {
				*uaaUsers = append(*uaaUsers, usr.UserName)
			} else {
				*ldapUsers = append(*ldapUsers, usr.UserName)
			}
		} else {
			lo.G.Infof("CFUser [%s] not found in uaa user list", cfUser)
		}
	}
}

func getCFUsers(cc cc.Manager, entityGUID, entityType, role string) (map[string]string, error) {
	return cc.GetCFUsers(entityGUID, entityType, role)
}
