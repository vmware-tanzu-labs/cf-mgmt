package space

//go:generate counterfeiter -o fakes/fake_space_client.go types.go CFSpaceClient
//go:generate counterfeiter -o fakes/fake_space_feature_client.go types.go CFSpaceFeatureClient
//go:generate counterfeiter -o fakes/fake_job_client.go types.go CFJobClient
//go:generate counterfeiter -o fakes/fake_org_client.go types.go CFOrganizationClient
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
