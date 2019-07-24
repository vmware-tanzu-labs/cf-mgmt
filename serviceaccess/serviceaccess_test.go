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
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableServiceAccess: true,
				ServiceAccess: []config.ServiceVisibility{
					config.ServiceVisibility{
						Service: "p-mysql",
						Plan:    "small",
						Orgs:    []string{"test-org"},
					},
				},
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
			privatePlanArgs := fakeCFClient.MakeServicePlanPrivateArgsForCall(0)
			Expect(privatePlanArgs).To(Equal("small-guid"))
			servicePlanGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("small-guid"))
			Expect(orgGUID).To(Equal("test-org-guid"))
			Expect(fakeCFClient.DeleteServicePlanVisibilityByPlanAndOrgCallCount()).To(Equal(2))
		})
	})

	Context("EnableOrgServiceAccess", func() {
		It("Should add when no visibilities exist", func() {
			servicePlan := &ServicePlanInfo{
				GUID:        "10mb-guid",
				Name:        "10mb",
				ServiceName: "p-mysql",
			}

			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(servicePlan, []string{"test-org"}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(1))
			servicePlanGUID, orgGUID := fakeCFClient.CreateServicePlanVisibilityArgsForCall(0)
			Expect(servicePlanGUID).To(Equal("10mb-guid"))
			Expect(orgGUID).To(Equal("test-org-guid"))
		})

		It("Should not add when visibility already exist", func() {
			servicePlan := &ServicePlanInfo{
				GUID:        "10mb-guid",
				Name:        "10mb",
				ServiceName: "p-mysql",
			}
			servicePlan.AddOrg("test-org-guid", cfclient.ServicePlanVisibility{Guid: "visibility-guid"})

			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			err := manager.EnableOrgServiceAccess(servicePlan, []string{"test-org"}, nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when finding org", func() {
			servicePlan := &ServicePlanInfo{
				GUID:        "10mb-guid",
				Name:        "10mb",
				ServiceName: "p-mysql",
			}

			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, errors.New("Org not found"))
			err := manager.EnableOrgServiceAccess(servicePlan, []string{"test-org"}, nil)
			Expect(err).Should(MatchError("Org not found"))
			Expect(fakeCFClient.CreateServicePlanVisibilityCallCount()).To(Equal(0))
		})

		It("Should error when adding visiblity", func() {
			servicePlan := &ServicePlanInfo{
				GUID:        "10mb-guid",
				Name:        "10mb",
				ServiceName: "p-mysql",
			}
			fakeOrgMgr.FindOrgReturns(cfclient.Org{Guid: "test-org-guid", Name: "test-org"}, nil)
			fakeCFClient.CreateServicePlanVisibilityReturns(cfclient.ServicePlanVisibility{}, errors.New("Error creating visibility"))
			err := manager.EnableOrgServiceAccess(servicePlan, []string{"test-org"}, nil)
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
				plans, err := servicesPlanInfo.GetPlans(service.Label, []string{"small", "large"})
				Expect(err).ToNot(HaveOccurred())
				Expect(len(plans)).To(Equal(2))
				for range plans {
					args := fakeCFClient.ListServicePlansByQueryArgsForCall(i)
					Expect(args).To(BeEquivalentTo(url.Values{
						"q": []string{fmt.Sprintf("%s:%s", "service_guid", service.Guid)},
					}))
				}
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
})
