package space

//{"name":"test2","organization_guid":"76c940e3-1c4c-411b-8672-3edc0651cae7"}

//Manager -
type Manager interface {
	CreateSpace(orgGUID, spaceName string) (space Resource, err error)
	FindSpace(orgGUID, spaceName string) (space Resource, err error)
	CreateSpaces(configFile string) (err error)
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

//InputSpaces -
type InputSpaces struct {
	Org    string   `yml:"org"`
	Spaces []string `yml:"spaces"`
}

//Entity -
type Entity struct {
	Name string `json:"name"`
	//SpacesURL          string `json:"spaces_url"`
	//QuotaURL           string `json:"quota_definition_url"`
	//SpaceQuoteURL      string `json:"space_quota_definitions_url"`
	//UsersURL           string `json:"users_url"`
	//ManagersURL        string `json:"managers_url"`
	//BillingManagersURL string `json:"billing_managers_url"`
	//AuditorsURL        string `json:"auditors_url"`
}

//DefaultSpaceManager -
type DefaultSpaceManager struct {
	Token     string
	SysDomain string
	Spaces    []Resource
}
