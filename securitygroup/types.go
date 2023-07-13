package securitygroup

import (
	"context"

	cfclient "github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type Manager interface {
	ListNonDefaultSecurityGroups() (map[string]*resource.SecurityGroup, error)
	ListDefaultSecurityGroups() (map[string]*resource.SecurityGroup, error)
	ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error)
	GetSecurityGroupRules(sgGUID string) ([]byte, error)
	CreateApplicationSecurityGroups() error
	CreateGlobalSecurityGroups() error
	AssignDefaultSecurityGroups() error
}

type CFSecurityGroupClient interface {
	ListAll(ctx context.Context, opts *cfclient.SecurityGroupListOptions) ([]*resource.SecurityGroup, error)
	Create(ctx context.Context, r *resource.SecurityGroupCreate) (*resource.SecurityGroup, error)
	Update(ctx context.Context, guid string, r *resource.SecurityGroupUpdate) (*resource.SecurityGroup, error)
	BindRunningSecurityGroup(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error)
	BindStagingSecurityGroup(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error)
	UnBindRunningSecurityGroup(ctx context.Context, guid string, spaceGUID string) error
	UnBindStagingSecurityGroup(ctx context.Context, guid string, spaceGUID string) error
	Get(ctx context.Context, guid string) (*resource.SecurityGroup, error)
	ListRunningForSpaceAll(ctx context.Context, spaceGUID string, opts *cfclient.SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, error)
}
