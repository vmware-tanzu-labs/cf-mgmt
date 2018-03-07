package space

import (
	"net/url"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

//Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (cfclient.Space, error)
	CreateSpaces(configDir, ldapBindPassword string) error
	UpdateSpaces(configDir string) (err error)
	UpdateSpaceUsers(configDir, ldapBindPassword string) error
	CreateQuotas(configDir string) error
	CreateApplicationSecurityGroups(configDir string) error
	DeleteSpaces(configFile string) (err error)
	ListSpaces(orgGUID string) ([]cfclient.Space, error)
	ListSpaceAuditors(spaceGUID string) (map[string]string, error)
	ListSpaceDevelopers(spaceGUID string) (map[string]string, error)
	ListSpaceManagers(spaceGUID string) (map[string]string, error)
}

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

// UserMgr - interface type encapsulating Update space users behavior
type UserMgr interface {
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
	ListOrgSpaceQuotas(orgGUID string) ([]cfclient.SpaceQuota, error)
	UpdateSpaceQuota(spaceQuotaGUID string, spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	AssignSpaceQuota(quotaGUID, spaceGUID string) error
	CreateSpaceQuota(spaceQuote cfclient.SpaceQuotaRequest) (*cfclient.SpaceQuota, error)
	GetSpaceByGuid(spaceGUID string) (cfclient.Space, error)
	UpdateSpace(spaceGUID string, req cfclient.SpaceRequest) (cfclient.Space, error)
	ListSpacesByQuery(query url.Values) ([]cfclient.Space, error)
	CreateSpace(req cfclient.SpaceRequest) (cfclient.Space, error)
	DeleteSpace(guid string, recursive, async bool) error
}
