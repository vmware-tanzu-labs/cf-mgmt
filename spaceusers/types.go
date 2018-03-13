package spaceusers

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
	ListUsers                                   func(spaceGUID string) (map[string]string, error)
	AddUser                                     func(orgGUID, spaceGUID, userName string) error
	RemoveUser                                  func(spaceGUID, userName string) error
}

// Manager - interface type encapsulating Update space users behavior
type Manager interface {
	UpdateSpaceUsers(configDir, ldapBindPassword string) error
	RemoveSpaceAuditorByUsername(spaceGUID, userName string) error
	RemoveSpaceDeveloperByUsername(spaceGUID, userName string) error
	RemoveSpaceManagerByUsername(spaceGUID, userName string) error
	ListSpaceAuditors(spaceGUID string) (map[string]string, error)
	ListSpaceDevelopers(spaceGUID string) (map[string]string, error)
	ListSpaceManagers(spaceGUID string) (map[string]string, error)
	AssociateSpaceAuditorByUsername(orgGUID, spaceGUID, userName string) error
	AssociateSpaceDeveloperByUsername(orgGUID, spaceGUID, userName string) error
	AssociateSpaceManagerByUsername(orgGUID, spaceGUID, userName string) error
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
}
