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
	DefaultServiceAccessNone      bool                    `yaml:"default-service-access-none"`
	ServiceAccess                 []Broker                `yaml:"service-access"`
}

type PlanInfo struct {
	Limited   bool
	AllAccess bool
	NoAccess  bool
	Orgs      []string
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
	Services []Service
}

type Service struct {
	Name               string           `yaml:"service"`
	AllAccessPlans     []string         `yaml:"all_access_plans,omitempty"`
	LimitedAccessPlans []PlanVisibility `yaml:"limited_access_plans,omitempty"`
	NoAccessPlans      []string         `yaml:"no_access_plans,omitempty"`
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
