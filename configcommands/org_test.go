package configcommands_test

import (
	"os"
	"path"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	. "github.com/vmwarepivotallabs/cf-mgmt/configcommands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Org", func() {
	var (
		configManager config.Manager
		command       *OrgConfigurationCommand
		pwd, _        = os.Getwd()
		configDir     = path.Join(pwd, "_testGenOrgs")
	)
	BeforeEach(func() {
		configManager = config.NewManager(configDir)
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
			orgs, err := configManager.Orgs()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(orgs.Orgs).Should(ConsistOf("test"))
		})
	})
	Context("Update Org that does exist", func() {
		It("Should Succeed", func() {
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("Allow setting metadata labels", func() {
		It("Should Add Labels", func() {
			command.Metadata = Metadata{
				LabelKey:   []string{"hello", "foo"},
				LabelValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err := configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Labels).Should(HaveKeyWithValue("hello", "world"))
			Expect(org.Metadata.Labels).Should(HaveKeyWithValue("foo", "bar"))
		})

		It("Should Remove Existings Labels", func() {
			command.Metadata = Metadata{
				LabelKey:   []string{"hello", "foo"},
				LabelValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err := configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Labels).Should(HaveKeyWithValue("hello", "world"))
			Expect(org.Metadata.Labels).Should(HaveKeyWithValue("foo", "bar"))

			command.Metadata = Metadata{
				LabelKey:       []string{},
				LabelValue:     []string{},
				LabelsToRemove: []string{"hello", "foo"},
			}

			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err = configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Labels).ShouldNot(HaveKey("hello"))
			Expect(org.Metadata.Labels).ShouldNot(HaveKey("foo"))
		})
	})

	Context("Allow setting metadata annotations", func() {
		It("Should Add Annotations", func() {
			command.Metadata = Metadata{
				AnnotationKey:   []string{"hello", "foo"},
				AnnotationValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err := configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Annotations).Should(HaveKeyWithValue("hello", "world"))
			Expect(org.Metadata.Annotations).Should(HaveKeyWithValue("foo", "bar"))
		})

		It("Should remove Annotations", func() {
			command.Metadata = Metadata{
				AnnotationKey:   []string{"hello", "foo"},
				AnnotationValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err := configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Annotations).Should(HaveKeyWithValue("hello", "world"))
			Expect(org.Metadata.Annotations).Should(HaveKeyWithValue("foo", "bar"))

			command.Metadata = Metadata{
				AnnotationKey:       []string{},
				AnnotationValue:     []string{},
				AnnotationsToRemove: []string{"hello", "foo"},
			}
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			org, err = configManager.GetOrgConfig("test")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(org.Metadata.Annotations).ShouldNot(HaveKey("hello"))
			Expect(org.Metadata.Annotations).ShouldNot(HaveKey("foo"))
		})
	})
})
