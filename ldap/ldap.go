package ldap

import (
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	l "github.com/nmcclain/ldap"
	"github.com/xchapter7x/lo"
)

var (
	attributes  = []string{"*"}
	groupFilter = "(cn=%s)"
	userFilter  = "(%s)"
)

//NewDefaultManager -
func NewDefaultManager(configDir string) (mgr Manager, err error) {
	var config Config
	m := &DefaultManager{}
	if config, err = m.getConfig(configDir); err == nil {
		m.Config = config
	}
	mgr = m
	return
}

func (m *DefaultManager) getConfig(configDir string) (config Config, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(configDir + "/ldap.yml"); err == nil {
		c := &Config{}
		if err = yaml.Unmarshal(data, c); err == nil {
			config = *c
		}
	}
	return
}

//GetUserIDs -
func (m *DefaultManager) GetUserIDs(groupName string) (users []User, err error) {
	var groupEntry *l.Entry
	var userEntry *l.Entry
	if groupEntry, err = m.getGroup(groupName); err == nil {
		if groupEntry != nil {
			userDNList := groupEntry.GetAttributeValues(m.Config.GroupAttribute)
			for _, userDN := range userDNList {
				if userEntry, err = m.getUser(userDN); err == nil {
					if userEntry != nil {
						user := User{
							UserDN: userEntry.DN,
							UserID: userEntry.GetAttributeValue(m.Config.UserNameAttribute),
							Email:  userEntry.GetAttributeValue(m.Config.UserMailAttribute),
						}
						users = append(users, user)
					} else {
						lo.G.Info("User entry not found", userDN)
					}
				}
			}
		} else {
			lo.G.Info("Group not found", groupName)
		}
	}
	return
}

func (m *DefaultManager) getUser(userDN string) (entry *l.Entry, err error) {
	var ldapConnection *l.Conn
	var sr *l.SearchResult
	ldapURL := fmt.Sprintf("%s:%d", m.Config.LdapHost, m.Config.LdapPort)
	lo.G.Info("Connecting to", ldapURL)
	if ldapConnection, err = l.Dial("tcp", ldapURL); err == nil {
		// be sure to add error checking!
		defer ldapConnection.Close()
		if err = ldapConnection.Bind(m.Config.BindDN, m.Config.BindPassword); err == nil {
			//filter := fmt.Sprintf(userFilter, userObjectClass, userID)
			lo.G.Debug("User DN:", userDN)
			userCNTemp := strings.Replace(userDN, ","+m.Config.UserSearchBase, "", 1)
			userCN := strings.Replace(userCNTemp, "\\,", ",", 1)
			filter := fmt.Sprintf(userFilter, userCN)
			lo.G.Debug("Using user search filter", filter)
			search := l.NewSearchRequest(
				m.Config.UserSearchBase,
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
		} else {
			lo.G.Error(err)
		}
	} else {
		lo.G.Error(err)
	}
	return
}
func (m *DefaultManager) getGroup(groupName string) (entry *l.Entry, err error) {
	var ldapConnection *l.Conn
	var sr *l.SearchResult
	ldapURL := fmt.Sprintf("%s:%d", m.Config.LdapHost, m.Config.LdapPort)
	lo.G.Info("Connecting to", ldapURL)
	if ldapConnection, err = l.Dial("tcp", ldapURL); err == nil {
		// be sure to add error checking!
		defer ldapConnection.Close()
		filter := fmt.Sprintf(groupFilter, groupName)
		lo.G.Debug("Using group filter", filter)
		lo.G.Debug("Using group search base:", m.Config.GroupSearchBase)
		if err = ldapConnection.Bind(m.Config.BindDN, m.Config.BindPassword); err == nil {
			search := l.NewSearchRequest(
				m.Config.GroupSearchBase,
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
		} else {
			lo.G.Error(err)
		}
	} else {
		lo.G.Error(err)
	}
	return
}
