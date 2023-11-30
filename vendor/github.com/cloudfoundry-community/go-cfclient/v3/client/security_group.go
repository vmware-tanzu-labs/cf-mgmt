package client

import (
	"context"
	"net/url"

	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

type SecurityGroupClient commonClient

// SecurityGroupListOptions list filters
type SecurityGroupListOptions struct {
	*ListOptions

	GUIDs             Filter `qs:"guids"`               // list of security group guids to filter by
	Names             Filter `qs:"names"`               // list of security group names to filter by
	RunningSpaceGUIDs Filter `qs:"running_space_guids"` // list of space guids to filter by
	StagingSpaceGUIDs Filter `qs:"staging_space_guids"` // list of space guids to filter by

	GloballyEnabledRunning *bool `qs:"globally_enabled_running"` // If true, only include the security groups that are enabled for running
	GloballyEnabledStaging *bool `qs:"globally_enabled_staging"` // If true, only include the security groups that are enabled for staging
}

// NewSecurityGroupListOptions creates new options to pass to list
func NewSecurityGroupListOptions() *SecurityGroupListOptions {
	return &SecurityGroupListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SecurityGroupListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// SecurityGroupSpaceListOptions list filters
type SecurityGroupSpaceListOptions struct {
	*ListOptions

	GUIDs Filter `qs:"guids"` // list of security group guids to filter by
	Names Filter `qs:"names"` // list of security group names to filter by
}

// NewSecurityGroupSpaceListOptions creates new options to pass to list
func NewSecurityGroupSpaceListOptions() *SecurityGroupSpaceListOptions {
	return &SecurityGroupSpaceListOptions{
		ListOptions: NewListOptions(),
	}
}

func (o SecurityGroupSpaceListOptions) ToQueryString() (url.Values, error) {
	return o.ListOptions.ToQueryString(o)
}

// BindRunningSecurityGroup binds one or more spaces to a security group with the running lifecycle and returns
// the space GUIDs bound to the security group
//
// Running app containers within these spaces will inherit the rules specified by this security group. Apps within
// these spaces must be restarted for these changes to take effect. Unless a security group is globally-enabled,
// an admin must add it to a space for it to be visible for the org and space managers. Once it’s visible, org and
// space managers can add it to additional spaces.
func (c *SecurityGroupClient) BindRunningSecurityGroup(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/security_groups/%s/relationships/running_spaces", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// BindStagingSecurityGroup binds one or more spaces to a security group with the staging lifecycle and returns
// the space GUIDs bound to the security group
//
// Staging app containers within these spaces will inherit the rules specified by this security group. Apps within
// these spaces must be restaged for these changes to take effect. Unless a security group is globally-enabled,
// an admin must add it to a space for it to be visible for the org and space managers. Once it’s visible, org and
// space managers can add it to additional spaces.
func (c *SecurityGroupClient) BindStagingSecurityGroup(ctx context.Context, guid string, spaceGUIDs []string) ([]string, error) {
	req := resource.NewToManyRelationships(spaceGUIDs)
	var relation resource.ToManyRelationships
	_, err := c.client.post(ctx, path.Format("/v3/security_groups/%s/relationships/staging_spaces", guid), req, &relation)
	if err != nil {
		return nil, err
	}
	var guids []string
	for _, r := range relation.Data {
		guids = append(guids, r.GUID)
	}
	return guids, nil
}

// Create a new domain
func (c *SecurityGroupClient) Create(ctx context.Context, r *resource.SecurityGroupCreate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.post(ctx, "/v3/security_groups", r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// Delete the specified security group asynchronously and return a jobGUID
func (c *SecurityGroupClient) Delete(ctx context.Context, guid string) (string, error) {
	return c.client.delete(ctx, path.Format("/v3/security_groups/%s", guid))
}

// First returns the first security group matching the options or an error when less than 1 match
func (c *SecurityGroupClient) First(ctx context.Context, opts *SecurityGroupListOptions) (*resource.SecurityGroup, error) {
	return First[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Get the specified security group
func (c *SecurityGroupClient) Get(ctx context.Context, guid string) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	err := c.client.get(ctx, path.Format("/v3/security_groups/%s", guid), &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// List pages SecurityGroups the user has access to
func (c *SecurityGroupClient) List(ctx context.Context, opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
	if opts == nil {
		opts = NewSecurityGroupListOptions()
	}
	var res resource.SecurityGroupList
	err := c.client.list(ctx, "/v3/security_groups", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListAll retrieves all SecurityGroups the user has access to
func (c *SecurityGroupClient) ListAll(ctx context.Context, opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, error) {
	if opts == nil {
		opts = NewSecurityGroupListOptions()
	}
	return AutoPage[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// Single returns a single security group matching the options or an error if not exactly 1 match
func (c *SecurityGroupClient) Single(ctx context.Context, opts *SecurityGroupListOptions) (*resource.SecurityGroup, error) {
	return Single[*SecurityGroupListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.List(ctx, opts)
	})
}

// ListRunningForSpace pages security groups that are enabled for running globally or at the space level for the given space
func (c *SecurityGroupClient) ListRunningForSpace(ctx context.Context, spaceGUID string, opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, *Pager, error) {
	if opts == nil {
		opts = NewSecurityGroupSpaceListOptions()
	}
	var res resource.SecurityGroupList
	err := c.client.list(ctx, "/v3/spaces/"+spaceGUID+"/running_security_groups", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListRunningForSpaceAll retrieves all security groups that are enabled for running globally or at the space level for the given space
func (c *SecurityGroupClient) ListRunningForSpaceAll(ctx context.Context, spaceGUID string, opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, error) {
	if opts == nil {
		opts = NewSecurityGroupSpaceListOptions()
	}
	return AutoPage[*SecurityGroupSpaceListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.ListRunningForSpace(ctx, spaceGUID, opts)
	})
}

// ListStagingForSpace pages security groups that are enabled for staging globally or at the space level for the given space
func (c *SecurityGroupClient) ListStagingForSpace(ctx context.Context, spaceGUID string, opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, *Pager, error) {
	if opts == nil {
		opts = NewSecurityGroupSpaceListOptions()
	}
	var res resource.SecurityGroupList
	err := c.client.list(ctx, "/v3/spaces/"+spaceGUID+"/staging_security_groups", opts.ToQueryString, &res)
	if err != nil {
		return nil, nil, err
	}
	pager := NewPager(res.Pagination)
	return res.Resources, pager, nil
}

// ListStagingForSpaceAll retrieves all security groups that are enabled for staging globally or at the space level for the given space
func (c *SecurityGroupClient) ListStagingForSpaceAll(ctx context.Context, spaceGUID string, opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, error) {
	if opts == nil {
		opts = NewSecurityGroupSpaceListOptions()
	}
	return AutoPage[*SecurityGroupSpaceListOptions, *resource.SecurityGroup](opts, func(opts *SecurityGroupSpaceListOptions) ([]*resource.SecurityGroup, *Pager, error) {
		return c.ListStagingForSpace(ctx, spaceGUID, opts)
	})
}

// UnBindRunningSecurityGroup removes a space from a security group with the running lifecycle
//
// Apps within this space must be restarted for these changes to take effect.
func (c *SecurityGroupClient) UnBindRunningSecurityGroup(ctx context.Context, guid string, spaceGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/security_groups/%s/relationships/running_spaces/%s", guid, spaceGUID))
	return err
}

// UnBindStagingSecurityGroup removes a space from a security group with the staging lifecycle
//
// Apps within this space must be restarted for these changes to take effect.
func (c *SecurityGroupClient) UnBindStagingSecurityGroup(ctx context.Context, guid string, spaceGUID string) error {
	_, err := c.client.delete(ctx, path.Format("/v3/security_groups/%s/relationships/staging_spaces/%s", guid, spaceGUID))
	return err
}

// Update the specified attributes of the app
func (c *SecurityGroupClient) Update(ctx context.Context, guid string, r *resource.SecurityGroupUpdate) (*resource.SecurityGroup, error) {
	var d resource.SecurityGroup
	_, err := c.client.patch(ctx, path.Format("/v3/security_groups/%s", guid), r, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
