package ldap

import (
	l "github.com/go-ldap/ldap"
)

//Manager -
type Manager interface {
	GetUserIDs(config *Config, groupName string) (users []User, err error)
	GetUser(config *Config, userID string) (*User, error)
	GetConfig(configDir, ldapBindPassword string) (*Config, error)
	GetLdapUser(config *Config, userDN string) (*User, error)
	LdapConnection(config *Config) (*l.Conn, error)
}

//DefaultManager -
type DefaultManager struct {
}

//Config -
type Config struct {
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

//User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
