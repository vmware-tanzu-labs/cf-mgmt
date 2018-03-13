package securitygroup

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
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

type Manager interface {
	CreateApplicationSecurityGroups(configDir string) error
	CreateGlobalSecurityGroups() error
	AssignDefaultSecurityGroups() error
	ListNonDefaultSecurityGroups() (map[string]cfclient.SecGroup, error)
	ListDefaultSecurityGroups() (map[string]cfclient.SecGroup, error)
	AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error
	UpdateSecurityGroup(sgGUID, sgName, contents string) error
	CreateSecurityGroup(sgName, contents string) (*cfclient.SecGroup, error)
	ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error)
	GetSecurityGroupRules(sgGUID string) ([]byte, error)
}

type CFClient interface {
	ListSecGroups() ([]cfclient.SecGroup, error)
	CreateSecGroup(name string, rules []cfclient.SecGroupRule, spaceGuids []string) (*cfclient.SecGroup, error)
	UpdateSecGroup(guid, name string, rules []cfclient.SecGroupRule, spaceGuids []string) (*cfclient.SecGroup, error)
	BindSecGroup(secGUID, spaceGUID string) error
	BindRunningSecGroup(secGUID string) error
	BindStagingSecGroup(secGUID string) error
	UnbindRunningSecGroup(secGUID string) error
	UnbindStagingSecGroup(secGUID string) error
	GetSecGroup(guid string) (*cfclient.SecGroup, error)
	ListSpaceSecGroups(spaceGUID string) ([]cfclient.SecGroup, error)
}
