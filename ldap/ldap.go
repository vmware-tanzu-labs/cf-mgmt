package ldap

import (
	"fmt"
	"strings"

	l "github.com/go-ldap/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

var (
	attributes = []string{"*"}
)

const (
	groupFilter                 = "(cn=%s)"
	groupFilterWithObjectClass  = "(&(objectclass=%s)(%s))"
	userFilter                  = "(%s=%s)"
	userFilterWithObjectClass   = "(&(objectclass=%s)(%s=%s))"
	userDNFilter                = "(%s)"
	userDNFilterWithObjectClass = "(&(objectclass=%s)(%s))"
)

func NewManager(ldapConfig *config.LdapConfig) (*Manager, error) {
	conn, err := CreateConnection(ldapConfig)
	if err != nil {
		return nil, err
	}
	return &Manager{
		Config:     ldapConfig,
		Connection: conn,
		userMap:    make(map[string]*User),
		groupMap:   make(map[string][]string),
	}, nil
}

func (m *Manager) groupInCache(groupName string) bool {
	if m.groupMap == nil {
		return false
	}
	if _, ok := m.groupMap[groupName]; ok {
		return true
	}

	return false
}

func (m *Manager) addGroupToCache(groupName string, result []string) {
	if m.groupMap == nil {
		m.groupMap = make(map[string][]string)
	}
	m.groupMap[groupName] = result
}

func (m *Manager) userInCache(userFilter string) bool {
	if m.userMap == nil {
		return false
	}
	if _, ok := m.userMap[userFilter]; ok {
		return true
	}
	return false
}

func (m *Manager) addUserToCache(userFilter string, result *User) {
	if m.userMap == nil {
		m.userMap = make(map[string]*User)
	}
	m.userMap[userFilter] = result
}

func (m *Manager) GetUserDNs(groupName string) ([]string, error) {
	if m.groupInCache(groupName) {
		lo.G.Debugf("Group %s found in cache", groupName)
		return m.groupMap[groupName], nil
	}
	filter := fmt.Sprintf(groupFilter, l.EscapeFilter(groupName))
	var groupEntry *l.Entry
	lo.G.Debug("Searching for group:", filter)
	lo.G.Debug("Using group search base:", m.Config.GroupSearchBase)

	search := l.NewSearchRequest(
		m.Config.GroupSearchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)
	sr, err := m.doSearch(search)
	if err != nil {
		lo.G.Error(err)
		return nil, err
	}

	if len(sr.Entries) == 0 {
		lo.G.Errorf("group not found: %s", groupName)
		return []string{}, nil
	}
	if len(sr.Entries) > 1 {
		lo.G.Errorf("multiple groups found for: %s", groupName)
		return []string{}, nil
	}

	groupEntry = sr.Entries[0]
	userDNList := groupEntry.GetAttributeValues(m.Config.GroupAttribute)
	if len(userDNList) == 0 {
		lo.G.Warningf("No users found under group: %s", groupName)
	}

	userMap := make(map[string]string)
	for _, userDN := range userDNList {
		isGroup, nestedGroupName, err := m.IsGroup(userDN)
		if err != nil {
			return nil, err
		}
		if isGroup {
			if err != nil {
				return nil, err
			}
			nestedUsers, err := m.GetUserDNs(nestedGroupName)
			if err != nil {
				return nil, err
			}
			for _, nestedUser := range nestedUsers {
				userMap[nestedUser] = nestedUser
			}
		} else {
			userMap[userDN] = userDN
		}
	}
	var userList []string
	for _, userDN := range userMap {
		userList = append(userList, userDN)
	}
	m.addGroupToCache(groupName, userList)
	return userList, nil
}

func (m *Manager) doSearch(search *l.SearchRequest) (*l.SearchResult, error) {
	if m.Connection.IsClosing() {
		conn, err := CreateConnection(m.Config)
		if err != nil {
			return nil, err
		}
		m.Connection = conn
	}
	sr, err := m.Connection.Search(search)
	if err != nil {
		return nil, err
	}
	return sr, nil
}

func (m *Manager) GroupFilter(userDN string) (string, error) {
	cn, _, err := ParseUserCN(userDN)
	if err != nil {
		return "", err
	}
	cnTemp := UnescapeFilterValue(cn)
	lo.G.Debug("CN unescaped:", cnTemp)

	escapedCN := l.EscapeFilter(strings.Replace(cnTemp, "\\", "", -1))
	lo.G.Debug("CN escaped:", escapedCN)
	groupObjectFilter := "groupOfNames"
	if m.Config.GroupObjectClass != "" {
		groupObjectFilter = m.Config.GroupObjectClass
	}
	return fmt.Sprintf(groupFilterWithObjectClass, groupObjectFilter, escapedCN), nil
}
func (m *Manager) IsGroup(DN string) (bool, string, error) {
	if strings.Contains(DN, m.Config.GroupSearchBase) {
		filter, err := m.GroupFilter(DN)
		if err != nil {
			return false, "", err
		}
		search := l.NewSearchRequest(
			m.Config.GroupSearchBase,
			l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
			filter,
			attributes,
			nil)
		sr, err := m.doSearch(search)
		if err != nil {
			return false, "", err
		}
		lo.G.Debugf("Found %d entries for group filter %s", len(sr.Entries), filter)
		if len(sr.Entries) == 1 {
			return true, sr.Entries[0].GetAttributeValue("cn"), nil
		}
		return false, "", nil
	} else {
		return false, "", nil
	}
}

func (m *Manager) GetUserByDN(userDN string) (*User, error) {
	cn, searchBase, err := ParseUserCN(userDN)
	if err != nil {
		return nil, err
	}
	userCNTemp := UnescapeFilterValue(cn)
	lo.G.Debug("CN unescaped:", userCNTemp)

	userCN := EscapeFilterValue(userCNTemp)
	lo.G.Debug("CN escaped:", userCN)

	filter := m.getUserFilterWithCN(userCN)
	return m.searchUser(filter, searchBase, "")
}

func (m *Manager) GetUserByID(userID string) (*User, error) {
	filter := m.getUserFilter(userID)
	return m.searchUser(filter, m.Config.UserSearchBase, userID)
}

func (m *Manager) searchUser(filter, searchBase, userID string) (*User, error) {
	if m.userInCache(filter) {
		lo.G.Debugf("User with filter %s found in cache", filter)
		return m.userMap[filter], nil
	}
	lo.G.Debugf("Searching with filter [%s]", filter)
	lo.G.Debugf("Using user search base: [%s]", searchBase)
	search := l.NewSearchRequest(
		searchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)

	sr, err := m.doSearch(search)
	if err != nil {
		lo.G.Error(err)
		return nil, err
	}

	if (len(sr.Entries)) == 1 {
		entry := sr.Entries[0]
		user := &User{
			UserDN: entry.DN,
			Email:  entry.GetAttributeValue(m.Config.UserMailAttribute),
		}
		if userID != "" {
			user.UserID = userID
		} else {
			user.UserID = entry.GetAttributeValue(m.Config.UserNameAttribute)
		}
		lo.G.Debugf("Search filter %s returned userDN [%s], email [%s], userID [%s]", filter, user.UserDN, user.Email, user.UserID)
		m.addUserToCache(filter, user)
		return user, nil
	}
	lo.G.Errorf("Found %d number of entries for filter %s", len(sr.Entries), filter)
	return nil, nil
}

func mustEscape(c byte) bool {
	return c > 0x7f || c == '(' || c == ')' || c == '\\' || c == '*' || c == 0
}

func (m *Manager) getUserFilter(userID string) string {
	if m.Config.UserObjectClass == "" {
		return fmt.Sprintf(userFilter, m.Config.UserNameAttribute, userID)
	}
	return fmt.Sprintf(userFilterWithObjectClass, m.Config.UserObjectClass, m.Config.UserNameAttribute, userID)
}

func (m *Manager) getUserFilterWithCN(cn string) string {
	if m.Config.UserObjectClass == "" {
		return fmt.Sprintf(userDNFilter, cn)
	}
	return fmt.Sprintf(userDNFilterWithObjectClass, m.Config.UserObjectClass, cn)
}

func (m *Manager) Close() {
	if m.Connection != nil {
		m.Connection.Close()
		m.userMap = nil
		m.groupMap = nil
	}
}
