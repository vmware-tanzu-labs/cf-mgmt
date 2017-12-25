package cloudcontroller_test

import (
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cf-mgmt/cloudcontroller"
	httpmanager "github.com/pivotalservices/cf-mgmt/http"
)

var _ = Describe("given CloudControllerManager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("https://api.test.com", "token", false)
			Ω(manager).ShouldNot(BeNil())
			cloudControllerManager := manager.(*DefaultManager)
			Ω(cloudControllerManager.Host).Should(Equal("https://api.test.com"))
		})
	})

	var (
		server    *Server
		manager   DefaultManager
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
		manager = DefaultManager{
			Host:  server.URL(),
			Token: token,
			HTTP:  httpmanager.NewManager(),
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Context("AddUserToSpaceRole()", func() {

		It("should add user to space role", func() {
			bodyBytes := []byte(`{"username":"cwashburn"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/spaces/1234-5678/SpaceDeveloper"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AddUserToSpaceRole(userName, "SpaceDeveloper", spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek add user to space role", func() {
			manager.Peek = true
			err := manager.AddUserToSpaceRole(userName, "SpaceDeveloper", spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/spaces/1234-5678/SpaceDeveloper"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AddUserToSpaceRole(userName, "SpaceDeveloper", spaceGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("AddUserToOrg()", func() {

		It("should be successful", func() {
			bodyBytes := []byte(`{"username":"cwashburn"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234/users"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AddUserToOrg(userName, orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			manager.Peek = true
			err := manager.AddUserToOrg(userName, orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234/users"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AddUserToOrg(userName, orgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("UpdateSpaceSSH()", func() {

		It("should be successful", func() {
			bodyBytes := []byte(`{"allow_ssh":true}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/spaces/1234-5678"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.UpdateSpaceSSH(true, spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.UpdateSpaceSSH(true, spaceGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/spaces/1234-5678"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.UpdateSpaceSSH(true, spaceGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("CreateSpace()", func() {

		It("should be successful", func() {
			bodyBytes := []byte(`{"name":"test","organization_guid":"5678-1234"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/spaces"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
			err := manager.CreateSpace("test", orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/spaces"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.CreateSpace("test", orgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("ListSpaces()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/spaces.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			spaces, err := manager.ListSpaces(orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(spaces).ShouldNot(BeNil())
			Ω(spaces).Should(HaveLen(3))
			for _, space := range spaces {
				Ω(space.Entity).ShouldNot(BeNil())
				Ω(space.Entity.AllowSSH).ShouldNot(BeNil())
				Ω(space.Entity.Name).ShouldNot(BeNil())
				Ω(space.Entity.OrgGUID).ShouldNot(BeNil())
				Ω(space.MetaData).ShouldNot(BeNil())
				Ω(space.MetaData.GUID).ShouldNot(BeNil())
			}
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should paginate through all results", func() {
			bytes, err := ioutil.ReadFile("fixtures/spaces-with-paging.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			bytes, err = ioutil.ReadFile("fixtures/spaces.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			spaces, err := manager.ListSpaces(orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(spaces).ShouldNot(BeNil())
			Ω(spaces).Should(HaveLen(4))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListSpaces(orgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListSecurityGroups()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/security-groups.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/security_groups"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			securityGroups, err := manager.ListSecurityGroups()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(securityGroups).ShouldNot(BeNil())
			Ω(securityGroups).Should(HaveLen(6))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/security_groups"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListSecurityGroups()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("GetSecurityGroupRules()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/security-group.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/security_groups/sg-guid"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			rules, err := manager.GetSecurityGroupRules("sg-guid")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(rules).Should(MatchJSON(`[
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
	    ]`))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/security_groups/sg-guid"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.GetSecurityGroupRules("sg-guid")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListSpaceSecurityGroups()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/space-security-groups.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/spaces/space-guid/security_groups"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			ruleNames, err := manager.ListSpaceSecurityGroups("space-guid")
			keys := make([]string, 0, len(ruleNames))
			for k := range ruleNames {
				keys = append(keys, k)
			}
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ruleNames).Should(HaveLen(3))
			Ω(keys).Should(ConsistOf("public_networks", "dns", "all_pcfdev"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/spaces/space-guid/security_groups"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListSpaceSecurityGroups("space-guid")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UpdateSecurityGroup()", func() {

		It("should be successful", func() {
			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			bodyBytes := []byte(`{"name":"test","rules":[{"destination":"10.68.192.1-10.68.192.49","protocol":"all"}]}`)

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/security_groups/SG-1234"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err = manager.UpdateSecurityGroup(sgGUID, "test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			err = manager.UpdateSecurityGroup(sgGUID, "test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/security_groups/SG-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.UpdateSecurityGroup(sgGUID, "test", "contents")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreateSecurityGroup()", func() {

		It("should be successful", func() {
			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			bodyBytes := []byte(`{"name":"test","rules":[{"destination":"10.68.192.1-10.68.192.49","protocol":"all"}]}`)

			responsebytes, err := ioutil.ReadFile("fixtures/create-asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/security_groups"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, string(responsebytes)),
				),
			)
			guid, err := manager.CreateSecurityGroup("test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("601d30e6-f16f-4c3d-88ab-723f7c51184a"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			contentbytes, err := ioutil.ReadFile("fixtures/asg.json")
			Ω(err).ShouldNot(HaveOccurred())
			guid, err := manager.CreateSecurityGroup("test", string(contentbytes))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("dry-run-security-group-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/security_groups"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.CreateSecurityGroup("test", "contents")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("AssignSecurityGroupToSpace()", func() {

		It("should be successful", func() {

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/security_groups/SG-1234/spaces/1234-5678"),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AssignSecurityGroupToSpace(spaceGUID, sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			manager.Peek = true
			err := manager.AssignSecurityGroupToSpace(spaceGUID, sgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/security_groups/SG-1234/spaces/1234-5678"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AssignSecurityGroupToSpace(spaceGUID, sgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("AssignQuotaToSpace()", func() {

		It("should be successful", func() {

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234/spaces/1234-5678"),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AssignQuotaToSpace(spaceGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			manager.Peek = true
			err := manager.AssignQuotaToSpace(spaceGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234/spaces/1234-5678"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AssignQuotaToSpace(spaceGUID, quotaGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreateSpaceQuota()", func() {

		It("should be successful", func() {

			bodyBytes := `{
					"organization_guid": "5678-1234",
					"name": "name",
					"memory_limit": 1,
					"instance_memory_limit": 2,
					"total_routes": 3,
					"total_services": 4,
					"total_private_domains": 5,
					"total_reserved_route_ports": 6,
					"total_service_keys": 7,
					"app_instance_limit": 8,
					"non_basic_services_allowed": false
				}`
			responsebytes, err := ioutil.ReadFile("fixtures/create-quota.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/space_quota_definitions"),
					VerifyContentType("application/json"),
					VerifyJSON(bodyBytes),
					RespondWith(http.StatusCreated, string(responsebytes)),
				),
			)
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
			Ω(guid).Should(Equal("601d30e6-f16f-4c3d-88ab-723f7c51184a"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			manager.Peek = true
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
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/space_quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			spaceQuota := SpaceQuotaEntity{
				OrgGUID: orgGUID,
			}
			_, err := manager.CreateSpaceQuota(spaceQuota)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UpdateSpaceQuota()", func() {

		It("should be successful", func() {

			bodyBytes := `{
					"organization_guid": "5678-1234",
					"name": "name",
					"memory_limit": 1,
					"instance_memory_limit": 2,
					"total_routes": 3,
					"total_services": 4,
					"total_private_domains": 5,
					"total_reserved_route_ports": 6,
					"total_service_keys": 7,
					"app_instance_limit": 8,
					"non_basic_services_allowed": false
				}`
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234"),
					VerifyContentType("application/json"),
					VerifyJSON(bodyBytes),
					RespondWith(http.StatusCreated, ""),
				),
			)
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
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
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
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			spaceQuota := SpaceQuotaEntity{
				OrgGUID: orgGUID,
			}
			err := manager.UpdateSpaceQuota(quotaGUID, spaceQuota)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListSpaceQuotas()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/space-quotas.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/space_quota_definitions"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			spaceQuotas, err := manager.ListAllSpaceQuotasForOrg(orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(spaceQuotas).ShouldNot(BeNil())
			Ω(spaceQuotas).Should(HaveLen(2))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/space_quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListAllSpaceQuotasForOrg(orgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreateOrg()", func() {

		It("should be successful", func() {
			bodyBytes := []byte(`{"name":"test"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/organizations"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
			err := manager.CreateOrg("test")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.CreateOrg("test")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/organizations"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.CreateOrg("test")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("DeleteOrg()", func() {
		It("should be successful", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/organizations/22d428d0-014a-473b-87b2-131367a31248", "recursive=true"),
					RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
			err := manager.DeleteOrg("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.DeleteOrg("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/organizations/22d428d0-014a-473b-87b2-131367a31248"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.DeleteOrg("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("DeleteSpace()", func() {
		It("should be successful", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/spaces/some-guid-for-a-space", "recursive=true"),
					RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
			err := manager.DeleteSpace("some-guid-for-a-space")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.DeleteSpace("some-guid-for-a-space")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/spaces/some-guid-for-a-space"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.DeleteSpace("some-guid-for-a-space")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListOrgs()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/orgs.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			orgs, err := manager.ListOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(orgs).ShouldNot(BeNil())
			Ω(orgs).Should(HaveLen(1))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should retrieve all results when there are more results than the initial size of 100", func() {
			bytes, err := ioutil.ReadFile("fixtures/orgs-with-paging.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			bytes, err = ioutil.ReadFile("fixtures/orgs.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			orgs, err := manager.ListOrgs()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(orgs).ShouldNot(BeNil())
			Ω(orgs).Should(HaveLen(3))
			Ω(server.ReceivedRequests()).Should(HaveLen(2))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListOrgs()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("AddUserToOrgRole()", func() {

		It("should add user to space role", func() {
			bodyBytes := []byte(`{"username":"cwashburn"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234/OrgManager"),
					VerifyBody(bodyBytes),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AddUserToOrgRole(userName, "OrgManager", orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			manager.Peek = true
			err := manager.AddUserToOrgRole(userName, "OrgManager", orgGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234/OrgManager"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AddUserToOrgRole(userName, "OrgManager", orgGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListQuotas()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/quotas.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/quota_definitions"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			quotas, err := manager.ListAllOrgQuotas()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(quotas).ShouldNot(BeNil())
			Ω(quotas).Should(HaveLen(19))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListAllOrgQuotas()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreateQuota()", func() {

		It("should be successful", func() {

			bodyBytes := `{
					"name": "name",
					"memory_limit": 1,
					"instance_memory_limit": 2,
					"total_routes": 3,
					"total_services": 4,
					"total_private_domains": 5,
					"total_reserved_route_ports": 6,
					"total_service_keys": 7,
					"app_instance_limit": 8,
					"non_basic_services_allowed": false
				}`
			responsebytes, err := ioutil.ReadFile("fixtures/create-quota.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/quota_definitions"),
					VerifyContentType("application/json"),
					VerifyJSON(bodyBytes),
					RespondWith(http.StatusCreated, string(responsebytes)),
				),
			)
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
			Ω(guid).Should(Equal("601d30e6-f16f-4c3d-88ab-723f7c51184a"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
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
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			quota := QuotaEntity{}
			_, err := manager.CreateQuota(quota)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UpdateQuota()", func() {

		It("should be successful", func() {

			bodyBytes := `{
					"name": "name",
					"memory_limit": 1,
					"instance_memory_limit": 2,
					"total_routes": 3,
					"total_services": 4,
					"total_private_domains": 5,
					"total_reserved_route_ports": 6,
					"total_service_keys": 7,
					"app_instance_limit": 8,
					"non_basic_services_allowed": false
				}`
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/quota_definitions/Quota-1234"),
					VerifyContentType("application/json"),
					VerifyJSON(bodyBytes),
					RespondWith(http.StatusCreated, ""),
				),
			)
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
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true

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
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/quota_definitions/Quota-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			quota := QuotaEntity{}
			err := manager.UpdateQuota(quotaGUID, quota)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("AssignQuotaToOrg()", func() {

		It("should be successful", func() {

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234"),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.AssignQuotaToOrg(orgGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.AssignQuotaToOrg(orgGUID, quotaGUID)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/5678-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.AssignQuotaToOrg(orgGUID, quotaGUID)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("GetCFUser()", func() {
		It("should retrieve all results when there are more results", func() {
			bytes, err := ioutil.ReadFile("fixtures/space-developers-paging.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/spaces/2ae52bf0-6f0a-4461-b683-8fa96c15d350/developers"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			bytes, err = ioutil.ReadFile("fixtures/space-developers.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/spaces/2ae52bf0-6f0a-4461-b683-8fa96c15d350/developers"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			devs, err := manager.GetCFUsers("2ae52bf0-6f0a-4461-b683-8fa96c15d350", "spaces", "developers")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(devs).ShouldNot(BeNil())
			Ω(devs).Should(HaveLen(4))
			Ω(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Context("ListAllPrivateDomains()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/all-private-domains.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListAllPrivateDomains()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(4))
			Ω(privateDomains).Should(HaveKeyWithValue("vcap.me", PrivateDomainInfo{OrgGUID: "4cf3bc47-eccd-4662-9322-7833c3bdcded", PrivateDomainGUID: "b2a35f0c-d5ad-4a59-bea7-461711d96b0d"}))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-61.example.com", PrivateDomainInfo{OrgGUID: "c262280e-0ccc-4e13-918a-6852f2d1e3a0", PrivateDomainGUID: "28db6393-cc6f-4318-a63c-f4009e8842bc"}))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-62.example.com", PrivateDomainInfo{OrgGUID: "68f69961-f751-4b52-907c-4469009fdf74", PrivateDomainGUID: "a16ffec7-5fab-4447-861e-c38da6548c6d"}))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-63.example.com", PrivateDomainInfo{OrgGUID: "8d8ed1ba-f7f3-48f1-8d9a-2dfaad91335b", PrivateDomainGUID: "4168cdaf-1586-41a6-9e5f-d8c715c332f5"}))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/private_domains"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListAllPrivateDomains()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("ListOrgOwnedPrivateDomains()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5afc3416-50ba-46db-ae2c-5e0c88cdac3f/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgOwnedPrivateDomains("5afc3416-50ba-46db-ae2c-5e0c88cdac3f")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(1))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-28.example.com", "ffcf939a-22ed-4ae5-9371-2e737bd1eb48"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-OrgA", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgA.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--A/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgOwnedPrivateDomains("org--A")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(3))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.A.example.com", "domain--A"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.B.example.com", "domain--B"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.C.example.com", "domain--C"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-orgB", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgB.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--B/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgOwnedPrivateDomains("org--B")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(3))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.D.example.com", "domain--D"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.E.example.com", "domain--E"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.F.example.com", "domain--F"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-orgC", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgC.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--C/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgOwnedPrivateDomains("org--C")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(3))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.G.example.com", "domain--G"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.H.example.com", "domain--H"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.I.example.com", "domain--I"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org_guid/private_domains"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListOrgOwnedPrivateDomains("org_guid")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	// new
	Context("ListOrgSharedPrivateDomains()", func() {

		It("should be successful", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5afc3416-50ba-46db-ae2c-5e0c88cdac3f/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgSharedPrivateDomains("5afc3416-50ba-46db-ae2c-5e0c88cdac3f")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(2))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-29.example.com", "718a97e7-db8a-46a7-874a-187b6b2a6ff1"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain-30.example.com", "f29112ca-de2e-4a8d-b343-161ba8cae556"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-orgA", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgA.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--A/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgSharedPrivateDomains("org--A")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(3))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.E.example.com", "domain--E"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.G.example.com", "domain--G"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.H.example.com", "domain--H"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-orgB", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgB.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--B/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgSharedPrivateDomains("org--B")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(1))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.B.example.com", "domain--B"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should be successful-orgC", func() {
			bytes, err := ioutil.ReadFile("fixtures/org-private-domains-orgC.json")

			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org--C/private_domains"),
					RespondWith(http.StatusOK, string(bytes)),
				),
			)
			privateDomains, err := manager.ListOrgSharedPrivateDomains("org--C")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(privateDomains).ShouldNot(BeNil())
			Ω(privateDomains).Should(HaveLen(2))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.C.example.com", "domain--C"))
			Ω(privateDomains).Should(HaveKeyWithValue("domain.E.example.com", "domain--E"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/org_guid/private_domains"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListOrgSharedPrivateDomains("org_guid")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	// new

	Context("DeletePrivateDomain()", func() {
		It("should be successful", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/private_domains/22d428d0-014a-473b-87b2-131367a31248", "async=false"),
					RespondWithJSONEncoded(http.StatusOK, ""),
				),
			)
			err := manager.DeletePrivateDomain("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.DeletePrivateDomain("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/private_domains/22d428d0-014a-473b-87b2-131367a31248", "async=false"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.DeletePrivateDomain("22d428d0-014a-473b-87b2-131367a31248")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreatePrivateDomain()", func() {
		responseBytes, _ := ioutil.ReadFile("fixtures/create-private-domain-result.json")

		It("should be successful", func() {
			bodyBytes := []byte(`{"name":"test.com","owning_organization_guid":"5678-1234"}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/private_domains"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, string(responseBytes)),
				),
			)
			guid, err := manager.CreatePrivateDomain(orgGUID, "test.com")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(BeEquivalentTo("b98aeca1-22b9-49f9-8428-3ace9ea2ba11"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			guid, err := manager.CreatePrivateDomain(orgGUID, "test.com")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(BeEquivalentTo("dry-run-private-domain-guid"))
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/private_domains"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.CreatePrivateDomain(orgGUID, "test.com")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("SharePrivateDomain()", func() {
		It("should be successful", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/1234o/private_domains/1234d"),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.SharePrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))

		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.SharePrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))

		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/organizations/1234o/private_domains/1234d"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.SharePrivateDomain("1234o", "1234d")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UnsharePrivateDomain()", func() {
		It("should be successful", func() {
			server.AppendHandlers(
				CombineHandlers(),
				VerifyRequest("DELETE", "/v2/organizations/1234o/private_domains/1234d"),
				RespondWithJSONEncoded(http.StatusNoContent, ""),
			)
			err := manager.RemoveSharedPrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should peek", func() {
			manager.Peek = true
			err := manager.RemoveSharedPrivateDomain("1234o", "1234d")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})

		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("DELETE", "/v2/organizations/1234o/private_domains/1234d"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.RemoveSharedPrivateDomain("1234o", "1234d")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

})
