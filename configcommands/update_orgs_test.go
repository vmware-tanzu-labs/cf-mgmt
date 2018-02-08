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
		configuration UpdateOrgsConfigurationCommand
	)
	BeforeEach(func() {
		mockConfig = new(configfakes.FakeManager)
		configuration = UpdateOrgsConfigurationCommand{
			ConfigManager: mockConfig,
		}
	})
	Context("Update Orgs yaml file", func() {
		It("should keep values the same", func() {
			mockConfig.OrgsReturns(&config.Orgs{
				EnableDeleteOrgs: true,
				Orgs:             []string{"foo", "bar"},
				ProtectedOrgs:    []string{"system", "my-special-org"},
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgsCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgsArgsForCall(0)).To(BeEquivalentTo(&config.Orgs{
				EnableDeleteOrgs: true,
				Orgs:             []string{"foo", "bar"},
				ProtectedOrgs:    []string{"system", "my-special-org"},
			}))
		})

		It("should change enable-delete-orgs from true to false", func() {
			configuration.EnableDeleteOrgs = "false"
			mockConfig.OrgsReturns(&config.Orgs{
				EnableDeleteOrgs: true,
				Orgs:             []string{"foo", "bar"},
				ProtectedOrgs:    []string{"system", "my-special-org"},
			}, nil)
			err := configuration.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(mockConfig.SaveOrgsCallCount()).To(Equal(1))
			Expect(mockConfig.SaveOrgsArgsForCall(0)).To(BeEquivalentTo(&config.Orgs{
				EnableDeleteOrgs: false,
				Orgs:             []string{"foo", "bar"},
				ProtectedOrgs:    []string{"system", "my-special-org"},
			}))
		})

	})

	Context("Failures", func() {
		It("should fail retrieving orgs", func() {
			mockConfig.OrgsReturns(nil, errors.New("error retrieve"))
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(BeEquivalentTo("error retrieve"))
		})
		It("should fail savings orgs", func() {
			mockConfig.OrgsReturns(&config.Orgs{}, nil)
			mockConfig.SaveOrgsReturns(errors.New("error saving"))
			err := configuration.Execute(nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(BeEquivalentTo("error saving"))
		})
	})
})
