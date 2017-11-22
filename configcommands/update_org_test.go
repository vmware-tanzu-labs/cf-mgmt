package configcommands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"

	"github.com/pivotalservices/cf-mgmt/config/configfakes"
	. "github.com/pivotalservices/cf-mgmt/configcommands"
)

var _ = Describe("given update orgs config command", func() {
	var (
		mockConfig    *configfakes.FakeManager
		configuration UpdateOrgConfigurationCommand
	)
	orgName := "foo"
	BeforeEach(func() {
		mockConfig = new(configfakes.FakeManager)
		configuration = UpdateOrgConfigurationCommand{
			OrgName:       orgName,
			ConfigManager: mockConfig,
		}
	})
	Context("Updating basic org config", func() {
		It("should succeed when updating private domains", func() {
			configuration.PrivateDomains = []string{"foo.com", "bar.io"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                     orgName,
				RemovePrivateDomains:    false,
				PrivateDomains:          []string{"foo.com", "bar.io"},
				EnableOrgQuota:          false,
				PaidServicePlansAllowed: false,
			}))
		})

		It("should succeed when deleting private domains", func() {
			configuration.PrivateDomainsToRemove = []string{"foo.com"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org:            orgName,
				PrivateDomains: []string{"foo.com", "bar.io"},
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                     orgName,
				RemovePrivateDomains:    false,
				PrivateDomains:          []string{"bar.io"},
				EnableOrgQuota:          false,
				PaidServicePlansAllowed: false,
			}))
		})

		It("should enable remove of private domains", func() {
			configuration.EnableRemovePrivateDomains = "true"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org:                  orgName,
				RemovePrivateDomains: false,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                  orgName,
				RemovePrivateDomains: true,
			}))
		})

		It("should disable remove of private domains", func() {
			configuration.EnableRemovePrivateDomains = "false"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org:                  orgName,
				RemovePrivateDomains: true,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                  orgName,
				RemovePrivateDomains: false,
			}))
		})
		It("should fail when enable is not a valid boolean", func() {
			configuration.EnableRemovePrivateDomains = "asdfasf"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org:                  orgName,
				RemovePrivateDomains: true,
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("--enable-remove-private-domains must be an boolean instead of [asdfasf]"))
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(0))
		})

		It("should succeed when updating shared private domains", func() {
			configuration.SharedPrivateDomains = []string{"foo.com", "bar.io"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                     orgName,
				RemovePrivateDomains:    false,
				SharedPrivateDomains:    []string{"foo.com", "bar.io"},
				EnableOrgQuota:          false,
				PaidServicePlansAllowed: false,
			}))
		})

		It("should succeed when deleting shared private domains", func() {
			configuration.SharedPrivateDomainsToRemove = []string{"foo.com"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org:                  orgName,
				SharedPrivateDomains: []string{"foo.com", "bar.io"},
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                     orgName,
				RemovePrivateDomains:    false,
				SharedPrivateDomains:    []string{"bar.io"},
				EnableOrgQuota:          false,
				PaidServicePlansAllowed: false,
			}))
		})

		It("should enable remove of shared private domains", func() {
			configuration.EnableRemoveSharedPrivateDomains = "true"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				RemoveSharedPrivateDomains: false,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				RemoveSharedPrivateDomains: true,
			}))
		})

		It("should disable remove of private domains", func() {
			configuration.EnableRemoveSharedPrivateDomains = "false"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				RemoveSharedPrivateDomains: true,
			}, nil)
			mockConfig.SaveOrgConfigReturns(nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				RemoveSharedPrivateDomains: false,
			}))
		})
		It("should fail when enable is not a valid boolean", func() {
			configuration.EnableRemoveSharedPrivateDomains = "asdfasf"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				RemoveSharedPrivateDomains: true,
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("enable-remove-shared-private-domains must be an boolean instead of [asdfasf]"))
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(0))
		})
	})
	Context("Updating Quotas", func() {
		It("should succeed", func() {
			configuration.Quota.EnableOrgQuota = "true"
			configuration.Quota.MemoryLimit = "1"
			configuration.Quota.InstanceMemoryLimit = "2"
			configuration.Quota.TotalRoutes = "3"
			configuration.Quota.TotalServices = "4"
			configuration.Quota.PaidServicesAllowed = "true"
			configuration.Quota.TotalPrivateDomains = "5"
			configuration.Quota.TotalReservedRoutePorts = "6"
			configuration.Quota.TotalServiceKeys = "7"
			configuration.Quota.AppInstanceLimit = "8"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org:                     orgName,
				RemovePrivateDomains:    false,
				EnableOrgQuota:          true,
				MemoryLimit:             1,
				InstanceMemoryLimit:     2,
				TotalRoutes:             3,
				TotalServices:           4,
				PaidServicePlansAllowed: true,
				TotalPrivateDomains:     5,
				TotalReservedRoutePorts: 6,
				TotalServiceKeys:        7,
				AppInstanceLimit:        8,
			}))
		})

		It("should fail with non integer value", func() {
			configuration.Quota.EnableOrgQuota = "true"
			configuration.Quota.MemoryLimit = "asdfasfasf"
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("--memory-limit must be an integer instead of [asdfasfasf]"))
		})

	})

	Context("Update Users", func() {
		It("should add users to empty list", func() {
			configuration.Manager.Users = []string{"foo", "bar"}
			configuration.BillingManager.Users = []string{"hello", "world"}
			configuration.Auditor.Users = []string{"test", "value"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					Users: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test", "value"},
				},
			}))
		})

		It("should not add users that already exist", func() {
			configuration.Manager.Users = []string{"bar"}
			configuration.BillingManager.Users = []string{"world"}
			configuration.Auditor.Users = []string{"value"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					Users: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("--value [world] already exists in [hello world]"))
			Expect(err.Error()).Should(ContainSubstring("--value [value] already exists in [test value]"))
			Expect(err.Error()).Should(ContainSubstring("--value [bar] already exists in [foo bar]"))
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(0))
		})

		It("should not duplicates", func() {
			configuration.Manager.Users = []string{"bar", "bar", "foo"}
			configuration.BillingManager.Users = []string{"world", "world", "hello"}
			configuration.Auditor.Users = []string{"value", "value", "test"}
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("value [bar] cannot be specified more than once [bar bar foo]"))
			Expect(err.Error()).Should(ContainSubstring("value [world] cannot be specified more than once [world world hello]"))
			Expect(err.Error()).Should(ContainSubstring("value [value] cannot be specified more than once [value value test]"))
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(0))
		})
		It("should remove users from existing", func() {
			configuration.Manager.UsersToRemove = []string{"bar"}
			configuration.BillingManager.UsersToRemove = []string{"world"}
			configuration.Auditor.UsersToRemove = []string{"value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					Users: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					Users: []string{"foo"},
				},
				BillingManager: config.UserMgmt{
					Users: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test"},
				},
			}))
		})

		It("should add saml users to empty list", func() {
			configuration.Manager.SamlUsers = []string{"foo", "bar"}
			configuration.BillingManager.SamlUsers = []string{"hello", "world"}
			configuration.Auditor.SamlUsers = []string{"test", "value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					SamlUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test", "value"},
				},
			}))
		})

		It("should remove saml users from existing", func() {
			configuration.Manager.SamlUsersToRemove = []string{"bar"}
			configuration.BillingManager.SamlUsersToRemove = []string{"world"}
			configuration.Auditor.SamlUsersToRemove = []string{"value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					SamlUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo"},
				},
				BillingManager: config.UserMgmt{
					SamlUsers: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test"},
				},
			}))
		})

		It("should add ldap users to empty list", func() {
			configuration.Manager.LDAPUsers = []string{"foo", "bar"}
			configuration.BillingManager.LDAPUsers = []string{"hello", "world"}
			configuration.Auditor.LDAPUsers = []string{"test", "value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					LDAPUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test", "value"},
				},
			}))
		})

		It("should remove ldap users from existing", func() {
			configuration.Manager.LDAPUsersToRemove = []string{"bar"}
			configuration.BillingManager.LDAPUsersToRemove = []string{"world"}
			configuration.Auditor.LDAPUsersToRemove = []string{"value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					LDAPUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo"},
				},
				BillingManager: config.UserMgmt{
					LDAPUsers: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test"},
				},
			}))
		})

		It("should add ldap groups to empty list", func() {
			configuration.Manager.LDAPGroups = []string{"foo", "bar"}
			configuration.BillingManager.LDAPGroups = []string{"hello", "world"}
			configuration.Auditor.LDAPGroups = []string{"test", "value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					LDAPGroups: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPGroups: []string{"test", "value"},
				},
			}))
		})

		It("should remove ldap groups from existing", func() {
			configuration.Manager.LDAPGroupsToRemove = []string{"bar"}
			configuration.BillingManager.LDAPGroupsToRemove = []string{"world"}
			configuration.Auditor.LDAPGroupsToRemove = []string{"value"}

			mockConfig.GetOrgConfigReturns(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo", "bar"},
				},
				BillingManager: config.UserMgmt{
					LDAPGroups: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPGroups: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveOrgConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgConfigArgsForCall(0)).To(BeEquivalentTo(&config.OrgConfig{
				Org: orgName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo"},
				},
				BillingManager: config.UserMgmt{
					LDAPGroups: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					LDAPGroups: []string{"test"},
				},
			}))
		})
	})
	Context("Failures", func() {
		It("should fail retrieving config", func() {
			mockConfig.GetOrgConfigReturns(nil, errors.New("error retrieve"))
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(BeEquivalentTo("error retrieve"))
		})
		It("should fail saving config", func() {
			mockConfig.GetOrgConfigReturns(&config.OrgConfig{}, nil)
			mockConfig.SaveOrgConfigReturns(errors.New("error save"))

			err := configuration.Execute(nil)
			Expect(err.Error()).Should(BeEquivalentTo("error save"))
		})
	})
})
