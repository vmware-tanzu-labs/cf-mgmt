package commands

type AssignDefaultSecurityGroups struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates security groups
func (c *AssignDefaultSecurityGroups) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SecurityGroupManager.AssignDefaultSecurityGroups()
	}
	return err
}
