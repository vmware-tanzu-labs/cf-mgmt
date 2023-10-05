package quota

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/organizationreader"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(
	spaceQuotaClient CFSpaceQuotaClient,
	orgQuotaClient CFOrgQuotaClient,
	spaceMgr space.Manager,
	orgReader organizationreader.Reader,
	cfg config.Reader, peek bool) *Manager {
	return &Manager{
		Cfg:              cfg,
		SpaceQuoteClient: spaceQuotaClient,
		OrgQuoteClient:   orgQuotaClient,
		SpaceMgr:         spaceMgr,
		OrgReader:        orgReader,
		Peek:             peek,
	}
}

// Manager -
type Manager struct {
	Cfg              config.Reader
	SpaceQuoteClient CFSpaceQuotaClient
	OrgQuoteClient   CFOrgQuotaClient
	SpaceMgr         space.Manager
	OrgReader        organizationreader.Reader
	Peek             bool
	SpaceQuotas      map[string]map[string]*resource.SpaceQuota
}

// CreateSpaceQuotas -
func (m *Manager) CreateSpaceQuotas() error {
	m.SpaceQuotas = nil
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		if input.NamedQuota != "" && input.EnableSpaceQuota {
			return fmt.Errorf("cannot have named quota %s and enable-space-quota for org/space %s/%s", input.NamedQuota, input.Org, input.Space)
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

			orgQuotas, err := m.ListAllOrgQuotas()
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

			if spaceQuota != nil && (space.Relationships.Quota == nil || space.Relationships.Quota.Data.GUID != spaceQuota.GUID) {
				if err = m.AssignQuotaToSpace(space, spaceQuota); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) createSpaceQuota(input config.SpaceQuota, space *resource.Space, quotas map[string]*resource.SpaceQuota, orgQuotas map[string]*resource.OrganizationQuota) error {
	quota := &resource.SpaceQuotaCreateOrUpdate{
		Name:     &input.Name,
		Apps:     &resource.SpaceQuotaApps{},
		Services: &resource.SpaceQuotaServices{},
		Routes:   &resource.SpaceQuotaRoutes{},
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
				if orgQuota.Apps.TotalMemoryInMB == nil {
					memoryLimit = nil
				} else {
					memoryLimit = orgQuota.Apps.TotalMemoryInMB
				}
			}
		}
	}

	logRateLimit, err := config.ToInteger(input.LogRateLimitBytesPerSecond)
	if err != nil {
		return err
	}

	quota.Apps.TotalInstances = appInstanceLimit
	quota.Apps.PerAppTasks = appTaskLimit
	quota.Apps.TotalMemoryInMB = memoryLimit
	quota.Apps.PerProcessMemoryInMB = instanceMemoryLimit
	quota.Apps.LogRateLimitInBytesPerSecond = logRateLimit
	quota.Routes.TotalReservedPorts = totalReservedRoutePorts
	quota.Routes.TotalRoutes = totalRoutes
	quota.Services.PaidServicesAllowed = &input.PaidServicePlansAllowed
	quota.Services.TotalServiceInstances = totalServices
	quota.Services.TotalServiceKeys = totalServiceKeys

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
	existingAppQuota := quota.Apps
	newAppQuota := newQuota.Apps
	if !reflect.DeepEqual(existingAppQuota, *newAppQuota) {
		m.debugCompareOutput("Apps Quota has changed from %s to %s", existingAppQuota, *newAppQuota)
		return true
	}
	existingRoutesQuota := quota.Routes
	newRoutesQuota := newQuota.Routes
	if !reflect.DeepEqual(existingRoutesQuota, *newRoutesQuota) {
		m.debugCompareOutput("Routes Quota has changed from %s to %s", existingRoutesQuota, *newRoutesQuota)
		return true
	}

	existingServicesQuota := quota.Services
	newServicesQuota := newQuota.Services
	if !reflect.DeepEqual(existingServicesQuota, *newServicesQuota) {
		m.debugCompareOutput("Services Quota has changed from %s to %s", existingServicesQuota, *newServicesQuota)
		return true
	}
	return false
}

