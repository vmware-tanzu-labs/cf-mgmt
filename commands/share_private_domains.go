package commands

type SharePrivateDomainsCommand struct {
	BaseCFConfigCommand
}

//Execute - creates private domains
func (c *SharePrivateDomainsCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.SharePrivateDomains()
	}
	return err
}
