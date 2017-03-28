package space_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	l "github.com/pivotalservices/cf-mgmt/ldap"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"

	. "github.com/pivotalservices/cf-mgmt/space"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
)

var _ = Describe("given SpaceManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", "uaacToken")
			Ω(manager).ShouldNot(BeNil())
		})
	})

	var (
		ctrl                *gomock.Controller
		mockCloudController *cc.MockManager
		mockLdap            *ldap.MockManager
		mockUaac            *uaac.MockManager
		userManager         UserMgr
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		userManager = NewUserManager(mockCloudController, mockLdap, mockUaac)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("UpdateSpaceUsers()", func() {
		It("update ldap group users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID:     "my-space-guid",
				OrgGUID:       "my-org-guid",
				Role:          "my-role",
				LdapGroupName: "ldap-group-name",
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "user-dn",
				UserID: "user-id",
				Email:  "user@test.com",
			}}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			mockUaac.EXPECT().CreateExternalUser("user-id", "user@test.com", "user-dn", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-id", "my-role", "my-space-guid").Return(nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})
		It("update ldap group users where users are not uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-id"] = "user-id"
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID:     "my-space-guid",
				OrgGUID:       "my-org-guid",
				Role:          "my-role",
				LdapGroupName: "ldap-group-name",
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "user-dn",
				UserID: "user-id",
				Email:  "user@test.com",
			}}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-id", "my-role", "my-space-guid").Return(nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})

		It("update ldap users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockLdap.EXPECT().GetUser(config, "ldap-user-1").Return(&l.User{
				UserDN: "user-1-dn",
				UserID: "user-1-id",
				Email:  "user1@test.com",
			}, nil)
			mockLdap.EXPECT().GetUser(config, "ldap-user-2").Return(&l.User{
				UserDN: "user-2-dn",
				UserID: "user-2-id",
				Email:  "user2@test.com",
			}, nil)

			mockUaac.EXPECT().CreateExternalUser("user-1-id", "user1@test.com", "user-1-dn", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("user-1-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-1-id", "my-role", "my-space-guid").Return(nil)

			mockUaac.EXPECT().CreateExternalUser("user-2-id", "user2@test.com", "user-2-dn", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-2-id", "my-role", "my-space-guid").Return(nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})

		It("update ldap users where users are in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-1-id"] = "user-1-id"
			uaacUsers["user-2-id"] = "user-2-id"
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockLdap.EXPECT().GetUser(config, "ldap-user-1").Return(&l.User{
				UserDN: "user-1-dn",
				UserID: "user-1-id",
				Email:  "user1@test.com",
			}, nil)
			mockLdap.EXPECT().GetUser(config, "ldap-user-2").Return(&l.User{
				UserDN: "user-2-dn",
				UserID: "user-2-id",
				Email:  "user2@test.com",
			}, nil)

			mockCloudController.EXPECT().AddUserToOrg("user-1-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-1-id", "my-role", "my-space-guid").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-2-id", "my-role", "my-space-guid").Return(nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})

		It("update users where users are in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-1"] = "user-1"
			uaacUsers["user-2"] = "user-2"
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				Users:     []string{"user-1", "user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockCloudController.EXPECT().AddUserToOrg("user-1", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-1", "my-role", "my-space-guid").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("user-2", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("user-2", "my-role", "my-space-guid").Return(nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})

		It("update users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			spaceUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				Users:     []string{"user-1"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(Not(BeNil()))
			Ω(err.Error()).Should(BeEquivalentTo("User user-1 doesn't exist in cloud foundry, so must add internal user first"))
		})

		It("remove users that in space but not in config", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			spaceUsers := make(map[string]string)
			spaceUsers["cwashburn"] = "cwashburn"
			spaceUsers["cwashburn1"] = "cwashburn1"
			spaceUsers["cwashburn2"] = "cwashburn2"
			updateUsersInput := UpdateUsersInput{
				SpaceGUID:   "my-space-guid",
				OrgGUID:     "my-org-guid",
				Role:        "my-role",
				RemoveUsers: true,
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			mockCloudController.EXPECT().RemoveCFUser("my-space-guid", "spaces", "cwashburn", "my-role").Return(nil)
			mockCloudController.EXPECT().RemoveCFUser("my-space-guid", "spaces", "cwashburn1", "my-role").Return(nil)
			mockCloudController.EXPECT().RemoveCFUser("my-space-guid", "spaces", "cwashburn2", "my-role").Return(nil)
			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})
		It("don't remove users that in space but not in config", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			spaceUsers := make(map[string]string)
			spaceUsers["cwashburn"] = "cwashburn"
			spaceUsers["cwashburn1"] = "cwashburn1"
			spaceUsers["cwashburn2"] = "cwashburn2"
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})
	})
})
