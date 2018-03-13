package commands

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		err = cfMgmt.SpaceUserManager.UpdateSpaceUsers(c.ConfigDirectory, c.LdapPassword)
	}
	return err
}
