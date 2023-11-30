package commands

import "fmt"

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BaseAzureADCommand
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
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, ldapMgr); err == nil {
		errs := cfMgmt.UserManager.UpdateSpaceUsers()
		if len(errs) > 0 {
			return fmt.Errorf("got errors processing update space users %v", errs)
		}

		if err := cfMgmt.UserManager.InitializeAzureAD(c.AadTenantId, c.AadClientId, c.AadSecret, c.AADUserOrigin); err != nil {
			return err
		}

		return nil
	}
	return nil
}
