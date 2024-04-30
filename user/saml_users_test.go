package user_test

import (
	"errors"

	uaaclient "github.com/cloudfoundry-community/go-uaa"
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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SamlUsers", func() {
	var (
		userManager *DefaultManager
		ldapFake    *fakes.FakeLdapManager
		uaaFake     *uaafakes.FakeUaa
		fakeReader  *configfakes.FakeReader
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeReader
		roleMgrFake *rolefakes.FakeManager
	)
	BeforeEach(func() {
		ldapFake = new(fakes.FakeLdapManager)
		uaaFake = new(uaafakes.FakeUaa)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
		roleMgrFake = new(rolefakes.FakeManager)
		userManager = &DefaultManager{
			Cfg:        fakeReader,
			UAAMgr:     &uaa.DefaultUAAManager{Client: uaaFake},
			LdapMgr:    ldapFake,
			SpaceMgr:   spaceFake,
			OrgReader:  orgFake,
			Peek:       false,
			RoleMgr:    roleMgrFake,
			LdapConfig: &config.LdapConfig{Origin: "saml_origin"},
		}
		roleMgrFake.ListOrgUsersByRoleReturns(role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), nil)
		fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{}, nil)
	})
	Context("SyncSamlUsers", func() {
		var roleUsers *role.RoleUsers
		BeforeEach(func() {
			userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}

			uaaUsers := []uaaclient.User{}
			uaaUsers = append(uaaUsers, uaaclient.User{Username: "Test.Test@test.com", Emails: []uaaclient.Email{{Value: "test.test@test.com"}}, ExternalID: "Test.Test@test.com", Origin: "saml_origin", ID: "test-id"})
			uaaUsers = append(uaaUsers, uaaclient.User{Username: "test2.test2@test.com", Emails: []uaaclient.Email{{Value: "test2.test2@test.com"}}, ExternalID: "test2.test2@test.com", Origin: "saml_origin", ID: "test2-id"})
			uaaFake.ListUsersReturns(uaaUsers, uaaclient.Page{StartIndex: 1, TotalResults: 2, ItemsPerPage: 500}, nil)

			users, err := userManager.UAAMgr.ListUsers()
			Expect(err).ShouldNot(HaveOccurred())
			roleUsers, _ = role.NewRoleUsers(
				[]*uaa.User{
					{Username: "Test.Test@test.com", GUID: "test-id"},
				},
				users,
			)

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
			Expect(roleUsers.HasUserForOrigin("test.test@test.com", "saml_origin")).Should(BeTrue())
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(roleUsers.HasUserForOrigin("test.test@test.com", "saml_origin")).Should(BeFalse())
			Expect(uaaFake.CreateUserCallCount()).Should(Equal(0))
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
			Expect(roleUsers.HasUserForOrigin("Test.Test@test.com", "saml_origin")).Should(BeFalse())
			Expect(uaaFake.CreateUserCallCount()).Should(Equal(0))
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
			uaaFake.CreateUserReturns(&uaaclient.User{ID: "user-guid"}, nil)
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(uaaFake.CreateUserCallCount()).Should(Equal(1))
			user := uaaFake.CreateUserArgsForCall(0)
			Expect(user.Username).Should(Equal("test3.test3@test.com"))
			Expect(user.Emails[0].Value).Should(Equal("test3.test3@test.com"))
			Expect(user.ExternalID).Should(Equal("test3.test3@test.com"))
			Expect(user.Origin).Should(Equal("saml_origin"))
		})

		It("Should not error when create external user errors", func() {
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test.test@foo.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			uaaFake.ListUsersReturns([]uaaclient.User{}, uaaclient.Page{StartIndex: 1, TotalResults: 0, ItemsPerPage: 500}, nil)
			uaaFake.CreateUserReturns(nil, errors.New("error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).Should(HaveOccurred())
			Expect(uaaFake.CreateUserCallCount()).Should(Equal(1))
		})

		It("Should return error", func() {
			roleUsers := role.InitRoleUsers()
			roleUsers.AddUsers([]role.RoleUser{
				{UserName: "test"},
			})

			uaaUsers := []uaaclient.User{}
			uaaUsers = append(uaaUsers, uaaclient.User{Username: "test.test@test.com", Origin: "saml_origin"})
			updateUsersInput := UsersInput{
				SamlUsers: []string{"test.test@test.com"},
				SpaceGUID: "space_guid",
				OrgGUID:   "org_guid",
				AddUser:   roleMgrFake.AssociateSpaceAuditor,
				RoleUsers: role.InitRoleUsers(),
			}
			uaaFake.ListUsersReturns(uaaUsers, uaaclient.Page{StartIndex: 1, TotalResults: 0, ItemsPerPage: 500}, nil)
			roleMgrFake.AssociateSpaceAuditorReturns(errors.New("Got an error"))
			err := userManager.SyncSamlUsers(roleUsers, updateUsersInput)
			Expect(err).Should(HaveOccurred())
			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
		})
	})
	Context("Change Saml Origin", func() {
		var roleUsers *role.RoleUsers
		BeforeEach(func() {
			userManager.LdapConfig = &config.LdapConfig{Origin: "saml_origin"}
			uaaUsers := []uaaclient.User{}
			uaaUsers = append(uaaUsers, uaaclient.User{Username: "test.test@test.com", Emails: []uaaclient.Email{{Value: "test.test@test.com"}}, ExternalID: "test.test@test.com", Origin: "saml_original_origin", ID: "test-id"})
			uaaFake.ListUsersReturns(uaaUsers, uaaclient.Page{StartIndex: 1, TotalResults: 1, ItemsPerPage: 500}, nil)

			users, err := userManager.UAAMgr.ListUsers()
			Expect(err).ShouldNot(HaveOccurred())
			roleUsers, _ = role.NewRoleUsers(
				[]*uaa.User{
					{Username: "test.test@test.com", GUID: "test-id", Origin: "saml_original_origin"},
				},
				users,
			)
		})

		It("Should add saml new user to role with new origin and remove old user from role", func() {
			updateUsersInput := UsersInput{
				SamlUsers:   []string{"test.test@test.com"},
				SpaceGUID:   "space_guid",
				OrgGUID:     "org_guid",
				OrgName:     "test-org",
				SpaceName:   "test-space",
				RoleUsers:   roleUsers,
				AddUser:     roleMgrFake.AssociateSpaceAuditor,
				RemoveUser:  roleMgrFake.RemoveSpaceAuditor,
				RemoveUsers: true,
			}
			uaaFake.CreateUserReturns(&uaaclient.User{ID: "new-user-guid"}, nil)
			err := userManager.SyncUsers(updateUsersInput)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(uaaFake.CreateUserCallCount()).Should(Equal(1))
			user := uaaFake.CreateUserArgsForCall(0)
			Expect(user.Username).Should(Equal("test.test@test.com"))
			Expect(user.Emails[0].Value).Should(Equal("test.test@test.com"))
			Expect(user.ExternalID).Should(Equal("test.test@test.com"))
			Expect(user.Origin).Should(Equal("saml_origin"))

			Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
			orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
			Expect(orgGUID).Should(Equal("org_guid"))
			Expect(userGUID).Should(Equal("new-user-guid"))
			Expect(spaceGUID).Should(Equal("space_guid"))
			Expect(spaceName).Should(Equal("test-org/test-space"))
			Expect(userName).Should(Equal("test.test@test.com"))

			Expect(roleMgrFake.RemoveSpaceAuditorCallCount()).Should(Equal(1))
			spaceName, spaceGUID, userName, userGUID = roleMgrFake.RemoveSpaceAuditorArgsForCall(0)
			Expect(userGUID).Should(Equal("test-id"))
			Expect(spaceGUID).Should(Equal("space_guid"))
			Expect(spaceName).Should(Equal("test-org/test-space"))
			Expect(userName).Should(Equal("test.test@test.com"))
		})

	})
})
