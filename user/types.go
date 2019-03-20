package user

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/uaa"
)

func InitRoleUsers() *RoleUsers {
	return &RoleUsers{
		users: make(map[string][]RoleUser),
	}
}

type RoleUsers struct {
	users map[string][]RoleUser
}
type RoleUser struct {
	UserName string
	GUID     string
	Origin   string
}

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	InitializeLdap(ldapBindPassword string) error
	DeinitializeLdap() error
	UpdateSpaceUsers() error
	UpdateOrgUsers() error
	CleanupOrgUsers() error
	ListSpaceAuditors(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListSpaceDevelopers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListSpaceManagers(spaceGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgAuditors(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgBillingManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
	ListOrgManagers(orgGUID string, uaaUsers *uaa.Users) (*RoleUsers, error)
}

type CFClient interface {
	ListSpaceAuditors(spaceGUID string) ([]cfclient.User, error)
	ListSpaceManagers(spaceGUID string) ([]cfclient.User, error)
	ListSpaceDevelopers(spaceGUID string) ([]cfclient.User, error)
	ListOrgAuditors(orgGUID string) ([]cfclient.User, error)
	ListOrgManagers(orgGUID string) ([]cfclient.User, error)
	ListOrgBillingManagers(orgGUID string) ([]cfclient.User, error)
	ListOrgUsers(orgGUID string) ([]cfclient.User, error)
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)

	RemoveSpaceAuditor(spaceGUID, userGUID string) error
	RemoveSpaceDeveloper(spaceGUID, userGUID string) error
	RemoveSpaceManager(spaceGUID, userGUID string) error
	AssociateOrgUser(orgGUID, userGUID string) (cfclient.Org, error)
	AssociateSpaceAuditor(spaceGUID, userGUID string) (cfclient.Space, error)
	AssociateSpaceDeveloper(spaceGUID, userGUID string) (cfclient.Space, error)
	AssociateSpaceManager(spaceGUID, userGUID string) (cfclient.Space, error)
	RemoveOrgUser(orgGUID, userGUID string) error
	RemoveOrgAuditor(orgGUID, userGUID string) error
	RemoveOrgBillingManager(orgGUID, userGUID string) error
	RemoveOrgManager(orgGUID, userGUID string) error
	AssociateOrgAuditor(orgGUID, userGUID string) (cfclient.Org, error)
	AssociateOrgManager(orgGUID, userGUID string) (cfclient.Org, error)
	AssociateOrgBillingManager(orgGUID, userGUID string) (cfclient.Org, error)
	DeleteUser(userGuid string) error
}

type LdapManager interface {
	GetUserDNs(groupName string) ([]string, error)
	GetUserByDN(userDN string) (*ldap.User, error)
	GetUserByID(userID string) (*ldap.User, error)
	Close()
}
