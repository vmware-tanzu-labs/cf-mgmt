package export

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/privatedomain"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/pivotalservices/cf-mgmt/user"
	"github.com/xchapter7x/lo"
)

//NewExportManager Creates a new instance of the ImportConfig manager
func NewExportManager(
	configDir string,
	uaaMgr uaa.Manager,
	spaceManager space.Manager,
	userManager user.Manager,
	orgManager organization.Manager,
	securityGroupManager securitygroup.Manager,
	isoSegmentMgr isosegment.Manager,
	privateDomainMgr privatedomain.Manager) Manager {
	return &DefaultImportManager{
		ConfigDir:            configDir,
		UAAMgr:               uaaMgr,
		SpaceManager:         spaceManager,
		UserManager:          userManager,
		OrgManager:           orgManager,
		SecurityGroupManager: securityGroupManager,
		IsoSegmentManager:    isoSegmentMgr,
		PrivateDomainManager: privateDomainMgr,
	}
}

//DefaultImportManager  -
type DefaultImportManager struct {
	ConfigDir            string
	UAAMgr               uaa.Manager
	SpaceManager         space.Manager
	UserManager          user.Manager
	OrgManager           organization.Manager
	SecurityGroupManager securitygroup.Manager
	IsoSegmentManager    isosegment.Manager
	PrivateDomainManager privatedomain.Manager
}

