package configcommands

import (
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type DeleteOrgConfigurationCommand struct {
	BaseConfigCommand
	OrgName         string `long:"org" description:"Org name to delete" required:"true"`
	ConfirmDeletion bool   `long:"confirm-deletion" description:"Confirm Deletion" required:"true"`
}

// Execute - deletes org from config
func (c *DeleteOrgConfigurationCommand) Execute([]string) error {
	if err := config.NewManager(c.ConfigDirectory).DeleteOrgConfig(c.OrgName); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("The org %s was successfully deleted", c.OrgName))
	return nil
}
