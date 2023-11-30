package serviceaccess_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	. "github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"

	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serviceaccess", func() {
	var fakeCFClient *fakes.FakeCFClient
	var fakeOrgReader *orgfakes.FakeReader
	var fakeReader *configfakes.FakeReader
	var manager *Manager
	BeforeEach(func() {
		fakeCFClient = &fakes.FakeCFClient{}
		fakeOrgReader = &orgfakes.FakeReader{}
		fakeReader = &configfakes.FakeReader{}
		manager = NewManager(fakeCFClient, fakeOrgReader, fakeReader, false)
	})

	Context("UpdateServiceAccess", func() {
		It("Will do nothing as not enabled", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: false,
			}
			serviceInfo := &ServiceInfo{}
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
		It("Will do nothing as all plans are already public", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: true,
				ServiceAccess:       []*config.Broker{},
			}
			serviceInfo := &ServiceInfo{}
			broker := &ServiceBroker{Name: "mysql"}
			serviceInfo.AddBroker(broker)
			service := &Service{Name: "p-mysql"}
			broker.AddService(service)
			servicePlan := &ServicePlanInfo{Name: "small", ServiceName: "p-mysql", GUID: "small-guid", Public: true}
			service.AddPlan(servicePlan)
			protectedOrgs := []string{"system"}
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
		It("Will change private plan to public", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: true,
				ServiceAccess:       []*config.Broker{},
			}
			serviceInfo := &ServiceInfo{}
			broker := &ServiceBroker{Name: "mysql"}
			serviceInfo.AddBroker(broker)
			service := &Service{Name: "p-mysql"}
			broker.AddService(service)
			servicePlan := &ServicePlanInfo{Name: "small", ServiceName: "p-mysql", GUID: "small-guid", Public: false}
			service.AddPlan(servicePlan)
			protectedOrgs := []string{"system"}
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(1))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})

		It("Will change public plan to private with no access", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: true,
				ServiceAccess: []*config.Broker{
					{
						Name: "mysql-broker",
						Services: []*config.Service{
							{
								Name:          "p-mysql",
								NoAccessPlans: []string{"small"},
							},
						},
					},
				},
			}
			serviceInfo := &ServiceInfo{}
			broker := &ServiceBroker{Name: "mysql-broker"}
			serviceInfo.AddBroker(broker)
			service := &Service{Name: "p-mysql"}
			broker.AddService(service)
			servicePlan := &ServicePlanInfo{Name: "small", ServiceName: "p-mysql", GUID: "small-guid", Public: true}
			service.AddPlan(servicePlan)
			protectedOrgs := []string{"system"}
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})

		It("Will change public plan to private with access to 2 orgs", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: true,
				ServiceAccess: []*config.Broker{
					{
						Name: "mysql-broker",
						Services: []*config.Service{
							{
								Name: "p-mysql",
								LimitedAccessPlans: []*config.PlanVisibility{
									{
										Name: "small",
										Orgs: []string{"test-org"},
									},
								},
							},
						},
					},
				},
			}
			serviceInfo := &ServiceInfo{}
			broker := &ServiceBroker{Name: "mysql-broker"}
			serviceInfo.AddBroker(broker)
			service := &Service{Name: "p-mysql"}
			broker.AddService(service)
			servicePlan := &ServicePlanInfo{Name: "small", ServiceName: "p-mysql", GUID: "small-guid", Public: true}
			service.AddPlan(servicePlan)
			protectedOrgs := []string{"system"}

			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(2))
			planGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(planGUID).Should(Equal("small-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})

		When("The broker is space-scoped", func() {
			It("does not attempt to make the broker's plans public", func() {
				globalCfg := &config.GlobalConfig{
					EnableServiceAccess: true,
				}

				spaceScopedBroker := &ServiceBroker{
					Name:      "space-scoped-broker",
					SpaceGUID: "9f815177-3eee-4d9c-a1d0-42db4f4b49a5",
				}
				service := &Service{
					Name: "some-service",
				}
				service.AddPlan(&ServicePlanInfo{})
				spaceScopedBroker.AddService(service)

				serviceInfo := &ServiceInfo{}
				serviceInfo.AddBroker(spaceScopedBroker)
				protectedOrgs := []string{}

				err := manager.UpdateServiceAccess(globalCfg, serviceInfo, protectedOrgs)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			})
		})
	})

	Context("EnsurePublicAccess", func() {
		It("Should make 1 plan public", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      false,
				ServiceName: "a-service",
			}
			err := manager.EnsurePublicAccess(plan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(1))
			planGUID := fakeCFClient.MakeServicePlanPublicArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
		})
		It("Should return an error", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      false,
				ServiceName: "a-service",
			}
			fakeCFClient.MakeServicePlanPublicReturns(errors.New("error making plan public"))
			err := manager.EnsurePublicAccess(plan)
			Expect(err).Should(MatchError("error making plan public"))
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(1))
			planGUID := fakeCFClient.MakeServicePlanPublicArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
		})

		It("Should peek 1 plan public", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      false,
				ServiceName: "a-service",
			}
			manager.Peek = true
			err := manager.EnsurePublicAccess(plan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
		})
	})

	Context("EnsureNoAccessAccess", func() {
		It("Should make 1 plan noaccess", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			err := manager.EnsureNoAccess(plan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			planGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
		})
		It("Should return an error", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeCFClient.MakeServicePlanPrivateReturns(errors.New("error making plan private"))
			err := manager.EnsureNoAccess(plan)
			Expect(err).Should(MatchError("error making plan private"))
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			planGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
		})
		It("Should peek 1 plan noaccess", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			manager.Peek = true
			err := manager.EnsureNoAccess(plan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
		})
	})

	Context("EnsureLimitedAccess", func() {
		It("Should make 1 plan limited access", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			planGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(1))
			planGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
		It("Should not change existing limited access", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      false,
				ServiceName: "a-service",
			}
			plan.AddOrg(&Visibility{
				OrgGUID:         "test-org-guid",
				ServicePlanGUID: "a-plan-guid",
			})
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
		It("Should make 0 orgs limited access", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			err := manager.EnsureLimitedAccess(plan, []string{}, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			planGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})

		It("Should peek 1 plan limited access", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			manager.Peek = true
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(0))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
		It("Should return an error making private", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeCFClient.MakeServicePlanPrivateReturns(errors.New("error making private"))
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).Should(MatchError("error making private"))
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
		It("Should return an error retrieving org", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeOrgReader.FindOrgReturns(nil, errors.New("error getting org"))
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).Should(MatchError("error getting org"))
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).Should(Equal(1))
			Expect(fakeCFClient.MakeServicePlanPublicCallCount()).Should(Equal(0))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
		})
	})
	Context("ProtectedOrgList", func() {
		It("Should return a list", func() {
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{Name: "foo"},
				{Name: "system"},
				{Name: "bar"},
			}, nil)
			fakeReader.OrgsReturns(&config.Orgs{}, nil)
			protectedOrgsList, err := manager.ProtectedOrgList()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(protectedOrgsList)).Should(BeEquivalentTo(1))
		})
		It("Should error getting org config", func() {
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{Name: "foo"},
				{Name: "system"},
				{Name: "bar"},
			}, nil)
			fakeReader.OrgsReturns(&config.Orgs{}, errors.New("Getting org config"))
			protectedOrgsList, err := manager.ProtectedOrgList()
			Expect(err).Should(MatchError("Getting org config"))
			Expect(len(protectedOrgsList)).Should(BeEquivalentTo(0))
		})
		It("Should error getting orgs", func() {
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{Name: "foo"},
				{Name: "system"},
				{Name: "bar"},
			}, errors.New("Getting orgs"))
			fakeReader.OrgsReturns(&config.Orgs{}, nil)
			protectedOrgsList, err := manager.ProtectedOrgList()
			Expect(err).Should(MatchError("Getting orgs"))
			Expect(len(protectedOrgsList)).Should(BeEquivalentTo(0))
		})
	})
	Context("CreateServiceVisibility", func() {
		It("Creates visibility from org that doesn't have access", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			err := manager.CreatePlanVisibility(servicePlan, "test-org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(1))
			planGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})

		It("Skips creating visibility from org that already has access", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid"})
			err := manager.CreatePlanVisibility(servicePlan, "test-org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(0))
			Expect(servicePlan.OrgHasAccess("test-org-guid")).Should(BeFalse())
		})

		It("errors creating visibility from org that doesn't have access", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			fakeCFClient.CreateServicePlanVisibilityReturns(cfclient.ServicePlanVisibility{}, errors.New("creating visiblity"))
			err := manager.CreatePlanVisibility(servicePlan, "test-org")
			Expect(err).Should(MatchError("creating visiblity"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).Should(Equal(1))
			planGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
	})

	Context("RemoveServiceVisibility", func() {
		It("Removes visibility from org that shouldn't have access", func() {
			fakeOrgReader.FindOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).Should(Equal(1))
			planGUID, orgGUID, async := fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgArgsForCall(0)
			Expect(planGUID).Should(Equal("service-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(async).Should(BeFalse())
		})

		It("Peeks Removes visibility from org that shouldn't have access", func() {
			fakeOrgReader.FindOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			manager.Peek = true
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).Should(Equal(0))
		})

		It("Errors getting org", func() {
			fakeOrgReader.FindOrgByGUIDReturns(nil, errors.New("getting org by guid"))
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).Should(MatchError("getting org by guid"))
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).Should(Equal(0))
		})

		It("errors removing visibility from org that shouldn't have access", func() {
			fakeOrgReader.FindOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgReturns(errors.New("deleting visibility"))
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).Should(MatchError("deleting visibility"))
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).Should(Equal(1))
			planGUID, orgGUID, async := fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgArgsForCall(0)
			Expect(planGUID).Should(Equal("service-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(async).Should(BeFalse())
		})
	})
})
