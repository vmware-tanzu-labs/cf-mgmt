package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	ldap "github.com/pivotalservices/cf-mgmt/ldap"
	ldapfakes "github.com/pivotalservices/cf-mgmt/ldap/fakes"
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
		userList    []cfclient.User
	)
	BeforeEach(func() {
		client = new(fakes.FakeCFClient)
		ldapFake = new(ldapfakes.FakeManager)
		uaaFake = new(uaafakes.FakeManager)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				Client:  client,
				Cfg:     nil,
				UAAMgr:  uaaFake,
				LdapMgr: ldapFake,
				Peek:    false}
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
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
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
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
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
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
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
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
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
				uaaUsers := make(map[string]string)
				uaaUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers).ShouldNot(HaveKey("test"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
				}
				err := userManager.SyncInternalUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("user test doesn't exist in cloud foundry, so must add internal user first"))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]string)
				uaaUsers["test"] = "test"
				updateUsersInput := UpdateUsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test@test.com"] = "test@test.com"
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test@test.com"] = "test@test.com"
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
				}
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleUsers).ShouldNot(HaveKey("test@test.com"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
				}
				uaaFake.CreateExternalUserReturns(errors.New("error"))
				err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).ShouldNot(HaveKey("test@test.com"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]string)
				uaaUsers["test@test.com"] = "test@test.com"
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test@test.com"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test_ldap"] = "test_ldap"
				updateUsersInput := UpdateUsersInput{
					LdapUsers:      []string{"test_ldap"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test_ldap"] = "test_ldap"
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SamlUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test_ldap"] = "test_ldap"
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
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
				uaaUsers := make(map[string]string)
				uaaUsers["test_ldap"] = "test_ldap"
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditorByUsername,
				}
				ldapFake.GetLdapUsersReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
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
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceDeveloperByUsername", func() {
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on RemoveSpaceManagerByUsername", func() {
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.RemoveSpaceManagerByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
			It("Should succeed on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).NotTo(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(0))
			})
		})
		Context("Error", func() {
			It("Should error on RemoveSpaceAuditorByUsername", func() {
				client.RemoveSpaceAuditorByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceAuditorByUsername("foo", "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceAuditorByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceAuditorByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceDeveloperByUsername", func() {
				client.RemoveSpaceDeveloperByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceDeveloperByUsername("foo", "bar")
				Expect(err).Should(HaveOccurred())
				Expect(client.RemoveSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				spaceGUID, userName := client.RemoveSpaceDeveloperByUsernameArgsForCall(0)
				Expect(spaceGUID).To(Equal("foo"))
				Expect(userName).To(Equal("bar"))
			})
			It("Should error on RemoveSpaceManagerByUsername", func() {
				client.RemoveSpaceManagerByUsernameReturns(errors.New("error"))
				err := userManager.RemoveSpaceManagerByUsername("foo", "bar")
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
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceAuditorByUsername", func() {
				client.AssociateSpaceAuditorByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceAuditorByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceDeveloperByUsername", func() {
				client.AssociateSpaceDeveloperByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceDeveloperByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceDeveloperByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, errors.New("error"))
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(1))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
			It("Should error on AssociateSpaceManagerByUsername", func() {
				client.AssociateSpaceManagerByUsernameReturns(cfclient.Space{}, nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.AssociateSpaceManagerByUsername("orgGUID", "spaceGUID", "userName")
				Expect(err).Should(HaveOccurred())
				Expect(client.AssociateSpaceManagerByUsernameCallCount()).To(Equal(0))
				Expect(client.AssociateOrgUserByUsernameCallCount()).To(Equal(1))
			})
		})
		Context("Add Users", func() {
			/*It("update ldap group users where users are not in uaac", func() {
				config := &config.LdapConfig{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID:      "my-space-guid",
					OrgGUID:        "my-org-guid",
					LdapGroupNames: []string{"ldap-group-name", "ldap-group-name-2"},
				}

				ldapGroupUsers := []ldap.User{ldap.User{
					UserDN: "user-dn",
					UserID: "user-id",
					Email:  "user@test.com",
				}}

				ldapGroupUsers2 := []ldap.User{ldap.User{
					UserDN: "user-dn2",
					UserID: "user-id2",
					Email:  "user2@test.com",
				}}


				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)
				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name-2").Return(ldapGroupUsers2, nil)

				mockUaac.EXPECT().CreateExternalUser("user-id", "user@test.com", "user-dn", "ldap").Return(nil)
				mockUaac.EXPECT().CreateExternalUser("user-id2", "user2@test.com", "user-dn2", "ldap").Return(nil)
				mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToOrg("user-id2", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-id", "my-role", "my-space-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-id2", "my-role", "my-space-guid").Return(nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
				_, ok := uaacUsers["user-id"]
				Ω(ok).Should(BeTrue())

				_, ok = uaacUsers["user-id2"]
				Ω(ok).Should(BeTrue())
			})*/
			/*It("update ldap group users where users are in uaac", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user-id"] = "user-id"
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID:      "my-space-guid",
					OrgGUID:        "my-org-guid",
					Role:           "my-role",
					LdapGroupNames: []string{"ldap-group-name"},
				}

				ldapGroupUsers := []ldap.User{ldap.User{
					UserDN: "user-dn",
					UserID: "user-id",
					Email:  "user@test.com",
				}}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

				mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-id", "my-role", "my-space-guid").Return(nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
				_, ok := uaacUsers["user-id"]
				Ω(ok).Should(BeTrue())
			})

			It("update ldap group users where users are in uaac and already in space", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user-id"] = "user-id"
				spaceUsers := make(map[string]string)
				spaceUsers["user-id"] = "user-id"
				updateUsersInput := UpdateUsersInput{
					SpaceGUID:      "my-space-guid",
					OrgGUID:        "my-org-guid",
					Role:           "my-role",
					LdapGroupNames: []string{"ldap-group-name"},
				}

				ldapGroupUsers := []ldap.User{ldap.User{
					UserDN: "user-dn",
					UserID: "user-id",
					Email:  "user@test.com",
				}}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
				_, ok := uaacUsers["user-id"]
				Ω(ok).Should(BeTrue())
			})

			It("update other origin users where users are in uaac and already in space", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "other",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user@test.com"] = "user@test.com"
				spaceUsers := make(map[string]string)
				spaceUsers["user@test.com"] = "user@test.com"
				updateUsersInput := UpdateUsersInput{
					SpaceGUID:      "my-space-guid",
					OrgGUID:        "my-org-guid",
					Role:           "my-role",
					LdapGroupNames: []string{"ldap-group-name"},
				}

				ldapGroupUsers := []ldap.User{ldap.User{
					UserDN: "user-dn",
					UserID: "user-id",
					Email:  "user@test.com",
				}}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
				_, ok := uaacUsers["user@test.com"]
				Ω(ok).Should(BeTrue())
			})

			It("update ldap users where users are not in uaac", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
					LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUser(config, "ldap-user-1").Return(&ldap.User{
					UserDN: "user-1-dn",
					UserID: "user-1-id",
					Email:  "user1@test.com",
				}, nil)
				mockLdap.EXPECT().GetUser(config, "ldap-user-2").Return(&ldap.User{
					UserDN: "user-2-dn",
					UserID: "user-2-id",
					Email:  "user2@test.com",
				}, nil)

				mockUaac.EXPECT().CreateExternalUser("user-1-id", "user1@test.com", "user-1-dn", "ldap").Return(nil)
				mockCloudController.EXPECT().AddUserToOrg("user-1-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-1-id", "my-role", "my-space-guid").Return(nil)

				mockUaac.EXPECT().CreateExternalUser("user-2-id", "user2@test.com", "user-2-dn", "ldap").Return(nil)
				mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-2-id", "my-role", "my-space-guid").Return(nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
				_, ok := uaacUsers["user-1-id"]
				Ω(ok).Should(BeTrue())
				_, ok = uaacUsers["user-2-id"]
				Ω(ok).Should(BeTrue())
			})

			It("update ldap users where users are in uaac", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user-1-id"] = "user-1-id"
				uaacUsers["user-2-id"] = "user-2-id"
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
					LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockLdap.EXPECT().GetUser(config, "ldap-user-1").Return(&ldap.User{
					UserDN: "user-1-dn",
					UserID: "user-1-id",
					Email:  "user1@test.com",
				}, nil)
				mockLdap.EXPECT().GetUser(config, "ldap-user-2").Return(&ldap.User{
					UserDN: "user-2-dn",
					UserID: "user-2-id",
					Email:  "user2@test.com",
				}, nil)

				mockCloudController.EXPECT().AddUserToOrg("user-1-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-1-id", "my-role", "my-space-guid").Return(nil)

				mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-2-id", "my-role", "my-space-guid").Return(nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())

				Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
				_, ok := uaacUsers["user-1-id"]
				Ω(ok).Should(BeTrue())
				_, ok = uaacUsers["user-2-id"]
				Ω(ok).Should(BeTrue())
			})

			It("update users where users are in uaac", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user-1"] = "user-1"
				uaacUsers["user-2"] = "user-2"
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
					Users:     []string{"user-1", "user-2"},
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockCloudController.EXPECT().AddUserToOrg("user-1", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-1", "my-role", "my-space-guid").Return(nil)

				mockCloudController.EXPECT().AddUserToOrg("user-2", "my-org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("user-2", "my-role", "my-space-guid").Return(nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())

				Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
				_, ok := uaacUsers["user-1"]
				Ω(ok).Should(BeTrue())
				_, ok = uaacUsers["user-2"]
				Ω(ok).Should(BeTrue())
			})

			It("update users where users are in uaac and in a space", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				uaacUsers["user-1"] = "user-1"
				uaacUsers["user-2"] = "user-2"
				spaceUsers := make(map[string]string)
				spaceUsers["user-1"] = "asfdsdf-1"
				spaceUsers["user-2"] = "asdfsaf-2"
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
					Users:     []string{"USER-1", "user-2"},
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)

				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())

				Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
				_, ok := uaacUsers["user-1"]
				Ω(ok).Should(BeTrue())
				_, ok = uaacUsers["user-2"]
				Ω(ok).Should(BeTrue())
			})

			It("update users where users are not in uaac", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				spaceUsers := make(map[string]string)
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
					Users:     []string{"user-1"},
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)

				Ω(userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)).ShouldNot(Succeed())
				Ω(len(uaacUsers)).Should(BeEquivalentTo(0))
			})

			It("remove users that in space but not in config", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				spaceUsers := make(map[string]string)
				spaceUsers["cwashburn"] = "cwashburn"
				spaceUsers["cwashburn1"] = "cwashburn1"
				spaceUsers["cwashburn2"] = "cwashburn2"
				updateUsersInput := UpdateUsersInput{
					SpaceGUID:   "my-space-guid",
					OrgGUID:     "my-org-guid",
					Role:        "my-role",
					RemoveUsers: true,
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn", "my-role").Return(nil)
				mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn1", "my-role").Return(nil)
				mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn2", "my-role").Return(nil)
				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
			})

			It("remove orphaned LDAP users while leaving existing group members - GH issue 33", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "https://saml.example.com",
				}

				uaacUsers := make(map[string]string)
				uaacUsers["chris.a.washburn@example.com"] = "cwashburn-uaac-guid"
				uaacUsers["joe.h.fitzy@example.com"] = "jfitzy-uaac-guid"
				uaacUsers["alex.j.smith@example.com"] = "asmith-uaac-guid" // <-- user in uaac, but not ldap group

				spaceUsers := make(map[string]string)
				spaceUsers["chris.a.washburn@example.com"] = "cwashburn-space-user-guid"
				spaceUsers["joe.h.fitzy@example.com"] = "jfitzy-space-user-guid"
				spaceUsers["alex.j.smith@example.com"] = "asmith-space-user-guid" // <-- user in space, but not ldap group

				updateUsersInput := UpdateUsersInput{
					SpaceName:      "space-name",
					SpaceGUID:      "space-guid",
					OrgName:        "org-name",
					OrgGUID:        "org-guid",
					Role:           "space-role-name",
					LdapGroupNames: []string{"ldap-group-name"},
					RemoveUsers:    true,
				}

				ldapGroupUsers := []ldap.User{ldap.User{
					UserDN: "CN=Washburn, Chris,OU=End Users,OU=Accounts,DC=add,DC=example,DC=com",
					UserID: "u-cwashburn",
					Email:  "Chris.A.Washburn@example.com",
				}, ldap.User{
					UserDN: "CN=Fitzy, Joe,OU=End Users,OU=Accounts,DC=ad,DC=example,DC=com",
					UserID: "u-jfitzy",
					Email:  "Joe.H.Fitzy@example.com",
				}}

				mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

				mockCloudController.EXPECT().GetCFUsers("space-guid", "spaces", "space-role-name").Return(spaceUsers, nil)
				mockCloudController.EXPECT().RemoveCFUserByUserName("space-guid", "spaces", "alex.j.smith@example.com", "space-role-name").Return(nil)
				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
			})

			It("don't remove users that in space but not in config", func() {
				config := &l.Config{
					Enabled: true,
					Origin:  "ldap",
				}
				uaacUsers := make(map[string]string)
				spaceUsers := make(map[string]string)
				spaceUsers["cwashburn"] = "cwashburn"
				spaceUsers["cwashburn1"] = "cwashburn1"
				spaceUsers["cwashburn2"] = "cwashburn2"
				updateUsersInput := UpdateUsersInput{
					SpaceGUID: "my-space-guid",
					OrgGUID:   "my-org-guid",
					Role:      "my-role",
				}

				mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
			})

			It("adding users to uaac based on saml", func() {
				config := &l.Config{
					Enabled: false,
					Origin:  "https://saml.example.com",
				}

				uaacUsers := make(map[string]string)
				uaacUsers["chris.a.washburn@example.com"] = "cwashburn-uaac-guid"
				uaacUsers["joe.h.fitzy@example.com"] = "jfitzy-uaac-guid"

				spaceUsers := make(map[string]string)
				spaceUsers["chris.a.washburn@example.com"] = "cwashburn-space-user-guid"
				spaceUsers["joe.h.fitzy@example.com"] = "jfitzy-space-user-guid"

				updateUsersInput := UpdateUsersInput{
					SpaceName:   "space-name",
					SpaceGUID:   "space-guid",
					OrgName:     "org-name",
					OrgGUID:     "org-guid",
					Role:        "space-role-name",
					SamlUsers:   []string{"chris.a.washburn@example.com", "joe.h.fitzy@example.com", "test@test.com"},
					RemoveUsers: true,
				}

				mockCloudController.EXPECT().GetCFUsers("space-guid", "spaces", "space-role-name").Return(spaceUsers, nil)
				mockUaac.EXPECT().CreateExternalUser("test@test.com", "test@test.com", "test@test.com", "https://saml.example.com").Return(nil)
				mockCloudController.EXPECT().AddUserToOrg("test@test.com", "org-guid").Return(nil)
				mockCloudController.EXPECT().AddUserToSpaceRole("test@test.com", "space-role-name", "space-guid").Return(nil)
				err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
				Ω(err).Should(BeNil())
				Ω(uaacUsers).Should(HaveKey("test@test.com"))
			})*/
		})
	})
})
