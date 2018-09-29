package ldap

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	l "github.com/go-ldap/ldap"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

var (
	attributes = []string{"*"}
)

var (
	userRegexp          = regexp.MustCompile(",[A-Z]+=")
	escapeFilterRegex   = regexp.MustCompile(`([\\\(\)\*\0-\37\177-\377])`)
	unescapeFilterRegex = regexp.MustCompile(`\\([\da-fA-F]{2}|[()\\*])`) // only match \[)*\] or \xx x=a-fA-F
)

const (
	groupFilter                 = "(cn=%s)"
	userFilter                  = "(%s=%s)"
	userFilterWithObjectClass   = "(&(objectclass=%s)(%s=%s))"
	userDNFilter                = "(%s)"
	userDNFilterWithObjectClass = "(&(objectclass=%s)(%s))"
)

func NewManager(ldapConfig *config.LdapConfig) (Manager, error) {
	conn, err := CreateConnection(ldapConfig)
	if err != nil {
		return nil, err
	}
	return &DefaultManager{
		Config:     ldapConfig,
		Connection: conn,
	}, nil
}

func (m *DefaultManager) GetUserDNs(groupName string) ([]string, error) {
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
	sr, err := m.Connection.Search(search)
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
	return userDNList, nil
}

func (m *DefaultManager) GetUserByDN(userDN string) (*User, error) {
	lo.G.Debug("User DN:", userDN)
	indexes := userRegexp.FindStringIndex(strings.ToUpper(userDN))
	if len(indexes) == 0 {
		return nil, fmt.Errorf("cannot find CN for user DN: %s", userDN)
	}
	index := indexes[0]
	userCNTemp := UnescapeFilterValue(userDN[:index])
	lo.G.Debug("CN unescaped:", userCNTemp)

	userCN := l.EscapeFilter(strings.Replace(userCNTemp, "\\", "", -1))
	lo.G.Debug("CN escaped:", userCN)
	filter := m.getUserFilterWithDN(userCN)
	return m.searchUser(filter, userDN[index+1:], "")
}

func (m *DefaultManager) GetUserByID(userID string) (*User, error) {
	filter := m.getUserFilter(userID)
	lo.G.Debug("Searching for user:", filter)
	lo.G.Debug("Using user search base:", m.Config.UserSearchBase)
	return m.searchUser(filter, m.Config.UserSearchBase, userID)
}

func (m *DefaultManager) searchUser(filter, searchBase, userID string) (*User, error) {
	search := l.NewSearchRequest(
		searchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)

	sr, err := m.Connection.Search(search)
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
		return user, nil
	}
	lo.G.Errorf("Found %d number of entries for filter %s", len(sr.Entries), filter)
	return nil, nil
}

func UnescapeFilterValue(filter string) string {
	repl := unescapeFilterRegex.ReplaceAllFunc(
		[]byte(filter),
		func(match []byte) []byte {
			// \( \) \\ \*
			if len(match) == 2 {
				return []byte{match[1]}
			}
			// had issues with Decode, TODO fix to use Decode?.
			res, _ := hex.DecodeString(string(match[1:]))
			return res
		},
	)
	return string(repl)
}

func (m *DefaultManager) getUserFilter(userID string) string {
	if m.Config.UserObjectClass == "" {
		return fmt.Sprintf(userFilter, m.Config.UserNameAttribute, userID)
	}
	return fmt.Sprintf(userFilterWithObjectClass, m.Config.UserObjectClass, m.Config.UserNameAttribute, userID)
}

func (m *DefaultManager) getUserFilterWithDN(userDN string) string {
	if m.Config.UserObjectClass == "" {
		return fmt.Sprintf(userDNFilter, userDN)
	}
	return fmt.Sprintf(userDNFilterWithObjectClass, m.Config.UserObjectClass, userDN)
}

func (m *DefaultManager) Close() {
	if m.Connection != nil {
		m.Connection.Close()
	}
}
