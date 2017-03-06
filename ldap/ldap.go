package ldap

import (
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"errors"

	l "github.com/go-ldap/ldap"
	"github.com/op/go-logging"
	"github.com/xchapter7x/lo"
)

var (
	attributes  = []string{"*"}
	groupFilter = "(cn=%s)"
	userFilter  = "(%s)"
)

func NewManager() Manager {
	return &DefaultManager{}
}

func (m *DefaultManager) GetConfig(configDir, ldapBindPassword string) (*Config, error) {

	if data, err := ioutil.ReadFile(configDir + "/ldap.yml"); err == nil {
		config := &Config{}
		if err = yaml.Unmarshal(data, &config); err == nil {
			if ldapBindPassword != "" {
				config.BindPassword = ldapBindPassword
			} else {
				lo.G.Warning("Ldap bind password should be removed from ldap.yml as this will be deprecated in a future release.  Use --ldap-password flag instead.")
			}
			if config.Origin == "" {
				config.Origin = "ldap"
			}
			return config, nil
		} else {
			lo.G.Error(err)
			return nil, err
		}
	} else {
		lo.G.Error(err)
		return nil, err
	}
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
	if connection != nil {
		if logging.GetLevel(lo.LOG_MODULE) == logging.DEBUG {
			connection.Debug = true
		}
		if err = connection.Bind(config.BindDN, config.BindPassword); err != nil {
			connection.Close()
			lo.G.Error(fmt.Sprintf("Error binding with %s", config.BindDN), err)
			return nil, err
		}
	}
	return connection, err

}

//GetUserIDs -
func (m *DefaultManager) GetUserIDs(config *Config, groupName string) (users []User, err error) {
	var ldapConnection *l.Conn
	if ldapConnection, err = m.LdapConnection(config); err == nil {
		defer ldapConnection.Close()
		var groupEntry *l.Entry
		var user *User
		if groupEntry, err = m.getGroup(ldapConnection, groupName, config.GroupSearchBase); err == nil {
			if groupEntry != nil {
				userDNList := groupEntry.GetAttributeValues(config.GroupAttribute)
				if len(userDNList) == 0 {
					lo.G.Info("No users found under group : ", config.GroupAttribute)
				} else {
					for _, userDN := range userDNList {
						lo.G.Info("Getting details about user : ", userDN)
						if user, err = m.getLdapUser(ldapConnection, userDN, config); err == nil {
							if user != nil {
								users = append(users, *user)
							} else {
								lo.G.Info("User entry not found", userDN)
							}
						}
					}
				}
			} else {
				lo.G.Info("Group not found : ", groupName)
			}
		}
	} else {
		return nil, err
	}
	return
}

func (m *DefaultManager) GetLdapUser(config *Config, userDN, userSearchBase string) (*User, error) {
	if ldapConnection, err := m.LdapConnection(config); err == nil {
		defer ldapConnection.Close()
		if user, err := m.getLdapUser(ldapConnection, userDN, config); err == nil {
			return user, nil
		} else {
			lo.G.Info("User not found :", user)
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (m *DefaultManager) getLdapUser(ldapConnection *l.Conn, userDN string, config *Config) (user *User, err error) {
	var sr *l.SearchResult
	lo.G.Debug("User DN:", userDN)
	regex, _ := regexp.Compile(",[A-Z]+=")
	indexes := regex.FindStringIndex(strings.ToUpper(userDN))
	if len(indexes) == 0 {
		return nil, errors.New(fmt.Sprintf("%s %s ", "Can't find CN for user DN:", userDN))
	}
	index := indexes[0]
	userCNTemp := m.UnescapeFilterValue(userDN[:index])
	lo.G.Debug("CN unescaped:", userCNTemp)
	userCN := l.EscapeFilter(strings.Replace(userCNTemp, "\\", "", -1))
	//userCN := l.EscapeFilter(unEscapeLDAPValue(userDN[:index]))
	lo.G.Debug("CN escaped", userCN)
	filter := fmt.Sprintf(userFilter, userCN)
	lo.G.Info("Searching for user:", filter)
	search := l.NewSearchRequest(
		config.UserSearchBase,
		l.ScopeWholeSubtree,
		l.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		attributes,
		nil)

	if sr, err = ldapConnection.Search(search); err == nil {
		if (len(sr.Entries)) == 1 {
			userEntry := sr.Entries[0]
			user = &User{
				UserDN: userEntry.DN,
				UserID: userEntry.GetAttributeValue(config.UserNameAttribute),
				Email:  userEntry.GetAttributeValue(config.UserMailAttribute),
			}
		}
	} else {
		lo.G.Error(err)
	}
	return
}
func (m *DefaultManager) getGroup(ldapConnection *l.Conn, groupName, groupSearchBase string) (entry *l.Entry, err error) {

	var sr *l.SearchResult
	filter := fmt.Sprintf(groupFilter, l.EscapeFilter(groupName))

	lo.G.Info("Searching for group : ", filter)

	lo.G.Info("Using group search base : ", groupSearchBase)

	search := l.NewSearchRequest(
		groupSearchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)
	if sr, err = ldapConnection.Search(search); err == nil {
		if (len(sr.Entries)) == 1 {
			entry = sr.Entries[0]
		} else {
			lo.G.Info("Group not found : ", groupName)
		}
	} else {
		lo.G.Error(err)
	}

	return
}

func (m *DefaultManager) GetUser(config *Config, userID string) (*User, error) {
	if ldapConnection, err := m.LdapConnection(config); err != nil {
		return nil, err
	} else {
		defer ldapConnection.Close()

		theUserFilter := "(" + config.UserNameAttribute + "=%s)"

		lo.G.Debug("User filter before escape:", theUserFilter)

		filter := fmt.Sprintf(theUserFilter, l.EscapeFilter(userID))

		lo.G.Info("Searching for user:", filter)

		lo.G.Info("Using user search base:", config.UserSearchBase)

		search := l.NewSearchRequest(
			config.UserSearchBase,
			l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
			filter,
			attributes,
			nil)
		if sr, err := ldapConnection.Search(search); err == nil {

			lo.G.Info(fmt.Sprintf("Found %d number of entries for filter %s", len(sr.Entries), filter))

			if (len(sr.Entries)) == 1 {
				entry := sr.Entries[0]
				user := &User{
					UserDN: entry.DN,
					UserID: userID,
					Email:  entry.GetAttributeValue(config.UserMailAttribute),
				}
				return user, nil
			}
		} else {
			lo.G.Error(err)
			return nil, err
		}
	}
	return nil, nil
}

func (m *DefaultManager) EscapeFilterValue(filter string) string {
	var escapeFilterRegex *regexp.Regexp
	escapeFilterRegex = regexp.MustCompile(`([\\\(\)\*\0-\37\177-\377])`)
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
	var unescapeFilterRegex *regexp.Regexp
	unescapeFilterRegex = regexp.MustCompile(`\\([\da-fA-F]{2}|[()\\*])`)
	// regex wil only match \[)*\] or \xx x=a-fA-F
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
