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
			It("Should succeed on RemoveSpaceAuditorByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceDeveloperByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})
			It("Should succeed on RemoveSpaceManagerByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceManagerByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should succeed on AssociateSpaceAuditorByUsernameAndOrigin", func() {
				client.AssociateSpaceAuditorByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.AssociateSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should succeed on AssociateSpaceDeveloperByUsernameAndOrigin", func() {
				client.AssociateSpaceDeveloperByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.AssociateSpaceDeveloperByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should succeed on AssociateSpaceManagerByUsernameAndOrigin", func() {
				client.AssociateSpaceManagerByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})
		})

		Context("SyncInternalUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers map[string]uaa.User
			BeforeEach(func() {
				uaaUsers = make(map[string]uaa.User)
				uaaUsers["test"] = uaa.User{Username: "test", Origin: "uaa"}
				uaaUsers["test-id"] = uaa.User{Username: "test", Origin: "uaa"}
				uaaUsers["test-existing-id"] = uaa.User{Username: "test-existing", Origin: "uaa"}
				uaaUsers["test-existing"] = uaa.User{Username: "test-existing", Origin: "uaa"}
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test-existing", Guid: "test-existing-id"},
				}, uaaUsers)
			})
			It("Should add internal user to role", func() {
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))

				spaceGUID, userName, origin := client.AssociateSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should not add existing internal user to role", func() {
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test-existing"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				updateUsersInput := UpdateUsersInput{
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
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})

		})

		Context("SyncSamlUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers map[string]uaa.User
			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}
				uaaUsers = map[string]uaa.User{
					"test-id":       uaa.User{Username: "test@test.com", Origin: "saml_origin"},
					"test@test.com": uaa.User{Username: "test@test.com", Origin: "saml_origin"},
				}
				roleUsers, _ = NewRoleUsers(
					[]cfclient.User{
						cfclient.User{Username: "test@test.com", Guid: "test-id"},
					},
					uaaUsers,
				)
			})
			It("Should add saml user to role", func() {
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test2.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test@test2.com"))
				Expect(origin).Should(Equal("saml_origin"))

				spaceGUID, userName, origin := client.AssociateSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test@test2.com"))
				Expect(origin).Should(Equal("saml_origin"))
			})

			It("Should not add existing saml user to role", func() {
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers.HasUser("test@test.com")).Should(BeFalse())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test2.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test@test2.com"))
				Expect(arg2).Should(Equal("test@test2.com"))
				Expect(arg3).Should(Equal("test@test2.com"))
				Expect(origin).Should(Equal("saml_origin"))
			})

			It("Should not error when create external user errors", func() {
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				uaaFake.CreateExternalUserReturns(errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, nil, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				roleUsers := &RoleUsers{}
				roleUsers.AddUsers([]RoleUser{
					RoleUser{UserName: "test"},
				})
				uaaUsers := make(map[string]uaa.User)
				uaaUsers["test@test.com"] = uaa.User{Username: "test@test.com"}
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})
		})

		Context("Remove Users", func() {
			var roleUsers *RoleUsers
			BeforeEach(func() {
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test", Guid: "test-id"},
				}, map[string]uaa.User{
					"test-id": uaa.User{Username: "test", Origin: "uaa"},
				})
			})

			It("Should remove users", func() {
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))

				spaceGUID, userName, origin := client.RemoveSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should not remove users", func() {
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: false,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}
				client.RemoveSpaceAuditorByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))
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
			It("Should succeed on RemoveSpaceAuditorByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceManagerByUsernameAndOrigin", func() {
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameAndOriginCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceAuditorByUsernameAndOrigin", func() {
				client.AssociateSpaceAuditorByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceDeveloperByUsernameAndOrigin", func() {
				client.AssociateSpaceDeveloperByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceManagerByUsernameAndOrigin", func() {
				client.AssociateSpaceManagerByUsernameAndOriginReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditorByUsernameAndOrigin", func() {
				client.RemoveSpaceAuditorByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})
			It("Should error on RemoveSpaceDeveloperByUsernameAndOrigin", func() {
				client.RemoveSpaceDeveloperByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceDeveloperByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})
			It("Should error on RemoveSpaceManagerByUsernameAndOrigin", func() {
				client.RemoveSpaceManagerByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameAndOriginCallCount()).To(Equal(1))
				spaceGUID, userName, origin := client.RemoveSpaceManagerByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
				Expect(origin).Should(Equal("uaa"))
			})
			It("Should error on AssociateSpaceAuditorByUsernameAndOrigin", func() {
				client.AssociateSpaceAuditorByUsernameAndOriginReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditorByUsernameAndOrigin", func() {
				client.AssociateSpaceAuditorByUsernameAndOriginReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsernameAndOrigin", func() {
				client.AssociateSpaceDeveloperByUsernameAndOriginReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsernameAndOrigin", func() {
				client.AssociateSpaceDeveloperByUsernameAndOriginReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsernameAndOrigin", func() {
				client.AssociateSpaceManagerByUsernameAndOriginReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsernameAndOrigin", func() {
				client.AssociateSpaceManagerByUsernameAndOriginReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
			})
		})
		Context("AddUserToOrg", func() {
			It("should associate user", func() {
				err := userManager.AddUserToOrg("test", "uaa", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.AddUserToOrg("test", "uaa", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AddUserToOrg("test", "uaa", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})
		})
		Context("RemoveOrgAuditorByUsernameAndOrigin", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgAuditorByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameAndOriginCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgAuditorByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgAuditorByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})
		})

		Context("RemoveOrgBillingManagerByUsernameAndOrigin", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgBillingManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgBillingManagerByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgBillingManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})
		})

		Context("RemoveOrgManagerByUsernameAndOrigin", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameAndOriginCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgManagerByUsernameAndOriginReturns(errors.New("error"))
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test", "uaa")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.RemoveOrgManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
				Expect(origin).Should(Equal("uaa"))
			})
		})

		Context("AssociateOrgAuditorByUsernameAndOrigin", func() {
			It("Should succeed", func() {
				client.AssociateOrgAuditorByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgAuditorByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should fail", func() {
				client.AssociateOrgAuditorByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgAuditorByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgAuditorByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
		})
		Context("AssociateOrgBillingManagerByUsernameAndOrigin", func() {
			It("Should succeed", func() {
				client.AssociateOrgBillingManagerByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgBillingManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should fail", func() {
				client.AssociateOrgBillingManagerByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgBillingManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgBillingManagerByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgManagerByUsernameAndOrigin", func() {
			It("Should succeed", func() {
				client.AssociateOrgManagerByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should fail", func() {
				client.AssociateOrgManagerByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameAndOriginCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgManagerByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))

				orgGUID, userName, origin = client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(1))

				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
				Expect(origin).Should(Equal("uaa"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgManagerByUsernameAndOriginReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName", "uaa")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameAndOriginCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).To(Equal(0))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := make(map[string]uaa.User)
				userMap["test-user"] = uaa.User{Username: "test-user-guid"}
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
				userMap := make(map[string]uaa.User)
				userMap["test-user"] = uaa.User{Username: "test-user-guid"}
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
