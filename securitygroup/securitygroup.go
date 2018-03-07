package securitygroup

import (
	"encoding/json"
	"fmt"
	"strings"

	cfclient "github.com/cloudfoundry-community/go-cfclient"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

//NewManager -
func NewManager(client CFClient, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:    cfg,
		Client: client,
		Peek:   peek,
	}
}

//DefaultSecurityGroupManager -
type DefaultManager struct {
	Cfg         config.Reader
	FilePattern string
	FilePaths   []string
	Client      CFClient
	Peek        bool
}

func (m *DefaultManager) ListSecurityGroups() (map[string]cfclient.SecGroup, error) {
	securityGroups := make(map[string]cfclient.SecGroup)
	secGroups, err := m.Client.ListSecGroups()
	if err != nil {
		return securityGroups, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	for _, sg := range secGroups {
		securityGroups[sg.Name] = sg
	}
	return securityGroups, nil
}

//CreateApplicationSecurityGroups -
func (m *DefaultManager) CreateApplicationSecurityGroups() error {
	sgs, err := m.ListSecurityGroups()
	if err != nil {
		return err
	}
	securityGroupConfigs, err := m.Cfg.GetASGConfigs()
	if err != nil {
		return err
	}
	defaultSecurityGroupConfigs, err := m.Cfg.GetDefaultASGConfigs()
	if err != nil {
		return err
	}
	err = m.processSecurityGroups(securityGroupConfigs, sgs)
	if err != nil {
		return err
	}
	err = m.processSecurityGroups(defaultSecurityGroupConfigs, sgs)
	if err != nil {
		return err
	}

	return nil
}

//AssignDefaultSecurityGroups -
func (m *DefaultManager) AssignDefaultSecurityGroups() error {
	sgs, err := m.ListSecurityGroups()
	if err != nil {
		return err
	}
	globalConfig, err := m.Cfg.GetGlobalConfig()
	if err != nil {
		return err
	}

	for _, runningGroup := range globalConfig.RunningSecurityGroups {
		if group, ok := sgs[runningGroup]; ok {
			if !group.Running {
				lo.G.Infof("assigning security group %s as running security group", runningGroup)
				err = m.AssignRunningSecurityGroup(group.Guid)
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("Running security group [%s] does not exist", runningGroup)
		}
	}

	for _, stagingGroup := range globalConfig.StagingSecurityGroups {
		if group, ok := sgs[stagingGroup]; ok {
			if !group.Staging {
				lo.G.Infof("assigning security group [%s] as staging security group", stagingGroup)
				err = m.AssignStagingSecurityGroup(group.Guid)
				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("Staging security group %s does not exist", stagingGroup)
		}
	}

	if globalConfig.EnableUnassignSecurityGroups {
		for groupName, group := range sgs {
			if group.Running && !m.Contains(globalConfig.RunningSecurityGroups, groupName) {
				lo.G.Infof("unassigning security group %s as running security group", groupName)
				err = m.UnassignRunningSecurityGroup(group.Guid)
				if err != nil {
					return err
				}
			}
			if group.Staging && !m.Contains(globalConfig.StagingSecurityGroups, groupName) {
				lo.G.Infof("unassigning security group %s as staging security group", groupName)
				err = m.UnassignStagingSecurityGroup(group.Guid)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *DefaultManager) Contains(list []string, groupName string) bool {
	groupNameToUpper := strings.ToUpper(groupName)
	for _, v := range list {
		if strings.ToUpper(v) == groupNameToUpper {
			return true
		}
	}
	return false
}

func (m *DefaultManager) processSecurityGroups(securityGroupConfigs []config.ASGConfig, sgs map[string]cfclient.SecGroup) error {
	for _, input := range securityGroupConfigs {
		sgName := input.Name

		// For every named security group
		// Check if it's a new group or Update
		if sgInfo, ok := sgs[sgName]; ok {
			jsonBytes, err := json.Marshal(sgInfo.Rules)
			if err != nil {
				return err
			}
			match, err := DoesJsonMatch(string(jsonBytes), input.Rules)
			if err != nil {
				return err
			}
			if !match {
				lo.G.Info("Updating security group", sgName)
				if err := m.UpdateSecurityGroup(sgInfo.Guid, sgName, input.Rules); err != nil {
					return err
				}
			}
		} else {
			lo.G.Info("Creating security group", sgName)
			if _, err := m.CreateSecurityGroup(sgName, input.Rules); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DefaultManager) AssignSecurityGroupToSpace(spaceGUID, sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s to spaceGUID %s", sgGUID, spaceGUID)
		return nil
	}
	return m.Client.BindSecGroup(sgGUID, spaceGUID)
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (*cfclient.SecGroup, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating securityGroup %s with contents %s", sgName, contents)
		return nil, nil
	}
	securityGroup := &cfclient.SecGroup{}
	err := json.Unmarshal([]byte(contents), &securityGroup)
	if err != nil {
		return nil, err
	}
	securityGroup, err = m.Client.CreateSecGroup(sgName, securityGroup.Rules, nil)
	return securityGroup, err
}

func (m *DefaultManager) UpdateSecurityGroup(sgGUID, sgName, contents string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: updating securityGroup %s with guid %s with contents %s", sgName, sgGUID, contents)
		return nil
	}
	securityGroup := &cfclient.SecGroup{}
	err := json.Unmarshal([]byte(contents), &securityGroup)
	if err != nil {
		return err
	}
	_, err = m.Client.UpdateSecGroup(sgGUID, sgName, securityGroup.Rules, nil)
	return err
}
func (m *DefaultManager) ListNonDefaultSecurityGroups() (map[string]cfclient.SecGroup, error) {
	securityGroups := make(map[string]cfclient.SecGroup)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.Running == false && groupMap.Staging == false {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}

func (m *DefaultManager) ListDefaultSecurityGroups() (map[string]cfclient.SecGroup, error) {
	securityGroups := make(map[string]cfclient.SecGroup)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.Running == true || groupMap.Staging == true {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}

func (m *DefaultManager) AssignRunningSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s as running security group", sgGUID)
		return nil
	}
	return m.Client.BindRunningSecGroup(sgGUID)
}
func (m *DefaultManager) AssignStagingSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sgGUID %s as staging security group", sgGUID)
		return nil
	}
	return m.Client.BindStagingSecGroup(sgGUID)
}
func (m *DefaultManager) UnassignRunningSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassinging sgGUID %s as running security group", sgGUID)
		return nil
	}
	return m.Client.UnbindRunningSecGroup(sgGUID)
}
func (m *DefaultManager) UnassignStagingSecurityGroup(sgGUID string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassigning sgGUID %s as staging security group", sgGUID)
		return nil
	}
	return m.Client.UnbindStagingSecGroup(sgGUID)
}

func (m *DefaultManager) GetSecurityGroupRules(sgGUID string) ([]byte, error) {
	secGroup, err := m.Client.GetSecGroup(sgGUID)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(secGroup.Rules, "", "\t")
}

func (m *DefaultManager) ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error) {
	secGroups, err := m.Client.ListSpaceSecGroups(spaceGUID)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	names := make(map[string]string)
	for _, sg := range secGroups {
		if sg.Running == false && sg.Staging == false {
			names[sg.Name] = sg.Guid
		}
	}
	return names, nil
}
