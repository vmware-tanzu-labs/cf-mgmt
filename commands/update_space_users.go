package commands

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.SpaceUserManager.InitializeLdap(c.LdapPassword); err != nil {
			return err
		}
		return cfMgmt.SpaceUserManager.UpdateSpaceUsers()
	}
	return nil
}
