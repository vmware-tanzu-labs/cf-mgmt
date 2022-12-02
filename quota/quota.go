package quota

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"reflect"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organization"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(spaceQuotaClient CFSpaceQuotaClient,
	organizationQuotaClient CFOrganizationQuotaClient,
	spaceMgr space.Manager,
	orgReader organizationreader.Reader,
	orgMgr organization.Manager,
	cfg config.Reader, peek bool) *Manager {
	return &Manager{
		Cfg:                     cfg,
		SpaceQuotaClient:        spaceQuotaClient,
		OrganizationQuotaClient: organizationQuotaClient,
		SpaceMgr:                spaceMgr,
		OrgReader:               orgReader,
		OrgMgr:                  orgMgr,
		Peek:                    peek,
	}
}

// Manager -
type Manager struct {
	Cfg                     config.Reader
	SpaceQuotaClient        CFSpaceQuotaClient
	OrganizationQuotaClient CFOrganizationQuotaClient
	SpaceMgr                space.Manager
	OrgReader               organizationreader.Reader
	OrgMgr                  organization.Manager
	Peek                    bool
}

// CreateSpaceQuotas -
func (m *Manager) CreateSpaceQuotas() error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		if input.NamedQuota != "" && input.EnableSpaceQuota {
			return fmt.Errorf("Cannot have named quota %s and enable-space-quota for org/space %s/%s", input.NamedQuota, input.Org, input.Space)
		}
		if input.NamedQuota != "" || input.EnableSpaceQuota {
			space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
			if err != nil {
				return errors.Wrap(err, "Finding spaces")
			}
			quotas, err := m.ListAllSpaceQuotasForOrg(space.Relationships.Organization.Data.GUID)
			if err != nil {
				return errors.Wrap(err, "ListAllSpaceQuotasForOrg")
			}

			orgQuotas, err := m.OrganizationQuotaClient.ListAll(context.Background(), nil)
			if err != nil {
				return err
			}
			if input.NamedQuota != "" {
				spaceQuotas, err := m.Cfg.GetSpaceQuotas(input.Org)
				if err != nil {
					return err
				}

				for _, spaceQuotaConfig := range spaceQuotas {
					err = m.createSpaceQuota(spaceQuotaConfig, space, quotas, orgQuotas)
					if err != nil {
						return err
					}
				}
			} else {
				if input.EnableSpaceQuota {
					quotaDef := input.GetQuota()
					err = m.createSpaceQuota(quotaDef, space, quotas, orgQuotas)
					if err != nil {
						return err
					}
					input.NamedQuota = input.Space
				}
			}
			spaceQuota := quotas[input.NamedQuota]
			if spaceQuotaGUID(space) != spaceQuota.GUID {
				if err = m.AssignQuotaToSpace(space, spaceQuota); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) createSpaceQuota(input config.SpaceQuota, space *resource.Space, quotas map[string]*resource.SpaceQuota, orgQuotas []*resource.OrganizationQuota) error {

	instanceMemoryLimit, err := config.ToMegabytes(input.InstanceMemoryLimit)
	if err != nil {
		return err
	}

	totalRoutes, err := config.ToInteger(input.TotalRoutes)
	if err != nil {
		return err
	}

	totalServices, err := config.ToInteger(input.TotalServices)
	if err != nil {
		return err
	}

	totalReservedRoutePorts, err := config.ToInteger(input.TotalReservedRoutePorts)
	if err != nil {
		return err
	}

	totalServiceKeys, err := config.ToInteger(input.TotalServiceKeys)
	if err != nil {
		return err
	}

	appInstanceLimit, err := config.ToInteger(input.AppInstanceLimit)
	if err != nil {
		return err
	}

	appTaskLimit, err := config.ToInteger(input.AppTaskLimit)
	if err != nil {
		return err
	}

	memoryLimit, err := config.ToMegabytes(input.MemoryLimit)
	if err != nil {
		return err
	}

	if input.IsUnlimitedMemory() {
		org, err := m.OrgReader.FindOrg(input.Org)
		if err != nil {
			return err
		}
		for _, orgQuota := range orgQuotas {
			if org.Relationships.Quota.Data.GUID == orgQuota.GUID {
				memoryLimit = *orgQuota.Apps.TotalMemoryInMB
			}
		}
	}
	quota := resource.NewSpaceQuotaUpdate().
		WithName(input.Name).
		WithTotalMemoryInMB(memoryLimit).
		WithPerProcessMemoryInMB(instanceMemoryLimit).
		WithTotalRoutes(totalRoutes).
		WithTotalServiceInstances(totalServices).
		WithPaidServicesAllowed(input.PaidServicePlansAllowed).
		WithTotalReservedPorts(totalReservedRoutePorts).
		WithTotalServiceKeys(totalServiceKeys).
		WithTotalInstances(appInstanceLimit).
		WithPerAppTasks(appTaskLimit)

	if spaceQuota, ok := quotas[input.Name]; ok {
		if m.hasSpaceQuotaChanged(spaceQuota, quota) {
			if err := m.UpdateSpaceQuota(spaceQuota.GUID, quota); err != nil {
				return err
			}
		}
	} else {
		createdQuota, err := m.CreateSpaceQuota(quota)
		if err != nil {
			return err
		}
		quotas[input.Name] = createdQuota
	}
	return nil
}

func (m *Manager) hasSpaceQuotaChanged(quota *resource.SpaceQuota, newQuota *resource.SpaceQuotaCreateOrUpdate) bool {
	// the v2 API used to allow you to update the quota org id, v3 does not
	quoteRequest := resource.NewSpaceQuotaUpdate().
		WithName(quota.Name).
		WithTotalMemoryInMB(*quota.Apps.TotalMemoryInMB).
		WithPerProcessMemoryInMB(*quota.Apps.PerProcessMemoryInMB).
		WithTotalInstances(*quota.Apps.TotalInstances).
		WithPerAppTasks(*quota.Apps.PerAppTasks).
		WithTotalServiceInstances(*quota.Services.TotalServiceInstances).
		WithPaidServicesAllowed(*quota.Services.PaidServicesAllowed).
		WithTotalServiceKeys(*quota.Services.TotalServiceKeys).
		WithTotalReservedPorts(*quota.Routes.TotalReservedPorts).
		WithTotalRoutes(*quota.Routes.TotalRoutes)

	if reflect.DeepEqual(quoteRequest, newQuota) {
		return false
	} else {
		lo.G.Debugf("Quota has changed from %v to %v", quoteRequest, newQuota)
		return true
	}
}

func (m *Manager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]*resource.SpaceQuota, error) {
	quotas := make(map[string]*resource.SpaceQuota)
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return quotas, nil
	}
	opts := cfclient.NewSpaceQuotaListOptions()
	opts.OrganizationGUIDs.EqualTo(orgGUID)
	spaceQuotas, err := m.SpaceQuotaClient.ListAll(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	for _, quota := range spaceQuotas {
		quotas[quota.Name] = quota
	}
	return quotas, nil
}

