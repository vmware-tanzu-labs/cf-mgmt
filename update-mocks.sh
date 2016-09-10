#! /bin/bash -e

mockgen github.com/pivotalservices/cf-mgmt/cloudcontroller Manager \
  > cloudcontroller/mocks/mock_manager.go
gofmt -w cloudcontroller/mocks/mock_manager.go

mockgen github.com/pivotalservices/cf-mgmt/ldap Manager \
  > ldap/mocks/mock_manager.go
gofmt -w ldap/mocks/mock_manager.go

mockgen github.com/pivotalservices/cf-mgmt/utils Manager \
  > utils/mocks/mock_manager.go
gofmt -w utils/mocks/mock_manager.go

mockgen github.com/pivotalservices/cf-mgmt/uaac Manager \
  > uaac/mocks/mock_manager.go
gofmt -w uaac/mocks/mock_manager.go

echo >&2 "OK"