func (m *Manager) debugCompareOutput(msg string, a interface{}, b interface{}) {
	aOutput, _ := json.Marshal(a)
	bOutput, _ := json.Marshal(b)
	lo.G.Debugf(msg, string(aOutput), string(bOutput))
}

func (m *Manager) ListAllSpaceQuotasForOrg(orgGUID string) (map[string]*resource.SpaceQuota, error) {
	if m.Peek && strings.Contains(orgGUID, "dry-run-org-guid") {
		return make(map[string]*resource.SpaceQuota), nil
	}
	if m.SpaceQuotas == nil {
		spaceQuotas, err := m.SpaceQuoteClient.ListAll(context.Background(), &client.SpaceQuotaListOptions{
			ListOptions: &client.ListOptions{
				PerPage: 5000,
			},
		})
		if err != nil {
			return nil, err
		}
		spaceQuotaMap := make(map[string]map[string]*resource.SpaceQuota)
		for _, spaceQuota := range spaceQuotas {
			orgGUID := spaceQuota.Relationships.Organization.Data.GUID
			if orgSpaceQuotaMap, ok := spaceQuotaMap[orgGUID]; ok {
				orgSpaceQuotaMap[spaceQuota.Name] = spaceQuota
			} else {
				orgSpaceQuotaMap := make(map[string]*resource.SpaceQuota)
				orgSpaceQuotaMap[spaceQuota.Name] = spaceQuota
				spaceQuotaMap[orgGUID] = orgSpaceQuotaMap
			}
		}
		m.SpaceQuotas = spaceQuotaMap
	}
	spaceQuotas := m.SpaceQuotas[orgGUID]
	if spaceQuotas == nil {
		spaceQuotas = make(map[string]*resource.SpaceQuota)
	}
	lo.G.Debug("Total space quotas returned :", len(spaceQuotas))
	return spaceQuotas, nil
}

func (m *Manager) UpdateSpaceQuota(quotaGUID string, quota *resource.SpaceQuotaCreateOrUpdate) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update space quota %s", *quota.Name)
		return nil
	}
	lo.G.Infof("Updating space quota %s", *quota.Name)
	_, err := m.SpaceQuoteClient.Update(context.Background(), quotaGUID, quota)
	return err
}

func (m *Manager) AssignQuotaToSpace(space *resource.Space, quota *resource.SpaceQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning quota %s to space %s", quota.Name, space.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to %s", quota.Name, space.Name)
	_, err := m.SpaceQuoteClient.Apply(context.Background(), quota.GUID, []string{space.GUID})
	return err
}

