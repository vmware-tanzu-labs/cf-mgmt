package space

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/ldap"
	"github.com/pivotalservices/cf-mgmt/organization"
	"github.com/pivotalservices/cf-mgmt/uaac"
	"github.com/pivotalservices/cf-mgmt/utils"
)

//Manager -
type Manager interface {
	FindSpace(orgName, spaceName string) (space *cloudcontroller.Space, err error)
	CreateSpaces(configDir, ldapBindPassword string) (err error)
	UpdateSpaces(configDir string) (err error)
	UpdateSpaceUsers(configDir, ldapBindPassword string) (err error)
	CreateQuotas(configDir string) (err error)
	CreateApplicationSecurityGroups(configDir string) (err error)
}

//InputCreateSpaces -
type InputCreateSpaces struct {
	Org    string   `yaml:"org"`
	Spaces []string `yaml:"spaces"`
}

//Contains -
func (s *InputCreateSpaces) Contains(spaceName string) bool {
	set := make(map[string]bool)
	for _, v := range s.Spaces {
		set[v] = true
	}
	_, ok := set[spaceName]
	return ok
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
	} else {
		return i.DeveloperGroup
	}
}

func (i *InputUpdateSpaces) GetManagerGroup() string {
	if i.Manager.LdapGroup != "" {
		return i.Manager.LdapGroup
	} else {
		return i.ManagerGroup
	}
}

func (i *InputUpdateSpaces) GetAuditorGroup() string {
	if i.Auditor.LdapGroup != "" {
		return i.Auditor.LdapGroup
	} else {
		return i.AuditorGroup
	}
}

type UserMgmt struct {
	LdapUser  []string `yaml:"ldap_users"`
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
