package user

//go:generate counterfeiter -o fakes/fake_space_client.go types.go CFSpaceClient
//go:generate counterfeiter -o fakes/fake_job_client.go types.go CFJobClient
//go:generate counterfeiter -o fakes/fake_role_client.go types.go CFRoleClient
//go:generate counterfeiter -o fakes/fake_user_client.go types.go CFUserClient
//go:generate counterfeiter -o fakes/fake_mgr.go types.go Manager
//go:generate counterfeiter -o fakes/fake_ldap_mgr.go types.go LdapManager
