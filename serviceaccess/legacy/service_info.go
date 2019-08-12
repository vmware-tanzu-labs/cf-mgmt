package legacy

import (
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
)

type ServiceInfo struct {
	m map[string]map[string]*ServicePlanInfo
}

type ServicePlanInfo struct {
	GUID   string
	Name   string
	Public bool
	m      map[string]*cfclient.ServicePlanVisibility
}

func (s *ServicePlanInfo) ListVisibilities() []cfclient.ServicePlanVisibility {
	var result []cfclient.ServicePlanVisibility
	for _, visibility := range s.m {
		result = append(result, *visibility)
	}
	return result
}
func (s *ServicePlanInfo) AddOrg(orgGUID string, visibility cfclient.ServicePlanVisibility) {
	if s.m == nil {
		s.m = make(map[string]*cfclient.ServicePlanVisibility)
	}
	s.m[orgGUID] = &visibility
}

func (s *ServicePlanInfo) RemoveOrg(orgGUID string) {
	delete(s.m, orgGUID)
}

func (s *ServicePlanInfo) OrgHasAccess(orgGUID string) bool {
	_, ok := s.m[orgGUID]
	return ok
}

func (s *ServiceInfo) AddPlan(serviceName string, servicePlan cfclient.ServicePlan) *ServicePlanInfo {
	if s.m == nil {
		s.m = make(map[string]map[string]*ServicePlanInfo)
	}
	plans, ok := s.m[serviceName]
	if !ok {
		plans = make(map[string]*ServicePlanInfo)
		s.m[serviceName] = plans
	}
	servicePlanInfo := &ServicePlanInfo{GUID: servicePlan.Guid, Name: servicePlan.Name, Public: servicePlan.Public}
	plans[servicePlan.Name] = servicePlanInfo
	return servicePlanInfo
}

func (s *ServiceInfo) GetPlans(serviceName string, planNames []string) ([]*ServicePlanInfo, error) {
	plans, ok := s.m[serviceName]
	if !ok {
		return nil, fmt.Errorf("Service %s does not exist", serviceName)
	}
	var servicePlans []*ServicePlanInfo
	for planName, plan := range plans {
		if Matches(planName, planNames) {
			servicePlans = append(servicePlans, plan)
		}
	}

	if len(servicePlans) == 0 {
		return nil, fmt.Errorf("No plans for for service %s with expected plans %v", serviceName, planNames)
	}

	return servicePlans, nil
}

func Matches(planName string, planList []string) bool {
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

func (s *ServiceInfo) AllPlans() map[string][]*ServicePlanInfo {
	allPlans := make(map[string][]*ServicePlanInfo)
	for serviceName, planMap := range s.m {
		var plans []*ServicePlanInfo
		for _, plan := range planMap {
			plans = append(plans, plan)
		}
		allPlans[serviceName] = plans
	}
	return allPlans
}
