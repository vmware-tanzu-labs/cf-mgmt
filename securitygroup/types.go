package securitygroup

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Resource -
type Resource struct {
	MetaData MetaData `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

//MetaData -
type MetaData struct {
	GUID string `json:"guid"`
}

//Entity -
type Entity struct {
	Name    string `json:"name"`
	Rules   []Rule `json:"rules"`
	Running bool   `json:"running_default"`
	Staging bool   `json:"staging_default"`
}

//Rule -
type Rule struct {
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	Ports       string `json:"ports"`
}

//DefaultSecurityGroupManager -
type DefaultSecurityGroupManager struct {
	Cfg             config.Reader
	FilePattern     string
	FilePaths       []string
	CloudController cloudcontroller.Manager
	//UAACMgr         uaac.Manager
	//OrgMgr          organization.Manager
	//LdapMgr         ldap.Manager
	UtilsMgr utils.Manager
	//	UserMgr         UserMgr
}

type Manager interface {
	//FindSpace(orgName, spaceName string) (*cloudcontroller.Space, error)
	//CreateSpaces(configDir, ldapBindPassword string) error
	//	UpdateSpaces(configDir string) (err error)
	//	UpdateSpaceUsers(configDir, ldapBindPassword string) error
	//	CreateQuotas(configDir string) error
	CreateApplicationSecurityGroups(configDir string) error
	//DeleteSpaces(configFile string, peekDeletion bool) (err error)
}

/*

package space

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
)

const SPACES = "spaces"
const ROLE_SPACE_DEVELOPERS = "developers"
const ROLE_SPACE_MANAGERS = "managers"
const ROLE_SPACE_AUDITORS = "auditors"

//Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (*cloudcontroller.Space, error)
	CreateSpaces(configDir, ldapBindPassword string) error
	UpdateSpaces(configDir string) (err error)
	UpdateSpaceUsers(configDir, ldapBindPassword string) error
	CreateQuotas(configDir string) error
	CreateApplicationSecurityGroups(configDir string) error
	DeleteSpaces(configFile string, peekDeletion bool) (err error)
}

//DefaultSpaceManager -
type DefaultSpaceManager struct {
	Cfg             config.Reader
	FilePattern     string
	FilePaths       []string
	CloudController cloudcontroller.Manager
	UAACMgr         uaac.Manager
	OrgMgr          organization.Manager
	LdapMgr         ldap.Manager
	UtilsMgr        utils.Manager
	UserMgr         UserMgr
}
*/
