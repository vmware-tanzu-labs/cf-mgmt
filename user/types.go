package user

import (
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
	AddUser                                     func(updateUserInput UpdateUsersInput, userName string) error
	RemoveUser                                  func(updateUserInput UpdateUsersInput, userName string) error
}

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	InitializeLdap(ldapBindPassword string) error
	UpdateSpaceUsers() error
	UpdateOrgUsers() error
	ListSpaceAuditors(spaceGUID string) (map[string]string, error)
	ListSpaceDevelopers(spaceGUID string) (map[string]string, error)
	ListSpaceManagers(spaceGUID string) (map[string]string, error)
	ListOrgAuditors(orgGUID string) (map[string]string, error)
	ListOrgBillingManagers(orgGUID string) (map[string]string, error)
	ListOrgManagers(orgGUID string) (map[string]string, error)
}

type CFClient interface {
	RemoveSpaceAuditorByUsername(spaceGUID, userName string) error
	RemoveSpaceDeveloperByUsername(spaceGUID, userName string) error
	RemoveSpaceManagerByUsername(spaceGUID, userName string) error
	ListSpaceAuditors(spaceGUID string) ([]cfclient.User, error)
	ListSpaceManagers(spaceGUID string) ([]cfclient.User, error)
	ListSpaceDevelopers(spaceGUID string) ([]cfclient.User, error)
	AssociateOrgUserByUsername(orgGUID, userName string) (cfclient.Org, error)
	AssociateSpaceAuditorByUsername(spaceGUID, userName string) (cfclient.Space, error)
	AssociateSpaceDeveloperByUsername(spaceGUID, userName string) (cfclient.Space, error)
	AssociateSpaceManagerByUsername(spaceGUID, userName string) (cfclient.Space, error)

	RemoveOrgUserByUsername(orgGUID, name string) error
	RemoveOrgAuditorByUsername(orgGUID, name string) error
	RemoveOrgBillingManagerByUsername(orgGUID, name string) error
	RemoveOrgManagerByUsername(orgGUID, name string) error
	ListOrgAuditors(orgGUID string) ([]cfclient.User, error)
	ListOrgManagers(orgGUID string) ([]cfclient.User, error)
	ListOrgBillingManagers(orgGUID string) ([]cfclient.User, error)
	AssociateOrgAuditorByUsername(orgGUID, name string) (cfclient.Org, error)
	AssociateOrgManagerByUsername(orgGUID, name string) (cfclient.Org, error)
	AssociateOrgBillingManagerByUsername(orgGUID, name string) (cfclient.Org, error)
}
