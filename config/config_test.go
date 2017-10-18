package config_test

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	. "github.com/pivotalservices/cf-mgmt/config/test_data"
	kOrg "github.com/pivotalservices/cf-mgmt/organization/constants"
	kSpace "github.com/pivotalservices/cf-mgmt/space/constants"
	mock "github.com/pivotalservices/cf-mgmt/utils/mocks"
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
		Context("GetOrgConfigs", func() {
			var utilsMgrMock *mock.MockUtilsManager
			BeforeEach(func() {
				utilsMgrMock = mock.NewMockUtilsManager()
				PopulateWithTestData(utilsMgrMock)
			})
			It("should return a list of 2", func() {
				m := config.NewManager("./fixtures/config", utilsMgrMock)
				c, err := m.GetOrgConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(c).Should(HaveLen(2))
			})

			It("should return a list of 1", func() {
				m := config.NewManager("./fixtures/user_config", utilsMgrMock)
				c, err := m.GetOrgConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(c).Should(HaveLen(1))

				org := c[0]
				Ω(org.GetAuditorGroups()).Should(ConsistOf([]string{"test_org_auditors"}))
				Ω(org.GetManagerGroups()).Should(ConsistOf([]string{"test_org_managers"}))
				Ω(org.GetBillingManagerGroups()).Should(ConsistOf([]string{"test_billing_managers", "test_billing_managers_2"}))
			})

			It("should fail when given an invalid config dir", func() {
				m := config.NewManager("./fixtures/blah", utilsMgrMock)
				c, err := m.GetOrgConfigs()
				Ω(err).Should(HaveOccurred())
				Ω(c).Should(BeEmpty())
			})
		})

		Context("GetSpaceConfigs", func() {
			var utilsMgrMock *mock.MockUtilsManager
			BeforeEach(func() {
				utilsMgrMock = mock.NewMockUtilsManager()
				PopulateWithTestData(utilsMgrMock)
			})
			It("should return a single space", func() {
				m := config.NewManager("./fixtures/space-defaults", utilsMgrMock)
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
				m := config.NewManager("./fixtures/config", utilsMgrMock)
				configs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(configs).Should(HaveLen(2))
			})

			It("should return configs for user info", func() {
				utilsMgrMock.MockFileData = map[string]interface{}{}
				PopulateWithTestData(utilsMgrMock)
				m := config.NewManager("./fixtures/user_config", utilsMgrMock)
				configs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(configs).Should(HaveLen(1))
			})

			It("should return configs for user info", func() {
				m := config.NewManager("./fixtures/user_config_multiple_groups", utilsMgrMock)
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
					m := config.NewManager("./fixtures/no-security-json", utilsMgrMock)
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})

				It("should return an error when malformed yaml", func() {
					utilsMgrMock.MockFileDataHasError = true
					m := config.NewManager("./fixtures/bad-yml", utilsMgrMock)
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})

				It("should return an error when path does not exist", func() {
					m := config.NewManager("./fixtures/blah", utilsMgrMock)
					configs, err := m.GetSpaceConfigs()
					Ω(err).Should(HaveOccurred())
					Ω(configs).Should(BeNil())
				})
			})

		})
	})

	Context("Adding Users", func() {
		Context("AddUserToSpaceConfig", func() {
			var utilsMgrMock *mock.MockUtilsManager
			var configDir string
			var randomUserName string
			var orgName string
			var spaceName string
			BeforeEach(func() {
				utilsMgrMock = mock.NewMockUtilsManager()
				PopulateWithTestData(utilsMgrMock)
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)

				firstName := make([]byte, 5)
				lastName := make([]byte, 5)

				r1.Read(firstName)
				r1.Read(lastName)

				randomUserName = fmt.Sprintf("%X.%X", firstName, lastName)
				configDir = "./fixtures/user_update"
				orgName = "test"
				spaceName = "space1"
			})

			It("should be able to insert an LDAP space developer", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_DEVELOPERS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, LDAPUser := range spaceConfig.Developer.LDAPUsers {
							if LDAPUser == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert an LDAP space auditor", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_AUDITORS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, LDAPUser := range spaceConfig.Auditor.LDAPUsers {
							if LDAPUser == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert an LDAP space manager", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_MANAGERS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, LDAPUser := range spaceConfig.Manager.LDAPUsers {
							if LDAPUser == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account space developer", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_DEVELOPERS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, User := range spaceConfig.Developer.Users {
							if User == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account space auditor", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_AUDITORS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, User := range spaceConfig.Auditor.Users {
							if User == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account space manager", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToSpaceConfig(randomUserName, kSpace.ROLE_SPACE_MANAGERS, spaceName, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the space config and check that our randomUserName exists in the target role
				spaceConfigs, err := m.GetSpaceConfigs()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(spaceConfigs).ShouldNot(BeNil())

				foundUserName := false
				for _, spaceConfig := range spaceConfigs {
					if spaceConfig.Org == orgName && spaceConfig.Space == spaceName {
						for _, User := range spaceConfig.Manager.Users {
							if User == randomUserName {
								foundUserName = true
								break
							}
						}
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
		})

		Context("AddUserToOrgConfig", func() {
			var utilsMgrMock *mock.MockUtilsManager
			var randomUserName string
			var configDir string
			var orgName string
			BeforeEach(func() {
				utilsMgrMock = mock.NewMockUtilsManager()
				PopulateWithTestData(utilsMgrMock)
				s1 := rand.NewSource(time.Now().UnixNano())
				r1 := rand.New(s1)

				firstName := make([]byte, 5)
				lastName := make([]byte, 5)

				r1.Read(firstName)
				r1.Read(lastName)

				randomUserName = fmt.Sprintf("%X.%X", firstName, lastName)
				configDir = "./fixtures/user_update"
				orgName = "test"
			})

			It("should be able to insert an LDAP Org Manager", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_MANAGERS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, LDAPUser := range orgConfig.Manager.LDAPUsers {
					if LDAPUser == randomUserName {
						foundUserName = true
						break
					}
				}

				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert an LDAP Org Auditor", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_AUDITORS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, LDAPUser := range orgConfig.Auditor.LDAPUsers {
					if LDAPUser == randomUserName {
						foundUserName = true
						break
					}
				}

				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert an LDAP Org Billing Manager", func() {
				isLdapUser := true
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_BILLING_MANAGERS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, LDAPUser := range orgConfig.BillingManager.LDAPUsers {
					if LDAPUser == randomUserName {
						foundUserName = true
						break
					}
				}

				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account Org Manager", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_MANAGERS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				m = config.NewManager(configDir, utilsMgrMock)
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, User := range orgConfig.Manager.Users {
					if User == randomUserName {
						foundUserName = true
						break
					}
				}

				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account Org auditor", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_AUDITORS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				m = config.NewManager(configDir, utilsMgrMock)
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, User := range orgConfig.Auditor.Users {
					if User == randomUserName {
						foundUserName = true
						break
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
			It("should be able to insert a service account Org billing manager", func() {
				isLdapUser := false
				m := config.NewManager(configDir, utilsMgrMock)
				err := m.AddUserToOrgConfig(randomUserName, kOrg.ROLE_ORG_BILLING_MANAGERS, orgName, isLdapUser)
				Ω(err).ShouldNot(HaveOccurred())

				// Get the org config and check that our randomUserName exists in the target role
				orgConfig, err := m.GetAnOrgConfig(orgName)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(orgConfig).ShouldNot(BeNil())

				foundUserName := false
				for _, User := range orgConfig.BillingManager.Users {
					if User == randomUserName {
						foundUserName = true
						break
					}
				}
				Ω(foundUserName).Should(BeTrue())
			})
		})

		Context("Update Org Quota", func() {
			var utilsMgrMock *mock.MockUtilsManager
			Context("UpdateQuotasInOrgConfig", func() {
				var targetOrgName string
				var configDir string
				BeforeEach(func() {
					utilsMgrMock = mock.NewMockUtilsManager()
					PopulateWithTestData(utilsMgrMock)
					targetOrgName = "test"
					configDir = "./fixtures/config"
				})
				It("should be able to update org quota with with the enable org quotas off", func() {
					// Some initial setup and assumptions
					enableOrgQuota := false
					m := config.NewManager(configDir, utilsMgrMock)

					// Get a copy of the original Org Config
					originalOrgConfig, err := m.GetAnOrgConfig(targetOrgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(originalOrgConfig).ShouldNot(BeNil())

					newQuotaSettings := map[string]string{
						"MemoryLimit":             strconv.FormatInt(int64(originalOrgConfig.MemoryLimit+0*10*1024), 10),
						"InstanceMemoryLimit":     strconv.FormatInt(int64(originalOrgConfig.InstanceMemoryLimit+0*1024), 10),
						"TotalRoutes":             strconv.FormatInt(int64(originalOrgConfig.TotalRoutes+0*1000), 10),
						"TotalServices":           strconv.FormatInt(int64(originalOrgConfig.TotalServices+0*1000), 10),
						"PaidServicePlansAllowed": strconv.FormatBool(!originalOrgConfig.PaidServicePlansAllowed),
						"TotalPrivateDomains":     strconv.FormatInt(int64(originalOrgConfig.TotalPrivateDomains+0*1000), 10),
						"TotalReservedRoutePorts": strconv.FormatInt(int64(originalOrgConfig.TotalReservedRoutePorts+0*1000), 10),
						"TotalServiceKeys":        strconv.FormatInt(int64(originalOrgConfig.TotalServiceKeys+0*1000), 10),
						"AppInstanceLimit":        strconv.FormatInt(int64(originalOrgConfig.AppInstanceLimit+0*1000), 10),
					}

					err = m.UpdateQuotasInOrgConfig(targetOrgName, enableOrgQuota, newQuotaSettings)
					Ω(err).ShouldNot(HaveOccurred())

					// Check the values to see if they match
					newOrgConfig, err := m.GetAnOrgConfig(targetOrgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(newOrgConfig).ShouldNot(BeNil())

					Ω(newOrgConfig.EnableOrgQuota).Should(Equal(false))
					Ω(strconv.FormatInt(int64(newOrgConfig.MemoryLimit), 10)).Should(Equal(newQuotaSettings["MemoryLimit"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.InstanceMemoryLimit), 10)).Should(Equal(newQuotaSettings["InstanceMemoryLimit"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalRoutes), 10)).Should(Equal(newQuotaSettings["TotalRoutes"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalServices), 10)).Should(Equal(newQuotaSettings["TotalServices"]))
					Ω(strconv.FormatBool(newOrgConfig.PaidServicePlansAllowed)).Should(Equal(newQuotaSettings["PaidServicePlansAllowed"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalPrivateDomains), 10)).Should(Equal(newQuotaSettings["TotalPrivateDomains"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalReservedRoutePorts), 10)).Should(Equal(newQuotaSettings["TotalReservedRoutePorts"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalServiceKeys), 10)).Should(Equal(newQuotaSettings["TotalServiceKeys"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.AppInstanceLimit), 10)).Should(Equal(newQuotaSettings["AppInstanceLimit"]))
				})

				It("should be able to update org quota with with the enable org quotas on", func() {
					// Some initial setup and assumptions
					enableOrgQuota := true
					m := config.NewManager(configDir, utilsMgrMock)

					// Get a copy of the original Org Config
					originalOrgConfig, err := m.GetAnOrgConfig(targetOrgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(originalOrgConfig).ShouldNot(BeNil())

					newQuotaSettings := map[string]string{
						"MemoryLimit":             strconv.FormatInt(int64(originalOrgConfig.MemoryLimit+0*10*1024), 10),
						"InstanceMemoryLimit":     strconv.FormatInt(int64(originalOrgConfig.InstanceMemoryLimit+0*1024), 10),
						"TotalRoutes":             strconv.FormatInt(int64(originalOrgConfig.TotalRoutes+0*1000), 10),
						"TotalServices":           strconv.FormatInt(int64(originalOrgConfig.TotalServices+0*1000), 10),
						"PaidServicePlansAllowed": strconv.FormatBool(!originalOrgConfig.PaidServicePlansAllowed),
						"TotalPrivateDomains":     strconv.FormatInt(int64(originalOrgConfig.TotalPrivateDomains+0*1000), 10),
						"TotalReservedRoutePorts": strconv.FormatInt(int64(originalOrgConfig.TotalReservedRoutePorts+0*1000), 10),
						"TotalServiceKeys":        strconv.FormatInt(int64(originalOrgConfig.TotalServiceKeys+0*1000), 10),
						"AppInstanceLimit":        strconv.FormatInt(int64(originalOrgConfig.AppInstanceLimit+0*1000), 10),
					}

					err = m.UpdateQuotasInOrgConfig(targetOrgName, enableOrgQuota, newQuotaSettings)
					Ω(err).ShouldNot(HaveOccurred())

					// Check the values to see if they match
					newOrgConfig, err := m.GetAnOrgConfig(targetOrgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(newOrgConfig).ShouldNot(BeNil())

					Ω(newOrgConfig.EnableOrgQuota).Should(Equal(true))
					Ω(strconv.FormatInt(int64(newOrgConfig.MemoryLimit), 10)).Should(Equal(newQuotaSettings["MemoryLimit"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.InstanceMemoryLimit), 10)).Should(Equal(newQuotaSettings["InstanceMemoryLimit"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalRoutes), 10)).Should(Equal(newQuotaSettings["TotalRoutes"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalServices), 10)).Should(Equal(newQuotaSettings["TotalServices"]))
					Ω(strconv.FormatBool(newOrgConfig.PaidServicePlansAllowed)).Should(Equal(newQuotaSettings["PaidServicePlansAllowed"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalPrivateDomains), 10)).Should(Equal(newQuotaSettings["TotalPrivateDomains"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalReservedRoutePorts), 10)).Should(Equal(newQuotaSettings["TotalReservedRoutePorts"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.TotalServiceKeys), 10)).Should(Equal(newQuotaSettings["TotalServiceKeys"]))
					Ω(strconv.FormatInt(int64(newOrgConfig.AppInstanceLimit), 10)).Should(Equal(newQuotaSettings["AppInstanceLimit"]))
				})
			})
		})
		Context("Update Space Quota", func() {
			var utilsMgrMock *mock.MockUtilsManager
			Context("UpdateQuotasInSpaceConfig", func() {
				var targetOrgName string
				var targetSpaceName string
				var configDir string
				BeforeEach(func() {
					utilsMgrMock = mock.NewMockUtilsManager()
					PopulateWithTestData(utilsMgrMock)
					targetOrgName = "test"
					targetSpaceName = "space1"
					configDir = "./fixtures/config"
				})
				It("should be able to update space quota with with the enable space quotas off", func() {
					// Some initial setup and assumptions
					enableSpaceQuota := false
					m := config.NewManager(configDir, utilsMgrMock)

					// Get a copy of the original Space Config
					loadSpaceDefaults := false
					originalSpaceConfig, err := m.GetASpaceConfig(targetOrgName, targetSpaceName, loadSpaceDefaults)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(originalSpaceConfig).ShouldNot(BeNil())

					newQuotaSettings := map[string]string{
						"MemoryLimit":             strconv.FormatInt(int64(originalSpaceConfig.MemoryLimit+0*10*1024), 10),
						"InstanceMemoryLimit":     strconv.FormatInt(int64(originalSpaceConfig.InstanceMemoryLimit+0*1024), 10),
						"TotalRoutes":             strconv.FormatInt(int64(originalSpaceConfig.TotalRoutes+0*1000), 10),
						"TotalServices":           strconv.FormatInt(int64(originalSpaceConfig.TotalServices+0*1000), 10),
						"PaidServicePlansAllowed": strconv.FormatBool(!originalSpaceConfig.PaidServicePlansAllowed),
						"TotalPrivateDomains":     strconv.FormatInt(int64(originalSpaceConfig.TotalPrivateDomains+0*1000), 10),
						"TotalReservedRoutePorts": strconv.FormatInt(int64(originalSpaceConfig.TotalReservedRoutePorts+0*1000), 10),
						"TotalServiceKeys":        strconv.FormatInt(int64(originalSpaceConfig.TotalServiceKeys+0*1000), 10),
						"AppInstanceLimit":        strconv.FormatInt(int64(originalSpaceConfig.AppInstanceLimit+0*1000), 10),
					}

					err = m.UpdateQuotasInSpaceConfig(targetOrgName, targetSpaceName, enableSpaceQuota, newQuotaSettings)
					Ω(err).ShouldNot(HaveOccurred())

					// Check the values to see if they match
					newSpaceConfig, err := m.GetASpaceConfig(targetOrgName, targetSpaceName, loadSpaceDefaults)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(newSpaceConfig).ShouldNot(BeNil())

					Ω(newSpaceConfig.EnableSpaceQuota).Should(Equal(false))
					Ω(strconv.FormatInt(int64(newSpaceConfig.MemoryLimit), 10)).Should(Equal(newQuotaSettings["MemoryLimit"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.InstanceMemoryLimit), 10)).Should(Equal(newQuotaSettings["InstanceMemoryLimit"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalRoutes), 10)).Should(Equal(newQuotaSettings["TotalRoutes"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalServices), 10)).Should(Equal(newQuotaSettings["TotalServices"]))
					Ω(strconv.FormatBool(newSpaceConfig.PaidServicePlansAllowed)).Should(Equal(newQuotaSettings["PaidServicePlansAllowed"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalPrivateDomains), 10)).Should(Equal(newQuotaSettings["TotalPrivateDomains"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalReservedRoutePorts), 10)).Should(Equal(newQuotaSettings["TotalReservedRoutePorts"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalServiceKeys), 10)).Should(Equal(newQuotaSettings["TotalServiceKeys"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.AppInstanceLimit), 10)).Should(Equal(newQuotaSettings["AppInstanceLimit"]))
				})

				It("should be able to update space quota with with the enable space quotas on", func() {
					// Some initial setup and assumptions
					enableSpaceQuota := true
					m := config.NewManager(configDir, utilsMgrMock)

					// Get a copy of the original Space Config
					loadSpaceDefaults := false
					originalSpaceConfig, err := m.GetASpaceConfig(targetOrgName, targetSpaceName, loadSpaceDefaults)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(originalSpaceConfig).ShouldNot(BeNil())

					newQuotaSettings := map[string]string{
						"MemoryLimit":             strconv.FormatInt(int64(originalSpaceConfig.MemoryLimit+0*10*1024), 10),
						"InstanceMemoryLimit":     strconv.FormatInt(int64(originalSpaceConfig.InstanceMemoryLimit+0*1024), 10),
						"TotalRoutes":             strconv.FormatInt(int64(originalSpaceConfig.TotalRoutes+0*1000), 10),
						"TotalServices":           strconv.FormatInt(int64(originalSpaceConfig.TotalServices+0*1000), 10),
						"PaidServicePlansAllowed": strconv.FormatBool(!originalSpaceConfig.PaidServicePlansAllowed),
						"TotalPrivateDomains":     strconv.FormatInt(int64(originalSpaceConfig.TotalPrivateDomains+0*1000), 10),
						"TotalReservedRoutePorts": strconv.FormatInt(int64(originalSpaceConfig.TotalReservedRoutePorts+0*1000), 10),
						"TotalServiceKeys":        strconv.FormatInt(int64(originalSpaceConfig.TotalServiceKeys+0*1000), 10),
						"AppInstanceLimit":        strconv.FormatInt(int64(originalSpaceConfig.AppInstanceLimit+0*1000), 10),
					}

					err = m.UpdateQuotasInSpaceConfig(targetOrgName, targetSpaceName, enableSpaceQuota, newQuotaSettings)
					Ω(err).ShouldNot(HaveOccurred())

					// Check the values to see if they match
					newSpaceConfig, err := m.GetASpaceConfig(targetOrgName, targetSpaceName, loadSpaceDefaults)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(newSpaceConfig).ShouldNot(BeNil())

					Ω(newSpaceConfig.EnableSpaceQuota).Should(Equal(true))
					Ω(strconv.FormatInt(int64(newSpaceConfig.MemoryLimit), 10)).Should(Equal(newQuotaSettings["MemoryLimit"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.InstanceMemoryLimit), 10)).Should(Equal(newQuotaSettings["InstanceMemoryLimit"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalRoutes), 10)).Should(Equal(newQuotaSettings["TotalRoutes"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalServices), 10)).Should(Equal(newQuotaSettings["TotalServices"]))
					Ω(strconv.FormatBool(newSpaceConfig.PaidServicePlansAllowed)).Should(Equal(newQuotaSettings["PaidServicePlansAllowed"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalPrivateDomains), 10)).Should(Equal(newQuotaSettings["TotalPrivateDomains"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalReservedRoutePorts), 10)).Should(Equal(newQuotaSettings["TotalReservedRoutePorts"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.TotalServiceKeys), 10)).Should(Equal(newQuotaSettings["TotalServiceKeys"]))
					Ω(strconv.FormatInt(int64(newSpaceConfig.AppInstanceLimit), 10)).Should(Equal(newQuotaSettings["AppInstanceLimit"]))
				})
			})
		})
		Context("Adding Private Domains", func() {
			Context("AddOrgPrivateDomainToConfig", func() {
				var utilsMgrMock *mock.MockUtilsManager
				var randomDomainName string
				var configDir string
				var orgName string
				BeforeEach(func() {
					utilsMgrMock = mock.NewMockUtilsManager()
					PopulateWithTestData(utilsMgrMock)
					s1 := rand.NewSource(time.Now().UnixNano())
					r1 := rand.New(s1)

					firstName := make([]byte, 5)
					lastName := make([]byte, 5)

					r1.Read(firstName)
					r1.Read(lastName)

					randomDomainName = fmt.Sprintf("%X-%X.com", firstName, lastName)
					configDir = "./fixtures/user_update"
					orgName = "test"
				})

				It("should be able to insert a private domain name", func() {
					m := config.NewManager(configDir, utilsMgrMock)
					err := m.AddPrivateDomainToOrgConfig(orgName, randomDomainName)
					Ω(err).ShouldNot(HaveOccurred())

					orgConfig, err := m.GetAnOrgConfig(orgName)
					Ω(err).ShouldNot(HaveOccurred())
					Ω(orgConfig).ShouldNot(BeNil())

					foundPrivateDomain := false
					for _, privateDomainName := range orgConfig.PrivateDomains {
						if privateDomainName == randomDomainName {
							foundPrivateDomain = true
						}
					}
					Ω(foundPrivateDomain).Should(BeTrue())
				})

			})

		})
	})

})
