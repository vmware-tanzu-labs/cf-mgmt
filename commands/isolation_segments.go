package commands

type IsolationSegmentsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

//Execute - updates spaces
func (c *IsolationSegmentsCommand) Execute([]string) error {
	cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek)
	if err != nil {
		return err
	}

	u := cfMgmt.IsolationSegmentManager
	if err := u.Ensure(); err != nil {
		return err
	}
	if err := u.Entitle(); err != nil {
		return err
	}
	if err := u.UpdateOrgs(); err != nil {
		return err
	}
	if err := u.UpdateSpaces(); err != nil {
		return err
	}

	return nil
}
