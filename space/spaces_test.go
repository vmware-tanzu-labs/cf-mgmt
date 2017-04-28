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
	s "github.com/pivotalservices/cf-mgmt/space/mocks"
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
		mockUserMgr         *s.MockUserMgr
		spaceManager        DefaultSpaceManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		mockUserMgr = s.NewMockUserMgr(ctrl)

		spaceManager = DefaultSpaceManager{
			CloudController: mockCloudController,
			UAACMgr:         mockUaac,
			UtilsMgr:        utils.NewDefaultManager(),
			LdapMgr:         mockLdap,
			OrgMgr:          mockOrgMgr,
			UserMgr:         mockUserMgr,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("GetSpaceConfigs()", func() {
		Context("for default_config", func() {
			var config *InputUpdateSpaces
			BeforeEach(func() {
				configs, err := spaceManager.GetSpaceConfigs("./fixtures/space-defaults")
				Ω(err).Should(BeNil())
				Ω(configs).ShouldNot(BeNil())
				Ω(configs).Should(HaveLen(1))
				config = configs[0]
			})
			It("should return test space", func() {
				Ω(config.Space).Should(BeEquivalentTo("space1"))
			})
			It("should return ldap users from space and space defaults", func() {
				Ω(config.Developer.LdapUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
			})
			It("should return users from space and space defaults", func() {
				Ω(config.Developer.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
			})
			It("should return ldap group from space config only", func() {
				Ω(config.Developer.LdapGroup).Should(BeEquivalentTo("space-1-ldap-group"))
			})

			It("should return ldap users from space and space defaults", func() {
				Ω(config.Auditor.LdapUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
			})
			It("should return users from space and space defaults", func() {
				Ω(config.Auditor.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
			})
			It("should return ldap group from space config only", func() {
				Ω(config.Auditor.LdapGroup).Should(BeEquivalentTo("space-1-ldap-group"))
			})

			It("should return ldap users from space and space defaults", func() {
				Ω(config.Manager.LdapUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
			})
			It("should return users from space and space defaults", func() {
				Ω(config.Manager.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
			})
			It("should return ldap group from space config only", func() {
				Ω(config.Manager.LdapGroup).Should(BeEquivalentTo("space-1-ldap-group"))
			})
		})
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
			spaces := []*cloudcontroller.Space{
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
			spaces := []*cloudcontroller.Space{
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
			spaces := []*cloudcontroller.Space{}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().CreateSpace("space2", "testOrgGUID").Return(nil)
			err := spaceManager.CreateSpaces("./fixtures/config", "")
			Ω(err).Should(BeNil())
		})
		It("should create 1 space", func() {
			spaces := []*cloudcontroller.Space{
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
			spaces := []*cloudcontroller.Space{{
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
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return([]*cloudcontroller.Space{}, nil)
			mockLdap.EXPECT().GetConfig("./fixtures/default_config", "test_pwd").Return(config, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(config, uaacUsers,
				UpdateUsersInput{
					SpaceName:     "space1",
					SpaceGUID:     "space1GUID",
					OrgName:       "test",
					OrgGUID:       "testOrgGUID",
					Role:          "developers",
					LdapGroupName: "default_test_space1_developers",
					LdapUsers:     []string{"cwashburndefault1"},
					Users:         []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(config, uaacUsers,
				UpdateUsersInput{
					SpaceName:     "space1",
					SpaceGUID:     "space1GUID",
					OrgName:       "test",
					OrgGUID:       "testOrgGUID",
					Role:          "managers",
					LdapGroupName: "default_test_space1_managers",
					LdapUsers:     []string{"cwashburndefault1"},
					Users:         []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(config, uaacUsers,
				UpdateUsersInput{
					SpaceName:     "space1",
					SpaceGUID:     "space1GUID",
					OrgName:       "test",
					OrgGUID:       "testOrgGUID",
					Role:          "auditors",
					LdapGroupName: "default_test_space1_auditors",
					LdapUsers:     []string{"cwashburndefault1"},
					Users:         []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)

			err := spaceManager.CreateSpaces("./fixtures/default_config", "test_pwd")
			Ω(err).Should(BeNil())
		})
	})

	Context("CreateApplicationSecurityGroups()", func() {
		It("should create 1 asg", func() {
			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())
			spaces := []*cloudcontroller.Space{
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
			spaces := []*cloudcontroller.Space{
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
			spaces := []*cloudcontroller.Space{
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
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().CreateSpaceQuota("testOrgGUID", "space1", 10240, -1, 10, -1, true).Return("space1QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().CreateSpaceQuota("testOrgGUID", "space2", 10240, -1, 10, -1, true).Return("space2QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
			err := spaceManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should update 2 quota", func() {
			spaces := []*cloudcontroller.Space{
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
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().UpdateSpaceQuota("testOrgGUID", "space1QuotaGUID", "space1", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().UpdateSpaceQuota("testOrgGUID", "space2QuotaGUID", "space2", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space2GUID", "space2QuotaGUID")
			err := spaceManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})
	})

	Context("UpdateSpaces()", func() {
		spaces := []*cloudcontroller.Space{
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
})
