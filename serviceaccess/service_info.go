package serviceaccess

import (
	"context"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
)

// GetServiceInfo - returns broker, it's services and their plans
func GetServiceInfo(
	servicePlanClient CFServicePlanClient,
	servicePlanVisibilityClient CFServicePlanVisibilityClient,
	serviceOfferingClient CFServiceOfferingClient,
	serviceBrokerClient CFServiceBrokerClient) (*ServiceInfo, error) {

	serviceInfo := &ServiceInfo{}
	brokers, err := serviceBrokerClient.ListAll(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	for _, broker := range brokers {
		spaceGUID := ""
		if broker.Relationships.Space.Data != nil {
			spaceGUID = broker.Relationships.Space.Data.GUID
		}
		serviceBroker := &ServiceBroker{
			Name:      broker.Name,
			SpaceGUID: spaceGUID,
		}

		serviceOfferingOpts := client.NewServiceOfferingListOptions()
		serviceOfferingOpts.ServiceBrokerGUIDs.EqualTo(broker.GUID)
		serviceOfferings, err := serviceOfferingClient.ListAll(context.Background(), serviceOfferingOpts)
		if err != nil {
			return nil, err
		}
		for _, serviceOffering := range serviceOfferings {
			service := &Service{
				Name: serviceOffering.Name,
			}

			servicePlanOpts := client.NewServicePlanListOptions()
			servicePlanOpts.ServiceOfferingGUIDs.EqualTo(serviceOffering.GUID)
			servicePlans, err := servicePlanClient.ListAll(context.Background(), servicePlanOpts)
			if err != nil {
				return nil, err
			}
			for _, plan := range servicePlans {
				servicePlan := &ServicePlanInfo{
					Name:        plan.Name,
					GUID:        plan.GUID,
					ServiceName: service.Name,
					Public:      plan.VisibilityType == resource.ServicePlanVisibilityPublic.String(),
				}
				planVisibility, err := servicePlanVisibilityClient.Get(context.Background(), servicePlan.GUID)
				if err != nil {
					return nil, err
				}
				for _, org := range planVisibility.Organizations {
					orgVisibility := &Visibility{
						OrgGUID:         org.GUID,
						ServicePlanGUID: servicePlan.GUID,
					}
					servicePlan.AddOrg(orgVisibility)
				}

				service.AddPlan(servicePlan)
			}

			serviceBroker.AddService(service)
		}

		serviceInfo.AddBroker(serviceBroker)
	}

	return serviceInfo, nil
}

type ServiceInfo struct {
	brokers map[string]*ServiceBroker
}

func (s *ServiceInfo) StandardBrokers() []*ServiceBroker {
	brokerList := []*ServiceBroker{}
	for _, broker := range s.brokers {
		if !isSpaceScopedServiceBroker(broker) {
			brokerList = append(brokerList, broker)
		}
	}
	return brokerList
}

func isSpaceScopedServiceBroker(broker *ServiceBroker) bool {
	return broker.SpaceGUID != ""
}

func (s *ServiceInfo) GetBroker(brokerName string) (*ServiceBroker, error) {
	if broker, ok := s.brokers[brokerName]; ok {
		return broker, nil
	}
	return nil, fmt.Errorf("Broker %s is not found", brokerName)
}
func (s *ServiceInfo) AddBroker(serviceBroker *ServiceBroker) {
	if s.brokers == nil {
		s.brokers = make(map[string]*ServiceBroker)
	}
	s.brokers[serviceBroker.Name] = serviceBroker
}

func (s *ServiceInfo) GetServicePlans(brokerName, serviceName string, plans []string) ([]*ServicePlanInfo, error) {
	servicePlans := []*ServicePlanInfo{}
	broker, err := s.GetBroker(brokerName)
	if err != nil {
		return nil, err
	}
	service, err := broker.GetService(serviceName)
	if err != nil {
		return nil, err
	}
	for _, plan := range plans {
		servicePlan, err := service.GetPlan(plan)
		if err != nil {
			return nil, err
		}
		servicePlans = append(servicePlans, servicePlan)
	}

	return servicePlans, nil
}
