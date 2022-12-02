package privatedomain

//go:generate counterfeiter -o fakes/fake_domain_client.go types.go CFDomainClient
//go:generate counterfeiter -o fakes/fake_job_client.go types.go CFJobClient
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
