package configcommands_test

import (
	"os"
	"path"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	. "github.com/vmwarepivotallabs/cf-mgmt/configcommands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Space", func() {
	var (
		command       *SpaceConfigurationCommand
		configManager config.Manager
		pwd, _        = os.Getwd()
		configDir     = path.Join(pwd, "_testGenSpaces")
	)
	BeforeEach(func() {
		configManager = config.NewManager(configDir)
		err := configManager.CreateConfigIfNotExists("uaa")
		Expect(err).ShouldNot(HaveOccurred())
		orgCommand := &OrgConfigurationCommand{
			OrgName: "test-org",
		}
		orgCommand.ConfigDirectory = configDir
		err = orgCommand.Execute(nil)
		Expect(err).ShouldNot(HaveOccurred())
		command = &SpaceConfigurationCommand{
			OrgName:   "test-org",
			SpaceName: "test-space",
		}
		command.ConfigDirectory = configDir

	})
	AfterEach(func() {
		err := os.RemoveAll(configDir)
		Expect(err).ShouldNot(HaveOccurred())
	})
	Context("Create Space that doesn't exist", func() {
		It("Should Succeed", func() {
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			spaces, err := configManager.OrgSpaces("test-org")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(spaces.Spaces).Should(ConsistOf("test-space"))
		})
	})
	Context("Update Space that does exist", func() {
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
			space, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Labels).Should(HaveKeyWithValue("hello", "world"))
			Expect(space.Metadata.Labels).Should(HaveKeyWithValue("foo", "bar"))
		})

		It("Should Remove Existings Labels", func() {
			command.Metadata = Metadata{
				LabelKey:   []string{"hello", "foo"},
				LabelValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			space, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Labels).Should(HaveKeyWithValue("hello", "world"))
			Expect(space.Metadata.Labels).Should(HaveKeyWithValue("foo", "bar"))

			command.Metadata = Metadata{
				LabelKey:       []string{},
				LabelValue:     []string{},
				LabelsToRemove: []string{"hello", "foo"},
			}

			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			space, err = configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Labels).ShouldNot(HaveKey("hello"))
			Expect(space.Metadata.Labels).ShouldNot(HaveKey("foo"))
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
			space, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Annotations).Should(HaveKeyWithValue("hello", "world"))
			Expect(space.Metadata.Annotations).Should(HaveKeyWithValue("foo", "bar"))
		})

		It("Should remove Annotations", func() {
			command.Metadata = Metadata{
				AnnotationKey:   []string{"hello", "foo"},
				AnnotationValue: []string{"world", "bar"},
			}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			space, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Annotations).Should(HaveKeyWithValue("hello", "world"))
			Expect(space.Metadata.Annotations).Should(HaveKeyWithValue("foo", "bar"))

			command.Metadata = Metadata{
				AnnotationKey:       []string{},
				AnnotationValue:     []string{},
				AnnotationsToRemove: []string{"hello", "foo"},
			}
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())
			space, err = configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(space.Metadata.Annotations).ShouldNot(HaveKey("hello"))
			Expect(space.Metadata.Annotations).ShouldNot(HaveKey("foo"))
		})
	})

	Context("Should add a space developer to a space that has a period in name", func() {
		It("Should Succeed", func() {
			command.Developer.LDAPUsers = []string{"xxx.yyy"}
			err := command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			space, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(space.Developer.LDAPUsers)).To(Equal(1))
			Expect(space.Developer.LDAPUsers).To(ConsistOf([]string{"xxx.yyy"}))
		})
	})

	Context("No options provided", func() {
		It("Nothing should change", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:                        "test-org",
				Space:                      "test-space",
				MemoryLimit:                "unlimited",
				InstanceMemoryLimit:        "unlimited",
				TotalRoutes:                "100",
				TotalServices:              "100",
				PaidServicePlansAllowed:    true,
				TotalReservedRoutePorts:    "unlimited",
				TotalServiceKeys:           "unlimited",
				AppInstanceLimit:           "unlimited",
				AppTaskLimit:               "unlimited",
				EnableSpaceQuota:           true,
				LogRateLimitBytesPerSecond: "unlimited",
			})
			Expect(err).ShouldNot(HaveOccurred())
			spaceBefore, err := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceBefore).Should(BeEquivalentTo(spaceAfter))
		})
	})
	Context("Allow ssh", func() {
		It("allow ssh should stay true when nothing specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: true,
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeTrue())
		})
		It("allow ssh should stay true when true is specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: true,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.AllowSSH = "true"
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeTrue())
		})
		It("allow ssh should change to true when true is specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: false,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.AllowSSH = "true"
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeTrue())
		})
		It("allow ssh should stay false when nothing specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: false,
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeFalse())
		})
		It("allow ssh should stay false when false specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: false,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.AllowSSH = "false"
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeFalse())
		})
		It("allow ssh should change to false when false specified", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:      "test-org",
				Space:    "test-space",
				AllowSSH: true,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.AllowSSH = "false"
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.AllowSSH).Should(BeFalse())
		})
	})

	Context("Named Quota", func() {
		It("if present removes quota elements", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:                     "test-org",
				Space:                   "test-space",
				MemoryLimit:             "unlimited",
				InstanceMemoryLimit:     "unlimited",
				TotalRoutes:             "100",
				TotalServices:           "100",
				PaidServicePlansAllowed: true,
				TotalReservedRoutePorts: "unlimited",
				TotalServiceKeys:        "unlimited",
				AppInstanceLimit:        "unlimited",
				AppTaskLimit:            "unlimited",
				EnableSpaceQuota:        true,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.NamedQuota = "small"
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo("small"))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.MemoryLimit).Should(BeEquivalentTo(""))
			Expect(spaceAfter.InstanceMemoryLimit).Should(BeEquivalentTo(""))
			Expect(spaceAfter.TotalRoutes).Should(BeEquivalentTo(""))
			Expect(spaceAfter.TotalServices).Should(BeEquivalentTo(""))
			Expect(spaceAfter.PaidServicePlansAllowed).Should(BeFalse())
			Expect(spaceAfter.TotalReservedRoutePorts).Should(BeEquivalentTo(""))
			Expect(spaceAfter.TotalServiceKeys).Should(BeEquivalentTo(""))
			Expect(spaceAfter.AppInstanceLimit).Should(BeEquivalentTo(""))
			Expect(spaceAfter.AppTaskLimit).Should(BeEquivalentTo(""))
		})

		It("if cleared named quotes sets back to defaults", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:              "test-org",
				Space:            "test-space",
				NamedQuota:       "small",
				EnableSpaceQuota: false,
			})
			Expect(err).ShouldNot(HaveOccurred())
			command.ClearNamedQuota = true
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo(""))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.MemoryLimit).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.InstanceMemoryLimit).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.TotalRoutes).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.TotalServices).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.PaidServicePlansAllowed).Should(BeFalse())
			Expect(spaceAfter.TotalReservedRoutePorts).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.TotalServiceKeys).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.AppInstanceLimit).Should(BeEquivalentTo("unlimited"))
			Expect(spaceAfter.AppTaskLimit).Should(BeEquivalentTo("unlimited"))
		})
	})

	Context("Space Quota", func() {
		It("if memory limit not set default is set", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:   "test-org",
				Space: "test-space",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo(""))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.MemoryLimit).Should(BeEquivalentTo("unlimited"))
		})
		It("if memory limit is set left alone", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:         "test-org",
				Space:       "test-space",
				MemoryLimit: "100MB",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo(""))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.MemoryLimit).Should(BeEquivalentTo("100MB"))
		})

		It("if total services not set default is set", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:   "test-org",
				Space: "test-space",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo(""))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.TotalServices).Should(BeEquivalentTo("unlimited"))
		})

		It("if total services set don't change", func() {
			err := configManager.SaveSpaceConfig(&config.SpaceConfig{
				Org:           "test-org",
				Space:         "test-space",
				TotalServices: "50",
			})
			Expect(err).ShouldNot(HaveOccurred())
			err = command.Execute(nil)
			Expect(err).ShouldNot(HaveOccurred())

			spaceAfter, errAfter := configManager.GetSpaceConfig("test-org", "test-space")
			Expect(errAfter).ShouldNot(HaveOccurred())
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.NamedQuota).Should(BeEquivalentTo(""))
			Expect(spaceAfter.EnableSpaceQuota).Should(BeFalse())
			Expect(spaceAfter.TotalServices).Should(BeEquivalentTo("50"))
		})
	})
})
