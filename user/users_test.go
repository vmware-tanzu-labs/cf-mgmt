package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	ldapfakes "github.com/pivotalservices/cf-mgmt/ldap/fakes"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	spacefakes "github.com/pivotalservices/cf-mgmt/space/fakes"
	"github.com/pivotalservices/cf-mgmt/uaa"
	uaafakes "github.com/pivotalservices/cf-mgmt/uaa/fakes"
	. "github.com/pivotalservices/cf-mgmt/user"
	"github.com/pivotalservices/cf-mgmt/user/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager *DefaultManager
		client      *fakes.FakeCFClient
		ldapFake    *ldapfakes.FakeManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeManager
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		ldapFake = new(ldapfakes.FakeManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeManager)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				Client:     client,
				Cfg:        fakeReader,
				UAAMgr:     uaaFake,
				LdapMgr:    ldapFake,
				SpaceMgr:   spaceFake,
				OrgMgr:     orgFake,
				Peek:       false,
				LdapConfig: &config.LdapConfig{Origin: "ldap"}}
		})

		Context("Success", func() {
			It("Should succeed on RemoveSpaceAuditor", func() {
				err := userManager.RemoveSpaceAuditor(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("test-guid"))
			})
			It("Should succeed on RemoveSpaceDeveloper", func() {
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceDeveloperArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("test-guid"))
			})
			It("Should succeed on RemoveSpaceManager", func() {
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "test-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceManagerArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("test-guid"))
			})

			It("Should succeed on AssociateSpaceAuditor", func() {
				client.AssociateSpaceAuditorReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should succeed on AssociateSpaceDeveloper", func() {
				client.AssociateSpaceDeveloperReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.AssociateSpaceDeveloperArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should succeed on AssociateSpaceManager", func() {
				client.AssociateSpaceManagerReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})
		})

		Context("SyncInternalUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers *uaa.Users
			BeforeEach(func() {
				uaaUsers = &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-user-guid"})
				uaaUsers.Add(uaa.User{Username: "test-existing", Origin: "uaa", GUID: "test-existing-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test-existing", Guid: "test-existing-id"},
				}, uaaUsers)
			})
			It("Should add internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})

			It("Should not add existing internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test-existing"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("user test2 doesn't exist in origin uaa, so must add internal user first"))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

		})

		Context("SyncSamlUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers *uaa.Users
			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}
				uaaUsers = &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test@test.com", Origin: "saml_origin", GUID: "test-id"})
				uaaUsers.Add(uaa.User{Username: "test@test2.com", Origin: "saml_origin", GUID: "test2-id"})
				roleUsers, _ = NewRoleUsers(
					[]cfclient.User{
						cfclient.User{Username: "test@test.com", Guid: "test-id"},
					},
					uaaUsers,
				)
			})
			It("Should add saml user to role", func() {
				updateUsersInput := UsersInput{
					SamlUsers: []string{"test@test2.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test2-id"))
				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test2-id"))
			})

			It("Should not add existing saml user to role", func() {
				updateUsersInput := UsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers.HasUser("test@test.com")).Should(BeFalse())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					SamlUsers: []string{"test@test3.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test@test3.com"))
				Expect(arg2).Should(Equal("test@test3.com"))
				Expect(arg3).Should(Equal("test@test3.com"))
				Expect(origin).Should(Equal("saml_origin"))
			})

			It("Should not error when create external user errors", func() {
				updateUsersInput := UsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, &uaa.Users{}, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				roleUsers := &RoleUsers{}
				roleUsers.AddUsers([]RoleUser{
					RoleUser{UserName: "test"},
				})
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test@test.com"})
				updateUsersInput := UsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
		})

		Context("Remove Users", func() {
			var roleUsers *RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test", Guid: "test-id"},
				}, uaaUsers)
			})

			It("Should remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorCallCount()).Should(Equal(1))

				spaceGUID, userGUID := client.RemoveSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test-id"))

			})

			It("Should not remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: false,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}
				client.RemoveSpaceAuditorReturns(errors.New("error"))
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorCallCount()).Should(Equal(1))
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
				Expect(client.RemoveSpaceAuditorCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceDeveloper", func() {
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceManager", func() {
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceAuditor", func() {
				client.AssociateSpaceAuditorReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceDeveloper", func() {
				client.AssociateSpaceDeveloperReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceManager", func() {
				client.AssociateSpaceManagerReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditor", func() {
				client.RemoveSpaceAuditorReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditor(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("user-guid"))
			})
			It("Should error on RemoveSpaceDeveloper", func() {
				client.RemoveSpaceDeveloperReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloper(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceDeveloperArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("user-guid"))
			})
			It("Should error on RemoveSpaceManager", func() {
				client.RemoveSpaceManagerReturns(errors.New("error"))
				err := userManager.RemoveSpaceManager(UsersInput{SpaceGUID: "foo"}, "bar", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceManagerCallCount()).To(Equal(1))
				spaceGUID, userGUID := client.RemoveSpaceManagerArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userGUID).To(Equal("user-guid"))
			})
			It("Should error on AssociateSpaceAuditor", func() {
				client.AssociateSpaceAuditorReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditor", func() {
				client.AssociateSpaceAuditorReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloper", func() {
				client.AssociateSpaceDeveloperReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloper", func() {
				client.AssociateSpaceDeveloperReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManager", func() {
				client.AssociateSpaceManagerReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManager", func() {
				client.AssociateSpaceManagerReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
			})
		})
		Context("AddUserToOrg", func() {
			It("should associate user", func() {
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})
		})
		Context("RemoveOrgAuditor", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgAuditorReturns(errors.New("error"))
				err := userManager.RemoveOrgAuditor(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgAuditorCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})
		})

		Context("RemoveOrgBillingManager", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgBillingManagerArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgBillingManagerReturns(errors.New("error"))
				err := userManager.RemoveOrgBillingManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgBillingManagerArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})
		})

		Context("RemoveOrgManager", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgManagerArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgManagerReturns(errors.New("error"))
				err := userManager.RemoveOrgManager(UsersInput{OrgGUID: "test-org-guid"}, "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgManagerCallCount()).To(Equal(1))
				orgGUID, userGUID := client.RemoveOrgManagerArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))

			})
		})

		Context("AssociateOrgAuditor", func() {
			It("Should succeed", func() {
				client.AssociateOrgAuditorReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgAuditorArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should fail", func() {
				client.AssociateOrgAuditorReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgAuditorArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})
			It("Should fail", func() {
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgAuditorReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
		})
		Context("AssociateOrgBillingManager", func() {
			It("Should succeed", func() {
				client.AssociateOrgBillingManagerReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgBillingManagerArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should fail", func() {
				client.AssociateOrgBillingManagerReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgBillingManagerArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgBillingManagerReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgManager", func() {
			It("Should succeed", func() {
				client.AssociateOrgManagerReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgManagerArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should fail", func() {
				client.AssociateOrgManagerReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))
				orgGUID, userGUID := client.AssociateOrgManagerArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))

				orgGUID, userGUID = client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})
			It("Should fail", func() {
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(1))

				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgManagerReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UsersInput{OrgGUID: "orgGUID"}, "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserCallCount()).To(Equal(0))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := &uaa.Users{}
				userMap.Add(uaa.User{Username: "test-user-guid", GUID: "test-user"})
				uaaFake.ListUsersReturns(userMap, nil)
				fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
					config.SpaceConfig{
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
					config.OrgConfig{
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
