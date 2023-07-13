package resource

// Root links to other resources, endpoints, and external services that are relevant to API clients.
type Root struct {
	Links RootLinks `json:"links"`
}

type RootLinks struct {
	Self Link `json:"self"`

	CloudControllerV2 RootCloudController `json:"cloud_controller_v2"`
	CloudControllerV3 RootCloudController `json:"cloud_controller_v3"`

	NetworkPolicyV0 Link       `json:"network_policy_v0"`
	NetworkPolicyV1 Link       `json:"network_policy_v1"`
	Login           Link       `json:"login"`
	Uaa             Link       `json:"uaa"`
	Credhub         Link       `json:"credhub"`
	Routing         Link       `json:"routing"`
	Logging         Link       `json:"logging"`
	LogCache        Link       `json:"log_cache"`
	LogStream       Link       `json:"log_stream"`
	AppSSH          RootAppSSH `json:"app_ssh"`
}

type RootCloudController struct {
	Link
	Meta Meta `json:"meta"`
}

type RootAppSSHMeta struct {
	HostKeyFingerprint string `json:"host_key_fingerprint"`
	OauthClient        string `json:"oauth_client"`
}
type RootAppSSH struct {
	Link
	Meta RootAppSSHMeta `json:"meta"`
}
