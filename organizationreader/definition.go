package organizationreader

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_reader.go types.go Reader
//counterfeiter:generate -o fakes/fake_cf_client.go types.go CFClient
//counterfeiter:generate -o fakes/fake_cf_org_client.go types.go CFOrgClient
