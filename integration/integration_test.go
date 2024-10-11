package integration_test

import (
	"bytes"
	"fmt"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	configDir         = "./fixture"
	exportedConfigDir = "./exported-config"
)

// cf runs the cf CLI with the specified args.
func cf(args ...string) ([]byte, error) {
	cmd := exec.Command("cf", args...)

	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("cf %s: %v", strings.Join(args, " "), err)
	}

	return out, nil
}

var (
	outPath      string
	systemDomain string
	userID       string
	password     string
	clientSecret string
)

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(time.Second * 30)

	systemDomain = os.Getenv("SYSTEM_DOMAIN")
	userID = "admin"
	password = os.Getenv("CF_ADMIN_PASSWORD")
	clientSecret = os.Getenv("ADMIN_CLIENT_SECRET")

	_, err := cf("login", "--skip-ssl-validation", "-a", "https://api."+systemDomain, "-u", userID, "-p", password)
	Expect(err).ShouldNot(HaveOccurred())

	outPath, err = Build("github.com/vmwarepivotallabs/cf-mgmt/cmd/cf-mgmt")
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	CleanupBuildArtifacts()
})

var _ = Describe("cf-mgmt cli", func() {
	Describe("running against pcfdev", func() {
		Describe("orgs, spaces, isolation segments", func() {
			BeforeEach(func() {
				fmt.Println("********   Before called *********")

				cf("delete-org", "-f", "test1")
				cf("delete-org", "-f", "test2")
				cf("delete-org", "-f", "rogue-org1")
				cf("delete-org", "-f", "rogue-org2")
			})

			AfterEach(func() {
				fmt.Println("********   after called *********")
				os.RemoveAll("./config")
				cf("delete-org", "-f", "test1")
				cf("delete-org", "-f", "test2")
				cf("delete-org", "-f", "rogue-org1")
				cf("delete-org", "-f", "rogue-org2")
			})

			It("should complete successfully", func() {
				orgs, err := cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(orgs).ShouldNot(ContainElement("test1"))
				Expect(orgs).ShouldNot(ContainElement("test2"))

				By("creating orgs")
				createOrgsCommand := exec.Command(outPath, "create-orgs",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", userID,
					"--password", password,
					"--client-secret", clientSecret)
				session, err := Start(createOrgsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				orgs, err = cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("test1"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("test2"))).Should(BeTrue())

				By("deleting unused orgs")
				deleteOrgsCommand := exec.Command(outPath, "delete-orgs",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", userID,
					"--password", password,
					"--client-secret", clientSecret)
				session, err = Start(deleteOrgsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				orgs, err = cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("system"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org1"))).ShouldNot(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org1"))).ShouldNot(BeTrue())

				By("creating spaces")
				createSpacesCommand := exec.Command(outPath, "create-spaces",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", userID,
					"--password", password,
					"--client-secret", clientSecret)
				session, err = Start(createSpacesCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				_, err = cf("target", "-o", "test1")
				Expect(err).ShouldNot(HaveOccurred())
				spaces, err := cf("spaces")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaces, []byte("dev"))).Should(BeTrue())
				Expect(bytes.Contains(spaces, []byte("prod"))).Should(BeTrue())

				_, err = cf("target", "-o", "test2")
				Expect(err).ShouldNot(HaveOccurred())
				spaces, err = cf("spaces")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaces, []byte("No spaces found"))).Should(BeTrue())

				By("updating isolation segments")
				updateIsoSegmentsCommand := exec.Command(outPath, "isolation-segments",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", userID,
					"--password", password,
					"--client-secret", clientSecret)
				session, err = Start(updateIsoSegmentsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				is, err := cf("isolation-segments")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(is, []byte("test1-iso-segment"))).Should(BeTrue())
				Expect(bytes.Contains(is, []byte("test2-iso-segment"))).Should(BeTrue())

				// test1-iso-segment should be default for org test1, space dev
				cf("target", "-o", "test1")
				spaceInfo, err := cf("space", "dev")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaceInfo, []byte("test1-iso-segment"))).Should(BeTrue())

				// test2-iso-segment should be default for all of org test2
				orgInfo, err := cf("org", "test2")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgInfo, []byte("test2-iso-segment"))).Should(BeTrue())
			})

			It("should complete successfully without password", func() {
				orgs, err := cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("test1"))).ShouldNot(BeTrue())
				Expect(bytes.Contains(orgs, []byte("test2"))).ShouldNot(BeTrue())

				By("creating orgs")
				createOrgsCommand := exec.Command(outPath, "create-orgs",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", "cf-mgmt",
					"--client-secret", "cf-mgmt-secret")
				session, err := Start(createOrgsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				orgs, err = cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("test1"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("test2"))).Should(BeTrue())

				By("deleting unused orgs")
				deleteOrgsCommand := exec.Command(outPath, "delete-orgs",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", "cf-mgmt",
					"--client-secret", "cf-mgmt-secret")
				session, err = Start(deleteOrgsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				orgs, err = cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("system"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org1"))).ShouldNot(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org2"))).ShouldNot(BeTrue())

				By("creating spaces")
				createSpacesCommand := exec.Command(outPath, "create-spaces",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", "cf-mgmt",
					"--client-secret", "cf-mgmt-secret")
				session, err = Start(createSpacesCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				_, err = cf("target", "-o", "test1")
				Expect(err).ShouldNot(HaveOccurred())
				spaces, err := cf("spaces")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaces, []byte("dev"))).Should(BeTrue())
				Expect(bytes.Contains(spaces, []byte("prod"))).Should(BeTrue())

				_, err = cf("target", "-o", "test2")
				Expect(err).ShouldNot(HaveOccurred())
				spaces, err = cf("spaces")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaces, []byte("No spaces found"))).Should(BeTrue())

				By("updating isolation segments")
				updateIsoSegmentsCommand := exec.Command(outPath, "isolation-segments",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", "cf-mgmt",
					"--client-secret", "cf-mgmt-secret")
				session, err = Start(updateIsoSegmentsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				is, err := cf("isolation-segments")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(is, []byte("test1-iso-segment"))).Should(BeTrue())
				Expect(bytes.Contains(is, []byte("test2-iso-segment"))).Should(BeTrue())

				// test1-iso-segment should be default for org test1, space dev
				cf("target", "-o", "test1")
				spaceInfo, err := cf("space", "dev")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(spaceInfo, []byte("test1-iso-segment"))).Should(BeTrue())

				// test2-iso-segment should be default for all of org test2
				orgInfo, err := cf("org", "test2")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgInfo, []byte("test2-iso-segment"))).Should(BeTrue())
			})
		})

		Describe("export-config", func() {
			BeforeEach(func() {
				fmt.Println("********   Before called *********")

				err := os.MkdirAll(exportedConfigDir, os.ModePerm)
				Expect(err).ShouldNot(HaveOccurred())
				// There seems to be an existing issue where if this ldap.yml does not exist, "cf-mgmt export-config" will fail with:
				// "Unable to initialize cf-mgmt. Error : Error reading file exported-config/ldap.yml: open exported-config/ldap.yml: no such file or directory"
				_, err = os.Create(exportedConfigDir + "/ldap.yml")
				Expect(err).ShouldNot(HaveOccurred())

				cf("create-org", "test-org")
				cf("target", "-o", "test-org")
				cf("create-space", "test-space")
				cf("target", "-o", "test-org", "-s", "test-space")
			})

			AfterEach(func() {
				fmt.Println("********   after called *********")

				err := os.RemoveAll(exportedConfigDir)
				Expect(err).ShouldNot(HaveOccurred())

				cf("delete-space", "test-space", "-o", "test-org", "-f")
				cf("delete-org", "-f", "test-org")
			})

			It("should complete successfully", func() {
				orgs, err := cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(orgs)).Should(ContainSubstring("test-org"))

				By("exporting space roles")
				exportConfigCommand := exec.Command(outPath, "export-config",
					"--config-dir", exportedConfigDir,
					"--system-domain", systemDomain,
					"--user-id", userID,
					"--password", password,
					"--client-secret", clientSecret)
				session, err := Start(exportConfigCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				spaceConfigFile, err := os.ReadFile(exportedConfigDir + "/test-org/test-space/spaceConfig.yml")
				Expect(err).ShouldNot(HaveOccurred())
				var spaceConfig config.SpaceConfig
				yaml.Unmarshal(spaceConfigFile, &spaceConfig)

				// the user that created the space is auto-assigned as a space developer; the exported config should reflect that
				Expect(spaceConfig.Developer.Users).Should(ContainElement("admin"))
			})
		})
	})
})
