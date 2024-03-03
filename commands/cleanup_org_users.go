package commands

import "fmt"

type CleanupOrgUsersCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - removes org users
func (c *CleanupOrgUsersCommand) Execute([]string) error {
	cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil)
	if err != nil {
		return err
	}
	errs := cfMgmt.UserManager.CleanupOrgUsers()
	if len(errs) > 0 {
		return fmt.Errorf("got errors processing cleanup users %v", errs)
	}
	return nil
}
