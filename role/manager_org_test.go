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

		Context("AddUserToOrg", func() {
			It("should associate user", func() {
				err := roleManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

			})

			It("should peek", func() {
				roleManager.Peek = true
				err := roleManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
			})

			It("should error", func() {
				roleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
				err := roleManager.AddUserToOrg("test-org-guid", "test", "test-user-guid")
				Expect(err).Should(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).Should(Equal("test-org-guid"))
				Expect(userGUID).Should(Equal("test-user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))
			})
		})

		Context("AssociateOrgAuditor", func() {
			It("Should succeed", func() {
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgAuditor("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(2))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, orgGUID, userGUID, role = roleClient.CreateOrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleAuditor))
			})

			It("Should fail", func() {
				roleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
				err := roleManager.AssociateOrgAuditor("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))
			})

			It("Should peek", func() {
				roleManager.Peek = true
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgAuditor("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgBillingManager", func() {
			It("Should succeed", func() {
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgBillingManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(2))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, orgGUID, userGUID, role = roleClient.CreateOrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleBillingManager))
			})

			It("Should fail", func() {
				roleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
				err := roleManager.AssociateOrgBillingManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))

				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))
			})

			It("Should peek", func() {
				roleManager.Peek = true
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgBillingManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
			})
		})

		Context("AssociateOrgManager", func() {
			It("Should succeed", func() {
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(2))
				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))

				_, orgGUID, userGUID, role = roleClient.CreateOrganizationRoleArgsForCall(1)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleManager))
			})

			It("Should fail", func() {
				roleClient.CreateOrganizationRoleReturns(nil, errors.New("error"))
				err := roleManager.AssociateOrgManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).To(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(1))

				_, orgGUID, userGUID, role := roleClient.CreateOrganizationRoleArgsForCall(0)
				Expect(orgGUID).To(Equal("orgGUID"))
				Expect(userGUID).To(Equal("user-guid"))
				Expect(role).To(Equal(resource.OrganizationRoleUser))
			})

			It("Should peek", func() {
				roleManager.Peek = true
				roleClient.CreateOrganizationRoleReturns(nil, nil)
				err := roleManager.AssociateOrgManager("orgGUID", "org", "", "userName", "user-guid")
				Expect(err).NotTo(HaveOccurred())
				Expect(roleClient.CreateOrganizationRoleCallCount()).To(Equal(0))
			})
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
						Type: resource.OrganizationRoleAuditor.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							User: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-manager",
						Type: resource.OrganizationRoleManager.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							User: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-org-user",
						Type: resource.OrganizationRoleUser.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							User: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
					{
						GUID: "role-guid-org-billing-manager",
						Type: resource.OrganizationRoleBillingManager.String(),
						Relationships: resource.RoleSpaceUserOrganizationRelationships{
							Org:  resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-org-guid"}},
							User: resource.ToOneRelationship{Data: &resource.Relationship{GUID: "test-user-guid"}},
						},
					},
				}, nil)
				err := roleManager.InitializeOrgUserRolesMap()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(roleClient.ListAllCallCount()).To(Equal(1))
			})
			Context("RemoveOrgAuditor", func() {
				It("should succeed", func() {
					err := roleManager.RemoveOrgAuditor("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-auditor"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveOrgAuditor("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveOrgAuditor("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-auditor"))
				})
			})

			Context("RemoveOrgBillingManager", func() {
				It("should succeed", func() {
					err := roleManager.RemoveOrgBillingManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-org-billing-manager"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveOrgBillingManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveOrgBillingManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-org-billing-manager"))
				})
			})

			Context("RemoveOrgManager", func() {
				It("should succeed", func() {
					roleClient.ListAllReturns([]*resource.Role{}, nil)
					err := roleManager.RemoveOrgManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-manager"))
				})

				It("should peek", func() {
					roleManager.Peek = true
					err := roleManager.RemoveOrgManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(0))
				})

				It("should error", func() {
					roleClient.DeleteReturns("", errors.New("error"))
					err := roleManager.RemoveOrgManager("orgName", "test-org-guid", "test", "test-user-guid")
					Expect(err).Should(HaveOccurred())
					Expect(roleClient.DeleteCallCount()).To(Equal(1))
					_, roleGUID := roleClient.DeleteArgsForCall(0)
					Expect(roleGUID).Should(Equal("role-guid-manager"))
				})
			})
		})

	})
})
