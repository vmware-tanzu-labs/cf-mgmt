package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
)

var _ = Describe("CF-Mgmt Config", func() {
	Context("Default Config Updater", func() {
		Context("GetOrgConfigs", func() {
			It("should return a list of 2", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(c).Should(HaveLen(2))
			})

			It("should return a list of 1", func() {
				m := config.NewManager("./fixtures/user_config")
				c, err := m.GetOrgConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(c).Should(HaveLen(1))

				org := c[0]
				Ω(org.GetAuditorGroups()).Should(ConsistOf([]string{"test_org_auditors"}))
				Ω(org.GetManagerGroups()).Should(ConsistOf([]string{"test_org_managers"}))
				Ω(org.GetBillingManagerGroups()).Should(ConsistOf([]string{"test_billing_managers", "test_billing_managers_2"}))
			})

			It("should fail when given an invalid config dir", func() {
				m := config.NewManager("./fixtures/blah")
				c, err := m.GetOrgConfigs()
				Ω(err).Should(HaveOccurred())
				Ω(c).Should(BeEmpty())
			})
		})
	})
})
