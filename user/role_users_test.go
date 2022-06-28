package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"

	. "github.com/onsi/ginkgo"

	. "github.com/onsi/gomega"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("RoleUsers", func() {
	var (
		userManager *DefaultManager
		client      *fakes.FakeCFClient
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		userList    []cfclient.V3User
		uaaUsers    *uaa.Users
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeReader
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
		userManager = &DefaultManager{
			Client:     client,
			Cfg:        fakeReader,
			UAAMgr:     uaaFake,
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgReader:  orgFake,
			Peek:       false,
			LdapConfig: &config.LdapConfig{Origin: "ldap"},
		}
		userList = []cfclient.V3User{
			{
				Username: "hello",
				GUID:     "world",
			},
			{
				Username: "hello2",
				GUID:     "world2",
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
			client.ListV3SpaceRolesByGUIDAndTypeReturns(userList, nil)
			users, err := userManager.ListSpaceAuditors("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_AUDITOR))
		})
		It("Should remove orphaned users from ListSpaceAuditors", func() {
			userList = []cfclient.V3User{
				{
					Username: "hello",
					GUID:     "world",
				},
				{
					Username: "hello2",
					GUID:     "world2",
				},
				{
					Username: "orphaned_user",
					GUID:     "orphaned_user_guid",
				},
			}
			client.ListV3SpaceRolesByGUIDAndTypeReturns(userList, nil)
			users, err := userManager.ListSpaceManagers("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			Expect(len(users.OrphanedUsers())).Should(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_MANAGER))
		})
		It("Should error on ListSpaceAuditors", func() {
			client.ListV3SpaceRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceAuditors("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_AUDITOR))
		})
	})
	Context("Space Develpers", func() {
		It("Should succeed on ListSpaceDevelopers", func() {
			client.ListV3SpaceRolesByGUIDAndTypeReturns(userList, nil)
			users, err := userManager.ListSpaceDevelopers("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_DEVELOPER))
		})

		It("Should error on ListSpaceDevelopers", func() {
			client.ListV3SpaceRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceDevelopers("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_DEVELOPER))
		})
	})

	Context("Space Managers", func() {
		It("Should succeed on ListSpaceManagers", func() {
			client.ListV3SpaceRolesByGUIDAndTypeReturns(userList, nil)
			users, err := userManager.ListSpaceManagers("foo", uaaUsers)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(users.Users())).Should(Equal(2))
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_MANAGER))
		})

		It("Should error on ListSpaceManagers", func() {
			client.ListV3SpaceRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListSpaceManagers("foo", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3SpaceRolesByGUIDAndTypeCallCount()).To(Equal(1))
			spaceGUID, role := client.ListV3SpaceRolesByGUIDAndTypeArgsForCall(0)
			Expect(spaceGUID).To(Equal("foo"))
			Expect(role).To(Equal(SPACE_MANAGER))
		})
	})

	Context("ListOrgManager", func() {
		It("should succeed", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns([]cfclient.V3User{
				{
					Username: "test",
					GUID:     "test-guid",
				},
				{
					Username: "test2",
					GUID:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgManagers("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_MANAGER))
		})
		It("should error", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgManagers("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_MANAGER))
		})
	})

	Context("ListOrgAuditors", func() {
		It("should succeed", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns([]cfclient.V3User{
				{
					Username: "test",
					GUID:     "test-guid",
				},
				{
					Username: "test2",
					GUID:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgAuditors("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_AUDITOR))
		})
		It("should error", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgAuditors("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_AUDITOR))
		})
	})

	Context("ListOrgBillingManager", func() {
		It("should succeed", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns([]cfclient.V3User{
				{
					Username: "test",
					GUID:     "test-guid",
				},
				{
					Username: "test2",
					GUID:     "test2-guid",
				},
			}, nil)
			users, err := userManager.ListOrgBillingManagers("test-org-guid", uaaUsers)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(users.Users())).To(Equal(2))
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_BILLING_MANAGER))
		})
		It("should error", func() {
			client.ListV3OrganizationRolesByGUIDAndTypeReturns(nil, errors.New("error"))
			_, err := userManager.ListOrgBillingManagers("test-org-guid", uaaUsers)
			Expect(err).Should(HaveOccurred())
			Expect(client.ListV3OrganizationRolesByGUIDAndTypeCallCount()).To(Equal(1))
			orgGUID, role := client.ListV3OrganizationRolesByGUIDAndTypeArgsForCall(0)
			Expect(orgGUID).Should(Equal("test-org-guid"))
			Expect(role).To(Equal(ORG_BILLING_MANAGER))
		})
	})
})
