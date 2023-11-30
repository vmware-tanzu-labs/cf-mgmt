package commands

type UpdateOrgQuotasCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates orgs quotas
func (c *UpdateOrgQuotasCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		err = cfMgmt.QuotaManager.CreateOrgQuotas()
	}
	return err
}
