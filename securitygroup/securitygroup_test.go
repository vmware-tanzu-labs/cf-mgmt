package securitygroup_test

import (
	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller/fakes"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/config/configfakes"
	. "github.com/pivotalservices/cf-mgmt/securitygroup"
)

var _ = Describe("given SecurityGroupManager", func() {
	var (
		mockCloudController *fakes.FakeManager
		mockConfig          *configfakes.FakeManager
		securityManager     DefaultSecurityGroupManager
	)

	BeforeEach(func() {
		mockCloudController = new(fakes.FakeManager)
		mockConfig = new(configfakes.FakeManager)
		securityManager = DefaultSecurityGroupManager{
			Cfg:             config.NewManager("./fixtures/asg-config"),
			CloudController: mockCloudController,
		}
	})
	Context("CreateApplicationSecurityGroups()", func() {
		It("should create 2 asg", func() {
			test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
			Expect(e).Should(BeNil())
			dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/default_asgs/dns.json")
			Expect(e).Should(BeNil())
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)

			mockCloudController.CreateSecurityGroupReturns("SGGUID", nil)
			mockCloudController.CreateSecurityGroupReturns("SGGUID", nil)
			err := securityManager.CreateApplicationSecurityGroups()
			Expect(err).Should(BeNil())

			name, rules := mockCloudController.CreateSecurityGroupArgsForCall(0)
			Expect(name).Should(BeEquivalentTo("test-asg"))
			Expect(rules).Should(MatchJSON(string(test_asg_bytes)))

			name, rules = mockCloudController.CreateSecurityGroupArgsForCall(1)
			Expect(name).Should(BeEquivalentTo("dns"))
			Expect(rules).Should(MatchJSON(string(dns_bytes)))
		})

		It("should create 1 asg and update 1 asg", func() {
			test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
			Expect(e).Should(BeNil())
			dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/default_asgs/dns.json")
			Expect(e).Should(BeNil())
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-asg-guid", Rules: "[]"}
			sgs["test-default-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-default-asg-guid", Rules: "[]"}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)
			mockCloudController.UpdateSecurityGroupReturns(nil)
			mockCloudController.CreateSecurityGroupReturns("SGGUID", nil)
			err := securityManager.CreateApplicationSecurityGroups()
			Expect(err).Should(BeNil())
			Expect(mockCloudController.UpdateSecurityGroupCallCount()).Should(Equal(1))
			Expect(mockCloudController.CreateSecurityGroupCallCount()).Should(Equal(1))

			guid, name, rules := mockCloudController.UpdateSecurityGroupArgsForCall(0)
			Expect(guid).Should(BeEquivalentTo("test-asg-guid"))
			Expect(name).Should(BeEquivalentTo("test-asg"))
			Expect(rules).Should(MatchJSON(string(test_asg_bytes)))

			name, rules = mockCloudController.CreateSecurityGroupArgsForCall(0)
			Expect(name).Should(BeEquivalentTo("dns"))
			Expect(rules).Should(MatchJSON(string(dns_bytes)))
		})

		It("should not update any and create 1 asg", func() {
			test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
			Expect(e).Should(BeNil())
			dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/default_asgs/dns.json")
			Expect(e).Should(BeNil())
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-asg-guid", Rules: string(test_asg_bytes)}
			sgs["test-default-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-default-asg-guid", Rules: "[]"}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)
			mockCloudController.CreateSecurityGroupReturns("SGGUID", nil)
			err := securityManager.CreateApplicationSecurityGroups()
			Expect(err).Should(BeNil())

			name, rules := mockCloudController.CreateSecurityGroupArgsForCall(0)
			Expect(name).Should(BeEquivalentTo("dns"))
			Expect(rules).Should(MatchJSON(string(dns_bytes)))
		})

		It("should only update asgs", func() {
			test_asg_bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
			Expect(e).Should(BeNil())
			dns_bytes, e := ioutil.ReadFile("./fixtures/asg-config/default_asgs/dns.json")
			Expect(e).Should(BeNil())
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "test-asg-guid", Rules: "[]"}
			sgs["dns"] = cloudcontroller.SecurityGroupInfo{GUID: "dns-asg-guid", Rules: "[]"}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)
			mockCloudController.CreateSecurityGroupReturns("SGGUID", nil)
			err := securityManager.CreateApplicationSecurityGroups()
			Expect(err).Should(BeNil())

			guid, name, rules := mockCloudController.UpdateSecurityGroupArgsForCall(0)
			Expect(guid).Should(BeEquivalentTo("test-asg-guid"))
			Expect(name).Should(BeEquivalentTo("test-asg"))
			Expect(rules).Should(MatchJSON(string(test_asg_bytes)))

			guid, name, rules = mockCloudController.UpdateSecurityGroupArgsForCall(1)
			Expect(guid).Should(BeEquivalentTo("dns-asg-guid"))
			Expect(name).Should(BeEquivalentTo("dns"))
			Expect(rules).Should(MatchJSON(string(dns_bytes)))
		})
	})
	Context("AssignDefaultSecurityGroups()", func() {
		BeforeEach(func() {
			securityManager = DefaultSecurityGroupManager{
				Cfg:             mockConfig,
				CloudController: mockCloudController,
			}
		})
		It("should assign 1 running group and 1 staging group", func() {
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-running"] = cloudcontroller.SecurityGroupInfo{GUID: "test-running-guid", DefaultRunning: false}
			sgs["test-staging"] = cloudcontroller.SecurityGroupInfo{GUID: "test-staging-guid", DefaultStaging: false}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)

			globalConfig := &config.GlobalConfig{
				RunningSecurityGroups: []string{"test-running"},
				StagingSecurityGroups: []string{"test-staging"},
			}
			mockConfig.GetGlobalConfigReturns(globalConfig, nil)

			mockCloudController.AssignRunningSecurityGroupReturns(nil)
			mockCloudController.AssignStagingSecurityGroupReturns(nil)
			err := securityManager.AssignDefaultSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(mockCloudController.AssignRunningSecurityGroupCallCount()).Should(Equal(1))
			sgGUID := mockCloudController.AssignRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-running-guid"))
			Expect(mockCloudController.AssignStagingSecurityGroupCallCount()).Should(Equal(1))
			sgGUID = mockCloudController.AssignStagingSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-staging-guid"))
		})

		It("should not assign anything", func() {
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-running"] = cloudcontroller.SecurityGroupInfo{GUID: "test-running-guid", DefaultRunning: true}
			sgs["test-staging"] = cloudcontroller.SecurityGroupInfo{GUID: "test-staging-guid", DefaultStaging: true}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)

			globalConfig := &config.GlobalConfig{
				RunningSecurityGroups: []string{"test-running"},
				StagingSecurityGroups: []string{"test-staging"},
			}
			mockConfig.GetGlobalConfigReturns(globalConfig, nil)

			err := securityManager.AssignDefaultSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(mockCloudController.AssignRunningSecurityGroupCallCount()).Should(Equal(0))
			Expect(mockCloudController.AssignStagingSecurityGroupCallCount()).Should(Equal(0))
		})

		It("should assign same group to both running and staging", func() {
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-group"] = cloudcontroller.SecurityGroupInfo{GUID: "test-group-guid", DefaultRunning: false, DefaultStaging: false}

			mockCloudController.ListSecurityGroupsReturns(sgs, nil)

			globalConfig := &config.GlobalConfig{
				RunningSecurityGroups: []string{"test-group"},
				StagingSecurityGroups: []string{"test-group"},
			}
			mockConfig.GetGlobalConfigReturns(globalConfig, nil)

			err := securityManager.AssignDefaultSecurityGroups()
			mockCloudController.AssignRunningSecurityGroupReturns(nil)
			mockCloudController.AssignStagingSecurityGroupReturns(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(mockCloudController.AssignRunningSecurityGroupCallCount()).Should(Equal(1))
			sgGUID := mockCloudController.AssignRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-group-guid"))
			Expect(mockCloudController.AssignStagingSecurityGroupCallCount()).Should(Equal(1))
			sgGUID = mockCloudController.AssignStagingSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-group-guid"))
		})

		It("should unassign 1 running group and 1 staging group", func() {
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-running"] = cloudcontroller.SecurityGroupInfo{GUID: "test-running-guid", DefaultRunning: true}
			sgs["test-staging"] = cloudcontroller.SecurityGroupInfo{GUID: "test-staging-guid", DefaultStaging: true}
			mockCloudController.ListSecurityGroupsReturns(sgs, nil)

			globalConfig := &config.GlobalConfig{
				RunningSecurityGroups:        []string{},
				StagingSecurityGroups:        []string{},
				EnableUnassignSecurityGroups: true,
			}
			mockConfig.GetGlobalConfigReturns(globalConfig, nil)

			mockCloudController.UnassignRunningSecurityGroupReturns(nil)
			mockCloudController.UnassignStagingSecurityGroupReturns(nil)
			err := securityManager.AssignDefaultSecurityGroups()
			Expect(err).NotTo(HaveOccurred())
			Expect(mockCloudController.UnassignRunningSecurityGroupCallCount()).Should(Equal(1))
			sgGUID := mockCloudController.UnassignRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-running-guid"))
			Expect(mockCloudController.UnassignStagingSecurityGroupCallCount()).Should(Equal(1))
			sgGUID = mockCloudController.UnassignStagingSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("test-staging-guid"))
		})
	})
})
