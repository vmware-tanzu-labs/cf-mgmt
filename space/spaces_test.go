package space_test

import (
	"fmt"
	"io/ioutil"

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
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetSpaceConfigs()", func() {
		It("should return list of 2", func() {
			configs, err := spaceManager.GetSpaceConfigs("./fixtures/config")
			Ω(err).Should(BeNil())
			Ω(configs).ShouldNot(BeNil())
			Ω(configs).Should(HaveLen(2))
		})
		It("should return configs for user info", func() {
			configs, err := spaceManager.GetSpaceConfigs("./fixtures/user_config")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(configs).ShouldNot(BeNil())
			Ω(configs).Should(HaveLen(1))
		})
		It("should return an error when no security.json file is provided", func() {
			configs, err := spaceManager.GetSpaceConfigs("./fixtures/no-security-json")
			Ω(err).Should(HaveOccurred())
			Ω(configs).Should(BeNil())
		})
		It("should return an error when malformed yaml", func() {
			configs, err := spaceManager.GetSpaceConfigs("./fixtures/bad-yml")
			Ω(err).Should(HaveOccurred())
			Ω(configs).Should(BeNil())
		})
		It("should return an error when path does not exist", func() {
			configs, err := spaceManager.GetSpaceConfigs("./fixtures/blah")
			Ω(err).Should(HaveOccurred())
			Ω(configs).Should(BeNil())
		})
	})
	Context("FindSpace()", func() {
		It("should return an space", func() {
			spaces := []cloudcontroller.Space{
				{
									Entity: cloudcontroller.SpaceEntity{
										Name: "testSpace",
									},
									MetaData: cloudcontroller.SpaceMetaData{},
								},
			}
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace")
			Ω(err).Should(BeNil())
			Ω(space).ShouldNot(BeNil())
			Ω(space.Entity.Name).Should(Equal("testSpace"))
		})
		It("should return an error if space not found", func() {
			spaces := []cloudcontroller.Space{
				{
									Entity: cloudcontroller.SpaceEntity{
										Name: "testSpace",
									},
									MetaData: cloudcontroller.SpaceMetaData{},
								},
			}
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
			Ω(space).Should(BeNil())
		})
		It("should return an error if unable to get OrgGUID", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("", fmt.Errorf("test"))
			space, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
			Ω(space).Should(BeNil())
		})
		It("should return an error if unable to get Spaces", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(nil, fmt.Errorf("test"))
			space, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
			Ω(space).Should(BeNil())
		})
	})

	Context("CreateSpaces()", func() {
		It("should create 2 spaces", func() {
			spaces := []cloudcontroller.Space{}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().CreateSpace("space2", "testOrgGUID").Return(nil)
			err := spaceManager.CreateSpaces("./fixtures/config", "")
			Ω(err).Should(BeNil())
		})
		It("should create 1 space", func() {
			spaces := []cloudcontroller.Space{
				{
									Entity: cloudcontroller.SpaceEntity{
										Name: "space1",
									},
									MetaData: cloudcontroller.SpaceMetaData{},
								},
			}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().CreateSpace("space2", "testOrgGUID").Return(nil)
			err := spaceManager.CreateSpaces("./fixtures/config", "")
			Ω(err).Should(BeNil())
		})
		It("should create error if unable to get orgGUID", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("", fmt.Errorf("test"))
			err := spaceManager.CreateSpaces("./fixtures/config", "")
			Ω(err).Should(HaveOccurred())
		})
	})
	Context("CreateSpaces()", func() {
		It("should create 1 spaces with default users", func() {
			config := &l.Config{
				Enabled: true,
				Origin:  "ldap",
			}
			spaces := []cloudcontroller.Space{{
							Entity: cloudcontroller.SpaceEntity{
								Name:    "space1",
								OrgGUID: "testOrgGUID",
							},
							MetaData: cloudcontroller.SpaceMetaData{
								GUID: "space1GUID",
							},
						},
			}
			uaacUsers := make(map[string]string)
			uaacUsers["cwashburn"] = "cwashburn"
			users := []l.User{}
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return([]cloudcontroller.Space{}, nil)
			mockLdap.EXPECT().GetConfig("./fixtures/default_config", "test_pwd").Return(config, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)

			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)

			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockLdap.EXPECT().GetUserIDs(config, "default_test_space1_developers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburndefault1").Return(&l.User{UserID: "cwashburndefault1", UserDN: "cn=cwashburndefault1", Email: "cwashburndefault1@test.io"}, nil)

			mockUaac.EXPECT().CreateExternalUser("cwashburndefault1", "cwashburndefault1@test.io", "cn=cwashburndefault1", "ldap").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1", "developers", "space1GUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1@testdomain.com", "developers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "default_test_space1_managers").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburndefault1").Return(&l.User{UserID: "cwashburndefault1", UserDN: "cn=cwashburndefault1", Email: "cwashburndefault1@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1", "managers", "space1GUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1@testdomain.com", "managers", "space1GUID").Return(nil)

			mockLdap.EXPECT().GetUserIDs(config, "default_test_space1_auditors").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburndefault1").Return(&l.User{UserID: "cwashburndefault1", UserDN: "cn=cwashburndefault1", Email: "cwashburndefault1@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1", "auditors", "space1GUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburndefault1@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToSpaceRole("cwashburndefault1@testdomain.com", "auditors", "space1GUID").Return(nil)

			err := spaceManager.CreateSpaces("./fixtures/default_config", "test_pwd")
			Ω(err).Should(BeNil())
		})
	})

	Context("CreateApplicationSecurityGroups()", func() {
		It("should create 1 asg", func() {
			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())
			spaces := []cloudcontroller.Space{
				{
									Entity: cloudcontroller.SpaceEntity{
										Name: "space1",
									},
									MetaData: cloudcontroller.SpaceMetaData{
										GUID: "space1GUID",
									},
								},
			}
			sgs := make(map[string]string)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
			mockCloudController.EXPECT().CreateSecurityGroup("test-space1", string(bytes)).Return("SGGUID", nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should create update 1 asg", func() {
			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())
			spaces := []cloudcontroller.Space{
				{
									Entity: cloudcontroller.SpaceEntity{
										Name: "space1",
									},
									MetaData: cloudcontroller.SpaceMetaData{
										GUID: "space1GUID",
									},
								},
			}
			sgs := make(map[string]string)
			sgs["test-space1"] = "SGGUID"
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
			mockCloudController.EXPECT().UpdateSecurityGroup("SGGUID", "test-space1", string(bytes)).Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
		})
	})

	Context("CreateQuotas()", func() {
		It("should create 2 quotas", func() {
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
			quotas := make(map[string]string)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSpaceQuotas("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().CreateSpaceQuota("testOrgGUID", "space1", 10240, -1, 10, -1, true).Return("space1QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSpaceQuotas("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().CreateSpaceQuota("testOrgGUID", "space2", 10240, -1, 10, -1, true).Return("space2QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
			err := spaceManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should update 2 quota", func() {
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
			quotas := make(map[string]string)
			quotas["space1"] = "space1QuotaGUID"
			quotas["space2"] = "space2QuotaGUID"
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSpaceQuotas("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().UpdateSpaceQuota("testOrgGUID", "space1QuotaGUID", "space1", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListSpaceQuotas("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().UpdateSpaceQuota("testOrgGUID", "space2QuotaGUID", "space2", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
			err := spaceManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})
	})

	Context("UpdateSpaces()", func() {
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
		It("should turn on allow ssh", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space2GUID").Return(nil)

			err := spaceManager.UpdateSpaces("./fixtures/config")
			Ω(err).Should(BeNil())
		})
		It("should not modify anything", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)

			err := spaceManager.UpdateSpaces("./fixtures/config-sshoff")
			Ω(err).Should(BeNil())
		})
		It("should error when UpdateSpaceSSH errors", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(fmt.Errorf("test"))

			err := spaceManager.UpdateSpaces("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})
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
