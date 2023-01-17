package azureAD

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

const SCOPE = "https://graph.microsoft.com/.default"
const GRANT_TYPE = "client_credentials"
const TOKEN_LOGIN_URL = "https://login.microsoftonline.com/" 
const TOKEN_OAUTH_URI =  "/oauth2/v2.0/token"
const GRAPH_URL = "https://graph.microsoft.com/v1.0/"

//Manager -
type Manager struct {
	Config		*config.AzureADConfig
	Token		Token
	groupMap	map[string][]string
	userMap		map[string]*UserType
}

type AADGroupMemberListType struct  {
	Data_type string `json:"@odata.context"`
	Value []UserType `json:"value"`
}

type UserType struct {
	Upn string `json:"userPrincipalName"`
}

type AADGroupType struct  {
	Data_type string `json:"@odata.context"`
	Value []GroupType `json:"value"`
}

type GroupType struct {
	Id string `json:"id"`
}