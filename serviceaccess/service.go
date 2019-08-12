package serviceaccess

import (
	"fmt"
	"strings"
)

type Service struct {
	Name  string
	plans map[string]*ServicePlanInfo
}

func (s *Service) AddPlan(servicePlan *ServicePlanInfo) {
	if s.plans == nil {
		s.plans = make(map[string]*ServicePlanInfo)
	}
	s.plans[servicePlan.Name] = servicePlan
}

func (s *Service) GetPlans(planNames []string) ([]*ServicePlanInfo, error) {
	var servicePlans []*ServicePlanInfo
	for planName, plan := range s.plans {
		if matches(planName, planNames) {
			servicePlans = append(servicePlans, plan)
		}
	}

	if len(servicePlans) == 0 {
		return nil, fmt.Errorf("No plans for for service %s with expected plans %v", s.Name, planNames)
	}

	return servicePlans, nil
}

func (s *Service) GetPlan(plan string) (*ServicePlanInfo, error) {
	for planName, servicePlan := range s.plans {
		if strings.EqualFold(planName, plan) {
			return servicePlan, nil
		}
	}
	return nil, nil
}

func matches(planName string, planList []string) bool {
	for _, name := range planList {
		if name == "*" {
			return true
		}
		if strings.EqualFold(planName, name) {
			return true
		}
	}
	return false
}

func (s *Service) GetPlanNames() ([]string, error) {
	planNames := []string{}
	for planName := range s.plans {
		planNames = append(planNames, planName)
	}
	return planNames, nil
}

func (s *Service) Plans() []*ServicePlanInfo {
	plans := []*ServicePlanInfo{}
	for _, plan := range s.plans {
		plans = append(plans, plan)
	}
	return plans
}
