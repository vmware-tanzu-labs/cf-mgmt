package configcommands

import (
	"errors"
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type UpdateOrgsConfigurationCommand struct {
	ConfigManager config.Manager
	BaseConfigCommand
	EnableDeleteOrgs string `long:"enable-delete-orgs" description:"Enable delete orgs option" choice:"true" choice:"false"`
}

//Execute - updates org configuration`
func (c *UpdateOrgsConfigurationCommand) Execute(args []string) error {
	c.initConfig()
	orgs, err := c.ConfigManager.Orgs()
	if err != nil {
		return err
	}
	errorString := ""
	convertToBool("enable-delete-orgs", &orgs.EnableDeleteOrgs, c.EnableDeleteOrgs, &errorString)

	if errorString != "" {
		return errors.New(errorString)
	}

	if err := c.ConfigManager.SaveOrgs(orgs); err != nil {
		return err
	}
	fmt.Println("The orgs.yml has been updated")
	return nil
}

func (c *UpdateOrgsConfigurationCommand) initConfig() {
	if c.ConfigManager == nil {
		c.ConfigManager = config.NewManager(c.ConfigDirectory)
	}
}
