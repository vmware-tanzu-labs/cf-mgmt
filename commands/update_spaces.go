package commands

type UpdateSpacesCommand struct {
	BaseCFConfigCommand
}

//Execute - updates spaces
func (c *UpdateSpacesCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpaces(c.ConfigDirectory)
	}
	return err
}
