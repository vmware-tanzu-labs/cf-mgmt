package ldap

import (
	"crypto/tls"
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

func NewManager(cfg config.Reader, ldapBindPassword string) (Manager, error) {
	ldapConfig, err := cfg.LdapConfig(ldapBindPassword)
	if err != nil {
		return nil, err
	}
	return &DefaultManager{
		Config: ldapConfig,
	}, nil
}

func (m *DefaultManager) LdapConnection() (*l.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%d", m.Config.LdapHost, m.Config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	var connection *l.Conn
	var err error
	if m.Config.TLS {
		connection, err = l.DialTLS("tcp", ldapURL, &tls.Config{InsecureSkipVerify: true})
	} else {
		connection, err = l.Dial("tcp", ldapURL)
	}
	if err != nil {
		return nil, err
	}
	if connection != nil {
		if err = connection.Bind(m.Config.BindDN, m.Config.BindPassword); err != nil {
			connection.Close()
			return nil, fmt.Errorf("cannot bind with %s: %v", m.Config.BindDN, err)
		}
	}
	return connection, err

}

//GetUserIDs -
func (m *DefaultManager) GetUserIDs(groupName string) ([]User, error) {
	ldapConnection, err := m.LdapConnection()
	if err != nil {
		return nil, err
	}
	defer ldapConnection.Close()

	groupEntry, err := m.getGroup(ldapConnection, groupName, m.Config.GroupSearchBase)
	if err != nil || groupEntry == nil {
		lo.G.Errorf("group not found: %s", groupName)
		return nil, err
	}

	userDNList := groupEntry.GetAttributeValues(m.Config.GroupAttribute)
	if len(userDNList) == 0 {
		lo.G.Warningf("No users found under group: %s", groupName)
		return nil, nil
	}

	var users []User
	for _, userDN := range userDNList {
		user, err := m.GetLdapUser(userDN)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, *user)
		} else {
			lo.G.Warningf("User entry: %s not found", userDN)
		}
	}
	return users, nil
}

func (m *DefaultManager) GetLdapUser(userDN string) (*User, error) {
	lo.G.Debug("User DN:", userDN)
	indexes := userRegexp.FindStringIndex(strings.ToUpper(userDN))
	if len(indexes) == 0 {
		return nil, fmt.Errorf("cannot find CN for user DN: %s", userDN)
	}
	index := indexes[0]
	userCNTemp := m.UnescapeFilterValue(userDN[:index])
	lo.G.Debug("CN unescaped:", userCNTemp)

	userCN := l.EscapeFilter(strings.Replace(userCNTemp, "\\", "", -1))
	lo.G.Debug("CN escaped:", userCN)
	filter := m.getUserFilterWithDN(userCN)
	return m.searchUser(filter, userDN[index+1:], "")
}

func (m *DefaultManager) getGroup(ldapConnection *l.Conn, groupName, groupSearchBase string) (*l.Entry, error) {
	filter := fmt.Sprintf(groupFilter, l.EscapeFilter(groupName))

	lo.G.Debug("Searching for group:", filter)
	lo.G.Debug("Using group search base:", groupSearchBase)

	search := l.NewSearchRequest(
		groupSearchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)
	sr, err := ldapConnection.Search(search)
	if err != nil {
		lo.G.Error(err)
		return nil, err
	}
	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("group not found: %s", groupName)
	}
	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("multiple groups found for: %s", groupName)
	}
	return sr.Entries[0], nil
}

func (m *DefaultManager) GetUser(userID string) (*User, error) {
	filter := m.getUserFilter(userID)
	return m.searchUser(filter, m.Config.UserSearchBase, userID)
}

func (m *DefaultManager) searchUser(filter, searchBase, userID string) (*User, error) {
	lo.G.Debug("Searching for user:", filter)
	lo.G.Debug("Using user search base:", searchBase)
	ldapConnection, err := m.LdapConnection()
	if err != nil {
		return nil, err
	}
	defer ldapConnection.Close()
	search := l.NewSearchRequest(
		searchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)

	sr, err := ldapConnection.Search(search)
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

func (m *DefaultManager) EscapeFilterValue(filter string) string {
	repl := escapeFilterRegex.ReplaceAllFunc(
		[]byte(filter),
		func(match []byte) []byte {
			if len(match) == 2 {
				return []byte(fmt.Sprintf("\\%02x", match[1]))
			}
			return []byte(fmt.Sprintf("\\%02x", match[0]))
		},
	)
	return string(repl)
}
func (m *DefaultManager) UnescapeFilterValue(filter string) string {
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

func (m *DefaultManager) GetLdapUsers(groupNames []string, userList []string) ([]User, error) {
	uniqueUsers := make(map[string]string)
	users := []User{}
	for _, groupName := range groupNames {
		if groupName != "" {
			lo.G.Debug("Finding LDAP user for group:", groupName)
			if groupUsers, err := m.GetUserIDs(groupName); err == nil {
				for _, user := range groupUsers {
					if _, ok := uniqueUsers[strings.ToLower(user.UserDN)]; !ok {
						users = append(users, user)
						uniqueUsers[strings.ToLower(user.UserDN)] = user.UserDN
					} else {
						lo.G.Debugf("User %+v is already added to list", user)
					}
				}
			} else {
				lo.G.Warning(err)
			}
		}
	}
	for _, user := range userList {
		if ldapUser, err := m.GetUser(user); err == nil {
			if ldapUser != nil {
				if _, ok := uniqueUsers[strings.ToLower(ldapUser.UserDN)]; !ok {
					users = append(users, *ldapUser)
					uniqueUsers[strings.ToLower(ldapUser.UserDN)] = ldapUser.UserDN
				} else {
					lo.G.Debugf("User %+v is already added to list", ldapUser)
				}
			}
		} else {
			lo.G.Warning(err)
		}
	}
	return users, nil
}

func (m *DefaultManager) LdapConfig() *config.LdapConfig {
	return m.Config
}
