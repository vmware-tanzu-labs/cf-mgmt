package commands

type IsolationSegmentsCommand struct {
	BaseCFConfigCommand
}

//Execute - updates spaces
func (c *IsolationSegmentsCommand) Execute([]string) error {
	cfMgmt, err := InitializeManagers(c.BaseCFConfigCommand)
	if err != nil {
		return err
	}

	u := cfMgmt.IsolationSegmentUpdater
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
