package securitygroup

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_security_group_client.go types.go CFSecurityGroupClient
//counterfeiter:generate-o fakes/fake_mgr.go types.go Manager
