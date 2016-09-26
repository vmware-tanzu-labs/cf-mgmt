package ldap

import (
	"fmt"
	"io/ioutil"
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

func (m *DefaultManager) GetConfig(configDir, ldapBindPassword string) (config *Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configDir + "/ldap.yml"); err == nil {
		config := &Config{}
		if err = yaml.Unmarshal(data, &config); err == nil {
			if ldapBindPassword != "" {
				config.BindPassword = ldapBindPassword
			} else {
				lo.G.Warning("Ldap bind password should be removed from ldap.yml as this will be deprecated in a future release.  Use --ldap-password flag instead.")
			}
		}
	}
	return
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
	userCN := l.EscapeFilter(unEscapeLDAPValue(userDN[:index]))
	lo.G.Debug("userCN", userCN)
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
		lo.G.Info("Searching for group:", filter)
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
					UserID: entry.GetAttributeValue(config.UserNameAttribute),
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
