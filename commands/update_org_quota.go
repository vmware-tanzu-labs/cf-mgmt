package commands

type UpdateOrgQuotasCommand struct {
	BaseCFConfigCommand
}

//Execute - updates orgs quotas
func (c *UpdateOrgQuotasCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.CreateQuotas()
	}
	return err
}
