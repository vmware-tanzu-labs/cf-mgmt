package configcommands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type DeleteSpaceConfigurationCommand struct {
	BaseConfigCommand
	OrgName         string `long:"org" description:"Org name to delete" required:"true"`
	SpaceName       string `long:"space" description:"Space name to delete" required:"true"`
	ConfirmDeletion bool   `long:"confirm-deletion" default:"false" description:"REQUIRED: Confirm Deletion" required:"true"`
}

//Execute - deletes space from config
func (c *DeleteSpaceConfigurationCommand) Execute([]string) error {
	if err := config.NewManager(c.ConfigDirectory).DeleteSpaceConfig(c.OrgName, c.SpaceName); err != nil {
		return err
	}

	fmt.Printf("The org/space %s/%s was successfully deleted", c.OrgName, c.SpaceName)
	return nil
}
