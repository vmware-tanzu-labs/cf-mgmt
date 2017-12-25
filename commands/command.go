package commands

//BaseConfigCommand - commmand that specifies config-dir
type BaseConfigCommand struct {
	ConfigDirectory string `long:"config-dir" env:"CONFIG_DIR" default:"config" description:"Name of the config directory"`
}

//BaseCFConfigCommand - base command that has details to connect to cloud foundry instance
type BaseCFConfigCommand struct {
	BaseConfigCommand
	SystemDomain string `long:"system-domain" env:"SYSTEM_DOMAIN"  description:"system domain"`
	UserID       string `long:"user-id" env:"USER_ID"  description:"user id that has privileges to create/update/delete users, orgs and spaces"`
	Password     string `long:"password" env:"PASSWORD"  description:"password for user account [optional if client secret is provided]"`
	ClientSecret string `long:"client-secret" env:"CLIENT_SECRET" description:"secret for user account that has sufficient privileges to create/update/delete users, orgs and spaces]"`
}

//BaseLDAPCommand - base command that has ldap password
type BaseLDAPCommand struct {
	LdapPassword string `long:"ldap-password" env:"LDAP_PASSWORD"  description:"LDAP password for binding"`
}

//BasePeekCommand - base command for non read-only operations
type BasePeekCommand struct {
	Peek bool `long:"peek" env:"PEEK"  description:"Preview entities to change without modifying"`
}
