package azureAD

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

const Scope = "https://graph.microsoft.com/.default"
const GrantType = "client_credentials"
const TokenLoginURL = "https://login.microsoftonline.com/"
const TokenOAuthURI = "/oauth2/v2.0/token"
const GraphURL = "https://graph.microsoft.com/v1.0/"

// Manager -
type Manager struct {
	Config   *config.AzureADConfig
	Token    Token
	groupMap map[string][]string
	userMap  map[string]*UserType
}

type AADGroupMemberListType struct {
	DataType string     `json:"@odata.context"`
	Value    []UserType `json:"value"`
}

type UserType struct {
	Upn string `json:"userPrincipalName"`
}

type AADGroupType struct {
	DataType string      `json:"@odata.context"`
	Value    []GroupType `json:"value"`
}

type GroupType struct {
	Id string `json:"id"`
}
