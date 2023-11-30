package serviceaccess_test

import (
	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/fakes"
)

var _ = Describe("ServiceInfo", func() {
	Context("StandardBrokers", func() {
		fakeClient := new(fakes.FakeCFClient)

		It("returns standard brokers", func() {
			standardBrokers := []cfclient.ServiceBroker{
				{
					Guid: "some-guid",
					Name: "some-name",
				},
				{
					Guid: "some-guid-2",
					Name: "some-name-2",
				},
				{
					Guid: "some-guid-3",
					Name: "some-name-3",
				},
			}
			var error error
			fakeClient.ListServiceBrokersReturns(standardBrokers, error)

			serviceInfo, _ := serviceaccess.GetServiceInfo(fakeClient)
			Expect(serviceInfo.StandardBrokers()).Should(HaveLen(3))
		})

		It("does not return space scoped brokers", func() {
			brokers := []cfclient.ServiceBroker{
				{
					Guid: "some-guid",
					Name: "some-name",
				},
				{
					Guid: "some-guid-2",
					Name: "some-name-2",
				},
				{
					Guid:      "some-guid-3",
					Name:      "some-space-broker-name",
					SpaceGUID: "non-empty-guid",
				},
			}
			var error error
			fakeClient.ListServiceBrokersReturns(brokers, error)

			serviceInfo, _ := serviceaccess.GetServiceInfo(fakeClient)
			Expect(serviceInfo.StandardBrokers()).Should(HaveLen(2))
		})
	})
})
