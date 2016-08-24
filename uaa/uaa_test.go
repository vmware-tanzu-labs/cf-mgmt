package uaa_test

import (
	"fmt"
	"net/http"

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
			bodyBytes := []byte(fmt.Sprintf("grant_type=password&password=%s&response_type=token&username=%s", password, userID))
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					VerifyBasicAuth("cf", ""),
					VerifyBody(bodyBytes),
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
		bodyBytes := []byte("grant_type=client_credentials&response_type=token")
		It("should successfully get a token", func() {
			server.AppendHandlers(
				CombineHandlers(
					VerifyRequest("POST", "/oauth/token"),
					VerifyBasicAuth(userID, password),
					VerifyBody(bodyBytes),
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
