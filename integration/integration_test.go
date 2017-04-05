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
			Ω(err).ShouldNot(HaveOccurred())
		})
		AfterEach(func() {
			os.RemoveAll("./config")
			ccManager.DeleteOrg("test1")
			ccManager.DeleteOrg("test2")
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

			/*quotas, err := ccManager.ListQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			_, ok := quotas["test1"]
			Ω(ok).Should(BeFalse())

			_, ok = quotas["test2"]
			Ω(ok).Should(BeFalse())

			updateOrgQuotasCommand := exec.Command(outPath, "update-org-quotas", "--config-dir", configDir,
				"--system-domain", systemDomain, "--user-id", userId, "--password",
				password, "--client-secret", clientSecret)
			session, err = Start(updateOrgQuotasCommand, GinkgoWriter, GinkgoWriter)
			Expect(err).ShouldNot(HaveOccurred())
			Eventually(session).Should(Exit(0))

			quotas, err = ccManager.ListQuotas()
			Expect(err).ShouldNot(HaveOccurred())
			_, ok = quotas["test1"]
			Ω(ok).Should(BeTrue())

			_, ok = quotas["test2"]
			Ω(ok).Should(BeFalse())*/

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
