package organization

import (
	"fmt"
	"path/filepath"

	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
	"github.com/xchapter7x/lo"
)

func NewManager(sysDomain, token, uaacToken string) (mgr Manager) {
	cloudController := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", sysDomain), token)
	ldapMgr := ldap.NewManager()
	uaacMgr := uaac.NewManager(sysDomain, uaacToken)
	UserMgr := NewUserManager(cloudController, ldapMgr, uaacMgr)

	return &DefaultOrgManager{
		CloudController: cloudController,
		UAACMgr:         uaacMgr,
		UtilsMgr:        utils.NewDefaultManager(),
		LdapMgr:         ldapMgr,
		UserMgr:         UserMgr,
	}
}

func (m *DefaultOrgManager) GetOrgConfigs(configDir string) ([]*InputUpdateOrgs, error) {
	orgConfigs := []*InputUpdateOrgs{}
	files, err := m.UtilsMgr.FindFiles(configDir, "orgConfig.yml")
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		input := &InputUpdateOrgs{
			AppInstanceLimit:        -1,
			TotalReservedRoutePorts: 0,
			TotalPrivateDomains:     -1,
			TotalServiceKeys:        -1,
		}
		if err = m.UtilsMgr.LoadFile(f, input); err != nil {
			lo.G.Error(err)
			return nil, err
		}
		orgConfigs = append(orgConfigs, input)
	}
	return orgConfigs, nil
}

//CreateQuotas -
func (m *DefaultOrgManager) CreateQuotas(configDir string) error {
	orgs, err := m.GetOrgConfigs(configDir)
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
func (m *DefaultOrgManager) CreateOrgs(configDir string) error {
	configFile := filepath.Join(configDir, "orgs.yml")
	lo.G.Info("Processing org file", configFile)
	input := &InputOrgs{}
	if err := m.UtilsMgr.LoadFile(configFile, input); err != nil {
		return err
	}
	orgs, err := m.CloudController.ListOrgs()
	if err != nil {
		return err
	}

	for _, orgName := range input.Orgs {
		if m.DoesOrgExist(orgName, orgs) {
			lo.G.Infof("[%s] org already exists", orgName)
			continue
		}
		lo.G.Infof("Creating [%s] org", orgName)
		if err := m.CloudController.CreateOrg(orgName); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultOrgManager) CreatePrivateDomains(configDir string) error {
	orgConfigs, err := m.GetOrgConfigs(configDir)
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
			if existingOrgGUID, ok := allPrivateDomains[privateDomain]; ok {
				if orgGUID != existingOrgGUID {
					existingOrgName, _ := m.getOrgName(orgs, existingOrgGUID)
					msg := fmt.Sprintf("Private Domain %s already exists in org [%s]", privateDomain, existingOrgName)
					lo.G.Error(msg)
					return fmt.Errorf(msg)
				}
				lo.G.Infof("Private Domain %s already exists for Org %s", privateDomain, orgConfig.Org)
			} else {
				lo.G.Infof("Creating Private Domain %s for Org %s", privateDomain, orgConfig.Org)
				err = m.CloudController.CreatePrivateDomain(orgGUID, privateDomain)
				if err != nil {
					return err
				}
				allPrivateDomains[privateDomain] = orgGUID
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemovePrivateDomains {
			lo.G.Infof("Looking for private domains to remove for org [%s]", orgConfig.Org)
			orgPrivateDomains, err := m.CloudController.ListOrgPrivateDomains(orgGUID)
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
func (m *DefaultOrgManager) DeleteOrgs(configDir string, peekDeletion bool) error {
	configFile := filepath.Join(configDir, "orgs.yml")
	lo.G.Info("Processing org file", configFile)
	input := &InputOrgs{}
	if err := m.UtilsMgr.LoadFile(configFile, input); err != nil {
		return err
	}

	if !input.EnableDeleteOrgs {
		lo.G.Info("Org deletion is not enabled.  Set enable-delete-orgs: true")
		return nil
	}

	configuredOrgs := make(map[string]bool)
	for _, orgName := range input.Orgs {
		configuredOrgs[orgName] = true
	}

	protectedOrgs := make(map[string]bool)
	//never allow accidental deletion of system org
	protectedOrgs["system"] = true
	for _, orgName := range input.ProtectedOrgs {
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

func (m *DefaultOrgManager) DoesOrgExist(orgName string, orgs []*cloudcontroller.Org) bool {
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

	orgConfigs, err := m.GetOrgConfigs(configDir)
	if err != nil {
		lo.G.Error(err)
		return err
	}

	for _, input := range orgConfigs {
		if err := m.updateOrgUsers(config, input, uaacUsers); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultOrgManager) updateOrgUsers(config *ldap.Config, input *InputUpdateOrgs, uaacUsers map[string]string) error {
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
			LdapUsers:      input.BillingManager.LdapUsers,
			Users:          input.BillingManager.Users,
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
			LdapUsers:      input.Auditor.LdapUsers,
			Users:          input.Auditor.Users,
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
			LdapUsers:      input.Manager.LdapUsers,
			Users:          input.Manager.Users,
			RemoveUsers:    input.RemoveUsers,
		})
}
