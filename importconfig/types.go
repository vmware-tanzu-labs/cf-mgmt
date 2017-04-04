package importconfig

import (
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/uaac"
)

//LDAP Represents constant for value ldap
const LDAP string = "ldap"

//SAML Represents constant for value saml
const SAML string = "saml"

//Manager -
type Manager interface {
	ImportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error
}

//DefaultImportManager  -
type DefaultImportManager struct {
	ConfigDir       string
	UAACMgr         uaac.Manager
	CloudController cloudcontroller.Manager
}
