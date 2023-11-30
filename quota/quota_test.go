package quota_test

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgreaderfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/quota"
	quotafakes "github.com/vmwarepivotallabs/cf-mgmt/quota/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/util"
)

var _ = Describe("given QuotaManager", func() {
	var (
		fakeReader           *configfakes.FakeReader
		fakeOrgReader        *orgreaderfakes.FakeReader
		fakeSpaceMgr         *spacefakes.FakeManager
		quotaMgr             *quota.Manager
		fakeSpaceQuotaClient *quotafakes.FakeCFSpaceQuotaClient
		fakeOrgQuotaClient   *quotafakes.FakeCFOrgQuotaClient
	)

	BeforeEach(func() {
		fakeReader = new(configfakes.FakeReader)
		fakeOrgReader = new(orgreaderfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		fakeSpaceQuotaClient = new(quotafakes.FakeCFSpaceQuotaClient)
		fakeOrgQuotaClient = new(quotafakes.FakeCFOrgQuotaClient)
		quotaMgr = &quota.Manager{
			Cfg:              fakeReader,
			SpaceQuoteClient: fakeSpaceQuotaClient,
			OrgQuoteClient:   fakeOrgQuotaClient,
			OrgReader:        fakeOrgReader,
			SpaceMgr:         fakeSpaceMgr,
			Peek:             false,
		}
	})

	Context("ListAllSpaceQuotasForOrg()", func() {
		It("should return 2 quotas", func() {
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "quota-1",
					GUID: "quota-1-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "orgGUID",
							},
						},
					},
				},
				{
					Name: "quota-2",
					GUID: "quota-2-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "orgGUID",
							},
						},
					},
				},
				{
					Name: "quota-3",
					GUID: "quota-3-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "orgGUID-other",
							},
						},
					},
				},
			}, nil)
			quotas, err := quotaMgr.ListAllSpaceQuotasForOrg("orgGUID")
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ListAllCallCount()).Should(Equal(1))
			Expect(len(quotas)).Should(Equal(2))
			Expect(quotas).Should(HaveKey("quota-1"))
			Expect(quotas).Should(HaveKey("quota-2"))

			quotas, err = quotaMgr.ListAllSpaceQuotasForOrg("orgGUID-other")
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ListAllCallCount()).Should(Equal(1))
			Expect(len(quotas)).Should(Equal(1))
			Expect(quotas).Should(HaveKey("quota-3"))
		})
		It("should return an error", func() {
			fakeSpaceQuotaClient.ListAllReturns(nil, errors.New("error"))
			_, err := quotaMgr.ListAllSpaceQuotasForOrg("orgGUID")
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.ListAllCallCount()).Should(Equal(1))
		})
	})

	Context("CreateSpaceQuotas()", func() {

		BeforeEach(func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSpaceQuota:        true,
					Space:                   "space1",
					Org:                     "org1",
					MemoryLimit:             "10G",
					InstanceMemoryLimit:     "unlimited",
					TotalRoutes:             "1000",
					TotalServices:           "100",
					PaidServicePlansAllowed: true,
					TotalReservedRoutePorts: "0",
					TotalServiceKeys:        "unlimited",
					AppInstanceLimit:        "unlimited",
					AppTaskLimit:            "unlimited",
				},
				{
					EnableSpaceQuota:        false,
					Space:                   "space2",
					Org:                     "org1",
					MemoryLimit:             "10G",
					InstanceMemoryLimit:     "unlimited",
					TotalRoutes:             "1000",
					TotalServices:           "100",
					PaidServicePlansAllowed: true,
					TotalReservedRoutePorts: "0",
					TotalServiceKeys:        "unlimited",
					AppInstanceLimit:        "unlimited",
					AppTaskLimit:            "unlimited",
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "space1",
				GUID: "space1-guid",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
				},
			}, nil)
		})
		It("should create a quota and assign it", func() {
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(quotaRequest.Relationships).Should(BeNil())
			Expect(quotaRequest.Apps.TotalInstances).Should(BeNil())
			Expect(quotaRequest.Apps.PerAppTasks).Should(BeNil())
			Expect(quotaRequest.Apps.TotalMemoryInMB).ShouldNot(BeNil())
			Expect(*quotaRequest.Apps.TotalMemoryInMB).Should(Equal(10240))
			Expect(quotaRequest.Apps.PerProcessMemoryInMB).Should(BeNil())
			Expect(quotaRequest.Routes.TotalRoutes).ShouldNot(BeNil())
			Expect(*quotaRequest.Routes.TotalRoutes).Should(Equal(1000))
			Expect(quotaRequest.Services.TotalServiceInstances).ShouldNot(BeNil())
			Expect(*quotaRequest.Services.TotalServiceInstances).Should(Equal(100))
			Expect(quotaRequest.Routes.TotalReservedPorts).ShouldNot(BeNil())
			Expect(*quotaRequest.Routes.TotalReservedPorts).Should(Equal(0))
			Expect(quotaRequest.Services.TotalServiceKeys).Should(BeNil())
			Expect(*quotaRequest.Services.PaidServicesAllowed).Should(BeTrue())
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs).Should(ContainElement("space1-guid"))
		})

		It("should create a quota that has unlimited memory specified and assign it", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					EnableSpaceQuota: true,
					Space:            "space1",
					Org:              "org1",
					MemoryLimit:      "unlimited",
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					GUID: "org1-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: util.GetIntPointer(1000),
					},
				},
				{
					GUID: "org2-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: util.GetIntPointer(1000),
					},
				},
			}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(quotaRequest.Apps).ShouldNot(BeNil())
			Expect(quotaRequest.Apps.TotalMemoryInMB).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs).Should(ContainElement("space1-guid"))
		})

		It("should error creating a quota", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
		})

		It("should update a quota and assign it", func() {
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "space1",
					GUID: "space-quota-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "org1-guid",
							},
						},
					},
				},
			}, nil)
			fakeSpaceQuotaClient.UpdateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(quotaRequest.Relationships).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs).Should(ContainElement("space1-guid"))
		})

		It("should update a quota and not assign it", func() {
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "space1",
				GUID: "space1-guid",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
					Quota: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "space-quota-guid",
						},
					},
				},
			}, nil)
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "space1",
					GUID: "space-quota-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "org1-guid",
							},
						},
					},
				},
			}, nil)
			fakeSpaceQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should not update a quota or assign it", func() {
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "space1",
				GUID: "space1-guid",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
					Quota: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "space-quota-guid",
						}},
				},
			}, nil)
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "space1",
					GUID: "space-quota-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "org1-guid",
							},
						},
					},
					Apps: resource.SpaceQuotaApps{
						TotalInstances:       nil,
						PerAppTasks:          nil,
						TotalMemoryInMB:      util.GetIntPointer(10240),
						PerProcessMemoryInMB: nil,
					},
					Routes: resource.SpaceQuotaRoutes{
						TotalRoutes:        util.GetIntPointer(1000),
						TotalReservedPorts: util.GetIntPointer(0),
					},
					Services: resource.SpaceQuotaServices{
						TotalServiceInstances: util.GetIntPointer(100),
						TotalServiceKeys:      nil,
						PaidServicesAllowed:   util.GetBooleanPointer(true),
					},
				},
			}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should error updating a quota", func() {
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "space1",
					GUID: "space-quota-guid",
					Relationships: resource.SpaceQuotaRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "org1-guid",
							},
						},
					},
				},
			}, nil)
			fakeSpaceQuotaClient.UpdateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
			Expect(err).ShouldNot(BeNil())
			_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should create a quota and fail to assign it", func() {
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			fakeSpaceQuotaClient.ApplyReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs).Should(ContainElement("space1-guid"))
		})

		It("should peek create a quota and peek assign it", func() {
			quotaMgr.Peek = true
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("Should error getting configs", func() {
			fakeReader.GetSpaceConfigsReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
		})
		It("Should error finding space", func() {
			fakeSpaceMgr.FindSpaceReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
		})
		It("Should error listing space quotas", func() {
			fakeSpaceQuotaClient.ListAllReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
		})

		It("Should convert -1 in configuration to unlimited", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					EnableSpaceQuota: true,
					Space:            "space1",
					Org:              "org1",
					MemoryLimit:      "-1",
					AppInstanceLimit: "-1",
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					GUID: "org1-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: util.GetIntPointer(1000),
					},
				},
				{
					GUID: "org2-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: util.GetIntPointer(1000),
					},
				},
			}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(quotaRequest.Apps).ShouldNot(BeNil())
			Expect(quotaRequest.Apps.TotalMemoryInMB).Should(BeNil())
			Expect(quotaRequest.Apps.TotalInstances).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs).Should(ContainElement("space1-guid"))
		})

	})

	Context("CreateOrgQuotas()", func() {

		BeforeEach(func() {
			orgConfigs := []config.OrgConfig{
				{
					EnableOrgQuota:          true,
					Org:                     "org1",
					MemoryLimit:             "unlimited",
					InstanceMemoryLimit:     "unlimited",
					TotalRoutes:             "unlimited",
					TotalServices:           "unlimited",
					PaidServicePlansAllowed: true,
					TotalReservedRoutePorts: "unlimited",
					TotalServiceKeys:        "unlimited",
					AppInstanceLimit:        "unlimited",
					AppTaskLimit:            "unlimited",
				},
				{
					EnableOrgQuota:          false,
					Org:                     "org2",
					MemoryLimit:             "unlimited",
					InstanceMemoryLimit:     "unlimited",
					TotalRoutes:             "unlimited",
					TotalServices:           "unlimited",
					PaidServicePlansAllowed: true,
					TotalReservedRoutePorts: "unlimited",
					TotalServiceKeys:        "unlimited",
					AppInstanceLimit:        "unlimited",
					AppTaskLimit:            "unlimited",
				},
			}
			fakeReader.GetOrgConfigsReturns(orgConfigs, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org-guid"}, nil)
		})
		It("should create a quota and assign it", func() {
			fakeOrgQuotaClient.CreateReturns(&resource.OrganizationQuota{Name: "org1", GUID: "org-quota-guid"}, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrgQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeOrgQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, orgGUID, orgRequest := fakeOrgQuotaClient.ApplyArgsForCall(0)
			Expect(orgGUID).Should(Equal("org-quota-guid"))
			Expect(orgRequest).Should(ContainElement("org-guid"))
		})

		It("should error creating a quota", func() {
			fakeOrgQuotaClient.CreateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrgQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeOrgQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("org1"))
		})

		It("should update a quota and assign it", func() {
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					Name: "org1",
					GUID: "org-quota-guid",
					Routes: resource.OrganizationQuotaRoutes{
						TotalRoutes: util.GetIntPointer(100),
					},
				},
			}, nil)
			fakeOrgQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrgQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrgQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, orgGUID, orgRequest := fakeOrgQuotaClient.ApplyArgsForCall(0)
			Expect(orgGUID).Should(Equal("org-quota-guid"))
			Expect(orgRequest).Should(ContainElement("org-guid"))
		})

		It("should update a quota and not assign it", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Name: "org1", GUID: "org-guid",
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					Name: "org1",
					GUID: "org-quota-guid",
					Routes: resource.OrganizationQuotaRoutes{
						TotalRoutes: util.GetIntPointer(100),
					},
				},
			}, nil)
			fakeOrgQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrgQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrgQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(quotaRequest.Apps.TotalInstances).Should(BeNil())
			Expect(quotaRequest.Apps.PerAppTasks).Should(BeNil())
			Expect(quotaRequest.Apps.TotalMemoryInMB).Should(BeNil())
			Expect(quotaRequest.Apps.PerProcessMemoryInMB).Should(BeNil())
			Expect(quotaRequest.Routes.TotalRoutes).Should(BeNil())
			Expect(quotaRequest.Services.TotalServiceInstances).Should(BeNil())
			Expect(quotaRequest.Routes.TotalReservedPorts).Should(BeNil())
			Expect(quotaRequest.Services.TotalServiceKeys).Should(BeNil())
			Expect(*quotaRequest.Services.PaidServicesAllowed).Should(BeTrue())
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should not update a quota or assign it", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Name: "org1", GUID: "org-guid",
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					Name: "org1",
					GUID: "org-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalInstances:       nil,
						PerAppTasks:          nil,
						TotalMemoryInMB:      nil,
						PerProcessMemoryInMB: nil,
					},
					Routes: resource.OrganizationQuotaRoutes{
						TotalRoutes:        nil,
						TotalReservedPorts: nil,
					},
					Services: resource.OrganizationQuotaServices{
						TotalServiceInstances: nil,
						TotalServiceKeys:      nil,
						PaidServicesAllowed:   util.GetBooleanPointer(true),
					},
					Domains: resource.OrganizationQuotaDomains{
						TotalDomains: nil,
					},
				},
			}, nil)
			fakeOrgQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrgQuotaClient.UpdateCallCount()).Should(Equal(0))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should error updating quota", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Name: "org1", GUID: "org-guid",
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					Name: "org1",
					GUID: "org-quota-guid",
					Routes: resource.OrganizationQuotaRoutes{
						TotalRoutes: util.GetIntPointer(10),
					},
				},
			}, nil)
			fakeOrgQuotaClient.UpdateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrgQuotaClient.UpdateCallCount()).Should(Equal(1))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should error assigning quota", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Name: "org1", GUID: "org-guid",
				Relationships: resource.QuotaRelationship{
					Quota: resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org-quota-guid",
						},
					},
				},
			}, nil)
			fakeOrgQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					Name: "org1",
					GUID: "org-quota-guid2",
					Routes: resource.OrganizationQuotaRoutes{
						TotalRoutes: util.GetIntPointer(100),
					},
				},
			}, nil)
			fakeOrgQuotaClient.UpdateReturns(nil, nil)
			fakeOrgQuotaClient.ApplyReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrgQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrgQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid2"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(1))
		})
		It("should peek create a quota and peek assign it", func() {
			quotaMgr.Peek = true
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrgQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("Should error getting configs", func() {
			fakeReader.GetOrgConfigsReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
		})
		It("Should error finding org", func() {
			fakeOrgReader.FindOrgReturns(&resource.Organization{}, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
		})
		It("Should error listing org quotas", func() {
			fakeOrgQuotaClient.ListAllReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("UpdateSpaceQuota()", func() {
		It("should update a quota", func() {
			fakeSpaceQuotaClient.UpdateReturns(nil, nil)

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", &resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
		})
		It("should peek and not update a quota", func() {
			quotaMgr.Peek = true
			fakeSpaceQuotaClient.UpdateReturns(nil, nil)

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", &resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(0))
		})
		It("should return an error", func() {
			fakeSpaceQuotaClient.UpdateReturns(nil, errors.New("error"))

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", &resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("CreateSpaceQuota()", func() {
		It("should create a quota", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, nil)

			_, err := quotaMgr.CreateSpaceQuota(&resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
		})
		It("should peek and not create a quota", func() {
			quotaMgr.Peek = true
			fakeSpaceQuotaClient.CreateReturns(nil, nil)

			_, err := quotaMgr.CreateSpaceQuota(&resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
		})
		It("should return an error", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, errors.New("error"))

			_, err := quotaMgr.CreateSpaceQuota(&resource.SpaceQuotaCreateOrUpdate{Name: util.GetStringPointer("quota")})
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("CreateNamedOrgQuotas()", func() {
		It("Should create a named quota and assign it to org", func() {
			fakeReader.GetOrgQuotasReturns([]config.OrgQuota{
				{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				{
					Org:        "test",
					NamedQuota: "my-named-quota",
				},
			}, nil)
			fakeOrgQuotaClient.CreateReturns(&resource.OrganizationQuota{GUID: "my-named-quota-guid", Name: "my-named-quota"}, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{Name: "test"}, nil)

			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeOrgQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeOrgQuotaClient.ApplyCallCount()).Should(Equal(1))
		})
	})

	Context("CreateSpaceQuotas()", func() {
		It("Should create a named quota and assign it to space", func() {
			fakeReader.GetSpaceQuotasReturns([]config.SpaceQuota{
				{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "my-named-quota",
					EnableSpaceQuota: false,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{GUID: "my-named-quota-guid", Name: "my-named-quota"}, nil)
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "test-space",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
				},
			}, nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
		})

		It("Should create a space specfic quota", func() {
			fakeReader.GetSpaceQuotasReturns(nil, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "",
					EnableSpaceQuota: true,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{GUID: "test-space-quota-guid", Name: "test-space"}, nil)
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "test-space",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
				},
			}, nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			_, createQuotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*createQuotaRequest.Name).Should(Equal("test-space"))
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
		})

		It("should optimize calls if named quota is empty and enable space quotas if false", func() {
			fakeReader.GetSpaceQuotasReturns([]config.SpaceQuota{
				{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "",
					EnableSpaceQuota: false,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{GUID: "my-named-quota-guid", Name: "my-named-quota"}, nil)
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{Name: "test-space"}, nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeReader.GetSpaceQuotasCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})
	})

})
