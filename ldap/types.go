package ldap

import (
	"github.com/pivotalservices/cf-mgmt/config"
)

//Manager -
type Manager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByID(userID string) (*User, error)
	GetUserByDN(userDN string) (*User, error)
	Close()
}

//DefaultManager -
type DefaultManager struct {
	Config     *config.LdapConfig
	Connection Connection
}

//User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
