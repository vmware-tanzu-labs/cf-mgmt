package export_test

/*import (
	"io/ioutil"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller"
	ccmock "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
	. "github.com/pivotalservices/cf-mgmt/export"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/uaa"
	uaamock "github.com/pivotalservices/cf-mgmt/uaa/mocks"
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
			manager := NewExportManager("config", uaamock.NewMockManager(ctrl), ccmock.NewMockManager(ctrl))
			Ω(manager).ShouldNot(BeNil())
		})
	})
	var (
		ctrl           *gomock.Controller
		mockController *ccmock.MockManager
		mockUaa        *uaamock.MockManager
		exportManager  Manager
		configManager  config.Manager
		excludedOrgs   map[string]string
		excludedSpaces map[string]string
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockController = ccmock.NewMockManager(ctrl)
		mockUaa = uaamock.NewMockManager(ctrl)
		exportManager = NewExportManager("test/config", mockUaa, mockController)
		configManager = config.NewManager("test/config")
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
			userIDToUserMap := make(map[string]uaa.User, 0)
			orgs := make([]cfclient.Org, 0)
			user1 := uaa.User{ID: "1", Origin: "ldap", UserName: "user1"}
			user2 := uaa.User{ID: "2", Origin: "uaa", UserName: "user2"}
			userIDToUserMap["user1"] = user1
			userIDToUserMap["user2"] = user2

			org1 := cfclient.Org{Name: "org1", Guid: orgId}
			space := cfclient.Space{Name: "dev", Guid: spaceId}
			orgs = append(orgs, org1)
			spaces := make([]cfclient.Space, 0)
			spaces = append(spaces, space)

			securityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			defaultSecurityGroups := make(map[string]cc.SecurityGroupInfo, 0)

			mockUaa.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListIsolationSegments().Return([]cfclient.IsolationSegment{}, nil)
			mockController.EXPECT().ListOrgOwnedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListOrgSharedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListNonDefaultSecurityGroups().Return(securityGroups, nil)
			mockController.EXPECT().ListDefaultSecurityGroups().Return(defaultSecurityGroups, nil)
			mockController.EXPECT().ListSpaces(orgId).Return(spaces, nil)
			mockController.EXPECT().ListSpaceSecurityGroups(spaceId).Return(map[string]string{
				"foo": "foo-guid",
				"bar": "bar-guid",
			}, nil)
			cloudControllerOrgUserMock(mockController, orgId, map[string]string{"user1": "1", "user2": "2"}, map[string]string{}, map[string]string{})
			cloudControllerSpaceUserMock(mockController, spaceId, map[string]string{}, map[string]string{"user1": "1", "user2": "2"}, map[string]string{})

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)
			Ω(err).Should(BeNil())
			orgDetails := &config.OrgConfig{}

			orgDetails, err = configManager.GetOrgConfig("org1")
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))
			Ω(len(orgDetails.Manager.Users)).Should(BeEquivalentTo(1))
			Ω(orgDetails.Manager.Users[0]).Should(Equal("user2"))
			Ω(len(orgDetails.Manager.LDAPUsers)).Should(BeEquivalentTo(1))
			Ω(orgDetails.Manager.LDAPUsers[0]).Should(Equal("user1"))
			Ω(len(orgDetails.BillingManager.Users)).Should(BeEquivalentTo(0))
			Ω(len(orgDetails.Auditor.Users)).Should(BeEquivalentTo(0))

			spaceDetails, err := configManager.GetSpaceConfig("org1", "dev")
			Ω(err).Should(BeNil())
			Ω(spaceDetails.Org).Should(Equal("org1"))
			Ω(spaceDetails.Space).Should(Equal("dev"))
			Ω(spaceDetails.ASGs).Should(ConsistOf("foo", "bar"))

			Ω(len(spaceDetails.Developer.Users)).Should(BeEquivalentTo(1))
			Ω(spaceDetails.Developer.Users[0]).Should(Equal("user2"))
			Ω(len(spaceDetails.Developer.LDAPUsers)).Should(BeEquivalentTo(1))
			Ω(spaceDetails.Developer.LDAPUsers[0]).Should(Equal("user1"))
			Ω(len(spaceDetails.Manager.Users)).Should(BeEquivalentTo(0))
			Ω(len(spaceDetails.Auditor.Users)).Should(BeEquivalentTo(0))
		})

		XIt("Exports Quota definition", func() {

			orgId := "org1-1234"
			spaceId := "dev-1234"
			userIDToUserMap := make(map[string]uaa.User, 0)
			orgs := make([]cfclient.Org, 0)
			user1 := uaa.User{ID: "1", Origin: "ldap", UserName: "user1"}
			userIDToUserMap["user1"] = user1
			orgQuotaGUID := "54gdgf45454"
			spaceQuotaGUID := "75gdgf45454"
			org1 := cfclient.Org{Name: "org1", QuotaDefinitionGuid: orgQuotaGUID, Guid: orgId}
			space := cfclient.Space{Name: "dev", QuotaDefinitionGuid: spaceQuotaGUID, AllowSSH: true, Guid: spaceId}
			orgs = append(orgs, org1)
			spaces := make([]cfclient.Space, 0)
			spaces = append(spaces, space)

			securityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			defaultSecurityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			mockUaa.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListIsolationSegments().Return([]cfclient.IsolationSegment{}, nil)
			mockController.EXPECT().ListOrgOwnedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListOrgSharedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListNonDefaultSecurityGroups().Return(securityGroups, nil)
			mockController.EXPECT().ListDefaultSecurityGroups().Return(defaultSecurityGroups, nil)
			mockController.EXPECT().ListSpaces(orgId).Return(spaces, nil)
			mockController.EXPECT().ListSpaceSecurityGroups(spaceId).Return(map[string]string{}, nil)
			cloudControllerOrgUserMock(mockController, orgId, map[string]string{"user1": "1", "user2": "2"}, map[string]string{}, map[string]string{})
			cloudControllerSpaceUserMock(mockController, spaceId, map[string]string{}, map[string]string{"user1": "1", "user2": "2"}, map[string]string{})

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)

			Ω(err).Should(BeNil())
			orgDetails, err := configManager.GetOrgConfig("org1")
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))
			Ω(orgDetails.MemoryLimit).Should(Equal(2))
			Ω(orgDetails.InstanceMemoryLimit).Should(Equal(5))

			spaceDetails, err := configManager.GetSpaceConfig("org1", "dev")
			Ω(err).Should(BeNil())
			Ω(spaceDetails.Org).Should(Equal("org1"))
			Ω(spaceDetails.Space).Should(Equal("dev"))
			Ω(spaceDetails.MemoryLimit).Should(Equal(1))
			Ω(spaceDetails.InstanceMemoryLimit).Should(Equal(6))
			Ω(spaceDetails.AllowSSH).Should(BeTrue())
		})

		It("Exports Space security group definition", func() {
			sgRules := `[
  {
    "protocol": "udp",
    "ports": "8080",
    "destination": "198.41.191.47/1"
  },
  {
    "protocol": "tcp",
    "ports": "8080",
    "destination": "198.41.191.47/1"
  }
]`
			orgId := "org1-1234"
			spaceId := "dev-1234"
			userIDToUserMap := make(map[string]uaa.User, 0)
			orgs := make([]cfclient.Org, 0)
			user1 := uaa.User{ID: "1", Origin: "ldap", UserName: "user1"}
			userIDToUserMap["user1"] = user1
			org1 := cfclient.Org{Name: "org1", Guid: orgId}
			space := cfclient.Space{Name: "dev", AllowSSH: true, Guid: spaceId}
			orgs = append(orgs, org1)
			spaces := make([]cfclient.Space, 0)
			spaces = append(spaces, space)

			securityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			securityGroups["org1-dev"] = cc.SecurityGroupInfo{GUID: "sgGUID"}

			defaultSecurityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			defaultSecurityGroups["default"] = cc.SecurityGroupInfo{GUID: "sg-default-GUID"}

			mockUaa.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListIsolationSegments().Return([]cfclient.IsolationSegment{}, nil)
			mockController.EXPECT().ListOrgOwnedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListOrgSharedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListNonDefaultSecurityGroups().Return(securityGroups, nil)
			mockController.EXPECT().ListDefaultSecurityGroups().Return(defaultSecurityGroups, nil)
			mockController.EXPECT().GetSecurityGroupRules("sgGUID").Return([]byte(sgRules), nil)
			mockController.EXPECT().GetSecurityGroupRules("sg-default-GUID").Return([]byte(sgRules), nil)
			mockController.EXPECT().ListSpaceSecurityGroups(spaceId).Return(map[string]string{}, nil)
			mockController.EXPECT().ListSpaces(orgId).Return(spaces, nil)
			cloudControllerOrgUserMock(mockController, orgId, map[string]string{"user1": "1", "user2": "2"}, map[string]string{}, map[string]string{})
			cloudControllerSpaceUserMock(mockController, spaceId, map[string]string{}, map[string]string{"user1": "1", "user2": "2"}, map[string]string{})

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)

			Ω(err).Should(BeNil())
			orgDetails, err := configManager.GetOrgConfig("org1")
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))

			spaceDetails, err := configManager.GetSpaceConfig("org1", "dev")
			Ω(err).Should(BeNil())
			Ω(spaceDetails.Org).Should(Equal("org1"))
			Ω(spaceDetails.Space).Should(Equal("dev"))
			Ω(spaceDetails.AllowSSH).Should(BeTrue())

			data, err := ioutil.ReadFile("test/config/org1/dev/security-group.json")
			Ω(err).Should(BeNil())
			Ω(data).Should(MatchJSON(sgRules))
		})

		It("Exports global security group definition", func() {
			sgRules := `[
  {
    "protocol": "udp",
    "ports": "8080",
    "destination": "198.41.191.47/1"
  },
  {
    "protocol": "tcp",
    "ports": "8080",
    "destination": "198.41.191.47/1"
  }
]`
			orgId := "org1-1234"
			spaceId := "dev-1234"
			userIDToUserMap := make(map[string]uaa.User, 0)
			orgs := make([]cfclient.Org, 0)
			user1 := uaa.User{ID: "1", Origin: "ldap", UserName: "user1"}
			userIDToUserMap["user1"] = user1
			org1 := cfclient.Org{Name: "org1", Guid: orgId}
			space := cfclient.Space{Name: "dev", AllowSSH: true, Guid: spaceId}
			orgs = append(orgs, org1)
			spaces := make([]cfclient.Space, 0)
			spaces = append(spaces, space)

			securityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			securityGroups["test-asg"] = cc.SecurityGroupInfo{GUID: "sgGUID"}

			defaultSecurityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			defaultSecurityGroups["test-default-asg"] = cc.SecurityGroupInfo{GUID: "sg-default-GUID"}

			mockUaa.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListIsolationSegments().Return([]cfclient.IsolationSegment{}, nil)
			mockController.EXPECT().ListOrgOwnedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListOrgSharedPrivateDomains(orgId).Return(make(map[string]string), nil)
			mockController.EXPECT().ListNonDefaultSecurityGroups().Return(securityGroups, nil)
			mockController.EXPECT().ListDefaultSecurityGroups().Return(defaultSecurityGroups, nil)
			mockController.EXPECT().GetSecurityGroupRules("sgGUID").Return([]byte(sgRules), nil)
			mockController.EXPECT().GetSecurityGroupRules("sg-default-GUID").Return([]byte(sgRules), nil)
			mockController.EXPECT().ListSpaceSecurityGroups(spaceId).Return(map[string]string{}, nil)
			mockController.EXPECT().ListSpaces(orgId).Return(spaces, nil)
			cloudControllerOrgUserMock(mockController, orgId, map[string]string{"user1": "1", "user2": "2"}, map[string]string{}, map[string]string{})
			cloudControllerSpaceUserMock(mockController, spaceId, map[string]string{}, map[string]string{"user1": "1", "user2": "2"}, map[string]string{})

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)

			Ω(err).Should(BeNil())
			orgDetails, err := configManager.GetOrgConfig("org1")
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))

			spaceDetails, err := configManager.GetSpaceConfig("org1", "dev")
			Ω(err).Should(BeNil())
			Ω(spaceDetails.Org).Should(Equal("org1"))
			Ω(spaceDetails.Space).Should(Equal("dev"))
			Ω(spaceDetails.AllowSSH).Should(BeTrue())

			data, err := ioutil.ReadFile("test/config/asgs/test-asg.json")
			Ω(err).Should(BeNil())
			Ω(data).Should(MatchJSON(sgRules))
		})

		It("Skips excluded orgs from export", func() {

			orgId1 := "org1"
			orgId2 := "org2"
			userIDToUserMap := make(map[string]uaa.User, 0)
			orgs := make([]cfclient.Org, 0)
			user1 := uaa.User{ID: "1", Origin: "ldap", UserName: "user1"}
			userIDToUserMap["user1"] = user1

			org1 := cfclient.Org{Name: "org1", Guid: orgId1}
			org2 := cfclient.Org{Name: "org2", Guid: orgId2}

			orgs = append(orgs, org1)
			orgs = append(orgs, org2)

			securityGroups := make(map[string]cc.SecurityGroupInfo, 0)
			defaultSecurityGroups := make(map[string]cc.SecurityGroupInfo, 0)

			mockUaa.EXPECT().UsersByID().Return(userIDToUserMap, nil)
			mockController.EXPECT().ListOrgs().Return(orgs, nil)
			mockController.EXPECT().ListIsolationSegments().Return([]cfclient.IsolationSegment{}, nil)
			mockController.EXPECT().ListOrgOwnedPrivateDomains(orgId1).Return(make(map[string]string), nil)
			mockController.EXPECT().ListOrgSharedPrivateDomains(orgId1).Return(make(map[string]string), nil)
			mockController.EXPECT().ListNonDefaultSecurityGroups().Return(securityGroups, nil)
			mockController.EXPECT().ListDefaultSecurityGroups().Return(defaultSecurityGroups, nil)
			mockController.EXPECT().ListSpaces(orgId1).Return([]cfclient.Space{}, nil)
			cloudControllerOrgUserMock(mockController, orgId1, map[string]string{}, map[string]string{}, map[string]string{})
			excludedOrgs = map[string]string{orgId2: orgId2}

			err := exportManager.ExportConfig(excludedOrgs, excludedSpaces)

			Ω(err).Should(BeNil())
			orgDetails, err := configManager.GetOrgConfig("org1")
			Ω(err).Should(BeNil())
			Ω(orgDetails.Org).Should(Equal("org1"))

			_, err = configManager.GetOrgConfig("org2")
			Ω(err).Should(Not(BeNil()))

		})
	})
})
*/
