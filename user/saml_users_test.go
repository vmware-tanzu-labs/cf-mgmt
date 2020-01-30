package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	configfakes "github.com/pivotalservices/cf-mgmt/config/fakes"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	spacefakes "github.com/pivotalservices/cf-mgmt/space/fakes"
	"github.com/pivotalservices/cf-mgmt/uaa"
	uaafakes "github.com/pivotalservices/cf-mgmt/uaa/fakes"
	. "github.com/pivotalservices/cf-mgmt/user"
	"github.com/pivotalservices/cf-mgmt/user/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SamlUsers", func() {
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
		userManager = &DefaultManager{
			Client:     client,
			Cfg:        fakeReader,
			UAAMgr:     uaaFake,
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgReader:  orgFake,
			Peek:       false,
			LdapConfig: &config.LdapConfig{Origin: "saml_origin"}}
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
				OrgName:   "test-org",
				SpaceName: "test-space",
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
			roleUsers := InitRoleUsers()
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
})
