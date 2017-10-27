package securitygroup_test

import (
	"io/ioutil"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cc "github.com/pivotalservices/cf-mgmt/cloudcontroller/mocks"
	"github.com/pivotalservices/cf-mgmt/config"
	ldap "github.com/pivotalservices/cf-mgmt/ldap/mocks"
	o "github.com/pivotalservices/cf-mgmt/organization/mocks"
	. "github.com/pivotalservices/cf-mgmt/securitygroup"
	s "github.com/pivotalservices/cf-mgmt/space/mocks"
	uaac "github.com/pivotalservices/cf-mgmt/uaac/mocks"
	"github.com/pivotalservices/cf-mgmt/utils"
)

var _ = Describe("given SecurityGroupManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token", "uaacToken", config.NewManager("./fixtures/asg-config"))
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
		securityManager     DefaultSecurityGroupManager
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(test)
		mockCloudController = cc.NewMockManager(ctrl)
		mockLdap = ldap.NewMockManager(ctrl)
		mockUaac = uaac.NewMockManager(ctrl)
		mockOrgMgr = o.NewMockManager(ctrl)
		mockUserMgr = s.NewMockUserMgr(ctrl)

		securityManager = DefaultSecurityGroupManager{
			Cfg:             config.NewManager("./fixtures/asg-config"),
			CloudController: mockCloudController,
			//			UAACMgr:         mockUaac,
			UtilsMgr: utils.NewDefaultManager(),
			//		LdapMgr:         mockLdap,
			//	OrgMgr:          mockOrgMgr,
			//			UserMgr:         mockUserMgr,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	// Find securityGroup
	// Create SecurityGroup
	// Update SecurityGroup
	// Delete SecurityGroup????? - question if this is required.

	//  package space_test
	/*
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
	*/
	var _ = Describe("given SecurityGroupManager", func() {
		Describe("create new manager", func() {
			It("should return new manager", func() {
				manager := NewManager("test.com", "token", "uaacToken", config.NewManager("./fixtures/asg-config"))
				Ω(manager).ShouldNot(BeNil())
			})
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		Context("FindSecurityGroup()", func() {
			It("should return a SecurityGroup", func() {
				/*	securityGroups := []*cloudcontroller.SecurityGroup{
					{
						Entity: cloudcontroller.SecurityGroupEntity{
							Name: "testSecurityGroup",
						},
						MetaData: cloudcontroller.SecurityGroupMetaData{},
					},
				}*/

				m := make(map[string]string)
				m["testSecurityGroup"] = "testSecurityGroup"
				mockCloudController.EXPECT().ListSecurityGroups().Return(m, nil)
				securityGroup, err := securityManager.FindSecurityGroup("testSecurityGroup")
				Expect(err).Should(BeNil())
				Expect(securityGroup).ShouldNot(BeNil())
				//Ω(securityGroup.Entity.Name).Should(Equal("testSecurityGroup"))
			})

			It("should return an error if security group not found", func() {
				m := make(map[string]string)
				m["testSecurityGroup"] = "testSecurityGroup"
				mockCloudController.EXPECT().ListSecurityGroups().Return(m, nil)
				securityGroup, err := securityManager.FindSecurityGroup("NotThere")
				Expect(err).Should(HaveOccurred())
				Expect(securityGroup).ShouldNot(BeNil())
			})
		})

		Context("CreateApplicationSecurityGroups()", func() {

			It("should create 1 asg", func() {
				bytes, e := ioutil.ReadFile("./fixtures/asg-config/asgs/test-asg.json")
				Expect(e).Should(BeNil())
				sgs := make(map[string]string)
				//			mockOrgMgr.EXPECT().GetOrgGUID("test").Return("testOrgGUID", nil)
				mockCloudController.EXPECT().ListSecurityGroups().Return(sgs, nil)
				mockCloudController.EXPECT().CreateSecurityGroup("test-asg", string(bytes)).Return("SGGUID", nil)
				//				mockCloudController.EXPECT().AssignSecurityGroupToSpace("space1GUID", "SGGUID").Return(nil)
				err := securityManager.CreateApplicationSecurityGroups("./fixtures/asg-config")
				Expect(err).Should(BeNil())
			})

			/*	It("should create update 1 asg", func() {
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
				err := securityManager.CreateApplicationSecurityGroups("./fixtures/config")
				Ω(err).Should(BeNil())
			})*/
		})

	})

	/*


		Context("CreateSpaces()", func() {
			BeforeEach(func() {
				spaceManager.Cfg = config.NewManager("./fixtures/config")
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

				spaceManager.Cfg = config.NewManager("./fixtures/default_config")
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

	})*/

})
