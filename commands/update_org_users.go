package commands

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.OrgManager.InitializeLdap(c.LdapPassword); err != nil {
			return err
		}
		return cfMgmt.OrgManager.UpdateOrgUsers()
	}
	return nil
}
