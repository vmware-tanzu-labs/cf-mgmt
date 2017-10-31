package commands

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
}

//Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaceUsers(c.ConfigDirectory, c.LdapPassword)
	}
	return err
}
