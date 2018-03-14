package config

//Config -
type LdapConfig struct {
	Enabled           bool   `yaml:"enabled"`
	LdapHost          string `yaml:"ldapHost"`
	LdapPort          int    `yaml:"ldapPort"`
	TLS               bool   `yaml:"use_tls"`
	BindDN            string `yaml:"bindDN"`
	BindPassword      string `yaml:"bindPwd,omitempty"`
	UserSearchBase    string `yaml:"userSearchBase"`
	UserNameAttribute string `yaml:"userNameAttribute"`
	UserMailAttribute string `yaml:"userMailAttribute"`
	UserObjectClass   string `yaml:"userObjectClass"`
	GroupSearchBase   string `yaml:"groupSearchBase"`
	GroupAttribute    string `yaml:"groupAttribute"`
	Origin            string `yaml:"origin"`
}
