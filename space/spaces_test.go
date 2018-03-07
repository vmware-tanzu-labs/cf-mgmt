package space_test

/*import (
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
)*/

/*var _ = Describe("given SpaceManager", func() {
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
})*/

/*It("update ldap group users where users are not in uaac", func() {
	config := &l.Config{
		Enabled: true,
		Origin:  "ldap",
	}
	uaacUsers := make(map[string]string)
	spaceUsers := make(map[string]string)
	updateUsersInput := UpdateUsersInput{
		SpaceGUID:      "my-space-guid",
		OrgGUID:        "my-org-guid",
		Role:           "my-role",
		LdapGroupNames: []string{"ldap-group-name", "ldap-group-name-2"},
	}

	ldapGroupUsers := []l.User{l.User{
		UserDN: "user-dn",
		UserID: "user-id",
		Email:  "user@test.com",
	}}

	ldapGroupUsers2 := []l.User{l.User{
		UserDN: "user-dn2",
		UserID: "user-id2",
		Email:  "user2@test.com",
	}}

	mockCloudController.EXPECT().GetCFUsers("my-space-guid", "spaces", "my-role").Return(spaceUsers, nil)
	mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name").Return(ldapGroupUsers, nil)
	mockLdap.EXPECT().GetUserIDs(config, "ldap-group-name-2").Return(ldapGroupUsers2, nil)

	mockUaac.EXPECT().CreateExternalUser("user-id", "user@test.com", "user-dn", "ldap").Return(nil)
	mockUaac.EXPECT().CreateExternalUser("user-id2", "user2@test.com", "user-dn2", "ldap").Return(nil)
	mockCloudController.EXPECT().AddUserToOrg("user-id", "my-org-guid").Return(nil)
	mockCloudController.EXPECT().AddUserToOrg("user-id2", "my-org-guid").Return(nil)
	mockCloudController.EXPECT().AddUserToSpaceRole("user-id", "my-role", "my-space-guid").Return(nil)
	mockCloudController.EXPECT().AddUserToSpaceRole("user-id2", "my-role", "my-space-guid").Return(nil)

	err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
	Ω(err).Should(BeNil())
	Ω(len(uaacUsers)).Should(BeEquivalentTo(2))
	_, ok := uaacUsers["user-id"]
	Ω(ok).Should(BeTrue())

	_, ok = uaacUsers["user-id2"]
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
		SpaceGUID:      "my-space-guid",
		OrgGUID:        "my-org-guid",
		Role:           "my-role",
		LdapGroupNames: []string{"ldap-group-name"},
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
		SpaceGUID:      "my-space-guid",
		OrgGUID:        "my-org-guid",
		Role:           "my-role",
		LdapGroupNames: []string{"ldap-group-name"},
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
		SpaceGUID:      "my-space-guid",
		OrgGUID:        "my-org-guid",
		Role:           "my-role",
		LdapGroupNames: []string{"ldap-group-name"},
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
	mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn", "my-role").Return(nil)
	mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn1", "my-role").Return(nil)
	mockCloudController.EXPECT().RemoveCFUserByUserName("my-space-guid", "spaces", "cwashburn2", "my-role").Return(nil)
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
		SpaceName:      "space-name",
		SpaceGUID:      "space-guid",
		OrgName:        "org-name",
		OrgGUID:        "org-guid",
		Role:           "space-role-name",
		LdapGroupNames: []string{"ldap-group-name"},
		RemoveUsers:    true,
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
	mockCloudController.EXPECT().RemoveCFUserByUserName("space-guid", "spaces", "alex.j.smith@example.com", "space-role-name").Return(nil)
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

It("adding users to uaac based on saml", func() {
	config := &l.Config{
		Enabled: false,
		Origin:  "https://saml.example.com",
	}

	uaacUsers := make(map[string]string)
	uaacUsers["chris.a.washburn@example.com"] = "cwashburn-uaac-guid"
	uaacUsers["joe.h.fitzy@example.com"] = "jfitzy-uaac-guid"

	spaceUsers := make(map[string]string)
	spaceUsers["chris.a.washburn@example.com"] = "cwashburn-space-user-guid"
	spaceUsers["joe.h.fitzy@example.com"] = "jfitzy-space-user-guid"

	updateUsersInput := UpdateUsersInput{
		SpaceName:   "space-name",
		SpaceGUID:   "space-guid",
		OrgName:     "org-name",
		OrgGUID:     "org-guid",
		Role:        "space-role-name",
		SamlUsers:   []string{"chris.a.washburn@example.com", "joe.h.fitzy@example.com", "test@test.com"},
		RemoveUsers: true,
	}

	mockCloudController.EXPECT().GetCFUsers("space-guid", "spaces", "space-role-name").Return(spaceUsers, nil)
	mockUaac.EXPECT().CreateExternalUser("test@test.com", "test@test.com", "test@test.com", "https://saml.example.com").Return(nil)
	mockCloudController.EXPECT().AddUserToOrg("test@test.com", "org-guid").Return(nil)
	mockCloudController.EXPECT().AddUserToSpaceRole("test@test.com", "space-role-name", "space-guid").Return(nil)
	err := userManager.UpdateSpaceUsers(config, uaacUsers, updateUsersInput)
	Ω(err).Should(BeNil())
	Ω(uaacUsers).Should(HaveKey("test@test.com"))
})*/
