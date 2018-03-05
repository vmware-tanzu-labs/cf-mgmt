package space_test

import (
	"fmt"
	"io/ioutil"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
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
	uaa "github.com/pivotalservices/cf-mgmt/uaa/mocks"
)

var _ = Describe("given SpaceManager", func() {
	var (
		ctrl                *gomock.Controller
		mockCloudController *cc.MockManager
		mockLdap            *ldap.MockManager
		mockuaa             *uaa.MockManager
		mockOrgMgr          *o.MockManager
		mockUserMgr         *s.MockUserMgr
		spaceManager        DefaultSpaceManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockuaa = uaa.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		mockUserMgr = s.NewMockUserMgr(ctrl)

		spaceManager = DefaultSpaceManager{
			Cfg:             config.NewManager("./fixtures/config"),
			CloudController: mockCloudController,
			UAAMgr:          mockuaa,
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
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			space, err := spaceManager.FindSpace("testOrg", "testSpace")
			Ω(err).Should(BeNil())
			Ω(space).ShouldNot(BeNil())
			Ω(space.Name).Should(Equal("testSpace"))
		})
		It("should return an error if space not found", func() {
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error if unable to get OrgGUID", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("", fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error if unable to get Spaces", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("testOrg").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(nil, fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())

		})
	})

	Context("CreateSpaces()", func() {
		BeforeEach(func() {
			spaceManager.Cfg = config.NewManager("./fixtures/config")
		})

		It("should create 2 spaces", func() {
			spaces := []cfclient.Space{}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().CreateSpace("space2", "testOrgGUID").Return(nil)
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).Should(Succeed())
		})

		It("should create 1 space", func() {
			spaces := []cfclient.Space{
				{
					Name: "space1",
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
			spaces := []cfclient.Space{{
				Name:             "space1",
				OrganizationGuid: "testOrgGUID",
				Guid:             "space1GUID",
			},
			}
			uaaUsers := make(map[string]string)
			uaaUsers["cwashburn"] = "cwashburn"
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return([]cfclient.Space{}, nil)
			mockLdap.EXPECT().GetConfig("./fixtures/default_config", "test_pwd").Return(ldapCfg, nil)
			mockCloudController.EXPECT().CreateSpace("space1", "testOrgGUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockuaa.EXPECT().ListUsers().Return(uaaUsers, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaaUsers,
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
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaaUsers,
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
			mockUserMgr.EXPECT().UpdateSpaceUsers(ldapCfg, uaaUsers,
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

			spaceManager.Cfg = config.NewManager("./fixtures/default_config")
			Ω(spaceManager.CreateSpaces("./fixtures/default_config", "test_pwd")).Should(Succeed())
		})
	})

	Context("CreateApplicationSecurityGroups()", func() {
		It("should bind a named asg", func() {

			spaceManager = DefaultSpaceManager{
				Cfg:             config.NewManager("./fixtures/asg-config"),
				CloudController: mockCloudController,
				UAAMgr:          mockuaa,
				LdapMgr:         mockLdap,
				OrgMgr:          mockOrgMgr,
				UserMgr:         mockUserMgr,
			}

			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())

			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-asg"] = cloudcontroller.SecurityGroupInfo{GUID: "SGGZZUID"}
			sgs["test-space1"] = cloudcontroller.SecurityGroupInfo{GUID: "SGGUID"}

			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil).Times(2)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil).Times(2)
			mockCloudController.EXPECT().ListNonDefaultSecurityGroups().Return(sgs, nil)
			mockCloudController.EXPECT().UpdateSecurityGroup("SGGUID", "test-space1", string(bytes)).Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGZZUID").Return(nil)

			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should create 1 asg", func() {
			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())
			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["foo"] = cloudcontroller.SecurityGroupInfo{GUID: "SG-FOO-GUID"}
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil).Times(2)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil).Times(2)
			mockCloudController.EXPECT().ListNonDefaultSecurityGroups().Return(sgs, nil)
			mockCloudController.EXPECT().CreateSecurityGroup("test-space1", string(bytes)).Return("SGGUID", nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space2GUID", "SG-FOO-GUID").Return(nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should create update 1 asg", func() {
			bytes, e := ioutil.ReadFile("./fixtures/config/test/space1/security-group.json")
			Ω(e).Should(BeNil())
			spaces := []cfclient.Space{
				{
					Name: "space1",
					Guid: "space1GUID",
				},
				{
					Name: "space2",
					Guid: "space2GUID",
				},
			}
			sgs := make(map[string]cloudcontroller.SecurityGroupInfo)
			sgs["test-space1"] = cloudcontroller.SecurityGroupInfo{GUID: "SGGUID"}
			sgs["foo"] = cloudcontroller.SecurityGroupInfo{GUID: "SG-FOO-GUID"}
			mockCloudController.EXPECT().ListNonDefaultSecurityGroups().Return(sgs, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil).Times(2)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil).Times(2)
			mockCloudController.EXPECT().UpdateSecurityGroup("SGGUID", "test-space1", string(bytes)).Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
			mockCloudController.EXPECT().AssignSecurityGroupToSpace("space2GUID", "SG-FOO-GUID").Return(nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")

			Ω(err).Should(BeNil())
		})
	})

	Context("CreateQuotas()", func() {
		It("should create 2 quotas", func() {
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
				},
				{
					Name:             "space2",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space2GUID",
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
			spaces := []cfclient.Space{
				{
					Name:             "space1",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space1GUID",
				},
				{
					Name:             "space2",
					OrganizationGuid: "testOrgGUID",
					Guid:             "space2GUID",
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
		spaces := []cfclient.Space{
			{
				Name:             "space1",
				OrganizationGuid: "testOrgGUID",
				Guid:             "space1GUID",
			},
			{
				Name:             "space2",
				OrganizationGuid: "testOrgGUID",
				Guid:             "space2GUID",
			},
		}
		It("should turn on allow ssh", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space1GUID").Return(nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockCloudController.EXPECT().UpdateSpaceSSH(true, "space2GUID").Return(nil)

			spaceManager.Cfg = config.NewManager("./fixtures/config")
			err := spaceManager.UpdateSpaces("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("should not modify anything", func() {
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)
			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
			mockCloudController.EXPECT().ListSpaces("testOrgGUID").Return(spaces, nil)

			spaceManager.Cfg = config.NewManager("./fixtures/config-sshoff")
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
			spaceManager.Cfg = config.NewManager("./fixtures/config-delete")
		})

		It("should delete 1 and skip 1", func() {
			spaces := []cfclient.Space{
				cfclient.Space{
					Name: "space1",
					Guid: "space1-guid",
				},
				cfclient.Space{
					Name: "space2",
					Guid: "space2-guid",
				},
				cfclient.Space{
					Name: "space3",
					Guid: "space3-guid",
				},
			}
			mockOrgMgr.EXPECT().FindOrg("test2").Return(cfclient.Org{
				Name: "test2",
				Guid: "test2-org-guid",
			}, nil)
			mockCloudController.EXPECT().ListSpaces("test2-org-guid").Return(spaces, nil)
			mockCloudController.EXPECT().DeleteSpace("space3-guid").Return(nil)
			Ω(spaceManager.DeleteSpaces("./fixtures/config-delete")).Should(Succeed())
		})
	})
})