func (m *Manager) UpdateSpaceQuota(quotaGUID string, quota *resource.SpaceQuotaCreateOrUpdate) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update space quota %s", quota.Name)
		return nil
	}
	lo.G.Infof("Updating space quota %s", quota.Name)
	_, err := m.SpaceQuotaClient.Update(context.Background(), quotaGUID, quota)
	return err
}

func (m *Manager) AssignQuotaToSpace(space *resource.Space, quota *resource.SpaceQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quota %s to space %s", quota.Name, space.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to %s", quota.Name, space.Name)
	_, err := m.SpaceQuotaClient.Apply(context.Background(), quota.GUID, []string{space.GUID})
	return err
}

func (m *Manager) CreateSpaceQuota(quota *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %s", quota.Name)
		return &resource.SpaceQuota{Name: "dry-run-quota", GUID: "dry-run-guid"}, nil
	}
	lo.G.Infof("Creating quota %s", quota.Name)
	spaceQuota, err := m.SpaceQuotaClient.Create(context.Background(), quota)
	if err != nil {
		return nil, err
	}
	return spaceQuota, nil
}

func (m *Manager) SpaceQuotaByName(name string) (*resource.SpaceQuota, error) {
	opts := cfclient.NewSpaceQuotaListOptions()
	opts.Names.EqualTo(name)
	return m.SpaceQuotaClient.Single(context.Background(), opts)
}

