package commands

type CreateSpacesCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates spaces
func (c *CreateSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SpaceManager.CreateSpaces()
	}
	return err
}
