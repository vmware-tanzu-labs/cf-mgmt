package commands

import (
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/export"
	"github.com/xchapter7x/lo"
)

type ExportConfigurationCommand struct {
	BaseCFConfigCommand
	ExcludedOrgs      []string `long:"excluded-org" description:"Org to be excluded from export. Repeat the flag to specify multiple orgs"`
	ExcludedSpaces    []string `long:"excluded-space" description:"Space to be excluded from export. Repeat the flag to specify multiple spaces"`
	SkipSpaces        bool     `long:"skip-spaces" description:"Will not export space configurations"`
	SkipRoutingGroups bool     `long:"skip-routing-groups" description:"Will not export routing groups. Set to true if tcp routing is not configured"`
}

// Execute - initializes cf-mgmt configuration
func (c *ExportConfigurationCommand) Execute([]string) error {
	if cfMgmt, err := InitializeManagers(c.BaseCFConfigCommand); err != nil {
		lo.G.Errorf("Unable to initialize cf-mgmt. Error : %s", err)
		return err
	} else {
		exportManager := export.NewExportManager(c.ConfigDirectory,
			cfMgmt.UAAManager,
			cfMgmt.SpaceManager,
			cfMgmt.UserManager,
			cfMgmt.OrgReader,
			cfMgmt.SecurityGroupManager,
			cfMgmt.IsolationSegmentManager,
			cfMgmt.PrivateDomainManager,
			cfMgmt.SharedDomainManager,
			cfMgmt.ServiceAccessManager,
			cfMgmt.QuotaManager,
			cfMgmt.RoleManager,
		)
		exportManager.SkipRoutingGroups = c.SkipRoutingGroups
		excludedOrgs := make(map[string]string)
		for _, org := range config.DefaultProtectedOrgs {
			excludedOrgs[org] = org
		}
		for _, org := range c.ExcludedOrgs {
			excludedOrgs[org] = org
		}
		excludedSpaces := make(map[string]string)
		for _, space := range c.ExcludedSpaces {
			excludedSpaces[space] = space
		}
		lo.G.Infof("Orgs excluded from export by default: %v ", config.DefaultProtectedOrgs)
		lo.G.Infof("Orgs excluded from export by user:  %v ", c.ExcludedOrgs)
		lo.G.Infof("Spaces excluded from export by user:  %v ", c.ExcludedSpaces)
		err = exportManager.ExportConfig(excludedOrgs, excludedSpaces, c.SkipSpaces)
		if err != nil {
			lo.G.Errorf("Export failed with error:  %s", err)
			return err
		}
	}
	return nil
}
