package commands

type CreateOrgsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates organizations
func (c *CreateOrgsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.OrgManager.CreateOrgs()
	}
	return err
}
