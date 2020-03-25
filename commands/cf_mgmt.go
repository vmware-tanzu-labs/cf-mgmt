package commands

import (
	"github.com/vmwarepivotallabs/cf-mgmt/configcommands"
)

type CfMgmtCommand struct {
	Version                          configcommands.VersionCommand    `command:"version" description:"Print version information and exit"`
	CreateOrgsCommand                CreateOrgsCommand                `command:"create-orgs" description:"creates organizations for each orgConfig.yml"`
	CreateSecurityGroupsCommand      CreateSecurityGroupsCommand      `command:"create-security-groups" description:"creates named security groups that can be assigned to spaces"`
	AssignDefaultSecurityGroups      AssignDefaultSecurityGroups      `command:"assign-default-security-groups" description:"assigns security groups to default running or default staging"`
	CreatePrivateDomainsCommand      CreatePrivateDomainsCommand      `command:"create-org-private-domains" description:"creates private domains for an org"`
	DeleteOrgsCommand                DeleteOrgsCommand                `command:"delete-orgs" description:"deletes orgs not in the configuration"`
	UpdateOrgQuotasCommand           UpdateOrgQuotasCommand           `command:"update-org-quotas" description:"updates org quotas"`
	UpdateOrgUsersCommand            UpdateOrgUsersCommand            `command:"update-org-users" description:"update org user roles"`
	CleanupOrgUsersCommand           CleanupOrgUsersCommand           `command:"cleanup-org-users" description:"removes any users from org that don't have a role"`
	CreateSpacesCommand              CreateSpacesCommand              `command:"create-spaces" description:"creates spaces in configuration"`
	DeleteSpacesCommand              DeleteSpacesCommand              `command:"delete-spaces" description:"deletes spaces not in configurtion"`
	UpdateSpacesCommand              UpdateSpacesCommand              `command:"update-spaces" description:"enables/disables ssh access at space level"`
	UpdateSpacesMetadataCommand      UpdateSpacesMetadataCommand      `command:"update-spaces-metadata" description:"adds metadata for a space"`
	UpdateSpaceQuotasCommand         UpdateSpaceQuotasCommand         `command:"update-space-quotas" description:"updates spaces quotas"`
	UpdateSpaceUsersCommand          UpdateSpaceUsersCommand          `command:"update-space-users" description:"update space user roles"`
	CreateSpaceSecurityGroupsCommand CreateSpaceSecurityGroupsCommand `command:"update-space-security-groups" description:"updates space specific security groups"`
	IsolationSegmentsCommand         IsolationSegmentsCommand         `command:"isolation-segments" description:"assigns isolations segments to orgs and spaces"`
	SharePrivateDomainsCommand       SharePrivateDomainsCommand       `command:"share-org-private-domains" description:"shares an existing private domain with the specified org"`
	ServiceAccessCommand             ServiceAccessCommand             `command:"service-access" description:"enables/disables service access for orgs"`
	SharedDomainsCommand             SharedDomainsCommand             `command:"shared-domains" description:"adds/removes shared domains"`
	UpdateOrgsMetadataCommand        UpdateOrgsMetadataCommand        `command:"update-orgs-metadata" description:"updates organizations metadata for each orgConfig.yml"`
	ApplyCommand                     ApplyCommand                     `command:"apply" description:"applies the configuration to your target foundation"`
	ExportConfigurationCommand       ExportConfigurationCommand       `command:"export-config" description:"Exports org and space configurations from an existing Cloud Foundry instance. [Warning: This operation will delete existing config folder]"`
	ExportServiceAccessCommand       ExportServiceAccessCommand       `command:"export-service-access-config" description:"reverse engineer service access into cf-mgmt.yml and remove from orgConfig.yml(s) if present"`
}

var CfMgmt CfMgmtCommand
