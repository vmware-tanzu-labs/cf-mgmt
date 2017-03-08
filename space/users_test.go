package space_test

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	l "github.com/pivotalservices/cf-mgmt/ldap"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	o "github.com/pivotalservices/cf-mgmt/organization/mocks"
	. "github.com/pivotalservices/cf-mgmt/space"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	"github.com/pivotalservices/cf-mgmt/utils"
)

var _ = XDescribe("given SpaceManager", func() {
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
		mockOrgMgr          *o.MockManager
		spaceManager        DefaultSpaceManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		spaceManager = DefaultSpaceManager{
			CloudController: mockCloudController,
			UAACMgr:         mockUaac,
			UtilsMgr:        utils.NewDefaultManager(),
			LdapMgr:         mockLdap,
			OrgMgr:          mockOrgMgr,
			UserMgr:         NewUserManager(mockCloudController, mockLdap, mockUaac),
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("UpdateSpaceUsers()", func() {
		spaces := []cloudcontroller.Space{
			{
				Entity: cloudcontroller.SpaceEntity{
					Name:    "space1",
					OrgGUID: "testOrgGUID",
				},
				MetaData: cloudcontroller.SpaceMetaData{
					GUID: "space1GUID",
				},
			},
			{
				Entity: cloudcontroller.SpaceEntity{
					Name:    "space2",
					OrgGUID: "testOrgGUID",
				},
				MetaData: cloudcontroller.SpaceMetaData{
					GUID: "space2GUID",
				},
			},
		}
		It("update org users where users are already in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			uaacUsers["cwashburn"] = "cwashburn"
			uaacUsers["cwashburn1"] = "cwashburn1"
			uaacUsers["cwashburn2"] = "cwashburn2"
			users := []l.User{
				{UserID: "cwashburn", UserDN: "cn=cwashburn", Email: "cwashburn@testdomain.com"},
			}
			mockLdap.EXPECT().GetConfig("./fixtures/user_config", "test").Return(config, nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)

			mockLdap.EXPECT().GetUserIDs(config, "test_space1_developers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "developers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "managers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "auditors").Return(make(map[string]string), nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "developers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "test_space1_managers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "managers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "test_space1_auditors").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "auditors", "space1GUID").Return(nil)

			err := spaceManager.UpdateSpaceUsers("./fixtures/user_config", "test")
			Ω(err).Should(BeNil())
		})
		It("update org users where users aren't in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			uaacUsers := make(map[string]string)
			users := []l.User{
				{UserID: "cwashburn", UserDN: "cn=cwashburn", Email: "cwashburn@testdomain.com"},
			}
			mockLdap.EXPECT().GetConfig("./fixtures/user_config", "test").Return(config, nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_space1_developers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "developers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "managers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "auditors").Return(make(map[string]string), nil)

			mockUaac.EXPECT().CreateExternalUser("cwashburn", "cwashburn@testdomain.com", "cn=cwashburn", "ldap").Return(nil)
			mockUaac.EXPECT().CreateExternalUser("cwashburn1", "cwashburn1@test.io", "cn=cwashburn1", "ldap").Return(nil)
			mockUaac.EXPECT().CreateExternalUser("cwashburn2", "cwashburn2@test.io", "cn=cwashburn2", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "developers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "developers", "space1GUID").Return(nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_space1_managers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "managers", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "managers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "test_space1_auditors").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn1", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@testdomain.com", "auditors", "space1GUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn2@testdomain.com", "auditors", "space1GUID").Return(nil)

			err := spaceManager.UpdateSpaceUsers("./fixtures/user_config", "test")
			Ω(err).Should(BeNil())
		})
	})
	Context("UpdateSpaceUsers() for SAML", func() {
		spaces := []cloudcontroller.Space{
			{
				Entity: cloudcontroller.SpaceEntity{
					Name:    "space1",
					OrgGUID: "testOrgGUID",
				},
				MetaData: cloudcontroller.SpaceMetaData{
					GUID: "space1GUID",
				},
			},
			{
				Entity: cloudcontroller.SpaceEntity{
					Name:    "space2",
					OrgGUID: "testOrgGUID",
				},
				MetaData: cloudcontroller.SpaceMetaData{
					GUID: "space2GUID",
				},
			},
		}
		It("update org users where users aren't in uaac", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "saml",
			}
			uaacUsers := make(map[string]string)
			users := []l.User{
				{UserID: "cwashburn", UserDN: "cn=cwashburn", Email: "cwashburn@test.io"},
			}
			mockLdap.EXPECT().GetConfig("./fixtures/user_saml_config", "test").Return(config, nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_space1_developers").Return(users, nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "developers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "managers").Return(make(map[string]string), nil)
			mockCloudController.EXPECT().GetSpaceUsers("space1GUID", "auditors").Return(make(map[string]string), nil)

			mockUaac.EXPECT().CreateExternalUser("cwashburn@test.io", "cwashburn@test.io", "cwashburn@test.io", "saml").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@test.io", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@test.io", "developers", "space1GUID").Return(nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_space1_managers").Return(users, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn@test.io", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@test.io", "managers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "test_space1_auditors").Return(users, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn@test.io", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburn@test.io", "auditors", "space1GUID").Return(nil)

			err := spaceManager.UpdateSpaceUsers("./fixtures/user_saml_config", "test")
			Ω(err).Should(BeNil())
		})
	})
})
