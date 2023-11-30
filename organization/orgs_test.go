package organization_test

import (
	"fmt"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
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
		orgManager    DefaultManager
		fakeReader    *configfakes.FakeReader
		fakeOrgReader *orgreaderfakes.FakeReader
		fakeSpaceMgr  *spacefakes.FakeManager
		fakeOrgClient *orgfakes.FakeCFOrgClient
	)

	BeforeEach(func() {
		fakeReader = new(configfakes.FakeReader)
		fakeOrgReader = new(orgreaderfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		fakeOrgClient = new(orgfakes.FakeCFOrgClient)
		orgManager = DefaultManager{
			Cfg:       fakeReader,
			OrgClient: fakeOrgClient,
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
				{Org: "test"},
				{Org: "test2"},
			}, nil)
		})
		It("should create 2", func() {
			orgs := []*resource.Organization{}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeOrgClient.CreateCallCount()).Should(Equal(2))
		})
		It("should error on list orgs", func() {
			fakeOrgReader.ListOrgsReturns(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should error on create org", func() {
			orgs := []*resource.Organization{}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			fakeOrgClient.CreateReturns(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should not create any orgs", func() {
			orgs := []*resource.Organization{
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
			Expect(fakeOrgClient.CreateCallCount()).Should(Equal(0))
		})
		It("should create test2 org", func() {
			orgs := []*resource.Organization{
				{
					Name: "test",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeOrgClient.CreateCallCount()).Should(Equal(1))
			_, orgRequest := fakeOrgClient.CreateArgsForCall(0)
			Expect(orgRequest.Name).Should(Equal("test2"))
		})
		It("should not create org if renamed from an org that exists", func() {
			fakeReader.OrgsReturns(&config.Orgs{
				Orgs: []string{"test", "new-org"},
			}, nil)
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				{Org: "test"},
				{Org: "new-org", OriginalOrg: "test2"},
			}, nil)
			orgs := []*resource.Organization{
				{
					Name: "test",
					GUID: "test-guid",
				},
				{
					Name: "test2",
					GUID: "test2-guid",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			fakeOrgReader.FindOrgReturns(&resource.Organization{
				Name: "test2",
				GUID: "test2-guid",
			}, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeOrgClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeOrgClient.UpdateCallCount()).Should(Equal(1))
			_, orgGUID, orgRequest := fakeOrgClient.UpdateArgsForCall(0)
			Expect(orgGUID).To(Equal("test2-guid"))
			Expect(orgRequest.Name).To(Equal("new-org"))
		})

		When("the orgs.yml orgs list cannot be fetched", func() {
			It("errors", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{}, nil)
				fakeOrgReader.ListOrgsReturns([]*resource.Organization{}, nil)
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
				fakeOrgReader.ListOrgsReturns([]*resource.Organization{
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
				fakeOrgReader.ListOrgsReturns([]*resource.Organization{
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
			orgs := []*resource.Organization{
				{
					Name: "system",
					GUID: "system-guid",
				},
				{
					Name: "some-other-system-org",
					GUID: "some-other-system-guid",
				},
				{
					Name: "test",
					GUID: "test-guid",
				},
				{
					Name: "test2",
					GUID: "test2-guid",
				},
				{
					Name: "redis-test-ORG-1-2017_10_04-20h06m33.481s",
					GUID: "redis-guid",
				},
				{
					Name: "mop-bucket",
					GUID: "some-org-that-matches-p-",
				},
				{
					Name: "p-some-tile",
					GUID: "p-tile-guid",
				},
				{
					Name: "papaya-org",
					GUID: "papaya-guid",
				},
			}
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeOrgClient.DeleteCallCount()).Should(Equal(4))
			_, orgGUID := fakeOrgClient.DeleteArgsForCall(0)
			Expect(orgGUID).Should(Equal("some-other-system-guid"))
			_, orgGUID = fakeOrgClient.DeleteArgsForCall(1)
			Expect(orgGUID).Should(Equal("test2-guid"))
			_, orgGUID = fakeOrgClient.DeleteArgsForCall(2)
			Expect(orgGUID).Should(Equal("some-org-that-matches-p-"))
			_, orgGUID = fakeOrgClient.DeleteArgsForCall(3)
			Expect(orgGUID).Should(Equal("papaya-guid"))
		})
	})

	Context("DeleteOrgByName()", func() {
		var (
			orgs []*resource.Organization
		)

		BeforeEach(func() {
			orgs = []*resource.Organization{
				{
					Name: "system",
					GUID: "system-guid",
				},
				{
					Name: "test",
					GUID: "test-guid",
				},
				{
					Name: "test2",
					GUID: "test2-guid",
				},
				{
					Name: "redis-test-ORG-1-2017_10_04-20h06m33.481s",
					GUID: "redis-guid",
				},
			}
		})

		It("should delete 1", func() {
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeOrgClient.DeleteCallCount()).Should(Equal(1))
			_, orgGUID := fakeOrgClient.DeleteArgsForCall(0)
			Expect(orgGUID).Should(Equal("test2-guid"))
		})

		It("should error deleting org that doesn't exist", func() {
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("foo")
			Ω(err).Should(HaveOccurred())
			Expect(fakeOrgClient.DeleteCallCount()).Should(Equal(0))
		})

		It("should not delete any org", func() {
			orgManager.Peek = true
			fakeOrgReader.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeOrgClient.DeleteCallCount()).Should(Equal(0))
		})
	})
})
