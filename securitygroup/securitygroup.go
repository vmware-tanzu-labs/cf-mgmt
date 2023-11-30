package securitygroup

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/go-cfclient/v3/resource"

	"github.com/pkg/errors"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/vmwarepivotallabs/cf-mgmt/space"
	"github.com/xchapter7x/lo"
)

// NewManager -
func NewManager(client CFSecurityGroupClient, spaceMgr space.Manager, cfg config.Reader, peek bool) Manager {
	return &DefaultManager{
		Cfg:          cfg,
		Client:       client,
		SpaceManager: spaceMgr,
		Peek:         peek,
	}
}

// DefaultSecurityGroupManager -
type DefaultManager struct {
	Cfg          config.Reader
	SpaceManager space.Manager
	Client       CFSecurityGroupClient
	Peek         bool
}

// CreateApplicationSecurityGroups -
func (m *DefaultManager) CreateApplicationSecurityGroups() error {
	spaceConfigs, err := m.Cfg.GetSpaceConfigs()
	if err != nil {
		return errors.Wrap(err, "Getting space configs")
	}
	sgs, err := m.ListNonDefaultSecurityGroups()
	if err != nil {
		return err
	}

	for _, input := range spaceConfigs {
		space, err := m.SpaceManager.FindSpace(input.Org, input.Space)
		if err != nil {
			return errors.Wrapf(err, "Finding org/space %s/%s", input.Org, input.Space)
		}
		existingSpaceSecurityGroups, err := m.ListSpaceSecurityGroups(space.GUID)
		if err != nil {
			return errors.Wrapf(err, "Unabled to list existing space security groups for org/space [%s/%s]", input.Org, input.Space)
		}
		lo.G.Debugf("Existing space security groups %+v", existingSpaceSecurityGroups)
		// iterate through and assign named security groups to the space - ensuring that they are up to date is
		// done elsewhere.
		for _, securityGroupName := range input.ASGs {
			if sgInfo, ok := sgs[securityGroupName]; ok {
				if _, ok := existingSpaceSecurityGroups[securityGroupName]; !ok {
					err := m.AssignSecurityGroupToSpace(space, sgInfo)
					if err != nil {
						return err
					}
				} else {
					delete(existingSpaceSecurityGroups, securityGroupName)
				}
			} else {
				return fmt.Errorf("Security group [%s] does not exist as a non-running and non-staging security group [%v+]", securityGroupName, sgs)
			}
		}

		spaceSecurityGroupName := fmt.Sprintf("%s-%s", input.Org, input.Space)
		if input.EnableSecurityGroup {
			var sgInfo *resource.SecurityGroup
			var ok bool
			if sgInfo, ok = sgs[spaceSecurityGroupName]; ok {
				changed, err := m.hasSecurityGroupChanged(sgInfo, input.GetSecurityGroupContents())
				if err != nil {
					return errors.Wrapf(err, "Checking if security group %s has changed", spaceSecurityGroupName)
				}
				if changed {
					if err := m.UpdateSecurityGroup(sgInfo, input.GetSecurityGroupContents()); err != nil {
						return err
					}
				}
			} else {
				securityGroup, err := m.CreateSecurityGroup(spaceSecurityGroupName, input.GetSecurityGroupContents())
				if err != nil {
					return errors.Wrapf(err, "Creating security group %s for %s/%s security-group.json", spaceSecurityGroupName, input.Org, input.Space)
				}
				sgInfo = securityGroup
			}
			if _, ok := existingSpaceSecurityGroups[spaceSecurityGroupName]; !ok {
				err := m.AssignSecurityGroupToSpace(space, sgInfo)
				if err != nil {
					return err
				}
			} else {
				delete(existingSpaceSecurityGroups, spaceSecurityGroupName)
			}
		}

		if input.EnableUnassignSecurityGroup {
			lo.G.Debugf("Existing space security groups after %+v", existingSpaceSecurityGroups)
			for sgName, _ := range existingSpaceSecurityGroups {
				if sgInfo, ok := sgs[sgName]; ok {
					err := m.UnassignSecurityGroupToSpace(space, sgInfo)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("Security group [%s] does not exist as a non-running and non-staging security group [%v+]", sgName, sgs)
				}
			}
		}
	}
	return nil
}

