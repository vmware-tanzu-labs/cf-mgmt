package commands

import "github.com/vmwarepivotallabs/cf-mgmt/configcommands"

// BaseCFConfigCommand - base command that has details to connect to cloud foundry instance
type BaseCFConfigCommand struct {
	configcommands.BaseConfigCommand
	SystemDomain string `long:"system-domain" env:"SYSTEM_DOMAIN"  description:"system domain"`
	UserID       string `long:"user-id" env:"USER_ID"  description:"user id that has privileges to create/update/delete users, orgs and spaces"`
	Password     string `long:"password" env:"PASSWORD"  description:"password for user account [optional if client secret is provided]"`
	ClientSecret string `long:"client-secret" env:"CLIENT_SECRET" description:"secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces]"`
}

// BaseLDAPCommand - base command that has ldap password
type BaseLDAPCommand struct {
	LdapServer   string `long:"ldap-server" env:"LDAP_SERVER"  description:"LDAP server for binding"`
	LdapPassword string `long:"ldap-password" env:"LDAP_PASSWORD"  description:"LDAP password for binding"`
	LdapUser     string `long:"ldap-user" env:"LDAP_USER"  description:"LDAP user for binding"`
}

// BaseAzureADCommand - base command that has Azure AD info
type BaseAzureADCommand struct {
	AadTenantId   string `long:"aad-tenantid" env:"AAD_TENANT_ID"  description:"Azure AD Tenant id"`
	AadClientId   string `long:"aad-clientid" env:"AAD_CLIENT_ID"  description:"Azure AD Client Id"`
	AadSecret     string `long:"aad-secret" env:"AAD_SECRET"  description:"Azure AD Client secret"`
	AADUserOrigin string `long:"aad-origin" env:"AAD_ORIGIN"  description:"Azure AD Origin"`
}

// BasePeekCommand - base command for non read-only operations
type BasePeekCommand struct {
	Peek bool `long:"peek" env:"PEEK"  description:"Preview entities to change without modifying"`
}
