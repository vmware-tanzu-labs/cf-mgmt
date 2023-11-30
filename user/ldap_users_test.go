package user_test

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/ldap"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	rolefakes "github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager *DefaultManager
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeReader
		roleMgrFake *rolefakes.FakeManager
	)
	BeforeEach(func() {
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
		roleMgrFake = new(rolefakes.FakeManager)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				Cfg:        fakeReader,
				UAAMgr:     uaaFake,
				LdapMgr:    ldapFake,
				SpaceMgr:   spaceFake,
				OrgReader:  orgFake,
				Peek:       false,
				RoleMgr:    roleMgrFake,
				LdapConfig: &config.LdapConfig{Origin: "ldap"},
			}
			roleMgrFake.ListOrgUsersByRoleReturns(role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), nil)
		})
		Context("SyncLdapUsers", func() {
			var roleUsers *role.RoleUsers

			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{
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
					SpaceName:      "spaceName",
					OrgName:        "orgName",
					AddUser:        roleMgrFake.AssociateSpaceAuditor,
					RoleUsers:      role.InitRoleUsers(),
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
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(userName).Should(Equal("test_ldap2"))

			})

			It("Should add ldap user to role", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{
					{Username: "test_ldap", GUID: "test_ldap-id"},
				}, uaaUsers)

				userManager.UAAUsers = uaaUsers
				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap2"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					SpaceName:      "spaceName",
					OrgName:        "orgName",
					AddUser:        roleMgrFake.AssociateSpaceAuditor,
					RoleUsers:      role.InitRoleUsers(),
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
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(userName).Should(Equal("test_ldap2"))
			})

			It("Should add ldap group member to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					SpaceName:      "spaceName",
					OrgName:        "orgName",
					AddUser:        roleMgrFake.AssociateSpaceAuditor,
					RoleUsers:      role.InitRoleUsers(),
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
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(userName).Should(Equal("test_ldap2"))
			})

			It("Should not add existing ldap user to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
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
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap_new"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
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
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
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
				Expect(err).Should(HaveOccurred())
				Expect(userManager.UAAUsers.GetByNameAndOrigin("test_ldap3", "ldap")).Should(BeNil())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)
				roleMgrFake.AssociateSpaceAuditorReturns(errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("User test_ldap3 with origin ldap: error"))
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
			})

			It("Should not query ldap if user exists in UAA", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}

				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {

				updateUsersInput := UsersInput{
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        roleMgrFake.AssociateSpaceAuditor,
					RoleUsers:      role.InitRoleUsers(),
				}
				ldapFake.GetUserDNsReturns([]string{"cn=test_ldap2"}, nil)
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserDNsCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
				Expect(ldapFake.GetUserByDNCallCount()).Should(Equal(0))
			})
			It("Should return error", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				ldapFake.GetUserByIDReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(0))
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