func (m *DefaultManager) ListSecurityGroups() (map[string]*resource.SecurityGroup, error) {
	securityGroups := make(map[string]*resource.SecurityGroup)
	secGroups, err := m.Client.ListAll(context.Background(), nil)
	if err != nil {
		return securityGroups, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	for _, sg := range secGroups {
		securityGroups[sg.Name] = sg
	}
	return securityGroups, nil
}

// CreateGlobalSecurityGroups -
func (m *DefaultManager) CreateGlobalSecurityGroups() error {
	sgs, err := m.ListSecurityGroups()
	if err != nil {
		return err
	}
	securityGroupConfigs, err := m.Cfg.GetASGConfigs()
	if err != nil {
		return errors.Wrap(err, "Getting ASG Configs")
	}
	defaultSecurityGroupConfigs, err := m.Cfg.GetDefaultASGConfigs()
	if err != nil {
		return errors.Wrap(err, "Getting Default ASG Configs")
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

// AssignDefaultSecurityGroups -
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
			if !group.GloballyEnabled.Running {
				err = m.AssignSecurityGroupGlobalRunning(group)
				if err != nil {
					return err
				}
			}
		} else {
			if !m.Peek {
				return fmt.Errorf("Running security group [%s] does not exist", runningGroup)
			} else {
				lo.G.Infof("[dry-run]: assigning yet to be created sg %s as running security group", runningGroup)
			}
		}
	}

	for _, stagingGroup := range globalConfig.StagingSecurityGroups {
		if group, ok := sgs[stagingGroup]; ok {
			if !group.GloballyEnabled.Staging {
				err = m.AssignSecurityGroupGlobalStaging(group)
				if err != nil {
					return err
				}
			}
		} else {
			if !m.Peek {
				return fmt.Errorf("Staging security group [%s] does not exist", stagingGroup)
			} else {
				lo.G.Infof("[dry-run]: assigning yet to be created sg %s as staging security group", stagingGroup)
			}
		}
	}

	if globalConfig.EnableUnassignSecurityGroups {
		for groupName, group := range sgs {
			if group.GloballyEnabled.Running && !m.contains(globalConfig.RunningSecurityGroups, groupName) {
				err = m.UnassignSecurityGroupGlobalRunning(group)
				if err != nil {
					return err
				}
			}
			if group.GloballyEnabled.Staging && !m.contains(globalConfig.StagingSecurityGroups, groupName) {
				err = m.UnassignSecurityGroupGlobalStaging(group)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *DefaultManager) contains(list []string, groupName string) bool {
	groupNameToUpper := strings.ToUpper(groupName)
	for _, v := range list {
		if strings.ToUpper(v) == groupNameToUpper {
			return true
		}
	}
	return false
}

func (m *DefaultManager) processSecurityGroups(securityGroupConfigs []config.ASGConfig, sgs map[string]*resource.SecurityGroup) error {
	for _, input := range securityGroupConfigs {
		sgName := input.Name

		// For every named security group
		// Check if it's a new group or Update
		if sgInfo, ok := sgs[sgName]; ok {
			changed, err := m.hasSecurityGroupChanged(sgInfo, input.Rules)
			if err != nil {
				return errors.Wrapf(err, "Processing %s security group", sgName)
			}
			if changed {
				if err := m.UpdateSecurityGroup(sgInfo, input.Rules); err != nil {
					return err
				}
			}
		} else {
			if _, err := m.CreateSecurityGroup(sgName, input.Rules); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *DefaultManager) hasSecurityGroupChanged(sgInfo *resource.SecurityGroup, rules string) (bool, error) {
	sgInfo.Rules = m.updateSecurityRulesWithDefaults(sgInfo.Rules)
	jsonBytes, err := json.Marshal(sgInfo.Rules)
	if err != nil {
		return false, err
	}
	secRules := []resource.SecurityGroupRule{}
	err = json.Unmarshal([]byte(rules), &secRules)
	if err != nil {
		return false, err
	}
	secRules = m.updateSecurityRulesWithDefaults(secRules)

	jsonBytesToCompare, err := json.Marshal(secRules)
	if err != nil {
		return false, err
	}
	match, err := DoesJsonMatch(string(jsonBytes), string(jsonBytesToCompare))
	if err != nil {
		return false, err
	}
	if !match {
		lo.G.Infof("Security Group %s has changed from %s to %s", sgInfo.Name, string(jsonBytes), string(jsonBytesToCompare))
	}
	return !match, nil
}

func (m *DefaultManager) updateSecurityRulesWithDefaults(rules []resource.SecurityGroupRule) []resource.SecurityGroupRule {
	updatedRules := []resource.SecurityGroupRule{}
	for _, secRule := range rules {
		if secRule.Code == nil {
			intCode := new(int)
			*intCode = 0
			secRule.Code = intCode
		}
		if secRule.Type == nil {
			intType := new(int)
			*intType = 0
			secRule.Type = intType
		}
		if secRule.Ports == nil {
			sgPorts := new(string)
			*sgPorts = ""
			secRule.Ports = sgPorts
		}
		updatedRules = append(updatedRules, secRule)
	}
	return updatedRules
}

func (m *DefaultManager) isSecurityGroupAssignedToSpace(space *resource.Space, secGroup *resource.SecurityGroup) bool {
	for _, spaceRelation := range secGroup.Relationships.RunningSpaces.Data {
		if spaceRelation.GUID == space.GUID {
			return true
		}
	}
	for _, spaceRelation := range secGroup.Relationships.StagingSpaces.Data {
		if spaceRelation.GUID == space.GUID {
			return true
		}
	}
	return false
}
func (m *DefaultManager) AssignSecurityGroupToSpace(space *resource.Space, secGroup *resource.SecurityGroup) error {
	if m.isSecurityGroupAssignedToSpace(space, secGroup) {
		lo.G.Debugf("Security group %s is already assigned to space %s, skipping", secGroup.Name, space.Name)
		return nil
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning security group %s to space %s", secGroup.Name, space.Name)
		return nil
	}
	lo.G.Infof("assigning security group %s to space %s", secGroup.Name, space.Name)
	_, err := m.Client.BindRunningSecurityGroup(context.Background(), secGroup.GUID, []string{space.GUID})
	return err
}

func (m *DefaultManager) UnassignSecurityGroupToSpace(space *resource.Space, secGroup *resource.SecurityGroup) error {
	if !m.isSecurityGroupAssignedToSpace(space, secGroup) {
		lo.G.Debugf("Security group %s isn't assigned to space %s, skipping", secGroup.Name, space.Name)
		return nil
	}
	if m.Peek {
		lo.G.Infof("[dry-run]: unassigning security group %s to space %s", secGroup.Name, space.Name)
		return nil
	}
	lo.G.Infof("unassigning security group %s to space %s", secGroup.Name, space.Name)
	return m.Client.UnBindRunningSecurityGroup(context.Background(), secGroup.GUID, space.GUID)
}

func (m *DefaultManager) removeDestinationWhitespace(rules []*resource.SecurityGroupRule) []*resource.SecurityGroupRule {
	rulesToReturn := []*resource.SecurityGroupRule{}
	for _, rule := range rules {
		rulesToReturn = append(rulesToReturn, &resource.SecurityGroupRule{
			Protocol:    rule.Protocol,
			Ports:       rule.Ports,
			Destination: strings.Replace(rule.Destination, " ", "", -1),
			Description: rule.Description,
			Code:        rule.Code,
			Type:        rule.Type,
			Log:         rule.Log,
		})
	}
	return rulesToReturn
}

func (m *DefaultManager) CreateSecurityGroup(sgName, contents string) (*resource.SecurityGroup, error) {
	if m.Peek {
		lo.G.Infof("[dry-run]: creating securityGroup %s with contents %s", sgName, contents)
		return &resource.SecurityGroup{Name: "dry-run-name", GUID: "dry-run-guid"}, nil
	}
	securityGroupRules := []*resource.SecurityGroupRule{}
	err := json.Unmarshal([]byte(contents), &securityGroupRules)
	if err != nil {
		return nil, err
	}
	rulesToUse := m.removeDestinationWhitespace(securityGroupRules)
	lo.G.Infof("creating securityGroup %s with contents %+v", sgName, m.rulesAsString(rulesToUse))

	r := &resource.SecurityGroupCreate{
		Name: sgName,
		GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
			Running: false,
			Staging: false,
		},
		Rules: rulesToUse,
	}
	return m.Client.Create(context.Background(), r)
}

func (m *DefaultManager) rulesAsString(rules []*resource.SecurityGroupRule) string {
	var ruleString string
	for _, rule := range rules {
		var description string
		if rule.Description != nil {
			description = *rule.Description
		}
		var sgCode int
		if rule.Code == nil {
			sgCode = 0
		} else {
			sgCode = *rule.Code
		}
		var sgType int
		if rule.Type == nil {
			sgType = 0
		} else {
			sgType = *rule.Type
		}
		var sgLog bool
		if rule.Log == nil {
			sgLog = false
		} else {
			sgLog = *rule.Log
		}
		var sgPorts string
		if rule.Ports == nil {
			sgPorts = ""
		} else {
			sgPorts = *rule.Ports
		}
		ruleString = ruleString + fmt.Sprintf("[Protocol:%s Destination:%s Ports:%s Type:%d Code:%d Description:%s Log:%t]", rule.Protocol, rule.Destination, sgPorts, sgType, sgCode, description, sgLog)
	}
	return ruleString
}

func (m *DefaultManager) UpdateSecurityGroup(sg *resource.SecurityGroup, contents string) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: updating securityGroup %s with contents %s", sg.Name, contents)
		return nil
	}
	securityGroupRules := []*resource.SecurityGroupRule{}
	err := json.Unmarshal([]byte(contents), &securityGroupRules)
	if err != nil {
		return err
	}
	rulesToUse := m.removeDestinationWhitespace(securityGroupRules)
	lo.G.Infof("updating securityGroup %s with contents %+v", sg.Name, m.rulesAsString(rulesToUse))

	r := &resource.SecurityGroupUpdate{
		Name:  sg.Name,
		Rules: rulesToUse,
	}
	_, err = m.Client.Update(context.Background(), sg.GUID, r)
	return err
}

func (m *DefaultManager) ListNonDefaultSecurityGroups() (map[string]*resource.SecurityGroup, error) {
	securityGroups := make(map[string]*resource.SecurityGroup)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.GloballyEnabled.Running == false && groupMap.GloballyEnabled.Staging == false {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}

func (m *DefaultManager) ListDefaultSecurityGroups() (map[string]*resource.SecurityGroup, error) {
	securityGroups := make(map[string]*resource.SecurityGroup)
	groupMap, err := m.ListSecurityGroups()
	if err != nil {
		return nil, err
	}
	for key, groupMap := range groupMap {
		if groupMap.GloballyEnabled.Running == true || groupMap.GloballyEnabled.Staging == true {
			securityGroups[key] = groupMap
		}
	}
	return securityGroups, nil
}

func (m *DefaultManager) AssignSecurityGroupGlobalRunning(sg *resource.SecurityGroup) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sg %s as running security group", sg.Name)
		return nil
	}
	lo.G.Infof("assigning sg %s as running security group", sg.Name)
	r := &resource.SecurityGroupUpdate{
		GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
			Running: true,
		},
	}
	sg.GloballyEnabled.Running = true
	_, err := m.Client.Update(context.Background(), sg.GUID, r)
	return err
}

