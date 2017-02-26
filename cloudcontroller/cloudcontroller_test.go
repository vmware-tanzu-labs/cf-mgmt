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
			manager := NewManager("https://api.test.com", "token")
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
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces", "inline-relations-depth=1"),
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
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/v2/organizations/5678-1234/spaces", "inline-relations-depth=1"),
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

			bodyBytes := []byte(`{"instance_memory_limit":2,"memory_limit":1,"name":"name","non_basic_services_allowed":false,"organization_guid":"5678-1234","total_routes":3,"total_services":4}`)
			responsebytes, err := ioutil.ReadFile("fixtures/create-quota.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/space_quota_definitions"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, string(responsebytes)),
				),
			)
			guid, err := manager.CreateSpaceQuota(orgGUID, "name", 1, 2, 3, 4, false)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("601d30e6-f16f-4c3d-88ab-723f7c51184a"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/space_quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.CreateSpaceQuota(orgGUID, "name", 1, 2, 3, 4, false)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UpdateSpaceQuota()", func() {

		It("should be successful", func() {

			bodyBytes := []byte(`{"guid":"Quota-1234","instance_memory_limit":2,"memory_limit":1,"name":"name","non_basic_services_allowed":false,"organization_guid":"5678-1234","total_routes":3,"total_services":4}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, ""),
				),
			)
			err := manager.UpdateSpaceQuota(orgGUID, quotaGUID, "name", 1, 2, 3, 4, false)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/space_quota_definitions/Quota-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.UpdateSpaceQuota(orgGUID, quotaGUID, "name", 1, 2, 3, 4, false)
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
			spaceQuotas, err := manager.ListSpaceQuotas(orgGUID)
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
			_, err := manager.ListSpaceQuotas(orgGUID)
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
			quotas, err := manager.ListQuotas()
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
			_, err := manager.ListQuotas()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("CreateQuota()", func() {

		It("should be successful", func() {

			bodyBytes := []byte(`{"instance_memory_limit":2,"memory_limit":1,"name":"name","non_basic_services_allowed":false,"total_routes":3,"total_services":4}`)
			responsebytes, err := ioutil.ReadFile("fixtures/create-quota.json")
			Ω(err).ShouldNot(HaveOccurred())
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/quota_definitions"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, string(responsebytes)),
				),
			)
			guid, err := manager.CreateQuota("name", 1, 2, 3, 4, false)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(guid).Should(Equal("601d30e6-f16f-4c3d-88ab-723f7c51184a"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/v2/quota_definitions"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.CreateQuota("name", 1, 2, 3, 4, false)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})

	Context("UpdateQuota()", func() {

		It("should be successful", func() {

			bodyBytes := []byte(`{"guid":"Quota-1234","instance_memory_limit":2,"memory_limit":1,"name":"name","non_basic_services_allowed":false,"total_routes":3,"total_services":4}`)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/quota_definitions/Quota-1234"),
					VerifyContentType("application/json"),
					VerifyBody(bodyBytes),
					RespondWith(http.StatusCreated, ""),
				),
			)
			err := manager.UpdateQuota(quotaGUID, "name", 1, 2, 3, 4, false)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("PUT", "/v2/quota_definitions/Quota-1234"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.UpdateQuota(quotaGUID, "name", 1, 2, 3, 4, false)
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

})
