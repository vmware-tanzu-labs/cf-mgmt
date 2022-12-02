package legacy_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/legacy"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/legacy/fakes"
)

var _ = Describe("Serviceaccess", func() {
	var fakeServicePlanClient *fakes.FakeCFServicePlanClient
	var fakeServicePlanVisibilityClient *fakes.FakeCFServicePlanVisibilityClient
	var fakeServiceOfferingClient *fakes.FakeCFServiceOfferingClient
	var fakeOrgReader *orgfakes.FakeReader
	var fakeReader *configfakes.FakeReader
	var manager *Manager
	BeforeEach(func() {
		fakeServicePlanClient = &fakes.FakeCFServicePlanClient{}
		fakeServicePlanVisibilityClient = &fakes.FakeCFServicePlanVisibilityClient{}
		fakeServiceOfferingClient = &fakes.FakeCFServiceOfferingClient{}
		fakeOrgReader = &orgfakes.FakeReader{}
		fakeReader = &configfakes.FakeReader{}
		manager = NewManager(fakeServicePlanClient, fakeServicePlanVisibilityClient, fakeServiceOfferingClient, fakeOrgReader, fakeReader, false)
	})

	Context("Apply", func() {
		It("Should succeed", func() {
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{EnableServiceAccess: true}, nil)
			serviceOfferings := []*resource.ServiceOffering{
				{Name: "p-mysql", GUID: "p-mysql-guid"},
			}
			plansToReturn := []*resource.ServicePlan{
				{Name: "small", GUID: "small-guid", VisibilityType: "public"},
			}
			visibilityToReturn := &resource.ServicePlanVisibility{
				Type: "organization",
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "org1-guid",
					},
					{
						GUID: "org2-guid",
					},
				},
			}
			fakeServiceOfferingClient.ListAllReturns(serviceOfferings, nil)
			fakeServicePlanClient.ListAllReturns(plansToReturn, nil)
			fakeServicePlanVisibilityClient.GetReturns(visibilityToReturn, nil)

			fakeReader.OrgsReturns(&config.Orgs{}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				{Org: "test-org", ServiceAccess: map[string][]string{
					"p-mysql": []string{"small"},
				}},
			}, nil)
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{Name: "system", GUID: "system-guid"},
				{Name: "test-org", GUID: "test-org-guid"},
			}, nil)

			fakeOrgReader.FindOrgReturns(&resource.Organization{Name: "test-org", GUID: "test-org-guid"}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(fakeServicePlanVisibilityClient.UpdateCallCount()).To(Equal(1))
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(2))
			Expect(fakeServicePlanVisibilityClient.DeleteCallCount()).To(Equal(2))

			_, planGUID, visibilityRequest := fakeServicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(planGUID).To(Equal("small-guid"))
			Expect(visibilityRequest.Type).To(Equal("admin"))

			_, planGUID, planVisibilityRequest := fakeServicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).To(Equal("small-guid"))
			Expect(planVisibilityRequest.Organizations[0].GUID).To(Equal("system-guid"))

			_, planGUID, planVisibilityRequest = fakeServicePlanVisibilityClient.ApplyArgsForCall(1)
			Expect(planGUID).To(Equal("small-guid"))
			Expect(planVisibilityRequest.Organizations[0].GUID).To(Equal("test-org-guid"))
		})
	})

	Context("RemoveUnknownVisibilities", func() {
		It("Should remove 1 visibility", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("unknown_org_guid", &resource.ServicePlanVisibility{
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "unknown_org_guid",
					},
				},
			})

			err := manager.RemoveUnknownVisibilities(serviceInfo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.DeleteCallCount()).To(Equal(1))
			_, planGUID, orgGUID := fakeServicePlanVisibilityClient.DeleteArgsForCall(0)
			Expect(planGUID).To(Equal("10mb-guid"))
			Expect(orgGUID).To(Equal("unknown_org_guid"))
		})
	})

	Context("EnableOrgServiceAccess", func() {
		It("Should add when no visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "20mb-guid", Name: "20mb"})
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{
					GUID: "system-org-guid",
					Name: "system",
				},
			}, nil)
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(2))
		})

		It("Should add when only 1 visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "20mb-guid", Name: "20mb"})
			servicePlanInfo.AddOrg("system-org-guid", &resource.ServicePlanVisibility{
				Type: "organization",
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "system-org-guid",
					},
				},
			})
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{GUID: "system-org-guid", Name: "system"},
			}, nil)

			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(1))
			_, planGUID, planVisibilityRequest := fakeServicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).To(Equal("10mb-guid"))
			Expect(planVisibilityRequest.Organizations[0].GUID).To(Equal("system-org-guid"))
		})

		It("Should not add when visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("system-org-guid", &resource.ServicePlanVisibility{
				Type: "organization",
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "system-org-guid",
					},
				},
			})
			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{GUID: "system-org-guid", Name: "system"},
			}, nil)

			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(0))
		})

		It("Should error when listing orgs", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			fakeOrgReader.ListOrgsReturns(nil, errors.New("Org not found"))
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).Should(MatchError("Org not found"))
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(0))
		})

		It("Should error when adding visiblity", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})

			fakeOrgReader.ListOrgsReturns([]*resource.Organization{
				{GUID: "system-org-guid", Name: "system"},
			}, nil)
			fakeServicePlanVisibilityClient.ApplyReturns(&resource.ServicePlanVisibility{}, errors.New("Error creating visibility"))

			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).Should(MatchError("Error creating visibility"))
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(1))
		})
	})

	Context("EnableOrgServiceAccess", func() {
		It("Should add when no visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			orgConfigs := []config.OrgConfig{
				{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": {"10mb"},
					},
				},
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid", Name: "test-org"}, nil)

			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(1))
			_, planGUID, planVisibilityRequest := fakeServicePlanVisibilityClient.ApplyArgsForCall(0)
			Expect(planGUID).To(Equal("10mb-guid"))
			Expect(planVisibilityRequest.Organizations[0].GUID).To(Equal("test-org-guid"))
		})

		It("Should not add when visibility already exist", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("test-org-guid", &resource.ServicePlanVisibility{
				Type: "organization",
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "test-org-guid",
					},
				},
			})
			orgConfigs := []config.OrgConfig{
				{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": {"10mb"},
					},
				},
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(0))
		})

		It("Should warn but not do anything when config doesn't match existing service names", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			orgConfigs := []config.OrgConfig{
				{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-random": {"10mb"},
					},
				},
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid", Name: "test-org"}, nil)

			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(0))
		})

		It("Should error when finding org", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			orgConfigs := []config.OrgConfig{
				{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-random": {"10mb"},
					},
				},
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid", Name: "test-org"}, errors.New("Org not found"))

			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).Should(MatchError("Org not found"))
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(0))
		})

		It("Should error when adding visiblity", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "10mb-guid", Name: "10mb"})
			orgConfigs := []config.OrgConfig{
				{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": {"10mb"},
					},
				},
			}
			fakeOrgReader.FindOrgReturns(&resource.Organization{GUID: "test-org-guid", Name: "test-org"}, nil)
			fakeServicePlanVisibilityClient.ApplyReturns(&resource.ServicePlanVisibility{}, errors.New("Error creating visibility"))

			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).Should(MatchError("Error creating visibility"))
			Expect(fakeServicePlanVisibilityClient.ApplyCallCount()).To(Equal(1))
		})
	})

	Context("ListServiceInfo", func() {
		It("Should return a map of services by name with guid", func() {
			serviceOfferings := []*resource.ServiceOffering{
				{Name: "p-mysql", GUID: "p-mysql-guid"},
				{Name: "p-rabbit", GUID: "p-rabbit-guid"},
				{Name: "p-redis", GUID: "p-redis-guid"},
			}
			plansToReturn := []*resource.ServicePlan{
				{Name: "small", GUID: "small-guid"},
				{Name: "large", GUID: "large-guid"},
			}
			visibilityToReturn := &resource.ServicePlanVisibility{
				Type: "organization",
				Organizations: []resource.ServicePlanVisibilityRelation{
					{
						GUID: "org1-guid",
					},
					{
						GUID: "org2-guid",
					},
					{
						GUID: "org3-guid",
					},
				},
			}
			fakeServiceOfferingClient.ListAllReturns(serviceOfferings, nil)
			fakeServicePlanClient.ListAllReturns(plansToReturn, nil)
			fakeServicePlanVisibilityClient.GetReturns(visibilityToReturn, nil)

			servicesPlanInfo, err := manager.ListServiceInfo()
			Expect(err).ToNot(HaveOccurred())
			Expect(servicesPlanInfo).ToNot(BeNil())

			for i, serviceOffering := range serviceOfferings {
				plans, err := servicesPlanInfo.GetPlans(serviceOffering.Name, []string{"small", "large"})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(plans)).To(Equal(2))
				for range plans {
					_, opts := fakeServicePlanClient.ListAllArgsForCall(i)
					Expect(opts.ServiceOfferingGUIDs.Values[0]).To(Equal(serviceOffering.GUID))
				}
			}
		})

		It("Should error listing services", func() {
			fakeServiceOfferingClient.ListAllReturns(nil, errors.New("error listing services"))
			_, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("error listing services"))
		})

		It("Should return a map of services by name with guid", func() {
			serviceOfferings := []*resource.ServiceOffering{
				{Name: "p-mysql", GUID: "p-mysql-guid"},
				{Name: "p-rabbit", GUID: "p-rabbit-guid"},
				{Name: "p-redis", GUID: "p-redis-guid"},
			}
			fakeServiceOfferingClient.ListAllReturns(serviceOfferings, nil)
			fakeServicePlanClient.ListAllReturns(nil, errors.New("error listing plans"))
			_, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("error listing plans"))
		})

		It("Should return an error listing visibilities", func() {
			serviceOfferings := []*resource.ServiceOffering{
				{Name: "p-mysql", GUID: "p-mysql-guid"},
				{Name: "p-rabbit", GUID: "p-rabbit-guid"},
				{Name: "p-redis", GUID: "p-redis-guid"},
			}
			plansToReturn := []*resource.ServicePlan{
				{Name: "small", GUID: "small-guid"},
				{Name: "large", GUID: "large-guid"},
			}
			fakeServiceOfferingClient.ListAllReturns(serviceOfferings, nil)
			fakeServicePlanClient.ListAllReturns(plansToReturn, nil)
			fakeServicePlanVisibilityClient.GetReturns(nil, errors.New("errors listing visibilities"))

			servicesPlanInfo, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("errors listing visibilities"))
			Expect(servicesPlanInfo).To(BeNil())
		})
	})
	Context("DisablePublicServiceAccess", func() {
		It("Disable plans that are public", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-1", Name: "10mb", VisibilityType: "organization"})
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-2", Name: "20mb", VisibilityType: "organization"})
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-3", Name: "30mb", VisibilityType: "public"})

			err := manager.DisablePublicServiceAccess(serviceInfo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeServicePlanVisibilityClient.UpdateCallCount()).To(Equal(1))
			_, planGUID, planVisibilityRequest := fakeServicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(planGUID).To(Equal("guid-3"))
			Expect(planVisibilityRequest.Type).To(Equal(resource.ServicePlanVisibilityAdmin.String()))
		})

		It("Should error disabling service plan", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-1", Name: "10mb", VisibilityType: "organization"})
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-2", Name: "20mb", VisibilityType: "organization"})
			serviceInfo.AddPlan("p-mysql", &resource.ServicePlan{GUID: "guid-3", Name: "30mb", VisibilityType: "public"})
			fakeServicePlanVisibilityClient.UpdateReturns(nil, errors.New("error disabling service plan"))

			err := manager.DisablePublicServiceAccess(serviceInfo)
			Expect(err).To(MatchError("error disabling service plan"))
			Expect(fakeServicePlanVisibilityClient.UpdateCallCount()).To(Equal(1))
			_, planGUID, planVisibilityRequest := fakeServicePlanVisibilityClient.UpdateArgsForCall(0)
			Expect(planGUID).To(Equal("guid-3"))
			Expect(planVisibilityRequest.Type).To(Equal(resource.ServicePlanVisibilityAdmin.String()))
		})
	})
})