func (m *DefaultManager) AssignSecurityGroupGlobalStaging(sg *resource.SecurityGroup) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: assigning sg %s as staging security group", sg.Name)
		return nil
	}
	lo.G.Infof("assigning sg %s as staging security group", sg.Name)
	r := &resource.SecurityGroupUpdate{
		GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
			Staging: true,
		},
	}
	sg.GloballyEnabled.Staging = true
	_, err := m.Client.Update(context.Background(), sg.GUID, r)
	return err
}

func (m *DefaultManager) UnassignSecurityGroupGlobalRunning(sg *resource.SecurityGroup) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassinging sg %s as running security group", sg.Name)
		return nil
	}
	lo.G.Infof("unassinging sg %s as running security group", sg.Name)
	r := &resource.SecurityGroupUpdate{
		GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
			Running: false,
		},
	}
	sg.GloballyEnabled.Running = false
	_, err := m.Client.Update(context.Background(), sg.GUID, r)
	return err
}

func (m *DefaultManager) UnassignSecurityGroupGlobalStaging(sg *resource.SecurityGroup) error {
	if m.Peek {
		lo.G.Infof("[dry-run]: unassigning sg %s as staging security group", sg.Name)
		return nil
	}
	lo.G.Infof("unassigning sg %s as staging security group", sg.Name)
	r := &resource.SecurityGroupUpdate{
		GloballyEnabled: &resource.SecurityGroupGloballyEnabled{
			Staging: false,
		},
	}
	sg.GloballyEnabled.Staging = false
	_, err := m.Client.Update(context.Background(), sg.GUID, r)
	return err
}

func (m *DefaultManager) GetSecurityGroupRules(sgGUID string) ([]byte, error) {
	secGroup, err := m.Client.Get(context.Background(), sgGUID)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(secGroup.Rules, "", "\t")
}

func (m *DefaultManager) ListSpaceSecurityGroups(spaceGUID string) (map[string]string, error) {
	names := make(map[string]string)
	if strings.Contains(spaceGUID, "dry-run-space-guid") {
		return names, nil
	}
	secGroups, err := m.Client.ListRunningForSpaceAll(context.Background(), spaceGUID, nil)
	if err != nil {
		return nil, err
	}
	lo.G.Debug("Total security groups returned :", len(secGroups))
	for _, sg := range secGroups {
		if sg.GloballyEnabled.Running == false && sg.GloballyEnabled.Staging == false {
			names[sg.Name] = sg.GUID
		}
	}
	return names, nil
}
