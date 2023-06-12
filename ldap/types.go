package ldap

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

// Manager -
type Manager struct {
	Config     *config.LdapConfig
	Connection Connection
	groupMap   map[string][]string
	userMap    map[string]*User
}

// User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
