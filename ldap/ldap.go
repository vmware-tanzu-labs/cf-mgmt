package ldap

import (
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

//GetUserIDs -
func (m *DefaultManager) GetUserIDs(config *Config, groupName string) (users []User, err error) {
	var ldapConnection *l.Conn
	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	if ldapConnection, err = l.Dial("tcp", ldapURL); err == nil {
		// be sure to add error checking!
		defer ldapConnection.Close()
		if err = ldapConnection.Bind(config.BindDN, config.BindPassword); err != nil {
			return
		}
		var groupEntry *l.Entry
		var user *User
		if groupEntry, err = m.getGroup(ldapConnection, groupName, config.GroupSearchBase); err == nil {
			if groupEntry != nil {
				userDNList := groupEntry.GetAttributeValues(config.GroupAttribute)
				for _, userDN := range userDNList {
					if user, err = m.getLdapUser(ldapConnection, userDN, config.UserSearchBase, config.UserNameAttribute, config.UserMailAttribute); err == nil {
						if user != nil {
							users = append(users, *user)
						} else {
							lo.G.Info("User entry not found", userDN)
						}
					}
				}
			} else {
				lo.G.Info("Group not found", groupName)
			}
		}
	}
	return
}

func (m *DefaultManager) GetLdapUser(config *Config, userDN, userSearchBase string) (*User, error) {
	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	if ldapConnection, err := l.Dial("tcp", ldapURL); err == nil {
		// be sure to add error checking!
		defer ldapConnection.Close()
		if err := ldapConnection.Bind(config.BindDN, config.BindPassword); err != nil {
			return nil, err
		}
		if user, err := m.getLdapUser(ldapConnection, userDN, config.UserSearchBase, config.UserNameAttribute, config.UserMailAttribute); err == nil {
			return user, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (m *DefaultManager) getLdapUser(ldapConnection *l.Conn, userDN, userSearchBase, userNameAttribute, userMailAttribute string) (user *User, err error) {
	var sr *l.SearchResult
	lo.G.Debug("User DN:", userDN)
	index := strings.Index(strings.ToUpper(userDN), ",OU=")
	userCNTemp := m.UnescapeFilterValue(userDN[:index])
	lo.G.Debug("CN unescaped:", userCNTemp)
	userCN := l.EscapeFilter(strings.Replace(userCNTemp, "\\", "", -1))
	//userCN := l.EscapeFilter(unEscapeLDAPValue(userDN[:index]))
	lo.G.Debug("CN escaped", userCN)
	filter := fmt.Sprintf(userFilter, userCN)
	lo.G.Info("Searching for user:", filter)
	search := l.NewSearchRequest(
		userSearchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)

	if sr, err = ldapConnection.Search(search); err == nil {
		if (len(sr.Entries)) == 1 {
			userEntry := sr.Entries[0]
			user = &User{
				UserDN: userEntry.DN,
				UserID: userEntry.GetAttributeValue(userNameAttribute),
				Email:  userEntry.GetAttributeValue(userMailAttribute),
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
	lo.G.Info("Searching for group:", filter)
	lo.G.Debug("Using group search base:", groupSearchBase)

	search := l.NewSearchRequest(
		groupSearchBase,
		l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
		filter,
		attributes,
		nil)
	if sr, err = ldapConnection.Search(search); err == nil {
		if (len(sr.Entries)) == 1 {
			entry = sr.Entries[0]
		}
	} else {
		lo.G.Error(err)
	}

	return
}

func (m *DefaultManager) GetUser(config *Config, userID string) (*User, error) {

	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	if ldapConnection, err := l.Dial("tcp", ldapURL); err != nil {
		return nil, err
	} else {
		// be sure to add error checking!
		defer ldapConnection.Close()
		if err := ldapConnection.Bind(config.BindDN, config.BindPassword); err != nil {
			lo.G.Error(err)
			return nil, err
		}
		theUserFilter := "(" + config.UserNameAttribute + "=%s)"
		lo.G.Debug("User filter before escape:", theUserFilter)
		filter := fmt.Sprintf(theUserFilter, l.EscapeFilter(userID))
		lo.G.Info("Searching for user:", filter)
		lo.G.Debug("Using user search base:", config.UserSearchBase)

		search := l.NewSearchRequest(
			config.UserSearchBase,
			l.ScopeWholeSubtree, l.NeverDerefAliases, 0, 0, false,
			filter,
			attributes,
			nil)
		if sr, err := ldapConnection.Search(search); err == nil {
			lo.G.Debug(fmt.Sprintf("Found %d number of entries for filter %s", len(sr.Entries), filter))
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

func unEscapeLDAPValue(input string) string {
	var returnString string
	returnString = strings.Replace(input, "2C", ",", 1)
	returnString = strings.Replace(returnString, "\\,", ",", 1)
	return returnString
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
