package user_test

import (
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"

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
				[]cfclient.V3User{
					{Username: "test@test.com", GUID: "test-id"},
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
				OrgUsers:  InitRoleUsers(),
				AddUser:   userManager.AssociateSpaceAuditor,
			}
			err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(client.CreateV3OrganizationRoleCallCount()).Should(Equal(1))
			orgGUID, userGUID, role := client.CreateV3OrganizationRoleArgsForCall(0)
			Expect(orgGUID).Should(Equal("org_guid"))
			Expect(userGUID).Should(Equal("test2-id"))
			Expect(role).To(Equal(ORG_USER))

			spaceGUID, userGUID, roleType := client.CreateV3SpaceRoleArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space_guid"))
			Expect(userGUID).Should(Equal("test2-id"))
			Expect(roleType).Should(Equal(SPACE_AUDITOR))
		})

		It("Should not add existing saml user to role", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   userManager.AssociateSpaceAuditor,
				OrgUsers:  InitRoleUsers(),
			}
			err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleUsers.HasUser("test@test.com")).Should(BeFalse())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
			Expect(client.CreateV3OrganizationRoleCallCount()).Should(Equal(0))
			Expect(client.CreateV3SpaceRoleCallCount()).Should(Equal(0))
		})
		It("Should create external user when user doesn't exist in uaa", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test3.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   userManager.AssociateSpaceAuditor,
				OrgUsers:  InitRoleUsers(),
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
				OrgUsers:  InitRoleUsers(),
			}
			uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, &uaa.Users{}, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
		})

		It("Should return error", func() {
			roleUsers := InitRoleUsers()
			roleUsers.AddUsers([]RoleUser{
				{UserName: "test"},
			})
			uaaUsers := &uaa.Users{}
			uaaUsers.Add(uaa.User{Username: "test@test.com"})
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   userManager.AssociateSpaceAuditor,
				OrgUsers:  InitRoleUsers(),
			}
			client.CreateV3OrganizationRoleReturns(&cfclient.V3Role{}, errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, uaaUsers, updateUsersInput)
			Expect(err).Should(HaveOccurred())
			Expect(client.CreateV3OrganizationRoleCallCount()).Should(Equal(1))
			Expect(client.CreateV3SpaceRoleCallCount()).Should(Equal(0))
		})
	})
})
