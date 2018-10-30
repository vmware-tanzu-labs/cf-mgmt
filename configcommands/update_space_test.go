package configcommands_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/config/fakes"

	. "github.com/pivotalservices/cf-mgmt/configcommands"
)

var _ = Describe("given update orgs config command", func() {
	var (
		mockConfig    *fakes.FakeManager
		configuration UpdateSpaceConfigurationCommand
	)
	orgName := "foo"
	spaceName := "bar"
	BeforeEach(func() {
		mockConfig = new(fakes.FakeManager)
		configuration = UpdateSpaceConfigurationCommand{
			OrgName:       orgName,
			SpaceName:     spaceName,
			ConfigManager: mockConfig,
		}
	})
	Context("Updating basic org config", func() {

	})
	Context("Updating Quotas", func() {
		It("should succeed", func() {
			configuration.Quota.EnableSpaceQuota = "true"
			configuration.Quota.MemoryLimit = "1"
			configuration.Quota.InstanceMemoryLimit = "2"
			configuration.Quota.TotalRoutes = "3"
			configuration.Quota.TotalServices = "4"
			configuration.Quota.PaidServicesAllowed = "true"
			configuration.Quota.TotalPrivateDomains = "5"
			configuration.Quota.TotalReservedRoutePorts = "6"
			configuration.Quota.TotalServiceKeys = "7"
			configuration.Quota.AppInstanceLimit = "8"
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:                     orgName,
				Space:                   spaceName,
				EnableSpaceQuota:        true,
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
			configuration.Quota.EnableSpaceQuota = "true"
			configuration.Quota.MemoryLimit = "asdfasfasf"
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("--memory-limit must be an integer instead of [asdfasfasf]"))
		})

	})

	Context("Update named asgs", func() {
		It("should add named asgs to empty list", func() {
			configuration.ASGs = []string{"hello", "world"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			mockConfig.GetASGConfigsReturns([]config.ASGConfig{
				config.ASGConfig{
					Name: "hello",
				},
				config.ASGConfig{
					Name: "world",
				},
			}, nil)
			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				ASGs:  []string{"hello", "world"},
			}))
		})

		It("should error when asg definition doesn't exist", func() {
			configuration.ASGs = []string{"hello"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			mockConfig.GetASGConfigsReturns([]config.ASGConfig{
				config.ASGConfig{
					Name: "world",
				},
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("[hello.json] does not exist in asgs directory"))
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(0))
		})

		It("should not add asgs that already exist", func() {
			configuration.ASGs = []string{"hello"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				ASGs:  []string{"hello", "world"},
			}, nil)

			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("--value [hello] already exists in [hello world]"))
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(0))
		})

		It("should not duplicates", func() {
			configuration.ASGs = []string{"hello", "hello", "world"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("value [hello] cannot be specified more than once [hello hello world]"))
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(0))
		})
	})
	Context("Update Users", func() {
		It("should add users to empty list", func() {
			configuration.Manager.Users = []string{"foo", "bar"}
			configuration.Developer.Users = []string{"hello", "world"}
			configuration.Auditor.Users = []string{"test", "value"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					Users: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test", "value"},
				},
			}))
		})

		It("should not add users that already exist", func() {
			configuration.Manager.Users = []string{"bar"}
			configuration.Developer.Users = []string{"world"}
			configuration.Auditor.Users = []string{"value"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
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
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(0))
		})

		It("should not duplicates", func() {
			configuration.Manager.Users = []string{"bar", "bar", "foo"}
			configuration.Developer.Users = []string{"world", "world", "hello"}
			configuration.Auditor.Users = []string{"value", "value", "test"}
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)

			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("value [bar] cannot be specified more than once [bar bar foo]"))
			Expect(err.Error()).Should(ContainSubstring("value [world] cannot be specified more than once [world world hello]"))
			Expect(err.Error()).Should(ContainSubstring("value [value] cannot be specified more than once [value value test]"))
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(0))
		})
		It("should remove users from existing", func() {
			configuration.Manager.UsersToRemove = []string{"bar"}
			configuration.Developer.UsersToRemove = []string{"world"}
			configuration.Auditor.UsersToRemove = []string{"value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					Users: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					Users: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					Users: []string{"foo"},
				},
				Developer: config.UserMgmt{
					Users: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					Users: []string{"test"},
				},
			}))
		})

		It("should add saml users to empty list", func() {
			configuration.Manager.SamlUsers = []string{"foo", "bar"}
			configuration.Developer.SamlUsers = []string{"hello", "world"}
			configuration.Auditor.SamlUsers = []string{"test", "value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					SamlUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test", "value"},
				},
			}))
		})

		It("should remove saml users from existing", func() {
			configuration.Manager.SamlUsersToRemove = []string{"bar"}
			configuration.Developer.SamlUsersToRemove = []string{"world"}
			configuration.Auditor.SamlUsersToRemove = []string{"value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					SamlUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					SamlUsers: []string{"foo"},
				},
				Developer: config.UserMgmt{
					SamlUsers: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					SamlUsers: []string{"test"},
				},
			}))
		})

		It("should add ldap users to empty list", func() {
			configuration.Manager.LDAPUsers = []string{"foo", "bar"}
			configuration.Developer.LDAPUsers = []string{"hello", "world"}
			configuration.Auditor.LDAPUsers = []string{"test", "value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					LDAPUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test", "value"},
				},
			}))
		})

		It("should remove ldap users from existing", func() {
			configuration.Manager.LDAPUsersToRemove = []string{"bar"}
			configuration.Developer.LDAPUsersToRemove = []string{"world"}
			configuration.Auditor.LDAPUsersToRemove = []string{"value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					LDAPUsers: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPUsers: []string{"foo"},
				},
				Developer: config.UserMgmt{
					LDAPUsers: []string{"hello"},
				},
				Auditor: config.UserMgmt{
					LDAPUsers: []string{"test"},
				},
			}))
		})

		It("should add ldap groups to empty list", func() {
			configuration.Manager.LDAPGroups = []string{"foo", "bar"}
			configuration.Developer.LDAPGroups = []string{"hello", "world"}
			configuration.Auditor.LDAPGroups = []string{"test", "value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
			}, nil)

			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))

			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					LDAPGroups: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPGroups: []string{"test", "value"},
				},
			}))
		})

		It("should remove ldap groups from existing", func() {
			configuration.Manager.LDAPGroupsToRemove = []string{"bar"}
			configuration.Developer.LDAPGroupsToRemove = []string{"world"}
			configuration.Auditor.LDAPGroupsToRemove = []string{"value"}

			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo", "bar"},
				},
				Developer: config.UserMgmt{
					LDAPGroups: []string{"hello", "world"},
				},
				Auditor: config.UserMgmt{
					LDAPGroups: []string{"test", "value"},
				},
			}, nil)

			err := configuration.Execute(nil)
			Expect(mockConfig.SaveSpaceConfigCallCount()).To(Equal(1))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveSpaceConfigArgsForCall(0)).To(BeEquivalentTo(&config.SpaceConfig{
				Org:   orgName,
				Space: spaceName,
				Manager: config.UserMgmt{
					LDAPGroups: []string{"foo"},
				},
				Developer: config.UserMgmt{
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
			mockConfig.GetSpaceConfigReturns(nil, errors.New("error retrieve"))
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(BeEquivalentTo("error retrieve"))
		})
		It("should fail saving config", func() {
			mockConfig.GetSpaceConfigReturns(&config.SpaceConfig{}, nil)
			mockConfig.SaveSpaceConfigReturns(errors.New("error save"))

			err := configuration.Execute(nil)
			Expect(err.Error()).Should(BeEquivalentTo("error save"))
		})
	})
})
