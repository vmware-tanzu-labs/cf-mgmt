package export_test

import (
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller"
	ccmock "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	. "github.com/pivotalservices/cf-mgmt/export"
	"github.com/pivotalservices/cf-mgmt/organization"
	sp "github.com/pivotalservices/cf-mgmt/space"

	"github.com/pivotalservices/cf-mgmt/uaac"
	uaacmock "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	"github.com/pivotalservices/cf-mgmt/utils"
)

func cloudControllerOrgUserMock(mockController *ccmock.MockManager, entityGUID string, mangers, billingManagers, auditors map[string]string) {
	mockController.EXPECT().GetCFUsers(entityGUID, "organizations", "managers").Return(mangers, nil)
	mockController.EXPECT().GetCFUsers(entityGUID, "organizations", "billing_managers").Return(billingManagers, nil)
	mockController.EXPECT().GetCFUsers(entityGUID, "organizations", "auditors").Return(auditors, nil)
}

func cloudControllerSpaceUserMock(mockController *ccmock.MockManager, entityGUID string, managers, developers, auditors map[string]string) {
	mockController.EXPECT().GetCFUsers(entityGUID, "spaces", "managers").Return(managers, nil)
	mockController.EXPECT().GetCFUsers(entityGUID, "spaces", "developers").Return(developers, nil)
	mockController.EXPECT().GetCFUsers(entityGUID, "spaces", "auditors").Return(auditors, nil)
}

var _ = Describe("Export manager", func() {
	Describe("Create new manager", func() {
		It("should return new manager", func() {
			ctrl := gomock.NewController(test)
			manager := NewExportManager("config", uaacmock.NewMockManager(ctrl), ccmock.NewMockManager(ctrl))
			Ω(manager).ShouldNot(BeNil())
		})
	})
	var (
		ctrl           *gomock.Controller
		mockController *ccmock.MockManager
		mockUaac       *uaacmock.MockManager
		exportManager  Manager
		excludedOrgs   map[string]string
		excludedSpaces map[string]string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockController = ccmock.NewMockManager(ctrl)
		mockUaac = uaacmock.NewMockManager(ctrl)
		exportManager = NewExportManager("test/config", mockUaac, mockController)
		excludedOrgs = make(map[string]string)
		excludedSpaces = make(map[string]string)
	})

	AfterEach(func() {
		ctrl.Finish()
		os.RemoveAll("test")
	})

	Context("Export Config", func() {
		It("Exports Org configuration", func() {

			orgId := "org1-1234"
			spaceId := "dev-1234"
			userIDToUserMap := make(map[string]uaac.User, 0)
			orgs := make([]*cc.Org, 0)
			user1 := uaac.User{ID: "1", Origin: "ldap", UserName: "user1"}
			user2 := uaac.User{ID: "2", Origin: "uaa", UserName: "user2"}
			userIDToUserMap["user1"] = user1
			userIDToUserMap["user2"] = user2

			org1 := &cc.Org{Entity: cc.OrgEntity{Name: "org1"}, MetaData: cc.OrgMetaData{GUID: orgId}}
			space := &cc.Space{Entity: cc.SpaceEntity{Name: "dev"}, MetaData: cc.SpaceMetaData{GUID: spaceId}}
			orgs = append(orgs, org1)
			spaces := make([]*cc.Space, 0)
			spaces = append(spaces, space)

			mockUaac.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListSpaces(orgId).Return(spaces, nil)
			cloudControllerOrgUserMock(mockController, orgId, map[string]string{"user1": "1", "user2": "2"}, map[string]string{}, map[string]string{})
			cloudControllerSpaceUserMock(mockController, spaceId, map[string]string{}, map[string]string{"user1": "1", "user2": "2"}, map[string]string{})

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)
			Ω(err).Should(BeNil())
			orgDetails := &organization.InputUpdateOrgs{}
			err = utils.NewDefaultManager().LoadFile("test/config/org1/orgConfig.yml", orgDetails)
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))
			Ω(len(orgDetails.Manager.Users)).Should(BeEquivalentTo(1))
			Ω(orgDetails.Manager.Users[0]).Should(Equal("user2"))
			Ω(len(orgDetails.Manager.LdapUsers)).Should(BeEquivalentTo(1))
			Ω(orgDetails.Manager.LdapUsers[0]).Should(Equal("user1"))
			Ω(len(orgDetails.BillingManager.Users)).Should(BeEquivalentTo(0))
			Ω(len(orgDetails.Auditor.Users)).Should(BeEquivalentTo(0))

			spaceDetails := &sp.InputUpdateSpaces{}
			err = utils.NewDefaultManager().LoadFile("test/config/org1/dev/spaceConfig.yml", spaceDetails)
			Ω(err).Should(BeNil())
			Ω(spaceDetails.Org).Should(Equal("org1"))
			Ω(spaceDetails.Space).Should(Equal("dev"))

			Ω(len(spaceDetails.Developer.Users)).Should(BeEquivalentTo(1))
			Ω(spaceDetails.Developer.Users[0]).Should(Equal("user2"))
			Ω(len(spaceDetails.Developer.LdapUsers)).Should(BeEquivalentTo(1))
			Ω(spaceDetails.Developer.LdapUsers[0]).Should(Equal("user1"))
			Ω(len(spaceDetails.Manager.Users)).Should(BeEquivalentTo(0))
			Ω(len(spaceDetails.Auditor.Users)).Should(BeEquivalentTo(0))
		})
	})
})
