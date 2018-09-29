package ldap

import (
	"crypto/tls"
	"fmt"

	l "github.com/go-ldap/ldap"
	"github.com/pivotalservices/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type Connection interface {
	Close()
	Search(*l.SearchRequest) (*l.SearchResult, error)
}

func CreateConnection(config *config.LdapConfig) (Connection, error) {
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