//ExportConfig Imports org and space configuration from an existing CF instance
//Entries part of excludedOrgs and excludedSpaces are not included in the import
func (im *DefaultImportManager) ExportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error {
	//Get all the users from the foundation
	userIDToUserMap, err := im.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error("Unable to retrieve users")
		return err
	}
	lo.G.Debugf("uaa user id map %v", userIDToUserMap)
	//Get all the orgs
	orgs, err := im.OrgManager.ListOrgs()
	if err != nil {
		lo.G.Errorf("Unable to retrieve orgs. Error : %s", err)
		return err
	}

	securityGroups, err := im.SecurityGroupManager.ListNonDefaultSecurityGroups()
	if err != nil {
		lo.G.Errorf("Unable to retrieve security groups. Error : %s", err)
		return err
	}

	defaultSecurityGroups, err := im.SecurityGroupManager.ListDefaultSecurityGroups()
	if err != nil {
		lo.G.Errorf("Unable to retrieve security groups. Error : %s", err)
		return err
	}

	isolationSegments, err := im.IsoSegmentManager.ListIsolationSegments()
	if err != nil {
		lo.G.Errorf("Unable to retrieve isolation segments. Error : %s", err)
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
		orgName := org.Name
		if _, ok := excludedOrgs[orgName]; ok {
			lo.G.Infof("Skipping org: %s as it is ignored from import", orgName)
			continue
		}

		lo.G.Infof("Processing org: %s ", orgName)
		orgConfig := &config.OrgConfig{Org: orgName}
		//Add users
		im.addOrgUsers(orgConfig, userIDToUserMap, org.Guid)
		//Add Quota definition if applicable
		if org.QuotaDefinitionGuid != "" {
			quota, err := org.Quota()
			if err != nil {
				return err
			}
			if quota != nil {
				if quota.Name == orgName {
					orgConfig.EnableOrgQuota = true
				}
				orgConfig.MemoryLimit = quota.MemoryLimit
				orgConfig.InstanceMemoryLimit = quota.InstanceMemoryLimit
				orgConfig.TotalRoutes = quota.TotalRoutes
				orgConfig.TotalServices = quota.TotalServices
				orgConfig.PaidServicePlansAllowed = quota.NonBasicServicesAllowed
				orgConfig.TotalPrivateDomains = quota.TotalPrivateDomains
				orgConfig.TotalReservedRoutePorts = quota.TotalReservedRoutePorts
				orgConfig.TotalServiceKeys = quota.TotalServiceKeys
				orgConfig.AppInstanceLimit = quota.AppInstanceLimit
			}
		}
		if org.DefaultIsolationSegmentGuid != "" {
			for _, isosegment := range isolationSegments {
				if isosegment.GUID == org.DefaultIsolationSegmentGuid {
					orgConfig.DefaultIsoSegment = isosegment.Name
				}
			}
		}

		privatedomains, err := im.PrivateDomainManager.ListOrgSharedPrivateDomains(org.Guid)
		if err != nil {
			return err
		}
		for privatedomain, _ := range privatedomains {
			orgConfig.SharedPrivateDomains = append(orgConfig.SharedPrivateDomains, privatedomain)
		}

		privatedomains, err = im.PrivateDomainManager.ListOrgOwnedPrivateDomains(org.Guid)
		if err != nil {
			return err
		}
		for privatedomain, _ := range privatedomains {
			orgConfig.PrivateDomains = append(orgConfig.PrivateDomains, privatedomain)
		}
		configMgr.AddOrgToConfig(orgConfig)

		lo.G.Infof("Done creating org %s", orgConfig.Org)
		lo.G.Infof("Listing spaces for org %s", orgConfig.Org)
		spaces, _ := im.SpaceManager.ListSpaces(org.Guid)
		lo.G.Infof("Found %d Spaces for org %s", len(spaces), orgConfig.Org)
		for _, orgSpace := range spaces {
			spaceName := orgSpace.Name
			if _, ok := excludedSpaces[spaceName]; ok {
				lo.G.Infof("Skipping space: %s as it is ignored from import", spaceName)
				continue
			}
			lo.G.Infof("Processing space: %s", spaceName)

			spaceConfig := &config.SpaceConfig{Org: org.Name, Space: spaceName}
			//Add users
			im.addSpaceUsers(spaceConfig, userIDToUserMap, orgSpace.Guid)
			//Add Quota definition if applicable
			if orgSpace.QuotaDefinitionGuid != "" {
				quota, err := orgSpace.Quota()
				if err != nil {
					return err
				}
				if quota != nil {
					if quota.Name == orgSpace.Name {
						spaceConfig.EnableSpaceQuota = true
					}
					spaceConfig.MemoryLimit = quota.MemoryLimit
					spaceConfig.InstanceMemoryLimit = quota.InstanceMemoryLimit
					spaceConfig.TotalRoutes = quota.TotalRoutes
					spaceConfig.TotalServices = quota.TotalServices
					spaceConfig.PaidServicePlansAllowed = quota.NonBasicServicesAllowed
					spaceConfig.TotalReservedRoutePorts = quota.TotalReservedRoutePorts
					spaceConfig.TotalServiceKeys = quota.TotalServiceKeys
					spaceConfig.AppInstanceLimit = quota.AppInstanceLimit
				}
			}

			if orgSpace.IsolationSegmentGuid != "" {
				for _, isosegment := range isolationSegments {
					if isosegment.GUID == orgSpace.IsolationSegmentGuid {
						spaceConfig.IsoSegment = isosegment.Name
					}
				}

			}
			if orgSpace.AllowSSH {
				spaceConfig.AllowSSH = true
			}

			spaceSGName := fmt.Sprintf("%s-%s", orgName, spaceName)
			if spaceSGNames, err := im.SecurityGroupManager.ListSpaceSecurityGroups(orgSpace.Guid); err == nil {
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
				if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.Guid); err == nil {
					configMgr.AddSecurityGroupToSpace(orgName, spaceName, rules)
				}
			}

		}
	}

	for sgName, sgInfo := range securityGroups {
		lo.G.Infof("Adding security group %s", sgName)
		if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.Guid); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			configMgr.AddSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	for sgName, sgInfo := range defaultSecurityGroups {
		lo.G.Infof("Adding default security group %s", sgName)
		if sgInfo.Running {
			globalConfig.RunningSecurityGroups = append(globalConfig.RunningSecurityGroups, sgName)
		}
		if sgInfo.Staging {
			globalConfig.StagingSecurityGroups = append(globalConfig.StagingSecurityGroups, sgName)
		}
		if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.Guid); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			configMgr.AddDefaultSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	return configMgr.SaveGlobalConfig(globalConfig)
}

/*func quotaDefinition(controller cc.Manager, quotaDefinitionGUID, entityType string) cc.QuotaEntity {
	quotaDef, _ := controller.Org(quotaDefinitionGUID, entityType)
	if quotaDef.Name != "default" {
		return quotaDef.Entity
	}
	return cc.QuotaEntity{}
}*/

