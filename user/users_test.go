package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager *DefaultManager
		client      *fakes.FakeCFClient
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
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
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				Client:     client,
				Cfg:        fakeReader,
				UAAMgr:     uaaFake,
				LdapMgr:    ldapFake,
				SpaceMgr:   spaceFake,
				OrgReader:  orgFake,
				Peek:       false,
				LdapConfig: &config.LdapConfig{Origin: "ldap"}}

			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{}, nil)
		})

		Context("Success", func() {
			It("Should succeed on RemoveSpaceAuditor", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveSpaceAuditor(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})
			It("Should succeed on RemoveSpaceDeveloper", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})
			It("Should succeed on RemoveSpaceManager", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})

			It("Should succeed on AssociateSpaceAuditor", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
				spaceGUID, userGUID, roleType := client.CreateV3SpaceRoleArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(roleType).Should(Equal(SPACE_AUDITOR))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should succeed on AssociateSpaceDeveloper", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
				spaceGUID, userGUID, roleType := client.CreateV3SpaceRoleArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(roleType).Should(Equal(SPACE_DEVELOPER))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should succeed on AssociateSpaceManager", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})
		})

		Context("SyncInternalUsers", func() {
			var roleUsers *RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-user-guid"})
				uaaUsers.Add(uaa.User{Username: "test-existing", Origin: "uaa", GUID: "test-existing-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.V3User{
					{Username: "test-existing", GUID: "test-existing-id"},
				}, uaaUsers)

				userManager.UAAUsers = uaaUsers
			})
			It("Should add internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(role).Should(Equal(ORG_USER))

				spaceGUID, userGUID, roleType := client.CreateV3SpaceRoleArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(roleType).Should(Equal(SPACE_AUDITOR))

			})

			It("Should not add existing internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test-existing"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).Should(Equal(0))
				Expect(client.CreateV3SpaceRoleCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("user test2 doesn't exist in origin uaa, so must add internal user first"))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).Should(Equal(1))
				Expect(client.CreateV3SpaceRoleCallCount()).Should(Equal(0))
			})

		})

		Context("Remove Users", func() {
			var roleUsers *RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.V3User{
					{Username: "test", GUID: "test-id"},
				}, uaaUsers)
			})

			It("Should remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
					RoleUsers:   InitRoleUsers(),
				}
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).Should(Equal(1))

				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})

			It("Should not remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: false,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
					RoleUsers:   InitRoleUsers(),
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).Should(Equal(0))
			})

			It("Should skip users that match protected user pattern", func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "abcd_123_0919191", Origin: "uaa", GUID: "test-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.V3User{
					{Username: "abcd_123_0919191", GUID: "test-id"},
				}, uaaUsers)
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
					RoleUsers:   InitRoleUsers(),
				}

				fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
					ProtectedUsers: []string{"abcd_123_*"},
				}, nil)

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
					RoleUsers:   InitRoleUsers(),
				}
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).Should(Equal(1))
			})
		})

		Context("Peek", func() {
			BeforeEach(func() {
				userManager = &DefaultManager{
					Client:  client,
					Cfg:     nil,
					UAAMgr:  nil,
					LdapMgr: nil,
					Peek:    true}
			})
			It("Should succeed on RemoveSpaceAuditor", func() {
				err := userManager.RemoveSpaceAuditor(UsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceDeveloper", func() {
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceManager", func() {
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceAuditor", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceDeveloper", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceManager", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditor", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditor(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})
			It("Should error on RemoveSpaceDeveloper", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})
			It("Should error on RemoveSpaceManager", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).To(Equal("role-guid"))
			})
			It("Should error on AssociateSpaceAuditor", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditor", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloper", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloper", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManager", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(1))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManager", func() {
				client.CreateV3SpaceRoleReturns(&cfclient.V3Role{}, nil)
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID",
					RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3SpaceRoleCallCount()).To(Equal(0))
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
			})
		})
		Context("AddUserToOrg", func() {
			It("should associate user", func() {
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(role).To(Equal(ORG_USER))

			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})
		})
		Context("RemoveOrgAuditor", func() {
			It("should succeed", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})
		})

		Context("RemoveOrgBillingManager", func() {
			It("should succeed", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})
		})

		Context("RemoveOrgManager", func() {
			It("should succeed", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.ListV3RolesByQueryReturns([]cfclient.V3Role{
					{GUID: "role-guid"},
				}, nil)
				client.DeleteV3RoleReturns(errors.New("error"))
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid", RoleUsers: InitRoleUsers()}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.DeleteV3RoleCallCount()).To(Equal(1))
				roleGUID := client.DeleteV3RoleArgsForCall(0)
				Expect(roleGUID).Should(Equal("role-guid"))
			})
		})

		Context("AssociateOrgAuditor", func() {
			It("Should succeed", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(2))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))

				orgGUID, userGUID, role = client.CreateV3OrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_AUDITOR))
			})

			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})
			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
		})
		Context("AssociateOrgBillingManager", func() {
			It("Should succeed", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(2))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))

				orgGUID, userGUID, role = client.CreateV3OrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_BILLING_MANAGER))
			})

			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgManager", func() {
			It("Should succeed", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(2))
				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))

				orgGUID, userGUID, role = client.CreateV3OrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_MANAGER))
			})

			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})
			It("Should fail", func() {
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(1))

				orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(ORG_USER))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, nil)
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID", RoleUsers: InitRoleUsers()}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.CreateV3OrganizationRoleCallCount()).To(Equal(0))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := &uaa.Users{}
				userMap.Add(uaa.User{Username: "test-user-guid", GUID: "test-user"})
				uaaFake.ListUsersReturns(userMap, nil)
				fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
					{
						Space: "test-space",
						Org:   "test-org",
					},
				}, nil)
				spaceFake.FindSpaceReturns(cfclient.Space{
					Name:             "test-space",
					OrganizationGuid: "test-org-guid",
					Guid:             "test-space-guid",
				}, nil)
				userManager.LdapConfig = &config.LdapConfig{Enabled: false}
				err := userManager.UpdateSpaceUsers()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := &uaa.Users{}
				userMap.Add(uaa.User{Username: "test-user-guid", GUID: "test-user"})
				uaaFake.ListUsersReturns(userMap, nil)
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org: "test-org",
					},
				}, nil)
				orgFake.FindOrgReturns(cfclient.Org{
					Name: "test-org",
					Guid: "test-org-guid",
				}, nil)
				userManager.LdapConfig = &config.LdapConfig{Enabled: false}
				err := userManager.UpdateOrgUsers()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
