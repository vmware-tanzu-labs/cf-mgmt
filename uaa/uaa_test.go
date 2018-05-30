package uaa_test

import (
	"io/ioutil"
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cf-mgmt/uaa"
)

var _ = Describe("given uaa manager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewDefaultUAAManager("test.com", "token", false)
			Ω(manager).ShouldNot(BeNil())
			uaaManager, ok := manager.(*DefaultUAAManager)
			Ω(ok).Should(BeTrue())
			Ω(uaaManager.Host).Should(Equal("https://uaa.test.com"))
			Ω(uaaManager.Token).Should(Equal("token"))
		})
	})
	var (
		server       *Server
		manager      DefaultUAAManager
		userID       string
		password     string
		token        Token
		controlToken string
		secret       string
	)

	BeforeEach(func() {
		controlToken = "basdfasdfd"
		userID = "myUSERID"
		password = "myPassword"
		secret = "my-secret"
		token = Token{
			AccessToken: controlToken,
		}
		server = NewServer()
		manager = DefaultUAAManager{
			Host:  server.URL(),
			Token: controlToken,
			Http:  NewHttpManager(),
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Context("GetCFToken()", func() {

		It("should successfully get a token", func() {
			expectedValues := url.Values{}
			expectedValues.Add("grant_type", "password")
			expectedValues.Add("password", password)
			expectedValues.Add("response_type", "token")
			expectedValues.Add("username", userID)
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					VerifyBasicAuth("cf", ""),
					VerifyForm(expectedValues),
					RespondWithJSONEncoded(http.StatusOK, token),
				),
			)
			token, err := GetCFToken(server.URL(), userID, password)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeNil())
			Ω(token).Should(Equal(controlToken))
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := GetCFToken(server.URL(), userID, password)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("GetUAACToken()", func() {
		expectedValues := url.Values{}
		expectedValues.Add("response_type", "token")
		expectedValues.Add("grant_type", "client_credentials")
		It("should successfully get a token", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					VerifyBasicAuth(userID, secret),
					VerifyForm(expectedValues),
					RespondWithJSONEncoded(http.StatusOK, token),
				),
			)
			token, err := GetUAACToken(server.URL(), userID, secret)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(token).ShouldNot(BeNil())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := GetUAACToken(server.URL(), userID, secret)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("ListUsers()", func() {
		It("should return list of users", func() {
			page1, _ := ioutil.ReadFile("fixtures/user-list-page1.json")
			page2, _ := ioutil.ReadFile("fixtures/user-list-page2.json")
			page3, _ := ioutil.ReadFile("fixtures/user-list-page3.json")
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/Users", "count=5000"),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWith(http.StatusOK, string(page1)),
				),
				CombineHandlers(
					VerifyRequest("GET", "/Users", "count=5000&startIndex=4"),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWith(http.StatusOK, string(page2)),
				),
				CombineHandlers(
					VerifyRequest("GET", "/Users", "count=5000&startIndex=7"),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWith(http.StatusOK, string(page3)),
				),
			)
			users, err := manager.ListUsers()
			Ω(err).ShouldNot(HaveOccurred())
			keys := make([]string, 0, len(users))
			for k := range users {
				keys = append(keys, k)
			}
			Ω(len(users)).Should(Equal(8))
			Ω(keys).Should(ConsistOf("foo4", "admin", "user", "cwashburn", "foo", "foo1", "foo2", "foo3"))
			Ω(server.ReceivedRequests()).Should(HaveLen(3))
		})
		It("should return an error", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("GET", "/Users", "count=5000"),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
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
			externalID := "userDN"

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/Users"),
					VerifyBody([]byte(`{"emails":[{"value":"email"}],"externalId":"userDN","origin":"ldap","userName":"user"}`)),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
		It("should successfully create user with complex dn", func() {
			userName := "asdfasdfsadf"
			userEmail := "caleb.washburn@test.com"
			externalID := `CN=Washburn\, Caleb\, asdfasdfsadf\,OU=NO-HOME-USERS,OU=BU-USA,DC=1DC,DC=com`

			server.AppendHandlers(
				CombineHandlers(
					VerifyContentType("application/json"),
					VerifyRequest("POST", "/Users"),
					VerifyBody([]byte(`{"emails":[{"value":"caleb.washburn@test.com"}],"externalId":"CN=Washburn, Caleb, asdfasdfsadf,OU=NO-HOME-USERS,OU=BU-USA,DC=1DC,DC=com","origin":"ldap","userName":"asdfasdfsadf"}`)),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})

		It("should peek", func() {
			userName := "user"
			userEmail := "email"
			externalID := "userDN"
			manager.Peek = true
			err := manager.CreateExternalUser(userName, userEmail, externalID, "ldap")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
		It("should not invoke post", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/Users"),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.CreateExternalUser("", "", "", "ldap")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(0))
		})
	})
	Context("CreateSamlUser()", func() {
		It("should successfully create user", func() {
			userName := "user@test.com"
			userEmail := "user@test.com"
			externalID := "user@test.com"
			origin := "saml"

			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/Users"),
					VerifyBody([]byte(`{"emails":[{"value":"user@test.com"}],"externalId":"user@test.com","origin":"saml","userName":"user@test.com"}`)),
					VerifyHeader(http.Header{
						"Authorization": []string{"BEARER basdfasdfd"},
					}),
					RespondWithJSONEncoded(http.StatusCreated, ""),
				),
			)
			err := manager.CreateExternalUser(userName, userEmail, externalID, origin)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
