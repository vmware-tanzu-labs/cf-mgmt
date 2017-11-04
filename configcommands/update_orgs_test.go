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
