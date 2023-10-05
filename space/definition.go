package space

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_cf_space_client.go types.go CFSpaceClient
//counterfeiter:generate -o fakes/fake_cf_space_feature_client.go types.go CFSpaceFeatureClient
//counterfeiter:generate -o fakes/fake_mgr.go types.go Manager
