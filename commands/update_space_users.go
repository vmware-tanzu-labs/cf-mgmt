package commands

type UpdateSpaceUsersCommand struct {
	BaseCFConfigCommand
	BaseLDAPCommand
	BaseAzureADCommand
	BasePeekCommand
}

// Execute - updates space users
func (c *UpdateSpaceUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		if err := cfMgmt.UserManager.InitializeLdap(c.LdapUser, c.LdapPassword, c.LdapServer); err != nil {
			return err
		}
		defer cfMgmt.UserManager.DeinitializeLdap()

		if err := cfMgmt.UserManager.InitializeAzureAD(c.AadTenantId, c.AadClientId, c.AadSecret, c.AADUserOrigin); err != nil {
			return err
		}
		return cfMgmt.UserManager.UpdateSpaceUsers()
	}
	return nil
}
