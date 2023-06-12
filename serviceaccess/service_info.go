package serviceaccess

import (
	"fmt"
	"net/url"
)

// GetServiceInfo - returns broker, it's services and their plans
func GetServiceInfo(client CFClient) (*ServiceInfo, error) {
	serviceInfo := &ServiceInfo{}
	brokers, err := client.ListServiceBrokers()
	if err != nil {
		return nil, err
	}
	for _, broker := range brokers {
		serviceBroker := &ServiceBroker{
			Name:      broker.Name,
			SpaceGUID: broker.SpaceGUID,
		}
		services, err := client.ListServicesByQuery(url.Values{
			"q": []string{fmt.Sprintf("%s:%s", "service_broker_guid", broker.Guid)},
		})
		if err != nil {
			return nil, err
		}
		for _, svc := range services {
			service := &Service{
				Name: svc.Label,
			}
			plans, err := client.ListServicePlansByQuery(url.Values{
				"q": []string{fmt.Sprintf("%s:%s", "service_guid", svc.Guid)},
			})
			if err != nil {
				return nil, err
			}
			for _, plan := range plans {
				servicePlan := &ServicePlanInfo{
					Name:        plan.Name,
					GUID:        plan.Guid,
					ServiceName: service.Name,
					Public:      plan.Public,
				}
				visibilities, err := client.ListServicePlanVisibilitiesByQuery(url.Values{
					"q": []string{fmt.Sprintf("%s:%s", "service_plan_guid", plan.Guid)},
				})
				if err != nil {
					return nil, err
				}
				for _, visibility := range visibilities {
					orgVisibility := &Visibility{
						OrgGUID:         visibility.OrganizationGuid,
						ServicePlanGUID: visibility.ServicePlanGuid,
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
