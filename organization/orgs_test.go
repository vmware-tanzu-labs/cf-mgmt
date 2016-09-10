package organization_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	. "github.com/pivotalservices/cf-mgmt/organization"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	utils "github.com/pivotalservices/cf-mgmt/utils/mocks"
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
		mockUtils           *utils.MockManager
		mockUaac            *uaac.MockManager
		orgManager          DefaultOrgManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUtils = utils.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		orgManager = DefaultOrgManager{
			CloudController: mockCloudController,
			UAACMgr:         mockUaac,
			UtilsMgr:        mockUtils,
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

})
