package privatedomain_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/privatedomain"
	"github.com/vmwarepivotallabs/cf-mgmt/privatedomain/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		manager    *DefaultManager
		client     *fakes.FakeCFDomainClient
		jobClient  *fakes.FakeCFJobClient
		fakeReader *configfakes.FakeReader
		orgFake    *orgfakes.FakeReader
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFDomainClient)
		jobClient = new(fakes.FakeCFJobClient)
		fakeReader = new(configfakes.FakeReader)
		orgFake = new(orgfakes.FakeReader)
	})
	Context("Manager()", func() {
		BeforeEach(func() {
			manager = &DefaultManager{
				DomainClient: client,
				JobClient:    jobClient,
				Cfg:          fakeReader,
				OrgReader:    orgFake,
				Peek:         false}
		})

		Context("CreatePrivateDomains", func() {
			BeforeEach(func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:            "test",
						PrivateDomains: []string{"test.com"},
					},
				}, nil)
				orgFake.FindOrgReturns(&resource.Organization{
					Name: "test",
					GUID: "test-guid",
				}, nil)
			})
			It("should succeed when no private domain doesn't exist", func() {
				client.CreateReturns(&resource.Domain{Name: "test.com", GUID: "test.com-guid"}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(1))
				_, domainCreate := client.CreateArgsForCall(0)
				Expect(domainCreate.Name).Should(Equal("test.com"))
				Expect(domainCreate.Organization.Data.GUID).Should(Equal("test-guid"))
			})

			It("should error when no private domain doesn't exist", func() {
				client.CreateReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(1))
				_, domainCreate := client.CreateArgsForCall(0)
				Expect(domainCreate.Name).Should(Equal("test.com"))
				Expect(domainCreate.Organization.Data.GUID).Should(Equal("test-guid"))
			})

			It("should succeed and not create already existing private domain", func() {
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("", "test.com", "test-guid"),
				}, nil)
				client.CreateReturns(&resource.Domain{Name: "test.com", GUID: "test.com-guid"}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})

			It("should error when private domain is owned by different org", func() {
				otherOrg := &resource.Organization{
					Name: "foo",
					GUID: "foo-guid",
				}
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("", "test.com", otherOrg.GUID),
				}, nil)
				orgFake.FindOrgByGUIDReturns(otherOrg, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})

			It("should try to remove shared domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
				Expect(client.DeleteCallCount()).Should(Equal(1))
				_, guid := client.DeleteArgsForCall(0)
				Expect(guid).Should(Equal("test.com-guid"))
			})

			It("should error trying to remove shared domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.DeleteReturns("job-guid", errors.New("error"))

				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
				Expect(client.DeleteCallCount()).Should(Equal(1))
				_, guid := client.DeleteArgsForCall(0)
				Expect(guid).Should(Equal("test.com-guid"))
			})

			It("should error getting org config", func() {
				fakeReader.GetOrgConfigsReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})

			It("should error listing orgs", func() {
				orgFake.FindOrgReturns(&resource.Organization{}, errors.New("org test does not exist"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})

			It("should error listing domains", func() {
				client.ListAllReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})

			It("should error when org doesn't exist", func() {
				orgFake.FindOrgReturns(&resource.Organization{}, errors.New("org test does not exist"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
				Expect(err.Error()).Should(Equal("org test does not exist"))
			})

			It("should error listing org shared domains", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
				Expect(client.DeleteCallCount()).Should(Equal(0))
			})
		})

		Context("SharePrivateDomains", func() {
			BeforeEach(func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                  "test2",
						SharedPrivateDomains: []string{"test.com"},
					},
				}, nil)
				orgFake.FindOrgReturns(
					&resource.Organization{
						Name: "test2",
						GUID: "test2-guid",
					}, nil)
			})
			It("should succeed when private domain exists in other org", func() {
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns(nil, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.ShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should error when private domain doesn't already exist", func() {
				client.ListAllReturns(nil, nil)
				client.ListForOrganizationAllReturns(nil, nil)
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("Private Domain [test.com] is not defined"))
				Expect(client.ShareCallCount()).Should(Equal(0))
			})

			It("should error trying to share private domain exists in other org", func() {
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns(nil, nil)
				client.ShareReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.ShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.ShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should succeed unsharing private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                        "test2",
						SharedPrivateDomains:       nil,
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.UnShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should do nothing when no new shared private domains", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                        "test2",
						SharedPrivateDomains:       []string{"test.com"},
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnShareCallCount()).Should(Equal(0))

			})

			It("should succeed unsharing private domain and sharing a private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                        "test2",
						SharedPrivateDomains:       []string{"test2.com"},
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
					newCFDomain("test2.com-guid", "test2.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.ShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test2.com-guid"))

				Expect(client.UnShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID = client.UnShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should error unsharing private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					config.OrgConfig{
						Org:                        "test2",
						SharedPrivateDomains:       nil,
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.ListForOrganizationAllReturns([]*resource.Domain{
					newCFDomain("test.com-guid", "test.com", "test-guid"),
				}, nil)
				client.UnShareReturns(errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.UnShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.UnShareArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should error getting org config", func() {
				fakeReader.GetOrgConfigsReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
			})

			It("should error listing orgs", func() {
				orgFake.ListOrgsReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
			})

			It("should error listing domains", func() {
				client.ListAllReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
			})

			It("should error when org doesn't exist", func() {
				orgFake.FindOrgReturns(&resource.Organization{}, errors.New("org test2 does not exist"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("org test2 does not exist"))
			})

			It("should error listing org private domains", func() {
				client.ListForOrganizationAllReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("CreatePrivateDomain", func() {
			It("should succeed", func() {

				_, err := manager.CreatePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, "test.com")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(1))
				_, domainCreate := client.CreateArgsForCall(0)
				Expect(domainCreate.Name).Should(Equal("test.com"))
				Expect(domainCreate.Organization.Data.GUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				_, err := manager.CreatePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, "test.com")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateCallCount()).Should(Equal(0))
			})
		})

		Context("SharePrivateDomain", func() {
			It("should succeed", func() {

				err := manager.SharePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, &resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.ShareArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.SharePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, &resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareCallCount()).Should(Equal(0))
			})
		})

		Context("DeletePrivateDomain", func() {
			It("should succeed", func() {

				err := manager.DeletePrivateDomain(&resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteCallCount()).Should(Equal(1))
				_, domainGUID := client.DeleteArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.DeletePrivateDomain(&resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteCallCount()).Should(Equal(0))
			})
		})

		Context("RemoveSharedPrivateDomain", func() {
			It("should succeed", func() {

				err := manager.RemoveSharedPrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, &resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnShareCallCount()).Should(Equal(1))
				_, domainGUID, orgGUID := client.UnShareArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.RemoveSharedPrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, &resource.Domain{Name: "test.com", GUID: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnShareCallCount()).Should(Equal(0))
			})
		})
	})
})

func newCFDomain(guid, name, orgGUID string) *resource.Domain {
	return &resource.Domain{
		GUID: guid,
		Name: name,
		Relationships: resource.DomainRelationships{
			Organization: resource.ToOneRelationship{
				Data: &resource.Relationship{
					GUID: orgGUID,
				},
			},
		},
	}
}
