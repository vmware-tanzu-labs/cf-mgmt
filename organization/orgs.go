package organization

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/xchapter7x/lo"
)

func NewManager(sysDomain, token, uaacToken string, cfg config.Reader) Manager {
	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)
	ldapMgr := ldap.NewManager()
	uaacMgr := uaac.NewManager(sysDomain, uaacToken)
	UserMgr := NewUserManager(cloudController, ldapMgr, uaacMgr)

	return &DefaultOrgManager{
		Cfg:             cfg,
		CloudController: cloudController,
		UAACMgr:         uaacMgr,
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
		quotaName := org.Entity.Name
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
			lo.G.Info("Updating quota", quotaName)

			if err = m.CloudController.UpdateQuota(quotaGUID, quota); err != nil {
				return err
			}
			lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
			if err = m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, quotaGUID); err != nil {
				return err
			}
		} else {
			lo.G.Info("Creating quota", quotaName)
			targetQuotaGUID, err := m.CloudController.CreateQuota(quota)
			if err != nil {
				return err
			}
			lo.G.Info("Assigning", quotaName, "to", org.Entity.Name)
			if err := m.CloudController.AssignQuotaToOrg(org.MetaData.GUID, targetQuotaGUID); err != nil {
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
	return org.MetaData.GUID, nil
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
			lo.G.Infof("[%s] org already exists", org.Org)
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
				lo.G.Infof("Private Domain %s already exists for Org %s", privateDomain, orgConfig.Org)
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
			lo.G.Infof("Looking for private domains to remove for org [%s]", orgConfig.Org)
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
			lo.G.Infof("Private domains will not be removed for org [%s], must set enable-remove-private-domains: true in orgConfig.yml", orgConfig.Org)
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
			lo.G.Infof("Looking for shared private domains to remove for org [%s]", orgConfig.Org)
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
			lo.G.Infof("Shared private domains will not be removed for org [%s], must set enable-remove-shared-private-domains: true in orgConfig.yml", orgConfig.Org)
		}
	}

	return nil
}

func (m *DefaultOrgManager) getOrgName(orgs []*cloudcontroller.Org, orgGUID string) (string, error) {
	for _, org := range orgs {
		if org.MetaData.GUID == orgGUID {
			return org.Entity.Name, nil
		}
	}
	return "", fmt.Errorf("org for GUID %s does not exist", orgGUID)
}

func (m *DefaultOrgManager) getOrgGUID(orgs []*cloudcontroller.Org, orgName string) (string, error) {
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			return org.MetaData.GUID, nil
		}
	}
	return "", fmt.Errorf("org %s does not exist", orgName)
}

//DeleteOrgs -
func (m *DefaultOrgManager) DeleteOrgs(peekDeletion bool) error {
	orgsConfig, err := m.Cfg.Orgs()
	if err != nil {
		return err
	}

	if !orgsConfig.EnableDeleteOrgs {
		lo.G.Info("Org deletion is not enabled.  Set enable-delete-orgs: true")
		return nil
	}

	configuredOrgs := make(map[string]bool)
	for _, orgName := range orgsConfig.Orgs {
		configuredOrgs[orgName] = true
	}
	protectedOrgs := config.DefaultProtectedOrgs
	for _, orgName := range orgsConfig.ProtectedOrgs {
		protectedOrgs[orgName] = true
	}

	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}

	orgsToDelete := make([]*cloudcontroller.Org, 0)
	for _, org := range orgs {
		if _, exists := configuredOrgs[org.Entity.Name]; !exists {
			if _, protected := protectedOrgs[org.Entity.Name]; !protected {
				orgsToDelete = append(orgsToDelete, org)
			} else {
				lo.G.Info(fmt.Sprintf("Protected org [%s] - will not be deleted", org.Entity.Name))
			}
		}
	}

	if peekDeletion {
		for _, org := range orgsToDelete {
			lo.G.Info(fmt.Sprintf("Peek - Would Delete [%s] org", org.Entity.Name))
		}
	} else {
		for _, org := range orgsToDelete {
			lo.G.Info(fmt.Sprintf("Deleting [%s] org", org.Entity.Name))
			if err := m.CloudController.DeleteOrg(org.MetaData.GUID); err != nil {
				return err
			}
		}
	}

	return nil
}

func doesOrgExist(orgName string, orgs []*cloudcontroller.Org) bool {
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			return true
		}
	}
	return false
}

//FindOrg -
func (m *DefaultOrgManager) FindOrg(orgName string) (*cloudcontroller.Org, error) {
	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return nil, err
	}
	for _, theOrg := range orgs {
		if theOrg.Entity.Name == orgName {
			return theOrg, nil
		}
	}
	return nil, fmt.Errorf("org %q not found", orgName)
}

//UpdateOrgUsers -
func (m *DefaultOrgManager) UpdateOrgUsers(configDir, ldapBindPassword string) error {
	config, err := m.LdapMgr.GetConfig(configDir, ldapBindPassword)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	uaacUsers, err := m.UAACMgr.ListUsers()
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
			OrgName:        org.Entity.Name,
			OrgGUID:        org.MetaData.GUID,
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
			OrgName:        org.Entity.Name,
			OrgGUID:        org.MetaData.GUID,
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
			OrgName:        org.Entity.Name,
			OrgGUID:        org.MetaData.GUID,
			Role:           "managers",
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
}
