package commands

type ServiceAccessCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - enables/disables service access
func (c *ServiceAccessCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.ServiceAccessManager.Apply()
	}
	return err
}
