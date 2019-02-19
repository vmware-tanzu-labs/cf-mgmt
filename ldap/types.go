package ldap

import (
	"github.com/pivotalservices/cf-mgmt/config"
)

//Manager -
type Manager struct {
	Config     *config.LdapConfig
	Connection Connection
}

//User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
