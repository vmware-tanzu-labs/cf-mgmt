package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/jwt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// OAuthSessionManager creates and manages OAuth http client instances
type OAuthSessionManager struct {
	config *config.Config

	oauthClient               *http.Client
	oauthClientNonRedirecting *http.Client

	tokenSource oauth2.TokenSource
	mutex       *sync.RWMutex
}

// NewOAuthSessionManager creates a new OAuth session manager
func NewOAuthSessionManager(config *config.Config) *OAuthSessionManager {
	return &OAuthSessionManager{
		config: config,
		mutex:  &sync.RWMutex{},
	}
}

// Client returns an authenticated OAuth http client
func (m *OAuthSessionManager) Client(ctx context.Context, followRedirects bool) (*http.Client, error) {
	err := m.init(ctx)
	if err != nil {
		return nil, err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if followRedirects {
		return m.oauthClient, nil
	}
	return m.oauthClientNonRedirecting, nil
}

// ReAuthenticate causes a new http.Client to be created with new a new authentication context,
// likely in response to a 401
//
// This won't work for userTokenAuth since we have no credentials to exchange for a new token.
func (m *OAuthSessionManager) ReAuthenticate(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.isUserTokenAuth() {
		return errors.New("cannot reauthenticate user token auth type, check your access and/or refresh token expiration date")
	}

	// attempt to create a new token source
	return m.newTokenSource(ctx)
}

// AccessToken returns the raw OAuth access token
func (m *OAuthSessionManager) AccessToken(ctx context.Context) (string, error) {
	err := m.init(ctx)
	if err != nil {
		return "", err
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	token, err := m.tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("error getting bearer token: %w", err)
	}
	return token.AccessToken, nil
}

func (m *OAuthSessionManager) init(ctx context.Context) error {
	// get a reader lock and check to see if the token source has been initialized already
	m.mutex.RLock()
	if m.tokenSource != nil {
		m.mutex.RUnlock()
		return nil
	}

	// not initialized, upgrade to write lock
	m.mutex.RUnlock()
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// attempt to create a new token source
	return m.newTokenSource(ctx)
}

// newTokenSource creates an appropriate OAuth token source based off the provided config
func (m *OAuthSessionManager) newTokenSource(ctx context.Context) error {
	if m.config.LoginEndpointURL == "" || m.config.UAAEndpointURL == "" {
		return errors.New("login and UAA endpoints must not be empty")
	}

	loginEndpoint := path.Join(m.config.LoginEndpointURL, "/oauth/auth")
	uaaEndpoint := path.Join(m.config.UAAEndpointURL, "/oauth/token")

	// this provides the http.Client instance that the oauth subsystem will use for token acquisition
	// and copy the base http transport from for new clients
	oauthCtx := context.WithValue(ctx, oauth2.HTTPClient, m.config.HTTPClient())

	switch {
	case m.isUserTokenAuth():
		return m.userTokenAuth(oauthCtx, loginEndpoint, uaaEndpoint)
	case m.isClientAuth():
		m.clientAuth(oauthCtx, uaaEndpoint)
	default:
		return m.userAuth(oauthCtx, loginEndpoint, uaaEndpoint)
	}

	return nil
}

// userAuth initializes a http client using standard username and password
func (m *OAuthSessionManager) userAuth(ctx context.Context, loginEndpoint, uaaEndpoint string) error {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  loginEndpoint,
			TokenURL: uaaEndpoint,
		},
	}
	if m.config.Origin != "" {
		type LoginHint struct {
			Origin string `json:"origin"`
		}
		loginHint := LoginHint{m.config.Origin}
		origin, err := json.Marshal(loginHint)
		if err != nil {
			return fmt.Errorf("error creating login_hint for user auth: %w", err)
		}
		val := url.Values{}
		val.Set("login_hint", string(origin))
		authConfig.Endpoint.TokenURL = path.Format("%s?%s", authConfig.Endpoint.TokenURL, val)
	}

	token, err := authConfig.PasswordCredentialsToken(ctx, m.config.Username, m.config.Password)
	if err != nil {
		return fmt.Errorf("error getting token for user auth: %w", err)
	}

	tokenSource := authConfig.TokenSource(ctx, token)
	m.initOAuthClient(ctx, tokenSource)

	return nil
}

// clientAuth initializes a http client using OAuth client id and secret
func (m *OAuthSessionManager) clientAuth(ctx context.Context, uaaEndpoint string) {
	authConfig := &clientcredentials.Config{
		ClientID:     m.config.ClientID,
		ClientSecret: m.config.ClientSecret,
		TokenURL:     uaaEndpoint,
	}
	tokenSource := authConfig.TokenSource(ctx)
	m.initOAuthClient(ctx, tokenSource)
}

// userTokenAuth initializes client credentials from existing bearer token.
func (m *OAuthSessionManager) userTokenAuth(ctx context.Context, loginEndpoint, uaaEndpoint string) (err error) {
	authConfig := &oauth2.Config{
		ClientID: "cf",
		Scopes:   []string{""},
		Endpoint: oauth2.Endpoint{
			AuthURL:  loginEndpoint,
			TokenURL: uaaEndpoint,
		},
	}

	// we could be given only a refresh token, so optionally parse
	var exp time.Time
	if m.config.AccessToken != "" {
		exp, err = jwt.AccessTokenExpiration(m.config.AccessToken)
		if err != nil {
			return err
		}
	}

	// AccessToken is expected to have no "bearer" prefix
	token := &oauth2.Token{
		RefreshToken: m.config.RefreshToken,
		AccessToken:  m.config.AccessToken,
		Expiry:       exp,
		TokenType:    "Bearer",
	}
	tokenSource := authConfig.TokenSource(ctx, token)
	m.initOAuthClient(ctx, tokenSource)

	return nil
}

func (m *OAuthSessionManager) initOAuthClient(ctx context.Context, tokenSource oauth2.TokenSource) {
	bc := m.config.HTTPClient()

	// oauth2.NewClient copies the underlying transport only, so explicitly copy other client values over
	// without modifying the oauth2 client that was returned (since that's unsupported)
	// https://github.com/golang/oauth2/issues/368
	oac := oauth2.NewClient(ctx, tokenSource)
	oauthClient := &http.Client{
		Transport:     oac.Transport,
		Timeout:       bc.Timeout,
		Jar:           bc.Jar,
		CheckRedirect: bc.CheckRedirect,
	}

	oauthClientNonRedirecting := &http.Client{
		Transport: oac.Transport,
		Timeout:   bc.Timeout,
		Jar:       bc.Jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// cache the created client and token source
	m.oauthClientNonRedirecting = oauthClientNonRedirecting
	m.oauthClient = oauthClient
	m.tokenSource = tokenSource
}

func (m *OAuthSessionManager) isUserTokenAuth() bool {
	return m.config.RefreshToken != "" || m.config.AccessToken != ""
}

func (m *OAuthSessionManager) isClientAuth() bool {
	return m.config.ClientID != ""
}
