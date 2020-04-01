package config

// UserMgmt specifies users and groups that can be associated to a particular org or space.
type UserMgmt struct {
	LDAPUsers  []string `yaml:"ldap_users"`
	Users      []string `yaml:"users"`
	SamlUsers  []string `yaml:"saml_users"`
	LDAPGroup  string   `yaml:"ldap_group,omitempty"`
	LDAPGroups []string `yaml:"ldap_groups"`
}

// UserOrigin is an enum type encoding from what source a user originated.
// Choices are: internal, saml, ldap. If you give a UserOrigin value that lies
// outside of these options, the behaviour is undefined.
type UserOrigin int

const (
	// InternalOrigin corresponds to a UAA user
	InternalOrigin UserOrigin = iota

	// SAMLOrigin corresponds to a SAML backed user
	SAMLOrigin

	// LDAPOrigin corresponds to a LDAP backed user
	LDAPOrigin
)

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

func (u *UserMgmt) hasUser(origin UserOrigin, username string) bool {
	var userList []string
	switch origin {
	case InternalOrigin:
		userList = u.Users
	case SAMLOrigin:
		userList = u.SamlUsers
	case LDAPOrigin:
		userList = u.LDAPUsers
	}

	for _, existentUser := range userList {
		if username == existentUser {
			return true
		}
	}

	return false
}

func (u *UserMgmt) addUser(origin UserOrigin, username string) {
	if u.hasUser(origin, username) {
		return
	}

	var userList *[]string

	switch origin {
	case InternalOrigin:
		userList = &u.Users
	case SAMLOrigin:
		userList = &u.SamlUsers
	case LDAPOrigin:
		userList = &u.LDAPUsers
	}

	*userList = append(*userList, username)
}
