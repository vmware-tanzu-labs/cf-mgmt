package commands

import (
	"github.com/pivotalservices/cf-mgmt/configcommands"
	"github.com/xchapter7x/lo"
)

type GenerateConcoursePipelineCommand struct {
}

//Execute - generates concourse pipeline and tasks
func (c *GenerateConcoursePipelineCommand) Execute(params []string) error {
	lo.G.Warning("This command has been deprecated use lastest cf-mgmt-config cli")
	cfg := configcommands.GenerateConcoursePipelineCommand{}
	return cfg.Execute(params)
}
