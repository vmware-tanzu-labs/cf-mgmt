package export

import (
	"fmt"

	"code.cloudfoundry.org/routing-api/models"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/isosegment"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/privatedomain"
	"github.com/vmwarepivotallabs/cf-mgmt/quota"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/securitygroup"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"
	"github.com/vmwarepivotallabs/cf-mgmt/shareddomain"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	"github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

// NewExportManager Creates a new instance of the ImportConfig manager
func NewExportManager(
	configDir string,
	uaaMgr uaa.Manager,
	spaceManager space.Manager,
	userManager user.Manager,
	orgReader organizationreader.Reader,
	securityGroupManager securitygroup.Manager,
	isoSegmentMgr isosegment.Manager,
	privateDomainMgr privatedomain.Manager,
	sharedDomainMgr *shareddomain.Manager,
	serviceAccessMgr *serviceaccess.Manager,
	quotaMgr *quota.Manager,
	roleMgr role.Manager) *Manager {
	return &Manager{
		ConfigMgr:            config.NewManager(configDir),
		UAAMgr:               uaaMgr,
		SpaceManager:         spaceManager,
		UserManager:          userManager,
		OrgReader:            orgReader,
		SecurityGroupManager: securityGroupManager,
		IsoSegmentManager:    isoSegmentMgr,
		PrivateDomainManager: privateDomainMgr,
		SharedDomainManager:  sharedDomainMgr,
		ServiceAccessManager: serviceAccessMgr,
		QuotaManager:         quotaMgr,
		RoleManager:          roleMgr,
	}
}

type Manager struct {
	ConfigMgr            config.Manager
	UAAMgr               uaa.Manager
	SpaceManager         space.Manager
	UserManager          user.Manager
	RoleManager          role.Manager
	OrgReader            organizationreader.Reader
	SecurityGroupManager securitygroup.Manager
	IsoSegmentManager    isosegment.Manager
	PrivateDomainManager privatedomain.Manager
	SharedDomainManager  *shareddomain.Manager
	ServiceAccessManager *serviceaccess.Manager
	QuotaManager         *quota.Manager
	SkipSpaces           bool
	SkipRoutingGroups    bool
}

func (im *Manager) ExportServiceAccess() error {

	orgConfigs, err := im.ConfigMgr.GetOrgConfigs()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		if orgConfig.ServiceAccess != nil {
			orgConfig.ServiceAccess = nil
			err = im.ConfigMgr.SaveOrgConfig(&orgConfig)
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Updated orgConfig.yml for org [%s] to remove service-access configuration", orgConfig.Org))
		}
	}

	globalConfig, err := im.ConfigMgr.GetGlobalConfig()
	if err != nil {
		return err
	}
	orgs, err := im.OrgReader.ListOrgs()
	if err != nil {
		lo.G.Errorf("Unable to retrieve orgs. Error : %s", err)
		return err
	}
	err = im.exportServiceAccess(globalConfig, orgs)
	if err != nil {
		return err
	}
	globalConfig.IgnoreLegacyServiceAccess = true
	err = im.ConfigMgr.SaveGlobalConfig(globalConfig)

	if err == nil {
		fmt.Println("Updated cf-mgmt.yml with service-access configuration")
	}
	return err
}

