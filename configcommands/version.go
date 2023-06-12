package configcommands

import (
	"fmt"
	"strings"
)

var VERSION = "0.0.0"
var COMMIT = "dev"

type VersionCommand struct {
}

// Execute - returns the version
func (c *VersionCommand) Execute([]string) error {
	fmt.Println(GetFormattedVersion())
	return nil
}

func GetVersion() string {
	return strings.Replace(VERSION, "v", "", 1)
}

func GetFormattedVersion() string {
	return fmt.Sprintf("Version: [%s], Commit: [%s]", GetVersion(), COMMIT)
}
