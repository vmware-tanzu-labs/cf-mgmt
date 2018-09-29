package ldap

//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
//go:generate counterfeiter -o fakes/fake_connection.go connection.go Connection