// ExportConfig Imports org and space configuration from an existing CF instance
// Entries part of excludedOrgs and excludedSpaces are not included in the import
func (im *Manager) ExportConfig(excludedOrgs, excludedSpaces map[string]string, skipSpaces bool) error {
	//Get all the users from the foundation
	uaaUsers, err := im.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error("Unable to retrieve users")
		return err
	}
	lo.G.Debugf("uaa user id map %v", uaaUsers)
	//Get all the orgs
	orgs, err := im.OrgReader.ListOrgs()
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
	lo.G.Info("Trying to delete existing config directory")
	//Delete existing config directory
	err = im.ConfigMgr.DeleteConfigIfExists()
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
	err = im.ConfigMgr.CreateConfigIfNotExists(uaaUserOrigin)
	if err != nil {
		return err
	}

	globalConfig, err := im.ConfigMgr.GetGlobalConfig()
	if err != nil {
		return err
	}

	lo.G.Debugf("Orgs to process: %s", orgs)

	globalConfig.EnableServiceAccess = true

	err = im.exportServiceAccess(globalConfig, orgs)
	if err != nil {
		return err
	}

	for _, org := range orgs {
		orgName := org.Name
		if _, ok := excludedOrgs[orgName]; ok {
			lo.G.Infof("Skipping org: %s as it is ignored from export", orgName)
			continue
		}

		lo.G.Infof("Processing org: %s ", orgName)
		orgConfig := &config.OrgConfig{Org: orgName}
		//Add users
		err = im.addOrgUsers(orgConfig, org.GUID)
		if err != nil {
			return err
		}
		//Add Quota definition if applicable
		if org.Relationships.Quota.Data != nil {
			orgQuota, err := im.QuotaManager.GetOrgQuota(org.Relationships.Quota.Data.GUID)
			if err != nil {
				return err
			}
			if orgQuota != nil {
				orgConfig.EnableOrgQuota = false
				orgConfig.NamedQuota = orgQuota.Name
			}
		}
		orgIsolationSegmentGUID, err := im.OrgReader.GetDefaultIsolationSegment(org)
		if err != nil {
			return err
		}
		if orgIsolationSegmentGUID != "" {
			for _, isosegment := range isolationSegments {
				if isosegment.GUID == orgIsolationSegmentGUID {
					orgConfig.DefaultIsoSegment = isosegment.Name
				}
			}
		}

		privatedomains, err := im.PrivateDomainManager.ListOrgSharedPrivateDomains(org.GUID)
		if err != nil {
			return err
		}
		for privatedomain := range privatedomains {
			orgConfig.SharedPrivateDomains = append(orgConfig.SharedPrivateDomains, privatedomain)
		}

		privatedomains, err = im.PrivateDomainManager.ListOrgOwnedPrivateDomains(org.GUID)
		if err != nil {
			return err
		}
		for privatedomain := range privatedomains {
			orgConfig.PrivateDomains = append(orgConfig.PrivateDomains, privatedomain)
		}

		err = im.ConfigMgr.SaveOrgConfig(orgConfig)
		if err != nil {
			return err
		}
		err = im.ConfigMgr.SaveOrgSpaces(&config.Spaces{Org: orgConfig.Org, EnableDeleteSpaces: !skipSpaces})
		if err != nil {
			return err
		}
		lo.G.Infof("Done creating org %s", orgConfig.Org)
		if !skipSpaces {
			err := im.processSpaces(orgConfig, org.GUID, excludedSpaces, isolationSegments, securityGroups)
			if err != nil {
				return errors.Wrapf(err, "Processing org %s", orgConfig.Org)
			}
		}
	}

	for sgName, sgInfo := range securityGroups {
		lo.G.Infof("Adding security group %s", sgName)
		if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.GUID); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			im.ConfigMgr.AddSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	for sgName, sgInfo := range defaultSecurityGroups {
		lo.G.Infof("Adding default security group %s", sgName)
		if sgInfo.GloballyEnabled.Running {
			globalConfig.RunningSecurityGroups = append(globalConfig.RunningSecurityGroups, sgName)
		}
		if sgInfo.GloballyEnabled.Staging {
			globalConfig.StagingSecurityGroups = append(globalConfig.StagingSecurityGroups, sgName)
		}
		if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.GUID); err == nil {
			lo.G.Infof("Adding rules for %s", sgName)
			im.ConfigMgr.AddDefaultSecurityGroup(sgName, rules)
		} else {
			lo.G.Error(err)
		}
	}

	orgQuotas, err := im.QuotaManager.ListAllOrgQuotas()
	if err != nil {
		return err
	}

	for _, orgQuota := range orgQuotas {

		err = im.ConfigMgr.AddOrgQuota(config.OrgQuota{
			Name:                       orgQuota.Name,
			AppInstanceLimit:           config.AsString(orgQuota.Apps.TotalInstances),
			TotalPrivateDomains:        config.AsString(orgQuota.Domains.TotalDomains),
			TotalReservedRoutePorts:    config.AsString(orgQuota.Routes.TotalReservedPorts),
			TotalServiceKeys:           config.AsString(orgQuota.Services.TotalServiceKeys),
			AppTaskLimit:               config.AsString(orgQuota.Apps.PerAppTasks),
			MemoryLimit:                config.ByteSize(orgQuota.Apps.TotalMemoryInMB),
			InstanceMemoryLimit:        config.ByteSize(orgQuota.Apps.PerProcessMemoryInMB),
			TotalRoutes:                config.AsString(orgQuota.Routes.TotalRoutes),
			TotalServices:              config.AsString(orgQuota.Services.TotalServiceInstances),
			PaidServicePlansAllowed:    *orgQuota.Services.PaidServicesAllowed,
			LogRateLimitBytesPerSecond: config.AsString(orgQuota.Apps.LogRateLimitInBytesPerSecond),
		})
		if err != nil {
			return err
		}
	}

	lo.G.Infof("Listing Shared Domains")
	sharedDomains, err := im.SharedDomainManager.CFClient.ListSharedDomains()
	if err != nil {
		return errors.Wrapf(err, "Getting shared domains")
	}

	var routerGroups []models.RouterGroup
	if !im.SkipRoutingGroups {
		lo.G.Infof("Listing Router Groups")
		routerGroups, err = im.SharedDomainManager.RoutingClient.RouterGroups()
		if err != nil {
			return errors.Wrapf(err, "Getting routing groups")
		}
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
	return im.ConfigMgr.SaveGlobalConfig(globalConfig)
}

