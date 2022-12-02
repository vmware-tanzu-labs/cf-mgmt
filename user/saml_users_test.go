package user_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

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
			LdapConfig:  &config.LdapConfig{Origin: "saml_origin"}}
	})
	Context("SyncSamlUsers", func() {
		var roleUsers *RoleUsers
		BeforeEach(func() {
			userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}
			uaaUsers := &uaa.Users{}
			uaaUsers.Add(uaa.User{Username: "test@test.com", Origin: "saml_origin", GUID: "test-id"})
			uaaUsers.Add(uaa.User{Username: "test@test2.com", Origin: "saml_origin", GUID: "test2-id"})
			roleUsers, _ = NewRoleUsers(
				[]*resource.User{
					{Username: "test@test.com", GUID: "test-id"},
				},
				uaaUsers,
			)
			userManager.UAAUsers = uaaUsers
		})
		It("Should add saml user to role", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test2.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				OrgName:   "test-org",
				SpaceName: "test-space",
				RoleUsers: InitRoleUsers(),
				AddUser:   userManager.AssociateSpaceAuditor,
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
			_, orgGUID, userGUID, role := fakeRoleClient.CreateOrganizationRoleArgsForCall(0)
			Expect(orgGUID).Should(Equal("org_guid"))
			Expect(userGUID).Should(Equal("test2-id"))
			Expect(role).To(Equal(resource.OrganizationRoleUser))

			_, spaceGUID, userGUID, roleType := fakeRoleClient.CreateSpaceRoleArgsForCall(0)
			Expect(spaceGUID).Should(Equal("space_guid"))
			Expect(userGUID).Should(Equal("test2-id"))
			Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))
		})

		It("Should not add existing saml user to role", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   userManager.AssociateSpaceAuditor,
				RoleUsers: InitRoleUsers(),
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleUsers.HasUser("test@test.com")).Should(BeFalse())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
			Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(0))
			Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(0))
		})
		It("Should create external user when user doesn't exist in uaa", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test@test3.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   userManager.AssociateSpaceAuditor,
				RoleUsers: InitRoleUsers(),
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
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
				RoleUsers: InitRoleUsers(),
			}
			userManager.UAAUsers = &uaa.Users{}
			uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
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
				RoleUsers: InitRoleUsers(),
			}
			userManager.UAAUsers = uaaUsers
			fakeRoleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).Should(HaveOccurred())
			Expect(fakeRoleClient.CreateOrganizationRoleCallCount()).Should(Equal(1))
			Expect(fakeRoleClient.CreateSpaceRoleCallCount()).Should(Equal(0))
		})
	})
})
