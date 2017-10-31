package commands

import (
	"fmt"
)

var VERSION = "0.0.0-dev"

type VersionCommand struct {
}

//Execute - returns the version
func (c *VersionCommand) Execute([]string) error {
	fmt.Println(VERSION)
	return nil
}
