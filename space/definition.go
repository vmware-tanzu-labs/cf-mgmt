package space

//go:generate counterfeiter -o fakes/fake_cf_client.go types.go CFClient
//go:generate counterfeiter -o fakes/fake_user_mgr.go types.go UserMgr
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
