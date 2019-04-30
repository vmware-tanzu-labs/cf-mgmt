package export

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/isosegment"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/privatedomain"
	"github.com/pivotalservices/cf-mgmt/quota"
	"github.com/pivotalservices/cf-mgmt/securitygroup"
	"github.com/pivotalservices/cf-mgmt/serviceaccess"
	"github.com/pivotalservices/cf-mgmt/shareddomain"
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
	privateDomainMgr privatedomain.Manager,
	sharedDomainMgr *shareddomain.Manager,
	serviceAccessMgr *serviceaccess.Manager,
	quotaMgr *quota.Manager) Manager {
	return &DefaultImportManager{
		ConfigDir:            configDir,
		UAAMgr:               uaaMgr,
		SpaceManager:         spaceManager,
		UserManager:          userManager,
		OrgManager:           orgManager,
		SecurityGroupManager: securityGroupManager,
		IsoSegmentManager:    isoSegmentMgr,
		PrivateDomainManager: privateDomainMgr,
		SharedDomainManager:  sharedDomainMgr,
		ServiceAccessManager: serviceAccessMgr,
		QuotaManager:         quotaMgr,
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
	SharedDomainManager  *shareddomain.Manager
	ServiceAccessManager *serviceaccess.Manager
	QuotaManager         *quota.Manager
}

//ExportConfig Imports org and space configuration from an existing CF instance
//Entries part of excludedOrgs and excludedSpaces are not included in the import
func (im *DefaultImportManager) ExportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error {
	//Get all the users from the foundation
	uaaUsers, err := im.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error("Unable to retrieve users")
		return err
	}
	lo.G.Debugf("uaa user id map %v", uaaUsers)
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
	for _, usr := range uaaUsers.List() {
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

	globalConfig.EnableServiceAccess = true
	for _, org := range orgs {
		orgName := org.Name
		if _, ok := excludedOrgs[orgName]; ok {
			lo.G.Infof("Skipping org: %s as it is ignored from export", orgName)
			continue
		}

		lo.G.Infof("Processing org: %s ", orgName)
		orgConfig := &config.OrgConfig{Org: orgName}
		//Add users
		im.addOrgUsers(orgConfig, uaaUsers, org.Guid)
		//Add Quota definition if applicable
		if org.QuotaDefinitionGuid != "" {
			orgQuota, err := org.Quota()
			if err != nil {
				return err
			}
			if orgQuota != nil {
				orgConfig.EnableOrgQuota = false
				orgConfig.NamedQuota = orgQuota.Name
				// orgConfig.MemoryLimit = config.ByteSize(orgQuota.MemoryLimit)
				// orgConfig.InstanceMemoryLimit = config.ByteSize(orgQuota.InstanceMemoryLimit)
				// orgConfig.TotalRoutes = config.AsString(orgQuota.TotalRoutes)
				// orgConfig.TotalServices = config.AsString(orgQuota.TotalServices)
				// orgConfig.PaidServicePlansAllowed = orgQuota.NonBasicServicesAllowed
				// orgConfig.TotalPrivateDomains = config.AsString(orgQuota.TotalPrivateDomains)
				// orgConfig.TotalReservedRoutePorts = config.AsString(orgQuota.TotalReservedRoutePorts)
				// orgConfig.TotalServiceKeys = config.AsString(orgQuota.TotalServiceKeys)
				// orgConfig.AppInstanceLimit = config.AsString(orgQuota.AppInstanceLimit)
				// orgConfig.AppTaskLimit = config.AsString(orgQuota.AppTaskLimit)

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

		orgConfig.ServiceAccess = make(map[string][]string)
		serviceInfo, err := im.ServiceAccessManager.ListServiceInfo()
		if err != nil {
			return err
		}
		for service, plans := range serviceInfo.AllPlans() {
			accessPlans := []string{}
			for _, plan := range plans {
				if plan.Public || plan.OrgHasAccess(org.Guid) {
					accessPlans = append(accessPlans, plan.Name)
				}
			}
			orgConfig.ServiceAccess[service] = accessPlans
		}

		spacesConfig := &config.Spaces{Org: orgConfig.Org, EnableDeleteSpaces: true}
		configMgr.AddOrgToConfig(orgConfig, spacesConfig)

		lo.G.Infof("Done creating org %s", orgConfig.Org)
		lo.G.Infof("Listing spaces for org %s", orgConfig.Org)
		spaces, _ := im.SpaceManager.ListSpaces(org.Guid)
		lo.G.Infof("Found %d Spaces for org %s", len(spaces), orgConfig.Org)

		spaceQuotas, err := im.QuotaManager.ListAllSpaceQuotasForOrg(org.Guid)
		if err != nil {
			return err
		}

		for _, spaceQuota := range spaceQuotas {
			if !im.doesSpaceExist(spaces, spaceQuota.Name) {
				err = configMgr.AddSpaceQuota(config.SpaceQuota{
					Org:                     org.Name,
					Name:                    spaceQuota.Name,
					AppInstanceLimit:        config.AsString(spaceQuota.AppInstanceLimit),
					TotalReservedRoutePorts: config.AsString(spaceQuota.TotalReservedRoutePorts),
					TotalServiceKeys:        config.AsString(spaceQuota.TotalServiceKeys),
					AppTaskLimit:            config.AsString(spaceQuota.AppTaskLimit),
					MemoryLimit:             config.ByteSize(spaceQuota.MemoryLimit),
					InstanceMemoryLimit:     config.ByteSize(spaceQuota.InstanceMemoryLimit),
					TotalRoutes:             config.AsString(spaceQuota.TotalRoutes),
					TotalServices:           config.AsString(spaceQuota.TotalServices),
					PaidServicePlansAllowed: spaceQuota.NonBasicServicesAllowed,
				})
				if err != nil {
					return err
				}
			}
		}

		for _, orgSpace := range spaces {
			spaceName := orgSpace.Name
			if _, ok := excludedSpaces[spaceName]; ok {
				lo.G.Infof("Skipping space: %s as it is ignored from export", spaceName)
				continue
			}
			lo.G.Infof("Processing space: %s", spaceName)

			spaceConfig := &config.SpaceConfig{Org: org.Name, Space: spaceName, EnableUnassignSecurityGroup: true}
			//Add users
			im.addSpaceUsers(spaceConfig, uaaUsers, orgSpace.Guid)
			//Add Quota definition if applicable
			if orgSpace.QuotaDefinitionGuid != "" {
				quota, err := orgSpace.Quota()
				if err != nil {
					return err
				}
				if quota != nil {
					if quota.Name == orgSpace.Name {
						spaceConfig.EnableSpaceQuota = true
						spaceConfig.MemoryLimit = config.ByteSize(quota.MemoryLimit)
						spaceConfig.InstanceMemoryLimit = config.ByteSize(quota.InstanceMemoryLimit)
						spaceConfig.TotalRoutes = config.AsString(quota.TotalRoutes)
						spaceConfig.TotalServices = config.AsString(quota.TotalServices)
						spaceConfig.PaidServicePlansAllowed = quota.NonBasicServicesAllowed
						spaceConfig.TotalReservedRoutePorts = config.AsString(quota.TotalReservedRoutePorts)
						spaceConfig.TotalServiceKeys = config.AsString(quota.TotalServiceKeys)
						spaceConfig.AppInstanceLimit = config.AsString(quota.AppInstanceLimit)
						spaceConfig.AppTaskLimit = config.AsString(quota.AppTaskLimit)
					} else {
						spaceConfig.NamedQuota = quota.Name
					}
				}
			} else {
				spaceConfig.MemoryLimit = orgConfig.MemoryLimit
				spaceConfig.InstanceMemoryLimit = orgConfig.InstanceMemoryLimit
				spaceConfig.TotalRoutes = orgConfig.TotalRoutes
				spaceConfig.TotalServices = orgConfig.TotalServices
				spaceConfig.PaidServicePlansAllowed = orgConfig.PaidServicePlansAllowed
				spaceConfig.TotalReservedRoutePorts = orgConfig.TotalReservedRoutePorts
				spaceConfig.TotalServiceKeys = orgConfig.TotalServiceKeys
				spaceConfig.AppInstanceLimit = orgConfig.AppInstanceLimit
				spaceConfig.AppTaskLimit = orgConfig.AppTaskLimit
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

	orgQuotas, err := im.QuotaManager.Client.ListOrgQuotas()
	if err != nil {
		return err
	}

	for _, orgQuota := range orgQuotas {

		err = configMgr.AddOrgQuota(config.OrgQuota{
			Name:                    orgQuota.Name,
			AppInstanceLimit:        config.AsString(orgQuota.AppInstanceLimit),
			TotalPrivateDomains:     config.AsString(orgQuota.TotalPrivateDomains),
			TotalReservedRoutePorts: config.AsString(orgQuota.TotalReservedRoutePorts),
			TotalServiceKeys:        config.AsString(orgQuota.TotalServiceKeys),
			AppTaskLimit:            config.AsString(orgQuota.AppTaskLimit),
			MemoryLimit:             config.ByteSize(orgQuota.MemoryLimit),
			InstanceMemoryLimit:     config.ByteSize(orgQuota.InstanceMemoryLimit),
			TotalRoutes:             config.AsString(orgQuota.TotalRoutes),
			TotalServices:           config.AsString(orgQuota.TotalServices),
			PaidServicePlansAllowed: orgQuota.NonBasicServicesAllowed,
		})
		if err != nil {
			return err
		}
	}

	sharedDomains, err := im.SharedDomainManager.CFClient.ListSharedDomains()
	if err != nil {
		return err
	}

	routerGroups, err := im.SharedDomainManager.RoutingClient.RouterGroups()
	if err != nil {
		return err
	}
	globalConfig.EnableDeleteSharedDomains = true
	globalConfig.SharedDomains = make(map[string]config.SharedDomain)
	for _, sharedDomain := range sharedDomains {
		sharedDomainConfig := config.SharedDomain{
			Internal: sharedDomain.Internal,
		}
		if sharedDomain.RouterGroupGuid != "" {
			for _, routerGroup := range routerGroups {
				if routerGroup.Guid == sharedDomain.RouterGroupGuid {
					sharedDomainConfig.RouterGroup = routerGroup.Name
					continue
				}
			}
		}
		globalConfig.SharedDomains[sharedDomain.Name] = sharedDomainConfig
	}
	return configMgr.SaveGlobalConfig(globalConfig)
}

func (im *DefaultImportManager) doesOrgExist(orgs []cfclient.Org, orgName string) bool {
	for _, org := range orgs {
		if org.Name == orgName {
			return true
		}
	}
	return false
}

func (im *DefaultImportManager) doesSpaceExist(spaces []cfclient.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Name == spaceName {
			return true
		}
	}
	return false
}

func (im *DefaultImportManager) addOrgUsers(orgConfig *config.OrgConfig, uaaUsers *uaa.Users, orgGUID string) {
	im.addOrgManagers(orgConfig, uaaUsers, orgGUID)
	im.addBillingManagers(orgConfig, uaaUsers, orgGUID)
	im.addOrgAuditors(orgConfig, uaaUsers, orgGUID)
}

func (im *DefaultImportManager) addSpaceUsers(spaceConfig *config.SpaceConfig, uaaUsers *uaa.Users, spaceGUID string) {
	im.addSpaceDevelopers(spaceConfig, uaaUsers, spaceGUID)
	im.addSpaceManagers(spaceConfig, uaaUsers, spaceGUID)
	im.addSpaceAuditors(spaceConfig, uaaUsers, spaceGUID)
}

func (im *DefaultImportManager) addOrgManagers(orgConfig *config.OrgConfig, uaaUsers *uaa.Users, orgGUID string) {
	orgMgrs, _ := im.UserManager.ListOrgManagers(orgGUID, uaaUsers)
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs.Users()), orgConfig.Org)
	doAddUsers(orgMgrs, &orgConfig.Manager.Users, &orgConfig.Manager.LDAPUsers, &orgConfig.Manager.SamlUsers)
}

func (im *DefaultImportManager) addBillingManagers(orgConfig *config.OrgConfig, uaaUsers *uaa.Users, orgGUID string) {
	orgBillingMgrs, _ := im.UserManager.ListOrgBillingManagers(orgGUID, uaaUsers)
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs.Users()), orgConfig.Org)
	doAddUsers(orgBillingMgrs, &orgConfig.BillingManager.Users, &orgConfig.BillingManager.LDAPUsers, &orgConfig.BillingManager.SamlUsers)
}

func (im *DefaultImportManager) addOrgAuditors(orgConfig *config.OrgConfig, uaaUsers *uaa.Users, orgGUID string) {
	orgAuditors, _ := im.UserManager.ListOrgAuditors(orgGUID, uaaUsers)
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors.Users()), orgConfig.Org)
	doAddUsers(orgAuditors, &orgConfig.Auditor.Users, &orgConfig.Auditor.LDAPUsers, &orgConfig.Auditor.SamlUsers)
}

