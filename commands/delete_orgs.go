package commands

type DeleteOrgsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - deletes orgs
func (c *DeleteOrgsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.OrgManager.DeleteOrgs()
	}
	return err
}
