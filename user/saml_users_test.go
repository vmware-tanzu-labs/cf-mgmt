package user_test

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	rolefakes "github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
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
		userManager = &DefaultManager{
			Cfg:        fakeReader,
			UAAMgr:     uaaFake,
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgReader:  orgFake,
			Peek:       false,
			RoleMgr:    roleMgrFake,
			LdapConfig: &config.LdapConfig{Origin: "saml_origin"},
		}
		roleMgrFake.ListOrgUsersByRoleReturns(role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), nil)
	})
	Context("SyncSamlUsers", func() {
		var roleUsers *role.RoleUsers
		BeforeEach(func() {
			userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}
			uaaUsers := &uaa.Users{}

			uaaUsers.Add(uaa.User{Username: "Test.Test@test.com", Email: "test.test@test.com", ExternalID: "Test.Test@test.com", Origin: "saml_origin", GUID: "test-id"})
			uaaUsers.Add(uaa.User{Username: "test2.test2@test.com", Email: "test2.test2@test.com", ExternalID: "test2.test2@test.com", Origin: "saml_origin", GUID: "test2-id"})
			roleUsers, _ = role.NewRoleUsers(
				[]*resource.User{
					{Username: "Test.Test@test.com", GUID: "test-id"},
				},
				uaaUsers,
			)
			userManager.UAAUsers = uaaUsers
		})
		It("Should add saml user to role", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test2.test2@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				OrgName:   "test-org",
				SpaceName: "test-space",
				RoleUsers: role.InitRoleUsers(),
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
			orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
			Expect(orgGUID).Should(Equal("org_guid"))
			Expect(userGUID).Should(Equal("test2-id"))
			Expect(spaceGUID).Should(Equal("space_guid"))
			Expect(spaceName).Should(Equal("test-org/test-space"))
			Expect(userName).Should(Equal("test2.test2@test.com"))
		})

		It("Should not add existing saml user to role", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test.test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			Expect(roleUsers.HasUser("test.test@test.com")).Should(BeTrue())
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleUsers.HasUser("test.test@test.com")).Should(BeFalse())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(0))

		})

		It("Should not add existing saml user to role due to mixed case match", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"Test.Test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleUsers.HasUser("Test.Test@test.com")).Should(BeFalse())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(0))
			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(0))
		})
		It("Should create external user when user doesn't exist in uaa", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test3.test3@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
			arg1, arg2, arg3, origin := uaaFake.CreateExternalUserArgsForCall(0)
			Expect(arg1).Should(Equal("test3.test3@test.com"))
			Expect(arg2).Should(Equal("test3.test3@test.com"))
			Expect(arg3).Should(Equal("test3.test3@test.com"))
			Expect(origin).Should(Equal("saml_origin"))
		})

		It("Should not error when create external user errors", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test.test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			userManager.UAAUsers = &uaa.Users{}
			uaaFake.CreateExternalUserReturns("guid", errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(uaaFake.CreateExternalUserCallCount()).Should(Equal(1))
		})

		It("Should return error", func() {
			roleUsers := role.InitRoleUsers()
			roleUsers.AddUsers([]role.RoleUser{
				{UserName: "test"},
			})
			uaaUsers := &uaa.Users{}
			uaaUsers.Add(uaa.User{Username: "test.test@test.com"})
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test.test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			userManager.UAAUsers = uaaUsers
			roleMgrFake.AssociateSpaceAuditorReturns(errors.New("Got an error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).Should(HaveOccurred())
			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
		})
	})
})
