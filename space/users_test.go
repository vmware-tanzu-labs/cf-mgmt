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
			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user-id"]
			Ω(ok).Should(BeTrue())
		})
		It("update ldap group users where users are in uaac", func() {
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
			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user-id"]
			Ω(ok).Should(BeTrue())
		})

		It("update ldap group users where users are in uaac and already in space", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-id"] = "user-id"
			spaceUsers := make(map[string]string)
			spaceUsers["user-id"] = "user-id"
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

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user-id"]
			Ω(ok).Should(BeTrue())
		})

		It("update other origin users where users are in uaac and already in space", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "other",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user@test.com"] = "user@test.com"
			spaceUsers := make(map[string]string)
			spaceUsers["user@test.com"] = "user@test.com"
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

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user@test.com"]
			Ω(ok).Should(BeTrue())
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
			Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
			_, ok := uaacUsers["user-1-id"]
			Ω(ok).Should(BeTrue())
			_, ok = uaacUsers["user-2-id"]
			Ω(ok).Should(BeTrue())
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

			Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
			_, ok := uaacUsers["user-1-id"]
			Ω(ok).Should(BeTrue())
			_, ok = uaacUsers["user-2-id"]
			Ω(ok).Should(BeTrue())
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

			Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
			_, ok := uaacUsers["user-1"]
			Ω(ok).Should(BeTrue())
			_, ok = uaacUsers["user-2"]
			Ω(ok).Should(BeTrue())
		})

		It("update users where users are in uaac and in a space", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-1"] = "user-1"
			uaacUsers["user-2"] = "user-2"
			spaceUsers := make(map[string]string)
			spaceUsers["user-1"] = "asfdsdf-1"
			spaceUsers["user-2"] = "asdfsaf-2"
			updateUsersInput := UpdateUsersInput{
				SpaceGUID: "my-space-guid",
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				Users:     []string{"USER-1", "user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)

			err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())

			Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
			_, ok := uaacUsers["user-1"]
			Ω(ok).Should(BeTrue())
			_, ok = uaacUsers["user-2"]
			Ω(ok).Should(BeTrue())
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

			Ω(userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)).ShouldNot(Succeed())
			Ω(len(uaacUsers)).Should(BeEquivalentTo(0))
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

		It("remove orphaned LDAP users while leaving existing group members - GH issue 33", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "https://saml.example.com",
			}

			uaacUsers := make(map[string]string)
			uaacUsers["chris.a.washburn@example.com"] = "cwashburn-uaac-guid"
			uaacUsers["joe.h.fitzy@example.com"] = "jfitzy-uaac-guid"
			uaacUsers["alex.j.smith@example.com"] = "asmith-uaac-guid" // <-- user in uaac, but not ldap group

			spaceUsers := make(map[string]string)
			spaceUsers["chris.a.washburn@example.com"] = "cwashburn-space-user-guid"
			spaceUsers["joe.h.fitzy@example.com"] = "jfitzy-space-user-guid"
			spaceUsers["alex.j.smith@example.com"] = "asmith-space-user-guid" // <-- user in space, but not ldap group

			updateUsersInput := UpdateUsersInput{
				SpaceName:     "space-name",
				SpaceGUID:     "space-guid",
				OrgName:       "org-name",
				OrgGUID:       "org-guid",
				Role:          "space-role-name",
				LdapGroupName: "ldap-group-name",
				RemoveUsers:   true,
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "CN=Washburn, Chris,OU=End Users,OU=Accounts,DC=add,DC=example,DC=com",
				UserID: "u-cwashburn",
				Email:  "Chris.A.Washburn@example.com",
			}, l.User{
				UserDN: "CN=Fitzy, Joe,OU=End Users,OU=Accounts,DC=ad,DC=example,DC=com",
				UserID: "u-jfitzy",
				Email:  "Joe.H.Fitzy@example.com",
			}}

			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			mockCloudController.EXPECT().GetCFUsers("space-guid", "spaces", "space-role-name").Return(spaceUsers, nil)
			mockCloudController.EXPECT().RemoveCFUser("space-guid", "spaces", "asmith-space-user-guid", "space-role-name").Return(nil)
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
