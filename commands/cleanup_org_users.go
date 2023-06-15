package commands

type CleanupOrgUsersCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - removes org users
func (c *CleanupOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		return cfMgmt.UserManager.CleanupOrgUsers()
	}
	return nil
}
