package integration_test

import (
	"fmt"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
	"github.com/pivotalservices/cf-mgmt/cloudcontroller"
	"github.com/pivotalservices/cf-mgmt/uaa"
)

const (
	systemDomain = "local.pcfdev.io"
	userId       = "admin"
	password     = "admin"
	clientSecret = "admin-client-secret"
	configDir    = "./fixture"
)

var _ = Describe("cf-mgmt cli", func() {
	Describe("running tests against pcfdev", func() {

		var (
			outPath   string
			err       error
			ccManager cloudcontroller.Manager
			cfToken   string
		)
		BeforeEach(func() {
			uaaManager := uaa.NewDefaultUAAManager(systemDomain, userId)
			cfToken, err = uaaManager.GetCFToken(password)
			Expect(err).ShouldNot(HaveOccurred())

			ccManager = cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", systemDomain), cfToken)
			outPath, err = Build("github.com/pivotalservices/cf-mgmt")

			ccManager.CreateOrg("rogue-org1")
			ccManager.CreateOrg("rogue-org2")

			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll("./config")
			ccManager.DeleteOrg("test1")
			ccManager.DeleteOrg("test2")
			ccManager.DeleteOrg("rogue-org1")
			ccManager.DeleteOrg("rogue-org2")
			CleanupBuildArtifacts()
		})
		It("should complete successfully", func() {
			var orgs []*cloudcontroller.Org
			var spaces []*cloudcontroller.Space
			orgs, err = ccManager.ListOrgs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getOrg(orgs, "test1")).Should(BeNil())
			Expect(getOrg(orgs, "test2")).Should(BeNil())

			createOrgsCommand := exec.Command(outPath, "create-orgs", "--config-dir", configDir,
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err := Start(createOrgsCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			orgs, err = ccManager.ListOrgs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getOrg(orgs, "test1")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "test2")).ShouldNot(BeNil())

			deleteOrgsCommand := exec.Command(outPath, "delete-orgs", "--config-dir", configDir,
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err = Start(deleteOrgsCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			orgs, err = ccManager.ListOrgs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getOrg(orgs, "system")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "pcfdev-org")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "rogue-org1")).Should(BeNil())
			Expect(getOrg(orgs, "rogue-org2")).Should(BeNil())

			ccManager.CreateOrg("rogue-org1")
			ccManager.CreateOrg("rogue-org2")
			peekDeleteOrgsCommand := exec.Command(outPath, "delete-orgs", "--peek", "--config-dir", configDir,
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err = Start(peekDeleteOrgsCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			orgs, err = ccManager.ListOrgs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getOrg(orgs, "system")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "pcfdev-org")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "rogue-org1")).ShouldNot(BeNil())
			Expect(getOrg(orgs, "rogue-org2")).ShouldNot(BeNil())

			createSpacesCommand := exec.Command(outPath, "create-spaces", "--config-dir", configDir,
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err = Start(createSpacesCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			spaces, err = ccManager.ListSpaces(getOrg(orgs, "test1").MetaData.GUID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(getSpace(spaces, "dev")).ShouldNot(BeNil())
			Expect(getSpace(spaces, "prod")).ShouldNot(BeNil())

			spaces, err = ccManager.ListSpaces(getOrg(orgs, "test2").MetaData.GUID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(spaces)).Should(BeEquivalentTo(0))
		})

		It("should export config with > 50 spaces", func() {

			ccManager := cloudcontroller.NewManager(fmt.Sprintf("https://api.%s", systemDomain), cfToken)

			ccManager.CreateOrg("test1")
			orgs, _ := ccManager.ListOrgs()
			for _, org := range orgs {
				if org.Entity.Name == "test1" {
					i := 1
					for i < 101 {
						ccManager.CreateSpace(fmt.Sprintf("space-%d", i), org.MetaData.GUID)
						i++
					}
				}
			}

			exportConfigCommand := exec.Command(outPath, "export-config", "--config-dir", "./config",
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err := Start(exportConfigCommand, GinkgoWriter, GinkgoWriter)
			session.Wait(20)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

		})

	})
})

func getSpace(spaces []*cloudcontroller.Space, spaceName string) *cloudcontroller.Space {
	for _, space := range spaces {
		if space.Entity.Name == spaceName {
			return space
		}
	}
	return nil
}

func getOrg(orgs []*cloudcontroller.Org, orgName string) *cloudcontroller.Org {
	for _, org := range orgs {
		if org.Entity.Name == orgName {
			return org
		}
	}
	return nil
}