// CreateOrgQuotas -
func (m *Manager) CreateOrgQuotas() error {
	quotas, err := m.ListAllOrgQuotas()
	if err != nil {
		return err
	}

	orgQuotas, err := m.Cfg.GetOrgQuotas()
	if err != nil {
		return err
	}
	for _, orgQuotaConfig := range orgQuotas {
		err = m.createOrgQuota(orgQuotaConfig, quotas)
		if err != nil {
			return err
		}
	}
	orgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	for _, input := range orgs {
		if input.NamedQuota != "" && input.EnableOrgQuota {
			return fmt.Errorf("Cannot have named quota %s and enable-org-quota for org %s", input.NamedQuota, input.Org)
		}
		if input.EnableOrgQuota || input.NamedQuota != "" {
			org, err := m.OrgReader.FindOrg(input.Org)
			if err != nil {
				return err
			}
			if input.EnableOrgQuota {
				quotaDef := input.GetQuota()
				err = m.createOrgQuota(quotaDef, quotas)
				if err != nil {
					return err
				}
				input.NamedQuota = input.Org
			}
			orgQuota := quotas[input.NamedQuota]
			if orgQuotaGUID(org) != orgQuota.GUID {
				if err = m.AssignQuotaToOrg(org, orgQuota); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) createOrgQuota(input config.OrgQuota, quotas map[string]*resource.OrganizationQuota) error {
	memoryLimit, err := config.ToMegabytes(input.MemoryLimit)
	if err != nil {
		return err
	}

	instanceMemoryLimit, err := config.ToMegabytes(input.InstanceMemoryLimit)
	if err != nil {
		return err
	}

	totalRoutes, err := config.ToInteger(input.TotalRoutes)
	if err != nil {
		return err
	}

	totalServices, err := config.ToInteger(input.TotalServices)
	if err != nil {
		return err
	}

	totalReservedRoutePorts, err := config.ToInteger(input.TotalReservedRoutePorts)
	if err != nil {
		return err
	}

	totalServiceKeys, err := config.ToInteger(input.TotalServiceKeys)
	if err != nil {
		return err
	}

	appInstanceLimit, err := config.ToInteger(input.AppInstanceLimit)
	if err != nil {
		return err
	}

	appTaskLimit, err := config.ToInteger(input.AppTaskLimit)
	if err != nil {
		return err
	}

	totalPrivateDomains, err := config.ToInteger(input.TotalPrivateDomains)
	if err != nil {
		return err
	}

	quota := resource.NewOrganizationQuotaUpdate().
		WithName(input.Name).
		WithDomains(totalPrivateDomains).
		WithAppsTotalMemoryInMB(memoryLimit).
		WithPerProcessMemoryInMB(instanceMemoryLimit).
		WithTotalRoutes(totalRoutes).
		WithTotalServiceInstances(totalServices).
		WithPaidServicesAllowed(input.PaidServicePlansAllowed).
		WithTotalReservedPorts(totalReservedRoutePorts).
		WithTotalServiceKeys(totalServiceKeys).
		WithTotalInstances(appInstanceLimit).
		WithPerAppTasks(appTaskLimit)
	if orgQuota, ok := quotas[input.Name]; ok {
		if m.hasOrgQuotaChanged(orgQuota, quota) {
			if err = m.UpdateOrgQuota(orgQuota.GUID, quota); err != nil {
				return err
			}
		}
	} else {
		createdQuota, err := m.CreateOrgQuota(quota)
		if err != nil {
			return err
		}
		quotas[input.Name] = createdQuota
	}

	return nil
}

func (m *Manager) hasOrgQuotaChanged(quota *resource.OrganizationQuota, newQuota *resource.OrganizationQuotaCreateOrUpdate) bool {
	quoteRequest := resource.NewOrganizationQuotaUpdate().
		WithName(quota.Name).
		WithDomains(*quota.Domains.TotalDomains).
		WithAppsTotalMemoryInMB(*quota.Apps.TotalMemoryInMB).
		WithPerProcessMemoryInMB(*quota.Apps.PerProcessMemoryInMB).
		WithTotalRoutes(*quota.Routes.TotalRoutes).
		WithTotalServiceInstances(*quota.Services.TotalServiceInstances).
		WithPaidServicesAllowed(*quota.Services.PaidServicesAllowed).
		WithTotalReservedPorts(*quota.Routes.TotalReservedPorts).
		WithTotalServiceKeys(*quota.Services.TotalServiceKeys).
		WithTotalInstances(*quota.Apps.TotalInstances).
		WithPerAppTasks(*quota.Apps.PerAppTasks)
	if reflect.DeepEqual(quoteRequest, newQuota) {
		return false
	} else {
		lo.G.Debugf("Quota has changed from %v to %v", quoteRequest, newQuota)
		return true
	}
}

func (m *Manager) ListAllOrgQuotas() (map[string]*resource.OrganizationQuota, error) {
	quotas := make(map[string]*resource.OrganizationQuota)
	orgQutotas, err := m.OrganizationQuotaClient.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total org quotas returned :", len(orgQutotas))
	for _, quota := range orgQutotas {
		quotas[quota.Name] = quota
	}
	return quotas, nil
}

func (m *Manager) CreateOrgQuota(quota *resource.OrganizationQuotaCreateOrUpdate) (*resource.OrganizationQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org quota %s", quota.Name)
		return &resource.OrganizationQuota{Name: "dry-run-quota", GUID: "dry-run-quota-guid"}, nil
	}

	lo.G.Infof("Creating org quota %s", quota.Name)
	orgQuota, err := m.OrganizationQuotaClient.Create(context.Background(), quota)
	if err != nil {
		return nil, err
	}
	return orgQuota, nil
}

func (m *Manager) UpdateOrgQuota(quotaGUID string, quota *resource.OrganizationQuotaCreateOrUpdate) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update org quota %s", quota.Name)
		return nil
	}
	lo.G.Infof("Updating org quota %s", quota.Name)
	_, err := m.OrganizationQuotaClient.Update(context.Background(), quotaGUID, quota)
	return err
}

func (m *Manager) AssignQuotaToOrg(org *resource.Organization, quota *resource.OrganizationQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota %s to org %s", quota.Name, org.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to org %s", quota.Name, org.Name)
	_, err := m.OrganizationQuotaClient.Apply(context.Background(), quota.GUID, []string{org.GUID})
	return err
}

func (m *Manager) OrgQuotaByName(name string) (*resource.OrganizationQuota, error) {
	opts := cfclient.NewOrganizationQuotaListOptions()
	opts.Names.EqualTo(name)
	return m.OrganizationQuotaClient.Single(context.Background(), opts)
}

func spaceQuotaGUID(space *resource.Space) string {
	if space.Relationships != nil && space.Relationships.Quota != nil {
		return space.Relationships.Quota.Data.GUID
	}
	return ""
}

func orgQuotaGUID(org *resource.Organization) string {
	if org.Relationships.Quota.Data != nil {
		return org.Relationships.Quota.Data.GUID
	}
	return ""
}
