package export

//LDAP Represents constant for value ldap
const LDAP string = "ldap"

//SAML Represents constant for value saml
const SAML string = "saml"

//Manager -
type Manager interface {
	ExportConfig(excludedOrgs map[string]string, excludedSpaces map[string]string) error
}
