package resource

import "time"

// SecurityGroup implements the security group object. Security groups are collections
// of egress traffic rules that can be applied to the staging or running state of applications.
type SecurityGroup struct {
	GUID            string                       `json:"guid"`
	CreatedAt       time.Time                    `json:"created_at"`
	UpdatedAt       time.Time                    `json:"updated_at"`
	Name            string                       `json:"name"`
	GloballyEnabled SecurityGroupGloballyEnabled `json:"globally_enabled"`
	Rules           []SecurityGroupRule          `json:"rules"`
	Relationships   SecurityGroupsRelationships  `json:"relationships"`
	Links           map[string]Link              `json:"links"`
}

type SecurityGroupList struct {
	Pagination Pagination       `json:"pagination,omitempty"`
	Resources  []*SecurityGroup `json:"resources,omitempty"`
}

// SecurityGroupCreate implements an object that is passed to Create method
type SecurityGroupCreate struct {
	Name            string                         `json:"name"`
	GloballyEnabled *SecurityGroupGloballyEnabled  `json:"globally_enabled,omitempty"`
	Rules           []*SecurityGroupRule           `json:"rules,omitempty"`
	Relationships   map[string]ToManyRelationships `json:"relationships,omitempty"`
}

// SecurityGroupUpdate implements an object that is passed to Update method
type SecurityGroupUpdate struct {
	Name            string                        `json:"name,omitempty"`
	GloballyEnabled *SecurityGroupGloballyEnabled `json:"globally_enabled,omitempty"`
	Rules           []*SecurityGroupRule          `json:"rules,omitempty"`
}

// SecurityGroupGloballyEnabled object controls if the group is applied globally to the lifecycle of all applications
type SecurityGroupGloballyEnabled struct {
	Running bool `json:"running"`
	Staging bool `json:"staging"`
}

type SecurityGroupsRelationships struct {
	StagingSpaces ToManyRelationships `json:"staging_spaces"`
	RunningSpaces ToManyRelationships `json:"running_spaces"`
}

// SecurityGroupRule is an object that provide a rule that will be applied by a security group
type SecurityGroupRule struct {
	Protocol    string  `json:"protocol"`
	Destination string  `json:"destination"`
	Ports       *string `json:"ports,omitempty"`
	Type        *int    `json:"type,omitempty"` // https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-types
	Code        *int    `json:"code,omitempty"` // https://www.iana.org/assignments/icmp-parameters/icmp-parameters.xhtml#icmp-parameters-codes
	Description *string `json:"description,omitempty"`
	Log         *bool   `json:"log,omitempty"`
}

func NewSecurityGroupRuleTCP(destination string, enableLogging bool) *SecurityGroupRule {
	return &SecurityGroupRule{
		Protocol:    "tcp",
		Destination: destination,
		Log:         &enableLogging,
	}
}

func NewSecurityGroupRuleUDP(destination string) *SecurityGroupRule {
	return &SecurityGroupRule{
		Protocol:    "udp",
		Destination: destination,
	}
}

func NewSecurityGroupRuleAll(destination string) *SecurityGroupRule {
	return &SecurityGroupRule{
		Protocol:    "all",
		Destination: destination,
	}
}

func NewSecurityGroupRuleICMP(destination string, icmpType, icmpCode int) *SecurityGroupRule {
	return &SecurityGroupRule{
		Protocol:    "icmp",
		Destination: destination,
		Type:        &icmpType,
		Code:        &icmpCode,
	}
}

func (sg *SecurityGroupRule) WithPorts(ports string) *SecurityGroupRule {
	sg.Ports = &ports
	return sg
}

func (sg *SecurityGroupRule) WithDescription(description string) *SecurityGroupRule {
	sg.Description = &description
	return sg
}
