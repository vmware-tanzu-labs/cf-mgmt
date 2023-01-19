package commands

type UpdateOrgUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BaseAzureADCommand
	BasePeekCommand
}

// Execute - updates orgs quotas
func (c *UpdateOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.UserManager.InitializeLdap(c.LdapUser, c.LdapPassword, c.LdapServer); err != nil {
			return err
		}
		defer cfMgmt.UserManager.DeinitializeLdap()

		if err := cfMgmt.UserManager.InitializeAzureAD(c.AadTennantId, c.AadClientId, c.AadSecret, c.AADUserOrigin); err != nil {
			return err
		}

		return cfMgmt.UserManager.UpdateOrgUsers()
	}
	return nil
}
