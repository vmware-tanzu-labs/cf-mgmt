package organization

//Manager -
type Manager interface {
	CreateOrg(orgName string) (org Resource, err error)
	FindOrg(orgName string) (org Resource, err error)
	SyncOrgs(configFile string) (err error)
}

//Resources -
type Resources struct {
	Resource []Resource `json:"resources"`
}

//Resource -
type Resource struct {
	MetaData MetaData `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

//MetaData -
type MetaData struct {
	GUID string `json:"guid"`
}

//InputOrgs -
type InputOrgs struct {
	Orgs []string `yaml:"orgs"`
}

//Entity -
type Entity struct {
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
	Orgs      []Resource
}
