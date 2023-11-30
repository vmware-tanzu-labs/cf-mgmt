package commands

import "fmt"

type CleanupOrgUsersCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - removes org users
func (c *CleanupOrgUsersCommand) Execute([]string) error {
	if cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil); err == nil {
		errs := cfMgmt.UserManager.CleanupOrgUsers()
		if len(errs) > 0 {
			return fmt.Errorf("got errors processing cleanup users %v", errs)
		}
		return nil
	}
	return nil
}
