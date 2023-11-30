package configcommands_test

import (
	"os"
	"path"

	"github.com/vmwarepivotallabs/cf-mgmt/config"
	. "github.com/vmwarepivotallabs/cf-mgmt/configcommands"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Global", func() {
	Context("Integration", func() {
		var (
			configManager config.Manager
			command       *GlobalConfigurationCommand
			pwd, _        = os.Getwd()
			configDir     = path.Join(pwd, "_testGen")
		)
		BeforeEach(func() {
			configManager = config.NewManager(configDir)
			err := configManager.CreateConfigIfNotExists("uaa")
			Expect(err).ShouldNot(HaveOccurred())
			command = &GlobalConfigurationCommand{}
			command.ConfigDirectory = configDir
		})
		AfterEach(func() {
			err := os.RemoveAll(configDir)
			Expect(err).ShouldNot(HaveOccurred())
		})
		Context("EnableDeleteIsolationSegments", func() {
			It("Should be true", func() {
				command.EnableDeleteIsolationSegments = "true"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableDeleteIsolationSegments).To(BeTrue())
			})
			It("Should be false", func() {
				command.EnableDeleteIsolationSegments = "false"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableDeleteIsolationSegments).To(BeFalse())
			})
		})
		Context("EnableDeleteSharedDomains", func() {
			It("Should be true", func() {
				command.EnableDeleteSharedDomains = "true"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableDeleteSharedDomains).To(BeTrue())
			})
			It("Should be false", func() {
				command.EnableDeleteSharedDomains = "false"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableDeleteSharedDomains).To(BeFalse())
			})
		})
		Context("EnableServiceAccess", func() {
			It("Should be true", func() {
				command.EnableServiceAccess = "true"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableServiceAccess).To(BeTrue())
			})
			It("Should be false", func() {
				command.EnableServiceAccess = "false"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableServiceAccess).To(BeFalse())
			})
		})

		Context("EnableUnassignSecurityGroups", func() {
			It("Should be true", func() {
				command.EnableUnassignSecurityGroups = "true"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableUnassignSecurityGroups).To(BeTrue())
			})
			It("Should be false", func() {
				command.EnableUnassignSecurityGroups = "false"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.EnableUnassignSecurityGroups).To(BeFalse())
			})
		})

		Context("MetadataPrefix", func() {
			It("Should be unset", func() {
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.MetadataPrefix).To(Equal("cf-mgmt.pivotal.io"))
			})
			It("Should be changed", func() {
				command.MetadataPrefix = "foo.bar"
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(globalConfig.MetadataPrefix).To(Equal("foo.bar"))
			})
		})

		Context("Staging Security Groups", func() {
			It("Should be unset", func() {
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.StagingSecurityGroups)).To(Equal(0))
			})
			It("Should add 2 staging sec groups", func() {
				command.StagingSecurityGroups = []string{"foo", "bar"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.StagingSecurityGroups)).To(Equal(2))
			})
			It("Should add 2 staging sec groups, and remove 1", func() {
				command.StagingSecurityGroups = []string{"foo", "bar"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.StagingSecurityGroups)).To(Equal(2))

				command.StagingSecurityGroups = []string{}
				command.RemoveStagingSecurityGroups = []string{"bar"}
				err = command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err = configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.StagingSecurityGroups)).To(Equal(1))
			})
		})

		Context("Running Security Groups", func() {
			It("Should be unset", func() {
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.RunningSecurityGroups)).To(Equal(0))
			})
			It("Should add 2 running sec groups", func() {
				command.RunningSecurityGroups = []string{"foo", "bar"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.RunningSecurityGroups)).To(Equal(2))
			})
			It("Should add 2 running sec groups, and remove 1", func() {
				command.RunningSecurityGroups = []string{"foo", "bar"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.RunningSecurityGroups)).To(Equal(2))

				command.RunningSecurityGroups = []string{}
				command.RemoveRunningSecurityGroups = []string{"bar"}
				err = command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err = configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.RunningSecurityGroups)).To(Equal(1))
			})
		})

		Context("Shared Domains", func() {
			It("Should be unset", func() {
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(0))
			})
			It("Should add 2 internal shared domains", func() {
				command.InternalSharedDomains = []string{"foo.io", "bar.io"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(2))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("foo.io", config.SharedDomain{Internal: true}))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("bar.io", config.SharedDomain{Internal: true}))
			})
			It("Should add 2 shared domains", func() {
				command.SharedDomains = []string{"foo.io", "bar.io"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(2))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("foo.io", config.SharedDomain{Internal: false}))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("bar.io", config.SharedDomain{Internal: false}))
			})

			It("Should add 2 router group shared domains", func() {
				command.RouterGroupSharedDomains = []string{"foo.io", "bar.io"}
				command.RouterGroupSharedDomainsGroups = []string{"grp1", "grp2"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(2))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("foo.io", config.SharedDomain{Internal: false, RouterGroup: "grp1"}))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("bar.io", config.SharedDomain{Internal: false, RouterGroup: "grp2"}))
			})

			It("Should add 2 shared domains and remove 1", func() {
				command.SharedDomains = []string{"foo.io", "bar.io"}
				err := command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err := configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(2))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("foo.io", config.SharedDomain{Internal: false}))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("bar.io", config.SharedDomain{Internal: false}))

				command.SharedDomains = []string{}
				command.RemoveSharedDomains = []string{"foo.io"}
				err = command.Execute(nil)
				Expect(err).ShouldNot(HaveOccurred())
				globalConfig, err = configManager.GetGlobalConfig()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(globalConfig.SharedDomains)).To(Equal(1))
				Expect(globalConfig.SharedDomains).To(HaveKeyWithValue("bar.io", config.SharedDomain{Internal: false}))
			})
		})
	})
	Context("ServiceAccess", func() {
		var (
			command      *GlobalConfigurationCommand
			globalConfig *config.GlobalConfig
		)
		BeforeEach(func() {
			command = &GlobalConfigurationCommand{}
			globalConfig = &config.GlobalConfig{}
		})
		It("Will make plan public", func() {
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.AllAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
		})

		It("Will not update anything as plan is already public", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name:           "test-service",
							AllAccessPlans: []string{"test-plan"},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.AllAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
		})

		It("Will make plan public that is current no-access", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name:          "test-service",
							NoAccessPlans: []string{"test-plan"},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.AllAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
		})

		It("Will make plan public that is current limited access", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name: "test-service",
							LimitedAccessPlans: []*config.PlanVisibility{
								{
									Name: "test-plan",
								},
							},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.AllAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
		})

		It("Will make plan no-access", func() {
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.NoAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will not update anything as plan is already no access", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name:          "test-service",
							NoAccessPlans: []string{"test-plan"},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.NoAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will make plan no access that is current limited access", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name: "test-service",
							LimitedAccessPlans: []*config.PlanVisibility{
								{
									Name: "test-plan",
								},
							},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.NoAccessPlan = "test-plan"
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans).Should(ConsistOf([]string{"test-plan"}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will make plan limited-access", func() {
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.LimitedAccessPlan = "test-plan"
			command.ServiceAccess.OrgsToAdd = []string{"org1", "org2"}
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans).Should(ConsistOf([]*config.PlanVisibility{
				{
					Name: "test-plan", Orgs: []string{"org1", "org2"},
				},
			}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will not update anything as plan is already limited", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name: "test-service",
							LimitedAccessPlans: []*config.PlanVisibility{
								{
									Name: "test-plan", Orgs: []string{"org1", "org2"},
								},
							},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.LimitedAccessPlan = "test-plan"
			command.ServiceAccess.OrgsToAdd = []string{"org1", "org2"}
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans).Should(ConsistOf([]*config.PlanVisibility{
				{
					Name: "test-plan", Orgs: []string{"org1", "org2"},
				},
			}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will make plan limited that is current public access", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name:           "test-service",
							AllAccessPlans: []string{"test-plan"},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.LimitedAccessPlan = "test-plan"
			command.ServiceAccess.OrgsToAdd = []string{"org1", "org2"}
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans).Should(ConsistOf([]*config.PlanVisibility{
				{
					Name: "test-plan", Orgs: []string{"org1", "org2"},
				},
			}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
		})

		It("Will add an org to already limited plan", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name: "test-service",
							LimitedAccessPlans: []*config.PlanVisibility{
								{
									Name: "test-plan", Orgs: []string{"org1", "org2"},
								},
							},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.LimitedAccessPlan = "test-plan"
			command.ServiceAccess.OrgsToAdd = []string{"org3"}
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans).Should(ConsistOf([]*config.PlanVisibility{
				{
					Name: "test-plan", Orgs: []string{"org1", "org2", "org3"},
				},
			}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})

		It("Will add an org and remove an org to already limited plan", func() {
			globalConfig.ServiceAccess = []*config.Broker{
				{
					Name: "test-broker",
					Services: []*config.Service{
						{
							Name: "test-service",
							LimitedAccessPlans: []*config.PlanVisibility{
								{
									Name: "test-plan", Orgs: []string{"org1", "org2"},
								},
							},
						},
					},
				},
			}
			command.ServiceAccess.Broker = "test-broker"
			command.ServiceAccess.Service = "test-service"
			command.ServiceAccess.LimitedAccessPlan = "test-plan"
			command.ServiceAccess.OrgsToAdd = []string{"org3"}
			command.ServiceAccess.OrgsToRemove = []string{"org2"}
			errList := command.UpdateServiceAccess(globalConfig)
			Expect(len(errList)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Name).Should(Equal("test-broker"))
			Expect(len(globalConfig.ServiceAccess[0].Services)).Should(Equal(1))
			Expect(globalConfig.ServiceAccess[0].Services[0].Name).Should(Equal("test-service"))
			Expect(globalConfig.ServiceAccess[0].Services[0].LimitedAccessPlans).Should(ConsistOf([]*config.PlanVisibility{
				{
					Name: "test-plan", Orgs: []string{"org1", "org3"},
				},
			}))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].NoAccessPlans)).Should(Equal(0))
			Expect(len(globalConfig.ServiceAccess[0].Services[0].AllAccessPlans)).Should(Equal(0))
		})
	})
})
