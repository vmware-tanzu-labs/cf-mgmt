package organization_test

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/organization"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organization/fakes"
	orgreaderfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
)

var _ = Describe("given OrgManager", func() {
	var (
		fakeClient    *orgfakes.FakeCFClient
		orgManager    DefaultManager
		fakeReader    *configfakes.FakeReader
		fakeOrgReader *orgreaderfakes.FakeReader
		fakeSpaceMgr  *spacefakes.FakeManager
	)

	BeforeEach(func() {
		fakeClient = new(orgfakes.FakeCFClient)
		fakeReader = new(configfakes.FakeReader)
		fakeOrgReader = new(orgreaderfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		orgManager = DefaultManager{
			Cfg:       fakeReader,
			Client:    fakeClient,
			Peek:      false,
			SpaceMgr:  fakeSpaceMgr,
			OrgReader: fakeOrgReader,
		}
	})

	Context("CreateOrgs()", func() {
		BeforeEach(func() {
			fakeReader.OrgsReturns(&config.Orgs{
				Orgs: []string{"test", "test2"},
			}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				config.OrgConfig{Org: "test"},
				config.OrgConfig{Org: "test2"},
			}, nil)
		})
		It("should create 2", func() {
			orgs := []cfclient.Org{}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(2))
		})
		It("should error on list orgs", func() {
			fakeOrgReader.ListOrgsReturns(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should error on create org", func() {
			orgs := []cfclient.Org{}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			fakeClient.CreateOrgReturns(cfclient.Org{}, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should not create any orgs", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
				},
				{
					Name: "test2",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(0))
		})
		It("should create test2 org", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(1))
			orgRequest := fakeClient.CreateOrgArgsForCall(0)
			Expect(orgRequest.Name).Should(Equal("test2"))
		})
		It("should not create org if renamed from an org that exists", func() {
			fakeReader.OrgsReturns(&config.Orgs{
				Orgs: []string{"test", "new-org"},
			}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				config.OrgConfig{Org: "test"},
				config.OrgConfig{Org: "new-org", OriginalOrg: "test2"},
			}, nil)
			orgs := []cfclient.Org{
				{
					Name: "test",
					Guid: "test-guid",
				},
				{
					Name: "test2",
					Guid: "test2-guid",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			fakeOrgReader.FindOrgReturns(cfclient.Org{
				Name: "test2",
				Guid: "test2-guid",
			}, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(0))
			Expect(fakeClient.UpdateOrgCallCount()).Should(Equal(1))
			orgGUID, orgRequest := fakeClient.UpdateOrgArgsForCall(0)
			Expect(orgGUID).To(Equal("test2-guid"))
			Expect(orgRequest.Name).To(Equal("new-org"))
		})

		When("the orgs.yml orgs list cannot be fetched", func() {
			It("errors", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{}, nil)
				fakeOrgReader.ListOrgsReturns([]cfclient.Org{}, nil)
				fakeReader.OrgsReturns(nil, fmt.Errorf("some error"))
				err := orgManager.CreateOrgs()
				Expect(err).Should(HaveOccurred())
			})
		})

		When("an org exists in an orgConfig, but not in orgs.yml", func() {
			It("errors", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{Org: "in-org-list"},
					{Org: "not-in-org-list"},
				}, nil)
				fakeReader.OrgsReturns(&config.Orgs{
					Orgs: []string{"in-org-list"},
				}, nil)
				fakeOrgReader.ListOrgsReturns([]cfclient.Org{
					{Name: "in-org-list"},
				}, nil)

				err := orgManager.CreateOrgs()
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError("[not-in-org-list] found in an orgConfig but not in orgs.yml"))
			})
		})

		When("an org has been renamed in an orgConfig, but not in orgs.yml", func() {
			It("errors", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{Org: "in-org-list"},
					{Org: "not-in-org-list", OriginalOrg: "was-in-org-list"},
				}, nil)
				fakeReader.OrgsReturns(&config.Orgs{
					Orgs: []string{"in-org-list", "was-in-org-list"},
				}, nil)
				fakeOrgReader.ListOrgsReturns([]cfclient.Org{
					{Name: "in-org-list"},
				}, nil)

				err := orgManager.CreateOrgs()
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(MatchError("[not-in-org-list] found in an orgConfig but not in orgs.yml"))
			})
		})
	})

	Context("DeleteOrgs()", func() {
		It("should delete 4", func() {
			fakeReader.OrgsReturns(&config.Orgs{
				EnableDeleteOrgs: true,
				Orgs:             []string{"test"},
			}, nil)

			fakeReader.GetOrgConfigReturns(&config.OrgConfig{}, nil)
			orgs := []cfclient.Org{
				cfclient.Org{
					Name: "system",
					Guid: "system-guid",
				},
				cfclient.Org{
					Name: "some-other-system-org",
					Guid: "some-other-system-guid",
				},
				cfclient.Org{
					Name: "test",
					Guid: "test-guid",
				},
				cfclient.Org{
					Name: "test2",
					Guid: "test2-guid",
				},
				cfclient.Org{
					Name: "redis-test-ORG-1-2017_10_04-20h06m33.481s",
					Guid: "redis-guid",
				},
				cfclient.Org{
					Name: "mop-bucket",
					Guid: "some-org-that-matches-p-",
				},
				cfclient.Org{
					Name: "p-some-tile",
					Guid: "p-tile-guid",
				},
				cfclient.Org{
					Name: "papaya-org",
					Guid: "papaya-guid",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(4))
			orgGUID, _, _ := fakeClient.DeleteOrgArgsForCall(0)
			Expect(orgGUID).Should(Equal("some-other-system-guid"))
			orgGUID, _, _ = fakeClient.DeleteOrgArgsForCall(1)
			Expect(orgGUID).Should(Equal("test2-guid"))
			orgGUID, _, _ = fakeClient.DeleteOrgArgsForCall(2)
			Expect(orgGUID).Should(Equal("some-org-that-matches-p-"))
			orgGUID, _, _ = fakeClient.DeleteOrgArgsForCall(3)
			Expect(orgGUID).Should(Equal("papaya-guid"))
		})
	})

	Context("DeleteOrgByName()", func() {
		var (
			orgs []cfclient.Org
		)

		BeforeEach(func() {
			orgs = []cfclient.Org{
				cfclient.Org{
					Name: "system",
					Guid: "system-guid",
				},
				cfclient.Org{
					Name: "test",
					Guid: "test-guid",
				},
				cfclient.Org{
					Name: "test2",
					Guid: "test2-guid",
				},
				cfclient.Org{
					Name: "redis-test-ORG-1-2017_10_04-20h06m33.481s",
					Guid: "redis-guid",
				},
			}
		})

		It("should delete 1", func() {
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(1))
			orgGUID, _, _ := fakeClient.DeleteOrgArgsForCall(0)
			Expect(orgGUID).Should(Equal("test2-guid"))
		})

		It("should error deleting org that doesn't exist", func() {
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("foo")
			Ω(err).Should(HaveOccurred())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(0))
		})

		It("should not delete any org", func() {
			orgManager.Peek = true
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(0))
		})
	})

	Context("ClearMetadata()", func() {
		It("should remove metadata from given org", func() {
			fakeClient.SupportsMetadataAPIReturns(true, nil)
			org := cfclient.Org{
				Guid: "org-guid",
			}
			err := orgManager.ClearMetadata(org)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.RemoveOrgMetadataCallCount()).Should(Equal(1))
			Expect(fakeClient.RemoveOrgMetadataArgsForCall(0)).Should(Equal("org-guid"))
		})
	})
})
