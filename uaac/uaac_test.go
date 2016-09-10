package uaac_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cf-mgmt/uaac"
)

var _ = Describe("given uaac manager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewManager("test.com", "token")
			Ω(manager).ShouldNot(BeNil())
			uaacManager, ok := manager.(*DefaultUAACManager)
			Ω(ok).Should(BeTrue())
			Ω(uaacManager.Host).Should(Equal("https://uaa.test.com"))
			Ω(uaacManager.UUACToken).Should(Equal("token"))
		})
	})
	var (
		server  *ghttp.Server
		manager DefaultUAACManager
		token   string
	)

	BeforeEach(func() {
		token = "secret"
		server = ghttp.NewServer()
		manager = DefaultUAACManager{
			Host:      server.URL(),
			UUACToken: token,
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Context("ListUsers()", func() {
		It("should return list of users", func() {
			userList := UserList{
				Users: []User{
					User{ID: "ID1", Name: "Test1"},
					User{ID: "ID2", Name: "Test2"},
				},
			}
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/Users"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusOK, userList),
				),
			)
			users, err := manager.ListUsers()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(users)).Should(Equal(2))
			Ω(users["test1"]).Should(Equal("ID1"))
			Ω(users["test2"]).Should(Equal("ID2"))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/Users"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.ListUsers()
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("CreateLdapUser()", func() {
		It("should successfully create user", func() {
			userName := "user"
			userEmail := "email"
			userDN := "userDN"

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/Users"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.CreateLdapUser(userName, userEmail, userDN)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/Users"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.CreateLdapUser("", "", "")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
