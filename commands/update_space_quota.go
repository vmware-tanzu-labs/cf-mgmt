package commands

type UpdateSpaceQuotasCommand struct {
	BaseCFConfigCommand
}

//Execute - updates space quotas
func (c *UpdateSpaceQuotasCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.SpaceManager.CreateQuotas(c.ConfigDirectory)
	}
	return err
}
