package ldap

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	l "github.com/go-ldap/ldap"
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

func NewManager() Manager {
	return &DefaultManager{}
}

func (m *DefaultManager) GetConfig(configDir, ldapBindPassword string) (*Config, error) {
	data, err := ioutil.ReadFile(configDir + "/ldap.yml")
	if err != nil {
		lo.G.Error(err)
		return nil, err
	}
	config := &Config{}
	if err = yaml.Unmarshal(data, &config); err != nil {
		lo.G.Error(err)
		return nil, err
	}
	if ldapBindPassword != "" {
		config.BindPassword = ldapBindPassword
	} else {
		lo.G.Warning("Ldap bind password should be removed from ldap.yml as this will be deprecated in a future release.  Use --ldap-password flag instead.")
	}
	if config.Origin == "" {
		config.Origin = "ldap"
	}
	return config, nil
}

func (m *DefaultManager) LdapConnection(config *Config) (*l.Conn, error) {
	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	var connection *l.Conn
	var err error
	if config.TLS {
		connection, err = l.DialTLS("tcp", ldapURL, &tls.Config{InsecureSkipVerify: true})
	} else {
		connection, err = l.Dial("tcp", ldapURL)
	}
	if err != nil {
		return nil, err
	}
	if connection != nil {
		if err = connection.Bind(config.BindDN, config.BindPassword); err != nil {
			connection.Close()
			return nil, fmt.Errorf("cannot bind with %s: %v", config.BindDN, err)
		}
	}
	return connection, err

}

//GetUserIDs -
func (m *DefaultManager) GetUserIDs(config *Config, groupName string) ([]User, error) {
	ldapConnection, err := m.LdapConnection(config)
	if err != nil {
		return nil, err
	}
	defer ldapConnection.Close()

	groupEntry, err := m.getGroup(ldapConnection, groupName, config.GroupSearchBase)
	if err != nil || groupEntry == nil {
		lo.G.Info("group not found:", groupName)
		return nil, err
	}

	userDNList := groupEntry.GetAttributeValues(config.GroupAttribute)
	if len(userDNList) == 0 {
		lo.G.Info("No users found under group:", config.GroupAttribute)
		return nil, nil
	}

	var users []User
	for _, userDN := range userDNList {
		user, err := m.GetLdapUser(config, userDN)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, *user)
		} else {
			lo.G.Infof("User entry: %s not found", userDN)
		}
	}
	return users, nil
}

func (m *DefaultManager) GetLdapUser(config *Config, userDN string) (*User, error) {
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
	filter := m.getUserFilterWithDN(config, userCN)
	return m.searchUser(filter, userDN[index+1:], config)
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

func (m *DefaultManager) GetUser(config *Config, userID string) (*User, error) {
	filter := m.getUserFilter(config, userID)
	return m.searchUser(filter, config.UserSearchBase, config)
}

func (m *DefaultManager) searchUser(filter, searchBase string, config *Config) (*User, error) {
	lo.G.Debug("Searching for user:", filter)
	lo.G.Debug("Using user search base:", searchBase)
	ldapConnection, err := m.LdapConnection(config)
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
			UserID: entry.GetAttributeValue(config.UserNameAttribute),
			Email:  entry.GetAttributeValue(config.UserMailAttribute),
		}
		return user, nil
	}
	lo.G.Infof("Found %d number of entries for filter %s", len(sr.Entries), filter)
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

func (m *DefaultManager) getUserFilter(config *Config, userID string) string {
	if config.UserObjectClass == "" {
		return fmt.Sprintf(userFilter, config.UserNameAttribute, userID)
	}
	return fmt.Sprintf(userFilterWithObjectClass, config.UserObjectClass, config.UserNameAttribute, userID)
}

func (m *DefaultManager) getUserFilterWithDN(config *Config, userDN string) string {
	if config.UserObjectClass == "" {
		return fmt.Sprintf(userDNFilter, userDN)
	}
	return fmt.Sprintf(userDNFilterWithObjectClass, config.UserObjectClass, userDN)
}
