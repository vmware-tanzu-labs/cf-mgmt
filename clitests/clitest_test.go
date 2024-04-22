package clitests_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("cf-mgmt cli", func() {
	Describe("running cli commands to create config files", func() {
		var (
			outPath string
			err     error
		)
		BeforeEach(func() {
			os.RemoveAll("./config")
			outPath, err = Build("github.com/vmwarepivotallabs/cf-mgmt/cmd/cf-mgmt-config")
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll("./config")
			CleanupBuildArtifacts()
		})
		It("should complete successfully", func() {

			initConfigCommand := exec.Command(outPath, "init")
			session, err := Start(initConfigCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			_, err = os.Stat("./config/ldap.yml")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat("./config/orgs.yml")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat("./config/spaceDefaults.yml")
			Expect(err).ShouldNot(HaveOccurred())

			addOrgToConfigCommand := exec.Command(outPath, "add-org", "--org", "test-org")
			session, err = Start(addOrgToConfigCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			_, err = os.Stat("./config/test-org/orgConfig.yml")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat("./config/test-org/spaces.yml")
			Expect(err).ShouldNot(HaveOccurred())

			addSpaceToConfigCommand := exec.Command(outPath, "add-space", "--org", "test-org", "--space", "test-space")
			session, err = Start(addSpaceToConfigCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			_, err = os.Stat("./config/test-org/test-space/security-group.json")
			Expect(err).ShouldNot(HaveOccurred())
			_, err = os.Stat("./config/test-org/test-space/spaceConfig.yml")
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Describe("running cli commands to create concourse pipeline", func() {
		var (
			outPath string
			err     error
		)
		BeforeEach(func() {
			outPath, err = Build("github.com/vmwarepivotallabs/cf-mgmt/cmd/cf-mgmt-config")
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll("./ci")
			os.Remove("pipeline.yml")
			os.Remove("vars.yml")
			CleanupBuildArtifacts()
		})
		It("should complete successfully", func() {

			initConfigCommand := exec.Command(outPath, "generate-concourse-pipeline")
			session, err := Start(initConfigCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			_, err = os.Stat("pipeline.yml")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = os.Stat("config/vars.yml")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = os.Stat("./ci/tasks/cf-mgmt.sh")
			Expect(err).ShouldNot(HaveOccurred())

			_, err = os.Stat("./ci/tasks/cf-mgmt.yml")
			Expect(err).ShouldNot(HaveOccurred())

		})
	})
})
