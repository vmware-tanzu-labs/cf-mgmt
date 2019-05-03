package configcommands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotalservices/cf-mgmt/generated"
	"github.com/xchapter7x/lo"
)

type GenerateConcoursePipelineCommand struct {
}

//Execute - generates concourse pipeline and tasks
func (c *GenerateConcoursePipelineCommand) Execute([]string) error {
	const varsFileName = "vars.yml"
	const pipelineFileName = "pipeline.yml"
	const cfMgmtSh = "cf-mgmt.sh"
	var targetFile string
	fmt.Println("Generating pipeline....")
	if err := createFile(pipelineFileName, pipelineFileName); err != nil {
		lo.G.Error("Error creating pipeline.yml", err)
		return err
	}
	if err := createFile(varsFileName, varsFileName); err != nil {
		lo.G.Error("Error creating vars.yml", err)
		return err
	}

	if err := os.MkdirAll("ci/tasks", 0755); err == nil {
		lo.G.Debug("Creating", targetFile)
		if err = createTaskYml(); err != nil {
			lo.G.Error("Error creating cf-mgmt.yml", err)
			return err
		}
		targetFile = filepath.Join("ci", "tasks", cfMgmtSh)
		lo.G.Debug("Creating", targetFile)
		if err = createFile(cfMgmtSh, targetFile); err != nil {
			lo.G.Error("Error creating cf-mgmt.sh", err)
			return err
		}
	} else {
		lo.G.Error("Error making directories", err)
		return err
	}
	fmt.Println("1) Update vars.yml with the appropriate values")
	fmt.Println("2) Using following command to set your pipeline in concourse after you have checked all files in to git")
	fmt.Println("fly -t lite set-pipeline -p cf-mgmt -c pipeline.yml --load-vars-from=vars.yml")
	return nil
}

func createFile(assetName, fileName string) error {
	bytes, err := generated.Asset(fmt.Sprintf("files/%s", assetName))
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	if strings.HasSuffix(fileName, ".sh") {
		perm = 0755
	}
	return ioutil.WriteFile(fileName, bytes, perm)
}

func createTaskYml() error {
	version := GetVersion()
	bytes, err := generated.Asset("files/cf-mgmt.yml")
	if err != nil {
		return err
	}
	perm := os.FileMode(0666)
	versioned := strings.Replace(string(bytes), "~VERSION~", version, -1)
	return ioutil.WriteFile(filepath.Join("ci", "tasks", "cf-mgmt.yml"), []byte(versioned), perm)
}
