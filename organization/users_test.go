package organization_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
	l "github.com/pivotalservices/cf-mgmt/ldap"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"

	. "github.com/pivotalservices/cf-mgmt/organization"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	"github.com/pivotalservices/cf-mgmt/utils/mocks"
)

var _ = Describe("given UserManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", "uaacToken", config.NewManager("./fixtures/config", mock_utils.NewMockUtilsManager()))
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

	Context("UpdateorgUsers()", func() {
		It("update ldap group users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID:        "my-org-guid",
				Role:           "my-role",
				LdapGroupNames: []string{"ldap-group-name"},
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "user-dn",
				UserID: "user-id",
				Email:  "user@test.com",
			}}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			mockUaac.EXPECT().CreateExternalUser("user-id", "user@test.com", "user-dn", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-id", "my-role", "my-org-guid").Return(nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
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
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID:        "my-org-guid",
				Role:           "my-role",
				LdapGroupNames: []string{"ldap-group-name"},
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "user-dn",
				UserID: "user-id",
				Email:  "user@test.com",
			}}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-id", "my-role", "my-org-guid").Return(nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())

			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user-id"]
			Ω(ok).Should(BeTrue())
		})

		It("update ldap group users where users are in uaac and in org", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-id"] = "user-id"
			orgUsers := make(map[string]string)
			orgUsers["user-id"] = "user-id"
			updateUsersInput := UpdateUsersInput{
				OrgGUID:        "my-org-guid",
				Role:           "my-role",
				LdapGroupNames: []string{"ldap-group-name"},
			}

			ldapGroupUsers := []l.User{l.User{
				UserDN: "user-dn",
				UserID: "user-id",
				Email:  "user@test.com",
			}}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())

			Ω(len(uaacUsers)).Should(BeEquivalentTo(1))
			_, ok := uaacUsers["user-id"]
			Ω(ok).Should(BeTrue())
		})

		It("update ldap users where users are not in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
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
			mockCloudController.EXPECT().AddUserToOrgRole("user-1-id", "my-role", "my-org-guid").Return(nil)

			mockUaac.EXPECT().CreateExternalUser("user-2-id", "user2@test.com", "user-2-dn", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-2-id", "my-role", "my-org-guid").Return(nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
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
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID:   "my-org-guid",
				Role:      "my-role",
				LdapUsers: []string{"ldap-user-1", "ldap-user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
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
			mockCloudController.EXPECT().AddUserToOrgRole("user-1-id", "my-role", "my-org-guid").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("user-2-id", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-2-id", "my-role", "my-org-guid").Return(nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
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
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID: "my-org-guid",
				Role:    "my-role",
				Users:   []string{"user-1", "user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			mockCloudController.EXPECT().AddUserToOrg("user-1", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-1", "my-role", "my-org-guid").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("user-2", "my-org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("user-2", "my-role", "my-org-guid").Return(nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())

			Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
			_, ok := uaacUsers["user-1"]
			Ω(ok).Should(BeTrue())
			_, ok = uaacUsers["user-2"]
			Ω(ok).Should(BeTrue())
		})

		It("update users where users are in uaac and in org", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["user-1"] = "user-1"
			uaacUsers["user-2"] = "user-2"
			orgUsers := make(map[string]string)
			orgUsers["user-1"] = "user-1"
			orgUsers["user-2"] = "user-2"
			updateUsersInput := UpdateUsersInput{
				OrgGUID: "my-org-guid",
				Role:    "my-role",
				Users:   []string{"USER-1", "user-2"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
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
			orgUsers := make(map[string]string)
			updateUsersInput := UpdateUsersInput{
				OrgGUID: "my-org-guid",
				Role:    "my-role",
				Users:   []string{"user-1"},
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)

			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(Not(BeNil()))
			Ω(err.Error()).Should(BeEquivalentTo("User user-1 doesn't exist in cloud foundry, so must add internal user first"))

			Ω(len(uaacUsers)).Should(BeEquivalentTo(0))

		})

		It("remove users that in space but not in config", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			orgUsers := make(map[string]string)
			orgUsers["cwashburn"] = "cwashburn"
			orgUsers["cwashburn1"] = "cwashburn1"
			orgUsers["cwashburn2"] = "cwashburn2"
			updateUsersInput := UpdateUsersInput{
				OrgGUID:     "my-org-guid",
				Role:        "my-role",
				RemoveUsers: true,
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			mockCloudController.EXPECT().RemoveCFUser("my-org-guid", "organizations", "cwashburn", "my-role").Return(nil)
			mockCloudController.EXPECT().RemoveCFUser("my-org-guid", "organizations", "cwashburn1", "my-role").Return(nil)
			mockCloudController.EXPECT().RemoveCFUser("my-org-guid", "organizations", "cwashburn2", "my-role").Return(nil)
			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})
		It("don't remove users that in space but not in config", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			orgUsers := make(map[string]string)
			orgUsers["cwashburn"] = "cwashburn"
			orgUsers["cwashburn1"] = "cwashburn1"
			orgUsers["cwashburn2"] = "cwashburn2"
			updateUsersInput := UpdateUsersInput{
				OrgGUID: "my-org-guid",
				Role:    "my-role",
			}

			mockCloudController.EXPECT().GetCFUsers("my-org-guid", "organizations", "my-role").Return(orgUsers, nil)
			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
		})
		It("adding users to uaac based on saml", func() {
			config := &l.Config{
				Enabled: false,
				Origin:  "https://saml.example.com",
			}

			uaacUsers := make(map[string]string)
			uaacUsers["chris.a.washburn@example.com"] = "cwashburn-uaac-guid"
			uaacUsers["joe.h.fitzy@example.com"] = "jfitzy-uaac-guid"

			orgUsers := make(map[string]string)
			orgUsers["chris.a.washburn@example.com"] = "cwashburn-space-user-guid"
			orgUsers["joe.h.fitzy@example.com"] = "jfitzy-space-user-guid"

			updateUsersInput := UpdateUsersInput{
				OrgName:     "org-name",
				OrgGUID:     "org-guid",
				Role:        "org-role-name",
				SamlUsers:   []string{"chris.a.washburn@example.com", "joe.h.fitzy@example.com", "test@test.com"},
				RemoveUsers: true,
			}

			mockCloudController.EXPECT().GetCFUsers("org-guid", "organizations", "org-role-name").Return(orgUsers, nil)
			mockUaac.EXPECT().CreateExternalUser("test@test.com", "test@test.com", "test@test.com", "https://saml.example.com").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("test@test.com", "org-guid").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("test@test.com", "org-role-name", "org-guid").Return(nil)
			err := userManager.UpdateOrgUsers(config, uaacUsers, updateUsersInput)
			Ω(err).Should(BeNil())
			Ω(uaacUsers).Should(HaveKey("test@test.com"))
		})
	})
})
