package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type AddOrgToConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	OrgName                 string      `long:"org" description:"Org name" required:"true"`
	PrivateDomains          []string    `long:"private-domain" description:"Private Domain(s) to add, specify multiple times"`
	SharedPrivateDomains    []string    `long:"shared-private-domain" description:"Shared Private Domain(s) to add, specify multiple times"`
	DefaultIsolationSegment string      `long:"default-isolation-segment" description:"Default isolation segment for org" `
	Quota                   OrgQuota    `group:"quota"`
	BillingManager          UserRoleAdd `group:"billing-manager" namespace:"billing-manager"`
	Manager                 UserRoleAdd `group:"manager" namespace:"manager"`
	Auditor                 UserRoleAdd `group:"auditor" namespace:"auditor"`
	NamedQuota              string      `long:"named-quota" description:"Named quota to assign to org"`
	ServiceAccess           struct {
		ServiceNames []string `long:"service" description:"*** Deprecated **** - Service Name to add, specify multiple times"`
	} `group:"service-access"`
	EnableRemoveSpaces string `long:"enable-remove-spaces" description:"Enable removing spaces" choice:"true" choice:"false"`
}

//Execute - adds a named org to the configuration
func (c *AddOrgToConfigurationCommand) Execute([]string) error {
	lo.G.Warning("*** Deprecated *** - Use `org` command instead for adding/updating org configurations")
	orgConfig := &config.OrgConfig{
		Org: c.OrgName,
	}

	c.initConfig()

	if c.Quota.EnableOrgQuota == "true" && c.NamedQuota != "" {
		return fmt.Errorf("cannot enable org quota and use named quotas")
	}
	errorString := ""

	if c.DefaultIsolationSegment != "" {
		orgConfig.DefaultIsoSegment = c.DefaultIsolationSegment
	}
	orgConfig.RemoveUsers = true
	orgConfig.RemovePrivateDomains = true
	orgConfig.PrivateDomains = addToSlice(orgConfig.PrivateDomains, c.PrivateDomains, &errorString)

	orgConfig.RemoveSharedPrivateDomains = true
	orgConfig.SharedPrivateDomains = addToSlice(orgConfig.SharedPrivateDomains, c.SharedPrivateDomains, &errorString)

	updateOrgQuotaConfig(orgConfig, c.Quota, &errorString)
	orgConfig.NamedQuota = c.NamedQuota

	c.updateUsers(orgConfig, &errorString)

	globalCfg, err := c.ConfigManager.GetGlobalConfig()
	if err != nil {
		return err
	}
	if len(c.ServiceAccess.ServiceNames) > 0 {
		if globalCfg.IgnoreLegacyServiceAccess {
			return fmt.Errorf("Service access is managed with 'cf-mgmt-config global service-access' command")
		} else {
			lo.G.Warning("Service access is managed with 'cf-mgmt-config global service-access' command")
		}
		orgConfig.ServiceAccess = make(map[string][]string)
	}

	for _, service := range c.ServiceAccess.ServiceNames {
		orgConfig.ServiceAccess[service] = []string{"*"}
	}

	orgSpaces := &config.Spaces{Org: orgConfig.Org, EnableDeleteSpaces: true}
	convertToBool("enable-remove-spaces", &orgSpaces.EnableDeleteSpaces, c.EnableRemoveSpaces, &errorString)
	if errorString != "" {
		return errors.New(errorString)
	}

	if err := config.NewManager(c.ConfigDirectory).AddOrgToConfig(orgConfig); err != nil {
		return err
	}
	if err := config.NewManager(c.ConfigDirectory).SaveOrgSpaces(orgSpaces); err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("The org [%s] has been added", c.OrgName))
	return nil
}

func (c *AddOrgToConfigurationCommand) updateUsers(orgConfig *config.OrgConfig, errorString *string) {
	addUsersBasedOnRole(&orgConfig.BillingManager, orgConfig.GetBillingManagerGroups(), &c.BillingManager, errorString)
	addUsersBasedOnRole(&orgConfig.Auditor, orgConfig.GetAuditorGroups(), &c.Auditor, errorString)
	addUsersBasedOnRole(&orgConfig.Manager, orgConfig.GetManagerGroups(), &c.Manager, errorString)

	orgConfig.BillingManagerGroup = ""
	orgConfig.ManagerGroup = ""
	orgConfig.AuditorGroup = ""
}

func (c *AddOrgToConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
