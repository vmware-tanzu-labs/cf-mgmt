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
	DeleteSpaces(configFile string, peekDeletion bool) (err error)
}

//InputSpaces -
type InputSpaces struct {
	Org                string   `yaml:"org"`
	Spaces             []string `yaml:"spaces"`
	EnableDeleteSpaces bool     `yaml:"enable-delete-space"`
}

//Contains -
func (s *InputSpaces) Contains(spaceName string) bool {
	for _, v := range s.Spaces {
		if v == spaceName {
			return true
		}
	}
	return false
}

//InputSpaceConfig -
type InputSpaceConfig struct {
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

func (i *InputSpaceConfig) GetDeveloperGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Developer.LdapGroups {
		groupMap[group] = group
	}
	if i.Developer.LdapGroup != "" {
		groupMap[i.Developer.LdapGroup] = i.Developer.LdapGroup
	}
	if i.DeveloperGroup != "" {
		groupMap[i.DeveloperGroup] = i.DeveloperGroup
	}
	return mapToKeys(groupMap)
}

func (i *InputSpaceConfig) GetManagerGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Manager.LdapGroups {
		groupMap[group] = group
	}
	if i.Manager.LdapGroup != "" {
		groupMap[i.Manager.LdapGroup] = i.Manager.LdapGroup
	}
	if i.ManagerGroup != "" {
		groupMap[i.ManagerGroup] = i.ManagerGroup
	}
	return mapToKeys(groupMap)
}

func (i *InputSpaceConfig) GetAuditorGroups() []string {
	groupMap := make(map[string]string)
	for _, group := range i.Auditor.LdapGroups {
		groupMap[group] = group
	}
	if i.Auditor.LdapGroup != "" {
		groupMap[i.Auditor.LdapGroup] = i.Auditor.LdapGroup
	}
	if i.AuditorGroup != "" {
		groupMap[i.AuditorGroup] = i.AuditorGroup
	}
	return mapToKeys(groupMap)
}

func mapToKeys(aMap map[string]string) []string {
	var keys []string
	for k := range aMap {
		keys = append(keys, k)
	}
	return keys
}

type UserMgmt struct {
	LdapUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	LdapGroup  string   `yaml:"ldap_group"`
	LdapGroups []string `yaml:"ldap_groups"`
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
