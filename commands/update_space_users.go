package commands

import "fmt"

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

// Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
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
	errs := cfMgmt.UserManager.UpdateSpaceUsers()
	if len(errs) > 0 {
		return fmt.Errorf("got errors processing update space users %v", errs)
	}
	return nil
}
