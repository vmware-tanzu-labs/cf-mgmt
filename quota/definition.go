package quota

//go:generate counterfeiter -o fakes/fake_space_quota_client.go types.go CFSpaceQuotaClient
//go:generate counterfeiter -o fakes/fake_org_quota_client.go types.go CFOrganizationQuotaClient
//--go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
