package commands

type DeleteSpacesCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

//Execute - deletes spaces
func (c *DeleteSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.DeleteSpaces(c.ConfigDirectory, c.Peek)
	}
	return err
}
