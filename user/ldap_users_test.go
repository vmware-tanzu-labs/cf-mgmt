package user_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	ldap "github.com/vmwarepivotallabs/cf-mgmt/ldap"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager     *DefaultManager
		fakeRoleClient  *fakes.FakeCFRoleClient
		fakeUserClient  *fakes.FakeCFUserClient
		fakeSpaceClient *fakes.FakeCFSpaceClient
		fakeJobClient   *fakes.FakeCFJobClient
		ldapFake        *fakes.FakeLdapManager
		uaaFake         *uaafakes.FakeManager
		fakeReader      *configfakes.FakeReader
		spaceFake       *spacefakes.FakeManager
		orgFake         *orgfakes.FakeReader
	)
	BeforeEach(func() {
		fakeRoleClient = new(fakes.FakeCFRoleClient)
		fakeUserClient = new(fakes.FakeCFUserClient)
		fakeSpaceClient = new(fakes.FakeCFSpaceClient)
		fakeJobClient = new(fakes.FakeCFJobClient)
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				RoleClient:  fakeRoleClient,
				UserClient:  fakeUserClient,
				SpaceClient: fakeSpaceClient,
				JobClient:   fakeJobClient,
				Cfg:         fakeReader,
				UAAMgr:      uaaFake,
				LdapMgr:     ldapFake,
				SpaceMgr:    spaceFake,
				OrgReader:   orgFake,
				Peek:        false,
				LdapConfig:  &config.LdapConfig{Origin: "ldap"}}
		})
		Context("SyncLdapUsers", func() {
			var roleUsers *RoleUsers

			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = NewRoleUsers([]*resource.User{
					{Username: "test_ldap", GUID: "test_ldap-id"},
				}, uaaUsers)
				userManager.UAAUsers = uaaUsers
			})
			It("Should add ldap user to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap2"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
					RoleUsers:      InitRoleUsers(),
				}

				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(1))
				_, orgGUID, userGUID, role := fakeRoleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, spaceGUID, userGUID, roleType := fakeRoleClient.CreateSpaceRoleArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))
			})

			It("Should add ldap user to role", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = NewRoleUsers([]*resource.User{
					{Username: "test_ldap", GUID: "test_ldap-id"},
				}, uaaUsers)

				userManager.UAAUsers = uaaUsers
				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap2"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
					RoleUsers:      InitRoleUsers(),
				}

				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(1))
				_, orgGUID, userGUID, role := fakeRoleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, spaceGUID, userGUID, roleType := fakeRoleClient.CreateSpaceRoleArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))
			})

			It("Should add ldap group member to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
					RoleUsers:      InitRoleUsers(),
				}

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap2",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(1))
				_, orgGUID, userGUID, role := fakeRoleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, spaceGUID, userGUID, roleType := fakeRoleClient.CreateSpaceRoleArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))
			})

			It("Should not add existing ldap user to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(0))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap_new"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap_new",
						Email:  "test@test.com",
					},
					nil)
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap_new"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("ldap_test_dn"))
				Expect(origin).Should(Equal("ldap"))
			})

			It("Should not error when create external user errors", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)
				uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(userManager.UAAUsers.GetByName("test_ldap3"))).Should(Equal(0))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)
				fakeRoleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("User test_ldap3 with origin ldap: error"))
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}

				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {

				updateUsersInput := UsersInput{
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
					RoleUsers:      InitRoleUsers(),
				}
				ldapFake.GetUserDNsReturns([]string{"cn=test_ldap2"}, nil)
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserDNsCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
				Expect(ldapFake.GetUserByDNCallCount()).Should(Equal(0))
			})
			It("Should return error", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
					RoleUsers: InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(0))
				Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(0))
			})
		})
		Context("UpdateUserInfo", func() {

			It("ldap origin with email", func() {
				userManager.LdapConfig.Origin = "ldap"
				userInfo := userManager.UpdateUserInfo(ldap.User{
					Email:  "test@test.com",
					UserID: "testUser",
					UserDN: "testUserDN",
				})
				Expect(userInfo.Email).Should(Equal("test@test.com"))
				Expect(userInfo.UserDN).Should(Equal("testUserDN"))
				Expect(userInfo.UserID).Should(Equal("testuser"))
			})

			It("ldap origin without email", func() {
				userManager.LdapConfig.Origin = "ldap"
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
				userManager.LdapConfig.Origin = "foo"
				userManager.LdapConfig = &config.LdapConfig{Origin: ""}
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
	})
})
