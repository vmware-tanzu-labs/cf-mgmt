package cloudcontroller

//go:generate mockgen -destination mocks/mock_manager.go "github.com/pivotalservices/cf-mgmt/cloudcontroller" Manager
//go:generate counterfeiter -o fakes/fake_manager.go cloudcontroller.go Manager
