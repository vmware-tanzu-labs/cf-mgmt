package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
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
				LdapConfig: &config.LdapConfig{Origin: "ldap", LdapOrigin: "ldap"}}
		})
		Context("SyncLdapUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers *uaa.Users
			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:     "ldap",
					LdapOrigin: "ldap",
					Enabled:    true,
				}
				uaaUsers = &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test_ldap", Guid: "test_ldap-id"},
				}, uaaUsers)
			})
			It("Should add ldap user to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap2"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
			})

			It("Should add ldap user to role", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:     "ldap",
					LdapOrigin: "ldap",
					Enabled:    true,
				}
				uaaUsers = &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap", GUID: "test_ldap-id"})
				uaaUsers.Add(uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2", GUID: "test_ldap2-id"})
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test_ldap", Guid: "test_ldap-id"},
				}, uaaUsers)

				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap2"},
					LdapGroupNames: []string{},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))
			})

			It("Should add ldap group member to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap2",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap2-id"))

			})

			It("Should add saml ldap group member to role", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					UseIDForSAMLUser: true,
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap3"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("test_ldap3"))
				Expect(origin).Should(Equal("saml"))

			})

			It("Should not add existing ldap user to role", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap_new"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap_new",
						Email:  "test@test.com",
					},
					nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap_new"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("ldap_test_dn"))
				Expect(origin).Should(Equal("ldap"))
			})

			It("Should create external user when user doesn't exist in uaa 2", func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:     "saml",
					LdapOrigin: "ldap",
					Enabled:    true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{"test_ldap_new"},
					SpaceGUID:      "space_guid",
					LdapGroupNames: []string{"test_group"},
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap_new",
						Email:  "test@test.com",
					},
					nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap_new"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("ldap_test_dn"))
				Expect(origin).Should(Equal("ldap"))

			})

			// saml_group tests
			It("Should add saml ldap group member to role", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					LdapOrigin:       "ldap",
					UseIDForSAMLUser: true,
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					SamlGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap3"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("test_ldap3"))
				Expect(origin).Should(Equal("saml"))

			})

			// saml group filter test
			It("Should not add saml group member to role, because of filter", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					LdapOrigin:       "ldap",
					UseIDForSAMLUser: true,
					SamlUserFilter:   "bla",
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					SamlGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

			// filter test
			It("Should add saml group member to role, because of filter", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					LdapOrigin:       "ldap",
					UseIDForSAMLUser: true,
					SamlUserFilter:   "ldap_test_dn",
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					SamlGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap3"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("test_ldap3"))
				Expect(origin).Should(Equal("saml"))
			})

			// saml group filter and mode test
			It("Should not add saml group member to role, because of filter and exclusion mode", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:             "saml",
					LdapOrigin:         "ldap",
					UseIDForSAMLUser:   true,
					SamlUserFilter:     "ldap_test_dn",
					SamlUserFilterMode: "exclude",
					Enabled:            true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					SamlGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

			// ldap group filter test

			It("Should not add saml ldap group member to role, because of filter", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					UseIDForSAMLUser: true,
					LdapUserFilter:   "bla",
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

			It("Should not add saml ldap group member to role, because of filter and exclusion mode", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:             "saml",
					UseIDForSAMLUser:   true,
					LdapUserFilter:     "ldap_test_dn",
					LdapUserFilterMode: "exclude",
					Enabled:            true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

			// filter test
			It("Should  add saml ldap group member to role, because of filter", func() {

				userManager.LdapConfig = &config.LdapConfig{
					Origin:           "saml",
					UseIDForSAMLUser: true,
					LdapUserFilter:   "ldap_test_dn",
					Enabled:          true,
				}
				updateUsersInput := UsersInput{
					LdapUsers:      []string{},
					LdapGroupNames: []string{"test_group2"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}

				uaaFake.CreateExternalUserReturns("test_ldap3-id", nil)

				ldapFake.GetUserDNsReturns([]string{"cn=ldap_test_dn"}, nil)
				ldapFake.GetUserByDNReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				orgGUID, userGUID := client.AssociateOrgUserArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				spaceGUID, userGUID := client.AssociateSpaceAuditorArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userGUID).Should(Equal("test_ldap3-id"))

				arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
				Expect(arg1).Should(Equal("test_ldap3"))
				Expect(arg2).Should(Equal("test@test.com"))
				Expect(arg3).Should(Equal("test_ldap3"))
				Expect(origin).Should(Equal("saml"))
			})

			//

			It("Should not error when create external user errors", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)
				uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(uaaUsers.GetByName("test_ldap3"))).Should(Equal(0))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap3",
						Email:  "test@test.com",
					},
					nil)
				client.AssociateOrgUserReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("User test_ldap3 with origin ldap: error"))
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				updateUsersInput := UsersInput{
					LdapUsers: []string{"test_ldap2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {

				updateUsersInput := UsersInput{
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserDNsReturns([]string{"cn=test_ldap2"}, nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(1))
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
				}
				ldapFake.GetUserByIDReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
		})
		Context("UpdateUserInfo", func() {

			It("ldap origin with email", func() {
				userManager.LdapConfig.Origin = "ldap"
				userInfo := userManager.UpdateUserInfo(ldap.User{
					Email:  "test@test.com",
					UserID: "testUser",
					Origin: "ldap",
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
					Origin: "ldap",
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
					Origin: "foo",
				})
				Expect(userInfo.Email).Should(Equal("test@test.com"))
				Expect(userInfo.UserDN).Should(Equal("test@test.com"))
				Expect(userInfo.UserID).Should(Equal("test@test.com"))
			})
		})
	})
})
