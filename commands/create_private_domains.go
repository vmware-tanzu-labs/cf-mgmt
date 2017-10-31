package commands

type CreatePrivateDomainsCommand struct {
	BaseCFConfigCommand
}

//Execute - creates private domains
func (c *CreatePrivateDomainsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.CreatePrivateDomains()
	}
	return err
}
