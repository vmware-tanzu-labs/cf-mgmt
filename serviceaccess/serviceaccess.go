package serviceaccess

import (
	"fmt"
	"net/url"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient,
	orgMgr organization.Manager,
	cfg config.Reader, peek bool) *Manager {
	return &Manager{
		Client: client,
		OrgMgr: orgMgr,
		Cfg:    cfg,
		Peek:   peek,
	}
}

type Manager struct {
	Client CFClient
	Cfg    config.Reader
	OrgMgr organization.Manager
	Peek   bool
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

	protectedOrgList, err := m.ProtectedOrgList()
	if err != nil {
		return err
	}

	for _, serviceVisibility := range globalCfg.ServiceAccess {
		servicePlans, err := m.GetServicePlans(serviceInfo, serviceVisibility)
		if err != nil {
			return err
		}
		for _, servicePlan := range servicePlans {
			if serviceVisibility.Disable {
				err = m.DisableServiceAccess(servicePlan)
				if err != nil {
					return err
				}
			} else {
				err = m.EnableOrgServiceAccess(servicePlan, serviceVisibility.Orgs, protectedOrgList)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Manager) GetServicePlans(serviceInfo *ServiceInfo, serviceVisibility config.ServiceVisibility) ([]*ServicePlanInfo, error) {
	servicePlans := []*ServicePlanInfo{}
	var plans []string
	var err error
	if serviceVisibility.Plan != "" {
		plans = append(plans, serviceVisibility.Plan)
	} else {
		plans, err = serviceInfo.GetPlanNames(serviceVisibility.Service)
		if err != nil {
			return nil, err
		}
	}
	for _, plan := range plans {
		servicePlan, err := serviceInfo.GetPlan(serviceVisibility.Service, plan)
		if err != nil {
			return nil, err
		}
		servicePlans = append(servicePlans, servicePlan)
	}

	return servicePlans, nil
}

func (m *Manager) ProtectedOrgList() ([]string, error) {
	orgConfig, err := m.Cfg.Orgs()
	if err != nil {
		return nil, err
	}
	orgs, err := m.OrgMgr.ListOrgs()
	if err != nil {
		return nil, err
	}
	orgList := []string{}
	for _, org := range orgs {
		if organization.Matches(org.Name, orgConfig.ProtectedOrgList()) {
			orgList = append(orgList, org.Name)
		}
	}
	return orgList, nil
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

func (m *Manager) DisableServiceAccess(servicePlan *ServicePlanInfo) error {
	err := m.MakePrivate(servicePlan)
	if err != nil {
		return err
	}
	return m.RemoveVisibilities(servicePlan)

}
func (m *Manager) EnableOrgServiceAccess(servicePlan *ServicePlanInfo, orgs, protectedOrgs []string) error {

	if len(orgs) == 0 {
		return m.MakePublic(servicePlan)
	}

	err := m.MakePrivate(servicePlan)
	if err != nil {
		return err
	}

	orgs = append(orgs, protectedOrgs...)
	for _, orgName := range orgs {
		org, err := m.OrgMgr.FindOrg(orgName)
		if err != nil {
			return err
		}
		if !servicePlan.OrgHasAccess(org.Guid) {
			if m.Peek {
				lo.G.Infof("[dry-run]: adding plan %s for service %s to org %s", servicePlan.Name, servicePlan.ServiceName, org.Name)
				continue
			}
			lo.G.Infof("adding plan %s for service %s to org %s", servicePlan.Name, servicePlan.ServiceName, org.Name)
			_, err = m.Client.CreateServicePlanVisibility(servicePlan.GUID, org.Guid)
			if err != nil {
				return err
			}
		} else {
			servicePlan.RemoveOrg(org.Guid)
		}
	}

	return m.RemoveVisibilities(servicePlan)
}

func (m *Manager) RemoveVisibilities(servicePlan *ServicePlanInfo) error {
	for _, visibility := range servicePlan.ListVisibilities() {
		if m.Peek {
			lo.G.Infof("[dry-run]: removing plan %s for service %s to org with guid %s", servicePlan.Name, servicePlan.ServiceName, visibility.OrganizationGuid)
			continue
		}
		lo.G.Infof("removing plan %s for service %s to org with guid %s", servicePlan.Name, servicePlan.ServiceName, visibility.OrganizationGuid)
		err := m.Client.DeleteServicePlanVisibilityByPlanAndOrg(visibility.ServicePlanGuid, visibility.OrganizationGuid, false)
		if err != nil {
			return err
		}
	}
	return nil
}

//ListServiceInfo - returns services and their cooresponding plans
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
