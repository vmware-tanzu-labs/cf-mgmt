package commands

type SharePrivateDomainsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates private domains
func (c *SharePrivateDomainsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.PrivateDomainManager.SharePrivateDomains()
	}
	return err
}