func (im *Manager) processSpaces(orgConfig *config.OrgConfig, orgGUID string, excludedSpaces map[string]string, isolationSegments []cfclient.IsolationSegment, securityGroups map[string]*resource.SecurityGroup) error {
	lo.G.Infof("Listing spaces for org %s", orgConfig.Org)
	spaces, _ := im.SpaceManager.ListSpaces(orgGUID)
	lo.G.Infof("Found %d Spaces for org %s", len(spaces), orgConfig.Org)

	spaceQuotas, err := im.QuotaManager.ListAllSpaceQuotasForOrg(orgGUID)
	if err != nil {
		return err
	}

	for _, spaceQuota := range spaceQuotas {
		if !im.doesSpaceExist(spaces, spaceQuota.Name) {
			err = im.ConfigMgr.AddSpaceQuota(config.SpaceQuota{
				Org:                        orgConfig.Org,
				Name:                       spaceQuota.Name,
				AppInstanceLimit:           config.AsString(spaceQuota.Apps.TotalInstances),
				TotalReservedRoutePorts:    config.AsString(spaceQuota.Routes.TotalReservedPorts),
				TotalServiceKeys:           config.AsString(spaceQuota.Services.TotalServiceKeys),
				AppTaskLimit:               config.AsString(spaceQuota.Apps.PerAppTasks),
				MemoryLimit:                config.ByteSize(spaceQuota.Apps.TotalMemoryInMB),
				InstanceMemoryLimit:        config.ByteSize(spaceQuota.Apps.PerProcessMemoryInMB),
				TotalRoutes:                config.AsString(spaceQuota.Routes.TotalRoutes),
				TotalServices:              config.AsString(spaceQuota.Services.TotalServiceInstances),
				PaidServicePlansAllowed:    *spaceQuota.Services.PaidServicesAllowed,
				LogRateLimitBytesPerSecond: config.AsString(spaceQuota.Apps.LogRateLimitInBytesPerSecond),
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

		spaceConfig := &config.SpaceConfig{Org: orgConfig.Org, Space: spaceName, EnableUnassignSecurityGroup: true}
		//Add users
		err = im.addSpaceUsers(spaceConfig, orgSpace.Relationships.Organization.Data.GUID)
		if err != nil {
			return err
		}
		//Add Quota definition if applicable
		if orgSpace.Relationships.Quota.Data != nil {
			quota, err := im.QuotaManager.GetSpaceQuota(orgSpace.Relationships.Quota.Data.GUID)
			if err != nil {
				return err
			}
			if quota != nil {
				if quota.Name == orgSpace.Name {
					spaceConfig.EnableSpaceQuota = true
					spaceConfig.MemoryLimit = config.ByteSize(quota.Apps.TotalMemoryInMB)
					spaceConfig.InstanceMemoryLimit = config.ByteSize(quota.Apps.PerProcessMemoryInMB)
					spaceConfig.TotalRoutes = config.AsString(quota.Routes.TotalRoutes)
					spaceConfig.TotalServices = config.AsString(quota.Services.TotalServiceInstances)
					spaceConfig.PaidServicePlansAllowed = *quota.Services.PaidServicesAllowed
					spaceConfig.TotalReservedRoutePorts = config.AsString(quota.Routes.TotalReservedPorts)
					spaceConfig.TotalServiceKeys = config.AsString(quota.Services.TotalServiceKeys)
					spaceConfig.AppInstanceLimit = config.AsString(quota.Apps.TotalInstances)
					spaceConfig.AppTaskLimit = config.AsString(quota.Apps.PerAppTasks)
					spaceConfig.LogRateLimitBytesPerSecond = config.AsString(quota.Apps.LogRateLimitInBytesPerSecond)
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
			spaceConfig.LogRateLimitBytesPerSecond = orgConfig.LogRateLimitBytesPerSecond
		}
		isoSegGUID, err := im.SpaceManager.GetSpaceIsolationSegmentGUID(orgSpace)
		if err != nil {
			return err
		}
		if isoSegGUID != "" {
			for _, isosegment := range isolationSegments {
				if isosegment.GUID == isoSegGUID {
					spaceConfig.IsoSegment = isosegment.Name
				}
			}

		}
		sshEnabled, err := im.SpaceManager.IsSSHEnabled(orgSpace)
		if err != nil {
			return err
		}
		if sshEnabled {
			spaceConfig.AllowSSH = true
		}

		spaceSGName := fmt.Sprintf("%s-%s", orgConfig.Org, spaceName)
		if spaceSGNames, err := im.SecurityGroupManager.ListSpaceSecurityGroups(orgSpace.GUID); err == nil {
			for securityGroupName := range spaceSGNames {
				lo.G.Infof("Adding named security group [%s] to space [%s]", securityGroupName, spaceName)
				if securityGroupName != spaceSGName {
					spaceConfig.ASGs = append(spaceConfig.ASGs, securityGroupName)
				} else {
					spaceConfig.EnableSecurityGroup = true
				}
			}
		}

		im.ConfigMgr.AddSpaceToConfig(spaceConfig)

		if sgInfo, ok := securityGroups[spaceSGName]; ok {
			delete(securityGroups, spaceSGName)
			if rules, err := im.SecurityGroupManager.GetSecurityGroupRules(sgInfo.GUID); err == nil {
				err = im.ConfigMgr.AddSecurityGroupToSpace(orgConfig.Org, spaceName, rules)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (im *Manager) exportServiceAccess(globalConfig *config.GlobalConfig, orgs []*resource.Organization) error {
	globalConfig.ServiceAccess = nil
	serviceInfo, err := im.ServiceAccessManager.ListServiceInfo()
	if err != nil {
		return err
	}
	for _, broker := range serviceInfo.StandardBrokers() {
		brokerConfig := &config.Broker{
			Name: broker.Name,
		}
		for _, service := range broker.Services() {
			serviceVisibility := &config.Service{
				Name: service.Name,
			}
			for _, plan := range service.Plans() {
				if plan.Public {
					serviceVisibility.AllAccessPlans = append(serviceVisibility.AllAccessPlans, plan.Name)
					continue
				}
				if len(plan.ListVisibilities()) == 0 {
					serviceVisibility.NoAccessPlans = append(serviceVisibility.NoAccessPlans, plan.Name)
					continue
				}

				privatePlan := &config.PlanVisibility{
					Name: plan.Name,
				}
				for _, orgAccess := range plan.ListVisibilities() {
					orgName, err := im.getOrgName(orgs, orgAccess.OrgGUID)
					if err != nil {
						return err
					}

					if !util.Matches(orgName, config.DefaultProtectedOrgs) {
						privatePlan.Orgs = append(privatePlan.Orgs, orgName)
					}
				}
				serviceVisibility.LimitedAccessPlans = append(serviceVisibility.LimitedAccessPlans, privatePlan)
			}
			brokerConfig.Services = append(brokerConfig.Services, serviceVisibility)
		}
		globalConfig.ServiceAccess = append(globalConfig.ServiceAccess, brokerConfig)
	}
	return nil
}

func (im *Manager) getOrgName(orgs []*resource.Organization, orgGUID string) (string, error) {
	for _, org := range orgs {
		if org.GUID == orgGUID {
			return org.Name, nil
		}
	}
	return "", fmt.Errorf("no org exists for org guid %s", orgGUID)
}

func (im *Manager) doesSpaceExist(spaces []*resource.Space, spaceName string) bool {
	for _, space := range spaces {
		if space.Name == spaceName {
			return true
		}
	}
	return false
}

func (im *Manager) addOrgUsers(orgConfig *config.OrgConfig, orgGUID string) error {
	_, managerRoleUsers, billingManagerRoleUsers, auditorRoleUsers, err := im.RoleManager.ListOrgUsersByRole(orgGUID)
	if err != nil {
		return err
	}
	im.addOrgManagers(orgConfig, orgGUID, managerRoleUsers)
	im.addBillingManagers(orgConfig, orgGUID, billingManagerRoleUsers)
	im.addOrgAuditors(orgConfig, orgGUID, auditorRoleUsers)
	return nil
}

func (im *Manager) addSpaceUsers(spaceConfig *config.SpaceConfig, spaceGUID string) error {
	managerRoleUsers, developerRoleUsers, auditorRoleUsers, supporterRoleUsers, err := im.RoleManager.ListSpaceUsersByRole(spaceGUID)
	if err != nil {
		return err
	}
	im.addSpaceDevelopers(spaceConfig, spaceGUID, developerRoleUsers)
	im.addSpaceManagers(spaceConfig, spaceGUID, managerRoleUsers)
	im.addSpaceAuditors(spaceConfig, spaceGUID, auditorRoleUsers)
	im.addSpaceSupporters(spaceConfig, spaceGUID, supporterRoleUsers)
	return nil
}

func (im *Manager) addOrgManagers(orgConfig *config.OrgConfig, orgGUID string, orgMgrs *role.RoleUsers) {
	lo.G.Debugf("Found %d Org Managers for Org: %s", len(orgMgrs.Users()), orgConfig.Org)
	doAddUsers(orgMgrs, &orgConfig.Manager.Users, &orgConfig.Manager.LDAPUsers, &orgConfig.Manager.SamlUsers)
}

func (im *Manager) addBillingManagers(orgConfig *config.OrgConfig, orgGUID string, orgBillingMgrs *role.RoleUsers) {
	lo.G.Debugf("Found %d Org Billing Managers for Org: %s", len(orgBillingMgrs.Users()), orgConfig.Org)
	doAddUsers(orgBillingMgrs, &orgConfig.BillingManager.Users, &orgConfig.BillingManager.LDAPUsers, &orgConfig.BillingManager.SamlUsers)
}

func (im *Manager) addOrgAuditors(orgConfig *config.OrgConfig, orgGUID string, orgAuditors *role.RoleUsers) {
	lo.G.Debugf("Found %d Org Auditors for Org: %s", len(orgAuditors.Users()), orgConfig.Org)
	doAddUsers(orgAuditors, &orgConfig.Auditor.Users, &orgConfig.Auditor.LDAPUsers, &orgConfig.Auditor.SamlUsers)
}

func (im *Manager) addSpaceManagers(spaceConfig *config.SpaceConfig, spaceGUID string, spaceMgrs *role.RoleUsers) {
	lo.G.Debugf("Found %d Space Managers for Org: %s and  Space:  %s", len(spaceMgrs.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceMgrs, &spaceConfig.Manager.Users, &spaceConfig.Manager.LDAPUsers, &spaceConfig.Manager.SamlUsers)
}

func (im *Manager) addSpaceDevelopers(spaceConfig *config.SpaceConfig, spaceGUID string, spaceDevs *role.RoleUsers) {
	lo.G.Debugf("Found %d Space Developers for Org: %s and  Space:  %s", len(spaceDevs.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceDevs, &spaceConfig.Developer.Users, &spaceConfig.Developer.LDAPUsers, &spaceConfig.Developer.SamlUsers)
}

func (im *Manager) addSpaceAuditors(spaceConfig *config.SpaceConfig, spaceGUID string, spaceAuditors *role.RoleUsers) {
	lo.G.Debugf("Found %d Space Auditors for Org: %s and  Space:  %s", len(spaceAuditors.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceAuditors, &spaceConfig.Auditor.Users, &spaceConfig.Auditor.LDAPUsers, &spaceConfig.Auditor.SamlUsers)
}
func (im *Manager) addSpaceSupporters(spaceConfig *config.SpaceConfig, spaceGUID string, spaceSupporters *role.RoleUsers) {
	lo.G.Debugf("Found %d Space Supporters for Org: %s and  Space:  %s", len(spaceSupporters.Users()), spaceConfig.Org, spaceConfig.Space)
	doAddUsers(spaceSupporters, &spaceConfig.Supporter.Users, &spaceConfig.Supporter.LDAPUsers, &spaceConfig.Supporter.SamlUsers)
}

func doAddUsers(roleUser *role.RoleUsers, uaaUsers *[]string, ldapUsers *[]string, samlUsers *[]string) {
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
