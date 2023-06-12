package configcommands

import (
	"errors"
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type GlobalConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	EnableDeleteIsolationSegments  string              `long:"enable-delete-isolation-segments" description:"Enable removing isolation segments" choice:"true" choice:"false"`
	EnableDeleteSharedDomains      string              `long:"enable-delete-shared-domains" description:"Enable removing shared domains" choice:"true" choice:"false"`
	EnableServiceAccess            string              `long:"enable-service-access" description:"Enable managing service access" choice:"true" choice:"false"`
	EnableUnassignSecurityGroups   string              `long:"enable-unassign-security-groups" description:"Enable unassigning security groups" choice:"true" choice:"false"`
	MetadataPrefix                 string              `long:"metadata-prefix" description:"Prefix for org/space metadata"`
	StagingSecurityGroups          []string            `long:"staging-security-group" description:"Staging Security Group to add"`
	RemoveStagingSecurityGroups    []string            `long:"remove-staging-security-group" description:"Staging Security Group to remove"`
	RunningSecurityGroups          []string            `long:"running-security-group" description:"Running Security Group to add"`
	RemoveRunningSecurityGroups    []string            `long:"remove-running-security-group" description:"Running Security Group to remove"`
	SharedDomains                  []string            `long:"shared-domain" description:"Shared Domain to add"`
	RouterGroupSharedDomains       []string            `long:"router-group-shared-domain" description:"Router Group Shared Domain to add"`
	RouterGroupSharedDomainsGroups []string            `long:"router-group-shared-domain-group" description:"Router Group Shared Domain group"`
	InternalSharedDomains          []string            `long:"internal-shared-domain" description:"Internal Shared Domain to add"`
	RemoveSharedDomains            []string            `long:"remove-shared-domain" description:"Shared Domain to remove"`
	ServiceAccess                  GlobalServiceAccess `group:"service-access"`
}

// Execute - adds/updates a named asg to the configuration
func (c *GlobalConfigurationCommand) Execute([]string) error {
	c.initConfig()

	globalConfig, err := c.ConfigManager.GetGlobalConfig()
	if err != nil {
		return err
	}
	errorString := ""

	if globalConfig.SharedDomains == nil {
		globalConfig.SharedDomains = map[string]config.SharedDomain{}
	}

	convertToBool("enable-delete-isolation-segments", &globalConfig.EnableDeleteIsolationSegments, c.EnableDeleteIsolationSegments, &errorString)
	convertToBool("enable-delete-shared-domains", &globalConfig.EnableDeleteSharedDomains, c.EnableDeleteSharedDomains, &errorString)
	convertToBool("enable-service-access", &globalConfig.EnableServiceAccess, c.EnableServiceAccess, &errorString)
	convertToBool("enable-unassign-security-groups", &globalConfig.EnableUnassignSecurityGroups, c.EnableUnassignSecurityGroups, &errorString)
	if c.MetadataPrefix != "" {
		globalConfig.MetadataPrefix = c.MetadataPrefix
	}

	globalConfig.StagingSecurityGroups = c.updateSecGroups(globalConfig.StagingSecurityGroups, c.StagingSecurityGroups, c.RemoveStagingSecurityGroups)
	globalConfig.RunningSecurityGroups = c.updateSecGroups(globalConfig.RunningSecurityGroups, c.RunningSecurityGroups, c.RemoveRunningSecurityGroups)

	for _, domain := range c.SharedDomains {
		globalConfig.SharedDomains[domain] = config.SharedDomain{Internal: false}
	}
	for _, domain := range c.InternalSharedDomains {
		globalConfig.SharedDomains[domain] = config.SharedDomain{Internal: true}
	}

	if len(c.RouterGroupSharedDomains) > 0 {
		if len(c.RouterGroupSharedDomains) != len(c.RouterGroupSharedDomainsGroups) {
			return fmt.Errorf("Must specify same number of router-group-shared-domain args as router-group-shared-domain-group args")
		}

		for index, domain := range c.RouterGroupSharedDomains {
			globalConfig.SharedDomains[domain] = config.SharedDomain{
				Internal:    false,
				RouterGroup: c.RouterGroupSharedDomainsGroups[index],
			}
		}
	}

	for _, domain := range c.RemoveSharedDomains {
		delete(globalConfig.SharedDomains, domain)
	}

	errorList := c.UpdateServiceAccess(globalConfig)
	for _, err := range errorList {
		errorString += "\n--" + err.Error()
	}

	if errorString != "" {
		return errors.New(errorString)
	}

	err = c.ConfigManager.SaveGlobalConfig(globalConfig)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The cf-mgmt.yml has been updated"))
	return nil
}

func (c *GlobalConfigurationCommand) updateSecGroups(current, additions, removals []string) []string {
	secGroupMap := make(map[string]string)
	for _, secGroup := range current {
		secGroupMap[secGroup] = secGroup
	}
	for _, secGroup := range additions {
		secGroupMap[secGroup] = secGroup
	}
	for _, secGroup := range removals {
		delete(secGroupMap, secGroup)
	}

	var result []string
	for _, secGroup := range secGroupMap {
		result = append(result, secGroup)
	}

	return result
}

func (c *GlobalConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}

func (c *GlobalConfigurationCommand) UpdateServiceAccess(globalConfig *config.GlobalConfig) []error {
	var errorList []error
	if len(c.ServiceAccess.AllAccessPlan) > 0 || len(c.ServiceAccess.LimitedAccessPlan) > 0 || len(c.ServiceAccess.NoAccessPlan) > 0 {
		if len(c.ServiceAccess.Broker) == 0 {
			errorList = append(errorList, fmt.Errorf("must specify --broker arg"))
		}
		if len(c.ServiceAccess.Service) == 0 {
			errorList = append(errorList, fmt.Errorf("must specify --service arg"))
		}
		if len(errorList) > 0 {
			return errorList
		}
		broker := globalConfig.GetBroker(c.ServiceAccess.Broker)
		service := broker.GetService(c.ServiceAccess.Service)
		if len(c.ServiceAccess.AllAccessPlan) > 0 {
			service.AddAllAccessPlan(c.ServiceAccess.AllAccessPlan)
		}
		if len(c.ServiceAccess.NoAccessPlan) > 0 {
			service.AddNoAccessPlan(c.ServiceAccess.NoAccessPlan)
		}

		if len(c.ServiceAccess.LimitedAccessPlan) > 0 {
			service.AddLimitedAccessPlan(c.ServiceAccess.LimitedAccessPlan, c.ServiceAccess.OrgsToAdd, c.ServiceAccess.OrgsToRemove)
		}
	}

	return nil
}

type GlobalServiceAccess struct {
	Broker  string `long:"broker" description:"Name of Broker"`
	Service string `long:"service" description:"Name of Service"`

	AllAccessPlan     string   `long:"all-access-plan" description:"Plan to give access to all orgs"`
	LimitedAccessPlan string   `long:"limited-access-plan" description:"Plan to give limited access to, must also provide org list"`
	OrgsToAdd         []string `long:"org" description:"Orgs to add to limited plan"`
	OrgsToRemove      []string `long:"remove-org" description:"Orgs to remove from limited plan"`
	NoAccessPlan      string   `long:"no-access-plan" description:"Plan to give access to all orgs"`
}
