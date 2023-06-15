package commands

import "fmt"

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

// Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.UserManager.InitializeLdap(c.LdapUser, c.LdapPassword, c.LdapServer); err != nil {
			return err
		}
		defer cfMgmt.UserManager.DeinitializeLdap()
		errs := cfMgmt.UserManager.UpdateOrgUsers()
		if len(errs) > 0 {
			return fmt.Errorf("got errors processing update org users %v", errs)
		}
		return nil
	}
	return nil
}
