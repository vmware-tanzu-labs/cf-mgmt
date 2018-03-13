package space_test

import (
	"errors"
	"fmt"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"

	ldapfakes "github.com/pivotalservices/cf-mgmt/ldap/fakes"
	orgfakes "github.com/pivotalservices/cf-mgmt/organization/fakes"
	"github.com/pivotalservices/cf-mgmt/space"
	spacefakes "github.com/pivotalservices/cf-mgmt/space/fakes"
	uaafakes "github.com/pivotalservices/cf-mgmt/uaa/fakes"
)

var _ = Describe("given SpaceManager", func() {
	var (
		fakeLdap     *ldapfakes.FakeManager
		fakeUaa      *uaafakes.FakeManager
		fakeOrgMgr   *orgfakes.FakeManager
		fakeClient   *spacefakes.FakeCFClient
		spaceManager space.DefaultManager
	)

	BeforeEach(func() {
		fakeLdap = new(ldapfakes.FakeManager)
		fakeUaa = new(uaafakes.FakeManager)
		fakeOrgMgr = new(orgfakes.FakeManager)
		fakeClient = new(spacefakes.FakeCFClient)
		spaceManager = space.DefaultManager{
			Cfg:     config.NewManager("./fixtures/config"),
			Client:  fakeClient,
			UAAMgr:  fakeUaa,
			OrgMgr:  fakeOrgMgr,
			LdapMgr: fakeLdap,
			Peek:    false,
		}
	})

	Context("FindSpace()", func() {
		It("should return an space", func() {
			spaces := []cfclient.Space{
				{
					Name: "testSpace",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
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
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})

		It("should return an error if unable to get OrgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", fmt.Errorf("test"))
			_, err := spaceManager.FindSpace("testOrg", "testSpace2")
			Ω(err).Should(HaveOccurred())
		})
		It("should return an error if unable to get Spaces", func() {
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(nil, fmt.Errorf("test"))
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
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).Should(Succeed())

			Expect(fakeClient.CreateSpaceCallCount()).Should(Equal(2))
			var spaceNames []string
			spaceRequest := fakeClient.CreateSpaceArgsForCall(0)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			spaceRequest = fakeClient.CreateSpaceArgsForCall(1)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			spaceNames = append(spaceNames, spaceRequest.Name)
			Expect(spaceNames).Should(ConsistOf([]string{"space1", "space2"}))
		})

		It("should create 1 space", func() {
			spaces := []cfclient.Space{
				{
					Name: "space1",
				},
			}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)

			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).Should(Succeed())
			Expect(fakeClient.CreateSpaceCallCount()).Should(Equal(1))
			spaceRequest := fakeClient.CreateSpaceArgsForCall(0)
			Expect(spaceRequest.OrganizationGuid).Should(Equal("testOrgGUID"))
			Expect(spaceRequest.Name).Should(Equal("space2"))
		})

		It("should create error if unable to get orgGUID", func() {
			fakeOrgMgr.GetOrgGUIDReturns("", errors.New("error1"))
			Ω(spaceManager.CreateSpaces("./fixtures/config", "")).ShouldNot(Succeed())
		})
	})

	/*Context("CreateApplicationSecurityGroups()", func() {
		It("should bind a named asg", func() {
			spaceManager.Cfg = config.NewManager("./fixtures/asg-config")

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
			sgs := make(map[string]cfclient.SecGroup)
			sgs["test-asg"] = cfclient.SecGroup{Name: "test-asg", Guid: "SGGZZUID"}
			sgs["test-space1"] = cfclient.SecGroup{Name: "test-space1", Guid: "SGGUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)

			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.UpdateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(2))
		})

		It("should create 1 asg", func() {
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
			sgs := make(map[string]cfclient.SecGroup)
			sgs["foo"] = cfclient.SecGroup{Name: "foo", Guid: "SG-FOO-GUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)
			fakeSecMgr.CreateSecurityGroupReturns(&cfclient.SecGroup{}, nil)
			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.CreateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(1))
		})

		It("should create update 1 asg", func() {
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
			sgs := make(map[string]cfclient.SecGroup)
			sgs["test-space1"] = cfclient.SecGroup{Name: "test-space1", Guid: "SGGUID"}
			sgs["foo"] = cfclient.SecGroup{Name: "foo", Guid: "SG-FOO-GUID"}
			fakeOrgMgr.GetOrgGUIDReturns("testOrgGUID", nil)
			fakeClient.ListSpacesByQueryReturns(spaces, nil)
			fakeSecMgr.ListNonDefaultSecurityGroupsReturns(sgs, nil)

			err := spaceManager.CreateApplicationSecurityGroups("./fixtures/config")
			Ω(err).Should(BeNil())
			Expect(fakeSecMgr.UpdateSecurityGroupCallCount()).Should(Equal(1))
			Expect(fakeSecMgr.AssignSecurityGroupToSpaceCallCount()).Should(Equal(1))
		})
	})*/

	/*Context("CreateQuotas()", func() {
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

	It("update ldap group users where users are not in uaac", func() {
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
	*/
})
