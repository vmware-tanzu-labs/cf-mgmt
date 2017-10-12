package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
)

var _ = Describe("CF-Mgmt Config", func() {
	Context("Protected Org Defaults", func() {
		Describe("Defaults", func() {
			It("should setup default protected orgs", func() {
				Ω(config.DefaultProtectedOrgs).Should(HaveKey("system"))
				Ω(config.DefaultProtectedOrgs).Should(HaveKey("p-spring-cloud-services"))
				Ω(config.DefaultProtectedOrgs).Should(HaveKey("splunk-nozzle-org"))
				Ω(config.DefaultProtectedOrgs).Should(HaveLen(3))
			})
		})
	})

	Context("Default Config Reader", func() {
		Context("GetASGConfigs", func() {
			It("should return a single ASG", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfgs).Should(HaveLen(1))

				cfg := cfgs[0]
				Expect(cfg.Rules).Should(BeEquivalentTo("[{\"protocol\": \"icmp\",\"destination\": \"0.0.0.0/0\"}]\n"))
			})

			It("should have a name based on the ASG filename", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfgs).Should(HaveLen(1))

				cfg := cfgs[0]
				Expect(cfg.Name).Should(BeEquivalentTo("test-asg"))

			})

			It("can optionally have a ASG name in the spaced config.", func() {
				m := config.NewManager("./fixtures/asg-defaults")

				// Get ASG Config
				asgcfgs, err := m.GetASGConfigs()
				Ω(err).ShouldNot(HaveOccurred())

				// Get space config
				cfgs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())

				cfg := cfgs[0]
				Ω(cfg.Space).Should(BeEquivalentTo("space1"))
				Expect(cfg.ASGs).Should(ConsistOf("test-asg"))
				// We expect the ASGs to reference actually defined asgs

				namedList := make([]string, len(asgcfgs))
				for i, asg := range asgcfgs {
					namedList[i] = asg.Name
				}

				for _, asg := range cfg.ASGs {
					Expect([]string{asg}).Should(ConsistOf(namedList))
				}

			})

			It("should fail when spaceconfig gives an ASG name that doesn't exist in the global definition", func() {
				m := config.NewManager("./fixtures/asg-defaults-invalid")
				_, err := m.GetSpaceConfigs()
				Ω(err).Should(HaveOccurred())

			})

		})

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

		Context("GetSpaceConfigs", func() {
			It("should return a single space", func() {
				m := config.NewManager("./fixtures/space-defaults")
				cfgs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfgs).Should(HaveLen(1))

				cfg := cfgs[0]
				Ω(cfg.Space).Should(BeEquivalentTo("space1"))
				Ω(cfg.Developer.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Ω(cfg.Developer.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Ω(cfg.Developer.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))

				Ω(cfg.Auditor.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Ω(cfg.Auditor.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Ω(cfg.Auditor.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))

				Ω(cfg.Manager.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Ω(cfg.Manager.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Ω(cfg.Manager.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))
			})

			It("should return a list of 2", func() {
				m := config.NewManager("./fixtures/config")
				configs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(configs).Should(HaveLen(2))
			})

			It("should return configs for user info", func() {
				m := config.NewManager("./fixtures/user_config")
				configs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(configs).Should(HaveLen(1))
			})

			It("should return configs for user info", func() {
				m := config.NewManager("./fixtures/user_config_multiple_groups")
				configs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(configs).Should(HaveLen(1))
				config := configs[0]
				Ω(config.GetDeveloperGroups()).Should(ConsistOf([]string{"test_space1_developers"}))
				Ω(config.GetAuditorGroups()).Should(ConsistOf([]string{"test_space1_auditors"}))
				Ω(config.GetManagerGroups()).Should(ConsistOf([]string{"test_space1_managers", "test_space1_managers_2"}))
			})

			Context("failure cases", func() {
				It("should return an error when no security.json file is provided", func() {
					m := config.NewManager("./fixtures/no-security-json")
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})

				It("should return an error when malformed yaml", func() {
					m := config.NewManager("./fixtures/bad-yml")
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})

				It("should return an error when path does not exist", func() {
					m := config.NewManager("./fixtures/blah")
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})
			})

		})
	})
})
