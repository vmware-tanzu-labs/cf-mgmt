package commands

type DeleteSpacesCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

//Execute - deletes spaces
func (c *DeleteSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		err = cfMgmt.SpaceManager.DeleteSpaces(c.ConfigDirectory)
	}
	return err
}
