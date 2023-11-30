package serviceaccess

import (
	"strings"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/legacy"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient,
	orgReader organizationreader.Reader,
	cfg config.Reader, peek bool) *Manager {
	return &Manager{
		Client:    client,
		OrgReader: orgReader,
		Cfg:       cfg,
		Peek:      peek,
		LegacyMgr: legacy.NewManager(client, orgReader, cfg, peek),
	}
}

type Manager struct {
	Client    CFClient
	Cfg       config.Reader
	OrgReader organizationreader.Reader
	Peek      bool
	LegacyMgr *legacy.Manager
}

func (m *Manager) Apply() error {
	globalCfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	if globalCfg.EnableServiceAccess {
		orgConfigs, err := m.Cfg.GetOrgConfigs()
		if err != nil {
			return err
		}
		orgList := []string{}
		for _, orgConfig := range orgConfigs {
			if len(orgConfig.ServiceAccess) > 0 {
				orgList = append(orgList, orgConfig.Org)
			}
		}

		if len(orgList) > 0 && !globalCfg.IgnoreLegacyServiceAccess {
			lo.G.Warning("**** Deprecated **** - run `cf-mgmt export-service-access-config` and check in configuration changes as services-access for orgs [%s] is no longer supported in orgConfig.yml", strings.Join(orgList, ","))
			return m.LegacyMgr.Apply()
		}
	}
	serviceInfo, err := m.ListServiceInfo()
	if err != nil {
		return err
	}
	protectedOrgs, err := m.ProtectedOrgList()
	if err != nil {
		return err
	}
	return m.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
}

func (m *Manager) UpdateServiceAccess(globalCfg *config.GlobalConfig, serviceInfo *ServiceInfo, protectedOrgs []string) error {
	if !globalCfg.EnableServiceAccess {
		lo.G.Info("Service Access is not enabled.  Set enable-service-access: true in cf-mgmt.yml")
		return nil
	}

	for _, broker := range serviceInfo.StandardBrokers() {
		for _, service := range broker.Services() {
			for _, plan := range service.Plans() {
				planInfo := globalCfg.GetPlanInfo(broker.Name, service.Name, plan.Name)
				if planInfo.NoAccess {
					err := m.EnsureNoAccess(plan)
					if err != nil {
						return err
					}
					continue
				}
				if planInfo.Limited {
					err := m.EnsureLimitedAccess(plan, planInfo.Orgs, protectedOrgs)
					if err != nil {
						return err
					}
					continue
				}
				err := m.EnsurePublicAccess(plan)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) EnsurePublicAccess(plan *ServicePlanInfo) error {
	if !plan.Public {
		err := m.MakePublic(plan)
		if err != nil {
			return err
		}
	}
	return m.RemoveVisibilities(plan)
}

func (m *Manager) EnsureLimitedAccess(plan *ServicePlanInfo, orgs, protectedOrgs []string) error {
	lo.G.Debugf("Current Visibilities %+v for plan %s", plan.ListVisibilities(), plan.ServiceName)
	err := m.MakePrivate(plan)
	if err != nil {
		return err
	}
	for _, orgName := range m.uniqueOrgs(orgs, protectedOrgs) {
		err = m.CreatePlanVisibility(plan, orgName)
		if err != nil {
			return err
		}
	}
	return m.RemoveVisibilities(plan)
}

func (m *Manager) uniqueOrgs(orgs, protectedOrgs []string) []string {
	orgMap := make(map[string]string)
	allOrgs := append(orgs, protectedOrgs...)
	for _, org := range allOrgs {
		orgLower := strings.ToLower(org)
		_, ok := orgMap[orgLower]
		if !ok {
			orgMap[orgLower] = orgLower
		} else {
			lo.G.Debugf("Duplicate org %s in %+v", orgLower, allOrgs)
		}
	}
	orgList := []string{}
	for _, org := range orgMap {
		orgList = append(orgList, org)
	}

	return orgList
}

func (m *Manager) EnsureNoAccess(plan *ServicePlanInfo) error {
	if plan.Public {
		err := m.MakePrivate(plan)
		if err != nil {
			return err
		}
	}
	return m.RemoveVisibilities(plan)
}

func (m *Manager) ProtectedOrgList() ([]string, error) {
	orgConfig, err := m.Cfg.Orgs()
	if err != nil {
		return nil, err
	}
	orgs, err := m.OrgReader.ListOrgs()
	if err != nil {
		return nil, err
	}
	orgList := []string{}
	for _, org := range orgs {
		if util.Matches(org.Name, orgConfig.ProtectedOrgList()) {
			orgList = append(orgList, org.Name)
		}
	}
	return orgList, nil
}

func (m *Manager) CreatePlanVisibility(servicePlan *ServicePlanInfo, orgName string) error {
	org, err := m.OrgReader.FindOrg(orgName)
	if err != nil {
		return err
	}
	if !servicePlan.OrgHasAccess(org.GUID) {
		if m.Peek {
			lo.G.Infof("[dry-run]: adding plan %s for service %s to org %s", servicePlan.Name, servicePlan.ServiceName, orgName)
			return nil
		}
		lo.G.Infof("adding plan %s for service %s to org %s", servicePlan.Name, servicePlan.ServiceName, orgName)
		_, err = m.Client.CreateServicePlanVisibility(servicePlan.GUID, org.GUID)
		if err != nil {
			return err
		}
	} else {
		lo.G.Debugf("plan %s for service %s already visible to org %s", servicePlan.Name, servicePlan.ServiceName, orgName)
		servicePlan.RemoveOrg(org.GUID)
	}
	return nil
}

func (m *Manager) MakePublic(servicePlan *ServicePlanInfo) error {
	if !servicePlan.Public {
		if m.Peek {
			lo.G.Infof("[dry-run]: Making plan %s for service %s public", servicePlan.Name, servicePlan.ServiceName)
			return nil
		}
		lo.G.Infof("Making plan %s for service %s public", servicePlan.Name, servicePlan.ServiceName)
		err := m.Client.MakeServicePlanPublic(servicePlan.GUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) MakePrivate(servicePlan *ServicePlanInfo) error {
	if servicePlan.Public {
		if m.Peek {
			lo.G.Infof("[dry-run]: Making plan %s for service %s private", servicePlan.Name, servicePlan.ServiceName)
			return nil
		}
		lo.G.Infof("Making plan %s for service %s private", servicePlan.Name, servicePlan.ServiceName)
		err := m.Client.MakeServicePlanPrivate(servicePlan.GUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) RemoveVisibilities(servicePlan *ServicePlanInfo) error {
	for _, visibility := range servicePlan.ListVisibilities() {
		org, err := m.OrgReader.FindOrgByGUID(visibility.OrgGUID)
		if err != nil {
			return err
		}
		if m.Peek {
			lo.G.Infof("[dry-run]: removing plan %s for service %s from org %s", servicePlan.Name, servicePlan.ServiceName, org.Name)
			continue
		}
		lo.G.Infof("removing plan %s for service %s from org %s", servicePlan.Name, servicePlan.ServiceName, org.Name)
		err = m.Client.DeleteServicePlanVisibilityByPlanAndOrg(visibility.ServicePlanGUID, visibility.OrgGUID, false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) ListServiceInfo() (*ServiceInfo, error) {
	return GetServiceInfo(m.Client)
}
