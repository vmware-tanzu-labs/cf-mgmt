package commands

type UpdateSpacesCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates spaces
func (c *UpdateSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaces()
	}
	return err
}
