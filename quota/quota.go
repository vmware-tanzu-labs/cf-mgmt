package quota

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/space"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(client CFClient,
	spaceMgr space.Manager,
	orgMgr organization.Manager,
	cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:      cfg,
		Client:   client,
		SpaceMgr: spaceMgr,
		OrgMgr:   orgMgr,
		Peek:     peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg      config.Reader
	Client   CFClient
	SpaceMgr space.Manager
	OrgMgr   organization.Manager
	Peek     bool
}

//CreateSpaceQuotas -
func (m *DefaultManager) CreateSpaceQuotas() error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}
	for _, input := range spaceConfigs {
		if !input.EnableSpaceQuota {
			continue
		}
		space, err := m.SpaceMgr.FindSpace(input.Org, input.Space)
		if err != nil {
			return err
		}
		quotas, err := m.ListAllSpaceQuotasForOrg(space.OrganizationGuid)
		if err != nil {
			return err
		}

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

		quota := cfclient.SpaceQuotaRequest{
			Name:                    space.Name,
			OrganizationGuid:        space.OrganizationGuid,
			MemoryLimit:             memoryLimit,
			InstanceMemoryLimit:     instanceMemoryLimit,
			TotalRoutes:             totalRoutes,
			TotalServices:           totalServices,
			NonBasicServicesAllowed: input.PaidServicePlansAllowed,
			TotalReservedRoutePorts: totalReservedRoutePorts,
			TotalServiceKeys:        totalServiceKeys,
			AppInstanceLimit:        appInstanceLimit,
			AppTaskLimit:            appTaskLimit,
		}
		var spaceQuota cfclient.SpaceQuota
		var ok bool
		if spaceQuota, ok = quotas[space.Name]; ok {
			if m.hasSpaceQuotaChanged(spaceQuota, quota) {
				if err := m.UpdateSpaceQuota(spaceQuota.Guid, quota); err != nil {
					return err
				}
			}
		} else {
			createdQuota, err := m.CreateSpaceQuota(quota)
			if err != nil {
				return err
			}
			spaceQuota = *createdQuota
		}
		if space.QuotaDefinitionGuid != spaceQuota.Guid {
			err := m.AssignQuotaToSpace(space, spaceQuota)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) hasSpaceQuotaChanged(quota cfclient.SpaceQuota, newQuota cfclient.SpaceQuotaRequest) bool {
	quoteRequest := cfclient.SpaceQuotaRequest{
		Name:                    quota.Name,
		OrganizationGuid:        quota.OrganizationGuid,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		TotalRoutes:             quota.TotalRoutes,
		TotalServices:           quota.TotalServices,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
		TotalServiceKeys:        quota.TotalServiceKeys,
		AppInstanceLimit:        quota.AppInstanceLimit,
		AppTaskLimit:            quota.AppTaskLimit,
	}
	if quoteRequest == newQuota {
		return false
	} else {
		lo.G.Debugf("Quota has changed from %v to %v", quoteRequest, newQuota)
		return true
	}
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]cfclient.SpaceQuota, error) {
	quotas := make(map[string]cfclient.SpaceQuota)
	spaceQuotas, err := m.Client.ListOrgSpaceQuotas(orgGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	for _, quota := range spaceQuotas {
		quotas[quota.Name] = quota
	}
	return quotas, nil
}

func (m *DefaultManager) UpdateSpaceQuota(quotaGUID string, quota cfclient.SpaceQuotaRequest) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update space quota %s", quota.Name)
		return nil
	}
	lo.G.Infof("Updating space quota %s", quota.Name)
	_, err := m.Client.UpdateSpaceQuota(quotaGUID, quota)
	return err
}

func (m *DefaultManager) AssignQuotaToSpace(space cfclient.Space, quota cfclient.SpaceQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quota %s to space %s", quota.Name, space.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to %s", quota.Name, space.Name)
	return m.Client.AssignSpaceQuota(quota.Guid, space.Guid)
}

func (m *DefaultManager) CreateSpaceQuota(quota cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %s", quota.Name)
		return &cfclient.SpaceQuota{Name: "dry-run-quota", Guid: "dry-run-guid"}, nil
	}
	lo.G.Infof("Creating quota %s", quota.Name)
	spaceQuota, err := m.Client.CreateSpaceQuota(quota)
	if err != nil {
		return nil, err
	}
	return spaceQuota, nil
}

