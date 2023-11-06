package commands

import "fmt"

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

// Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	ldapMgr, err := InitializeLdapManager(c.BaseCFConfigCommand, c.BaseLDAPCommand)
	if err != nil {
		return err
	}
	if ldapMgr != nil {
		defer ldapMgr.Close()
	}
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, ldapMgr); err == nil {
		errs := cfMgmt.UserManager.UpdateOrgUsers()
		if len(errs) > 0 {
			return fmt.Errorf("got errors processing update org users %v", errs)
		}
		return nil
	}
	return nil
}
