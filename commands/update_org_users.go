package commands

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
}

//Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializeManagers(c.BaseCFConfigCommand); err == nil {
		err = cfMgmt.OrgManager.UpdateOrgUsers(c.ConfigDirectory, c.LdapPassword)
	}
	return err
}
