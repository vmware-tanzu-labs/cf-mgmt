package serviceaccess_test

import (
	"errors"
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
	var fakeOrgReader *orgfakes.FakeReader
	var fakeReader *configfakes.FakeReader
	var servicePlanClient *fakes.FakeCFServicePlanClient
	var servicePlanVisibilityClient *fakes.FakeCFServicePlanVisibilityClient
	var serviceOfferingClient *fakes.FakeCFServiceOfferingClient
	var serviceBrokerClient *fakes.FakeCFServiceBrokerClient
	var manager *Manager
	BeforeEach(func() {
		servicePlanClient = &fakes.FakeCFServicePlanClient{}
		servicePlanVisibilityClient = &fakes.FakeCFServicePlanVisibilityClient{}
		serviceOfferingClient = &fakes.FakeCFServiceOfferingClient{}
		serviceBrokerClient = &fakes.FakeCFServiceBrokerClient{}
		fakeOrgReader = &orgfakes.FakeReader{}
		fakeReader = &configfakes.FakeReader{}
		manager = NewManager(servicePlanClient, servicePlanVisibilityClient, serviceOfferingClient, serviceBrokerClient,
			fakeOrgReader, fakeReader, false)
	})

	Context("UpdateServiceAccess", func() {
		It("Will do nothing as not enabled", func() {
			globalCfg := &config.GlobalConfig{
				EnableServiceAccess: false,
			}
			serviceInfo := &ServiceInfo{}
			err := manager.UpdateServiceAccess(globalCfg, serviceInfo, []string{})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(updateRequest.Type).To(Equal("public"))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(2))
			_, planGUID, applyRequest := servicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).Should(Equal("small-guid"))
			Expect(applyRequest.Organizations[0].GUID).Should(Equal("test-org-guid"))
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
				Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("public"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))
		})

		It("Should return an error", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      false,
				ServiceName: "a-service",
			}
			servicePlanVisibilityClient.UpdateReturns(nil, errors.New("error making plan public"))
			err := manager.EnsurePublicAccess(plan)
			Expect(err).Should(MatchError("error making plan public"))
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("public"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))
		})

		It("Should return an error", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			servicePlanVisibilityClient.UpdateReturns(nil, errors.New("error making plan private"))
			err := manager.EnsureNoAccess(plan)
			Expect(err).Should(MatchError("error making plan private"))
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))

			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(1))
			_, planGUID, applyRequest := servicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(applyRequest.Organizations[0].GUID).Should(Equal("test-org-guid"))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			_, servicePlanGUID, updateRequest := servicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(updateRequest.Type).To(Equal("admin"))
			Expect(servicePlanGUID).Should(Equal("a-plan-guid"))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(0))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
		})

		It("Should return an error making private", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			servicePlanVisibilityClient.UpdateReturns(nil, errors.New("error making private"))
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).Should(MatchError("error making private"))
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
		})

		It("Should return an error retrieving org", func() {
			plan := &ServicePlanInfo{
				Name:        "a-plan",
				GUID:        "a-plan-guid",
				Public:      true,
				ServiceName: "a-service",
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{}, errors.New("error getting org"))
			err := manager.EnsureLimitedAccess(plan, []string{"test-org"}, []string{})
			Expect(err).Should(MatchError("error getting org"))
			Expect(servicePlanVisibilityClient.UpdateCallCount()).Should(Equal(1))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
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
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(1))
			_, planGUID, applyRequest := servicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(applyRequest.Organizations[0].GUID).Should(Equal("test-org-guid"))
		})

		It("Skips creating visibility from org that already has access", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid"})
			err := manager.CreatePlanVisibility(servicePlan, "test-org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(0))
			Expect(servicePlan.OrgHasAccess("test-org-guid")).Should(BeFalse())
		})

		It("errors creating visibility from org that doesn't have access", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlanVisibilityClient.ApplyReturns(nil, errors.New("creating visiblity"))
			err := manager.CreatePlanVisibility(servicePlan, "test-org")
			Expect(err).Should(MatchError("creating visiblity"))
			Expect(servicePlanVisibilityClient.ApplyCallCount()).Should(Equal(1))
			_, planGUID, servicePlanRequest := servicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).Should(Equal("a-plan-guid"))
			Expect(servicePlanRequest.Organizations[0].GUID).Should(Equal("test-org-guid"))
		})
	})

	Context("RemoveServiceVisibility", func() {
		It("Removes visibility from org that shouldn't have access", func() {
			fakeOrgReader.GetOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(servicePlanVisibilityClient.DeleteCallCount()).Should(Equal(1))
			_, planGUID, orgGUID := servicePlanVisibilityClient.DeleteArgsForCall(0)
			Expect(planGUID).Should(Equal("service-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})

		It("Peeks Removes visibility from org that shouldn't have access", func() {
			fakeOrgReader.GetOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			manager.Peek = true
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(servicePlanVisibilityClient.DeleteCallCount()).Should(Equal(0))
		})

		It("Errors getting org", func() {
			fakeOrgReader.GetOrgByGUIDReturns(&resource.Organization{}, errors.New("getting org by guid"))
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).Should(MatchError("getting org by guid"))
			Expect(servicePlanVisibilityClient.DeleteCallCount()).Should(Equal(0))
		})

		It("errors removing visibility from org that shouldn't have access", func() {
			fakeOrgReader.GetOrgByGUIDReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			servicePlan := &ServicePlanInfo{
				GUID: "a-plan-guid",
			}
			servicePlanVisibilityClient.DeleteReturns(errors.New("deleting visibility"))
			servicePlan.AddOrg(&Visibility{OrgGUID: "test-org-guid", ServicePlanGUID: "service-plan-guid"})
			err := manager.RemoveVisibilities(servicePlan)
			Expect(err).Should(MatchError("deleting visibility"))
			Expect(servicePlanVisibilityClient.DeleteCallCount()).Should(Equal(1))
			_, planGUID, orgGUID := servicePlanVisibilityClient.DeleteArgsForCall(0)
			Expect(planGUID).Should(Equal("service-plan-guid"))
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
	})
})
