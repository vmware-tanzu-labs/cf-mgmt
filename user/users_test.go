package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	ldap "github.com/pivotalservices/cf-mgmt/ldap"
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
		userList    []cfclient.User
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
				Client:   client,
				Cfg:      fakeReader,
				UAAMgr:   uaaFake,
				LdapMgr:  ldapFake,
				SpaceMgr: spaceFake,
				OrgMgr:   orgFake,
				Peek:     false}
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
		})

		Context("Success", func() {
			It("Should succeed on RemoveSpaceAuditorByUsername", func() {
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceManagerByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on ListSpaceAuditors", func() {
				client.ListSpaceAuditorsReturns(userList, nil)
				users, err := userManager.ListSpaceAuditors("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should succeed on ListSpaceDevelopers", func() {
				client.ListSpaceDevelopersReturns(userList, nil)
				users, err := userManager.ListSpaceDevelopers("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should succeed on ListSpaceManagers", func() {
				client.ListSpaceManagersReturns(userList, nil)
				users, err := userManager.ListSpaceManagers("foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(len(users)).Should(Equal(2))
				Expect(users).Should(HaveKeyWithValue("hello", "world"))
				Expect(users).Should(HaveKeyWithValue("hello2", "world2"))
				Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceManagersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})

			It("Should succeed on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.AssociateSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should succeed on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.AssociateSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("spaceGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should succeed on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})
		})
		Context("UpdateUserInfo", func() {

			It("ldap origin with email", func() {
				ldapFake.LdapConfigReturns(&config.LdapConfig{Origin: "ldap"})
				userInfo := userManager.UpdateUserInfo(ldap.User{
					Email:  "test@test.com",
					UserID: "testUser",
					UserDN: "testUserDN",
				})
				Expect(userInfo.Email).Should(Equal("test@test.com"))
				Expect(userInfo.UserDN).Should(Equal("testUserDN"))
				Expect(userInfo.UserID).Should(Equal("testuser"))
			})

			It("ldap origin without email email", func() {
				ldapFake.LdapConfigReturns(&config.LdapConfig{Origin: "ldap"})
				userInfo := userManager.UpdateUserInfo(ldap.User{
					Email:  "",
					UserID: "testUser",
					UserDN: "testUserDN",
				})
				Expect(userInfo.Email).Should(Equal("testuser@user.from.ldap.cf"))
				Expect(userInfo.UserDN).Should(Equal("testUserDN"))
				Expect(userInfo.UserID).Should(Equal("testuser"))
			})

			It("non ldap origin should return same email", func() {
				ldapFake.LdapConfigReturns(&config.LdapConfig{Origin: ""})
				userInfo := userManager.UpdateUserInfo(ldap.User{
					Email:  "test@test.com",
					UserID: "testUser",
					UserDN: "testUserDN",
				})
				Expect(userInfo.Email).Should(Equal("test@test.com"))
				Expect(userInfo.UserDN).Should(Equal("test@test.com"))
				Expect(userInfo.UserID).Should(Equal("test@test.com"))
			})
		})

		Context("SyncInternalUsers", func() {
			It("Should add internal user to role", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test"] = &uaa.User{UserName: "test"}
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test"))

				spaceGUID, userName := client.AssociateSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("Should not add existing internal user to role", func() {
				roleUsers := make(map[string]string)
				roleUsers["test"] = "test"
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test"] = &uaa.User{UserName: "test"}
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers).ShouldNot(HaveKey("test"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("user test doesn't exist in cloud foundry, so must add internal user first"))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test"] = &uaa.User{UserName: "test"}
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})

		})

		Context("SyncSamlUsers", func() {
			BeforeEach(func() {
				ldapFake.LdapConfigReturns(&config.LdapConfig{Origin: "saml_origin"})
			})
			It("Should add saml user to role", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test@test.com"] = &uaa.User{UserName: "test@test.com"}
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test@test.com"))

				spaceGUID, userName := client.AssociateSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test@test.com"))
			})

			It("Should not add existing saml user to role", func() {
				roleUsers := make(map[string]string)
				roleUsers["test@test.com"] = "test@test.com"
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test@test.com"] = &uaa.User{UserName: "test@test.com"}
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers).ShouldNot(HaveKey("test@test.com"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).Should(HaveKey("test@test.com"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test@test.com"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("test@test.com"))
				Expect(origin).Should(Equal("saml_origin"))
			})

			It("Should not error when create external user errors", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				uaaFake.CreateExternalUserReturns(errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).ShouldNot(HaveKey("test@test.com"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test@test.com"] = &uaa.User{UserName: "test@test.com"}
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
		})

		Context("SyncLdapUsers", func() {
			BeforeEach(func() {
				ldapFake.LdapConfigReturns(&config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				})
			})
			It("Should add ldap user to role", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapUsers:      []string{"test_ldap"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				ldapFake.GetLdapUsersReturns([]ldap.User{
					ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
				}, nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(1))
				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test_ldap"))

				spaceGUID, userName := client.AssociateSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test_ldap"))
			})

			It("Should not add existing ldap user to role", func() {
				roleUsers := make(map[string]string)
				roleUsers["test_ldap"] = "test_ldap"
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetLdapUsersReturns([]ldap.User{
					ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
				}, nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers).ShouldNot(HaveKey("test_ldap"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetLdapUsersReturns([]ldap.User{
					ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
				}, nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).Should(HaveKey("test_ldap"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("ldap_test_dn"))
				Expect(origin).Should(Equal("ldap"))
			})

			It("Should not error when create external user errors", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetLdapUsersReturns([]ldap.User{
					ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
				}, nil)
				uaaFake.CreateExternalUserReturns(errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).ShouldNot(HaveKey("test_ldap"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetLdapUsersReturns([]ldap.User{
					ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
				}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetLdapUsersReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
		})

		Context("Remove Users", func() {
			It("Should remove users", func() {
				roleUsers := make(map[string]string)
				roleUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).Should(Equal(1))

				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("Should not remove users", func() {
				roleUsers := make(map[string]string)
				roleUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: false,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				roleUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  userManager.RemoveSpaceAuditor,
				}
				client.RemoveSpaceAuditorByUsernameReturns(errors.New("error"))
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).Should(Equal(1))
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
			It("Should succeed on RemoveSpaceAuditorByUsername", func() {
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditorByUsername", func() {
				client.RemoveSpaceAuditorByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditor(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceDeveloperByUsername", func() {
				client.RemoveSpaceDeveloperByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloper(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceManagerByUsername", func() {
				client.RemoveSpaceManagerByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceManager(UpdateUsersInput{SpaceGUID: "foo"}, "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceManagerByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on ListSpaceAuditors", func() {
				client.ListSpaceAuditorsReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceAuditors("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceAuditorsCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceAuditorsArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on ListSpaceDevelopers", func() {
				client.ListSpaceDevelopersReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceDevelopers("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceDevelopersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceDevelopersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on ListSpaceManagers", func() {
				client.ListSpaceManagersReturns(nil, errors.New("error"))
				_, err := userManager.ListSpaceManagers("foo")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListSpaceManagersCallCount()).To(Equal(1))
				spaceGUID := client.ListSpaceManagersArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
			})
			It("Should error on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditor(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloper(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceManager(UpdateUsersInput{
					OrgGUID:   "orgGUID",
					SpaceGUID: "spaceGUID"}, "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
		})
		Context("AddUserToOrg", func() {
			It("should associate user", func() {
				err := userManager.AddUserToOrg("test", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.AddUserToOrg("test", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AddUserToOrg("test", UpdateUsersInput{OrgGUID: "test-org-guid"})
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})
		})
		Context("RemoveOrgAuditorByUsername", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgAuditorByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgAuditorByUsernameReturns(errors.New("error"))
				err := userManager.RemoveOrgAuditor(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgAuditorByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgAuditorByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})
		})

		Context("RemoveOrgBillingManagerByUsername", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgBillingManagerByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgBillingManagerByUsernameReturns(errors.New("error"))
				err := userManager.RemoveOrgBillingManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgBillingManagerByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgBillingManagerByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})
		})

		Context("RemoveOrgManagerByUsername", func() {
			It("should succeed", func() {
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgManagerByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
			})

			It("should peek", func() {
				userManager.Peek = true
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameCallCount()).To(Equal(0))
			})

			It("should error", func() {
				client.RemoveOrgManagerByUsernameReturns(errors.New("error"))
				err := userManager.RemoveOrgManager(UpdateUsersInput{OrgGUID: "test-org-guid"}, "test")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveOrgManagerByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.RemoveOrgManagerByUsernameArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userName).Should(Equal("test"))
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
				users, err := userManager.ListOrgAuditors("test-org-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).To(Equal(2))
				Expect(client.ListOrgAuditorsCallCount()).To(Equal(1))
				orgGUID := client.ListOrgAuditorsArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
			})
			It("should error", func() {
				client.ListOrgAuditorsReturns(nil, errors.New("error"))
				_, err := userManager.ListOrgAuditors("test-org-guid")
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
				users, err := userManager.ListOrgBillingManagers("test-org-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).To(Equal(2))
				Expect(client.ListOrgBillingManagersCallCount()).To(Equal(1))
				orgGUID := client.ListOrgBillingManagersArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
			})
			It("should error", func() {
				client.ListOrgBillingManagersReturns(nil, errors.New("error"))
				_, err := userManager.ListOrgBillingManagers("test-org-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListOrgBillingManagersCallCount()).To(Equal(1))
				orgGUID := client.ListOrgBillingManagersArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
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
				users, err := userManager.ListOrgManagers("test-org-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(users)).To(Equal(2))
				Expect(client.ListOrgManagersCallCount()).To(Equal(1))
				orgGUID := client.ListOrgManagersArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
			})
			It("should error", func() {
				client.ListOrgManagersReturns(nil, errors.New("error"))
				_, err := userManager.ListOrgManagers("test-org-guid")
				Expect(err).Should(HaveOccurred())
				Expect(client.ListOrgManagersCallCount()).To(Equal(1))
				orgGUID := client.ListOrgManagersArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
			})
		})

		Context("AssociateOrgAuditorByUsername", func() {
			It("Should succeed", func() {
				client.AssociateOrgAuditorByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgAuditorByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgAuditorByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgAuditorByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgAuditorByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgAuditor(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
		})
		Context("AssociateOrgBillingManagerByUsername", func() {
			It("Should succeed", func() {
				client.AssociateOrgBillingManagerByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgBillingManagerByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgBillingManagerByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgBillingManagerByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgBillingManagerByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgBillingManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgBillingManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgManagerByUsername", func() {
			It("Should succeed", func() {
				client.AssociateOrgManagerByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgManagerByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgManagerByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
				orgGUID, userName := client.AssociateOrgManagerByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))

				orgGUID, userName = client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should fail", func() {
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).To(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))

				orgGUID, userName := client.AssociateOrgUserByUsernameArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userName).To(Equal("userName"))
			})

			It("Should peek", func() {
				userManager.Peek = true
				client.AssociateOrgManagerByUsernameReturns(cfclient.Org{}, nil)
				err := userManager.AssociateOrgManager(UpdateUsersInput{OrgGUID: "orgGUID"}, "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateOrgManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := make(map[string]*uaa.User)
				userMap["test-user"] = &uaa.User{UserName: "test-user-guid"}
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
				ldapFake.LdapConfigReturns(&config.LdapConfig{Enabled: false})
				err := userManager.UpdateSpaceUsers()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := make(map[string]*uaa.User)
				userMap["test-user"] = &uaa.User{UserName: "test-user-guid"}
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
				ldapFake.LdapConfigReturns(&config.LdapConfig{Enabled: false})
				err := userManager.UpdateOrgUsers()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
