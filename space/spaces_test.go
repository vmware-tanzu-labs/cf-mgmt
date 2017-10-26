package space_test

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
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
			manager := NewManager("test.com", "token", "uaacToken", config.NewManager("./fixtures/space-defaults", utils.NewDefaultManager()))
			Ω(manager).ShouldNot(BeNil())
		})
	})

	var (
		utilsMgr            utils.Manager
		ctrl                *gomock.Controller
		mockCloudController *cc.MockManager
		mockLdap            *ldap.MockManager
		mockUaac            *uaac.MockManager
		mockOrgMgr          *o.MockManager
		mockUserMgr         *s.MockUserMgr
		spaceManager        DefaultSpaceManager
	)

	BeforeEach(func() {
		utilsMgr = utils.NewDefaultManager()
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		mockUserMgr = s.NewMockUserMgr(ctrl)

		spaceManager = DefaultSpaceManager{
			Cfg:             config.NewManager("./fixtures/config", utilsMgr),
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
		BeforeEach(func() {
			spaceManager.Cfg = config.NewManager("./fixtures/config", utilsMgr)
		})

		It("should create 2 spaces", func() {
			spaces := []*cloudcontroller.Space{}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().CreateSpace("space2", "testOrgGUID").Return(nil)
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).Should(Succeed())
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
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).Should(Succeed())
		})

		It("should create error if unable to get orgGUID", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("", fmt.Errorf("test"))
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).ShouldNot(Succeed())
		})

		It("should create 1 spaces with default users", func() {
			ldapCfg := &l.Config{
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
			mockLdap.EXPECT().GetConfig("./fixtures/default_config", "test_pwd").Return(ldapCfg, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaacUsers,
				UpdateUsersInput{
					SpaceName:      "space1",
					SpaceGUID:      "space1GUID",
					OrgName:        "test",
					OrgGUID:        "testOrgGUID",
					Role:           "developers",
					LdapGroupNames: []string{"default_test_space1_developers"},
					LdapUsers:      []string{"cwashburndefault1"},
					Users:          []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaacUsers,
				UpdateUsersInput{
					SpaceName:      "space1",
					SpaceGUID:      "space1GUID",
					OrgName:        "test",
					OrgGUID:        "testOrgGUID",
					Role:           "managers",
					LdapGroupNames: []string{"default_test_space1_managers"},
					LdapUsers:      []string{"cwashburndefault1"},
					Users:          []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaacUsers,
				UpdateUsersInput{
					SpaceName:      "space1",
					SpaceGUID:      "space1GUID",
					OrgName:        "test",
					OrgGUID:        "testOrgGUID",
					Role:           "auditors",
					LdapGroupNames: []string{"default_test_space1_auditors"},
					LdapUsers:      []string{"cwashburndefault1"},
					Users:          []string{"cwashburndefault1@testdomain.com"},
				}).Return(nil)

			spaceManager.Cfg = config.NewManager("./fixtures/default_config", utilsMgr)
			Ω(spaceManager.CreateSpaces("./fixtures/default_config", "test_pwd")).Should(Succeed())
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
			mockCloudController.EXPECT().CreateSpaceQuota(cloudcontroller.SpaceQuotaEntity{
				OrgGUID: "testOrgGUID",
				QuotaEntity: cloudcontroller.QuotaEntity{
					Name:                    "space1",
					MemoryLimit:             10240,
					InstanceMemoryLimit:     -1,
					TotalRoutes:             10,
					TotalServices:           -1,
					PaidServicePlansAllowed: true,
					AppInstanceLimit:        -1,
					TotalReservedRoutePorts: 0,
					TotalPrivateDomains:     -1,
					TotalServiceKeys:        -1,
				}}).Return("space1QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().CreateSpaceQuota(cloudcontroller.SpaceQuotaEntity{
				OrgGUID: "testOrgGUID",
				QuotaEntity: cloudcontroller.QuotaEntity{
					Name:                    "space2",
					MemoryLimit:             10240,
					InstanceMemoryLimit:     -1,
					TotalRoutes:             10,
					TotalServices:           -1,
					PaidServicePlansAllowed: true,
					AppInstanceLimit:        -1,
					TotalReservedRoutePorts: 0,
					TotalPrivateDomains:     -1,
					TotalServiceKeys:        -1,
				}}).Return("space2QuotaGUID", nil)
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
			mockCloudController.EXPECT().UpdateSpaceQuota("space1QuotaGUID", cloudcontroller.SpaceQuotaEntity{
				OrgGUID: "testOrgGUID",
				QuotaEntity: cloudcontroller.QuotaEntity{
					Name:                    "space1",
					MemoryLimit:             10240,
					InstanceMemoryLimit:     -1,
					TotalRoutes:             10,
					TotalServices:           -1,
					PaidServicePlansAllowed: true,
					AppInstanceLimit:        -1,
					TotalReservedRoutePorts: 0,
					TotalPrivateDomains:     -1,
					TotalServiceKeys:        -1,
				}}).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToSpace("space1GUID", "space1QuotaGUID")

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().ListAllSpaceQuotasForOrg("testOrgGUID").Return(quotas, nil)
			mockCloudController.EXPECT().UpdateSpaceQuota("space2QuotaGUID", cloudcontroller.SpaceQuotaEntity{
				OrgGUID: "testOrgGUID",
				QuotaEntity: cloudcontroller.QuotaEntity{
					Name:                    "space2",
					MemoryLimit:             10240,
					InstanceMemoryLimit:     -1,
					TotalRoutes:             10,
					TotalServices:           -1,
					PaidServicePlansAllowed: true,
					AppInstanceLimit:        -1,
					TotalReservedRoutePorts: 0,
					TotalPrivateDomains:     -1,
					TotalServiceKeys:        -1,
				}}).Return(nil)
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

			spaceManager.Cfg = config.NewManager("./fixtures/config", utilsMgr)
			err := spaceManager.UpdateSpaces("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should not modify anything", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)

			spaceManager.Cfg = config.NewManager("./fixtures/config-sshoff", utilsMgr)
			Ω(spaceManager.UpdateSpaces("./fixtures/config-sshoff")).Should(Succeed())
		})

		It("should error when UpdateSpaceSSH errors", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(fmt.Errorf("test"))

			Ω(spaceManager.UpdateSpaces("./fixtures/config")).ShouldNot(Succeed())
		})
	})

	Context("DeleteSpaces()", func() {
		BeforeEach(func() {
			spaceManager.Cfg = config.NewManager("./fixtures/config-delete", utilsMgr)
		})

		It("should delete 1 and skip 1", func() {
			spaces := []*cloudcontroller.Space{
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space1",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space1-guid",
					},
				},
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space2",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space2-guid",
					},
				},
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space3",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space3-guid",
					},
				},
			}
			mockOrgMgr.EXPECT().FindOrg("test2").Return(&cloudcontroller.Org{
				Entity: cloudcontroller.OrgEntity{
					Name: "test2",
				},
				MetaData: cloudcontroller.OrgMetaData{
					GUID: "test2-org-guid",
				},
			}, nil)
			mockCloudController.EXPECT().ListSpaces("test2-org-guid").Return(spaces, nil)
			mockCloudController.EXPECT().DeleteSpace("space3-guid").Return(nil)
			Ω(spaceManager.DeleteSpaces("./fixtures/config-delete", false)).Should(Succeed())
		})

		It("should just peek", func() {
			spaces := []*cloudcontroller.Space{
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space1",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space1-guid",
					},
				},
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space2",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space2-guid",
					},
				},
				&cloudcontroller.Space{
					Entity: cloudcontroller.SpaceEntity{
						Name: "space3",
					},
					MetaData: cloudcontroller.SpaceMetaData{
						GUID: "space3-guid",
					},
				},
			}
			mockOrgMgr.EXPECT().FindOrg("test2").Return(&cloudcontroller.Org{
				Entity: cloudcontroller.OrgEntity{
					Name: "test",
				},
				MetaData: cloudcontroller.OrgMetaData{
					GUID: "test2-org-guid",
				},
			}, nil)
			mockCloudController.EXPECT().ListSpaces("test2-org-guid").Return(spaces, nil)
			Ω(spaceManager.DeleteSpaces("./fixtures/config-delete", true)).Should(Succeed())
		})
	})
})
