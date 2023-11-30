package commands

type CreateSecurityGroupsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates security groups
func (c *CreateSecurityGroupsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SecurityGroupManager.CreateGlobalSecurityGroups()
	}
	return err
}
