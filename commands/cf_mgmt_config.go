package commands

type CfMgmtConfigCommand struct {
	Version                          VersionCommand                   `command:"version" description:"Print version information and exit"`
	InitConfigurationCommand         InitConfigurationCommand         `command:"init" description:"Initializes folder structure for configuration"`
	AddOrgToConfigurationCommand     AddOrgToConfigurationCommand     `command:"add-org" description:"Adds specified org to configuration"`
	AddSpaceToConfigurationCommand   AddSpaceToConfigurationCommand   `command:"add-space" description:"Adds specified space to configuration for org"`
	ExportConfigurationCommand       ExportConfigurationCommand       `command:"export" description:"Exports org and space configurations from an existing Cloud Foundry instance. [Warning: This operation will delete existing config folder]"`
	GenerateConcoursePipelineCommand GenerateConcoursePipelineCommand `command:"generate-concourse-pipeline" description:"generates a concourse pipline to be used to drive cf-mgmt"`
}

var CfMgmtConfig CfMgmtConfigCommand
