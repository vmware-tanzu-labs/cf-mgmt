package organization_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	. "github.com/pivotalservices/cf-mgmt/organization"
	uaa "github.com/pivotalservices/cf-mgmt/uaa/mocks"
)

var _ = Describe("given OrgManager", func() {
	var (
		ctrl                *gomock.Controller
		mockCloudController *cc.MockManager
		mockLdap            *ldap.MockManager
		mockUaa             *uaa.MockManager
		orgManager          DefaultOrgManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaa = uaa.NewMockManager(ctrl)
		orgManager = DefaultOrgManager{
			Cfg:             config.NewManager("./fixtures/config"),
			CloudController: mockCloudController,
			UAAMgr:          mockUaa,
			LdapMgr:         mockLdap,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("FindOrg()", func() {
		It("should return an org", func() {
			orgs := []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
				{Entity: cloudcontroller.OrgEntity{
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
				{
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

	Context("CreateOrgs()", func() {
		It("should create 2", func() {
			orgs := []*cloudcontroller.Org{}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test").Return(nil)
			mockCloudController.EXPECT().CreateOrg("test2").Return(nil)
			err := orgManager.CreateOrgs()
			Ω(err).Should(BeNil())
		})
		It("should error on list orgs", func() {
			mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should error on create org", func() {
			orgs := []*cloudcontroller.Org{}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test").Return(fmt.Errorf("test"))
			err := orgManager.CreateOrgs()
			Ω(err).Should(HaveOccurred())
		})
		It("should not create any orgs", func() {
			orgs := []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
		})
		It("should not create test2 org", func() {
			orgs := []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateOrg("test2").Return(nil)
			err := orgManager.CreateOrgs()
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	Context("DeleteOrgs()", func() {
		BeforeEach(func() {
			orgManager.Cfg = config.NewManager("./fixtures/config-delete")
		})

		It("should delete 1", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "system",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "system-guid",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test-guid",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test2-guid",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "redis-test-ORG-1-2017_10_04-20h06m33.481s",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "redis-guid",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().DeleteOrg("test2-guid").Return(nil)
			err := orgManager.DeleteOrgs(false)
			Ω(err).Should(BeNil())
		})
		It("should just peek", func() {
			orgs := []*cloudcontroller.Org{
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "system",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "system-guid",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test-guid",
					},
				},
				&cloudcontroller.Org{
					Entity: cloudcontroller.OrgEntity{
						Name: "test2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "test2-guid",
					},
				},
			}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			err := orgManager.DeleteOrgs(true)
			Ω(err).Should(BeNil())
		})
	})

	Context("CreateQuotas()", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgManager.Cfg = config.NewManager("./fixtures/config")

			orgs = []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				{
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
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota(cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return("testQuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota(cloudcontroller.QuotaEntity{
				Name:                    "test2",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return("test2QuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("test2OrgGUID", "test2QuotaGUID").Return(nil)
			err := orgManager.CreateQuotas()
			Ω(err).Should(BeNil())
		})

		It("list quotas returns error", func() {
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})

		It("list orgs returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(nil, fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})

		It("create quota returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota(cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return("", fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})

		It("assign quota to org returns error", func() {
			quotas := make(map[string]string)
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().CreateQuota(cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return("testQuotaGUID", nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})

		It("should update 2 quotas", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			quotas["test2"] = "test2QuotaGUID"
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("test2QuotaGUID", cloudcontroller.QuotaEntity{
				Name:                    "test2",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("test2OrgGUID", "test2QuotaGUID").Return(nil)
			err := orgManager.CreateQuotas()
			Ω(err).Should(BeNil())
		})

		It("update quota errors", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})

		It("assign org to quota errors", func() {
			quotas := make(map[string]string)
			quotas["test"] = "testQuotaGUID"
			mockCloudController.EXPECT().ListAllOrgQuotas().Return(quotas, nil)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().UpdateQuota("testQuotaGUID", cloudcontroller.QuotaEntity{
				Name:                    "test",
				MemoryLimit:             10240,
				InstanceMemoryLimit:     -1,
				TotalRoutes:             10,
				TotalServices:           -1,
				PaidServicePlansAllowed: true,
				AppInstanceLimit:        -1,
				TotalReservedRoutePorts: 0,
				TotalPrivateDomains:     -1,
				TotalServiceKeys:        -1,
			}).Return(nil)
			mockCloudController.EXPECT().AssignQuotaToOrg("testOrgGUID", "testQuotaGUID").Return(fmt.Errorf("test"))
			err := orgManager.CreateQuotas()
			Ω(err).Should(HaveOccurred())
		})
	})

	Context("CreatePrivateDomains()", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgManager.Cfg = config.NewManager("./fixtures/config-private-domains")

			orgs = []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test-2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOtherOrgGUID",
					},
				},
			}
		})
		It("should create 2 private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			orgOwnedPrivateDomains := make(map[string]string)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil)
			mockCloudController.EXPECT().CreatePrivateDomain("testOrgGUID", "test.com").Return("test.com.guid", nil)
			mockCloudController.EXPECT().CreatePrivateDomain("testOrgGUID", "test2.com").Return("test2.com.guid", nil)
			mockCloudController.EXPECT().ListOrgOwnedPrivateDomains("testOrgGUID").Return(orgOwnedPrivateDomains, nil)
			mockCloudController.EXPECT().SharePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.CreatePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should create no private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test.com"] = cloudcontroller.PrivateDomainInfo{OrgGUID: "testOrgGUID"}
			allPrivateDomains["test2.com"] = cloudcontroller.PrivateDomainInfo{OrgGUID: "testOrgGUID"}
			orgPrivateDomains := make(map[string]string)
			orgPrivateDomains["test.com"] = "test.com.guid"
			orgPrivateDomains["test2.com"] = "test2.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil)
			mockCloudController.EXPECT().ListOrgOwnedPrivateDomains("testOrgGUID").Return(orgPrivateDomains, nil)
			mockCloudController.EXPECT().SharePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.CreatePrivateDomains()
			Ω(err).Should(BeNil())
		})

		It("should create 2 private domains and delete 2 domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test3.com"] = cloudcontroller.PrivateDomainInfo{OrgGUID: "testOrgGUID"}
			allPrivateDomains["test4.com"] = cloudcontroller.PrivateDomainInfo{OrgGUID: "testOrgGUID"}
			orgPrivateDomains := make(map[string]string)
			orgPrivateDomains["test.com"] = "test.com.guid"
			orgPrivateDomains["test2.com"] = "test2.com.guid"
			orgPrivateDomains["test3.com"] = "test3.com.guid"
			orgPrivateDomains["test4.com"] = "test4.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil)
			mockCloudController.EXPECT().CreatePrivateDomain("testOrgGUID", "test.com").Return("test.com.guid", nil)
			mockCloudController.EXPECT().CreatePrivateDomain("testOrgGUID", "test2.com").Return("test2.com.guid", nil)
			mockCloudController.EXPECT().ListOrgOwnedPrivateDomains("testOrgGUID").Return(orgPrivateDomains, nil)
			mockCloudController.EXPECT().DeletePrivateDomain("test3.com.guid").Return(nil)
			mockCloudController.EXPECT().DeletePrivateDomain("test4.com.guid").Return(nil)
			mockCloudController.EXPECT().SharePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.CreatePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should error as private domain exists in other org", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test.com"] = cloudcontroller.PrivateDomainInfo{OrgGUID: "testOtherOrgGUID"}
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil)
			mockCloudController.EXPECT().SharePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.CreatePrivateDomains()
			Ω(err).Should(Not(BeNil()))
		})
	})

	Context("SharePrivateDomainsDeleteEnabled", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgManager.Cfg = config.NewManager("./fixtures/config-private-domains")

			orgs = []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test-2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOtherOrgGUID",
					},
				},
			}
		})
		It("should error when private domain doesn't exist", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			orgSharedPrivateDomains := make(map[string]string)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(1)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(BeEquivalentTo("Private Domain [test-shared.com] is not defined"))
		})
		It("should share 2 private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test-shared.comGUID"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test-shared2.comGUID"}
			orgSharedPrivateDomains := make(map[string]string)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(2)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test-shared.comGUID").Return(nil)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test-shared2.comGUID").Return(nil)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should share no private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test-shared.comGUID"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test-shared2.comGUID"}
			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared.com"] = "test.shared.com.guid"
			orgSharedPrivateDomains["test-shared2.com"] = "test.shared2.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(2)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should unshare 2 domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared.com"] = "test.shared.com.guid"
			orgSharedPrivateDomains["test-shared2.com"] = "test.shared2.com.guid"
			orgSharedPrivateDomains["test-shared3.com"] = "test.shared3.com.guid"
			orgSharedPrivateDomains["test-shared4.com"] = "test.shared4.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(2)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain("testOrgGUID", "test.shared3.com.guid").Return(nil)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain("testOrgGUID", "test.shared4.com.guid").Return(nil)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should share 2 domains and unshare 2 domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			allPrivateDomains["test-shared3.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared3.com.guid"}
			allPrivateDomains["test-shared4.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared4.com.guid"}
			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared3.com"] = "test.shared3.com.guid"
			orgSharedPrivateDomains["test-shared4.com"] = "test.shared4.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(2)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared.com.guid").Return(nil)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared2.com.guid").Return(nil)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain("testOrgGUID", "test.shared3.com.guid").Return(nil)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain("testOrgGUID", "test.shared4.com.guid").Return(nil)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			mockCloudController.EXPECT().DeletePrivateDomain(gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
	})
	Context("SharePrivateDomainsDeleteDisabled", func() {
		var orgs []*cloudcontroller.Org
		BeforeEach(func() {
			orgManager.Cfg = config.NewManager("./fixtures/config-private-domains-no-delete")

			orgs = []*cloudcontroller.Org{
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOrgGUID",
					},
				},
				{
					Entity: cloudcontroller.OrgEntity{
						Name: "test-2",
					},
					MetaData: cloudcontroller.OrgMetaData{
						GUID: "testOtherOrgGUID",
					},
				},
			}
		})
		It("should share 2 private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			allPrivateDomains["test-shared3.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared3.com.guid"}
			allPrivateDomains["test-shared4.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared4.com.guid"}
			orgSharedPrivateDomains := make(map[string]string)
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared.com.guid").Return(nil)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared2.com.guid").Return(nil)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should share no private domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			allPrivateDomains["test-shared3.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared3.com.guid"}
			allPrivateDomains["test-shared4.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared4.com.guid"}
			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared.com"] = "test.shared.com.guid"
			orgSharedPrivateDomains["test-shared2.com"] = "test.shared2.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should NOT unshare 2 domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			allPrivateDomains["test-shared3.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared3.com.guid"}
			allPrivateDomains["test-shared4.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared4.com.guid"}
			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared.com"] = "test.shared.com.guid"
			orgSharedPrivateDomains["test-shared2.com"] = "test.shared2.com.guid"
			orgSharedPrivateDomains["test-shared3.com"] = "test.shared3.com.guid"
			orgSharedPrivateDomains["test-shared4.com"] = "test.shared4.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})
		It("should share 2 domains and NOT unshare 2 domains", func() {
			allPrivateDomains := make(map[string]cloudcontroller.PrivateDomainInfo)
			allPrivateDomains["test-shared.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared.com.guid"}
			allPrivateDomains["test-shared2.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared2.com.guid"}
			allPrivateDomains["test-shared3.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared3.com.guid"}
			allPrivateDomains["test-shared4.com"] = cloudcontroller.PrivateDomainInfo{PrivateDomainGUID: "test.shared4.com.guid"}

			orgSharedPrivateDomains := make(map[string]string)
			orgSharedPrivateDomains["test-shared3.com"] = "test.shared3.com.guid"
			orgSharedPrivateDomains["test-shared4.com"] = "test.shared4.com.guid"
			mockCloudController.EXPECT().ListOrgs().Return(orgs, nil)
			mockCloudController.EXPECT().ListAllPrivateDomains().Return(allPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().ListOrgSharedPrivateDomains("testOrgGUID").Return(orgSharedPrivateDomains, nil).Times(1)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared.com.guid").Return(nil)
			mockCloudController.EXPECT().SharePrivateDomain("testOrgGUID", "test.shared2.com.guid").Return(nil)
			mockCloudController.EXPECT().RemoveSharedPrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			mockCloudController.EXPECT().CreatePrivateDomain(gomock.Any().String(), gomock.Any().String()).Times(0)
			mockCloudController.EXPECT().DeletePrivateDomain(gomock.Any().String()).Times(0)
			err := orgManager.SharePrivateDomains()
			Ω(err).Should(BeNil())
		})

	})

})
