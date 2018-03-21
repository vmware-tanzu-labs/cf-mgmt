package securitygroup

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type Manager interface {
	ListNonDefaultSecurityGroups() (map[string]cfclient.SecGroup, error)
	ListDefaultSecurityGroups() (map[string]cfclient.SecGroup, error)
	ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error)
	GetSecurityGroupRules(sgGUID string) ([]byte, error)
	CreateApplicationSecurityGroups() error
	CreateGlobalSecurityGroups() error
	AssignDefaultSecurityGroups() error
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
