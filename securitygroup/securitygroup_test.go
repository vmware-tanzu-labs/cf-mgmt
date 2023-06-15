package securitygroup_test

import (
	"encoding/json"
	"errors"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	configfakes "github.com/vmwarepivotallabs/cf-mgmt/config/fakes"
	"github.com/vmwarepivotallabs/cf-mgmt/securitygroup"
	securitygroupfakes "github.com/vmwarepivotallabs/cf-mgmt/securitygroup/fakes"
	spacefakes "github.com/vmwarepivotallabs/cf-mgmt/space/fakes"
)

const asg_config = `[
  {
    "protocol": "icmp",
    "destination": "0.0.0.0/0",
    "type": 0,
    "code": 1
  },
  {
    "protocol": "tcp",
    "destination": "10.0.11.0/24",
    "ports": "80,443",
    "log": true,
    "description": "Allow http and https traffic from ZoneA"
  }
]`

const asg_config_with_trailing_space = `[
  {
    "protocol": "icmp",
    "destination": "0.0.0.0/0 ",
    "type": 0,
    "code": 1
  }
]`

const asg_config_with_space = `[
  {
    "protocol": "tcp",
    "destination": "10.0.11.0 - 10.0.11.2",
    "ports": "80,443",
    "log": true,
    "description": "Allow http and https traffic from ZoneA"
  }
]`

