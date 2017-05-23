package space

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
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
}

//InputCreateSpaces -
type InputCreateSpaces struct {
	Org    string   `yaml:"org"`
	Spaces []string `yaml:"spaces"`
}

//Contains -
func (s *InputCreateSpaces) Contains(spaceName string) bool {
	for _, v := range s.Spaces {
		if v == spaceName {
			return true
		}
	}
	return false
}

//InputUpdateSpaces -
type InputUpdateSpaces struct {
	Org                     string   `yaml:"org"`
	Space                   string   `yaml:"space"`
	Developer               UserMgmt `yaml:"space-developer"`
	Manager                 UserMgmt `yaml:"space-manager"`
	Auditor                 UserMgmt `yaml:"space-auditor"`
	DeveloperGroup          string   `yaml:"space-developer-group,omitempty"`
	ManagerGroup            string   `yaml:"space-manager-group,omitempty"`
	AuditorGroup            string   `yaml:"space-auditor-group,omitempty"`
	AllowSSH                bool     `yaml:"allow-ssh"`
	EnableSpaceQuota        bool     `yaml:"enable-space-quota"`
	MemoryLimit             int      `yaml:"memory-limit"`
	InstanceMemoryLimit     int      `yaml:"instance-memory-limit"`
	TotalRoutes             int      `yaml:"total-routes"`
	TotalServices           int      `yaml:"total-services"`
	PaidServicePlansAllowed bool     `yaml:"paid-service-plans-allowed"`
	EnableSecurityGroup     bool     `yaml:"enable-security-group"`
	SecurityGroupContents   string   `yaml:"security-group-contents,omitempty"`
	RemoveUsers             bool     `yaml:"enable-remove-users"`
}

//ConfigSpaceDefaults -
type ConfigSpaceDefaults struct {
	Developer UserMgmt `yaml:"space-developer"`
	Manager   UserMgmt `yaml:"space-manager"`
	Auditor   UserMgmt `yaml:"space-auditor"`
}

func (i *InputUpdateSpaces) GetDeveloperGroup() string {
	if i.Developer.LdapGroup != "" {
		return i.Developer.LdapGroup
	}
	return i.DeveloperGroup
}

func (i *InputUpdateSpaces) GetManagerGroup() string {
	if i.Manager.LdapGroup != "" {
		return i.Manager.LdapGroup
	}
	return i.ManagerGroup
}

func (i *InputUpdateSpaces) GetAuditorGroup() string {
	if i.Auditor.LdapGroup != "" {
		return i.Auditor.LdapGroup
	}
	return i.AuditorGroup
}

type UserMgmt struct {
	LdapUsers []string `yaml:"ldap_users"`
	Users     []string `yaml:"users"`
	LdapGroup string   `yaml:"ldap_group"`
}

//DefaultSpaceManager -
type DefaultSpaceManager struct {
	FilePattern     string
	FilePaths       []string
	CloudController cloudcontroller.Manager
	UAACMgr         uaac.Manager
	OrgMgr          organization.Manager
	LdapMgr         ldap.Manager
	UtilsMgr        utils.Manager
	UserMgr         UserMgr
}
