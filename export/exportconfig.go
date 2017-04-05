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
	var err error
	var orgs []*cc.Org
	var configMgr config.Manager
	var userIDToUserMap map[string]uaac.User
	var spaces []*cc.Space

	//Get all the users from the foundation
	userIDToUserMap, err = im.UAACMgr.UsersByID()
	if err != nil {
		lo.G.Error("Unable to retrieve users")
		return err
	}
	//Get all the orgs
	orgs, err = im.CloudController.ListOrgs()
	if err != nil {
		lo.G.Errorf("Unable to retrieve orgs. Error : %s", err)
		return err
	}
	configMgr = config.NewManager(im.ConfigDir)
	lo.G.Info("Trying to delete existing config directory")
	//Delete existing config directory
	configMgr.DeleteConfigIfExists()
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
		if _, ok := excludedOrgs[org.Entity.Name]; !ok {
			lo.G.Infof("Processing org: %s ", org.Entity.Name)
			orgConfig := &config.OrgConfig{OrgName: org.Entity.Name}
			//Add users
			addOrgUsers(orgConfig, im.CloudController, userIDToUserMap, org.MetaData.GUID)
			//Add Quota definition if applicable
			if org.Entity.QuotaDefinitionGUID != "" {
				orgConfig.OrgQuota = quotaDefinition(im.CloudController, org.Entity.QuotaDefinitionGUID, organization.ORGS)
			}
			configMgr.AddOrgToConfig(orgConfig)

			spaces, _ = im.CloudController.ListSpaces(org.MetaData.GUID)

			for _, orgSpace := range spaces {
				if _, ok := excludedSpaces[orgSpace.Entity.Name]; !ok {
					lo.G.Infof("Processing space: %s", orgSpace.Entity.Name)

					spaceConfig := &config.SpaceConfig{OrgName: org.Entity.Name, SpaceName: orgSpace.Entity.Name}
					//Add users
					addSpaceUsers(spaceConfig, im.CloudController, userIDToUserMap, orgSpace.MetaData.GUID)
					//Add Quota definition if applicable
					if orgSpace.Entity.QuotaDefinitionGUID != "" {
						spaceConfig.SpaceQuota = quotaDefinition(im.CloudController, orgSpace.Entity.QuotaDefinitionGUID, space.SPACES)
					}
					if orgSpace.Entity.AllowSSH {
						spaceConfig.AllowSSH = true
					}
					configMgr.AddSpaceToConfig(spaceConfig)
				} else {
					lo.G.Infof("Skipping space: %s as it is ignored from import", orgSpace.Entity.Name)
				}
			}
		} else {
			lo.G.Infof("Skipping org: %s as it is ignored from import", org.Entity.Name)
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
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs), orgConfig.OrgName)
	doAddUsers(orgMgrs, &orgConfig.OrgMgrUAAUsers, &orgConfig.OrgMgrLDAPUsers, userIDToUserMap)
}

func addBillingManagers(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	orgBillingMgrs, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_BILLING_MANAGERS)
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs), orgConfig.OrgName)
	doAddUsers(orgBillingMgrs, &orgConfig.OrgBillingMgrUAAUsers, &orgConfig.OrgBillingMgrLDAPUsers, userIDToUserMap)
}

func addOrgAuditors(orgConfig *config.OrgConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, orgGUID string) {
	orgAuditors, _ := getCFUsers(controller, orgGUID, organization.ORGS, organization.ROLE_ORG_AUDITORS)
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors), orgConfig.OrgName)
	doAddUsers(orgAuditors, &orgConfig.OrgAuditorUAAUsers, &orgConfig.OrgAuditorLDAPUsers, userIDToUserMap)
}

func addSpaceManagers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceMgrs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_MANAGERS)
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs), spaceConfig.OrgName, spaceConfig.SpaceName)
	doAddUsers(spaceMgrs, &spaceConfig.SpaceMgrUAAUsers, &spaceConfig.SpaceMgrLDAPUsers, userIDToUserMap)
}

func addSpaceDevelopers(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceDevs, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_DEVELOPERS)
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs), spaceConfig.OrgName, spaceConfig.SpaceName)
	doAddUsers(spaceDevs, &spaceConfig.SpaceDevUAAUsers, &spaceConfig.SpaceDevLDAPUsers, userIDToUserMap)
}

func addSpaceAuditors(spaceConfig *config.SpaceConfig, controller cc.Manager, userIDToUserMap map[string]uaac.User, spaceGUID string) {
	spaceAuditors, _ := getCFUsers(controller, spaceGUID, space.SPACES, space.ROLE_SPACE_MANAGERS)
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors), spaceConfig.OrgName, spaceConfig.SpaceName)
	doAddUsers(spaceAuditors, &spaceConfig.SpaceAuditorUAAUsers, &spaceConfig.SpaceAuditorLDAPUsers, userIDToUserMap)
}

func doAddUsers(cfUsers map[string]string, uaaUsers *[]string, ldapUsers *[]string, userIDToUserMap map[string]uaac.User) {
	for cfUser := range cfUsers {
		if usr, ok := userIDToUserMap[cfUser]; ok {
			if usr.Origin == "uaa" {
				*uaaUsers = append(*uaaUsers, usr.UserName)
			} else {
				*ldapUsers = append(*ldapUsers, usr.UserName)
			}
		}
	}
}

func getCFUsers(cc cc.Manager, entityGUID, entityType, role string) (map[string]string, error) {
	return cc.GetCFUsers(entityGUID, entityType, role)
}
