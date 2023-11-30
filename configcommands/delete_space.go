package configcommands

import (
	"fmt"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
)

type DeleteSpaceConfigurationCommand struct {
	BaseConfigCommand
	OrgName         string `long:"org" description:"Org name of space to delete" required:"true"`
	SpaceName       string `long:"space" description:"Space name to delete" required:"true"`
	ConfirmDeletion bool   `long:"confirm-deletion" description:"Confirm Deletion" required:"true"`
}

// Execute - deletes space from config
func (c *DeleteSpaceConfigurationCommand) Execute([]string) error {
	if err := config.NewManager(c.ConfigDirectory).DeleteSpaceConfig(c.OrgName, c.SpaceName); err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("The org/space %s/%s was successfully deleted", c.OrgName, c.SpaceName))
	return nil
}
