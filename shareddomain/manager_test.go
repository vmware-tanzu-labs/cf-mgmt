package shareddomain_test

import (
	"errors"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	"code.cloudfoundry.org/routing-api/models"
	. "github.com/vmwarepivotallabs/cf-mgmt/shareddomain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	fakeconfig "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/shareddomain/fakes"
)

var _ = Describe("Manager", func() {
	var (
		manager           *Manager
		fakeDomainClient  *fakes.FakeCFDomainClient
		fakeJobClient     *fakes.FakeCFJobClient
		fakeRoutingClient *fakes.FakeRoutingClient
		fakeCfg           *fakeconfig.FakeReader
	)
	BeforeEach(func() {
		fakeDomainClient = &fakes.FakeCFDomainClient{}
		fakeJobClient = &fakes.FakeCFJobClient{}
		fakeRoutingClient = &fakes.FakeRoutingClient{}
		fakeCfg = &fakeconfig.FakeReader{}
		manager = NewManager(fakeDomainClient, fakeJobClient, fakeRoutingClient, fakeCfg, false)
		fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
			SharedDomains: map[string]config.SharedDomain{
				"foo.bar":        config.SharedDomain{},
				"default.domain": config.SharedDomain{},
			},
		}, nil)
	})
	Context("Apply", func() {
		It("Should create 2 shared domains", func() {
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(2))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(0))
			for i := 0; i <= 1; i++ {
				_, domainCreate := fakeDomainClient.CreateArgsForCall(i)
				Expect(domainCreate.Name).NotTo(BeEmpty())
				Expect(*domainCreate.Internal).To(BeFalse())
				Expect(domainCreate.RouterGroup.GUID).To(BeEmpty())
			}
		})
		It("Should create 1 shared domain with routing group guid", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				SharedDomains: map[string]config.SharedDomain{
					"foo.bar": config.SharedDomain{
						Internal:    false,
						RouterGroup: "default-tcp",
					},
				},
			}, nil)
			fakeRoutingClient.RouterGroupWithNameReturns(models.RouterGroup{
				Guid: "default-tcp-guid",
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(1))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(0))
			_, domainCreate := fakeDomainClient.CreateArgsForCall(0)
			Expect(domainCreate.Name).NotTo(BeEmpty())
			Expect(*domainCreate.Internal).To(BeFalse())
			Expect(domainCreate.RouterGroup.GUID).NotTo(BeEmpty())

		})
		It("Should create no shared domains", func() {
			fakeDomainClient.ListAllReturns([]*resource.Domain{
				{
					Name: "foo.bar",
					GUID: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					GUID: "default.domain.guid",
				},
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(0))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(0))
		})

		It("Should delete 2 shared domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeDomainClient.ListAllReturns([]*resource.Domain{
				{
					Name: "foo.bar",
					GUID: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					GUID: "default.domain.guid",
				},
			}, nil)
			fakeDomainClient.DeleteReturns("job-guid", nil)
			fakeJobClient.PollCompleteReturns(nil)

			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(0))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(2))
			for i := 0; i <= 1; i++ {
				_, domainGUID := fakeDomainClient.DeleteArgsForCall(i)
				Expect(domainGUID).To(ContainSubstring("guid"))
			}
			_, jobGUID, _ := fakeJobClient.PollCompleteArgsForCall(0)
			Expect(jobGUID).To(Equal("job-guid"))
		})
	})

	Context("errors", func() {
		It("should error on getting config", func() {
			fakeCfg.GetGlobalConfigReturns(nil, errors.New("error getting config"))
			err := manager.Apply()
			Expect(err).To(MatchError("error getting config"))
		})

		It("should error listing domains", func() {
			fakeDomainClient.ListAllReturns(nil, errors.New("error getting shared domains"))
			err := manager.Apply()
			Expect(err).To(MatchError("error getting shared domains"))
		})
		It("should error creating domains", func() {
			fakeDomainClient.CreateReturns(nil, errors.New("error creating shared domain"))
			err := manager.Apply()
			Expect(err).To(MatchError("error creating shared domain"))
		})
		It("should error deleting domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeDomainClient.ListAllReturns([]*resource.Domain{
				{
					Name: "foo.bar",
					GUID: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					GUID: "default.domain.guid",
				},
			}, nil)
			fakeDomainClient.DeleteReturns("", errors.New("error deleting shared domain"))
			err := manager.Apply()
			Expect(err).To(MatchError("error deleting shared domain"))
		})
	})

	Context("peek", func() {
		BeforeEach(func() {
			manager = NewManager(fakeDomainClient, fakeJobClient, fakeRoutingClient, fakeCfg, true)
		})
		It("Should not create 2 shared domains", func() {
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(0))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(0))
		})
		It("Should not delete 2 shared domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeDomainClient.ListAllReturns([]*resource.Domain{
				{
					Name: "foo.bar",
					GUID: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					GUID: "default.domain.guid",
				},
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeDomainClient.CreateCallCount()).To(Equal(0))
			Expect(fakeDomainClient.DeleteCallCount()).To(Equal(0))
		})
	})
})
