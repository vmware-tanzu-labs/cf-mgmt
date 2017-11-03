package configcommands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddPrivateDomainToOrgConfigurationCommand struct {
	BaseConfigCommand
	OrgName       string `long:"org" description:"Org name to add" required:"true"`
	PrivateDomain string `long:"private-domain-name" env:"PRIVATE_DOMAIN_NAME" default:"false" description:"Private domain name. HTTP or HTTPS only" required:"true"`
}

//Execute - adds private domain for a given org
func (c *AddPrivateDomainToOrgConfigurationCommand) Execute([]string) error {

	cfg := config.NewManager(c.ConfigDirectory)
	orgConfig, err := cfg.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}

	orgConfig.PrivateDomains = append(orgConfig.PrivateDomains, c.PrivateDomain)
	if err := cfg.SaveOrgConfig(orgConfig); err != nil {
		return err
	}
	fmt.Printf("The private domain %s was successfully added into %s", c.PrivateDomain, c.OrgName)
	return nil

}
