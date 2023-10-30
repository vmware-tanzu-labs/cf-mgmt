package config

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const UserAgent = "Go-CF-Client/3.0"

// Config is used to configure the creation of a client
type Config struct {
	APIEndpointURL   string
	LoginEndpointURL string
	UAAEndpointURL   string

	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	UserAgent    string
	Origin       string
	AccessToken  string
	RefreshToken string

	baseHTTPClient    *http.Client
	requestTimeout    time.Duration
	skipTLSValidation bool
}

type cfHomeConfig struct {
	AccessToken           string
	RefreshToken          string
	Target                string
	AuthorizationEndpoint string
	OrganizationFields    struct {
		Name string
	}
	SpaceFields struct {
		Name string
	}
	SSLDisabled bool
}

// NewUserPassword creates a new config configured for regular user/password authentication
func NewUserPassword(apiRootURL, username, password string) (*Config, error) {
	if username == "" {
		return nil, errors.New("expected an non-empty CF API username")
	}
	if password == "" {
		return nil, errors.New("expected an non-empty CF API password")
	}

	c, err := newDefault(apiRootURL)
	if err != nil {
		return nil, err
	}
	c.Username = username
	c.Password = password

	return c, nil
}

// NewClientSecret creates a new config configured for client id and client secret authentication
func NewClientSecret(apiRoot, clientID, clientSecret string) (*Config, error) {
	if clientID == "" {
		return nil, errors.New("expected an non-empty CF API clientID")
	}
	if clientSecret == "" {
		return nil, errors.New("expected an non-empty CF API clientSecret")
	}

	c, err := newDefault(apiRoot)
	if err != nil {
		return nil, err
	}
	c.ClientID = clientID
	c.ClientSecret = clientSecret

	return c, nil
}

// NewToken creates a new config configured to use a static token
//
// This method of authentication does _not_ support re-authentication, the access token
// and/or refresh token must be valid and created externally to this client.
//
// If accessToken is empty, refreshToken must be non-empty and valid - an access token will be generated
// automatically using the refresh token.
func NewToken(apiRoot, accessToken, refreshToken string) (*Config, error) {
	if accessToken == "" && refreshToken == "" {
		return nil, errors.New("expected an non-empty CF API access token or refresh token")
	}

	c, err := newDefault(apiRoot)
	if err != nil {
		return nil, err
	}
	c.AccessToken = accessToken
	c.RefreshToken = refreshToken

	return c, nil
}

// NewFromCFHome is similar to NewToken but reads the access token from the CF_HOME config, which must
// exist and have a valid access token.
//
// This will use the currently configured CF_HOME env var if it exists, otherwise attempts to use the
// default CF_HOME directory.
func NewFromCFHome() (*Config, error) {
	dir, err := findCFHomeDir()
	if err != nil {
		return nil, err
	}
	return NewFromCFHomeDir(dir)
}

// NewFromCFHomeDir is similar to NewToken but reads the access token from the config in the specified directory
// which must exist and have a valid access token.
func NewFromCFHomeDir(cfHomeDir string) (*Config, error) {
	cfHomeConfig, err := loadCFHomeConfig(cfHomeDir)
	if err != nil {
		return nil, err
	}

	cfg, err := newDefault(cfHomeConfig.Target)
	if err != nil {
		return nil, err
	}
	cfg.AccessToken = cfHomeConfig.AccessToken
	cfg.RefreshToken = cfHomeConfig.RefreshToken
	cfg.skipTLSValidation = cfHomeConfig.SSLDisabled

	return cfg, nil
}

// WithHTTPClient overrides the default http.Client to be used as the base for all requests
//
// # The TLS and Timeout values on the http.Client will be set to match the config
//
// This is useful if you need to configure advanced http.Client or http.Transport settings,
// most consumers will not need to use this.
func (c *Config) WithHTTPClient(httpClient *http.Client) {
	c.baseHTTPClient = httpClient
	c.baseHTTPClient.Timeout = c.requestTimeout
	c.setTLSConfigOnHTTPClient()
}

// WithSkipTLSValidation sets the http.Client underlying transport InsecureSkipVerify
func (c *Config) WithSkipTLSValidation(skip bool) {
	c.skipTLSValidation = skip
	c.setTLSConfigOnHTTPClient()
}

// WithRequestTimeout overrides the http.Client underlying transport request timeout
func (c *Config) WithRequestTimeout(timeout time.Duration) {
	c.requestTimeout = timeout
	c.baseHTTPClient.Timeout = timeout
}

// HTTPClient returns the currently configured default base http.Client to be used as the base for all requests
func (c *Config) HTTPClient() *http.Client {
	return c.baseHTTPClient
}

// RequestTimeout returns the currently configured http.Client underlying transport request timeout
func (c *Config) RequestTimeout() time.Duration {
	return c.requestTimeout
}

// SkipTLSValidation returns the currently configured http.Client underlying transport InsecureSkipVerify
func (c *Config) SkipTLSValidation() bool {
	return c.skipTLSValidation
}

func (c *Config) setNewDefaultHTTPClient() {
	// use a copy of the default transport and it's settings
	transport := http.DefaultTransport.(*http.Transport).Clone()
	c.baseHTTPClient = &http.Client{
		Timeout:   c.requestTimeout,
		Transport: transport,
	}
}

func (c *Config) setTLSConfigOnHTTPClient() {
	// it's possible the consumer provided an oauth2.Transport instead of a http.Transport
	var tp *http.Transport
	switch t := c.baseHTTPClient.Transport.(type) {
	case *http.Transport:
		tp = t
	case *oauth2.Transport:
		if bt, ok := t.Base.(*http.Transport); ok {
			tp = bt
		}
	}

	// if we found a supported transport, set InsecureSkipVerify
	if tp != nil {
		if tp.TLSClientConfig == nil {
			tp.TLSClientConfig = &tls.Config{}
		}
		tp.TLSClientConfig.InsecureSkipVerify = c.skipTLSValidation
	}
}

func newDefault(apiRootURL string) (*Config, error) {
	u, err := url.ParseRequestURI(apiRootURL)
	if err != nil {
		return nil, fmt.Errorf("expected an http(s) CF API root URI, but got %s: %w", apiRootURL, err)
	}

	c := &Config{
		APIEndpointURL:    strings.TrimRight(u.String(), "/"),
		UserAgent:         UserAgent,
		skipTLSValidation: false,
		requestTimeout:    30 * time.Second,
	}
	c.setNewDefaultHTTPClient()
	c.setTLSConfigOnHTTPClient()
	return c, nil
}

func loadCFHomeConfig(cfHomeDir string) (*cfHomeConfig, error) {
	cfConfigDir := filepath.Join(cfHomeDir, ".cf")
	cfJSON, err := os.ReadFile(filepath.Join(cfConfigDir, "config.json"))
	if err != nil {
		return nil, err
	}

	var cfg cfHomeConfig
	err = json.Unmarshal(cfJSON, &cfg)
	if err == nil {
		if len(cfg.AccessToken) > len("bearer ") {
			cfg.AccessToken = cfg.AccessToken[len("bearer "):]
		}
	}

	return &cfg, nil
}

func findCFHomeDir() (string, error) {
	cfHomeDir := os.Getenv("CF_HOME")
	if cfHomeDir != "" {
		return cfHomeDir, nil
	}
	return os.UserHomeDir()
}
