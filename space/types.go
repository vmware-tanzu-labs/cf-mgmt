package space

import "github.com/pivotalservices/cf-mgmt/securitygroup"

//{"name":"test2","organization_guid":"76c940e3-1c4c-411b-8672-3edc0651cae7"}

//Manager -
type Manager interface {
	CreateSpace(orgName, spaceName string) (space Resource, err error)
	FindSpace(orgName, spaceName string) (space Resource, err error)
	CreateSpaces(configDir string) (err error)
	UpdateSpaces(configDir string) (err error)
	UpdateSpaceUsers(configDir string) (err error)
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
	URL  string `json:"url"`
}

//InputCreateSpaces -
type InputCreateSpaces struct {
	Org    string   `yaml:"org"`
	Spaces []string `yaml:"spaces"`
}

//InputUpdateSpaces -
type InputUpdateSpaces struct {
	Org            string `yaml:"org"`
	Space          string `yaml:"space"`
	DeveloperGroup string `yaml:"space-developer-group"`
	ManagerGroup   string `yaml:"space-manager-group"`
	AuditorGroup   string `yaml:"space-auditor-group"`
	AllowSSH       bool   `yaml:"allow-ssh"`
}

//Entity -
type Entity struct {
	Name           string                   `json:"name"`
	AllowSSH       bool                     `json:"allow_ssh"`
	SecurityGroups []securitygroup.Resource `json:"security_groups"`
	OrgGUID        string                   `json:"organization_guid"`
	Org            Org                      `json:"organization"`
	//SpacesURL          string `json:"spaces_url"`
	//QuotaURL           string `json:"quota_definition_url"`
	//SpaceQuoteURL      string `json:"space_quota_definitions_url"`
	//UsersURL           string `json:"users_url"`
	//ManagersURL        string `json:"managers_url"`
	//BillingManagersURL string `json:"billing_managers_url"`
	//AuditorsURL        string `json:"auditors_url"`
}

//Org -
type Org struct {
	OrgEntity OrgEntity `json:"entity"`
}

//OrgEntity -
type OrgEntity struct {
	Name string `json:"name"`
}

//DefaultSpaceManager -
type DefaultSpaceManager struct {
	Token       string
	UAACToken   string
	SysDomain   string
	Spaces      []Resource
	FilePattern string
	FilePaths   []string
}
