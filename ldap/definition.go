package ldap

//go:generate mockgen -destination mocks/mock_manager.go "github.com/pivotalservices/cf-mgmt/ldap" Manager
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager

// golang/mock doesn't support vendored dependencies:
// https://github.com/golang/mock/issues/30
//
// So for now fix the import path with sed:

//go:generate sed -i -e s/github.com\/pivotalservices\/cf-mgmt\/vendor\///g ./mocks/mock_manager.go