func (m *Manager) CreateSpaceQuota(quota *resource.SpaceQuotaCreateOrUpdate) (*resource.SpaceQuota, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating quota %s", *quota.Name)
		return &resource.SpaceQuota{Name: "dry-run-quota", GUID: "dry-run-guid"}, nil
	}
	lo.G.Infof("Creating quota %s", *quota.Name)
	spaceQuota, err := m.SpaceQuoteClient.Create(context.Background(), quota)
	if err != nil {
		return nil, err
	}
	return spaceQuota, nil
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
			return fmt.Errorf("cannot have named quota %s and enable-org-quota for org %s", input.NamedQuota, input.Org)
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
			if orgQuota != nil && (org.Relationships.Quota.Data == nil || org.Relationships.Quota.Data.GUID != orgQuota.GUID) {
				if err = m.AssignQuotaToOrg(org, orgQuota); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (m *Manager) createOrgQuota(input config.OrgQuota, quotas map[string]*resource.OrganizationQuota) error {

	quota := &resource.OrganizationQuotaCreateOrUpdate{
		Name:     &input.Name,
		Apps:     &resource.OrganizationQuotaApps{},
		Services: &resource.OrganizationQuotaServices{},
		Routes:   &resource.OrganizationQuotaRoutes{},
		Domains:  &resource.OrganizationQuotaDomains{},
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

	totalPrivateDomains, err := config.ToInteger(input.TotalPrivateDomains)
	if err != nil {
		return err
	}

	logRateLimit, err := config.ToInteger(input.LogRateLimitBytesPerSecond)
	if err != nil {
		return err
	}

	quota.Apps.TotalInstances = appInstanceLimit
	quota.Apps.PerAppTasks = appTaskLimit
	quota.Apps.TotalMemoryInMB = memoryLimit
	quota.Apps.PerProcessMemoryInMB = instanceMemoryLimit
	quota.Apps.LogRateLimitInBytesPerSecond = logRateLimit
	quota.Routes.TotalReservedPorts = totalReservedRoutePorts
	quota.Routes.TotalRoutes = totalRoutes
	quota.Services.PaidServicesAllowed = &input.PaidServicePlansAllowed
	quota.Services.TotalServiceInstances = totalServices
	quota.Services.TotalServiceKeys = totalServiceKeys
	quota.Domains.TotalDomains = totalPrivateDomains

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
	existingAppQuota := quota.Apps
	newAppQuota := newQuota.Apps
	if !reflect.DeepEqual(existingAppQuota, *newAppQuota) {
		m.debugCompareOutput("Apps Quota has changed from %s to %s", existingAppQuota, *newAppQuota)
		return true
	}
	existingRoutesQuota := quota.Routes
	newRoutesQuota := newQuota.Routes
	if !reflect.DeepEqual(existingRoutesQuota, *newRoutesQuota) {
		m.debugCompareOutput("Routes Quota has changed from %s to %s", existingRoutesQuota, *newRoutesQuota)
		return true
	}

	existingServicesQuota := quota.Services
	newServicesQuota := newQuota.Services
	if !reflect.DeepEqual(existingServicesQuota, *newServicesQuota) {
		m.debugCompareOutput("Services Quota has changed from %s to %s", existingServicesQuota, *newServicesQuota)
		return true
	}

	existingDomainsQuota := quota.Domains
	newDomainsQuota := newQuota.Domains
	if !reflect.DeepEqual(existingDomainsQuota, *newDomainsQuota) {
		m.debugCompareOutput("Domains Quota has changed from %s to %s", existingDomainsQuota, *newDomainsQuota)
		return true
	}
	return false
}

func (m *Manager) ListAllOrgQuotas() (map[string]*resource.OrganizationQuota, error) {
	quotas := make(map[string]*resource.OrganizationQuota)
	orgQutotas, err := m.OrgQuoteClient.ListAll(context.Background(), &client.OrganizationQuotaListOptions{
		ListOptions: &client.ListOptions{
			PerPage: 5000,
		},
	})
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
		lo.G.Infof("[dry-run]: create org quota %s", *quota.Name)
		return &resource.OrganizationQuota{Name: "dry-run-quota", GUID: "dry-run-quota-guid"}, nil
	}

	lo.G.Infof("Creating org quota %s", *quota.Name)
	orgQuota, err := m.OrgQuoteClient.Create(context.Background(), quota)
	if err != nil {
		return nil, err
	}
	return orgQuota, nil
}

func (m *Manager) UpdateOrgQuota(quotaGUID string, quota *resource.OrganizationQuotaCreateOrUpdate) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update org quota %s", *quota.Name)
		return nil
	}
	lo.G.Infof("Updating org quota %s", *quota.Name)
	_, err := m.OrgQuoteClient.Update(context.Background(), quotaGUID, quota)
	return err
}

func (m *Manager) AssignQuotaToOrg(org *resource.Organization, quota *resource.OrganizationQuota) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota %s to org %s", quota.Name, org.Name)
		return nil
	}
	lo.G.Infof("Assigning quota %s to org %s", quota.Name, org.Name)
	_, err := m.OrgQuoteClient.Apply(context.Background(), quota.GUID, []string{org.GUID})
	return err
}

func (m *Manager) GetSpaceQuota(guid string) (*resource.SpaceQuota, error) {
	return m.SpaceQuoteClient.Get(context.Background(), guid)
}

func (m *Manager) GetOrgQuota(guid string) (*resource.OrganizationQuota, error) {
	return m.OrgQuoteClient.Get(context.Background(), guid)
}
