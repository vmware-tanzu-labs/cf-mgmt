package commands

type UpdateSpacesMetadataCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates spaces metadata
func (c *UpdateSpacesMetadataCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SpaceManager.UpdateSpacesMetadata()
	}
	return err
}
