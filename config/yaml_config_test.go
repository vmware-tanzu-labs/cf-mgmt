package config_test

import (
	"fmt"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

var _ = Describe("CF-Mgmt Config", func() {
	Context("Protected Org Defaults", func() {
		Describe("Defaults", func() {
			It("should setup default protected orgs", func() {
				Expect(config.DefaultProtectedOrgs).Should(HaveLen(6))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("^system$"))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("splunk-nozzle-org"))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("redis-test-ORG"))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("appdynamics-org"))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("credhub-service-broker-org"))
				Expect(config.DefaultProtectedOrgs).Should(ContainElement("^p-"))
			})
		})
	})

	Context("Default Config Reader", func() {
		Context("Creating Configuration", func() {
			var (
				configManager config.Manager
				pwd, _        = os.Getwd()
				configDir     = path.Join(pwd, "_testGen")
			)
			BeforeEach(func() {
				configManager = config.NewManager(configDir)
			})

			AfterEach(func() {
				err := os.RemoveAll(configDir)
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("Should initialize the configuration", func() {
				err := configManager.CreateConfigIfNotExists("ldap")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(path.Join(configDir, "asgs")).To(BeADirectory())
				Expect(path.Join(configDir, "asgs", ".gitkeep")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "default_asgs")).To(BeADirectory())
				Expect(path.Join(configDir, "default_asgs", ".gitkeep")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "org_quotas")).To(BeADirectory())
				Expect(path.Join(configDir, "org_quotas", ".gitkeep")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "ldap.yml")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "cf-mgmt.yml")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "orgs.yml")).To(BeAnExistingFile())
				Expect(path.Join(configDir, "spaceDefaults.yml")).To(BeAnExistingFile())
			})
		})
		Context("GetASGConfigs", func() {
			It("should return a single ASG", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(cfgs).Should(HaveLen(2))

				cfg := cfgs[0]
				Expect(cfg.Rules).Should(BeEquivalentTo("[{\"protocol\": \"icmp\",\"destination\": \"0.0.0.0/0\"}]\n"))
			})

			It("should have a name based on the ASG filename", func() {
				m := config.NewManager("./fixtures/asg-defaults")
				cfgs, err := m.GetASGConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(cfgs).Should(HaveLen(2))

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
				Expect(err).ShouldNot(HaveOccurred())

				cfg := cfgs[0]
				Expect(cfg.Space).Should(BeEquivalentTo("space1"))
				Expect(cfg.ASGs).Should(ConsistOf("test-asg"))

			})

		})

		Context("GetOrgConfigs", func() {
			It("should return a list of 2", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(c).Should(HaveLen(2))
			})

			It("should return a list of 1", func() {
				m := config.NewManager("./fixtures/duplicate-files")
				c, err := m.GetOrgConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(c).Should(HaveLen(1))
			})

			It("should return a list of 1", func() {
				m := config.NewManager("./fixtures/user_config")
				c, err := m.GetOrgConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(c).Should(HaveLen(1))

				org := c[0]
				Expect(org.GetAuditorGroups()).Should(ConsistOf([]string{"test_org_auditors"}))
				Expect(org.GetManagerGroups()).Should(ConsistOf([]string{"test_org_managers"}))
				Expect(org.GetBillingManagerGroups()).Should(ConsistOf([]string{"test_billing_managers", "test_billing_managers_2"}))
			})

			It("should fail when given an invalid config dir", func() {
				m := config.NewManager("./fixtures/blah")
				c, err := m.GetOrgConfigs()
				Expect(err).Should(HaveOccurred())
				Expect(c).Should(BeEmpty())
			})
		})

		Context("GetOrgConfig", func() {
			It("should return a org", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfig("test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(c).ShouldNot(BeNil())
			})

			It("should return an error", func() {
				m := config.NewManager("./fixtures/config")
				c, err := m.GetOrgConfig("foo")
				Expect(err).Should(HaveOccurred())
				Expect(c).Should(BeNil())
				Expect(err.Error()).Should(BeEquivalentTo("Org [foo] not found in config"))
			})
		})

		Context("SaveOrgConfig", func() {
			var (
				err           error
				configManager config.Manager
				pwd, _        = os.Getwd()
				configDir     = path.Join(pwd, "_testGen")
			)
			BeforeEach(func() {
				configManager = config.NewManager(configDir)
				err = configManager.CreateConfigIfNotExists("uaa")
				Expect(err).ShouldNot(HaveOccurred())
			})
			AfterEach(func() {
				os.RemoveAll(configDir)
			})
			It("should succeed", func() {
				orgName := "foo"
				orgConfig := &config.OrgConfig{
					Org: orgName,
				}
				saveError := configManager.SaveOrgConfig(orgConfig)
				Expect(saveError).ShouldNot(HaveOccurred())
				retrieveConfig, retrieveError := configManager.GetOrgConfig(orgName)
				Expect(retrieveError).ShouldNot(HaveOccurred())
				Expect(retrieveConfig).ShouldNot(BeNil())
				orgs, err := configManager.Orgs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orgs.Orgs).Should(ConsistOf("foo"))
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
				tempDir, err = os.MkdirTemp("", "cf-mgmt")
				Expect(err).ShouldNot(HaveOccurred())
				configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
				configManager.CreateConfigIfNotExists("ldap")
				addError := configManager.AddOrgToConfig(orgConfig)
				Expect(addError).ShouldNot(HaveOccurred())
				addError = configManager.AddOrgToConfig(&config.OrgConfig{
					Org: "sdfasdfdf",
				})
				Expect(addError).ShouldNot(HaveOccurred())
				err := configManager.SaveOrgSpaces(spaces)
				Expect(err).ShouldNot(HaveOccurred())
			})
			AfterEach(func() {
				os.RemoveAll(tempDir)
			})
			It("should succeed", func() {
				err := configManager.DeleteOrgConfig(orgName)
				Expect(err).ShouldNot(HaveOccurred())
				orgs, err := configManager.Orgs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orgs.Orgs).ShouldNot(ConsistOf(orgName))
				_, err = configManager.GetOrgConfig(orgName)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(BeEquivalentTo("Org [foo] not found in config"))
			})
		})

		Context("GetSpaceConfigs", func() {
			It("should return a single space", func() {
				m := config.NewManager("./fixtures/space-defaults")
				cfgs, err := m.GetSpaceConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(cfgs).Should(HaveLen(1))

				cfg := cfgs[0]
				Expect(cfg.Space).Should(BeEquivalentTo("space1"))
				Expect(cfg.Developer.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Expect(cfg.Developer.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Expect(cfg.Developer.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))

				Expect(cfg.Auditor.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Expect(cfg.Auditor.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Expect(cfg.Auditor.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))

				Expect(cfg.Manager.LDAPUsers).Should(ConsistOf("default-ldap-user", "space1-ldap-user"))
				Expect(cfg.Manager.Users).Should(ConsistOf("default-user@test.com", "space-1-user@test.com"))
				Expect(cfg.Manager.LDAPGroup).Should(BeEquivalentTo("space-1-ldap-group"))
			})

			It("should return a list of 2", func() {
				m := config.NewManager("./fixtures/config")
				configs, err := m.GetSpaceConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(configs).Should(HaveLen(2))
			})

			It("should return configs for user info", func() {
				m := config.NewManager("./fixtures/user_config")
				configs, err := m.GetSpaceConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(configs).Should(HaveLen(1))
			})

			It("should return configs for user info", func() {
				m := config.NewManager("./fixtures/user_config_multiple_groups")
				configs, err := m.GetSpaceConfigs()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(configs).Should(HaveLen(1))
				config := configs[0]
				Expect(config.GetDeveloperGroups()).Should(ConsistOf([]string{"test_space1_developers"}))
				Expect(config.GetAuditorGroups()).Should(ConsistOf([]string{"test_space1_auditors"}))
				Expect(config.GetManagerGroups()).Should(ConsistOf([]string{"test_space1_managers", "test_space1_managers_2"}))
			})

			Context("GetSpaceConfig", func() {
				It("should return a space", func() {
					m := config.NewManager("./fixtures/config")
					c, err := m.GetSpaceConfig("test", "space1")
					Expect(err).ShouldNot(HaveOccurred())
					Expect(c).ShouldNot(BeNil())
				})

				It("should return an error", func() {
					m := config.NewManager("./fixtures/config")
					c, err := m.GetSpaceConfig("test", "foo")
					Expect(err).Should(HaveOccurred())
					Expect(c).Should(BeNil())
					Expect(err.Error()).Should(BeEquivalentTo("Space [foo] not found in org [test] config"))
				})
			})

			Context("SaveSpaceConfig", func() {
				var (
					configManager config.Manager
					pwd, _        = os.Getwd()
					configDir     = path.Join(pwd, "_testGen")
				)
				BeforeEach(func() {
					configManager = config.NewManager(configDir)
					err := configManager.CreateConfigIfNotExists("uaa")
					Expect(err).ShouldNot(HaveOccurred())
					err = configManager.SaveOrgConfig(&config.OrgConfig{
						Org: "foo",
					})
					Expect(err).ShouldNot(HaveOccurred())

					err = configManager.SaveOrgSpaces(&config.Spaces{Org: "foo"})
					Expect(err).ShouldNot(HaveOccurred())
				})
				AfterEach(func() {
					os.RemoveAll(configDir)
				})
				It("should succeed", func() {
					orgName := "foo"
					spaceName := "bar"
					spaceConfig := &config.SpaceConfig{
						Org:   orgName,
						Space: spaceName,
					}
					saveError := configManager.SaveSpaceConfig(spaceConfig)
					Expect(saveError).ShouldNot(HaveOccurred())
					retrieveConfig, retrieveError := configManager.GetSpaceConfig(orgName, spaceName)
					Expect(retrieveError).ShouldNot(HaveOccurred())
					Expect(retrieveConfig).ShouldNot(BeNil())
				})
			})

			Context("DeleteSpaceConfig", func() {
				var (
					configManager config.Manager
					pwd, _        = os.Getwd()
					configDir     = path.Join(pwd, "_testGen")
				)
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
					configManager = config.NewManager(configDir)
					configManager.CreateConfigIfNotExists("ldap")
					addError := configManager.AddOrgToConfig(orgConfig)
					Expect(addError).ShouldNot(HaveOccurred())
					err := configManager.SaveOrgSpaces(spaces)
					Expect(err).ShouldNot(HaveOccurred())
					addError = configManager.AddSpaceToConfig(spaceConfig)
					Expect(addError).ShouldNot(HaveOccurred())
					addError = configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   orgName,
						Space: "asdfsadfs",
					})
					Expect(addError).ShouldNot(HaveOccurred())
				})
				AfterEach(func() {
					os.RemoveAll(configDir)
				})
				It("should fail to find space", func() {
					err := configManager.DeleteSpaceConfig(orgName, spaceName)
					Expect(err).ShouldNot(HaveOccurred())
					spaces, err := configManager.OrgSpaces(orgName)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spaces.Spaces).ShouldNot(ConsistOf(spaceName))
					_, err = configManager.GetSpaceConfig(orgName, spaceName)
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).Should(Equal("Space [bar] not found in org [foo] config"))
				})
			})

			Context("AddOrgConfig", func() {
				var tempDir string
				var err error
				var configManager config.Manager
				BeforeEach(func() {
					tempDir, err = os.MkdirTemp("", "cf-mgmt")
					Expect(err).ShouldNot(HaveOccurred())
					configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
					configManager.CreateConfigIfNotExists("ldap")
				})
				AfterEach(func() {
					os.RemoveAll(tempDir)
				})
				It("should succeed adding an org that doesn't exist", func() {
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					})
					Expect(err).ShouldNot(HaveOccurred())
					orgs, err := configManager.Orgs()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(orgs.Orgs).Should(ConsistOf("foo"))
					_, err = configManager.GetOrgConfig("foo")
					Expect(err).Should(Not(HaveOccurred()))
				})
				It("should fail adding an org with different case", func() {
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					})
					Expect(err).ShouldNot(HaveOccurred())
					err = configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "Foo",
					})
					Expect(err).Should(HaveOccurred())
					orgs, err := configManager.Orgs()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(orgs.Orgs)).Should(BeEquivalentTo(1))
				})
			})

			Context("AddSpaceConfig", func() {
				var (
					configManager config.Manager
					pwd, _        = os.Getwd()
					configDir     = path.Join(pwd, "_testGen")
				)
				BeforeEach(func() {
					configManager = config.NewManager(configDir)
					configManager.CreateConfigIfNotExists("ldap")
					err := configManager.AddOrgToConfig(&config.OrgConfig{
						Org: "foo",
					})
					Expect(err).ShouldNot(HaveOccurred())
					err = configManager.SaveOrgSpaces(&config.Spaces{Org: "foo"})
					Expect(err).ShouldNot(HaveOccurred())
				})
				AfterEach(func() {
					os.RemoveAll(configDir)
				})
				It("should succeed adding an space that doesn't exist", func() {
					err := configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "bar",
					})
					Expect(err).ShouldNot(HaveOccurred())
					spaces, err := configManager.GetSpaceConfigs()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spaces[0].Space).Should(BeEquivalentTo("bar"))
				})
				It("should fail adding an space with different case", func() {
					err := configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "bar",
					})
					Expect(err).ShouldNot(HaveOccurred())
					err = configManager.AddSpaceToConfig(&config.SpaceConfig{
						Org:   "foo",
						Space: "Bar",
					})
					Expect(err).Should(HaveOccurred())
					spaces, err := configManager.GetSpaceConfigs()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(spaces[0].Space).Should(BeEquivalentTo("bar"))
				})
			})

			Context("failure cases", func() {
				It("should return an error when no security.json file is provided", func() {
					m := config.NewManager("./fixtures/no-security-json")
					configs, err := m.GetSpaceConfigs()
					Expect(err).Should(HaveOccurred())
					Expect(configs).Should(BeNil())
				})

				It("should return an error when malformed yaml", func() {
					m := config.NewManager("./fixtures/bad-yml")
					configs, err := m.GetSpaceConfigs()
					Expect(err).Should(HaveOccurred())
					Expect(configs).Should(BeNil())
				})

				It("should return an error when path does not exist", func() {
					m := config.NewManager("./fixtures/blah")
					configs, err := m.GetSpaceConfigs()
					Expect(err).Should(HaveOccurred())
					Expect(configs).Should(BeNil())
				})
			})

		})
	})

	Context("YAML Config Updater", func() {
		Context("User Actions", func() {
			var tempDir string
			var configManager config.Manager

			BeforeEach(func() {
				var err error
				tempDir, err = os.MkdirTemp("", "cf-mgmt")
				Expect(err).ShouldNot(HaveOccurred())
				configManager = config.NewManager(path.Join(tempDir, "cfmgmt"))
				configManager.CreateConfigIfNotExists("ldap")
			})

			AfterEach(func() {
				os.RemoveAll(tempDir)
			})

			Context("Associate Org Auditor", func() {
				When("the org exists", func() {
					const (
						orgName  = "the-org"
						userName = "the-user"
					)
					BeforeEach(func() {
						err := configManager.AddOrgToConfig(&config.OrgConfig{
							Org: orgName,
						})
						Expect(err).ShouldNot(HaveOccurred())
					})

					When("the user does not exist", func() {
						When("the internal origin is requested", func() {
							It("adds the user to the org", func() {
								o, err := configManager.GetOrgConfig(orgName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(o.Auditor.Users).Should(HaveLen(0))

								err = configManager.AssociateOrgAuditor(config.InternalOrigin, orgName, userName)
								Expect(err).ShouldNot(HaveOccurred())

								o, err = configManager.GetOrgConfig(orgName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(o.Auditor.Users).Should(HaveLen(1))
							})
						})

						When("the saml origin is requested", func() {
							It("adds the user to the org", func() {
								o, err := configManager.GetOrgConfig(orgName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(o.Auditor.SamlUsers).Should(HaveLen(0))

								err = configManager.AssociateOrgAuditor(config.SAMLOrigin, orgName, userName)
								Expect(err).ShouldNot(HaveOccurred())

								o, err = configManager.GetOrgConfig(orgName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(o.Auditor.SamlUsers).Should(HaveLen(1))
							})
						})
					})

					When("the user already exists", func() {
						BeforeEach(func() {
							o, err := configManager.GetOrgConfig(orgName)
							Expect(err).ShouldNot(HaveOccurred())

							o.Auditor.Users = append(o.Auditor.Users, userName)

							err = configManager.SaveOrgConfig(o)
							Expect(err).ShouldNot(HaveOccurred())
						})

						It("does nothing and returns nil", func() {
							o, err := configManager.GetOrgConfig(orgName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(o.Auditor.Users).Should(HaveLen(1))

							err = configManager.AssociateOrgAuditor(config.InternalOrigin, orgName, userName)
							Expect(err).ShouldNot(HaveOccurred())

							o, err = configManager.GetOrgConfig(orgName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(o.Auditor.Users).Should(HaveLen(1))
						})
					})
				})

				When("the org does not exist", func() {
					It("returns an error", func() {
						orgName := "org-that-does-not-exist"
						err := configManager.AssociateOrgAuditor(config.InternalOrigin, orgName, "my-user")
						Expect(err).Should(HaveOccurred())
						Expect(err).Should(MatchError(fmt.Sprintf("Org [%s] not found in config", orgName)))
					})
				})
			})

			Context("Associate Space Role", func() {
				When("the org and space exist", func() {
					const (
						orgName   = "the-org"
						spaceName = "the-space"
						userName  = "the-user"
					)
					BeforeEach(func() {
						err := configManager.AddOrgToConfig(&config.OrgConfig{
							Org: orgName,
						})
						Expect(err).ShouldNot(HaveOccurred())

						err = configManager.SaveOrgSpaces(&config.Spaces{
							Org: orgName,
						})
						Expect(err).ShouldNot(HaveOccurred())

						err = configManager.AddSpaceToConfig(&config.SpaceConfig{
							Org:   orgName,
							Space: spaceName,
						})
						Expect(err).ShouldNot(HaveOccurred())
					})

					When("the user does not exist in the space", func() {
						It("creates the user with a developer role", func() {
							s, err := configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Developer.Users).Should(HaveLen(0))

							err = configManager.AssociateSpaceDeveloper(config.InternalOrigin, orgName, spaceName, userName)
							Expect(err).ShouldNot(HaveOccurred())

							s, err = configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Developer.Users).Should(HaveLen(1))
						})

						It("creates the user with an auditor role", func() {
							s, err := configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Auditor.Users).Should(HaveLen(0))

							err = configManager.AssociateSpaceAuditor(config.InternalOrigin, orgName, spaceName, userName)
							Expect(err).ShouldNot(HaveOccurred())

							s, err = configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Auditor.Users).Should(HaveLen(1))
						})

						It("creates the user in the correct origin", func() {
							s, err := configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Auditor.SamlUsers).Should(HaveLen(0))

							err = configManager.AssociateSpaceAuditor(config.SAMLOrigin, orgName, spaceName, userName)
							Expect(err).ShouldNot(HaveOccurred())

							s, err = configManager.GetSpaceConfig(orgName, spaceName)
							Expect(err).ShouldNot(HaveOccurred())
							Expect(s.Auditor.SamlUsers).Should(HaveLen(1))
							Expect(s.Auditor.Users).Should(HaveLen(0))
						})
					})

					When("user already exists in space", func() {
						When("saml origin is requested", func() {
							BeforeEach(func() {
								s, err := configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())

								s.Developer.SamlUsers = append(s.Developer.SamlUsers, userName)

								err = configManager.SaveSpaceConfig(s)
								Expect(err).ShouldNot(HaveOccurred())
							})

							It("does nothing and returns nil", func() {
								s, err := configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(s.Developer.SamlUsers).Should(HaveLen(1))

								err = configManager.AssociateSpaceDeveloper(config.SAMLOrigin, orgName, spaceName, userName)
								Expect(err).ShouldNot(HaveOccurred())

								s, err = configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(s.Developer.SamlUsers).Should(HaveLen(1))
							})
						})

						When("internal origin is requested", func() {
							BeforeEach(func() {
								s, err := configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())

								s.Developer.Users = append(s.Developer.Users, userName)

								err = configManager.SaveSpaceConfig(s)
								Expect(err).ShouldNot(HaveOccurred())
							})

							It("does nothing and returns nil", func() {
								s, err := configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(s.Developer.Users).Should(HaveLen(1))

								err = configManager.AssociateSpaceDeveloper(config.InternalOrigin, orgName, spaceName, userName)
								Expect(err).ShouldNot(HaveOccurred())

								s, err = configManager.GetSpaceConfig(orgName, spaceName)
								Expect(err).ShouldNot(HaveOccurred())
								Expect(s.Developer.Users).Should(HaveLen(1))
							})
						})
					})
				})

				When("the space cannot be found", func() {
					It("returns an error", func() {
						const (
							orgName   = "org-that-maybe-exists"
							spaceName = "space-that-does-not-exist"
						)
						err := configManager.AssociateSpaceDeveloper(config.InternalOrigin, orgName, spaceName, "my-user")
						Expect(err).Should(HaveOccurred())
						Expect(err).Should(MatchError(fmt.Sprintf("Space [%s] not found in org [%s] config", spaceName, orgName)))
					})
				})
			})
		})
	})
})
