package privatedomain

//go:generate counterfeiter -o fakes/fake_cf_client.go types.go CFClient
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
