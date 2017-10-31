package commands

type CreateSpaceSecurityGroupsCommand struct {
	BaseCFConfigCommand
}

//Execute - creates space specific security groups
func (c *CreateSpaceSecurityGroupsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.CreateApplicationSecurityGroups(c.ConfigDirectory)
	}
	return err
}
