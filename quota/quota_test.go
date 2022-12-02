package quota_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organization/fakes"
	orgreaderfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/quota"
	quotafakes "github.com/vmwarepivotallabs/cf-mgmt/quota/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
)

var _ = Describe("given QuotaManager", func() {
	var (
		fakeReader                  *configfakes.FakeReader
		fakeOrgMgr                  *orgfakes.FakeManager
		fakeOrgReader               *orgreaderfakes.FakeReader
		fakeSpaceQuotaClient        *quotafakes.FakeCFSpaceQuotaClient
		fakeOrganizationQuotaClient *quotafakes.FakeCFOrganizationQuotaClient
		fakeSpaceMgr                *spacefakes.FakeManager
		quotaMgr                    *quota.Manager
	)

	BeforeEach(func() {
		fakeReader = new(configfakes.FakeReader)
		fakeOrgMgr = new(orgfakes.FakeManager)
		fakeOrgReader = new(orgreaderfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		fakeSpaceQuotaClient = new(quotafakes.FakeCFSpaceQuotaClient)
		fakeOrganizationQuotaClient = new(quotafakes.FakeCFOrganizationQuotaClient)
		quotaMgr = &quota.Manager{
			Cfg:                     fakeReader,
			SpaceQuotaClient:        fakeSpaceQuotaClient,
			OrganizationQuotaClient: fakeOrganizationQuotaClient,
			OrgMgr:                  fakeOrgMgr,
			OrgReader:               fakeOrgReader,
			SpaceMgr:                fakeSpaceMgr,
			Peek:                    false,
		}
	})

	Context("ListAllSpaceQuotasForOrg()", func() {
		It("should return 2 quotas", func() {
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
				{
					Name: "quota-1",
					GUID: "quota-1-guid",
				},
				{
					Name: "quota-2",
					GUID: "quota-2-guid",
				},
			}, nil)
			quotas, err := quotaMgr.ListAllSpaceQuotasForOrg("orgGUID")
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.ListAllCallCount()).Should(Equal(1))
			_, opts := fakeSpaceQuotaClient.ListAllArgsForCall(0)
			Expect(opts.OrganizationGUIDs.Values[0]).Should(Equal("orgGUID"))
			Expect(len(quotas)).Should(Equal(2))
			Expect(quotas).Should(HaveKey("quota-1"))
			Expect(quotas).Should(HaveKey("quota-2"))
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
				config.SpaceConfig{
					EnableSpaceQuota: true,
					Space:            "space1",
					Org:              "org1",
				},
				config.SpaceConfig{
					EnableSpaceQuota: false,
					Space:            "space2",
					Org:              "org1",
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeSpaceMgr.FindSpaceReturns(newSpace("space1-guid", "space1", "org1-guid", ""), nil)
		})

		It("should create a quota and assign it", func() {
			fakeSpaceQuotaClient.CreateReturns(&resource.SpaceQuota{Name: "space1", GUID: "space-quota-guid"}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		})

		It("should create a quota that has unlimited memory specified and assign it", func() {
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
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
			OneThousandMB := 1000
			OneHundredMB := 100
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{
				{
					GUID: "org1-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: &OneThousandMB,
					},
				},
				{
					GUID: "org2-quota-guid",
					Apps: resource.OrganizationQuotaApps{
						TotalMemoryInMB: &OneHundredMB,
					},
				},
			}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(*quotaRequest.Apps.TotalMemoryInMB).Should(Equal(1000))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		})

		It("should error creating a quota", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
		})

		// v3 API no longer lets you update a space quota's owning org id
		// hasSpaceQuotaChanged used to return true with the v2 api because of the orgID was empty vs populated
		//It("should update a quota and assign it", func() {
		//	fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{
		//		newSpaceQuota("space-quota-guid", "space1"),
		//	}, nil)
		//	fakeSpaceQuotaClient.UpdateReturns(nil, nil)
		//	err := quotaMgr.CreateSpaceQuotas()
		//	Expect(err).Should(BeNil())
		//	Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
		//	_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
		//	Expect(quotaGUID).Should(Equal("space-quota-guid"))
		//	Expect(*quotaRequest.Name).Should(Equal("space1"))
		//	Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
		//	_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
		//	Expect(quotaGUID).Should(Equal("space-quota-guid"))
		//	Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		//})

		// v3 API no longer lets you update a space quota's owning org id
		// hasSpaceQuotaChanged used to return true with the v2 api because of the orgID was empty vs populated
		//It("should update a quota and not assign it", func() {
		//	fakeSpaceMgr.FindSpaceReturns(newSpace("space1-guid", "space1", "org1-guid", "space-quota-guid"), nil)
		//	fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{newSpaceQuota("space-quota-guid", "space1")}, nil)
		//	fakeSpaceQuotaClient.UpdateReturns(nil, nil)
		//	err := quotaMgr.CreateSpaceQuotas()
		//	Expect(err).Should(BeNil())
		//	Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
		//	_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
		//	Expect(quotaGUID).Should(Equal("space-quota-guid"))
		//	Expect(*quotaRequest.Name).Should(Equal("space1"))
		//	Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		//})

		It("should not update a quota or assign it", func() {
			fakeSpaceMgr.FindSpaceReturns(newSpace("space1-guid", "space1", "org1-guid", "space-quota-guid"), nil)
			spaceQuota := newSpaceQuota("space-quota-guid", "space1")
			spaceQuota.Relationships = resource.SpaceQuotaRelationships{
				Organization: &resource.ToOneRelationship{
					Data: &resource.Relationship{
						GUID: "org1-guid",
					},
				},
			}
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{spaceQuota}, nil)
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should error updating a quota", func() {
			five := 5
			spaceQuota := newSpaceQuota("space-quota-guid", "space1")
			spaceQuota.Apps.TotalInstances = &five
			fakeSpaceQuotaClient.ListAllReturns([]*resource.SpaceQuota{spaceQuota}, nil)
			fakeSpaceQuotaClient.UpdateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeSpaceQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})

		It("should create a quota and fail to assign it", func() {
			fakeSpaceQuotaClient.CreateReturns(newSpaceQuota("space-quota-guid", "space1"), nil)
			fakeSpaceQuotaClient.ApplyReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("space1"))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
			_, quotaGUID, spaceGUIDs := fakeSpaceQuotaClient.ApplyArgsForCall(0)
			Expect(quotaGUID).Should(Equal("space-quota-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
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
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{}, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
		})
		It("Should error listing space quotas", func() {
			fakeSpaceQuotaClient.ListAllReturns(nil, errors.New("error"))
			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("CreateOrgQuotas()", func() {

		BeforeEach(func() {
			orgConfigs := []config.OrgConfig{
				config.OrgConfig{
					EnableOrgQuota: true,
					Org:            "org1",
				},
				config.OrgConfig{
					EnableOrgQuota: false,
					Org:            "org2",
				},
			}
			fakeReader.GetOrgConfigsReturns(orgConfigs, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{Name: "org1", GUID: "org-guid"}, nil)
		})
		It("should create a quota and assign it", func() {
			fakeOrganizationQuotaClient.CreateReturns(&resource.OrganizationQuota{Name: "org1", GUID: "org-quota-guid"}, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrganizationQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeOrganizationQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrganizationQuotaClient.ApplyCallCount()).To(Equal(1))
			_, quotaGUID, orgGUIDs := fakeOrganizationQuotaClient.ApplyArgsForCall(0)
			Expect(orgGUIDs[0]).Should(Equal("org-guid"))
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
		})

		It("should error creating a quota", func() {
			fakeOrganizationQuotaClient.CreateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrganizationQuotaClient.CreateCallCount()).Should(Equal(1))
			_, quotaRequest := fakeOrganizationQuotaClient.CreateArgsForCall(0)
			Expect(*quotaRequest.Name).Should(Equal("org1"))
		})

		It("should update a quota and assign it", func() {
			oneHundred := 100
			orgQuota := newOrganizationQuota("org-quota-guid", "org1")
			orgQuota.Routes.TotalRoutes = &oneHundred
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{orgQuota}, nil)
			fakeOrganizationQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrganizationQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrganizationQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))

			Expect(fakeOrganizationQuotaClient.ApplyCallCount()).To(Equal(1))
			_, quotaGUID, orgGUIDs := fakeOrganizationQuotaClient.ApplyArgsForCall(0)
			Expect(orgGUIDs[0]).Should(Equal("org-guid"))
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
		})

		It("should update a quota and not assign it", func() {
			oneHundred := 100
			fakeOrgReader.FindOrgReturns(newOrganization("org-guid", "org1", "org-quota-guid"), nil)
			orgQuota := newOrganizationQuota("org-quota-guid", "org1")
			orgQuota.Routes.TotalRoutes = &oneHundred
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{orgQuota}, nil)
			fakeOrganizationQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrganizationQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrganizationQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrgMgr.UpdateOrgCallCount()).Should(Equal(0))
		})

		It("should not update a quota or assign it", func() {
			fakeOrgReader.FindOrgReturns(newOrganization("org-guid", "org1", "org-quota-guid"), nil)
			orgQuota := newOrganizationQuota("org-quota-guid", "org1")
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{orgQuota}, nil)
			fakeOrganizationQuotaClient.UpdateReturns(nil, nil)
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrganizationQuotaClient.UpdateCallCount()).Should(Equal(0))
			Expect(fakeOrgMgr.UpdateOrgCallCount()).Should(Equal(0))
		})

		It("should error updating quota", func() {
			tenRoutes := 10
			fakeOrgReader.FindOrgReturns(newOrganization("org-guid", "org1", "org-quota-guid"), nil)
			orgQuota := newOrganizationQuota("org-quota-guid", "org1")
			orgQuota.Routes.TotalRoutes = &tenRoutes
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{orgQuota}, nil)
			fakeOrganizationQuotaClient.UpdateReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrganizationQuotaClient.UpdateCallCount()).Should(Equal(1))
			Expect(fakeOrgMgr.UpdateOrgCallCount()).Should(Equal(0))
		})

		It("should error assigning quota", func() {
			oneHundredRoutes := 100
			fakeOrgReader.FindOrgReturns(newOrganization("org-guid", "org1", "org-quota-guid"), nil)
			orgQuota := newOrganizationQuota("org-quota-guid2", "org1")
			orgQuota.Routes.TotalRoutes = &oneHundredRoutes
			fakeOrganizationQuotaClient.ListAllReturns([]*resource.OrganizationQuota{orgQuota}, nil)
			fakeOrganizationQuotaClient.UpdateReturns(nil, nil)
			fakeOrganizationQuotaClient.ApplyReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
			Expect(fakeOrganizationQuotaClient.UpdateCallCount()).Should(Equal(1))
			_, quotaGUID, quotaRequest := fakeOrganizationQuotaClient.UpdateArgsForCall(0)
			Expect(quotaGUID).Should(Equal("org-quota-guid2"))
			Expect(*quotaRequest.Name).Should(Equal("org1"))
			Expect(fakeOrganizationQuotaClient.ApplyCallCount()).To(Equal(1))
		})

		It("should peek create a quota and peek assign it", func() {
			quotaMgr.Peek = true
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).Should(BeNil())
			Expect(fakeOrganizationQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeOrgMgr.UpdateOrgCallCount()).Should(Equal(0))
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
			fakeOrganizationQuotaClient.ListAllReturns(nil, errors.New("error"))
			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("UpdateSpaceQuota()", func() {
		It("should update a quota", func() {
			fakeSpaceQuotaClient.UpdateReturns(nil, nil)

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", resource.NewSpaceQuotaUpdate().WithName("quota"))
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(1))
		})
		It("should peek and not update a quota", func() {
			quotaMgr.Peek = true
			fakeSpaceQuotaClient.UpdateReturns(nil, nil)

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", resource.NewSpaceQuotaUpdate().WithName("quota"))
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.UpdateCallCount()).Should(Equal(0))
		})
		It("should return an error", func() {
			fakeSpaceQuotaClient.UpdateReturns(nil, errors.New("error"))

			err := quotaMgr.UpdateSpaceQuota("quotaGUID", resource.NewSpaceQuotaUpdate())
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("CreateSpaceQuota()", func() {
		It("should create a quota", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, nil)

			_, err := quotaMgr.CreateSpaceQuota(resource.NewSpaceQuotaUpdate().WithName("quota"))
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
		})
		It("should peek and not create a quota", func() {
			quotaMgr.Peek = true
			fakeSpaceQuotaClient.CreateReturns(nil, nil)

			_, err := quotaMgr.CreateSpaceQuota(resource.NewSpaceQuotaUpdate().WithName("quota"))
			Expect(err).Should(BeNil())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
		})
		It("should return an error", func() {
			fakeSpaceQuotaClient.CreateReturns(nil, errors.New("error"))

			_, err := quotaMgr.CreateSpaceQuota(resource.NewSpaceQuotaUpdate().WithName("quota"))
			Expect(err).ShouldNot(BeNil())
		})
	})

	Context("CreateNamedOrgQuotas()", func() {
		It("Should create a named quota and assign it to org", func() {
			fakeReader.GetOrgQuotasReturns([]config.OrgQuota{
				config.OrgQuota{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				config.OrgConfig{
					Org:        "test",
					NamedQuota: "my-named-quota",
				},
			}, nil)
			fakeOrganizationQuotaClient.CreateReturns(&resource.OrganizationQuota{GUID: "my-named-quota-guid", Name: "my-named-quota"}, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{Name: "test"}, nil)

			err := quotaMgr.CreateOrgQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeOrganizationQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeOrganizationQuotaClient.ApplyCallCount()).To(Equal(1))
		})
	})

	Context("CreateSpaceQuotas()", func() {
		It("Should create a named quota and assign it to space", func() {
			fakeReader.GetSpaceQuotasReturns([]config.SpaceQuota{
				config.SpaceQuota{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "my-named-quota",
					EnableSpaceQuota: false,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(newSpaceQuota("my-named-quota-guid", "my-named-quota"), nil)
			fakeSpaceMgr.FindSpaceReturns(newSpace("test-space-guid", "test-space", "test-org-guid", ""), nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
		})

		It("Should create a space specfic quota", func() {
			fakeReader.GetSpaceQuotasReturns(nil, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "",
					EnableSpaceQuota: true,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(newSpaceQuota("test-space-quota-guid", "test-space"), nil)
			fakeSpaceMgr.FindSpaceReturns(newSpace("test-space-guid", "test-space", "test-org-guid", ""), nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			_, createQuotaRequest := fakeSpaceQuotaClient.CreateArgsForCall(0)
			Expect(*createQuotaRequest.Name).Should(Equal("test-space"))
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(1))
		})

		It("should optimize calls if named quota is empty and enable space quotas if false", func() {
			fakeReader.GetSpaceQuotasReturns([]config.SpaceQuota{
				config.SpaceQuota{
					Name: "my-named-quota",
				},
			}, nil)
			fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
				config.SpaceConfig{
					Org:              "test",
					Space:            "test-space",
					NamedQuota:       "",
					EnableSpaceQuota: false,
				},
			}, nil)
			fakeSpaceQuotaClient.CreateReturns(newSpaceQuota("my-named-quota-guid", "my-named-quota"), nil)
			fakeSpaceMgr.FindSpaceReturns(newSpace("test-space-guid", "test-space", "test-org-guid", ""), nil)

			err := quotaMgr.CreateSpaceQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeReader.GetSpaceQuotasCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeSpaceQuotaClient.ApplyCallCount()).Should(Equal(0))
		})
	})
})

func newSpaceQuota(guid, name string) *resource.SpaceQuota {
	zero := 0
	f := false
	return &resource.SpaceQuota{
		Name: name,
		GUID: guid,
		Apps: resource.SpaceQuotaApps{
			TotalMemoryInMB:      &zero,
			PerProcessMemoryInMB: &zero,
			TotalInstances:       &zero,
			PerAppTasks:          &zero,
		},
		Services: resource.SpaceQuotaServices{
			TotalServiceInstances: &zero,
			PaidServicesAllowed:   &f,
			TotalServiceKeys:      &zero,
		},
		Routes: resource.SpaceQuotaRoutes{
			TotalReservedPorts: &zero,
			TotalRoutes:        &zero,
		},
	}
}

func newSpace(guid, name, orgGUID, quotaGUID string) *resource.Space {
	return &resource.Space{
		Name: name,
		GUID: guid,
		Relationships: &resource.SpaceRelationships{
			Organization: &resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: orgGUID,
				},
			},
			Quota: &resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: quotaGUID,
				},
			},
		},
	}
}

func newOrganizationQuota(guid, name string) *resource.OrganizationQuota {
	zero := 0
	f := false
	return &resource.OrganizationQuota{
		Name: name,
		GUID: guid,
		Apps: resource.OrganizationQuotaApps{
			TotalMemoryInMB:      &zero,
			PerProcessMemoryInMB: &zero,
			TotalInstances:       &zero,
			PerAppTasks:          &zero,
		},
		Services: resource.OrganizationQuotaServices{
			TotalServiceInstances: &zero,
			PaidServicesAllowed:   &f,
			TotalServiceKeys:      &zero,
		},
		Routes: resource.OrganizationQuotaRoutes{
			TotalReservedPorts: &zero,
			TotalRoutes:        &zero,
		},
		Domains: resource.OrganizationQuotaDomains{
			TotalDomains: &zero,
		},
	}
}

func newOrganization(guid, name, quotaGUID string) *resource.Organization {
	return &resource.Organization{
		GUID:      guid,
		Name:      name,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Suspended: nil,
		Relationships: resource.QuotaRelationship{
			Quota: resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: quotaGUID,
				},
			},
		},
		Links:    nil,
		Metadata: nil,
	}
}
