package organization_test_data

import (
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	mock "github.com/pivotalservices/cf-mgmt/utils/mocks"
)

func PopulateWithTestData(utilsMgrMock *mock.MockUtilsManager) error {
	utilsMgrMock.MockFileData["./fixtures/config/ldap.yml"] = ldap.Config{Enabled: true, LdapHost: "127.0.0.1", LdapPort: 10389, TLS: false, BindDN: "uid=admin,ou=system", BindPassword: "secret", UserSearchBase: "ou=users,dc=example,dc=com", UserNameAttribute: "uid", UserMailAttribute: "mail", UserObjectClass: "", GroupSearchBase: "ou=groups,dc=example,dc=com", GroupAttribute: "member", Origin: ""}
	utilsMgrMock.MockFileData["./fixtures/config/orgs.yml"] = config.Orgs{Orgs: []string{"test", "test2"}, EnableDeleteOrgs: false, ProtectedOrgs: []string(nil)}
	utilsMgrMock.MockFileData["./fixtures/config/test/orgConfig.yml"] = config.OrgConfig{Org: "test", BillingManagerGroup: "test_billing_managers", ManagerGroup: "", AuditorGroup: "", BillingManager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Manager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Auditor: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, PrivateDomains: []string(nil), RemovePrivateDomains: false, EnableOrgQuota: true, MemoryLimit: 10240, InstanceMemoryLimit: -1, TotalRoutes: 10, TotalServices: -1, PaidServicePlansAllowed: true, RemoveUsers: false, TotalPrivateDomains: 0, TotalReservedRoutePorts: 0, TotalServiceKeys: 0, AppInstanceLimit: 0, DefaultIsoSegment: ""}
	utilsMgrMock.MockFileData["./fixtures/config/test/spaces.yml"] = config.Spaces{Org: "test", Spaces: []string{"space1", "space2", "space3", "space4"}, EnableDeleteSpaces: false}
	utilsMgrMock.MockFileData["./fixtures/config/test2/orgConfig.yml"] = config.OrgConfig{Org: "test2", BillingManagerGroup: "test2_billing_managers", ManagerGroup: "", AuditorGroup: "", BillingManager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Manager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Auditor: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, PrivateDomains: []string(nil), RemovePrivateDomains: false, EnableOrgQuota: true, MemoryLimit: 10240, InstanceMemoryLimit: -1, TotalRoutes: 10, TotalServices: -1, PaidServicePlansAllowed: true, RemoveUsers: false, TotalPrivateDomains: 0, TotalReservedRoutePorts: 0, TotalServiceKeys: 0, AppInstanceLimit: 0, DefaultIsoSegment: ""}
	utilsMgrMock.MockFileData["./fixtures/config/test2/spaces.yml"] = config.Spaces{Org: "test2", Spaces: []string{"space1a", "space2a", "space3a", "space4a"}, EnableDeleteSpaces: false}
	utilsMgrMock.MockFileData["./fixtures/config-delete/ldap.yml"] = ldap.Config{Enabled: true, LdapHost: "127.0.0.1", LdapPort: 10389, TLS: false, BindDN: "uid=admin,ou=system", BindPassword: "secret", UserSearchBase: "ou=users,dc=example,dc=com", UserNameAttribute: "uid", UserMailAttribute: "mail", UserObjectClass: "", GroupSearchBase: "ou=groups,dc=example,dc=com", GroupAttribute: "member", Origin: ""}
	utilsMgrMock.MockFileData["./fixtures/config-delete/orgs.yml"] = config.Orgs{Orgs: []string{"test"}, EnableDeleteOrgs: true, ProtectedOrgs: []string(nil)}
	utilsMgrMock.MockFileData["./fixtures/config-delete/test/orgConfig.yml"] = config.OrgConfig{Org: "test", BillingManagerGroup: "test_billing_managers", ManagerGroup: "", AuditorGroup: "", BillingManager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Manager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Auditor: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, PrivateDomains: []string(nil), RemovePrivateDomains: false, EnableOrgQuota: true, MemoryLimit: 10240, InstanceMemoryLimit: -1, TotalRoutes: 10, TotalServices: -1, PaidServicePlansAllowed: true, RemoveUsers: false, TotalPrivateDomains: 0, TotalReservedRoutePorts: 0, TotalServiceKeys: 0, AppInstanceLimit: 0, DefaultIsoSegment: ""}
	utilsMgrMock.MockFileData["./fixtures/config-delete/test/spaces.yml"] = config.Spaces{Org: "test", Spaces: []string{"space1", "space2", "space3", "space4"}, EnableDeleteSpaces: false}
	utilsMgrMock.MockFileData["./fixtures/config-private-domains/ldap.yml"] = ldap.Config{Enabled: true, LdapHost: "127.0.0.1", LdapPort: 10389, TLS: false, BindDN: "uid=admin,ou=system", BindPassword: "secret", UserSearchBase: "ou=users,dc=example,dc=com", UserNameAttribute: "uid", UserMailAttribute: "mail", UserObjectClass: "", GroupSearchBase: "ou=groups,dc=example,dc=com", GroupAttribute: "member", Origin: ""}
	utilsMgrMock.MockFileData["./fixtures/config-private-domains/orgs.yml"] = config.Orgs{Orgs: []string{"test", "test2"}, EnableDeleteOrgs: false, ProtectedOrgs: []string(nil)}
	utilsMgrMock.MockFileData["./fixtures/config-private-domains/test/orgConfig.yml"] = config.OrgConfig{Org: "test", BillingManagerGroup: "test_billing_managers", ManagerGroup: "", AuditorGroup: "", BillingManager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Manager: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, Auditor: config.UserMgmt{LDAPUsers: []string(nil), Users: []string(nil), SamlUsers: []string(nil), LDAPGroup: "", LDAPGroups: []string(nil)}, PrivateDomains: []string{"test.com", "test2.com"}, RemovePrivateDomains: true, EnableOrgQuota: true, MemoryLimit: 10240, InstanceMemoryLimit: -1, TotalRoutes: 10, TotalServices: -1, PaidServicePlansAllowed: true, RemoveUsers: false, TotalPrivateDomains: 0, TotalReservedRoutePorts: 0, TotalServiceKeys: 0, AppInstanceLimit: 0, DefaultIsoSegment: ""}
	utilsMgrMock.MockFileData["./fixtures/config-private-domains/test/spaces.yml"] = config.Spaces{Org: "test", Spaces: []string{"space1", "space2", "space3", "space4"}, EnableDeleteSpaces: false}

	return nil
}
