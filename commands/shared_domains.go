package commands

type SharedDomainsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - adds/removes shared domains
func (c *SharedDomainsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.SharedDomainManager.Apply()
	}
	return err
}
