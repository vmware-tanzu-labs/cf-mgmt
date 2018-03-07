package organization

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaa"
	"github.com/xchapter7x/lo"
)

func NewManager(client CFClient, uaaMgr uaa.Manager, cfg config.Reader, peek bool) Manager {
	ldapMgr := ldap.NewManager()
	UserMgr := NewUserManager(client, peek)

	return &DefaultManager{
		Cfg:     cfg,
		Client:  client,
		UAAMgr:  uaaMgr,
		LdapMgr: ldapMgr,
		UserMgr: UserMgr,
		Peek:    peek,
	}
}

//DefaultManager -
type DefaultManager struct {
	Cfg     config.Reader
	Client  CFClient
	UAAMgr  uaa.Manager
	LdapMgr ldap.Manager
	UserMgr UserMgr
	Peek    bool
}

//CreateQuotas -
func (m *DefaultManager) CreateQuotas() error {
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

func (m *DefaultManager) GetOrgGUID(orgName string) (string, error) {
	org, err := m.FindOrg(orgName)
	if err != nil {
		return "", err
	}
	return org.Guid, nil
}

//CreateOrgs -
func (m *DefaultManager) CreateOrgs() error {
	desiredOrgs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	currentOrgs, err := m.ListOrgs()
	if err != nil {
		return err
	}

	for _, org := range desiredOrgs {
		if doesOrgExist(org.Org, currentOrgs) {
			lo.G.Debugf("[%s] org already exists", org.Org)
			continue
		}
		lo.G.Infof("Creating [%s] org", org.Org)
		if err := m.CreateOrg(org.Org); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultManager) CreatePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		lo.G.Error(err)
		return err
	}

	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	allPrivateDomains, err := m.ListAllPrivateDomains()
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
				if orgGUID != existingPrivateDomain.OwningOrganizationGuid {
					existingOrgName, _ := m.getOrgName(orgs, existingPrivateDomain.OwningOrganizationGuid)
					msg := fmt.Sprintf("Private Domain %s already exists in org [%s]", privateDomain, existingOrgName)
					lo.G.Error(msg)
					return fmt.Errorf(msg)
				}
				lo.G.Debugf("Private Domain %s already exists for Org %s", privateDomain, orgConfig.Org)
			} else {
				lo.G.Infof("Creating Private Domain %s for Org %s", privateDomain, orgConfig.Org)
				privateDomain, err := m.CreatePrivateDomain(orgGUID, privateDomain)
				if err != nil {
					return err
				}
				allPrivateDomains[privateDomain.Name] = *privateDomain
			}
			privateDomainMap[privateDomain] = privateDomain
		}

		if orgConfig.RemovePrivateDomains {
			lo.G.Debugf("Looking for private domains to remove for org [%s]", orgConfig.Org)
			orgPrivateDomains, err := m.ListOrgOwnedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.DeletePrivateDomain(privateDomainGUID.Guid)
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

func (m *DefaultManager) SharePrivateDomains() error {
	orgConfigs, err := m.Cfg.GetOrgConfigs()
	if err != nil {
		return err
	}

	privateDomains, err := m.ListAllPrivateDomains()
	if err != nil {
		return err
	}
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, orgConfig := range orgConfigs {
		orgGUID, err := m.getOrgGUID(orgs, orgConfig.Org)
		if err != nil {
			return err
		}
		allSharedPrivateDomains, err := m.ListOrgSharedPrivateDomains(orgGUID)
		if err != nil {
			return err
		}

		privateDomainMap := make(map[string]string)
		for _, privateDomain := range orgConfig.SharedPrivateDomains {
			if _, ok := allSharedPrivateDomains[privateDomain]; !ok {
				if privateDomainGUID, ok := privateDomains[privateDomain]; ok {
					lo.G.Infof("Sharing Private Domain %s for Org %s", privateDomain, orgConfig.Org)
					err = m.SharePrivateDomain(orgGUID, privateDomainGUID.Guid)
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
			orgSharedPrivateDomains, err := m.ListOrgSharedPrivateDomains(orgGUID)
			if err != nil {
				return err
			}
			for existingPrivateDomain, privateDomainGUID := range orgSharedPrivateDomains {
				if _, ok := privateDomainMap[existingPrivateDomain]; !ok {
					lo.G.Infof("Removing Shared Private Domain %s for Org %s", existingPrivateDomain, orgConfig.Org)
					err = m.RemoveSharedPrivateDomain(orgGUID, privateDomainGUID.Guid)
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

func (m *DefaultManager) getOrgName(orgs []cfclient.Org, orgGUID string) (string, error) {
	for _, org := range orgs {
		if org.Guid == orgGUID {
			return org.Name, nil
		}
	}
	return "", fmt.Errorf("org for GUID %s does not exist", orgGUID)
}

func (m *DefaultManager) getOrgGUID(orgs []cfclient.Org, orgName string) (string, error) {
	for _, org := range orgs {
		if org.Name == orgName {
			return org.Guid, nil
		}
	}
	return "", fmt.Errorf("org %s does not exist", orgName)
}

//DeleteOrgs -
func (m *DefaultManager) DeleteOrgs() error {
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

	orgs, err := m.ListOrgs()
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
		if err := m.DeleteOrg(org.Guid); err != nil {
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
func (m *DefaultManager) FindOrg(orgName string) (cfclient.Org, error) {
	orgs, err := m.ListOrgs()
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
func (m *DefaultManager) UpdateOrgUsers(configDir, ldapBindPassword string) error {
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

func (m *DefaultManager) updateOrgUsers(config *ldap.Config, input *config.OrgConfig, uaacUsers map[string]string) error {
	org, err := m.FindOrg(input.Org)
	if err != nil {
		return err
	}

	err = m.syncOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetBillingManagerGroups(),
			LdapUsers:      input.BillingManager.LDAPUsers,
			Users:          input.BillingManager.Users,
			SamlUsers:      input.BillingManager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
	if err != nil {
		return err
	}

	err = m.syncOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetAuditorGroups(),
			LdapUsers:      input.Auditor.LDAPUsers,
			Users:          input.Auditor.Users,
			SamlUsers:      input.Auditor.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
	if err != nil {
		return err
	}

	return m.syncOrgUsers(
		config, uaacUsers, UpdateUsersInput{
			OrgName:        org.Name,
			OrgGUID:        org.Guid,
			LdapGroupNames: input.GetManagerGroups(),
			LdapUsers:      input.Manager.LDAPUsers,
			Users:          input.Manager.Users,
			SamlUsers:      input.Manager.SamlUsers,
			RemoveUsers:    input.RemoveUsers,
		})
}

//UpdateOrgUsers -
func (m *DefaultManager) syncOrgUsers(config *ldap.Config, uaacUsers map[string]string, updateUsersInput UpdateUsersInput) error {

	orgUsers, err := updateUsersInput.ListUsers(updateUsersInput.OrgGUID)

	if err != nil {
		return err
	}
	if config.Enabled {
		var ldapUsers []ldap.User
		ldapUsers, err = m.getLdapUsers(config, updateUsersInput)
		if err != nil {
			return err
		}
		for _, user := range ldapUsers {
			err = m.updateLdapUser(config, updateUsersInput, uaacUsers, user, orgUsers)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Debug("Skipping LDAP sync as LDAP is disabled (enable by updating config/ldap.yml)")
	}
	for _, userID := range updateUsersInput.Users {
		lowerUserID := strings.ToLower(userID)
		if _, ok := orgUsers[lowerUserID]; !ok {
			if _, userExists := uaacUsers[lowerUserID]; !userExists {
				return fmt.Errorf("User %s doesn't exist in cloud foundry, so must add internal user first", userID)
			}
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, userID); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(orgUsers, lowerUserID)
		}
	}

	for _, userEmail := range updateUsersInput.SamlUsers {
		lowerUserEmail := strings.ToLower(userEmail)
		if _, userExists := uaacUsers[lowerUserEmail]; !userExists {
			lo.G.Info("User", userEmail, "doesn't exist in cloud foundry, so creating user")
			if err = m.UAAMgr.CreateExternalUser(userEmail, userEmail, userEmail, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userEmail)
				return err
			} else {
				uaacUsers[userEmail] = userEmail
			}
		}
		if _, ok := orgUsers[lowerUserEmail]; !ok {
			if err = updateUsersInput.AddUser(updateUsersInput.OrgGUID, userEmail); err != nil {
				lo.G.Error(err)
				return err
			}
		} else {
			delete(orgUsers, lowerUserEmail)
		}
	}

	if updateUsersInput.RemoveUsers {
		lo.G.Debugf("Deleting users for org: %s", updateUsersInput.OrgName)
		for orgUser, _ := range orgUsers {
			err = updateUsersInput.RemoveUser(updateUsersInput.OrgGUID, orgUser)
			if err != nil {
				return err
			}
		}
	} else {
		lo.G.Debugf("Not removing users. Set enable-remove-users: true to orgConfig for org: %s", updateUsersInput.OrgName)
	}
	return nil
}

func (m *DefaultManager) updateLdapUser(config *ldap.Config, updateUsersInput UpdateUsersInput,
	uaacUsers map[string]string,
	user ldap.User, orgUsers map[string]string) error {

	userID := user.UserID
	externalID := user.UserDN
	if config.Origin != "ldap" {
		userID = user.Email
		externalID = user.Email
	} else {
		if user.Email == "" {
			user.Email = fmt.Sprintf("%s@user.from.ldap.cf", userID)
		}
	}
	userID = strings.ToLower(userID)

	if _, ok := orgUsers[userID]; !ok {
		if _, userExists := uaacUsers[userID]; !userExists {
			lo.G.Info("User", userID, "doesn't exist in cloud foundry, so creating user")
			if err := m.UAAMgr.CreateExternalUser(userID, user.Email, externalID, config.Origin); err != nil {
				lo.G.Error("Unable to create user", userID)
			} else {
				uaacUsers[userID] = userID
				if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, userID); err != nil {
					lo.G.Error(err)
					return err
				}
			}
		} else {
			if err := updateUsersInput.AddUser(updateUsersInput.OrgGUID, userID); err != nil {
				lo.G.Error(err)
				return err
			}
		}
	} else {
		delete(orgUsers, userID)
	}
	return nil
}

func (m *DefaultManager) getLdapUsers(config *ldap.Config, updateUsersInput UpdateUsersInput) ([]ldap.User, error) {
	users := []ldap.User{}
	for _, groupName := range updateUsersInput.LdapGroupNames {
		if groupName != "" {
			lo.G.Debug("Finding LDAP user for group:", groupName)
			if groupUsers, err := m.LdapMgr.GetUserIDs(config, groupName); err == nil {
				users = append(users, groupUsers...)
			} else {
				lo.G.Error(err)
				return nil, err
			}
		}
	}
	for _, user := range updateUsersInput.LdapUsers {
		if ldapUser, err := m.LdapMgr.GetUser(config, user); err == nil {
			if ldapUser != nil {
				users = append(users, *ldapUser)
			}
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	return users, nil
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

func (m *DefaultManager) CreateQuota(quota cloudcontroller.QuotaEntity) (string, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create quota %+v", quota)
		return "dry-run-quota-guid", nil
	}

	orgQuota, err := m.Client.CreateOrgQuota(cfclient.OrgQuotaRequest{
		Name: quota.GetName(),
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		TotalPrivateDomains:     quota.TotalPrivateDomains,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	if err != nil {
		return "", err
	}
	return orgQuota.Guid, nil
}

func (m *DefaultManager) UpdateQuota(quotaGUID string, quota cloudcontroller.QuotaEntity) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: update quota %+v with GUID %s", quota, quotaGUID)
		return nil
	}
	_, err := m.Client.UpdateOrgQuota(quotaGUID, cfclient.OrgQuotaRequest{
		Name: quota.GetName(),
		NonBasicServicesAllowed: quota.IsPaidServicesAllowed(),
		TotalServices:           quota.TotalServices,
		TotalRoutes:             quota.TotalRoutes,
		TotalPrivateDomains:     quota.TotalPrivateDomains,
		MemoryLimit:             quota.MemoryLimit,
		InstanceMemoryLimit:     quota.InstanceMemoryLimit,
		AppInstanceLimit:        quota.AppInstanceLimit,
		TotalServiceKeys:        quota.TotalServiceKeys,
		TotalReservedRoutePorts: quota.TotalReservedRoutePorts,
	})
	return err
}

func (m *DefaultManager) AssignQuotaToOrg(orgGUID, quotaGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assign quota GUID %s to org GUID %s", quotaGUID, orgGUID)
		return nil
	}
	org, err := m.Client.GetOrgByGuid(orgGUID)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateOrg(orgGUID, cfclient.OrgRequest{
		Name:                org.Name,
		QuotaDefinitionGuid: quotaGUID,
	})
	return err
}

//ListOrgs : Returns all orgs in the given foundation
func (m *DefaultManager) ListOrgs() ([]cfclient.Org, error) {
	orgs, err := m.Client.ListOrgs()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total orgs returned :", len(orgs))
	return orgs, nil
}

func (m *DefaultManager) CreateOrg(orgName string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: create org %s", orgName)
		return nil
	}
	_, err := m.Client.CreateOrg(cfclient.OrgRequest{
		Name: orgName,
	})
	return err
}

func (m *DefaultManager) DeleteOrg(orgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: delete org with GUID %s", orgGUID)
		return nil
	}
	return m.Client.DeleteOrg(orgGUID, true, true)
}

func (m *DefaultManager) DeleteOrgByName(orgName string) error {
	orgs, err := m.ListOrgs()
	if err != nil {
		return err
	}
	for _, org := range orgs {
		if org.Name == orgName {
			if m.Peek {
				lo.G.Infof("[dry-run]: delete org %s", orgName)
				return nil
			}
			m.DeleteOrg(org.Guid)
		}
	}
	return errors.New(fmt.Sprintf("org[%s] not found", orgName))
}

func (m *DefaultManager) ListAllPrivateDomains() (map[string]cfclient.Domain, error) {
	domains, err := m.Client.ListDomains()
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total private domains returned :", len(domains))
	privateDomainMap := make(map[string]cfclient.Domain)
	for _, privateDomain := range domains {
		privateDomainMap[privateDomain.Name] = privateDomain
	}
	return privateDomainMap, nil
}

func (m *DefaultManager) CreatePrivateDomain(orgGUID, privateDomain string) (*cfclient.Domain, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: create private domain %s for org GUID %s", privateDomain, orgGUID)
		return nil, nil
	}
	domain, err := m.Client.CreateDomain(privateDomain, orgGUID)
	if err != nil {
		return nil, err
	}
	return domain, nil
}
func (m *DefaultManager) SharePrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	_, err := m.Client.ShareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
	return err
}

func (m *DefaultManager) ListOrgSharedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error) {
	orgSharedPrivateDomainMap := make(map[string]cfclient.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID != privateDomain.OwningOrganizationGuid {
			orgSharedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgSharedPrivateDomainMap, nil
}

func (m *DefaultManager) listOrgPrivateDomains(orgGUID string) ([]cfclient.Domain, error) {
	privateDomains, err := m.Client.ListOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}

	lo.G.Debug("Total private domains returned :", len(privateDomains))
	return privateDomains, nil
}

func (m *DefaultManager) ListOrgOwnedPrivateDomains(orgGUID string) (map[string]cfclient.Domain, error) {
	orgOwnedPrivateDomainMap := make(map[string]cfclient.Domain)
	orgPrivateDomains, err := m.listOrgPrivateDomains(orgGUID)
	if err != nil {
		return nil, err
	}
	for _, privateDomain := range orgPrivateDomains {
		if orgGUID == privateDomain.OwningOrganizationGuid {
			orgOwnedPrivateDomainMap[privateDomain.Name] = privateDomain
		}
	}
	return orgOwnedPrivateDomainMap, nil
}

func (m *DefaultManager) DeletePrivateDomain(guid string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: Delete private domain %s", guid)
		return nil
	}
	return m.Client.DeleteDomain(guid)
}

func (m *DefaultManager) RemoveSharedPrivateDomain(sharedOrgGUID, privateDomainGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: remove share private domain %s for org GUID %s", privateDomainGUID, sharedOrgGUID)
		return nil
	}
	return m.Client.UnshareOrgPrivateDomain(sharedOrgGUID, privateDomainGUID)
}

func (m *DefaultManager) ListOrgAuditors(orgGUID string) (map[string]string, error) {
	return m.UserMgr.ListOrgAuditors(orgGUID)
}
func (m *DefaultManager) ListOrgBillingManager(orgGUID string) (map[string]string, error) {
	return m.UserMgr.ListOrgBillingManager(orgGUID)
}
func (m *DefaultManager) ListOrgManagers(orgGUID string) (map[string]string, error) {
	return m.UserMgr.ListOrgManagers(orgGUID)
}
