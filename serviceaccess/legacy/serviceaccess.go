package legacy

import (
	"fmt"
	"net/url"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
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
	}
}

type Manager struct {
	Client    CFClient
	Cfg       config.Reader
	OrgReader organizationreader.Reader
	Peek      bool
}

func (m *Manager) Apply() error {
	globalCfg, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	if !globalCfg.EnableServiceAccess {
		lo.G.Info("Service Access is not enabled.  Set enable-service-access: true in cf-mgmt.yml")
		return nil
	}

	serviceInfo, err := m.ListServiceInfo()
	if err != nil {
		return err
	}
	err = m.DisablePublicServiceAccess(serviceInfo)
	if err != nil {
		return err
	}
	orgConfig, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}
	err = m.EnableProtectedOrgServiceAccess(serviceInfo, orgConfig.ProtectedOrgList())
	if err != nil {
		return err
	}

	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	err = m.EnableOrgServiceAccess(serviceInfo, orgConfigs)
	if err != nil {
		return err
	}

	err = m.RemoveUnknownVisibilites(serviceInfo)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUnknownVisibilites - will remove any service plan visiblities that are not known by cf-mgmt
func (m *Manager) RemoveUnknownVisibilites(serviceInfo *ServiceInfo) error {
	for servicePlanName, servicePlan := range serviceInfo.AllPlans() {
		for _, plan := range servicePlan {
			for _, visibility := range plan.ListVisibilities() {
				if m.Peek {
					lo.G.Infof("[dry-run]: removing plan %s for service %s to org with guid %s", plan.Name, servicePlanName, visibility.OrganizationGuid)
					continue
				}
				lo.G.Infof("removing plan %s for service %s to org with guid %s", plan.Name, servicePlanName, visibility.OrganizationGuid)
				err := m.Client.DeleteServicePlanVisibilityByPlanAndOrg(visibility.ServicePlanGuid, visibility.OrganizationGuid, false)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// DisablePublicServiceAccess - will ensure any public plan is made private
func (m *Manager) DisablePublicServiceAccess(serviceInfo *ServiceInfo) error {
	for _, servicePlan := range serviceInfo.AllPlans() {
		for _, plan := range servicePlan {
			if plan.Public {
				err := m.Client.MakeServicePlanPrivate(plan.GUID)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ListServiceInfo - returns services and their cooresponding plans
func (m *Manager) ListServiceInfo() (*ServiceInfo, error) {
	serviceInfo := &ServiceInfo{}
	services, err := m.Client.ListServices()
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		plans, err := m.Client.ListServicePlansByQuery(url.Values{
			"q": []string{fmt.Sprintf("%s:%s", "service_guid", service.Guid)},
		})
		if err != nil {
			return nil, err
		}
		for _, plan := range plans {
			servicePlanInfo := serviceInfo.AddPlan(service.Label, plan)
			visibilities, err := m.Client.ListServicePlanVisibilitiesByQuery(url.Values{
				"q": []string{fmt.Sprintf("%s:%s", "service_plan_guid", plan.Guid)},
			})
			if err != nil {
				return nil, err
			}
			for _, visibility := range visibilities {
				servicePlanInfo.AddOrg(visibility.OrganizationGuid, visibility)
			}
		}
	}
	return serviceInfo, nil
}

func (m *Manager) EnableProtectedOrgServiceAccess(serviceInfo *ServiceInfo, protectedOrgs []string) error {
	orgs, err := m.OrgReader.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if util.Matches(org.Name, protectedOrgs) {
			for serviceName, plans := range serviceInfo.AllPlans() {
				for _, servicePlan := range plans {
					if !servicePlan.OrgHasAccess(org.GUID) {
						if m.Peek {
							lo.G.Infof("[dry-run]: adding plan %s for service %s to org %s", servicePlan.Name, serviceName, org.Name)
							continue
						}
						lo.G.Infof("adding plan %s for service %s to org %s", servicePlan.Name, serviceName, org.Name)
						_, err = m.Client.CreateServicePlanVisibility(servicePlan.GUID, org.GUID)
						if err != nil {
							return err
						}
					} else {
						servicePlan.RemoveOrg(org.GUID)
					}
				}
			}
		}
	}
	return nil
}

func (m *Manager) EnableOrgServiceAccess(serviceInfo *ServiceInfo, orgConfigs []config.OrgConfig) error {
	for _, orgConfig := range orgConfigs {
		if orgConfig.ServiceAccess != nil {
			org, err := m.OrgReader.FindOrg(orgConfig.Org)
			if err != nil {
				return err
			}
			for serviceName, plans := range orgConfig.ServiceAccess {
				servicePlans, err := serviceInfo.GetPlans(serviceName, plans)
				if err != nil {
					lo.G.Warning(err.Error())
					continue
				}
				for _, servicePlan := range servicePlans {
					if !servicePlan.OrgHasAccess(org.GUID) {
						if m.Peek {
							lo.G.Infof("[dry-run]: adding plan %s for service %s to org %s", servicePlan.Name, serviceName, org.Name)
							continue
						}
						lo.G.Infof("adding plan %s for service %s to org %s", servicePlan.Name, serviceName, org.Name)
						_, err = m.Client.CreateServicePlanVisibility(servicePlan.GUID, org.GUID)
						if err != nil {
							return err
						}
					} else {
						servicePlan.RemoveOrg(org.GUID)
					}
				}
			}
		}
	}

	return nil
}
