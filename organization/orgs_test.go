package organization_test

import (
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	. "github.com/pivotalservices/cf-mgmt/organization"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
)

var _ = Describe("given OrgManager", func() {
	var (
		fakeClient *orgfakes.FakeCFClient
		orgManager DefaultManager
		fakeReader *configfakes.FakeReader
	)

	BeforeEach(func() {
		fakeClient = new(orgfakes.FakeCFClient)
		fakeReader = new(configfakes.FakeReader)
		orgManager = DefaultManager{
			Cfg:    fakeReader,
			Client: fakeClient,
			Peek:   false,
		}
	})

	Context("FindOrg()", func() {
		It("should return an org", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
				},
				{
					Name: "test2",
				},
			}
			fakeClient.ListOrgsReturns(orgs, nil)
			org, err := orgManager.FindOrg("test")
			Ω(err).Should(BeNil())
			Ω(org).ShouldNot(BeNil())
			Ω(org.Name).Should(Equal("test"))
		})
	})
	It("should return an error for unfound org", func() {
		orgs := []cfclient.Org{}
		fakeClient.ListOrgsReturns(orgs, nil)
		_, err := orgManager.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})
	It("should return an error", func() {
		fakeClient.ListOrgsReturns(nil, fmt.Errorf("test"))
		_, err := orgManager.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
	})

	Context("GetOrgGUID()", func() {
		It("should return an GUID", func() {
			orgs := []cfclient.Org{
				{
					Name: "test",
					Guid: "theGUID",
				},
			}
			fakeClient.ListOrgsReturns(orgs, nil)
			guid, err := orgManager.GetOrgGUID("test")
			Ω(err).Should(BeNil())
			Ω(guid).ShouldNot(BeNil())
			Ω(guid).Should(Equal("theGUID"))
		})
	})

	It("should return an error", func() {
		fakeClient.ListOrgsReturns(nil, fmt.Errorf("test"))
		guid, err := orgManager.GetOrgGUID("test")
		Ω(err).ShouldNot(BeNil())
		Ω(guid).Should(Equal(""))
	})

	Context("CreateOrgs()", func() {
		BeforeEach(func() {
			fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
				config.OrgConfig{Org: "test"},
				config.OrgConfig{Org: "test2"},
			}, nil)
		})
		It("should create 2", func() {
			orgs := []cfclient.Org{}
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(2))
		})
		It("should error on list orgs", func() {
			fakeClient.ListOrgsReturns(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should error on create org", func() {
			orgs := []cfclient.Org{}
			fakeClient.ListOrgsReturns(orgs, nil)
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
			fakeClient.ListOrgsReturns(orgs, nil)
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
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(1))
			orgRequest := fakeClient.CreateOrgArgsForCall(0)
			Expect(orgRequest.Name).Should(Equal("test2"))
		})
		It("should not create org if renamed from an org that exists", func() {
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
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateOrgCallCount()).Should(Equal(0))
			Expect(fakeClient.UpdateOrgCallCount()).Should(Equal(1))
			orgGUID, orgRequest := fakeClient.UpdateOrgArgsForCall(0)
			Expect(orgGUID).To(Equal("test2-guid"))
			Expect(orgRequest.Name).To(Equal("new-org"))
		})
	})

	Context("DeleteOrgs()", func() {
		It("should delete 1", func() {
			fakeReader.OrgsReturns(&config.Orgs{
				EnableDeleteOrgs: true,
				Orgs:             []string{"test"},
			}, nil)
			orgs := []cfclient.Org{
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
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgs()
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(1))
			orgGUID, _, _ := fakeClient.DeleteOrgArgsForCall(0)
			Expect(orgGUID).Should(Equal("test2-guid"))
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
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(1))
			orgGUID, _, _ := fakeClient.DeleteOrgArgsForCall(0)
			Expect(orgGUID).Should(Equal("test2-guid"))
		})

		It("should error deleting org that doesn't exist", func() {
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("foo")
			Ω(err).Should(HaveOccurred())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(0))
		})

		It("should not delete any org", func() {
			orgManager.Peek = true
			fakeClient.ListOrgsReturns(orgs, nil)
			err := orgManager.DeleteOrgByName("test2")
			Ω(err).Should(BeNil())
			Expect(fakeClient.DeleteOrgCallCount()).Should(Equal(0))
		})
	})
})
