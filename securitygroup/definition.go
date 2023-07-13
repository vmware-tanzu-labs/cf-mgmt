package securitygroup

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//go:generate counterfeiter -o fakes/fake_security_group_client.go types.go CFSecurityGroupClient
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
