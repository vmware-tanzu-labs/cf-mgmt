package legacy

import (
	"context"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
	"github.com/xchapter7x/lo"
)

func NewManager(servicePlanClient CFServicePlanClient,
	servicePlanVisibilityClient CFServicePlanVisibilityClient,
	serviceOfferingClient CFServiceOfferingClient,
	orgReader organizationreader.Reader,
	cfg config.Reader, peek bool) *Manager {
	return &Manager{
		ServicePlanClient:           servicePlanClient,
		ServicePlanVisibilityClient: servicePlanVisibilityClient,
		ServiceOfferingClient:       serviceOfferingClient,
		OrgReader:                   orgReader,
		Cfg:                         cfg,
		Peek:                        peek,
	}
}

type Manager struct {
	ServicePlanClient           CFServicePlanClient
	ServicePlanVisibilityClient CFServicePlanVisibilityClient
	ServiceOfferingClient       CFServiceOfferingClient
	Cfg                         config.Reader
	OrgReader                   organizationreader.Reader
	Peek                        bool
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

	err = m.RemoveUnknownVisibilities(serviceInfo)
	if err != nil {
		return err
	}
	return nil
}

// RemoveUnknownVisibilities - will remove any service plan visibilities that are not known by cf-mgmt
func (m *Manager) RemoveUnknownVisibilities(serviceInfo *ServiceInfo) error {
	for servicePlanName, servicePlan := range serviceInfo.AllPlans() {
		for _, plan := range servicePlan {
			for _, orgVisibilityPair := range plan.ListVisibilitiesByOrg() {
				if m.Peek {
					lo.G.Infof("[dry-run]: removing plan %s for service %s to org with guid %s",
						plan.Name, servicePlanName, orgVisibilityPair.OrgGUID)
					continue
				}
				lo.G.Infof("removing plan %s for service %s to org with guid %s",
					plan.Name, servicePlanName, orgVisibilityPair.OrgGUID)
				err := m.ServicePlanVisibilityClient.Delete(context.Background(), plan.GUID, orgVisibilityPair.OrgGUID)
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
				// make the plan private by setting the type to admin
				r := resource.NewServicePlanVisibilityUpdate(resource.ServicePlanVisibilityAdmin)
				_, err := m.ServicePlanVisibilityClient.Update(context.Background(), plan.GUID, r)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// ListServiceInfo - returns services and their corresponding plans
func (m *Manager) ListServiceInfo() (*ServiceInfo, error) {
	serviceInfo := &ServiceInfo{}
	serviceOfferings, err := m.ServiceOfferingClient.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	for _, serviceOffering := range serviceOfferings {
		opts := client.NewServicePlanListOptions()
		opts.ServiceOfferingGUIDs.EqualTo(serviceOffering.GUID)
		plans, err := m.ServicePlanClient.ListAll(context.Background(), opts)
		if err != nil {
			return nil, err
		}
		for _, plan := range plans {
			servicePlanInfo := serviceInfo.AddPlan(serviceOffering.Name, plan)
			planVisibility, err := m.ServicePlanVisibilityClient.Get(context.Background(), plan.GUID)
			if err != nil {
				return nil, err
			}
			for _, orgVisibility := range planVisibility.Organizations {
				servicePlanInfo.AddOrg(orgVisibility.GUID, planVisibility)
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
						r := &resource.ServicePlanVisibility{
							Type: resource.ServicePlanVisibilityOrganization.String(),
							Organizations: []resource.ServicePlanVisibilityRelation{
								{
									GUID: org.GUID,
								},
							},
						}
						// TODO: This should be optimized to a single remote call for all orgs
						_, err = m.ServicePlanVisibilityClient.Apply(context.Background(), servicePlan.GUID, r)
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
						r := &resource.ServicePlanVisibility{
							Type: resource.ServicePlanVisibilityOrganization.String(),
							Organizations: []resource.ServicePlanVisibilityRelation{
								{
									GUID: org.GUID,
								},
							},
						}
						// TODO: This should be optimized to a single remote call for all orgs
						_, err = m.ServicePlanVisibilityClient.Apply(context.Background(), servicePlan.GUID, r)
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