func (im *DefaultImportManager) addOrgUsers(orgConfig *config.OrgConfig, userIDToUserMap map[string]uaa.User, orgGUID string) {
	im.addOrgManagers(orgConfig, userIDToUserMap, orgGUID)
	im.addBillingManagers(orgConfig, userIDToUserMap, orgGUID)
	im.addOrgAuditors(orgConfig, userIDToUserMap, orgGUID)
}

func (im *DefaultImportManager) addSpaceUsers(spaceConfig *config.SpaceConfig, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	im.addSpaceDevelopers(spaceConfig, userIDToUserMap, spaceGUID)
	im.addSpaceManagers(spaceConfig, userIDToUserMap, spaceGUID)
	im.addSpaceAuditors(spaceConfig, userIDToUserMap, spaceGUID)
}

func (im *DefaultImportManager) addOrgManagers(orgConfig *config.OrgConfig, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgMgrs, _ := im.UserManager.ListOrgManagers(orgGUID)
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs), orgConfig.Org)
	doAddUsers(orgMgrs, &orgConfig.Manager.Users, &orgConfig.Manager.LDAPUsers, &orgConfig.Manager.SamlUsers, userIDToUserMap)
}

func (im *DefaultImportManager) addBillingManagers(orgConfig *config.OrgConfig, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgBillingMgrs, _ := im.UserManager.ListOrgBillingManagers(orgGUID)
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs), orgConfig.Org)
	doAddUsers(orgBillingMgrs, &orgConfig.BillingManager.Users, &orgConfig.BillingManager.LDAPUsers, &orgConfig.BillingManager.SamlUsers, userIDToUserMap)
}

func (im *DefaultImportManager) addOrgAuditors(orgConfig *config.OrgConfig, userIDToUserMap map[string]uaa.User, orgGUID string) {
	orgAuditors, _ := im.UserManager.ListOrgAuditors(orgGUID)
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors), orgConfig.Org)
	doAddUsers(orgAuditors, &orgConfig.Auditor.Users, &orgConfig.Auditor.LDAPUsers, &orgConfig.Auditor.SamlUsers, userIDToUserMap)
}

func (im *DefaultImportManager) addSpaceManagers(spaceConfig *config.SpaceConfig, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceMgrs, _ := im.UserManager.ListSpaceManagers(spaceGUID)
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceMgrs, &spaceConfig.Manager.Users, &spaceConfig.Manager.LDAPUsers, &spaceConfig.Manager.SamlUsers, userIDToUserMap)
}

func (im *DefaultImportManager) addSpaceDevelopers(spaceConfig *config.SpaceConfig, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceDevs, _ := im.UserManager.ListSpaceDevelopers(spaceGUID)
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceDevs, &spaceConfig.Developer.Users, &spaceConfig.Developer.LDAPUsers, &spaceConfig.Developer.SamlUsers, userIDToUserMap)
}

func (im *DefaultImportManager) addSpaceAuditors(spaceConfig *config.SpaceConfig, userIDToUserMap map[string]uaa.User, spaceGUID string) {
	spaceAuditors, _ := im.UserManager.ListSpaceAuditors(spaceGUID)
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceAuditors, &spaceConfig.Auditor.Users, &spaceConfig.Auditor.LDAPUsers, &spaceConfig.Auditor.SamlUsers, userIDToUserMap)
}

func doAddUsers(cfUsers map[string]string, uaaUsers *[]string, ldapUsers *[]string, samlUsers *[]string, userIDToUserMap map[string]uaa.User) {
	for cfUser := range cfUsers {
		if usr, ok := userIDToUserMap[cfUser]; ok {
			if usr.Origin == "uaa" {
				*uaaUsers = append(*uaaUsers, usr.Username)
			} else if usr.Origin == "ldap" {
				*ldapUsers = append(*ldapUsers, usr.Username)
			} else {
				*samlUsers = append(*samlUsers, usr.Username)
			}
		} else {
			lo.G.Infof("CFUser [%s] not found in uaa user list", cfUser)
		}
	}
}
