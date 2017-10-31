package commands

type CreateSecurityGroupsCommand struct {
	BaseCFConfigCommand
}

//Execute - creates security groups
func (c *CreateSecurityGroupsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SecurityGroupManager.CreateApplicationSecurityGroups()
	}
	return err
}
