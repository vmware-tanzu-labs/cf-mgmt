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
	cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, ldapMgr)
	if err != nil {
		return err
	}
	errs := cfMgmt.UserManager.UpdateOrgUsers()
	if len(errs) > 0 {
		return fmt.Errorf("got errors processing update org users %v", errs)
	}
	return nil
}
