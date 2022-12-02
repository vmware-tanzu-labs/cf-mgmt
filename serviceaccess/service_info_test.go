package serviceaccess_test

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess"
	"github.com/vmwarepivotallabs/cf-mgmt/serviceaccess/fakes"
)

var _ = Describe("ServiceInfo", func() {
	Context("StandardBrokers", func() {
		fakeServicePlanClient := new(fakes.FakeCFServicePlanClient)
		fakeServicePlanVisibilityClient := new(fakes.FakeCFServicePlanVisibilityClient)
		fakeServiceOfferingClient := new(fakes.FakeCFServiceOfferingClient)
		fakeServiceBrokerClient := new(fakes.FakeCFServiceBrokerClient)

		It("returns standard brokers", func() {
			standardBrokers := []*resource.ServiceBroker{
				{
					GUID: "some-guid",
					Name: "some-name",
				},
				{
					GUID: "some-guid-2",
					Name: "some-name-2",
				},
				{
					GUID: "some-guid-3",
					Name: "some-name-3",
				},
			}
			fakeServiceBrokerClient.ListAllReturns(standardBrokers, nil)

			serviceInfo, _ := serviceaccess.GetServiceInfo(
				fakeServicePlanClient, fakeServicePlanVisibilityClient, fakeServiceOfferingClient, fakeServiceBrokerClient)
			Expect(serviceInfo.StandardBrokers()).Should(HaveLen(3))
		})

		It("does not return space scoped brokers", func() {
			brokers := []*resource.ServiceBroker{
				{
					GUID: "some-guid",
					Name: "some-name",
				},
				{
					GUID: "some-guid-2",
					Name: "some-name-2",
				},
				{
					GUID: "some-guid-3",
					Name: "some-space-broker-name",
					Relationships: resource.SpaceRelationship{
						Space: resource.ToOneRelationship{
							Data: &resource.Relationship{
								GUID: "non-empty-guid",
							},
						},
					},
				},
			}
			fakeServiceBrokerClient.ListAllReturns(brokers, nil)

			serviceInfo, _ := serviceaccess.GetServiceInfo(
				fakeServicePlanClient, fakeServicePlanVisibilityClient, fakeServiceOfferingClient, fakeServiceBrokerClient)
			Expect(serviceInfo.StandardBrokers()).Should(HaveLen(2))
		})
	})
})
