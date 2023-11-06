package commands

type IsolationSegmentsCommand struct {
	BaseCFConfigCommand
	BasePeekCommand
}

// Execute - updates spaces
func (c *IsolationSegmentsCommand) Execute([]string) error {
	cfMgmt, err := InitializePeekManagers(c.BaseCFConfigCommand, c.Peek, nil)
	if err != nil {
		return err
	}

	return cfMgmt.IsolationSegmentManager.Apply()

}
