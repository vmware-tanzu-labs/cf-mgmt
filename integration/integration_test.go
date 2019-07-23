package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

const (
	systemDomain = "dev.cfdev.sh"
	userId       = "admin"
	password     = "admin"
	clientSecret = "admin-client-secret"
	configDir    = "./fixture"
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

var outPath string

var _ = BeforeSuite(func() {
	_, err := cf("login", "--skip-ssl-validation", "-a", "https://api."+systemDomain, "-u", userId, "-p", password)
	Expect(err).ShouldNot(HaveOccurred())

	outPath, err = Build("github.com/pivotalservices/cf-mgmt/cmd/cf-mgmt")
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
					"--user-id", userId,
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
					"--user-id", userId,
					"--password", password,
					"--client-secret", clientSecret)
				session, err = Start(deleteOrgsCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				orgs, err = cf("orgs")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(bytes.Contains(orgs, []byte("system"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("cfdev-org"))).Should(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org1"))).ShouldNot(BeTrue())
				Expect(bytes.Contains(orgs, []byte("rogue-org1"))).ShouldNot(BeTrue())

				By("creating spaces")
				createSpacesCommand := exec.Command(outPath, "create-spaces",
					"--config-dir", configDir,
					"--system-domain", systemDomain,
					"--user-id", userId,
					"--password", password,
					"--client-secret", clientSecret)
				session, err = Start(createSpacesCommand, GinkgoWriter, GinkgoWriter)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

				_, err = cf("target", "-o", "test1")
				Expect(err).ShouldNot(HaveOccurred())
				spaces, err := cf("spaces")
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
					"--user-id", userId,
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
				Expect(bytes.Contains(orgs, []byte("cfdev-org"))).Should(BeTrue())
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

		Describe("export config", func() {
			/*It("should export config with > 50 spaces", func() {
				cfToken, err := uaa.GetCFToken(systemDomain, userId, password)
				Expect(err).ShouldNot(HaveOccurred())
				ccManager, err := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", systemDomain), cfToken, "1.0", false)
				Expect(err).ShouldNot(HaveOccurred())
				ccManager.CreateOrg("test1")
				orgs, _ := ccManager.ListOrgs()
				for _, org := range orgs {
					if org.Name == "test1" {
						for i := 1; i < 101; i++ {
							ccManager.CreateSpace(fmt.Sprintf("space-%d", i), org.Guid)
						}
					}
				}

				exportConfigCommand := exec.Command(outPath, "export-config",
					"--config-dir", "./config",
					"--system-domain", systemDomain,
					"--user-id", userId,
					"--password", password,
					"--client-secret", clientSecret)
				session, err := Start(exportConfigCommand, GinkgoWriter, GinkgoWriter)
				session.Wait(20)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))
			})

			It("should export config with > 50 spaces without a password", func() {
				cfToken, err := uaa.GetCFToken(systemDomain, userId, password)
				Expect(err).ShouldNot(HaveOccurred())
				ccManager, err := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", systemDomain), cfToken, "1.0", false)
				Expect(err).ShouldNot(HaveOccurred())
				ccManager.CreateOrg("test1")
				orgs, _ := ccManager.ListOrgs()
				for _, org := range orgs {
					if org.Name == "test1" {
						i := 1
						for i < 101 {
							ccManager.CreateSpace(fmt.Sprintf("space-%d", i), org.Guid)
							i++
						}
					}
				}

				exportConfigCommand := exec.Command(outPath, "export-config", "--config-dir", "./config",
					"--system-domain", systemDomain, "--user-id", "cf-mgmt", "--client-secret", "cf-mgmt-secret")
				session, err := Start(exportConfigCommand, GinkgoWriter, GinkgoWriter)
				session.Wait(20)
				Expect(err).ShouldNot(HaveOccurred())
				Eventually(session).Should(Exit(0))

			})*/

		})
	})
})
