package organization

//Manager -
type Manager interface {
	CreateOrg(orgName string) (org OrgResource, err error)
	FindOrg(orgName string) (org OrgResource, err error)
	SyncOrgs(configFile string) (err error)
}

//OrgResources -
type OrgResources struct {
	OrgResource []OrgResource `json:"resources"`
}

//OrgResource -
type OrgResource struct {
	OrgMetaData OrgMetaData `json:"metadata"`
	OrgEntity   OrgEntity   `json:"entity"`
}

//OrgMetaData -
type OrgMetaData struct {
	GUID string `json:"guid"`
}

//InputOrgs -
type InputOrgs struct {
	Orgs []string `json:"orgs"`
}

//OrgEntity -
type OrgEntity struct {
	Name               string `json:"name"`
	SpacesURL          string `json:"spaces_url"`
	QuotaURL           string `json:"quota_definition_url"`
	SpaceQuoteURL      string `json:"space_quota_definitions_url"`
	UsersURL           string `json:"users_url"`
	ManagersURL        string `json:"managers_url"`
	BillingManagersURL string `json:"billing_managers_url"`
	AuditorsURL        string `json:"auditors_url"`
}

//Org -
type Org struct {
	AccessToken string `json:"access_token"`
}

//DefaultOrgManager -
type DefaultOrgManager struct {
	Token     string
	SysDomain string
	Orgs      []OrgResource
}
