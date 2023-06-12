package shareddomain_test

import (
	"errors"

	"code.cloudfoundry.org/routing-api/models"
	cfclient "github.com/cloudfoundry-community/go-cfclient"
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
		fakeCFClient      *fakes.FakeCFClient
		fakeRoutingClient *fakes.FakeRoutingClient
		fakeCfg           *fakeconfig.FakeReader
	)
	BeforeEach(func() {
		fakeCFClient = &fakes.FakeCFClient{}
		fakeRoutingClient = &fakes.FakeRoutingClient{}
		fakeCfg = &fakeconfig.FakeReader{}
		manager = NewManager(fakeCFClient, fakeRoutingClient, fakeCfg, false)
		fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
			SharedDomains: map[string]config.SharedDomain{
				"foo.bar":        {},
				"default.domain": {},
			},
		}, nil)
	})
	Context("Apply", func() {
		It("Should create 2 shared domains", func() {
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(2))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(0))
			for i := 0; i <= 1; i++ {
				name, internal, routerGUID := fakeCFClient.CreateSharedDomainArgsForCall(i)
				Expect(name).To(Not(BeEmpty()))
				Expect(internal).To(BeFalse())
				Expect(routerGUID).To(BeEmpty())
			}
		})
		It("Should create 1 shared domain with routing group guid", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				SharedDomains: map[string]config.SharedDomain{
					"foo.bar": {
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
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(1))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(0))
			name, internal, routerGUID := fakeCFClient.CreateSharedDomainArgsForCall(0)
			Expect(name).To(Not(BeEmpty()))
			Expect(internal).To(BeFalse())
			Expect(routerGUID).To(Not(BeEmpty()))

		})
		It("Should create no shared domains", func() {
			fakeCFClient.ListSharedDomainsReturns([]cfclient.SharedDomain{
				{
					Name: "foo.bar",
					Guid: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					Guid: "default.domain.guid",
				},
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(0))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(0))
		})

		It("Should delete 2 shared domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeCFClient.ListSharedDomainsReturns([]cfclient.SharedDomain{
				{
					Name: "foo.bar",
					Guid: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					Guid: "default.domain.guid",
				},
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(0))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(2))
			for i := 0; i <= 1; i++ {
				guid, async := fakeCFClient.DeleteSharedDomainArgsForCall(i)
				Expect(guid).To(ContainSubstring("guid"))
				Expect(async).To(BeFalse())
			}
		})
	})

	Context("errors", func() {
		It("should error on getting config", func() {
			fakeCfg.GetGlobalConfigReturns(nil, errors.New("error getting config"))
			err := manager.Apply()
			Expect(err).To(MatchError("error getting config"))
		})

		It("should error listing domains", func() {
			fakeCFClient.ListSharedDomainsReturns(nil, errors.New("error getting shared domains"))
			err := manager.Apply()
			Expect(err).To(MatchError("error getting shared domains"))
		})
		It("should error creating domains", func() {
			fakeCFClient.CreateSharedDomainReturns(nil, errors.New("error creating shared domain"))
			err := manager.Apply()
			Expect(err).To(MatchError("error creating shared domain"))
		})
		It("should error deleting domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeCFClient.ListSharedDomainsReturns([]cfclient.SharedDomain{
				{
					Name: "foo.bar",
					Guid: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					Guid: "default.domain.guid",
				},
			}, nil)
			fakeCFClient.DeleteSharedDomainReturns(errors.New("error deleting shared domain"))
			err := manager.Apply()
			Expect(err).To(MatchError("error deleting shared domain"))
		})
	})

	Context("peek", func() {
		BeforeEach(func() {
			manager = NewManager(fakeCFClient, fakeRoutingClient, fakeCfg, true)
		})
		It("Should not create 2 shared domains", func() {
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(0))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(0))
		})
		It("Should not delete 2 shared domains", func() {
			fakeCfg.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableDeleteSharedDomains: true,
			}, nil)
			fakeCFClient.ListSharedDomainsReturns([]cfclient.SharedDomain{
				{
					Name: "foo.bar",
					Guid: "foo.bar.guid",
				},
				{
					Name: "default.domain",
					Guid: "default.domain.guid",
				},
			}, nil)
			err := manager.Apply()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeCFClient.CreateSharedDomainCallCount()).To(Equal(0))
			Expect(fakeCFClient.DeleteSharedDomainCallCount()).To(Equal(0))
		})
	})
})
