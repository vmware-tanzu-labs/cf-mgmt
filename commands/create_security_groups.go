package commands

type CreateSecurityGroupsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

//Execute - creates security groups
func (c *CreateSecurityGroupsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		err = cfMgmt.SecurityGroupManager.CreateApplicationSecurityGroups()
	}
	return err
}
