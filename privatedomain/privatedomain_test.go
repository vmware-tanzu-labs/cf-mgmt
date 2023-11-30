package privatedomain_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
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
		client     *fakes.FakeCFClient
		fakeReader *configfakes.FakeReader
		orgFake    *orgfakes.FakeReader
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		fakeReader = new(configfakes.FakeReader)
		orgFake = new(orgfakes.FakeReader)
	})
	Context("Manager()", func() {
		BeforeEach(func() {
			manager = &DefaultManager{
				Client:    client,
				Cfg:       fakeReader,
				OrgReader: orgFake,
				Peek:      false}
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
				client.CreateDomainReturns(&cfclient.Domain{Name: "test.com", Guid: "test.com-guid"}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(1))
				domain, orgGUID := client.CreateDomainArgsForCall(0)
				Expect(domain).Should(Equal("test.com"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should error when no private domain doesn't exist", func() {
				client.CreateDomainReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(1))
				domain, orgGUID := client.CreateDomainArgsForCall(0)
				Expect(domain).Should(Equal("test.com"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should succeed and not create already existing private domain", func() {
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.CreateDomainReturns(&cfclient.Domain{Name: "test.com", Guid: "test.com-guid"}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})

			It("should error when private domain is owned by different org", func() {
				orgFake.FindOrgByGUIDReturns(&resource.Organization{Name: "other-org"}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", OwningOrganizationGuid: "foo-guid"},
				}, nil)
				client.CreateDomainReturns(&cfclient.Domain{Name: "test.com", Guid: "test.com-guid"}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})

			It("should try to remove shared domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				err := manager.CreatePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
				Expect(client.DeleteDomainCallCount()).Should(Equal(1))
				guid := client.DeleteDomainArgsForCall(0)
				Expect(guid).Should(Equal("test.com-guid"))
			})

			It("should error trying to remove shared domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.DeleteDomainReturns(errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
				Expect(client.DeleteDomainCallCount()).Should(Equal(1))
				guid := client.DeleteDomainArgsForCall(0)
				Expect(guid).Should(Equal("test.com-guid"))
			})

			It("should error getting org config", func() {
				fakeReader.GetOrgConfigsReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})

			It("should error listing orgs", func() {
				orgFake.FindOrgReturns(&resource.Organization{}, errors.New("org test does not exist"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})

			It("should error listing domains", func() {
				client.ListDomainsReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})

			It("should error when org doesn't exist", func() {
				orgFake.FindOrgReturns(&resource.Organization{}, errors.New("org test does not exist"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
				Expect(err.Error()).Should(Equal("org test does not exist"))
			})

			It("should error listing org shared domains", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                  "test",
						PrivateDomains:       []string{},
						RemovePrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns(nil, errors.New("error"))
				err := manager.CreatePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
				Expect(client.DeleteDomainCallCount()).Should(Equal(0))
			})
		})

		Context("SharePrivateDomains", func() {
			BeforeEach(func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
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
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns(nil, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.ShareOrgPrivateDomainArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should error when private domain doesn't already exist", func() {
				client.ListDomainsReturns(nil, nil)
				client.ListOrgPrivateDomainsReturns(nil, nil)
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("Private Domain [test.com] is not defined"))
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(0))
			})

			It("should error trying to share private domain exists in other org", func() {
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns(nil, nil)
				client.ShareOrgPrivateDomainReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.ShareOrgPrivateDomainArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should succeed unsharing private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                        "test2",
						SharedPrivateDomains:       nil,
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.UnshareOrgPrivateDomainArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should do nothing when no new shared private domains", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                        "test2",
						SharedPrivateDomains:       []string{"test.com"},
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(0))

			})

			It("should succeed unsharing private domain and sharing a private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                        "test2",
						SharedPrivateDomains:       []string{"test2.com"},
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
					{Name: "test2.com", Guid: "test2.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				err := manager.SharePrivateDomains()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.ShareOrgPrivateDomainArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test2.com-guid"))

				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID = client.UnshareOrgPrivateDomainArgsForCall(0)
				Expect(orgGUID).Should(Equal("test2-guid"))
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should error unsharing private domain", func() {
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org:                        "test2",
						SharedPrivateDomains:       nil,
						RemoveSharedPrivateDomains: true,
					},
				}, nil)
				client.ListDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.ListOrgPrivateDomainsReturns([]cfclient.Domain{
					{Name: "test.com", Guid: "test.com-guid", OwningOrganizationGuid: "test-guid"},
				}, nil)
				client.UnshareOrgPrivateDomainReturns(errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.UnshareOrgPrivateDomainArgsForCall(0)
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
				client.ListDomainsReturns(nil, errors.New("error"))
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
				client.ListOrgPrivateDomainsReturns(nil, errors.New("error"))
				err := manager.SharePrivateDomains()
				Expect(err).Should(HaveOccurred())
			})
		})

		Context("CreatePrivateDomain", func() {
			It("should succeed", func() {

				_, err := manager.CreatePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, "test.com")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(1))
				domain, orgGUID := client.CreateDomainArgsForCall(0)
				Expect(domain).Should(Equal("test.com"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				_, err := manager.CreatePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, "test.com")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateDomainCallCount()).Should(Equal(0))
			})
		})

		Context("SharePrivateDomain", func() {
			It("should succeed", func() {

				err := manager.SharePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.ShareOrgPrivateDomainArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.SharePrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.ShareOrgPrivateDomainCallCount()).Should(Equal(0))
			})
		})

		Context("DeletePrivateDomain", func() {
			It("should succeed", func() {

				err := manager.DeletePrivateDomain(cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteDomainCallCount()).Should(Equal(1))
				domainGUID := client.DeleteDomainArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.DeletePrivateDomain(cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteDomainCallCount()).Should(Equal(0))
			})
		})

		Context("RemoveSharedPrivateDomain", func() {
			It("should succeed", func() {

				err := manager.RemoveSharedPrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(1))
				orgGUID, domainGUID := client.UnshareOrgPrivateDomainArgsForCall(0)
				Expect(domainGUID).Should(Equal("test.com-guid"))
				Expect(orgGUID).Should(Equal("test-guid"))
			})

			It("should peek", func() {
				manager.Peek = true
				err := manager.RemoveSharedPrivateDomain(&resource.Organization{Name: "test", GUID: "test-guid"}, cfclient.Domain{Name: "test.com", Guid: "test.com-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.UnshareOrgPrivateDomainCallCount()).Should(Equal(0))
			})
		})
	})
})
