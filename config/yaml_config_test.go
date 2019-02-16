package config_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
)

var _ = Describe("CF-Mgmt Config", func() {
	Context("Protected Org Defaults", func() {
		Describe("Defaults", func() {
			It("should setup default protected orgs", func() {
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("system"))
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("p-spring-cloud-services"))
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("splunk-nozzle-org"))
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("redis-test-ORG*"))
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("appdynamics-org"))
				Ω(config.DefaultProtectedOrgs).Should(ContainElement("credhub-service-broker-org"))
				Ω(config.DefaultProtectedOrgs).Should(HaveLen(6))
			})
		})
	})

	Context("Default Config Reader", func() {
		Context("GetASGConfigs", func() {
			It("should return a single ASG", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfgs).Should(HaveLen(2))

				cfg := cfgs[0]
				Expect(cfg.Rules).Should(BeEquivalentTo("[{\"protocol\": \"icmp\",\"destination\": \"0.0.0.0/0\"}]\n"))
			})

			It("should have a name based on the ASG filename", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cfgs).Should(HaveLen(2))

				namedList := make([]string, len(cfgs))
				for i, asg := range cfgs {
					namedList[i] = asg.Name
				}
				Expect(namedList).Should(ConsistOf("dns", "test-asg"))

			})

			It("can optionally have a ASG name in the spaced config.", func() {
				m := config.NewManager("./fixtures/asg-defaults")

				// Get space config
				cfgs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())

				cfg := cfgs[0]
				Ω(cfg.Space).Should(BeEquivalentTo("space1"))
				Expect(cfg.ASGs).Should(ConsistOf("test-asg"))

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

		Context("GetOrgConfig", func() {
			It("should return a org", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfig("test")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(c).ShouldNot(BeNil())
			})

			It("should return an error", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfig("foo")
				Ω(err).Should(HaveOccurred())
				Ω(c).Should(BeNil())
				Ω(err.Error()).Should(BeEquivalentTo("Org [foo] not found in config"))
			})
		})

		Context("SaveOrgConfig", func() {
			var tempDir string
			var err error
			var configManager config.Manager
			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "cf-mgmt")
				Ω(err).ShouldNot(HaveOccurred())
				configManager = config.NewManager(tempDir)
			})
			AfterEach(func() {
				os.RemoveAll(tempDir)
			})
			It("should succeed", func() {
				orgName := "foo"
				orgConfig := &config.OrgConfig{
					Org: orgName,
				}
				saveError := configManager.SaveOrgConfig(orgConfig)
				Ω(saveError).ShouldNot(HaveOccurred())
				retrieveConfig, retrieveError := configManager.GetOrgConfig(orgName)
				Ω(retrieveError).ShouldNot(HaveOccurred())
				Ω(retrieveConfig).ShouldNot(BeNil())
			})
		})

		Context("DeleteOrgConfig", func() {
			var tempDir string
			var err error
			var configManager config.Manager
			orgName := "foo"
			orgConfig := &config.OrgConfig{
				Org: orgName,
			}
			spaces := &config.Spaces{Org: orgName}
			BeforeEach(func() {
				tempDir, err = ioutil.TempDir("", "cf-mgmt")
				Ω(err).ShouldNot(HaveOccurred())
				configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
				configManager.CreateConfigIfNotExists("ldap")
				addError := configManager.AddOrgToConfig(orgConfig, spaces)
				Ω(addError).ShouldNot(HaveOccurred())
				addError = configManager.AddOrgToConfig(&config.OrgConfig{
					Org: "sdfasdfdf",
				}, &config.Spaces{Org: "sdfasdfdf"})
				Ω(addError).ShouldNot(HaveOccurred())
			})
			AfterEach(func() {
				os.RemoveAll(tempDir)
			})
			It("should succeed", func() {
				err := configManager.DeleteOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				orgs, err := configManager.Orgs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgs.Orgs).ShouldNot(ConsistOf(orgName))
				_, err = configManager.GetOrgConfig(orgName)
				Ω(err).Should(HaveOccurred())
				Ω(err.Error()).Should(BeEquivalentTo("Org [foo] not found in config"))
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

			Context("GetSpaceConfig", func() {
				It("should return a space", func() {
					m := config.NewManager("./fixtures/config")
					c, err := m.GetSpaceConfig("test", "space1")
					Ω(err).ShouldNot(HaveOccurred())
					Ω(c).ShouldNot(BeNil())
				})

				It("should return an error", func() {
					m := config.NewManager("./fixtures/config")
					c, err := m.GetSpaceConfig("test", "foo")
					Ω(err).Should(HaveOccurred())
					Ω(c).Should(BeNil())
					Ω(err.Error()).Should(BeEquivalentTo("Space [foo] not found in org [test] config"))
				})
			})

			Context("SaveSpaceConfig", func() {
				var tempDir string
				var err error
				var configManager config.Manager
				BeforeEach(func() {
					tempDir, err = ioutil.TempDir("", "cf-mgmt")
					Ω(err).ShouldNot(HaveOccurred())
					configManager = config.NewManager(tempDir)
				})
				AfterEach(func() {
					os.RemoveAll(tempDir)
				})
				It("should succeed", func() {
					orgName := "foo"
					spaceName := "bar"
					spaceConfig := &config.SpaceConfig{
						Org:   orgName,
						Space: spaceName,
					}
					saveError := configManager.SaveSpaceConfig(spaceConfig)
					Ω(saveError).ShouldNot(HaveOccurred())
					retrieveConfig, retrieveError := configManager.GetSpaceConfig(orgName, spaceName)
					Ω(retrieveError).ShouldNot(HaveOccurred())
					Ω(retrieveConfig).ShouldNot(BeNil())
				})
			})

			Context("DeleteSpaceConfig", func() {
				var tempDir string
				var err error
				var configManager config.Manager
				orgName := "foo"
				spaceName := "bar"
				orgConfig := &config.OrgConfig{
					Org: orgName,
				}
				spaces := &config.Spaces{
					Org:                orgName,
					EnableDeleteSpaces: true,
				}
				spaceConfig := &config.SpaceConfig{
					Org:   orgName,
					Space: spaceName,
				}
				BeforeEach(func() {
					tempDir, err = ioutil.TempDir("", "cf-mgmt")
					Ω(err).ShouldNot(HaveOccurred())
					configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
					configManager.CreateConfigIfNotExists("ldap")
					addError := configManager.AddOrgToConfig(orgConfig, spaces)
					Ω(addError).ShouldNot(HaveOccurred())
					addError = configManager.AddSpaceToConfig(spaceConfig)
					Ω(addError).ShouldNot(HaveOccurred())
					addError = configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   orgName,
						Space: "asdfsadfs",
					})
					Ω(addError).ShouldNot(HaveOccurred())
				})
				AfterEach(func() {
					os.RemoveAll(tempDir)
				})
				It("should fail to find space", func() {
					err := configManager.DeleteSpaceConfig(orgName, spaceName)
					Ω(err).ShouldNot(HaveOccurred())
					spaces, err := configManager.OrgSpaces(orgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(spaces.Spaces).ShouldNot(ConsistOf(spaceName))
					_, err = configManager.GetSpaceConfig(orgName, spaceName)
					Ω(err).Should(HaveOccurred())
					Ω(err.Error()).Should(Equal("Space [bar] not found in org [foo] config"))
				})
			})

			Context("AddOrgConfig", func() {
				var tempDir string
				var err error
				var configManager config.Manager
				BeforeEach(func() {
					tempDir, err = ioutil.TempDir("", "cf-mgmt")
					Ω(err).ShouldNot(HaveOccurred())
					configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
					configManager.CreateConfigIfNotExists("ldap")
				})
				AfterEach(func() {
					os.RemoveAll(tempDir)
				})
				It("should succeed adding an org that doesn't exist", func() {
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					}, &config.Spaces{
						Org:                "foo",
						EnableDeleteSpaces: true,
					})
					Ω(err).ShouldNot(HaveOccurred())
					orgs, err := configManager.Orgs()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(orgs.Orgs).Should(ConsistOf("foo"))
					_, err = configManager.GetOrgConfig("foo")
					Ω(err).Should(Not(HaveOccurred()))
				})
				It("should fail adding an org with different case", func() {
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					}, &config.Spaces{Org: "foo"})
					Ω(err).ShouldNot(HaveOccurred())
					err = configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "Foo",
					}, &config.Spaces{Org: "Foo"})
					Ω(err).Should(HaveOccurred())
					orgs, err := configManager.Orgs()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(len(orgs.Orgs)).Should(BeEquivalentTo(1))
				})
			})

			Context("AddSpaceConfig", func() {
				var tempDir string
				var err error
				var configManager config.Manager
				BeforeEach(func() {
					tempDir, err = ioutil.TempDir("", "cf-mgmt")
					Ω(err).ShouldNot(HaveOccurred())
					configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
					configManager.CreateConfigIfNotExists("ldap")
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					}, &config.Spaces{Org: "foo"})
					Ω(err).ShouldNot(HaveOccurred())
				})
				AfterEach(func() {
					os.RemoveAll(tempDir)
				})
				It("should succeed adding an space that doesn't exist", func() {
					err := configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "bar",
					})
					Ω(err).ShouldNot(HaveOccurred())
					spaces, err := configManager.GetSpaceConfigs()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(spaces[0].Space).Should(BeEquivalentTo("bar"))
				})
				It("should fail adding an space with different case", func() {
					err := configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "bar",
					})
					Ω(err).ShouldNot(HaveOccurred())
					err = configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "Bar",
					})
					Ω(err).Should(HaveOccurred())
					spaces, err := configManager.GetSpaceConfigs()
					Ω(err).ShouldNot(HaveOccurred())
					Ω(spaces[0].Space).Should(BeEquivalentTo("bar"))
				})
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
