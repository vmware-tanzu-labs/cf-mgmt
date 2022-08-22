package configcommands

type CfMgmtConfigCommand struct {
	Version                             VersionCommand                      `command:"version" description:"Print version information and exit"`
	InitConfigurationCommand            InitConfigurationCommand            `command:"init" description:"Initializes folder structure for configuration"`
	GlobalConfigurationCommand          GlobalConfigurationCommand          `command:"global" description:"Updates values in cf-mgmt.yml"`
	OrgConfigurationCommand             OrgConfigurationCommand             `command:"org" description:"Adds/updates specified org to configuration"`
	SpaceConfigurationCommand           SpaceConfigurationCommand           `command:"space" description:"adds/updates space configuration"`
	AddOrgToConfigurationCommand        AddOrgToConfigurationCommand        `command:"add-org" description:"*****DEPRECATED, use 'cf-mgmt-config org instead' - Adds specified org to configuration"`
	AddSpaceToConfigurationCommand      AddSpaceToConfigurationCommand      `command:"add-space" description:"*****DEPRECATED, use 'cf-mgmt-config space instead' - Adds specified space to configuration for org"`
	GenerateConcoursePipelineCommand    GenerateConcoursePipelineCommand    `command:"generate-concourse-pipeline" description:"generates a concourse pipeline to be used to drive cf-mgmt"`
	UpdateOrgConfigurationCommand       UpdateOrgConfigurationCommand       `command:"update-org" description:"*****DEPRECATED, use 'cf-mgmt-config org instead' - updates org configuration"`
	UpdateSpaceConfigurationCommand     UpdateSpaceConfigurationCommand     `command:"update-space" description:"*****DEPRECATED, use 'cf-mgmt-config space instead' - updates space configuration"`
	DeleteOrgConfigurationCommand       DeleteOrgConfigurationCommand       `command:"delete-org" description:"deletes org configuration"`
	DeleteSpaceConfigurationCommand     DeleteSpaceConfigurationCommand     `command:"delete-space" description:"deletes space configuration"`
	AddASGToConfigurationCommand        AddASGToConfigurationCommand        `command:"add-asg" description:"*****DEPRECATED, use 'cf-mgmt-config asg instead' - add a named asg to configuration"`
	ASGToConfigurationCommand           ASGToConfigurationCommand           `command:"asg" description:"creates/updates a named asg"`
	UpdateOrgsConfigurationCommand      UpdateOrgsConfigurationCommand      `command:"update-orgs" description:"updates orgs.yml"`
	RenameOrgConfigurationCommand       RenameOrgConfigurationCommand       `command:"rename-org" description:"renames an org"`
	RenameSpaceConfigurationCommand     RenameSpaceConfigurationCommand     `command:"rename-space" description:"renames a space for a given org"`
	OrgNamedQuotaConfigurationCommand   OrgNamedQuotaConfigurationCommand   `command:"named-org-quota" description:"creates/updates named org quota"`
	SpaceNamedQuotaConfigurationCommand SpaceNamedQuotaConfigurationCommand `command:"named-space-quota" description:"creates/updates named space quota"`
	ClearUsersCommand                   ClearUsersCommand                   `command:"clear-users" description:"updates all configuration but removes any user/group mapping"`
}

var CfMgmtConfig CfMgmtConfigCommand
