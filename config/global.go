package config

import "strings"

// GlobalConfig configuration for global settings
type GlobalConfig struct {
	EnableDeleteIsolationSegments bool                    `yaml:"enable-delete-isolation-segments"`
	EnableUnassignSecurityGroups  bool                    `yaml:"enable-unassign-security-groups"`
	RunningSecurityGroups         []string                `yaml:"running-security-groups"`
	StagingSecurityGroups         []string                `yaml:"staging-security-groups"`
	SharedDomains                 map[string]SharedDomain `yaml:"shared-domains"`
	EnableDeleteSharedDomains     bool                    `yaml:"enable-remove-shared-domains"`
	MetadataPrefix                string                  `yaml:"metadata-prefix"`
	EnableServiceAccess           bool                    `yaml:"enable-service-access"`
	// DefaultServiceAccessNone      bool                    `yaml:"default-service-access-none"`
	ServiceAccess []*Broker `yaml:"service-access"`
}

type PlanInfo struct {
	Limited   bool
	AllAccess bool
	NoAccess  bool
	Orgs      []string
}

func (g *GlobalConfig) GetBroker(brokerName string) *Broker {
	for _, broker := range g.ServiceAccess {
		if strings.EqualFold(brokerName, broker.Name) {
			return broker
		}
	}
	newBroker := &Broker{Name: brokerName}
	g.ServiceAccess = append(g.ServiceAccess, newBroker)
	return newBroker
}
func (b *Broker) GetService(serviceName string) *Service {
	for _, service := range b.Services {
		if strings.EqualFold(serviceName, service.Name) {
			return service
		}
	}
	newService := &Service{Name: serviceName}
	b.Services = append(b.Services, newService)
	return newService
}

func (g *GlobalConfig) GetPlanInfo(brokerName, serviceName, planName string) PlanInfo {
	planInfo := PlanInfo{}
	for _, broker := range g.ServiceAccess {
		if strings.EqualFold(brokerName, broker.Name) {
			for _, service := range broker.Services {
				if strings.EqualFold(serviceName, service.Name) {
					for _, plan := range service.NoAccessPlans {
						if strings.EqualFold(planName, plan) {
							planInfo.NoAccess = true
							return planInfo
						}
					}
					for _, plan := range service.LimitedAccessPlans {
						if strings.EqualFold(planName, plan.Name) {
							planInfo.Limited = true
							planInfo.Orgs = plan.Orgs
							return planInfo
						}
					}
				}
			}
		}
	}
	//default to always have plan enabled
	planInfo.AllAccess = true
	return planInfo
}

type SharedDomain struct {
	Internal    bool   `yaml:"internal"`
	RouterGroup string `yaml:"router-group,omitempty"`
}

type Broker struct {
	Name     string `yaml:"broker"`
	Services []*Service
}

type Service struct {
	Name               string            `yaml:"service"`
	AllAccessPlans     []string          `yaml:"all_access_plans,omitempty"`
	LimitedAccessPlans []*PlanVisibility `yaml:"limited_access_plans,omitempty"`
	NoAccessPlans      []string          `yaml:"no_access_plans,omitempty"`
}

func (s *Service) AddAllAccessPlan(planName string) {
	if s.contains(s.NoAccessPlans, planName) {
		s.NoAccessPlans = s.remove(s.NoAccessPlans, planName)
	}
	if s.contains(s.LimitedAccessPlanNames(), planName) {
		s.LimitedAccessPlans = s.removePlan(s.LimitedAccessPlans, planName)
	}
	if !s.contains(s.AllAccessPlans, planName) {
		s.AllAccessPlans = append(s.AllAccessPlans, planName)
	}
}

func (s *Service) AddNoAccessPlan(planName string) {
	if s.contains(s.AllAccessPlans, planName) {
		s.AllAccessPlans = s.remove(s.AllAccessPlans, planName)
	}
	if s.contains(s.LimitedAccessPlanNames(), planName) {
		s.LimitedAccessPlans = s.removePlan(s.LimitedAccessPlans, planName)
	}
	if !s.contains(s.NoAccessPlans, planName) {
		s.NoAccessPlans = append(s.NoAccessPlans, planName)
	}
}

func (s *Service) AddLimitedAccessPlan(planName string, orgsToAdd, orgsToRemove []string) {
	if s.contains(s.AllAccessPlans, planName) {
		s.AllAccessPlans = s.remove(s.AllAccessPlans, planName)
	}
	if s.contains(s.NoAccessPlans, planName) {
		s.NoAccessPlans = s.remove(s.NoAccessPlans, planName)
	}
	if !s.contains(s.LimitedAccessPlanNames(), planName) {
		s.LimitedAccessPlans = append(s.LimitedAccessPlans, &PlanVisibility{Name: planName, Orgs: orgsToAdd})
	} else {
		planVisibility := s.GetLimitedPlan(planName)
		for _, org := range orgsToAdd {
			if !s.contains(planVisibility.Orgs, org) {
				planVisibility.Orgs = append(planVisibility.Orgs, org)
			}
		}
		for _, org := range orgsToRemove {
			planVisibility.Orgs = s.remove(planVisibility.Orgs, org)
		}
	}

}

func (s *Service) contains(slice []string, e string) bool {
	for _, a := range slice {
		if a == e {
			return true
		}
	}
	return false
}

func (s *Service) remove(slice []string, e string) []string {
	sliceToReturn := []string{}
	for _, a := range slice {
		if a != e {
			sliceToReturn = append(sliceToReturn, a)
		}
	}
	return sliceToReturn
}

func (s *Service) removePlan(slice []*PlanVisibility, e string) []*PlanVisibility {
	sliceToReturn := []*PlanVisibility{}
	for _, a := range slice {
		if a.Name != e {
			sliceToReturn = append(sliceToReturn, a)
		}
	}
	return sliceToReturn
}

func (s *Service) GetLimitedPlan(planName string) *PlanVisibility {
	for _, plan := range s.LimitedAccessPlans {
		if strings.EqualFold(planName, plan.Name) {
			return plan
		}
	}
	newPlan := &PlanVisibility{Name: planName}
	s.LimitedAccessPlans = append(s.LimitedAccessPlans, newPlan)
	return newPlan
}

func (s *Service) LimitedAccessPlanNames() []string {
	names := []string{}
	for _, plan := range s.LimitedAccessPlans {
		names = append(names, plan.Name)
	}
	return names
}

type PlanVisibility struct {
	Name string   `yaml:"plan,omitempty"`
	Orgs []string `yaml:"orgs,omitempty"`
}
