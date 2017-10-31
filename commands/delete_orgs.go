package commands

type DeleteOrgsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

//Execute - deletes orgs
func (c *DeleteOrgsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.DeleteOrgs(c.Peek)
	}
	return err
}
