package commands

type CreateSpaceSecurityGroupsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates space specific security groups
func (c *CreateSpaceSecurityGroupsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SecurityGroupManager.CreateApplicationSecurityGroups()
	}
	return err
}
