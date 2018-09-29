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
				Client:     client,
				Cfg:        fakeReader,
				UAAMgr:     uaaFake,
				LdapMgr:    ldapFake,
				SpaceMgr:   spaceFake,
				OrgMgr:     orgFake,
				Peek:       false,
				LdapConfig: &config.LdapConfig{Origin: "ldap"}}
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
		Context("SyncLdapUsers", func() {
			BeforeEach(func() {
				userManager.LdapConfig = &config.LdapConfig{
					Origin:  "ldap",
					Enabled: true,
				}
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

				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

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

			It("Should add ldap group member to role", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
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
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)

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
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)
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
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)
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
				ldapFake.GetUserByIDReturns(
					&ldap.User{
						UserDN: "ldap_test_dn",
						UserID: "test_ldap",
						Email:  "test@test.com",
					},
					nil)
				client.AssociateOrgUserByUsernameReturns(cfclient.Org{}, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}

				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
			})

			It("Should not query ldap if user exists in UAA", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				uaaUsers["cn=test_ldap"] = &uaa.User{UserName: "test_ldap"}
				updateUsersInput := UpdateUsersInput{
					LdapGroupNames: []string{"test_group"},
					SpaceGUID:      "space_guid",
					OrgGUID:        "org_guid",
					AddUser:        userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserDNsReturns([]string{"cn=test_ldap"}, nil)
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(1))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserDNsCallCount()).Should(Equal(1))
				Expect(ldapFake.GetUserByIDCallCount()).Should(Equal(0))
				Expect(ldapFake.GetUserByDNCallCount()).Should(Equal(0))
			})
			It("Should return error", func() {
				roleUsers := make(map[string]string)
				uaaUsers := make(map[string]*uaa.User)
				updateUsersInput := UpdateUsersInput{
					LdapUsers: []string{"test_ldap"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   userManager.AssociateSpaceAuditor,
				}
				ldapFake.GetUserByIDReturns(nil, errors.New("error"))
				err := userManager.SyncLdapUsers(roleUsers, uaaUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("error"))
				Expect(client.AssociateOrgUserByUsernameCallCount()).Should(Equal(0))
				Expect(client.AssociateSpaceAuditorByUsernameCallCount()).Should(Equal(0))
			})
		})
	})
})
