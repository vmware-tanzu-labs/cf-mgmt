package user

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

// UpdateSpaceUserInput
type UpdateUsersInput struct {
	SpaceGUID                                   string
	OrgGUID                                     string
	LdapUsers, Users, LdapGroupNames, SamlUsers []string
	SpaceName                                   string
	OrgName                                     string
	RemoveUsers                                 bool
	ListUsers                                   func(updateUserInput UpdateUsersInput) (map[string]string, error)
	AddUser                                     func(updateUserInput UpdateUsersInput, userName, origin string) error
	RemoveUser                                  func(updateUserInput UpdateUsersInput, userName, origin string) error
}

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	InitializeLdap(ldapBindPassword string) error
	DeinitializeLdap() error
	UpdateSpaceUsers() error
	UpdateOrgUsers() error
	CleanupOrgUsers() error
	ListSpaceAuditors(spaceGUID string) (map[string]string, error)
	ListSpaceDevelopers(spaceGUID string) (map[string]string, error)
	ListSpaceManagers(spaceGUID string) (map[string]string, error)
	ListOrgAuditors(orgGUID string) (map[string]string, error)
	ListOrgBillingManagers(orgGUID string) (map[string]string, error)
	ListOrgManagers(orgGUID string) (map[string]string, error)
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
