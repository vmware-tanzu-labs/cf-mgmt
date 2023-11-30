package privatedomain

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_cf_client.go types.go CFClient
//counterfeiter:generate -o fakes/fake_mgr.go types.go Manager
