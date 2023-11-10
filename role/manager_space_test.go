package role_test

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/uaa"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("given RoleManager", func() {
	var (
		roleManager *DefaultManager
		roleClient  *fakes.FakeCFRoleClient
		userClient  *fakes.FakeCFUserClient
		jobClient   *fakes.FakeCFJobClient
		uaaFake     *uaafakes.FakeManager
	)
	BeforeEach(func() {
		roleClient = new(fakes.FakeCFRoleClient)
		userClient = new(fakes.FakeCFUserClient)
		jobClient = new(fakes.FakeCFJobClient)
		uaaFake = new(uaafakes.FakeManager)
	})
	Context("Role Manager", func() {
		BeforeEach(func() {
			roleManager = &DefaultManager{
				RoleClient: roleClient,
				UserClient: userClient,
				JobClient:  jobClient,
				UAAMgr:     uaaFake,
			}
		})

		Context("Remove", func() {
			BeforeEach(func() {
				uaaUsers := &uaa.Users{}
				uaaFake.ListUsersReturns(uaaUsers, nil)
				userClient.ListAllReturns([]*resource.User{
					{GUID: "test-user-guid"},
				}, nil)
				roleClient.ListAllReturns([]*resource.Role{
					{
						GUID: "role-guid-auditor",
						Type: resource.SpaceRoleAuditor.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:   resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							Space: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-space-guid"}},
							User:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-manager",
						Type: resource.SpaceRoleManager.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:   resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							Space: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-space-guid"}},
							User:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-developer",
						Type: resource.SpaceRoleDeveloper.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:   resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							Space: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-space-guid"}},
							User:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-supporter",
						Type: resource.SpaceRoleSupporter.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:   resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							Space: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-space-guid"}},
							User:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
				}, nil)
				err := roleManager.InitializeSpaceUserRolesMap()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleClient.ListAllCallCount()).To(Equal(1))
			})
			Context("RemoveSpaceAuditor", func() {
				It("should succeed", func() {
					err := roleManager.RemoveSpaceAuditor("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-auditor"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveSpaceAuditor("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveSpaceAuditor("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-auditor"))
				})
			})

			Context("RemoveSpaceManager", func() {
				It("should succeed", func() {
					err := roleManager.RemoveSpaceManager("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-manager"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveSpaceManager("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveSpaceManager("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-manager"))
				})
			})

			Context("RemoveSpaceDeveloper", func() {
				It("should succeed", func() {
					err := roleManager.RemoveSpaceDeveloper("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-developer"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveSpaceDeveloper("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveSpaceDeveloper("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-developer"))
				})
			})

			Context("RemoveSpaceSupporter", func() {
				It("should succeed", func() {
					err := roleManager.RemoveSpaceSupporter("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-supporter"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveSpaceSupporter("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveSpaceSupporter("orgName/spaceName", "test-space-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-supporter"))
				})
			})
		})

		Context("Add", func() {
			Context("SpaceAuditor", func() {
				It("Should succeed on AssociateSpaceAuditor", func() {
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceAuditor("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
				It("Should peek", func() {
					roleManager.Peek = true
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceAuditor("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(0))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
				})
				It("Should fail", func() {
					roleClient.CreateSpaceRoleReturns(nil, errors.New("error"))
					err := roleManager.AssociateSpaceAuditor("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleAuditor))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
			})

			Context("SpaceManager", func() {
				It("Should succeed on AssociateSpaceManager", func() {
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceManager("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleManager))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
				It("Should peek", func() {
					roleManager.Peek = true
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceAuditor("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(0))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
				})
				It("Should fail", func() {
					roleClient.CreateSpaceRoleReturns(nil, errors.New("error"))
					err := roleManager.AssociateSpaceManager("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleManager))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
			})

			Context("SpaceDeveloper", func() {
				It("Should succeed on AssociateSpaceDeveloper", func() {
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceDeveloper("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleDeveloper))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
				It("Should peek", func() {
					roleManager.Peek = true
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceDeveloper("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(0))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
				})
				It("Should fail", func() {
					roleClient.CreateSpaceRoleReturns(nil, errors.New("error"))
					err := roleManager.AssociateSpaceDeveloper("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleDeveloper))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
			})

			Context("SpaceSupporter", func() {
				It("Should succeed on AssociateSpaceSupporter", func() {
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceSupporter("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleSupporter))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
				It("Should peek", func() {
					roleManager.Peek = true
					roleClient.CreateSpaceRoleReturns(nil, nil)
					err := roleManager.AssociateSpaceSupporter("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).NotTo(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(0))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
				})
				It("Should fail", func() {
					roleClient.CreateSpaceRoleReturns(nil, errors.New("error"))
					err := roleManager.AssociateSpaceSupporter("orgGUID", "spaceName", "spaceGUID", "userName", "user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.CreateSpaceRoleCallCount()).To(Equal(1))
					Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
					_, spaceGUID, userGUID, roleType := roleClient.CreateSpaceRoleArgsForCall(0)
					Expect(spaceGUID).To(Equal("spaceGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(roleType).Should(Equal(resource.SpaceRoleSupporter))

					_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
					Expect(orgGUID).To(Equal("orgGUID"))
					Expect(userGUID).To(Equal("user-guid"))
					Expect(role).To(Equal(resource.OrganizationRoleUser))
				})
			})
		})
	})
})
