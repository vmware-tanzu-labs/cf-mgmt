package role

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_manager.go types.go Manager
//counterfeiter:generate -o fakes/fake_cf_role_client.go types.go CFRoleClient
//counterfeiter:generate -o fakes/fake_cf_user_client.go types.go CFUserClient
//counterfeiter:generate -o fakes/fake_cf_job_client.go types.go CFJobClient
