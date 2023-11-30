package securitygroup_test

import (
	"encoding/json"
	"errors"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

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
		fakeClient   *securitygroupfakes.FakeCFSecurityGroupClient
		fakeSpaceMgr *spacefakes.FakeManager
		securityMgr  securitygroup.DefaultManager
	)

	BeforeEach(func() {
		fakeReader = new(configfakes.FakeReader)
		fakeSpaceMgr = new(spacefakes.FakeManager)
		fakeClient = new(securitygroupfakes.FakeCFSecurityGroupClient)
		securityMgr = securitygroup.DefaultManager{
			Cfg:          fakeReader,
			Client:       fakeClient,
			SpaceManager: fakeSpaceMgr,
			Peek:         false,
		}
	})
	Context("ListNonDefaultSecurityGroups", func() {
		It("returns 2 security groups", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "group1",
					GUID: "group1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: false,
					},
				},
				{
					Name: "group2",
					GUID: "group2-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: false,
					},
				},
			}, nil)
			groups, err := securityMgr.ListNonDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(groups)).Should(Equal(2))
		})

		It("returns error", func() {
			fakeClient.ListAllReturns(nil, errors.New("error"))
			_, err := securityMgr.ListNonDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.ListAllCallCount()).Should(Equal(1))
		})
	})

	Context("ListDefaultSecurityGroups", func() {
		It("returns 2 security groups", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "group1",
					GUID: "group1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: true,
						Staging: false,
					},
				},
				{
					Name: "group2",
					GUID: "group2-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: true,
					},
				},
			}, nil)
			groups, err := securityMgr.ListDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(groups)).Should(Equal(2))
		})

		It("returns error", func() {
			fakeClient.ListAllReturns(nil, errors.New("error"))
			_, err := securityMgr.ListDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.ListAllCallCount()).Should(Equal(1))
		})
	})

	Context("CreateApplicationSecurityGroups", func() {
		BeforeEach(func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup:   true,
					Space:                 "space1",
					Org:                   "org1",
					SecurityGroupContents: asg_config,
				},
				config.SpaceConfig{
					EnableSecurityGroup: false,
					Space:               "space2",
					Org:                 "org1",
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeSpaceMgr.FindSpaceReturns(&resource.Space{
				Name: "space1",
				GUID: "space1-guid",
				Relationships: &resource.SpaceRelationships{
					Organization: &resource.ToOneRelationship{
						Data: &resource.Relationship{
							GUID: "org1-guid",
						},
					},
				},
			}, nil)
		})

		It("Should assign global group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(1))
			_, sgGUID, spaceGUIDs := fakeClient.BindRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("dns-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		})

		It("Should not assign global group to space that is already assigned", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
			}, nil)
			fakeClient.ListRunningForSpaceAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListRunningForSpaceAllCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should unbind global group to space that is not in configuration", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup:         false,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
					Relationships: resource.SecurityGroupsRelationships{
						RunningSpaces: resource.ToManyRelationships{
							Data: []resource.Relationship{
								{GUID: "space1-guid"},
							},
						},
					},
				},
				{
					Name: "ntp",
					GUID: "ntp-guid",
					Relationships: resource.SecurityGroupsRelationships{
						RunningSpaces: resource.ToManyRelationships{
							Data: []resource.Relationship{
								{GUID: "space1-guid"},
							},
						},
					},
				},
			}, nil)
			fakeClient.ListRunningForSpaceAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
				{
					Name: "ntp",
					GUID: "ntp-guid",
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListRunningForSpaceAllCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnBindRunningSecurityGroupCallCount()).Should(Equal(1))
			_, secGroupGUID, spaceGUID := fakeClient.UnBindRunningSecurityGroupArgsForCall(0)
			Expect(secGroupGUID).Should(Equal("ntp-guid"))
			Expect(spaceGUID).Should(Equal("space1-guid"))
		})

		It("Should not unbind space specific group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup:         true,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
				{
					Name:  "org1-space1",
					GUID:  "org1-space1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeClient.ListRunningForSpaceAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
				{
					Name:  "org1-space1",
					GUID:  "org1-space1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListRunningForSpaceAllCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnBindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should unbind space specific group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup:         false,
					Space:                       "space1",
					Org:                         "org1",
					ASGs:                        []string{"dns"},
					EnableUnassignSecurityGroup: true,
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
					Relationships: resource.SecurityGroupsRelationships{
						RunningSpaces: resource.ToManyRelationships{
							Data: []resource.Relationship{
								{GUID: "space1-guid"},
							},
						},
					},
				},
				{
					Name:  "org1-space1",
					GUID:  "org1-space1-guid",
					Rules: []resource.SecurityGroupRule{},
					Relationships: resource.SecurityGroupsRelationships{
						RunningSpaces: resource.ToManyRelationships{
							Data: []resource.Relationship{
								{GUID: "space1-guid"},
							},
						},
					},
				},
			}, nil)
			fakeClient.ListRunningForSpaceAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
				{
					Name:  "org1-space1",
					GUID:  "org1-space1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.ListRunningForSpaceAllCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
			Expect(fakeClient.UnBindRunningSecurityGroupCallCount()).Should(Equal(1))
		})

		It("Should error assigning global group to space", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
					EnableSecurityGroup: false,
					Space:               "space1",
					Org:                 "org1",
					ASGs:                []string{"dns"},
				},
			}
			fakeReader.GetSpaceConfigsReturns(spaceConfigs, nil)
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "dns",
					GUID: "dns-guid",
				},
			}, nil)
			fakeClient.BindRunningSecurityGroupReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(1))
		})

		It("Should error when group doesn't exist", func() {
			spaceConfigs := []config.SpaceConfig{
				config.SpaceConfig{
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
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should create and assign group to space", func() {
			fakeClient.ListAllReturns(nil, nil)
			fakeClient.CreateReturns(&resource.SecurityGroup{Name: "org1-space1", GUID: "org1-space1-guid"}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(1))
			_, sgGUID, spaceGUIDs := fakeClient.BindRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("org1-space1-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		})

		It("Should update and assign group to space", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "org1-space1",
					GUID: "org1-space1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: false,
					},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(1))
			_, sgGUID, spaceGUIDs := fakeClient.BindRunningSecurityGroupArgsForCall(0)
			Expect(sgGUID).Should(Equal("org1-space1-guid"))
			Expect(spaceGUIDs[0]).Should(Equal("space1-guid"))
		})

		It("Should update and not assign group to space", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "org1-space1",
					GUID: "org1-space1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: false,
					},
					Relationships: resource.SecurityGroupsRelationships{
						RunningSpaces: resource.ToManyRelationships{
							Data: []resource.Relationship{
								{GUID: "space1-guid"},
								{GUID: "space2-guid"},
							},
						},
					},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should not create and not assign group to space", func() {
			securityMgr.Peek = true
			fakeClient.ListAllReturns(nil, nil)
			fakeClient.CreateReturns(&resource.SecurityGroup{Name: "org1-space1", GUID: "org1-space1-guid"}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(0))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should not update and not assign group to space", func() {
			securityMgr.Peek = true
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "org1-space1",
					GUID: "org1-space1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: false,
					},
				},
			}, nil)
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
			Expect(fakeClient.BindRunningSecurityGroupCallCount()).Should(Equal(0))
		})

		It("Should error on get space config", func() {
			fakeReader.GetSpaceConfigsReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})

		It("Should error listing security groups", func() {
			fakeClient.ListAllReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
		It("Should error returning space", func() {
			fakeSpaceMgr.FindSpaceReturns(nil, errors.New("error"))
			err := securityMgr.CreateApplicationSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("CreateGlobalSecurityGroups", func() {
		It("should create 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
		})

		It("should create 1 asg from default asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetDefaultASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
		})

		It("should update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
		})

		It("should not update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			securityGroupRules := []resource.SecurityGroupRule{}
			err := json.Unmarshal([]byte(asg_config), &securityGroupRules)
			Expect(err).ShouldNot(HaveOccurred())
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: securityGroupRules,
				},
			}, nil)
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err = securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})

		It("should error create 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			fakeClient.CreateReturns(nil, errors.New("error"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
		})

		It("should error on update 1 asg from asg config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeClient.UpdateReturns(nil, errors.New("error"))
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
		})

		It("should error on update 1 asg from default config", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config,
				},
			}
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeClient.UpdateReturns(nil, errors.New("error"))
			fakeReader.GetDefaultASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
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
			fakeClient.ListAllReturns(nil, errors.New("errorr"))
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).Should(HaveOccurred())
		})

		It("should create 1 asg from asg config and remove trailing spaces", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config_with_trailing_space,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
			_, securityGroupCreate := fakeClient.CreateArgsForCall(0)
			Expect(securityGroupCreate.Name).Should(Equal("asg-1"))
			Expect(len(securityGroupCreate.Rules)).Should(Equal(1))
			Expect(securityGroupCreate.Rules[0].Destination).Should(Equal("0.0.0.0/0"))
		})

		It("should create 1 asg from asg config and remove trailing spaces", func() {
			asgConfigs := []config.ASGConfig{
				config.ASGConfig{
					Name:  "asg-1",
					Rules: asg_config_with_space,
				},
			}
			fakeReader.GetASGConfigsReturns(asgConfigs, nil)
			err := securityMgr.CreateGlobalSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.CreateCallCount()).Should(Equal(1))
			_, securityGroupCreate := fakeClient.CreateArgsForCall(0)
			Expect(securityGroupCreate.Name).Should(Equal("asg-1"))
			Expect(len(securityGroupCreate.Rules)).Should(Equal(1))
			Expect(securityGroupCreate.Rules[0].Destination).Should(Equal("10.0.11.0-10.0.11.2"))
		})
	})

	Context("AssignDefaultSecurityGroups", func() {
		It("should assign running security group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).To(Equal("asg-1-guid"))
			Expect(updateRequest.GloballyEnabled.Running).To(BeTrue())
		})

		It("should not assign running security group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "asg-1",
					GUID: "asg-1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: true,
						Staging: false,
					},
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})

		It("should error since group doesn't exist", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				RunningSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("Running security group [asg-1] does not exist"))
		})

		It("should assign running staging group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name:  "asg-1",
					GUID:  "asg-1-guid",
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).To(Equal("asg-1-guid"))
			Expect(updateRequest.GloballyEnabled.Staging).To(BeTrue())
		})

		It("should not assign staging security group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "asg-1",
					GUID: "asg-1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: true,
					},
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})

		It("should error since group doesn't exist", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				StagingSecurityGroups: []string{"asg-1"},
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("Staging security group [asg-1] does not exist"))
		})

		It("should unassign running security group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "asg-1",
					GUID: "asg-1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: true,
						Staging: false,
					},
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableUnassignSecurityGroups: true,
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).To(Equal("asg-1-guid"))
			Expect(updateRequest.GloballyEnabled.Running).To(BeFalse())
		})
		It("should unassign staging security group", func() {
			fakeClient.ListAllReturns([]*resource.SecurityGroup{
				{
					Name: "asg-1",
					GUID: "asg-1-guid",
					GloballyEnabled: resource.SecurityGroupGloballyEnabled{
						Running: false,
						Staging: true,
					},
					Rules: []resource.SecurityGroupRule{},
				},
			}, nil)
			fakeReader.GetGlobalConfigReturns(&config.GlobalConfig{
				EnableUnassignSecurityGroups: true,
			}, nil)
			err := securityMgr.AssignDefaultSecurityGroups()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).To(Equal("asg-1-guid"))
			Expect(updateRequest.GloballyEnabled.Staging).To(BeFalse())
		})
	})

	Context("ListSpaceSecurityGroups", func() {
		It("Should return 2", func() {
			fakeClient.ListRunningForSpaceAllReturns([]*resource.SecurityGroup{
				{Name: "1", GUID: "1-guid"},
				{Name: "2", GUID: "2-guid"},
			}, nil)
			secGroups, err := securityMgr.ListSpaceSecurityGroups("spaceGUID")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(secGroups)).Should(Equal(2))
		})
		It("Should error", func() {
			fakeClient.ListRunningForSpaceAllReturns(nil, errors.New("error"))
			_, err := securityMgr.ListSpaceSecurityGroups("spaceGUID")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("GetSecurityGroupRules", func() {
		It("Should succeed", func() {
			securityGroupRules := []resource.SecurityGroupRule{}
			err := json.Unmarshal([]byte(asg_config), &securityGroupRules)
			Expect(err).ShouldNot(HaveOccurred())
			fakeClient.GetReturns(&resource.SecurityGroup{
				Name:  "1",
				GUID:  "1-guid",
				Rules: securityGroupRules,
			}, nil)
			bytes, err := securityMgr.GetSecurityGroupRules("sgGUID")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bytes).ShouldNot(BeNil())
		})
		It("Should error", func() {
			fakeClient.GetReturns(nil, errors.New("error"))
			_, err := securityMgr.GetSecurityGroupRules("sgGUID")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("UnassignSecurityGroupGlobalStaging", func() {
		It("Should succeed", func() {
			err := securityMgr.UnassignSecurityGroupGlobalStaging(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, _ := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.UnassignSecurityGroupGlobalStaging(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})
	})

	Context("UnassignSecurityGroupGlobalRunning", func() {
		It("Should succeed", func() {
			err := securityMgr.UnassignSecurityGroupGlobalRunning(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
			Expect(updateRequest.GloballyEnabled.Running).To(BeFalse())
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.UnassignSecurityGroupGlobalRunning(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})
	})

	Context("AssignSecurityGroupGlobalStaging", func() {
		It("Should succeed", func() {
			err := securityMgr.AssignSecurityGroupGlobalStaging(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
			Expect(updateRequest.GloballyEnabled.Staging).To(BeTrue())
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.AssignSecurityGroupGlobalStaging(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})
	})

	Context("AssignSecurityGroupGlobalRunning", func() {
		It("Should succeed", func() {
			err := securityMgr.AssignSecurityGroupGlobalRunning(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(1))
			_, sgGUID, updateRequest := fakeClient.UpdateArgsForCall(0)
			Expect(sgGUID).Should(Equal("seg-group-guid"))
			Expect(updateRequest.GloballyEnabled.Running).To(BeTrue())
		})
		It("Should peek", func() {
			securityMgr.Peek = true
			err := securityMgr.AssignSecurityGroupGlobalRunning(&resource.SecurityGroup{
				Name: "sec-group",
				GUID: "seg-group-guid",
			})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(fakeClient.UpdateCallCount()).Should(Equal(0))
		})
	})
})
