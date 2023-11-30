package commands

type CreatePrivateDomainsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - creates private domains
func (c *CreatePrivateDomainsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.PrivateDomainManager.CreatePrivateDomains()
	}
	return err
}
