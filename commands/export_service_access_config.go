package commands

import (
	"github.com/vmwarepivotallabs/cf-mgmt/export"
)

type ExportServiceAccessCommand struct {
	BaseCFConfigCommand
}

// ExportServiceAccessCommand - updates commands to reverse engineer service access into cf-mgmt.yml and remove from orgConfig.yml if present
func (c *ExportServiceAccessCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
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
			cfMgmt.QuotaManager, cfMgmt.RoleManager)

		return exportManager.ExportServiceAccess()
	}
	return err
}
