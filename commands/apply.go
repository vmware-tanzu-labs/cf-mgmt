package commands

import (
	"fmt"
)

type ApplyCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
	BaseLDAPCommand
}

//Execute - applies all the config in order
func (c *ApplyCommand) Execute([]string) error {
	var cfMgmt *CFMgmt
	var err error
	if cfMgmt, err = InitializePeekManagers(c.BaseCFConfigCommand, c.Peek); err == nil {
		fmt.Println("*********  Creating Orgs")
		if err = cfMgmt.OrgManager.CreateOrgs(); err != nil {
			return err
		}

		fmt.Println("*********  Creating Application Security Groups")
		if err = cfMgmt.SecurityGroupManager.CreateApplicationSecurityGroups(); err != nil {
			return err
		}

	}
	return err
}
