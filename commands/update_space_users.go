package commands

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.UserManager.InitializeLdap(c.LdapPassword); err != nil {
			return err
		}
		defer cfMgmt.UserManager.DeinitializeLdap()
		return cfMgmt.UserManager.UpdateSpaceUsers()
	}
	return nil
}
