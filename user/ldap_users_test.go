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
		Context("SyncLdapUsers", func() {
			var roleUsers *RoleUsers
			var uaaUsers map[string]uaa.User
			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
				uaaUsers = make(map[string]uaa.User)
				uaaUsers["test_ldap"] = uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap"}
				uaaUsers["test_ldap-id"] = uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap"}
				uaaUsers["cn=test_ldap"] = uaa.User{Username: "test_ldap", Origin: "ldap", ExternalID: "cn=test_ldap"}
				uaaUsers["test_ldap2"] = uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2"}
				uaaUsers["test_ldap2-id"] = uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2"}
				uaaUsers["cn=test_ldap2"] = uaa.User{Username: "test_ldap2", Origin: "ldap", ExternalID: "cn=test_ldap2"}
				roleUsers, _ = NewRoleUsers([]cfclient.User{
					cfclient.User{Username: "test_ldap", Guid: "test_ldap-id"},
				}, uaaUsers)
			})
			It("Should add ldap user to role", func() {
				updateUsersInput := UpdateUsersInput{
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
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test_ldap2"))
				Expect(origin).Should(Equal("ldap"))

				spaceGUID, userName, origin := client.AssociateSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test_ldap2"))
				Expect(origin).Should(Equal("ldap"))
			})

			It("Should add ldap group member to role", func() {
				updateUsersInput := UpdateUsersInput{
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
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))
				orgGUID, userName, origin := client.AssociateOrgUserByUsernameAndOriginArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(userName).Should(Equal("test_ldap2"))
				Expect(origin).Should(Equal("ldap"))

				spaceGUID, userName, origin := client.AssociateSpaceAuditorByUsernameAndOriginArgsForCall(0)
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test_ldap2"))
				Expect(origin).Should(Equal("ldap"))
			})

			It("Should not add existing ldap user to role", func() {
				updateUsersInput := UpdateUsersInput{
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
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})
			It("Should create external user when user doesn't exist in uaa", func() {
				updateUsersInput := UpdateUsersInput{
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

			It("Should not error when create external user errors", func() {
				updateUsersInput := UpdateUsersInput{
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
				uaaFake.CreateExternalUserReturns(errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(uaaUsers).ShouldNot(HaveKey("test_ldap3"))
				Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			})

			It("Should return error", func() {
				updateUsersInput := UpdateUsersInput{
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
				client.AssociateOrgUserByUsernameAndOriginReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("User test_ldap3: error"))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				updateUsersInput := UpdateUsersInput{
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserDNsReturns([]string{"cn=test_ldap2"}, nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserDNsCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
				Expect(ldapFake.GetUserByDNCallCount()).Should(Equal(0))
			})
			It("Should return error", func() {
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap3"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameAndOriginCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameAndOriginCallCount()).Should(Equal(0))
			})
		})
	})
})
