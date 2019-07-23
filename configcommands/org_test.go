package configcommands_test

import (
	"os"
	"path"

	"github.com/pivotalservices/cf-mgmt/config"
	. "github.com/pivotalservices/cf-mgmt/configcommands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Org", func() {
	var (
		command   *OrgConfigurationCommand
		pwd, _    = os.Getwd()
		configDir = path.Join(pwd, "_testGen")
	)
	BeforeEach(func() {
		configManager := config.NewManager(configDir)
		err := configManager.CreateConfigIfNotExists("uaa")
		Expect(err).ShouldNot(HaveOccurred())
		command = &OrgConfigurationCommand{
			OrgName: "test",
		}
		command.ConfigDirectory = configDir
	})
	AfterEach(func() {
		err := os.RemoveAll(configDir)
		Expect(err).ShouldNot(HaveOccurred())
	})
	Context("Create Org that doesn't exist", func() {
		It("Should Succeed", func() {
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
	Context("Update Org that does exist", func() {
		It("Should Succeed", func() {
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
