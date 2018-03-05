package organization

import (
	"fmt"
	"regexp"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func NewManager(cloudController cloudcontroller.Manager, uaaMgr uaa.Manager, cfg config.Reader) Manager {
	ldapMgr := ldap.NewManager()
	UserMgr := NewUserManager(cloudController, ldapMgr, uaaMgr)

	return &DefaultOrgManager{
		Cfg:             cfg,
		CloudController: cloudController,
		UAAMgr:          uaaMgr,
		LdapMgr:         ldapMgr,
		UserMgr:         UserMgr,
	}
}

//CreateQuotas -
func (m *DefaultOrgManager) CreateQuotas() error {
	orgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	quotas, err := m.CloudController.ListAllOrgQuotas()
	if err != nil {
		return err
	}

	for _, input := range orgs {
		if !input.EnableOrgQuota {
			continue
		}

		org, err := m.FindOrg(input.Org)
		if err != nil {
			return err
		}
		quotaName := org.Name
		quota := cloudcontroller.QuotaEntity{
			Name:                    quotaName,
			MemoryLimit:             input.MemoryLimit,
			InstanceMemoryLimit:     input.InstanceMemoryLimit,
			TotalRoutes:             input.TotalRoutes,
			TotalServices:           input.TotalServices,
			PaidServicePlansAllowed: input.PaidServicePlansAllowed,
			TotalPrivateDomains:     input.TotalPrivateDomains,
			TotalReservedRoutePorts: input.TotalReservedRoutePorts,
			TotalServiceKeys:        input.TotalServiceKeys,
			AppInstanceLimit:        input.AppInstanceLimit,
		}
		if quotaGUID, ok := quotas[quotaName]; ok {
			lo.G.Debug("Updating quota", quotaName)

			if err = m.CloudController.UpdateQuota(quotaGUID, quota); err != nil {
				return err
			}
			lo.G.Debug("Assigning", quotaName, "to", org.Name)
			if err = m.CloudController.AssignQuotaToOrg(org.Guid, quotaGUID); err != nil {
				return err
			}
		} else {
			lo.G.Debug("Creating quota", quotaName)
			targetQuotaGUID, err := m.CloudController.CreateQuota(quota)
			if err != nil {
				return err
			}
			lo.G.Debug("Assigning", quotaName, "to", org.Name)
			if err := m.CloudController.AssignQuotaToOrg(org.Guid, targetQuotaGUID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *DefaultOrgManager) GetOrgGUID(orgName string) (string, error) {
	org, err := m.FindOrg(orgName)
	if err != nil {
		return "", err
	}
	return org.Guid, nil
}

//CreateOrgs -
func (m *DefaultOrgManager) CreateOrgs() error {
	desiredOrgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	currentOrgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}

	for _, org := range desiredOrgs {
		if doesOrgExist(org.Org, currentOrgs) {
			lo.G.Debugf("[%s] org already exists", org.Org)
			continue
		}
		lo.G.Infof("Creating [%s] org", org.Org)
		if err := m.CloudController.CreateOrg(org.Org); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultOrgManager) CreatePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}
	allPrivateDomains, err := m.CloudController.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		orgGUID, err := m.getOrgGUID(orgs, orgConfig.Org)
		if err != nil {
			return err
		}
		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.PrivateDomains {
			if existingPrivateDomain, ok := allPrivateDomains[privateDomain]; ok {
				if orgGUID != existingPrivateDomain.OrgGUID {
					existingOrgName, _ := m.getOrgName(orgs, existingPrivateDomain.OrgGUID)
					msg := fmt.Sprintf("Private Domain %s already exists in org [%s]", privateDomain, existingOrgName)
					lo.G.Error(msg)
					return fmt.Errorf(msg)
				}
				lo.G.Debugf("Private Domain %s already exists for Org %s", privateDomain, orgConfig.Org)
			} else {
				lo.G.Infof("Creating Private Domain %s for Org %s", privateDomain, orgConfig.Org)
				privateDomainGUID, err := m.CloudController.CreatePrivateDomain(orgGUID, privateDomain)
				if err != nil {
					return err
				}
				allPrivateDomains[privateDomain] = cloudcontroller.PrivateDomainInfo{OrgGUID: orgGUID, PrivateDomainGUID: privateDomainGUID}
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemovePrivateDomains {
			lo.G.Debugf("Looking for private domains to remove for org [%s]", orgConfig.Org)
			orgPrivateDomains, err := m.CloudController.ListOrgOwnedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.CloudController.DeletePrivateDomain(privateDomainGUID)
					if err != nil {
						return err
					}
				}
			}
		} else {
			lo.G.Debugf("Private domains will not be removed for org [%s], must set enable-remove-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultOrgManager) SharePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	privateDomains, err := m.CloudController.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		orgGUID, err := m.getOrgGUID(orgs, orgConfig.Org)
		if err != nil {
			return err
		}
		allSharedPrivateDomains, err := m.CloudController.ListOrgSharedPrivateDomains(orgGUID)
		if err != nil {
			return err
		}

		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.SharedPrivateDomains {
			if _, ok := allSharedPrivateDomains[privateDomain]; !ok {
				if privateDomainGUID, ok := privateDomains[privateDomain]; ok {
					lo.G.Infof("Sharing Private Domain %s for Org %s", privateDomain, orgConfig.Org)
					err = m.CloudController.SharePrivateDomain(orgGUID, privateDomainGUID.PrivateDomainGUID)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("Private Domain [%s] is not defined", privateDomain)
				}
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemoveSharedPrivateDomains {
			lo.G.Debugf("Looking for shared private domains to remove for org [%s]", orgConfig.Org)
			orgSharedPrivateDomains, err := m.CloudController.ListOrgSharedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgSharedPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Shared Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.CloudController.RemoveSharedPrivateDomain(orgGUID, privateDomainGUID)
					if err != nil {
						return err
					}
				}
			}
		} else {
			lo.G.Debugf("Shared private domains will not be removed for org [%s], must set enable-remove-shared-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultOrgManager) getOrgName(orgs []cfclient.Org, orgGUID string) (string, error) {
	for _, org := range orgs {
		if org.Guid == orgGUID {
			return org.Name, nil
		}
	}
	return "", fmt.Errorf("org for GUID %s does not exist", orgGUID)
}

func (m *DefaultOrgManager) getOrgGUID(orgs []cfclient.Org, orgName string) (string, error) {
	for _, org := range orgs {
		if org.Name == orgName {
			return org.Guid, nil
		}
	}
	return "", fmt.Errorf("org %s does not exist", orgName)
}

//DeleteOrgs -
func (m *DefaultOrgManager) DeleteOrgs() error {
	orgsConfig, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}

	if !orgsConfig.EnableDeleteOrgs {
		lo.G.Debug("Org deletion is not enabled.  Set enable-delete-orgs: true")
		return nil
	}

	configuredOrgs := make(map[string]bool)
	for _, orgName := range orgsConfig.Orgs {
		configuredOrgs[orgName] = true
	}
	protectedOrgs := append(config.DefaultProtectedOrgs, orgsConfig.ProtectedOrgs...)

	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}

	orgsToDelete := make([]cfclient.Org, 0)
	for _, org := range orgs {
		if _, exists := configuredOrgs[org.Name]; !exists {
			if shouldDeleteOrg(org.Name, protectedOrgs) {
				orgsToDelete = append(orgsToDelete, org)
			} else {
				lo.G.Infof("Protected org [%s] - will not be deleted", org.Name)
			}
		}
	}

	for _, org := range orgsToDelete {
		lo.G.Infof("Deleting [%s] org", org.Name)
		if err := m.CloudController.DeleteOrg(org.Guid); err != nil {
			return err
		}
	}

	return nil
}

func shouldDeleteOrg(orgName string, protectedOrgs []string) bool {
	for _, protectedOrgName := range protectedOrgs {
		match, _ := regexp.MatchString(protectedOrgName, orgName)
		if match {
			return false
		}
	}
	return true
}

func doesOrgExist(orgName string, orgs []cfclient.Org) bool {
	for _, org := range orgs {
		if org.Name == orgName {
			return true
		}
	}
	return false
}

//FindOrg -
func (m *DefaultOrgManager) FindOrg(orgName string) (cfclient.Org, error) {
	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return cfclient.Org{}, err
	}
	for _, theOrg := range orgs {
		if theOrg.Name == orgName {
			return theOrg, nil
		}
	}
	return cfclient.Org{}, fmt.Errorf("org %q not found", orgName)
}

//UpdateOrgUsers -
func (m *DefaultOrgManager) UpdateOrgUsers(configDir, ldapBindPassword string) error {
	config, err := m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	uaacUsers, err := m.UAAMgr.ListUsers()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(config, &input, uaacUsers); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultOrgManager) updateOrgUsers(config *ldap.Config, input *config.OrgConfig, uaacUsers map[string]string) error {
	org, err := m.FindOrg(input.Org)
	if err != nil {
		return err
	}

	err = m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			Role:           "billing_managers",
			LdapGroupNames: input.GetBillingManagerGroups(),
			LdapUsers:      input.BillingManager.LDAPUsers,
			Users:          input.BillingManager.Users,
			SamlUsers:      input.BillingManager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
	if err != nil {
		return err
	}

	err = m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			Role:           "auditors",
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
	if err != nil {
		return err
	}

	return m.UserMgr.UpdateOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			Role:           "managers",
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
}