var _ = Describe("given Security Group Manager", func() {
	var (
		fakeReader   *configfakes.FakeReader
		fakeClient   *securitygroupfakes.FakeCFClient
		fakeSpaceMgr *spacefakes.FakeManager
		securityMgr  securitygroup.DefaultManager
	)

	BeforeEach(func() {
		fakeReader = new(configfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		fakeClient = new(securitygroupfakes.FakeCFClient)
		securityMgr = securitygroup.DefaultManager{
			Cfg:          fakeReader,
			Client:       fakeClient,
			SpaceManager: fakeSpaceMgr,
			Peek:         false,
		}
	})
	Context("ListNonDefaultSecurityGroups", func() {
		It("returns 2 security groups", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "group1",
					Guid:    "group1-guid",
					Running: false,
					Staging: false,
				},
				{
					Name:    "group2",
					Guid:    "group2-guid",
					Running: false,
					Staging: false,
				},
			}, nil)
			groups, err := securityMgr.ListNonDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(groups)).Should(Equal(2))
		})

		It("returns error", func() {
			fakeClient.ListSecGroupsReturns(nil, errors.New("error"))
			_, err := securityMgr.ListNonDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.ListSecGroupsCallCount()).Should(Equal(1))
		})
	})

	Context("ListDefaultSecurityGroups", func() {
		It("returns 2 security groups", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "group1",
					Guid:    "group1-guid",
					Running: true,
					Staging: false,
				},
				{
					Name:    "group2",
					Guid:    "group2-guid",
					Running: false,
					Staging: true,
				},
			}, nil)
			groups, err := securityMgr.ListDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(groups)).Should(Equal(2))
		})

		It("returns error", func() {
			fakeClient.ListSecGroupsReturns(nil, errors.New("error"))
			_, err := securityMgr.ListDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.ListSecGroupsCallCount()).Should(Equal(1))
		})
	})

	Context("CreateApplicationSecurityGroups", func() {
		BeforeEach(func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup:   true,
					Space:                 "space1",
					Org:                   "org1",
					SecurityGroupContents: asg_config,
				},
				{
					EnableSecurityGroup: false,
					Space:               "space2",
					Org:                 "org1",
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeSpaceMgr.FindSpaceReturns(cfclient.Space{
				Name:             "space1",
				Guid:             "space1-guid",
				OrganizationGuid: "org1-guid",
			}, nil)
		})

		It("Should assign global group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(1))
			sgGUID, spaceGUID := fakeClient.BindSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("dns-guid"))
			Expect(spaceGUID).Should(Equal("space1-guid"))
		})

		It("Should not assign global group to space that is already assigned", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
			}, nil)
			fakeClient.ListSpaceSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListSpaceSecGroupsCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should unbind global group to space that is not in configuration", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup:         false,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
					SpacesData: []cfclient.SpaceResource{
						{
							Entity: cfclient.Space{
								Guid: "space1-guid",
							},
						},
					},
				},
				{
					Name: "ntp",
					Guid: "ntp-guid",
					SpacesData: []cfclient.SpaceResource{
						{
							Entity: cfclient.Space{
								Guid: "space1-guid",
							},
						},
					},
				},
			}, nil)
			fakeClient.ListSpaceSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
				{
					Name: "ntp",
					Guid: "ntp-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListSpaceSecGroupsCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnbindSecGroupCallCount()).Should(Equal(1))
			secGroupGUID, spaceGUID := fakeClient.UnbindSecGroupArgsForCall(0)
			Expect(secGroupGUID).Should(Equal("ntp-guid"))
			Expect(spaceGUID).Should(Equal("space1-guid"))
		})

		It("Should not unbind space specific group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup:         true,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
				{
					Name:  "org1-space1",
					Guid:  "org1-space1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeClient.ListSpaceSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
				{
					Name:  "org1-space1",
					Guid:  "org1-space1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListSpaceSecGroupsCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnbindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should unbind space specific group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup:         false,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
					SpacesData: []cfclient.SpaceResource{
						{
							Entity: cfclient.Space{
								Guid: "space1-guid",
							},
						},
					},
				},
				{
					Name:  "org1-space1",
					Guid:  "org1-space1-guid",
					Rules: []cfclient.SecGroupRule{},
					SpacesData: []cfclient.SpaceResource{
						{
							Entity: cfclient.Space{
								Guid: "space1-guid",
							},
						},
					},
				},
			}, nil)
			fakeClient.ListSpaceSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
				{
					Name:  "org1-space1",
					Guid:  "org1-space1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListSpaceSecGroupsCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnbindSecGroupCallCount()).Should(Equal(1))
		})

		It("Should error assigning global group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name: "dns",
					Guid: "dns-guid",
				},
			}, nil)
			fakeClient.BindSecGroupReturns(errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(1))
		})

		It("Should error when group doesn't exist", func() {
			spaceConfigs := []config.SpaceConfig{
				{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("Security group [dns] does not exist"))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should create and assign group to space", func() {
			fakeClient.ListSecGroupsReturns(nil, nil)
			fakeClient.CreateSecGroupReturns(&cfclient.SecGroup{Name: "org1-space1", Guid: "org1-space1-guid"}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(1))
			sgGUID, spaceGUID := fakeClient.BindSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("org1-space1-guid"))
			Expect(spaceGUID).Should(Equal("space1-guid"))
		})

		It("Should update and assign group to space", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "org1-space1",
					Guid:    "org1-space1-guid",
					Running: false,
					Staging: false,
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(1))
			sgGUID, spaceGUID := fakeClient.BindSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("org1-space1-guid"))
			Expect(spaceGUID).Should(Equal("space1-guid"))
		})

		It("Should update and not assign group to space", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "org1-space1",
					Guid:    "org1-space1-guid",
					Running: false,
					Staging: false,
					SpacesData: []cfclient.SpaceResource{
						{Entity: cfclient.Space{Guid: "space1-guid"}},
						{Entity: cfclient.Space{Guid: "space2-guid"}},
					},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(1))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should not create and not assign group to space", func() {
			securityMgr.Peek = true
			fakeClient.ListSecGroupsReturns(nil, nil)
			fakeClient.CreateSecGroupReturns(&cfclient.SecGroup{Name: "org1-space1", Guid: "org1-space1-guid"}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should not update and not assign group to space", func() {
			securityMgr.Peek = true
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "org1-space1",
					Guid:    "org1-space1-guid",
					Running: false,
					Staging: false,
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.BindSecGroupCallCount()).Should(Equal(0))
		})

		It("Should error on get space config", func() {
			fakeReader.GetSpaceConfigsReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})

		It("Should error listing security groups", func() {
			fakeClient.ListSecGroupsReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
		It("Should error returning space", func() {
			fakeSpaceMgr.FindSpaceReturns(cfclient.Space{}, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("CreateGlobalSecurityGroups", func() {
		It("should create 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
		})

		It("should create 1 asg from default asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetDefaultASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
		})

		It("should update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(1))
		})

		It("should not update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			securityGroupRules := []cfclient.SecGroupRule{}
			err := json.Unmarshal([]byte(asg_config), &securityGroupRules)
			Expect(err).ShouldNot(HaveOccurred())
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: securityGroupRules,
				},
			}, nil)
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err = securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(0))
		})

		It("should error create 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			fakeClient.CreateSecGroupReturns(nil, errors.New("error"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
		})

		It("should error on update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeClient.UpdateSecGroupReturns(nil, errors.New("error"))
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(1))
		})

		It("should error on update 1 asg from default config", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeClient.UpdateSecGroupReturns(nil, errors.New("error"))
			fakeReader.GetDefaultASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.UpdateSecGroupCallCount()).Should(Equal(1))
		})

		It("should error on getting asg config", func() {
			fakeReader.GetASGConfigsReturns(nil, errors.New("errorr"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
		It("should error on getting default asg config", func() {
			fakeReader.GetDefaultASGConfigsReturns(nil, errors.New("errorr"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
		It("should error on getting security groups", func() {
			fakeClient.ListSecGroupsReturns(nil, errors.New("errorr"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})

		It("should create 1 asg from asg config and remove trailing spaces", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config_with_trailing_space,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
			name, securityGroups, _ := fakeClient.CreateSecGroupArgsForCall(0)
			Expect(name).Should(Equal("asg-1"))
			Expect(len(securityGroups)).Should(Equal(1))
			Expect(securityGroups[0].Destination).Should(Equal("0.0.0.0/0"))
		})

		It("should create 1 asg from asg config and remove trailing spaces", func() {
			asgConfigs := []config.ASGConfig{
				{
					Name:  "asg-1",
					Rules: asg_config_with_space,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateSecGroupCallCount()).Should(Equal(1))
			name, securityGroups, _ := fakeClient.CreateSecGroupArgsForCall(0)
			Expect(name).Should(Equal("asg-1"))
			Expect(len(securityGroups)).Should(Equal(1))
			Expect(securityGroups[0].Destination).Should(Equal("10.0.11.0-10.0.11.2"))
		})
	})

	Context("AssignDefaultSecurityGroups", func() {
		It("should assign running security group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindRunningSecGroupCallCount()).Should(Equal(1))
		})

		It("should not assign running security group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "asg-1",
					Guid:    "asg-1-guid",
					Running: true,
					Rules:   []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindRunningSecGroupCallCount()).Should(Equal(0))
		})

		It("should error since group doesn't exist", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("Running security group [asg-1] does not exist"))
		})

		It("should assign running staging group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:  "asg-1",
					Guid:  "asg-1-guid",
					Rules: []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindStagingSecGroupCallCount()).Should(Equal(1))
		})

		It("should not assign staging security group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "asg-1",
					Guid:    "asg-1-guid",
					Staging: true,
					Rules:   []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindStagingSecGroupCallCount()).Should(Equal(0))
		})

		It("should error since group doesn't exist", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("Staging security group [asg-1] does not exist"))
		})

		It("should unassign running security group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "asg-1",
					Guid:    "asg-1-guid",
					Running: true,
					Rules:   []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableUnassignSecurityGroups: true,
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindRunningSecGroupCallCount()).Should(Equal(1))
		})
		It("should unassign staging security group", func() {
			fakeClient.ListSecGroupsReturns([]cfclient.SecGroup{
				{
					Name:    "asg-1",
					Guid:    "asg-1-guid",
					Staging: true,
					Rules:   []cfclient.SecGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableUnassignSecurityGroups: true,
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindStagingSecGroupCallCount()).Should(Equal(1))
		})
	})

	Context("ListSpaceSecurityGroups", func() {
		It("Should return 2", func() {
			fakeClient.ListSpaceSecGroupsReturns([]cfclient.SecGroup{
				{Name: "1", Guid: "1-guid"},
				{Name: "2", Guid: "2-guid"},
			}, nil)
			secGroups, err := securityMgr.ListSpaceSecurityGroups("spaceGUID")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(secGroups)).Should(Equal(2))
		})
		It("Should error", func() {
			fakeClient.ListSpaceSecGroupsReturns(nil, errors.New("error"))
			_, err := securityMgr.ListSpaceSecurityGroups("spaceGUID")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("GetSecurityGroupRules", func() {
		It("Should succeed", func() {
			securityGroupRules := []cfclient.SecGroupRule{}
			err := json.Unmarshal([]byte(asg_config), &securityGroupRules)
			Expect(err).ShouldNot(HaveOccurred())
			fakeClient.GetSecGroupReturns(&cfclient.SecGroup{
				Name:  "1",
				Guid:  "1-guid",
				Rules: securityGroupRules,
			}, nil)
			bytes, err := securityMgr.GetSecurityGroupRules("sgGUID")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bytes).ShouldNot(BeNil())
		})
		It("Should error", func() {
			fakeClient.GetSecGroupReturns(nil, errors.New("error"))
			_, err := securityMgr.GetSecurityGroupRules("sgGUID")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("UnassignStagingSecurityGroup", func() {
		It("Should succeed", func() {
			err := securityMgr.UnassignStagingSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindStagingSecGroupCallCount()).Should(Equal(1))
			sgGUID := fakeClient.UnbindStagingSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.UnassignStagingSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindStagingSecGroupCallCount()).Should(Equal(0))
		})
	})

	Context("UnassignRunningSecurityGroup", func() {
		It("Should succeed", func() {
			err := securityMgr.UnassignRunningSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindRunningSecGroupCallCount()).Should(Equal(1))
			sgGUID := fakeClient.UnbindRunningSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.UnassignRunningSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UnbindRunningSecGroupCallCount()).Should(Equal(0))
		})
	})

	Context("AssignStagingSecurityGroup", func() {
		It("Should succeed", func() {
			err := securityMgr.AssignStagingSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindStagingSecGroupCallCount()).Should(Equal(1))
			sgGUID := fakeClient.BindStagingSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.AssignStagingSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindStagingSecGroupCallCount()).Should(Equal(0))
		})
	})

	Context("AssignRunningSecurityGroup", func() {
		It("Should succeed", func() {
			err := securityMgr.AssignRunningSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindRunningSecGroupCallCount()).Should(Equal(1))
			sgGUID := fakeClient.BindRunningSecGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.AssignRunningSecurityGroup(cfclient.SecGroup{
				Name: "sec-group",
				Guid: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindRunningSecGroupCallCount()).Should(Equal(0))
		})
	})
})
