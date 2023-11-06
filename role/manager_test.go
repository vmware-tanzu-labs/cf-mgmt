package role_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
)

var _ = Describe("given RoleManager", func() {
	var (
		roleManager *DefaultManager
		roleClient  *fakes.FakeCFRoleClient
	)
	BeforeEach(func() {
		roleClient = new(fakes.FakeCFRoleClient)
	})
	Context("Role Manager", func() {
		BeforeEach(func() {
			roleManager = &DefaultManager{
				RoleClient: roleClient,
			}
		})

		Context("Space Roles", func() {
			It("Should succeed list space roles", func() {
				roleClient.ListAllReturns([]*resource.Role{
					{GUID: "role-guid"},
				}, nil)
				results, err := roleManager.ListSpaceRoles()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(results)).To(Equal(1))
			})
			It("Should fail as duplicate guids returned", func() {
				roleClient.ListAllReturns([]*resource.Role{
					{
						GUID: "role-guid",
						Type: "space-devloper",
					},
					{
						GUID: "role-guid",
						Type: "space-manager",
					},
				}, nil)
				results, err := roleManager.ListSpaceRoles()
				Expect(err).To(HaveOccurred())
				Expect(results).To(BeNil())
			})
		})

		Context("Org Roles", func() {
			It("Should succeed list space roles", func() {
				roleClient.ListAllReturns([]*resource.Role{
					{GUID: "role-guid"},
				}, nil)
				results, err := roleManager.ListOrgRoles()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(results)).To(Equal(1))
			})
		})
	})
})
