package configcommands

import (
	"fmt"

	"github.com/pivotalservices/cf-mgmt/config"
)

type AddUserToOrgConfigurationCommand struct {
	BaseConfigCommand
	OrgName  string `long:"org" description:"Org name" required:"true"`
	UserID   string `long:"user-id" description:"The user ID to add" required:"true"`
	RoleName string `long:"user-role" description:"The Org role name: managers, billing_managers or auditors" required:"true"`
	LdapUser bool   `long:"ldap-user" default:"false" description:"Boolean flag for whether the user is to be added into the LDAP Users. If blank, defaults to FALSE."`
}

//Execute - adds a user to appropriate role for a given org
func (c *AddUserToOrgConfigurationCommand) Execute([]string) error {

	cfg := config.NewManager(c.ConfigDirectory)
	orgConfig, err := cfg.GetOrgConfig(c.OrgName)
	if err != nil {
		return err
	}

	//TODO add user to role logic
	if err := cfg.SaveOrgConfig(orgConfig); err != nil {
		return err
	}
	userType := ""
	if c.LdapUser {
		userType = "LDAP "
	}
	fmt.Printf("%sUser %s was successfully added into %s with the %s role", userType, c.UserID, c.OrgName, c.RoleName)
	return nil

}
