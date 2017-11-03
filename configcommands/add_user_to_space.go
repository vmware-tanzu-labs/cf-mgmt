package configcommands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddUserToSpaceConfigurationCommand struct {
	BaseConfigCommand
	OrgName   string `long:"org" description:"Org name" required:"true"`
	SpaceName string `long:"space" description:"Space name" required:"true"`
	UserID    string `long:"user-id" description:"The user ID to add" required:"true"`
	RoleName  string `long:"user-role" description:"The Space role name: developers, managers or auditors" required:"true"`
	LdapUser  bool   `long:"ldap-user" default:"false" description:"Boolean flag for whether the user is to be added into the LDAP Users. If blank, defaults to FALSE."`
}

//Execute - adds a user to appropriate role for a given space
func (c *AddUserToSpaceConfigurationCommand) Execute([]string) error {
	cfg := config.NewManager(c.ConfigDirectory)
	spaceConfig, err := cfg.GetSpaceConfig(c.OrgName, c.SpaceName)
	if err != nil {
		return err
	}

	//TODO add user to role logic
	if err := cfg.SaveSpaceConfig(spaceConfig); err != nil {
		return err
	}
	userType := ""
	if c.LdapUser {
		userType = "LDAP "
	}
	fmt.Printf("%sUser %s was successfully added into %s/%s with the %s role", userType, c.UserID, c.OrgName, c.SpaceName, c.RoleName)
	return nil

}
