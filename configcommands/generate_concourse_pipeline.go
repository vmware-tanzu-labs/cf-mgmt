package configcommands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/embedded"
	"github.com/xchapter7x/lo"
)

type GenerateConcoursePipelineCommand struct {
	BaseConfigCommand
	TargetDirectory string `long:"target-dir" default:"." description:"Name of the target directory to generate into"`
}

// Execute - generates concourse pipeline and tasks
func (c *GenerateConcoursePipelineCommand) Execute([]string) error {
	const varsFileName = "vars.yml"
	const pipelineFileName = "pipeline.yml"
	const cfMgmtSh = "cf-mgmt.sh"
	var targetFile string
	if c.TargetDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.TargetDirectory = pwd
	}
	fmt.Println(fmt.Sprintf("Generating pipeline into %s", path.Join(c.TargetDirectory)))

	if err := os.MkdirAll(c.TargetDirectory, 0755); err != nil {
		return errors.Wrapf(err, "Error creating directory %s", c.TargetDirectory)
	}
	if err := os.MkdirAll(c.ConfigDirectory, 0755); err != nil {
		return errors.Wrapf(err, "Error creating directory %s", c.ConfigDirectory)
	}
	if err := c.createFile(pipelineFileName, pipelineFileName); err != nil {
		return errors.Wrap(err, "Error creating pipeline.yml")
	}
	if err := c.createVarsYml(); err != nil {
		return errors.Wrap(err, "Error creating vars.yml")
	}

	if err := os.MkdirAll(path.Join(c.TargetDirectory, "ci", "tasks"), 0755); err == nil {
		lo.G.Debug("Creating", targetFile)
		if err = c.createTaskYml(); err != nil {
			return errors.Wrap(err, "Error creating cf-mgmt.yml")
		}
		targetFile = filepath.Join("ci", "tasks", cfMgmtSh)
		lo.G.Debug("Creating", targetFile)
		if err = c.createFile(cfMgmtSh, targetFile); err != nil {
			return errors.Wrap(err, "Error creating cf-mgmt.sh")
		}
	} else {
		return errors.Wrap(err, "Error making directories")
	}
	fmt.Println(fmt.Sprintf("1) Update %s/vars.yml with the appropriate values", c.ConfigDirectory))
	fmt.Println("2) Using following command to set your pipeline in concourse after you have checked all files in to git")
	fmt.Println(fmt.Sprintf("fly -t <target> set-pipeline -p cf-mgmt -c %s/pipeline.yml --load-vars-from=%s/vars.yml", c.TargetDirectory, c.ConfigDirectory))
	return nil
}

func (c *GenerateConcoursePipelineCommand) createFile(assetName, fileName string) error {
	bytes, err := embedded.Files.ReadFile(fmt.Sprintf("files/%s", assetName))
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	if strings.HasSuffix(fileName, ".sh") {
		perm = 0755
	}
	return os.WriteFile(path.Join(c.TargetDirectory, fileName), bytes, perm)
}

func (c *GenerateConcoursePipelineCommand) createTaskYml() error {
	version := GetVersion()
	bytes, err := embedded.Files.ReadFile("files/cf-mgmt.yml")
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	versioned := strings.Replace(string(bytes), "~VERSION~", version, -1)
	return os.WriteFile(filepath.Join(c.TargetDirectory, "ci", "tasks", "cf-mgmt.yml"), []byte(versioned), perm)
}

func (c *GenerateConcoursePipelineCommand) createVarsYml() error {
	bytes, err := embedded.Files.ReadFile("files/vars-template.yml")
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	processedBytes := strings.Replace(string(bytes), "~CONFIGDIR~", c.ConfigDirectory, -1)
	return os.WriteFile(filepath.Join(c.ConfigDirectory, "vars.yml"), []byte(processedBytes), perm)
}
