package isosegment

//go:generate counterfeiter -o fakes/fake_iso_client.go types.go CFIsolationSegmentClient
//go:generate counterfeiter -o fakes/fake_org_client.go types.go CFOrganizationClient
//go:generate counterfeiter -o fakes/fake_space_client.go types.go CFSpaceClient
