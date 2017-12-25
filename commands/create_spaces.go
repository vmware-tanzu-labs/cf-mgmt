package commands

type CreateSpacesCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - creates spaces
func (c *CreateSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces(c.ConfigDirectory, c.LdapPassword)
	}
	return err
}
