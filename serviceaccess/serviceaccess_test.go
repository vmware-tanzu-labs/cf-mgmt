package serviceaccess_test

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/pivotalservices/cf-mgmt/config"
	. "github.com/pivotalservices/cf-mgmt/serviceaccess"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	"github.com/pivotalservices/cf-mgmt/serviceaccess/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Serviceaccess", func() {
	var fakeCFClient *fakes.FakeCFClient
	var fakeOrgMgr *orgfakes.FakeManager
	var fakeReader *configfakes.FakeReader
	var manager *Manager
	BeforeEach(func() {
		fakeCFClient = &fakes.FakeCFClient{}
		fakeOrgMgr = &orgfakes.FakeManager{}
		fakeReader = &configfakes.FakeReader{}
		manager = NewManager(fakeCFClient, fakeOrgMgr, fakeReader, false)
	})

	Context("Apply", func() {
		It("Should succeed", func() {
			servicesToReturn := []cfclient.Service{
				cfclient.Service{Label: "p-mysql", Guid: "p-mysql-guid"},
			}
			plansToReturn := []cfclient.ServicePlan{
				cfclient.ServicePlan{Name: "small", Guid: "small-guid", Public: true},
			}
			visibilitiesToReturn := []cfclient.ServicePlanVisibility{
				cfclient.ServicePlanVisibility{OrganizationGuid: "org1-guid", Guid: "org1-visibility-guid"},
				cfclient.ServicePlanVisibility{OrganizationGuid: "org2-guid", Guid: "org2-visibility-guid"},
			}
			fakeCFClient.ListServicesReturns(servicesToReturn, nil)
			fakeCFClient.ListServicePlansByQueryReturns(plansToReturn, nil)
			fakeCFClient.ListServicePlanVisibilitiesByQueryReturns(visibilitiesToReturn, nil)
			fakeReader.OrgsReturns(&config.Orgs{}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				config.OrgConfig{Org: "test-org", ServiceAccess: map[string][]string{
					"p-mysql": []string{"small"},
				}},
			}, nil)
			fakeOrgMgr.ListOrgsReturns([]cfclient.Org{
				cfclient.Org{Name: "system", Guid: "system-guid"},
				cfclient.Org{Name: "test-org", Guid: "test-org-guid"},
			}, nil)

			fakeOrgMgr.FindOrgReturns(cfclient.Org{Name: "test-org", Guid: "test-org-guid"}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).To(Equal(1))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(2))
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).To(Equal(2))
			privatePlanArgs := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(privatePlanArgs).To(Equal("small-guid"))
			servicePlanGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(orgGUID).To(Equal("system-guid"))
			servicePlanGUID, orgGUID = fakeCFClient.CreateServicePlanVisibilityArgsForCall(1)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(orgGUID).To(Equal("test-org-guid"))
			servicePlanGUID, orgGUID, _ = fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("org1-visibility-guid"))
			Expect(orgGUID).To(Equal("org1-guid"))
			servicePlanGUID, orgGUID, _ = fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgArgsForCall(1)
			Expect(servicePlanGUID).To(Equal("org2-visibility-guid"))
			Expect(orgGUID).To(Equal("org2-guid"))
		})
	})

	Context("RemoveUnknownVisibilites", func() {
		It("Should remove 1 visibility", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("system-org-guid", cfclient.ServicePlanVisibility{Guid: "visibility-guid", OrganizationGuid: "unknown_org_guid"})

			err := manager.RemoveUnknownVisibilites(serviceInfo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).To(Equal(1))
			visibilityGUID, orgGUID, async := fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgArgsForCall(0)
			Expect(visibilityGUID).To(Equal("visibility-guid"))
			Expect(orgGUID).To(Equal("unknown_org_guid"))
			Expect(async).To(Equal(false))
		})
	})

	Context("EnableOrgServiceAccess", func() {
		It("Should add when no visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "20mb-guid", Name: "20mb"})

			fakeOrgMgr.ListOrgsReturns([]cfclient.Org{cfclient.Org{Guid: "system-org-guid", Name: "system"}}, nil)
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(2))
		})
		It("Should add when only 1 visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "20mb-guid", Name: "20mb"})
			servicePlanInfo.AddOrg("system-org-guid", cfclient.ServicePlanVisibility{Guid: "visibility-guid"})

			fakeOrgMgr.ListOrgsReturns([]cfclient.Org{cfclient.Org{Guid: "system-org-guid", Name: "system"}}, nil)
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(1))
			servicePlanGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("10mb-guid"))
			Expect(orgGUID).To(Equal("system-org-guid"))
		})
		It("Should not add when visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("system-org-guid", cfclient.ServicePlanVisibility{Guid: "visibility-guid"})

			fakeOrgMgr.ListOrgsReturns([]cfclient.Org{cfclient.Org{Guid: "system-org-guid", Name: "system"}}, nil)
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when listing orgs", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			fakeOrgMgr.ListOrgsReturns(nil, errors.New("Org not found"))
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).Should(MatchError("Org not found"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when adding visiblity", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})

			fakeOrgMgr.ListOrgsReturns([]cfclient.Org{cfclient.Org{Guid: "system-org-guid", Name: "system"}}, nil)
			fakeCFClient.CreateServicePlanVisibilityReturns(cfclient.ServicePlanVisibility{}, errors.New("Error creating visibility"))
			err := manager.EnableProtectedOrgServiceAccess(serviceInfo, []string{"system"})
			Expect(err).Should(MatchError("Error creating visibility"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(1))
		})
	})
	Context("EnableOrgServiceAccess", func() {
		It("Should add when no visibilities exist", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})

			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": []string{"10mb"},
					},
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(1))
			servicePlanGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("10mb-guid"))
			Expect(orgGUID).To(Equal("test-org-guid"))
		})

		It("Should not add when visibility already exist", func() {
			serviceInfo := &ServiceInfo{}
			servicePlanInfo := serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})
			servicePlanInfo.AddOrg("test-org-guid", cfclient.ServicePlanVisibility{Guid: "visibility-guid"})

			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": []string{"10mb"},
					},
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should warn but not do anything when config doesn't match existing service names", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})

			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-random": []string{"10mb"},
					},
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when finding org", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})

			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-random": []string{"10mb"},
					},
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, errors.New("Org not found"))
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).Should(MatchError("Org not found"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when adding visiblity", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "10mb-guid", Name: "10mb"})

			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					Org: "test-org",
					ServiceAccess: map[string][]string{
						"p-mysql": []string{"10mb"},
					},
				},
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			fakeCFClient.CreateServicePlanVisibilityReturns(cfclient.ServicePlanVisibility{}, errors.New("Error creating visibility"))
			err := manager.EnableOrgServiceAccess(serviceInfo, orgConfigs)
			Expect(err).Should(MatchError("Error creating visibility"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(1))
		})
	})

	Context("ListServiceInfo", func() {
		It("Should return a map of services by name with guid", func() {
			servicesToReturn := []cfclient.Service{
				cfclient.Service{Label: "p-mysql", Guid: "p-mysql-guid"},
				cfclient.Service{Label: "p-rabbit", Guid: "p-rabbit-guid"},
				cfclient.Service{Label: "p-redis", Guid: "p-redis-guid"},
			}
			plansToReturn := []cfclient.ServicePlan{
				cfclient.ServicePlan{Name: "small", Guid: "small-guid"},
				cfclient.ServicePlan{Name: "large", Guid: "large-guid"},
			}
			visibilitiesToReturn := []cfclient.ServicePlanVisibility{
				cfclient.ServicePlanVisibility{OrganizationGuid: "org1-guid", Guid: "org1-visibility-guid"},
				cfclient.ServicePlanVisibility{OrganizationGuid: "org2-guid", Guid: "org2-visibility-guid"},
				cfclient.ServicePlanVisibility{OrganizationGuid: "org3-guid", Guid: "org3-visibility-guid"},
			}
			fakeCFClient.ListServicesReturns(servicesToReturn, nil)
			fakeCFClient.ListServicePlansByQueryReturns(plansToReturn, nil)
			fakeCFClient.ListServicePlanVisibilitiesByQueryReturns(visibilitiesToReturn, nil)
			servicesPlanInfo, err := manager.ListServiceInfo()
			Expect(err).ToNot(HaveOccurred())
			Expect(servicesPlanInfo).ToNot(BeNil())

			for i, service := range servicesToReturn {
				for _, planName := range []string{"small", "large"} {
					plan, err := servicesPlanInfo.GetPlan(service.Label, planName)
					Expect(err).ToNot(HaveOccurred())
					Expect(plan).ToNot(BeNil())
				}
				args := fakeCFClient.ListServicePlansByQueryArgsForCall(i)
				Expect(args).To(BeEquivalentTo(url.Values{
					"q": []string{fmt.Sprintf("%s:%s", "service_guid", service.Guid)},
				}))
			}

		})
		It("Should error listing services", func() {
			fakeCFClient.ListServicesReturns(nil, errors.New("error listing services"))
			_, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("error listing services"))
		})

		It("Should return a map of services by name with guid", func() {
			servicesToReturn := []cfclient.Service{
				cfclient.Service{Label: "p-mysql", Guid: "p-mysql-guid"},
				cfclient.Service{Label: "p-rabbit", Guid: "p-rabbit-guid"},
				cfclient.Service{Label: "p-redis", Guid: "p-redis-guid"},
			}

			fakeCFClient.ListServicesReturns(servicesToReturn, nil)
			fakeCFClient.ListServicePlansByQueryReturns(nil, errors.New("error listing plans"))
			_, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("error listing plans"))
		})

		It("Should return an error listing visibilities", func() {
			servicesToReturn := []cfclient.Service{
				cfclient.Service{Label: "p-mysql", Guid: "p-mysql-guid"},
				cfclient.Service{Label: "p-rabbit", Guid: "p-rabbit-guid"},
				cfclient.Service{Label: "p-redis", Guid: "p-redis-guid"},
			}
			plansToReturn := []cfclient.ServicePlan{
				cfclient.ServicePlan{Name: "small", Guid: "small-guid"},
				cfclient.ServicePlan{Name: "large", Guid: "large-guid"},
			}
			fakeCFClient.ListServicesReturns(servicesToReturn, nil)
			fakeCFClient.ListServicePlansByQueryReturns(plansToReturn, nil)
			fakeCFClient.ListServicePlanVisibilitiesByQueryReturns(nil, errors.New("errors listing visibilities"))
			servicesPlanInfo, err := manager.ListServiceInfo()
			Expect(err).To(MatchError("errors listing visibilities"))
			Expect(servicesPlanInfo).To(BeNil())
		})
	})
	Context("DisablePublicServiceAccess", func() {
		It("Disable plans that are public", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-1", Name: "10mb", Public: false})
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-2", Name: "20mb", Public: false})
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-3", Name: "30mb", Public: true})

			err := manager.DisablePublicServiceAccess(serviceInfo)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).To(Equal(1))
			servicePlanToDisableGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(servicePlanToDisableGUID).To(Equal("guid-3"))
		})

		It("Should error disabling service plan", func() {
			serviceInfo := &ServiceInfo{}
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-1", Name: "10mb", Public: false})
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-2", Name: "20mb", Public: false})
			serviceInfo.AddPlan("p-mysql", cfclient.ServicePlan{Guid: "guid-3", Name: "30mb", Public: true})

			fakeCFClient.MakeServicePlanPrivateReturns(errors.New("error disabling service plan"))
			err := manager.DisablePublicServiceAccess(serviceInfo)
			Expect(err).To(MatchError("error disabling service plan"))
			Expect(fakeCFClient.MakeServicePlanPrivateCallCount()).To(Equal(1))
			servicePlanToDisableGUID := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(servicePlanToDisableGUID).To(Equal("guid-3"))
		})
	})
})
