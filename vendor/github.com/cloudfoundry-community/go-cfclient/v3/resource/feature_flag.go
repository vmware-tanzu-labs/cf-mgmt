package resource

import "time"

type FeatureFlag struct {
	Name      string    `json:"name"`    // The name of the feature flag
	Enabled   bool      `json:"enabled"` // Whether the feature flag is enabled
	UpdatedAt time.Time `json:"updated_at"`

	// The error string returned by the API when a client performs an action disabled by the feature flag
	CustomErrorMessage string `json:"custom_error_message"`

	Links map[string]Link `json:"links"`
}

type FeatureFlagUpdate struct {
	Enabled *bool `json:"enabled,omitempty"` // Whether the feature flag is enabled

	// The error string returned by the API when a client performs an action disabled by the feature flag
	CustomErrorMessage *string `json:"custom_error_message,omitempty"`
}

type FeatureFlagList struct {
	Pagination Pagination     `json:"pagination,omitempty"`
	Resources  []*FeatureFlag `json:"resources,omitempty"`
}

// FeatureFlagType https://v3-apidocs.cloudfoundry.org/version/3.127.0/index.html#the-feature-flag-object
type FeatureFlagType int

const (
	FeatureFlagNone FeatureFlagType = iota

	// FeatureFlagAppBitsUpload When enabled, space developers can upload app bits. When disabled,
	// only admin users can upload app bits.
	FeatureFlagAppBitsUpload

	// FeatureFlagAppScaling When enabled, space developers can perform scaling operations
	// (i.e. change memory, disk, log rate, or instances). When disabled, only admins can perform scaling operations.
	FeatureFlagAppScaling

	// FeatureFlagDiegoDocker When enabled, Docker applications are supported by Diego. When disabled,
	// Docker applications will stop running. It will still be possible to stop and delete them and
	// update their configurations.
	FeatureFlagDiegoDocker

	// FeatureFlagEnvVarVisibility When enabled, all users can see their environment variables.
	// When disabled, no users can see environment variables.
	FeatureFlagEnvVarVisibility

	// FeatureFlagHideMarketPlaceFromUnauthenticatedUsers When enabled, service offerings available in
	// the marketplace will be hidden from unauthenticated users. When disabled, unauthenticated users
	// will be able to see the service offerings available in the marketplace.
	FeatureFlagHideMarketPlaceFromUnauthenticatedUsers

	// FeatureFlagPrivateDomainCreation When enabled, an organization manager can create private domains
	// for that organization. When disabled, only admin users can create private domains.
	FeatureFlagPrivateDomainCreation

	// FeatureFlagResourceMatching When enabled, any user can create resource matches. When disabled,
	// the resource match endpoint always returns an empty array of matches. The package upload endpoint
	// will not cache any uploaded packages for resource matching.
	FeatureFlagResourceMatching

	// FeatureFlagRouteCreation When enabled, a space developer can create routes in a space. When disabled,
	// only admin users can create routes.
	FeatureFlagRouteCreation

	// FeatureFlagRouteSharing When enabled, Space Developers can share routes between two spaces (even across orgs!)
	// in which they have the Space Developer role. When disabled, Space Developers cannot share routes between two spaces.
	FeatureFlagRouteSharing

	// FeatureFlagServiceInstanceCreation When enabled, a space developer can create service instances
	// in a space. When disabled, only admin users can create service instances.
	FeatureFlagServiceInstanceCreation

	// FeatureFlagServiceInstanceSharing When enabled, Space Developers can share service instances
	// between two spaces (even across orgs!) in which they have the Space Developer role. When
	// disabled, Space Developers cannot share service instances between two spaces.
	FeatureFlagServiceInstanceSharing

	// FeatureFlagSetRolesByUserName When enabled, Org Managers or Space Managers can add access roles by username.
	// In order for this feature to be enabled the CF operator must:
	// 1. Enable the /ids/users/ endpoint for UAA
	// 2. Create a UAA cloud_controller_username_lookup client with the scim.userids authority
	FeatureFlagSetRolesByUserName

	// FeatureFlagSpaceDeveloperEnvVarVisibility When enabled, space developers can perform a get on the
	// /v2/apps/:guid/env endpoint, and both space developers and space supporters can perform a get on
	// the /v3/apps/:guid/env and /v3/apps/:guid/environment_variables endpoints. When disabled, neither
	// space developers nor space supporters can access these endpoints.
	FeatureFlagSpaceDeveloperEnvVarVisibility

	// FeatureFlagSpaceScopedPrivateBrokerCreation When enabled, space developers can create space scoped
	// private brokers. When disabled, only admin users can create create space scoped private brokers.
	FeatureFlagSpaceScopedPrivateBrokerCreation

	// FeatureFlagTaskCreation When enabled, space developers can create tasks. When disabled,
	// only admin users can create tasks.
	FeatureFlagTaskCreation

	// FeatureFlagUnsetRolesByUsername When enabled, Org Managers or Space Managers can remove access
	// roles by username. In order for this feature to be enabled the CF operator must:
	// 1. Enable the /ids/users/ endpoint for UAA
	// 2. Create a UAA cloud_controller_username_lookup client with the scim.userids authority
	FeatureFlagUnsetRolesByUsername

	// FeatureFlagUserOrgCreation When enabled, any user can create an organization via the API.
	// When disabled, only admin users can create organizations via the API.
	FeatureFlagUserOrgCreation
)

func (a FeatureFlagType) String() string {
	switch a {
	case FeatureFlagAppBitsUpload:
		return "app_bits_upload"
	case FeatureFlagAppScaling:
		return "app_scaling"
	case FeatureFlagDiegoDocker:
		return "diego_docker"
	case FeatureFlagEnvVarVisibility:
		return "env_var_visibility"
	case FeatureFlagHideMarketPlaceFromUnauthenticatedUsers:
		return "hide_marketplace_from_unauthenticated_users"
	case FeatureFlagPrivateDomainCreation:
		return "private_domain_creation"
	case FeatureFlagResourceMatching:
		return "resource_matching"
	case FeatureFlagRouteCreation:
		return "route_creation"
	case FeatureFlagRouteSharing:
		return "route_sharing"
	case FeatureFlagServiceInstanceCreation:
		return "service_instance_creation"
	case FeatureFlagServiceInstanceSharing:
		return "service_instance_sharing"
	case FeatureFlagSetRolesByUserName:
		return "set_roles_by_username"
	case FeatureFlagSpaceDeveloperEnvVarVisibility:
		return "space_developer_env_var_visibility"
	case FeatureFlagSpaceScopedPrivateBrokerCreation:
		return "space_scoped_private_broker_creation"
	case FeatureFlagTaskCreation:
		return "task_creation"
	case FeatureFlagUnsetRolesByUsername:
		return "unset_roles_by_username"
	case FeatureFlagUserOrgCreation:
		return "user_org_creation"
	}
	return ""
}

func NewFeatureFlagUpdate() *FeatureFlagUpdate {
	return &FeatureFlagUpdate{}
}

func (ff *FeatureFlagUpdate) WithEnabled(enabled bool) *FeatureFlagUpdate {
	ff.Enabled = &enabled
	return ff
}

func (ff *FeatureFlagUpdate) WithCustomErrorMessage(msg string) *FeatureFlagUpdate {
	ff.CustomErrorMessage = &msg
	return ff
}