func (im *DefaultImportManager) addSpaceManagers(spaceConfig *config.SpaceConfig, uaaUsers *uaa.Users, spaceGUID string) {
	spaceMgrs, _ := im.UserManager.ListSpaceManagers(spaceGUID, uaaUsers)
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceMgrs, &spaceConfig.Manager.Users, &spaceConfig.Manager.LDAPUsers, &spaceConfig.Manager.SamlUsers)
}

func (im *DefaultImportManager) addSpaceDevelopers(spaceConfig *config.SpaceConfig, uaaUsers *uaa.Users, spaceGUID string) {
	spaceDevs, _ := im.UserManager.ListSpaceDevelopers(spaceGUID, uaaUsers)
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceDevs, &spaceConfig.Developer.Users, &spaceConfig.Developer.LDAPUsers, &spaceConfig.Developer.SamlUsers)
}

func (im *DefaultImportManager) addSpaceAuditors(spaceConfig *config.SpaceConfig, uaaUsers *uaa.Users, spaceGUID string) {
	spaceAuditors, _ := im.UserManager.ListSpaceAuditors(spaceGUID, uaaUsers)
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceAuditors, &spaceConfig.Auditor.Users, &spaceConfig.Auditor.LDAPUsers, &spaceConfig.Auditor.SamlUsers)
}

func doAddUsers(roleUser *user.RoleUsers, uaaUsers *[]string, ldapUsers *[]string, samlUsers *[]string) {
	for _, cfUser := range roleUser.Users() {
		if cfUser.Origin == "uaa" {
			*uaaUsers = append(*uaaUsers, cfUser.UserName)
		} else if cfUser.Origin == "ldap" {
			*ldapUsers = append(*ldapUsers, cfUser.UserName)
		} else {
			*samlUsers = append(*samlUsers, cfUser.UserName)
		}
	}
}
