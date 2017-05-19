package ldap

//go:generate mockgen -destination mocks/mock_manager.go "github.com/pivotalservices/cf-mgmt/ldap" Manager

// golang/mock doesn't support vendored dependencies:
// https://github.com/golang/mock/issues/30
//
// So for now fix the import path with sed:

//go:generate sed -i -e s/github.com\/pivotalservices\/cf-mgmt\/vendor\///g ./mocks/mock_manager.go
