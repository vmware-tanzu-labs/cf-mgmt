package commands

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BasePeekCommand
}

//Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.UserManager.InitializeLdap(c.LdapUser, c.LdapPassword, c.LdapServer); err != nil {
			return err
		}
		defer cfMgmt.UserManager.DeinitializeLdap()
		return cfMgmt.UserManager.UpdateOrgUsers()
	}
	return nil
}
