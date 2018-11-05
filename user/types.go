package user

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/uaa"
)

// UpdateSpaceUserInput
type UpdateUsersInput struct {
	SpaceGUID                                   string
	OrgGUID                                     string
	LdapUsers, Users, LdapGroupNames, SamlUsers []string
	SpaceName                                   string
	OrgName                                     string
	RemoveUsers                                 bool
	ListUsers                                   func(updateUserInput UpdateUsersInput, uaaUsers *uaa.Users) (*RoleUsers, error)
	AddUser                                     func(updateUserInput UpdateUsersInput, userName, origin string) error
	RemoveUser                                  func(updateUserInput UpdateUsersInput, userName, origin string) error
}

type RoleUsers struct {
	users map[string]map[string]RoleUser
}
type RoleUser struct {
	UserName string
	//GUID     string
	Origin string
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
	RemoveSpaceAuditorByUsernameAndOrigin(spaceGUID, userName, origin string) error
	RemoveSpaceDeveloperByUsernameAndOrigin(spaceGUID, userName, origin string) error
	RemoveSpaceManagerByUsernameAndOrigin(spaceGUID, userName, origin string) error
	ListSpaceAuditors(spaceGUID string) ([]cfclient.User, error)
	ListSpaceManagers(spaceGUID string) ([]cfclient.User, error)
	ListSpaceDevelopers(spaceGUID string) ([]cfclient.User, error)
	AssociateOrgUserByUsernameAndOrigin(orgGUID, userName, origin string) (cfclient.Org, error)
	AssociateSpaceAuditorByUsernameAndOrigin(spaceGUID, userName, origin string) (cfclient.Space, error)
	AssociateSpaceDeveloperByUsernameAndOrigin(spaceGUID, userName, origin string) (cfclient.Space, error)
	AssociateSpaceManagerByUsernameAndOrigin(spaceGUID, userName, origin string) (cfclient.Space, error)

	RemoveOrgUserByUsernameAndOrigin(orgGUID, name, origin string) error
	RemoveOrgAuditorByUsernameAndOrigin(orgGUID, name, origin string) error
	RemoveOrgBillingManagerByUsernameAndOrigin(orgGUID, name, origin string) error
	RemoveOrgManagerByUsernameAndOrigin(orgGUID, name, origin string) error
	ListOrgAuditors(orgGUID string) ([]cfclient.User, error)
	ListOrgManagers(orgGUID string) ([]cfclient.User, error)
	ListOrgBillingManagers(orgGUID string) ([]cfclient.User, error)
	AssociateOrgAuditorByUsernameAndOrigin(orgGUID, name, origin string) (cfclient.Org, error)
	AssociateOrgManagerByUsernameAndOrigin(orgGUID, name, origin string) (cfclient.Org, error)
	AssociateOrgBillingManagerByUsernameAndOrigin(orgGUID, name, origin string) (cfclient.Org, error)

	ListOrgUsers(orgGUID string) ([]cfclient.User, error)
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)
}
