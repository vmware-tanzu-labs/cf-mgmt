package organization_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	l "github.com/pivotalservices/cf-mgmt/ldap"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	. "github.com/pivotalservices/cf-mgmt/organization"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	"github.com/pivotalservices/cf-mgmt/utils"
)

var _ = Describe("given OrgManager", func() {
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
		orgManager          DefaultOrgManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		orgManager = DefaultOrgManager{
			CloudController: mockCloudController,
			UAACMgr:         mockUaac,
			UtilsMgr:        utils.NewDefaultManager(),
			LdapMgr:         mockLdap,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})
	Context("FindOrg()", func() {
		It("should return an org", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			org, err := orgManager.FindOrg("test")
			Ω(err).Should(BeNil())
			Ω(org).ShouldNot(BeNil())
			Ω(org.Entity.Name).Should(Equal("test"))
		})
	})
	It("should return an error for unfound org", func() {
		orgs := []*cloudcontroller.Org{}
		mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
		org, err := orgManager.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
		Ω(org).Should(BeNil())
	})
	It("should return an error", func() {
		mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
		org, err := orgManager.FindOrg("test")
		Ω(err).ShouldNot(BeNil())
		Ω(org).Should(BeNil())
	})

	Context("GetOrgGUID()", func() {
		It("should return an GUID", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "theGUID",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			guid, err := orgManager.GetOrgGUID("test")
			Ω(err).Should(BeNil())
			Ω(guid).ShouldNot(BeNil())
			Ω(guid).Should(Equal("theGUID"))
		})
	})
	It("should return an error", func() {
		mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
		guid, err := orgManager.GetOrgGUID("test")
		Ω(err).ShouldNot(BeNil())
		Ω(guid).Should(Equal(""))
	})

	Context("DoesOrgExist()", func() {
		It("should return true", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "theGUID",
					},
				},
			}
			exists := orgManager.DoesOrgExist("test", orgs)
			Ω(exists).Should(BeTrue())
		})
	})
	It("should return false", func() {
		orgs := []*cloudcontroller.Org{
			&cloudcontroller.Org{
				Entity: cloudcontroller.OrgEntity{
					Name: "test",
				},
				MetaData: cloudcontroller.OrgMetaData{
					GUID: "theGUID",
				},
			},
		}
		exists := orgManager.DoesOrgExist("blah", orgs)
		Ω(exists).Should(BeFalse())
	})

	Context("GetOrgConfigs()", func() {
		It("should return list of 2", func() {
			configs, err := orgManager.GetOrgConfigs("./fixtures/config")
			Ω(err).Should(BeNil())
			Ω(configs).ShouldNot(BeNil())
			Ω(configs).Should(HaveLen(2))
		})
		It("should return an error when path does not exist", func() {
			configs, err := orgManager.GetOrgConfigs("./fixtures/blah")
			Ω(err).Should(HaveOccurred())
			Ω(configs).Should(BeNil())
		})
	})

	Context("CreateOrgs()", func() {
		It("should create 2", func() {
			orgs := []*cloudcontroller.Org{}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test").Return(nil)
			mockCloudController.EXPECT().CreateOrg("test2").Return(nil)
			err := orgManager.CreateOrgs("./fixtures/config")
			Ω(err).Should(BeNil())
		})
		It("should error on list orgs", func() {
			mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})
		It("should error on create org", func() {
			orgs := []*cloudcontroller.Org{}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test").Return(fmt.Errorf("test"))
			err := orgManager.CreateOrgs("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})
		It("should not create any orgs", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			err := orgManager.CreateOrgs("./fixtures/config")
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should not create test2 org", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test2").Return(nil)
			err := orgManager.CreateOrgs("./fixtures/config")
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("CreateQuotas()", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgs = []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test2OrgGUID",
					},
				},
			}
		})
		It("should create 2 quotas", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota("test", 10240, -1, 10, -1, true).Return("testQuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota("test2", 10240, -1, 10, -1, true).Return("test2QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("test2OrgGUID", "test2QuotaGUID").Return(nil)
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("list quotas returns error", func() {
			mockCloudController.EXPECT().ListQuotas().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})

		It("list orgs returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})

		It("create quota returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota("test", 10240, -1, 10, -1, true).Return("", fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})

		It("assign quota to org returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota("test", 10240, -1, 10, -1, true).Return("testQuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})

		It("should update 2 quotas", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			quotas["test2"] = "test2QuotaGUID"
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", "test", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("test2QuotaGUID", "test2", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("test2OrgGUID", "test2QuotaGUID").Return(nil)
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(BeNil())
		})

		It("update quota errors", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", "test", 10240, -1, 10, -1, true).Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})
		It("assign org to quota errors", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			mockCloudController.EXPECT().ListQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", "test", 10240, -1, 10, -1, true).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas("./fixtures/config")
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("UpdateOrgUsers()", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgs = []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test2OrgGUID",
					},
				},
			}
		})
		It("update org users where users are already in uaac", func() {
			config := &l.Config{
				Enabled: true,
			}
			uaacUsers := make(map[string]string)
			uaacUsers["cwashburn"] = "cwashburn"
			uaacUsers["cwashburn1"] = "cwashburn1"
			uaacUsers["cwashburn2"] = "cwashburn2"

			users := []l.User{
				l.User{UserID: "cwashburn", UserDN: "cn=cwashburn", Email: "cwashburn@testdomain.com"},
			}
			mockLdap.EXPECT().GetConfig("./fixtures/user_config", "test").Return(config, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_billing_managers").Return(users, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_org_managers").Return(users, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_org_auditors").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "billing_managers", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "managers", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "auditors", "testOrgGUID").Return(nil)

			err := orgManager.UpdateOrgUsers("./fixtures/user_config", "test")
			Ω(err).Should(BeNil())
		})
		It("update org users where users aren't in uaac", func() {
			config := &l.Config{
				Enabled: true,
			}
			uaacUsers := make(map[string]string)
			users := []l.User{
				l.User{UserID: "cwashburn", UserDN: "cn=cwashburn", Email: "cwashburn@testdomain.com"},
			}
			mockLdap.EXPECT().GetConfig("./fixtures/user_config", "test").Return(config, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockUaac.EXPECT().ListUsers().Return(uaacUsers, nil)
			mockUaac.EXPECT().CreateLdapUser("cwashburn", "cwashburn@testdomain.com", "cn=cwashburn").Return(nil)
			mockUaac.EXPECT().CreateLdapUser("cwashburn1", "cwashburn1@test.io", "cn=cwashburn1").Return(nil)
			mockUaac.EXPECT().CreateLdapUser("cwashburn2", "cwashburn2@test.io", "cn=cwashburn2").Return(nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_billing_managers").Return(users, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_org_managers").Return(users, nil)
			mockLdap.EXPECT().GetUserIDs(config, "test_org_auditors").Return(users, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn1").Return(&l.User{UserID: "cwashburn1", UserDN: "cn=cwashburn1", Email: "cwashburn1@test.io"}, nil)
			mockLdap.EXPECT().GetUser(config, "cwashburn2").Return(&l.User{UserID: "cwashburn2", UserDN: "cn=cwashburn2", Email: "cwashburn2@test.io"}, nil)

			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn1", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn@testdomain.com", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrg("cwashburn2@testdomain.com", "testOrgGUID").Return(nil)

			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "billing_managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "managers", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn1", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn@testdomain.com", "auditors", "testOrgGUID").Return(nil)
			mockCloudController.EXPECT().AddUserToOrgRole("cwashburn2@testdomain.com", "auditors", "testOrgGUID").Return(nil)
			err := orgManager.UpdateOrgUsers("./fixtures/user_config", "test")
			Ω(err).Should(BeNil())
		})
	})
})
