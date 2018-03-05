package cloudcontroller_test

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cf-mgmt/cloudcontroller"
)

var _ = Describe("given CloudControllerManager", func() {
	var (
		server    *Server
		manager   Manager
		token     string
		userName  string
		spaceGUID string
		orgGUID   string
		sgGUID    string
		quotaGUID string
	)

	BeforeEach(func() {
		token = "token"
		userName = "cwashburn"
		spaceGUID = "1234-5678"
		orgGUID = "5678-1234"
		sgGUID = "SG-1234"
		quotaGUID = "Quota-1234"
		server = NewServer()
		infoResponse := fmt.Sprintf(`{
			"authorization_endpoint":"%s",
			"token_endpoint":"%s",
			"doppler_logging_endpoint":"wss://doppler.v3.pcfdev.io:443",
			"routing_endpoint":"%s/routing"
		}`, server.URL(), server.URL(), server.URL())

		server.AppendHandlers(
			CombineHandlers(
				VerifyRequest("GET", "/v2/info"),
				RespondWith(http.StatusOK, infoResponse),
			),
		)

		var err error
		manager, err = NewManager(server.URL(), token, "1.0", true)
		Ω(err).ShouldNot(HaveOccurred())
		server.Reset()
	})

	AfterEach(func() {
		server.Close()
	})
	Context("AddUserToSpaceRole()", func() {
		It("should peek add user to space role", func() {
			err := manager.AddUserToSpaceRole(userName, "SpaceDeveloper", spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("AddUserToOrg()", func() {

		It("should peek", func() {
			err := manager.AddUserToOrg(userName, orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

	})
	Context("UpdateSpaceSSH()", func() {
		It("should peek", func() {
			err := manager.UpdateSpaceSSH(true, spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("UpdateSecurityGroup()", func() {
		It("should peek", func() {

			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			err = manager.UpdateSecurityGroup(sgGUID, "test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("CreateSecurityGroup()", func() {

		It("should peek", func() {

			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			guid, err := manager.CreateSecurityGroup("test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("dry-run-security-group-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("AssignSecurityGroupToSpace()", func() {

		It("should peek", func() {
			err := manager.AssignSecurityGroupToSpace(spaceGUID, sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("AssignSecurityGroupToRunning()", func() {

		It("should peek", func() {
			err := manager.AssignRunningSecurityGroup(sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("UnassignSecurityGroupToRunning()", func() {

		It("should peek", func() {
			err := manager.UnassignRunningSecurityGroup(sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("AssignSecurityGroupToStaging()", func() {

		It("should peek", func() {
			err := manager.AssignStagingSecurityGroup(sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("UnassignSecurityGroupToStaging()", func() {

		It("should peek", func() {
			err := manager.UnassignStagingSecurityGroup(sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("AssignQuotaToSpace()", func() {

		It("should peek", func() {
			err := manager.AssignQuotaToSpace(spaceGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("CreateSpaceQuota()", func() {

		It("should peek", func() {
			spaceQuota := SpaceQuotaEntity{
				OrgGUID: orgGUID,
				QuotaEntity: QuotaEntity{
					Name:                    "name",
					MemoryLimit:             1,
					InstanceMemoryLimit:     2,
					TotalRoutes:             3,
					TotalServices:           4,
					TotalPrivateDomains:     5,
					TotalReservedRoutePorts: 6,
					TotalServiceKeys:        7,
					AppInstanceLimit:        8,
					PaidServicePlansAllowed: false,
				},
			}
			guid, err := manager.CreateSpaceQuota(spaceQuota)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("dry-run-space-quota-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("UpdateSpaceQuota()", func() {

		It("should peek", func() {
			spaceQuota := SpaceQuotaEntity{
				OrgGUID: orgGUID,
				QuotaEntity: QuotaEntity{
					Name:                    "name",
					MemoryLimit:             1,
					InstanceMemoryLimit:     2,
					TotalRoutes:             3,
					TotalServices:           4,
					TotalPrivateDomains:     5,
					TotalReservedRoutePorts: 6,
					TotalServiceKeys:        7,
					AppInstanceLimit:        8,
					PaidServicePlansAllowed: false,
				},
			}
			err := manager.UpdateSpaceQuota(quotaGUID, spaceQuota)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("CreateOrg()", func() {
		It("should peek", func() {
			err := manager.CreateOrg("test")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("DeleteOrg()", func() {
		It("should peek", func() {
			err := manager.DeleteOrg("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("DeleteSpace()", func() {
		It("should peek", func() {
			err := manager.DeleteSpace("some-guid-for-a-space")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

	})

	Context("AddUserToOrgRole()", func() {

		It("should peek", func() {
			err := manager.AddUserToOrgRole(userName, "OrgManager", orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

	})

	Context("CreateQuota()", func() {

		It("should peek", func() {
			quota := QuotaEntity{
				Name:                    "name",
				MemoryLimit:             1,
				InstanceMemoryLimit:     2,
				TotalRoutes:             3,
				TotalServices:           4,
				TotalPrivateDomains:     5,
				TotalReservedRoutePorts: 6,
				TotalServiceKeys:        7,
				AppInstanceLimit:        8,
				PaidServicePlansAllowed: false,
			}
			guid, err := manager.CreateQuota(quota)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("dry-run-quota-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("UpdateQuota()", func() {

		It("should peek", func() {

			quota := QuotaEntity{
				Name:                    "name",
				MemoryLimit:             1,
				InstanceMemoryLimit:     2,
				TotalRoutes:             3,
				TotalServices:           4,
				TotalPrivateDomains:     5,
				TotalReservedRoutePorts: 6,
				TotalServiceKeys:        7,
				AppInstanceLimit:        8,
				PaidServicePlansAllowed: false,
			}
			err := manager.UpdateQuota(quotaGUID, quota)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

	})

	Context("AssignQuotaToOrg()", func() {

		It("should peek", func() {
			err := manager.AssignQuotaToOrg(orgGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("DeletePrivateDomain()", func() {
		It("should peek", func() {
			err := manager.DeletePrivateDomain("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})

	Context("CreatePrivateDomain()", func() {

		It("should peek", func() {
			guid, err := manager.CreatePrivateDomain(orgGUID, "test.com")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(BeEquivalentTo("dry-run-private-domain-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

	})

	Context("SharePrivateDomain()", func() {
		It("should peek", func() {
			err := manager.SharePrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))

		})
	})

	Context("UnsharePrivateDomain()", func() {
		It("should peek", func() {
			err := manager.RemoveSharedPrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})
})
