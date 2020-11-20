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

type RefreshableConnection struct {
	Connection
	refreshConnection func() (Connection, error)
}

func (r *RefreshableConnection) Search(searchRequest *l.SearchRequest) (*l.SearchResult, error) {
	if r.Connection.IsClosing() {
		err := r.RefreshConnection()
		if err != nil {
			return nil, err
		}
	}
	return r.Connection.Search(searchRequest)
}

func (r *RefreshableConnection) RefreshConnection() error {
	connection, err := r.refreshConnection()
	if err != nil {
		lo.G.Error("Could not re-establish LDAP connection")
		return err
	}

	r.Connection = connection
	return nil
}

// NewRefreshableConnection creates a connection that will use the function
// `createConnection` to refresh the connection if it has been closed.
func NewRefreshableConnection(createConnection func() (Connection, error)) (*RefreshableConnection, error) {
	connection, err := createConnection()

	if err != nil {
		return nil, err
	}

	return &RefreshableConnection{
		Connection:        connection,
		refreshConnection: createConnection,
	}, nil
}

func setMaxTLSVersion(tlsMaxVersion int, tlsConfig *tls.Config) {
	switch tlsMaxVersion {
	case 1:
		tlsConfig.MaxVersion = tls.VersionTLS11
	case 2:
		tlsConfig.MaxVersion = tls.VersionTLS12
	}
}

func createConnection(config *config.LdapConfig) (Connection, error) {
	var connection *l.Conn
	var err error

	ldapURL := fmt.Sprintf("%s:%d", config.LdapHost, config.LdapPort)
	lo.G.Debug("Connecting to", ldapURL)

	if config.TLS {
		if config.TLSMaxVersion != 0 {
			tlsConfig := &tls.Config{}
			setMaxTLSVersion(config.TLSMaxVersion, tlsConfig)
		}
		if config.InsecureSkipVerify == "" || strings.EqualFold(config.InsecureSkipVerify, "true") {
			tlsConfig.InsecureSkipVerify = true
			connection, err = l.DialTLS("tcp", ldapURL, tlsConfig)
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
			tlsConfig.RootCAs = rootCAs
			tlsConfig.ServerName = config.LdapHost

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
