package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/internal/path"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

var errNilContext = errors.New("context cannot be nil")

// supportedHTTPMethods are the HTTP verbs this executor supports
var supportedHTTPMethods = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodPatch,
}

type unauthorizedError struct {
	Err error
}

func (e unauthorizedError) Error() string {
	return fmt.Sprintf("unable to get new access token: %s", e.Err)
}

// Executor handles executing HTTP requests
type Executor struct {
	userAgent      string
	apiAddress     string
	clientProvider ClientProvider
}

// NewExecutor creates a new HTTP Executor instance
func NewExecutor(clientProvider ClientProvider, apiAddress, userAgent string) *Executor {
	return &Executor{
		userAgent:      userAgent,
		apiAddress:     apiAddress,
		clientProvider: clientProvider,
	}
}

// ExecuteRequest executes the specified request using the http.Client provided by the client provider
func (c *Executor) ExecuteRequest(request *Request) (*http.Response, error) {
	followRedirects := request.followRedirects
	req, err := c.newHTTPRequest(request)
	if err != nil {
		return nil, err
	}

	// do the request to the remote API
	r, err := c.do(req, followRedirects)

	// it's possible the access token expired and the oauth subsystem could not obtain a new one because the
	// refresh token is expired or revoked. Attempt to get a new refresh and access token and retry the request.
	var authErr *unauthorizedError
	if errors.As(err, &authErr) {
		err = c.reAuthenticate(req.Context())
		if err != nil {
			return nil, err
		}
		r, err = c.do(req, followRedirects)
	}

	return r, err
}

// newHTTPRequest creates a new *http.Request instance from the internal model
func (c *Executor) newHTTPRequest(request *Request) (*http.Request, error) {
	if request.context == nil {
		return nil, errNilContext
	}
	if !isSupportedHTTPMethod(request.method) {
		return nil, fmt.Errorf("error executing request, found unsupported HTTP method %s", request.method)
	}

	// JSON encode the object and use that as the body if specified, otherwise use the body as-is
	reqBody := request.body
	if request.object != nil {
		b, err := encodeBody(request.object)
		if err != nil {
			return nil, fmt.Errorf("error executing request, failed to encode the request object to JSON: %w", err)
		}
		reqBody = b
	}
	u := path.Join(c.apiAddress, request.pathAndQuery)

	r, err := http.NewRequestWithContext(request.context, request.method, u, reqBody)
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed to create a new underlying HTTP request: %w", err)
	}
	r.Header.Set("User-Agent", c.userAgent)
	if request.contentType != "" {
		r.Header.Set("Content-type", request.contentType)
	}
	if request.contentLength != nil {
		r.ContentLength = *request.contentLength
	}
	for k, v := range request.headers {
		r.Header.Set(k, v)
	}

	return r, nil
}

// do will get the proper http.Client and calls Do on it using the specified http.Request
func (c *Executor) do(request *http.Request, followRedirects bool) (*http.Response, error) {
	client, err := c.clientProvider.Client(request.Context(), followRedirects)
	if err != nil {
		return nil, fmt.Errorf("error executing request, failed to get the underlying HTTP client: %w", err)
	}
	r, err := client.Do(request)
	if err != nil {
		// if we get an error because the context was cancelled, the context's error is more useful.
		ctx := request.Context()
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// see if the oauth subsystem was unable to use the refresh token to get a new access token
		var oauthErr *oauth2.RetrieveError
		if errors.As(err, &oauthErr) {
			if oauthErr.Response.StatusCode == http.StatusUnauthorized {
				return nil, &unauthorizedError{
					Err: err,
				}
			}
		}

		return nil, fmt.Errorf("error executing request, failed during HTTP request send: %w", err)
	}

	// perhaps the token looked valid, but was revoked etc
	if r.StatusCode == http.StatusUnauthorized {
		_ = r.Body.Close()
		return nil, &unauthorizedError{}
	}

	return r, nil
}

// reAuthenticate tells the client provider to restart authentication anew because we received a 401
func (c *Executor) reAuthenticate(ctx context.Context) error {
	err := c.clientProvider.ReAuthenticate(ctx)
	if err != nil {
		return fmt.Errorf("an error occurred attempting to reauthenticate "+
			"after initially receiving a 401 executing a request: %w", err)
	}
	return nil
}

// encodeBody is used to encode a request body
func encodeBody(obj any) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

// isSupportedHTTPMethod returns true if the executor supports this HTTP method
func isSupportedHTTPMethod(method string) bool {
	for _, v := range supportedHTTPMethods {
		if v == method {
			return true
		}
	}
	return false
}
