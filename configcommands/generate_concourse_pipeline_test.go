package configcommands_test

import (
	"path"

	. "github.com/vmwarepivotallabs/cf-mgmt/configcommands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GenerateConcoursePipeline", func() {
	var (
		command         *GenerateConcoursePipelineCommand
		targetDirectory = "./_testGen"
		configDirectory = "./_testGen/config"
	)
	BeforeEach(func() {
		command = &GenerateConcoursePipelineCommand{}
		command.TargetDirectory = targetDirectory
		command.ConfigDirectory = configDirectory
	})
	AfterEach(func() {
		// err := os.RemoveAll(configDirectory)
		// Expect(err).ShouldNot(HaveOccurred())
		// err = os.RemoveAll(targetDirectory)
		// Expect(err).ShouldNot(HaveOccurred())
	})
	Context("Should generate a concourse pipeline and tasks", func() {
		It("Should Succeed", func() {
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(path.Join(targetDirectory, "ci", "tasks", "cf-mgmt.sh")).To(BeAnExistingFile())
			Expect(path.Join(targetDirectory, "ci", "tasks", "cf-mgmt.yml")).To(BeAnExistingFile())
			Expect(path.Join(targetDirectory, "pipeline.yml")).To(BeAnExistingFile())
			Expect(path.Join(configDirectory, "vars.yml")).To(BeAnExistingFile())
		})
	})
})
