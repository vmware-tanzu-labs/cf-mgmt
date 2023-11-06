package user

import (
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
)

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	UpdateSpaceUsers() []error
	UpdateOrgUsers() []error
	CleanupOrgUsers() []error
}

type LdapManager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByDN(userDN string) (*ldap.User, error)
	GetUserByID(userID string) (*ldap.User, error)
	Close()
}
