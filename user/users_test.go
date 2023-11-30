package user_test

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	orgfakes "github.com/vmwarepivotallabs/cf-mgmt/organizationreader/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/role"
	rolefakes "github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
	. "github.com/vmwarepivotallabs/cf-mgmt/user"
	"github.com/vmwarepivotallabs/cf-mgmt/user/fakes"
)

var _ = Describe("given UserSpaces", func() {
	var (
		userManager *DefaultManager
		ldapFake    *fakes.FakeLdapManager
		azureADFake *fakes.FakeAzureADManager
		uaaFake     *uaafakes.FakeManager
		fakeReader  *configfakes.FakeReader
		spaceFake   *spacefakes.FakeManager
		orgFake     *orgfakes.FakeReader
		roleMgrFake *rolefakes.FakeManager
	)
	BeforeEach(func() {
		ldapFake = new(fakes.FakeLdapManager)
		azureADFake = new(fakes.FakeAzureADManager)
		uaaFake = new(uaafakes.FakeManager)
		fakeReader = new(configfakes.FakeReader)
		spaceFake = new(spacefakes.FakeManager)
		orgFake = new(orgfakes.FakeReader)
		roleMgrFake = new(rolefakes.FakeManager)
	})
	Context("User Manager()", func() {
		BeforeEach(func() {
			userManager = &DefaultManager{
				Cfg:           fakeReader,
				UAAMgr:        uaaFake,
				LdapMgr:       ldapFake,
				AzureADMgr:    azureADFake,
				SpaceMgr:      spaceFake,
				OrgReader:     orgFake,
				Peek:          false,
				RoleMgr:       roleMgrFake,
				LdapConfig:    &config.LdapConfig{Origin: "ldap"},
				AzureADConfig: &config.AzureADConfig{},
			}

			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{}, nil)
			roleMgrFake.ListOrgUsersByRoleReturns(role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), nil)
			roleMgrFake.ListSpaceUsersByRoleReturns(role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), role.InitRoleUsers(), nil)
		})

		Context("SyncInternalUsers", func() {
			var roleUsers *role.RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-user-guid"})
				uaaUsers.Add(uaa.User{Username: "test-existing", Origin: "uaa", GUID: "test-existing-id"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{
					{Username: "test-existing", GUID: "test-existing-id"},
				}, uaaUsers)

				userManager.UAAUsers = uaaUsers
			})
			It("Should add internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					SpaceName: "spaceName",
					OrgName:   "orgName",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput, false)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(userGUID).Should(Equal("test-user-guid"))
			})

			It("Should add internal user to role and ldap user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					SpaceName: "spaceName",
					OrgName:   "orgName",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput, false)
				Expect(err).ShouldNot(HaveOccurred())
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(userGUID).Should(Equal("test-user-guid"))
			})

			It("Should not add existing internal user to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test-existing"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput, false)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(0))
			})
			It("Should error when user doesn't exist in uaa", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test2"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput, false)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("user test2 doesn't exist in origin uaa, so must add internal user first"))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: role.InitRoleUsers(),
				}
				roleMgrFake.AssociateSpaceAuditorReturns(errors.New("error"))
				err := userManager.SyncInternalUsers(roleUsers, updateUsersInput, false)
				Expect(err).Should(HaveOccurred())
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(1))
			})

		})

		Context("Remove Users", func() {
			var roleUsers *role.RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-id"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{
					{Username: "test", GUID: "test-id"},
				}, uaaUsers)
			})

			It("Should remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					OrgName:     "orgName",
					SpaceName:   "spaceName",
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  roleMgrFake.RemoveSpaceAuditor,
					RoleUsers:   role.InitRoleUsers(),
					Role:        "auditor",
				}
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.RemoveSpaceAuditorCallCount()).Should(Equal(1))

				spaceName, spaceGUID, userName, userGUID := roleMgrFake.RemoveSpaceAuditorArgsForCall(0)
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(userGUID).Should(Equal("test-id"))
			})

			It("Should not remove users", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: false,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  roleMgrFake.RemoveSpaceAuditor,
					RoleUsers:   role.InitRoleUsers(),
				}

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.RemoveSpaceAuditorCallCount()).Should(Equal(0))
			})

			It("Should skip users that match protected user pattern", func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "abcd_123_0919191", Origin: "uaa", GUID: "test-id"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{
					{Username: "abcd_123_0919191", GUID: "test-id"},
				}, uaaUsers)
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  roleMgrFake.RemoveSpaceAuditor,
					RoleUsers:   role.InitRoleUsers(),
				}

				fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
					ProtectedUsers: []string{"abcd_123_*"},
				}, nil)

				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.RemoveSpaceAuditorCallCount()).Should(Equal(0))
			})

			It("Should return error", func() {
				updateUsersInput := UsersInput{
					RemoveUsers: true,
					SpaceGUID:   "space_guid",
					OrgGUID:     "org_guid",
					RemoveUser:  roleMgrFake.RemoveSpaceAuditor,
					RoleUsers:   role.InitRoleUsers(),
				}
				roleMgrFake.RemoveSpaceAuditorReturns(errors.New("error"))
				err := userManager.RemoveUsers(roleUsers, updateUsersInput)
				Expect(err).Should(HaveOccurred())
				Expect(roleMgrFake.RemoveSpaceAuditorCallCount()).Should(Equal(1))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := &uaa.Users{}
				userMap.Add(uaa.User{Username: "test-user-guid", GUID: "test-user"})
				uaaFake.ListUsersReturns(userMap, nil)
				fakeReader.GetSpaceConfigsReturns([]config.SpaceConfig{
					{
						Space: "test-space",
						Org:   "test-org",
					},
				}, nil)
				spaceFake.FindSpaceReturns(&resource.Space{
					Name: "test-space",
					Relationships: &resource.SpaceRelationships{
						Organization: &resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "test-org-guid",
							},
						},
					},
					GUID: "test-space-guid",
				}, nil)
				userManager.LdapConfig = &config.LdapConfig{Enabled: false}
				err := userManager.UpdateSpaceUsers()
				Expect(len(err)).To(Equal(0))
			})
		})

		Context("UpdateSpaceUsers", func() {
			It("Should succeed", func() {
				userMap := &uaa.Users{}
				userMap.Add(uaa.User{Username: "test-user-guid", GUID: "test-user"})
				uaaFake.ListUsersReturns(userMap, nil)
				fakeReader.GetOrgConfigsReturns([]config.OrgConfig{
					{
						Org: "test-org",
					},
				}, nil)
				orgFake.FindOrgReturns(&resource.Organization{
					Name: "test-org",
					GUID: "test-org-guid",
				}, nil)
				userManager.LdapConfig = &config.LdapConfig{Enabled: false}
				err := userManager.UpdateOrgUsers()
				Expect(len(err)).To(Equal(0))
			})
		})
		Context("Sync Users", func() {
			var roleUsers *role.RoleUsers
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaUsers.Add(uaa.User{Username: "test", Origin: "uaa", GUID: "test-user-guid"})
				uaaUsers.Add(uaa.User{Username: "test", Origin: "ldap", GUID: "test-ldap-user-guid", ExternalID: "cn=test"})
				roleUsers, _ = role.NewRoleUsers([]*resource.User{}, uaaUsers)

				userManager.UAAUsers = uaaUsers
			})
			It("Should add internal user to role and ldap user with same name to role", func() {
				updateUsersInput := UsersInput{
					Users:     []string{"test"},
					LdapUsers: []string{"test"},
					SpaceGUID: "space_guid",
					OrgGUID:   "org_guid",
					SpaceName: "spaceName",
					OrgName:   "orgName",
					AddUser:   roleMgrFake.AssociateSpaceAuditor,
					RoleUsers: roleUsers,
				}
				userManager.LdapConfig = &config.LdapConfig{Enabled: true, Origin: "ldap"}
				err := userManager.SyncUsers(updateUsersInput)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleMgrFake.AssociateSpaceAuditorCallCount()).Should(Equal(2))
				orgGUID, spaceName, spaceGUID, userName, userGUID := roleMgrFake.AssociateSpaceAuditorArgsForCall(0)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(userGUID).Should(Equal("test-ldap-user-guid"))

				orgGUID, spaceName, spaceGUID, userName, userGUID = roleMgrFake.AssociateSpaceAuditorArgsForCall(1)
				Expect(orgGUID).Should(Equal("org_guid"))
				Expect(spaceName).Should(Equal("orgName/spaceName"))
				Expect(spaceGUID).Should(Equal("space_guid"))
				Expect(userName).Should(Equal("test"))
				Expect(userGUID).Should(Equal("test-user-guid"))
			})
		})
	})
})
