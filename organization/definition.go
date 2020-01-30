package organization

//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
//go:generate counterfeiter -o fakes/fake_reader.go types.go Reader
//go:generate counterfeiter -o fakes/fake_cf_client.go types.go CFClient
