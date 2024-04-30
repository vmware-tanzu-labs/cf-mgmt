package role_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	uaaclient "github.com/cloudfoundry-community/go-uaa"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/vmwarepivotallabs/cf-mgmt/role"
	"github.com/vmwarepivotallabs/cf-mgmt/role/fakes"
	uaafakes "github.com/vmwarepivotallabs/cf-mgmt/uaa/fakes"
)

var _ = Describe("given RoleManager", func() {
	var (
		roleManager *DefaultManager
		roleClient  *fakes.FakeCFRoleClient
		uaaFake     = new(uaafakes.FakeUaa)
	)
	BeforeEach(func() {
		roleClient = new(fakes.FakeCFRoleClient)
		uaaFake = new(uaafakes.FakeUaa)
	})
	Context("Role Manager", func() {
		BeforeEach(func() {
			uaaFake.ListUsersReturns([]uaaclient.User{
				{ID: "test-user-guid", Username: "test"},
			}, uaaclient.Page{StartIndex: 1, TotalResults: 1, ItemsPerPage: 500}, nil)
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
