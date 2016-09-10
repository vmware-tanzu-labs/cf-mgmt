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
	CreateSpace(orgName, spaceName string) (space cloudcontroller.Space, err error)
	FindSpace(orgName, spaceName string) (space cloudcontroller.Space, err error)
	CreateSpaces(configDir string) (err error)
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
	Org                     string `yaml:"org"`
	Space                   string `yaml:"space"`
	DeveloperGroup          string `yaml:"space-developer-group"`
	ManagerGroup            string `yaml:"space-manager-group"`
	AuditorGroup            string `yaml:"space-auditor-group"`
	AllowSSH                bool   `yaml:"allow-ssh"`
	EnableSpaceQuota        bool   `yaml:"enable-space-quota"`
	MemoryLimit             int    `yaml:"memory-limit"`
	InstanceMemoryLimit     int    `yaml:"instance-memory-limit"`
	TotalRoutes             int    `yaml:"total-routes"`
	TotalServices           int    `yaml:"total-services"`
	PaidServicePlansAllowed bool   `yaml:"paid-service-plans-allowed"`
	EnableSecurityGroup     bool   `yaml:"enable-security-group"`
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
}
