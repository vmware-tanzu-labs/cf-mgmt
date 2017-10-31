package commands

type CreateSpacesCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
}

//Execute - creates spaces
func (c *CreateSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces(c.ConfigDirectory, c.LdapPassword)
	}
	return err
}
