package commands

type UpdateOrgsMetadataCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates organizations metadata
func (c *UpdateOrgsMetadataCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.OrgManager.UpdateOrgsMetadata()
	}
	return err
}
