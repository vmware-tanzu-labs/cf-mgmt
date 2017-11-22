package commands

import (
	"github.com/pivotalservices/cf-mgmt/configcommands"
)

type CfMgmtCommand struct {
	Version                          configcommands.VersionCommand    `command:"version" description:"Print version information and exit"`
	InitConfigurationCommand         InitConfigurationCommand         `command:"init-config" description:"Initializes folder structure for configuration"`
	AddOrgToConfigurationCommand     AddOrgToConfigurationCommand     `command:"add-org-to-config" description:"Adds specified org to configuration"`
	AddSpaceToConfigurationCommand   AddSpaceToConfigurationCommand   `command:"add-space-to-config" description:"Adds specified space to configuration for org"`
	GenerateConcoursePipelineCommand GenerateConcoursePipelineCommand `command:"generate-concourse-pipeline" description:"generates a concourse pipline to be used to drive cf-mgmt"`
	ExportConfigurationCommand       ExportConfigurationCommand       `command:"export-config" description:"Exports org and space configurations from an existing Cloud Foundry instance. [Warning: This operation will delete existing config folder]"`
	CreateOrgsCommand                CreateOrgsCommand                `command:"create-orgs" description:"creates organizations for each orgConfig.yml"`
	CreateSecurityGroupsCommand      CreateSecurityGroupsCommand      `command:"create-security-groups" description:"creates named security groups that can be assigned to spaces"`
	CreatePrivateDomainsCommand      CreatePrivateDomainsCommand      `command:"create-org-private-domains" description:"creates private domains for an org"`
	DeleteOrgsCommand                DeleteOrgsCommand                `command:"delete-orgs" description:"deletes orgs not in the configuration"`
	UpdateOrgQuotasCommand           UpdateOrgQuotasCommand           `command:"update-org-quotas" description:"updates org quotas"`
	UpdateOrgUsersCommand            UpdateOrgUsersCommand            `command:"update-org-users" description:"update org user roles"`
	CreateSpacesCommand              CreateSpacesCommand              `command:"create-spaces" description:"creates spaces in configuration"`
	DeleteSpacesCommand              DeleteSpacesCommand              `command:"delete-spaces" description:"deletes spaces not in configurtion"`
	UpdateSpacesCommand              UpdateSpacesCommand              `command:"update-spaces" description:"enables/disables ssh access at space level"`
	UpdateSpaceQuotasCommand         UpdateSpaceQuotasCommand         `command:"update-space-quotas" description:"updates spaces quotas"`
	UpdateSpaceUsersCommand          UpdateSpaceUsersCommand          `command:"update-space-users" description:"update space user roles"`
	CreateSpaceSecurityGroupsCommand CreateSpaceSecurityGroupsCommand `command:"update-space-security-groups" description:"updates space specific security groups"`
	IsolationSegmentsCommand         IsolationSegmentsCommand         `command:"isolation-segments" description:"assigns isolations segments to orgs and spaces"`
}

var CfMgmt CfMgmtCommand
