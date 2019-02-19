package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/uaa"
	. "github.com/pivotalservices/cf-mgmt/user"
	"github.com/pivotalservices/cf-mgmt/user/fakes"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	spacefakes "github.com/pivotalservices/cf-mgmt/space/fakes"
	uaafakes "github.com/pivotalservices/cf-mgmt/uaa/fakes"
)

var _ = Describe("RoleUsers", func() {
	var (
		userManager *DefaultManager
		client      *fakes.FakeCFClient
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		userList    []cfclient.User
		uaaUsers    *uaa.Users
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeManager
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeManager)
		userManager = &DefaultManager{
			Client:     client,
			Cfg:        fakeReader,
			UAAMgr:     uaaFake,
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgMgr:     orgFake,
			Peek:       false,
			LdapConfig: &config.LdapConfig{Origin: "ldap"},
		}
		userList = []cfclient.User{
			cfclient.User{
				Username: "hello",
				Guid:     "world",
			},
			cfclient.User{
				Username: "hello2",
				Guid:     "world2",
			},
		}
		uaaUsers = &uaa.Users{}
		uaaUsers.Add(uaa.User{
			Username: "test",
			Origin:   "uaa",
			GUID:     "test-guid",
		})
		uaaUsers.Add(uaa.User{
			Username: "test-2",
			Origin:   "uaa",
			GUID:     "test2-guid",
		})
		uaaUsers.Add(uaa.User{
			Username: "hello",
			Origin:   "uaa",
			GUID:     "world",
		})
		uaaUsers.Add(uaa.User{
			Username: "hello2",
			Origin:   "uaa",
			GUID:     "world2",
		})
	})
	Context("Space Auditors", func() {

		It("Should succeed on ListSpaceAuditors", func() {
			client.ListSpaceAuditorsReturns(userList, nil)
			users, err := userManager.ListSpaceAuditors("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})
		It("Should error on ListSpaceAuditors", func() {
			client.ListSpaceAuditorsReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceAuditors("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})
	})
	Context("Space Develpers", func() {
		It("Should succeed on ListSpaceDevelopers", func() {
			client.ListSpaceDevelopersReturns(userList, nil)
			users, err := userManager.ListSpaceDevelopers("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})

		It("Should error on ListSpaceDevelopers", func() {
			client.ListSpaceDevelopersReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceDevelopers("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})
	})

	Context("Space Managers", func() {
		It("Should succeed on ListSpaceManagers", func() {
			client.ListSpaceManagersReturns(userList, nil)
			users, err := userManager.ListSpaceManagers("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceManagersArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})
		It("Should error on ListSpaceManagers", func() {
			client.ListSpaceManagersReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceManagers("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
			spaceGUID := client.ListSpaceManagersArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
		})
	})

	Context("ListOrgManager", func() {
		It("should succeed", func() {
			client.ListOrgManagersReturns([]cfclient.User{
				cfclient.User{
					Username: "test",
					Guid:     "test-guid",
				},
				cfclient.User{
					Username: "test2",
					Guid:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgManagers("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListOrgManagersCallCount()).To(Equal(1))
			orgGUID := client.ListOrgManagersArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
		It("should error", func() {
			client.ListOrgManagersReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgManagers("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListOrgManagersCallCount()).To(Equal(1))
			orgGUID := client.ListOrgManagersArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
	})

	Context("ListOrgAuditors", func() {
		It("should succeed", func() {
			client.ListOrgAuditorsReturns([]cfclient.User{
				cfclient.User{
					Username: "test",
					Guid:     "test-guid",
				},
				cfclient.User{
					Username: "test2",
					Guid:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgAuditors("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListOrgAuditorsCallCount()).To(Equal(1))
			orgGUID := client.ListOrgAuditorsArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
		It("should error", func() {
			client.ListOrgAuditorsReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgAuditors("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListOrgAuditorsCallCount()).To(Equal(1))
			orgGUID := client.ListOrgAuditorsArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
	})

	Context("ListOrgBillingManager", func() {
		It("should succeed", func() {
			client.ListOrgBillingManagersReturns([]cfclient.User{
				cfclient.User{
					Username: "test",
					Guid:     "test-guid",
				},
				cfclient.User{
					Username: "test2",
					Guid:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgBillingManagers("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListOrgBillingManagersCallCount()).To(Equal(1))
			orgGUID := client.ListOrgBillingManagersArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
		It("should error", func() {
			client.ListOrgBillingManagersReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgBillingManagers("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListOrgBillingManagersCallCount()).To(Equal(1))
			orgGUID := client.ListOrgBillingManagersArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
		})
	})
})
