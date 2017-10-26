package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/pivotalservices/cf-mgmt/ldap"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <DIRECTORY-OF-CF-MGMT-CONFIG>", os.Args[0])
		return
	}
	argument := os.Args[1]

	_, err := os.Stat(argument)

	if err != nil {
		fmt.Printf("Error: %s is an invalid path\n", argument)
		return
	}

	fmt.Println("Executing Command ", os.Args)

	run(argument)
}

func run(targetDir string) error {

	var walkOutput string
	testDataFileTemplate := `package mock_test_data
	
	import (
		"github.com/pivotalservices/cf-mgmt/config"
		"github.com/pivotalservices/cf-mgmt/ldap"
		mock "github.com/pivotalservices/cf-mgmt/utils/mocks"
	)
	
	func PopulateWithTestData(utilsMgrMock *mock.MockUtilsManager) error {
%s
		return nil
	}`

	err := filepath.Walk(targetDir,
		func(path string, info os.FileInfo, e error) error {

			if info.IsDir() {
				return e
			}
			var item interface{}
			if strings.Contains(path, "ldap.yml") {
				item = &ldap.Config{}
				e = LoadFile(path, item)
			} else if strings.Contains(path, "orgs.yml") {
				item = &config.Orgs{}
				e = LoadFile(path, item)
			} else if strings.Contains(path, "spaces.yml") {
				item = &config.Spaces{}
				e = LoadFile(path, item)
			} else if strings.Contains(path, "orgConfig.yml") {
				item = &config.OrgConfig{}
				e = LoadFile(path, item)
			} else if strings.Contains(path, "security-group.json") {
				var data []byte
				data, e = ioutil.ReadFile(path)
				item = string(data)

			} else if strings.Contains(path, "spaceConfig.yml") || strings.Contains(path, "spaceDefaults.yml") {
				item = &config.SpaceConfig{}
				e = LoadFile(path, item)
			} else {
				e = fmt.Errorf("Unknown File at %s", path)
			}
			if e == nil {
				if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "/") {
					path = "./" + path
				}
				output := fmt.Sprintf("%#v", item)

				if strings.HasPrefix(output, "&") {
					output = strings.TrimPrefix(output, "&")
				}
				walkOutput += fmt.Sprintf("\t\tutilsMgrMock.MockFileData[\"%s\"] = %s\n", path, output)

			}

			return e
		})

	fmt.Printf(testDataFileTemplate, walkOutput)
	return err
}

func LoadFile(configFile string, dataType interface{}) error {
	var data []byte
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, dataType)
}
