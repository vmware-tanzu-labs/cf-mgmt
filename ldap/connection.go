package ldap

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"

	l "github.com/go-ldap/ldap"
	"github.com/vmwarepivotallabs/cf-mgmt/config"
	"github.com/xchapter7x/lo"
)

type Connection interface {
	Close()
	Search(*l.SearchRequest) (*l.SearchResult, error)
	IsClosing() bool
}

type ConnectionAdapter struct {
	Connection Connection
	Config     *config.LdapConfig
}

func (c ConnectionAdapter) Close() {
	c.Connection.Close()
}

func (c ConnectionAdapter) Search(searchRequest *l.SearchRequest) (*l.SearchResult, error) {
	if c.Connection.IsClosing() {
		err := c.RefreshConnection()
		if err != nil {
			return nil, err
		}
	}
	return c.Connection.Search(searchRequest)
}

func (c ConnectionAdapter) IsClosing() bool {
	return c.Connection.IsClosing()
}

func (c ConnectionAdapter) RefreshConnection() error {
	connection, err := CreateConnection(c.Config)
	if err != nil {
		lo.G.Error("Could not re-establish LDAP connection")
	} else {
		c.Connection = connection
	}
	return err
}

func NewConnectionAdapter(config *config.LdapConfig) (Connection, error) {
	connection, err := CreateConnection(config)

	if err != nil {
		return nil, err
	}

	return ConnectionAdapter{
		Connection: connection,
		Config:     config,
	}, nil
}

func CreateConnection(config *config.LdapConfig) (Connection, error) {
	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)
	var connection *l.Conn
	var err error
	if config.TLS {
		if config.InsecureSkipVerify == "" || strings.EqualFold(config.InsecureSkipVerify, "true") {
			connection, err = l.DialTLS("tcp", ldapURL, &tls.Config{InsecureSkipVerify: true})
		} else {
			// Get the SystemCertPool, continue with an empty pool on error
			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}

			// Append our cert to the system pool
			if ok := rootCAs.AppendCertsFromPEM([]byte(config.CACert)); !ok {
				log.Println("No certs appended, using system certs only")
			}

			// Trust the augmented cert pool in our client
			tlsConfig := &tls.Config{
				RootCAs:    rootCAs,
				ServerName: config.LdapHost,
			}

			connection, err = l.DialTLS("tcp", ldapURL, tlsConfig)
		}
	} else {
		connection, err = l.Dial("tcp", ldapURL)
	}
	if err != nil {
		return nil, err
	}
	if connection != nil {
		if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
			connection.Debug = true
		}
		if err = connection.Bind(config.BindDN, config.BindPassword); err != nil {
			connection.Close()
			return nil, fmt.Errorf("cannot bind with %s: %v", config.BindDN, err)
		}
	}
	return connection, err

}