func (m *DefaultManager) SpaceQuotaByName(name string) (cfclient.SpaceQuota, error) {
	return m.Client.GetSpaceQuotaByName(name)
}

//CreateOrgQuotas -
func (m *DefaultManager) CreateOrgQuotas() error {
	orgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	quotas, err := m.ListAllOrgQuotas()
	if err != nil {
		return err
	}

	for _, input := range orgs {
		if !input.EnableOrgQuota {
			continue
		}

		org, err := m.OrgMgr.FindOrg(input.Org)
		if err != nil {
			return err
		}
		quotaName := org.Name
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


		quota := cfclient.OrgQuotaRequest{
			Name:                    quotaName,
			MemoryLimit:             memoryLimit,
			InstanceMemoryLimit:     instanceMemoryLimit,
			TotalRoutes:             totalRoutes,
			TotalServices:           totalServices,
			NonBasicServicesAllowed: input.PaidServicePlansAllowed,
			TotalPrivateDomains:     totalPrivateDomains,
			TotalReservedRoutePorts: totalReservedRoutePorts,
			TotalServiceKeys:        totalServiceKeys,
			AppInstanceLimit:        appInstanceLimit,
			AppTaskLimit:            appTaskLimit,
		}
		var orgQuota cfclient.OrgQuota
		var ok bool
		if orgQuota, ok = quotas[quotaName]; ok {
			if m.hasOrgQuotaChanged(orgQuota, quota) {
				if err = m.UpdateOrgQuota(orgQuota.Guid, quota); err != nil {
					return err
				}
			}
		} else {
			createdQuota, err := m.CreateOrgQuota(quota)
			if err != nil {
				return err
			}
			orgQuota = *createdQuota
		}
		if org.QuotaDefinitionGuid != orgQuota.Guid {
			if err = m.AssignQuotaToOrg(org, orgQuota); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) hasOrgQuotaChanged(quota cfclient.OrgQuota, newQuota cfclient.OrgQuotaRequest) bool {
	quoteRequest := cfclient.OrgQuotaRequest{
		Name:                    quota.Name,
		TotalPrivateDomains:     quota.TotalPrivateDomains,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		TotalRoutes:             quota.TotalRoutes,
		TotalServices:           quota.TotalServices,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
		TotalServiceKeys:        quota.TotalServiceKeys,
		AppInstanceLimit:        quota.AppInstanceLimit,
		AppTaskLimit:            quota.AppTaskLimit,
	}
	if quoteRequest == newQuota {
		return false
	} else {
		lo.G.Debugf("Quota has changed from %v to %v", quoteRequest, newQuota)
		return true
	}
}

func (m *DefaultManager) ListAllOrgQuotas() (map[string]cfclient.OrgQuota, error) {
	quotas := make(map[string]cfclient.OrgQuota)
	orgQutotas, err := m.Client.ListOrgQuotas()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total org quotas returned :", len(orgQutotas))
	for _, quota := range orgQutotas {
		quotas[quota.Name] = quota
	}
	return quotas, nil
}

func (m *DefaultManager) CreateOrgQuota(quota cfclient.OrgQuotaRequest) (*cfclient.OrgQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org quota %s", quota.Name)
		return &cfclient.OrgQuota{Name: "dry-run-quota", Guid: "dry-run-quota-guid"}, nil
	}

	lo.G.Infof("Creating org quota %s", quota.Name)
	orgQuota, err := m.Client.CreateOrgQuota(quota)
	if err != nil {
		return nil, err
	}
	return orgQuota, nil
}

func (m *DefaultManager) UpdateOrgQuota(quotaGUID string, quota cfclient.OrgQuotaRequest) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update org quota %s", quota.Name)
		return nil
	}
	lo.G.Infof("Updating org quota %s", quota.Name)
	_, err := m.Client.UpdateOrgQuota(quotaGUID, quota)
	return err
}

func (m *DefaultManager) AssignQuotaToOrg(org cfclient.Org, quota cfclient.OrgQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota %s to org %s", quota.Name, org.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to org %s", quota.Name, org.Name)
	_, err := m.OrgMgr.UpdateOrg(org.Guid, cfclient.OrgRequest{
		Name:                org.Name,
		QuotaDefinitionGuid: quota.Guid,
	})
	return err
}

func (m *DefaultManager) OrgQuotaByName(name string) (cfclient.OrgQuota, error) {
	return m.Client.GetOrgQuotaByName(name)
}
