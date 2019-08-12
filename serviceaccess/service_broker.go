package serviceaccess

import "fmt"

type ServiceBroker struct {
	Name     string
	services map[string]*Service
}

func (s *ServiceBroker) Services() []*Service {
	serviceList := []*Service{}
	for _, service := range s.services {
		serviceList = append(serviceList, service)
	}
	return serviceList
}

func (s *ServiceBroker) GetService(serviceName string) (*Service, error) {
	if service, ok := s.services[serviceName]; ok {
		return service, nil
	}
	return nil, fmt.Errorf("Service %s is not found for broker %s", serviceName, s.Name)
}

func (s *ServiceBroker) AddService(service *Service) {
	if s.services == nil {
		s.services = make(map[string]*Service)
	}
	s.services[service.Name] = service
}
