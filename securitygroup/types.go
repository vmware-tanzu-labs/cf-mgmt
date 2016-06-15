package securitygroup

//Resource -
type Resource struct {
	MetaData MetaData `json:"metadata"`
	Entity   Entity   `json:"entity"`
}

//MetaData -
type MetaData struct {
	GUID string `json:"guid"`
}

//Entity -
type Entity struct {
	Name    string `json:"name"`
	Rules   []Rule `json:"rules"`
	Running bool   `json:"running_default"`
	Staging bool   `json:"staging_default"`
}

//Rule -
type Rule struct {
	Destination string `json:"destination"`
	Protocol    string `json:"protocol"`
	Ports       string `json:"ports"`
}
