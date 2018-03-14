package ldap

import (
	l "github.com/go-ldap/ldap"
	"github.com/pivotalservices/cf-mgmt/config"
)

//Manager -
type Manager interface {
	GetUserIDs(groupName string) (users []User, err error)
	GetUser(userID string) (*User, error)
	GetLdapUser(userDN string) (*User, error)
	LdapConnection() (*l.Conn, error)
	GetLdapUsers(groupNames []string, userList []string) ([]User, error)
	LdapConfig() *config.LdapConfig
}

//DefaultManager -
type DefaultManager struct {
	Config *config.LdapConfig
}

//User -
type User struct {
	UserDN string
	UserID string
	Email  string
}
