package uaa_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/pivotalservices/cf-mgmt/uaa"
)

type Sample struct {
	Value string `json:"resources"`
}

var _ = Describe("given utils manager", func() {
	Describe("create new manager", func() {
		It("should return new manager", func() {
			manager := NewHttpManager()
			Ω(manager).ShouldNot(BeNil())
		})
	})
	var (
		server  *ghttp.Server
		manager HttpManager
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		manager = NewHttpManager()
	})

	AfterEach(func() {
		server.Close()
	})
	Context("HTTPGet", func() {
		output := Sample{
			Value: "blah",
		}
		It("Should return instance of object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusOK, output),
				),
			)

			target := &Sample{}
			err := manager.Get(server.URL(), "secret", target)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
			Ω(target.Value).Should(Equal("blah"))
		})
		It("Should return error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			target := &Sample{}
			err := manager.Get(server.URL(), "secret", target)
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("HTTPPost", func() {
		output := Sample{
			Value: "blah",
		}
		It("Should return instance of object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusOK, output),
				),
			)

			payload, err := json.Marshal(output)
			Ω(err).ShouldNot(HaveOccurred())
			_, err = manager.Post(server.URL(), "secret", string(payload))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))

		})
		It("Should return error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			_, err := manager.Post(server.URL(), "secret", "")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
	Context("HTTPPut", func() {
		output := Sample{
			Value: "blah",
		}
		It("Should return instance of object", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusCreated, output),
				),
			)

			payload, err := json.Marshal(output)
			Ω(err).ShouldNot(HaveOccurred())
			err = manager.Put(server.URL(), "secret", string(payload))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))

		})
		It("Should return error", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/"),
					ghttp.VerifyHeader(http.Header{
						"Authorization": []string{"BEARER secret"},
					}),
					ghttp.RespondWithJSONEncoded(http.StatusServiceUnavailable, ""),
				),
			)
			err := manager.Put(server.URL(), "secret", "")
			Ω(err).Should(HaveOccurred())
			Ω(server.ReceivedRequests()).Should(HaveLen(1))
		})
	})
})
