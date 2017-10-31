package commands

type CreateOrgsCommand struct {
	BaseCFConfigCommand
}

//Execute - creates organizations
func (c *CreateOrgsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.CreateOrgs()
	}
	return err
}
