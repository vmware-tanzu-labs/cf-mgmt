package uaa_test

import (
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
			manager := NewDefaultUAAManager("test.com", "userID")
			Ω(manager).ShouldNot(BeNil())
			uaaManager, ok := manager.(*DefaultUAAManager)
			Ω(ok).Should(BeTrue())
			Ω(uaaManager.Host).Should(Equal("https://uaa.test.com"))
			Ω(uaaManager.UserID).Should(Equal("userID"))
		})
	})
	var (
		server       *Server
		manager      DefaultUAAManager
		userID       string
		password     string
		token        Token
		controlToken string
	)

	BeforeEach(func() {
		controlToken = "basdfasdfd"
		userID = "myUSERID"
		password = "myPassword"
		token = Token{
			AccessToken: controlToken,
		}
		server = NewServer()
		manager = DefaultUAAManager{
			Host:   server.URL(),
			UserID: userID,
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
			token, err := manager.GetCFToken(password)
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
			_, err := manager.GetCFToken(password)
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
					VerifyBasicAuth(userID, password),
					VerifyForm(expectedValues),
					RespondWithJSONEncoded(http.StatusOK, token),
				),
			)
			token, err := manager.GetUAACToken(password)
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
			_, err := manager.GetUAACToken(password)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
