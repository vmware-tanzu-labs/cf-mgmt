package commands

type UpdateSpaceQuotasCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates space quotas
func (c *UpdateSpaceQuotasCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.QuotaManager.CreateSpaceQuotas()
	}
	return err
}
