package legacy

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/resource"
	"strings"
)

type ServiceInfo struct {
	m map[string]map[string]*ServicePlanInfo
}

type ServicePlanInfo struct {
	GUID   string
	Name   string
	Public bool
	m      map[string]*resource.ServicePlanVisibility
}

type OrgServicePlanVisibilityPair struct {
	OrgGUID               string
	ServicePlanVisibility *resource.ServicePlanVisibility
}

func (s *ServicePlanInfo) ListVisibilitiesByOrg() []OrgServicePlanVisibilityPair {
	var result []OrgServicePlanVisibilityPair
	for org, visibility := range s.m {
		v := OrgServicePlanVisibilityPair{
			OrgGUID:               org,
			ServicePlanVisibility: visibility,
		}
		result = append(result, v)
	}
	return result
}

func (s *ServicePlanInfo) ListVisibilities() []*resource.ServicePlanVisibility {
	var result []*resource.ServicePlanVisibility
	for _, visibility := range s.m {
		result = append(result, visibility)
	}
	return result
}
func (s *ServicePlanInfo) AddOrg(orgGUID string, visibility *resource.ServicePlanVisibility) {
	if s.m == nil {
		s.m = make(map[string]*resource.ServicePlanVisibility)
	}
	s.m[orgGUID] = visibility
}

func (s *ServicePlanInfo) RemoveOrg(orgGUID string) {
	delete(s.m, orgGUID)
}

func (s *ServicePlanInfo) OrgHasAccess(orgGUID string) bool {
	_, ok := s.m[orgGUID]
	return ok
}

func (s *ServiceInfo) AddPlan(serviceName string, servicePlan *resource.ServicePlan) *ServicePlanInfo {
	if s.m == nil {
		s.m = make(map[string]map[string]*ServicePlanInfo)
	}
	plans, ok := s.m[serviceName]
	if !ok {
		plans = make(map[string]*ServicePlanInfo)
		s.m[serviceName] = plans
	}
	public := servicePlan.VisibilityType == resource.ServicePlanVisibilityPublic.String()
	servicePlanInfo := &ServicePlanInfo{GUID: servicePlan.GUID, Name: servicePlan.Name, Public: public}
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
