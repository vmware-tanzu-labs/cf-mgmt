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
			continue
		}
		quotaName := space.Name
		quotas, err := m.ListAllSpaceQuotasForOrg(space.OrganizationGuid)
		if err != nil {
			continue
		}

		quota := cfclient.SpaceQuota{
			OrganizationGuid: space.OrganizationGuid, Name: quotaName,
			MemoryLimit:             input.MemoryLimit,
			InstanceMemoryLimit:     input.InstanceMemoryLimit,
			TotalRoutes:             input.TotalRoutes,
			TotalServices:           input.TotalServices,
			NonBasicServicesAllowed: input.PaidServicePlansAllowed,
			TotalReservedRoutePorts: input.TotalReservedRoutePorts,
			TotalServiceKeys:        input.TotalServiceKeys,
			AppInstanceLimit:        input.AppInstanceLimit,
		}
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Debug("Updating quota", quotaName)
			if err := m.UpdateSpaceQuota(quotaGUID, quota); err != nil {
				continue
			}
			lo.G.Infof("Assigning %s to %s", quotaName, space.Name)
			return m.AssignQuotaToSpace(space.Guid, quotaGUID)
		} else {
			lo.G.Debug("Creating quota", quotaName)
			spaceQuota, err := m.CreateSpaceQuota(quota)
			if err != nil {
				continue
			}
			lo.G.Infof("Assigning %s to %s", quotaName, space.Name)
			return m.AssignQuotaToSpace(space.Guid, spaceQuota.Guid)
		}
	}
	return nil
}

func (m *DefaultManager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]string, error) {
	quotas := make(map[string]string)
	spaceQuotas, err := m.Client.ListOrgSpaceQuotas(orgGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	for _, quota := range spaceQuotas {
		quotas[quota.Name] = quota.Guid
	}
	return quotas, nil
}

func (m *DefaultManager) UpdateSpaceQuota(quotaGUID string, quota cfclient.SpaceQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %s with %+v", quotaGUID, quota)
		return nil
	}
	_, err := m.Client.UpdateSpaceQuota(quotaGUID, cfclient.SpaceQuotaRequest{
		Name:                    quota.Name,
		OrganizationGuid:        quota.OrganizationGuid,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	return err
}

func (m *DefaultManager) AssignQuotaToSpace(spaceGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quotaGUID %s to spaceGUID %s", quotaGUID, spaceGUID)
		return nil
	}
	return m.Client.AssignSpaceQuota(quotaGUID, spaceGUID)
}

func (m *DefaultManager) CreateSpaceQuota(quota cfclient.SpaceQuota) (*cfclient.SpaceQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %+v", quota)
		return nil, nil
	}
	spaceQuota, err := m.Client.CreateSpaceQuota(cfclient.SpaceQuotaRequest{
		Name:                    quota.Name,
		OrganizationGuid:        quota.OrganizationGuid,
		NonBasicServicesAllowed: quota.NonBasicServicesAllowed,
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
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
		quota := cfclient.OrgQuotaRequest{
			Name:                    quotaName,
			MemoryLimit:             input.MemoryLimit,
			InstanceMemoryLimit:     input.InstanceMemoryLimit,
			TotalRoutes:             input.TotalRoutes,
			TotalServices:           input.TotalServices,
			NonBasicServicesAllowed: input.PaidServicePlansAllowed,
			TotalPrivateDomains:     input.TotalPrivateDomains,
			TotalReservedRoutePorts: input.TotalReservedRoutePorts,
			TotalServiceKeys:        input.TotalServiceKeys,
			AppInstanceLimit:        input.AppInstanceLimit,
		}
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Debug("Updating quota", quotaName)

			if err = m.UpdateQuota(quotaGUID, quota); err != nil {
				return err
			}
			lo.G.Debug("Assigning", quotaName, "to", org.Name)
			if err = m.AssignQuotaToOrg(org.Guid, quotaGUID); err != nil {
				return err
			}
		} else {
			lo.G.Debug("Creating quota", quotaName)
			targetQuotaGUID, err := m.CreateQuota(quota)
			if err != nil {
				return err
			}
			lo.G.Debug("Assigning", quotaName, "to", org.Name)
			if err := m.AssignQuotaToOrg(org.Guid, targetQuotaGUID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultManager) ListAllOrgQuotas() (map[string]string, error) {
	quotas := make(map[string]string)
	orgQutotas, err := m.Client.ListOrgQuotas()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total org quotas returned :", len(orgQutotas))
	for _, quota := range orgQutotas {
		quotas[quota.Name] = quota.Guid
	}
	return quotas, nil
}

func (m *DefaultManager) CreateQuota(quota cfclient.OrgQuotaRequest) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create quota %+v", quota)
		return "dry-run-quota-guid", nil
	}

	orgQuota, err := m.Client.CreateOrgQuota(quota)
	if err != nil {
		return "", err
	}
	return orgQuota.Guid, nil
}

func (m *DefaultManager) UpdateQuota(quotaGUID string, quota cfclient.OrgQuotaRequest) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %+v with GUID %s", quota, quotaGUID)
		return nil
	}
	_, err := m.Client.UpdateOrgQuota(quotaGUID, quota)
	return err
}

func (m *DefaultManager) AssignQuotaToOrg(orgGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota GUID %s to org GUID %s", quotaGUID, orgGUID)
		return nil
	}
	org, err := m.OrgMgr.GetOrgByGUID(orgGUID)
	if err != nil {
		return err
	}
	_, err = m.OrgMgr.UpdateOrg(orgGUID, cfclient.OrgRequest{
		Name:                org.Name,
		QuotaDefinitionGuid: quotaGUID,
	})
	return err
}

func (m *DefaultManager) OrgQuotaByName(name string) (cfclient.OrgQuota, error) {
	return m.Client.GetOrgQuotaByName(name)
}
