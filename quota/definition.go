package quota

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_space_quota_client.go types.go CFSpaceQuotaClient
//counterfeiter:generate -o fakes/fake_org_quota_client.go types.go CFOrgQuotaClient
