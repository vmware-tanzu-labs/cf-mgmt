package securitygroup_test

import (
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	o "github.com/pivotalservices/cf-mgmt/organization/mocks"
	. "github.com/pivotalservices/cf-mgmt/securitygroup"
	s "github.com/pivotalservices/cf-mgmt/space/mocks"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
)

var _ = Describe("given SecurityGroupManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", config.NewManager("./fixtures/asg-config"))
			Ω(manager).ShouldNot(BeNil())
		})
	})

	var (
		ctrl                *gomock.Controller
		mockCloudController *cc.MockManager
		mockLdap            *ldap.MockManager
		mockUaac            *uaac.MockManager
		mockOrgMgr          *o.MockManager
		mockUserMgr         *s.MockUserMgr
		securityManager     DefaultSecurityGroupManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		mockUserMgr = s.NewMockUserMgr(ctrl)

		securityManager = DefaultSecurityGroupManager{
			Cfg:             config.NewManager("./fixtures/asg-config"),
			CloudController: mockCloudController,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	var _ = Describe("given SecurityGroupManager", func() {
		Describe("create new manager", func() {
			It("should return new manager", func() {
				manager := NewManager("test.com", "token", config.NewManager("./fixtures/asg-config"))
				Ω(manager).ShouldNot(BeNil())
			})
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		Context("CreateApplicationSecurityGroups()", func() {

			It("should create 2 asg", func() {
				test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
				Expect(e).Should(BeNil())
				dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/dns.json")
				Expect(e).Should(BeNil())
				sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
				mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
				mockCloudController.EXPECT().CreateSecurityGroup("test-asg", string(test_asg_bytes)).Return("SGGUID", nil)
				mockCloudController.EXPECT().CreateSecurityGroup("dns", string(dns_bytes)).Return("SGGUID", nil)
				err := securityManager.CreateApplicationSecurityGroups()
				Expect(err).Should(BeNil())
			})

			It("should create 1 asg and update 1 asg", func() {
				test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
				Expect(e).Should(BeNil())
				dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/dns.json")
				Expect(e).Should(BeNil())
				sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
				sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-asg-guid", Rules: "[]"}
				mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
				mockCloudController.EXPECT().UpdateSecurityGroup("test-asg-guid", "test-asg", string(test_asg_bytes)).Return(nil)
				mockCloudController.EXPECT().CreateSecurityGroup("dns", string(dns_bytes)).Return("SGGUID", nil)
				err := securityManager.CreateApplicationSecurityGroups()
				Expect(err).Should(BeNil())
			})

			It("should not update any and create 1 asg", func() {
				test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
				Expect(e).Should(BeNil())
				dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/dns.json")
				Expect(e).Should(BeNil())
				sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
				sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-asg-guid", Rules: string(test_asg_bytes)}
				mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
				mockCloudController.EXPECT().CreateSecurityGroup("dns", string(dns_bytes)).Return("SGGUID", nil)
				err := securityManager.CreateApplicationSecurityGroups()
				Expect(err).Should(BeNil())
			})
		})
	})
})
