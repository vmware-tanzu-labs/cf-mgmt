package commands

import (
	"github.com/pivotalservices/cf-mgmt/export"
	"github.com/xchapter7x/lo"
)

type ExportConfigurationCommand struct {
	BaseCFConfigCommand
	ExcludedOrgs   []string `long:"excluded-org" description:"Org to be excluded from export. Repeat the flag to specify multiple orgs"`
	ExcludedSpaces []string `long:"excluded-space" description:"Space to be excluded from export. Repeat the flag to specify multiple spaces"`
}

//Execute - initializes cf-mgmt configuration
func (c *ExportConfigurationCommand) Execute([]string) error {
	if cfMgmt, err := InitializeManagers(c.BaseCFConfigCommand); err != nil {
		lo.G.Errorf("Unable to initialize cf-mgmt. Error : %s", err)
		return err
	} else {
		exportManager := export.NewExportManager(c.ConfigDirectory, cfMgmt.UAACManager, cfMgmt.CloudController)
		excludedOrgs := make(map[string]string)
		excludedOrgs["system"] = "system"
		for _, org := range c.ExcludedOrgs {
			excludedOrgs[org] = org
		}
		excludedSpaces := make(map[string]string)
		for _, space := range c.ExcludedSpaces {
			excludedSpaces[space] = space
		}
		lo.G.Info("Orgs excluded from export by default: [system]")
		lo.G.Infof("Orgs excluded from export by user:  %v ", c.ExcludedOrgs)
		lo.G.Infof("Spaces excluded from export by user:  %v ", c.ExcludedSpaces)
		err = exportManager.ExportConfig(excludedOrgs, excludedSpaces)
		if err != nil {
			lo.G.Errorf("Export failed with error:  %s", err)
			return err
		}
	}
	return nil
}
