package config

// UserMgmt specifies users and groups that can be associated to a particular org or space.
type UserMgmt struct {
	LDAPUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	SamlUsers  []string `yaml:"saml_users"`
	LDAPGroup  string   `yaml:"ldap_group,omitempty"`
	LDAPGroups []string `yaml:"ldap_groups"`
}

func (u *UserMgmt) groups(groupName string) []string {
	groupMap := make(map[string]string)
	for _, group := range u.LDAPGroups {
		groupMap[group] = group
	}
	if u.LDAPGroup != "" {
		groupMap[u.LDAPGroup] = u.LDAPGroup
	}
	if groupName != "" {
		groupMap[groupName] = groupName
	}

	result := make([]string, 0, len(groupMap))
	for k := range groupMap {
		result = append(result, k)
	}
	return result
}
