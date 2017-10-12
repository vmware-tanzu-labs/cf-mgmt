package config

// OrgConfig describes configuration for an ASG.

// ASGRule describes the rules in a security group
type ASGRule struct {
	Protocol    string `json:"protocol"`
	Destination string `json:"destination"`
	Ruletype    int    `json:"type"`
	Code        int    `json:"code"`
	Ports       string `json:"ports"`
	Log         bool   `json:"log"`
	Description string `json:"description"`
}

// ASGConfig describes is an array of Rules
type ASGConfig struct {
	//Rules []ASGRule
	Rules string
	Name  string
}

// ASGs contains cf-mgmt configuration for all ASGs.
type ASGs struct {
	ASGs            []string `yaml:"asgs"`
	EnableDeleteASG bool     `yaml:"enable-delete-asg"`
}

// Contains determines whether an ASG is present in a list of ASG.
func (a *ASGs) Contains(asgName string) bool {
	for _, asg := range a.ASGs {
		if asg == asgName {
			return true
		}
	}
	return false
}

/*func (o *ASGConfig) GetBillingManagerGroups() []string {
	return o.BillingManager.groups(o.BillingManagerGroup)
}

func (o *OrgConfig) GetManagerGroups() []string {
	return o.Manager.groups(o.ManagerGroup)
}

func (o *OrgConfig) GetAuditorGroups() []string {
	return o.Auditor.groups(o.AuditorGroup)
}*/
