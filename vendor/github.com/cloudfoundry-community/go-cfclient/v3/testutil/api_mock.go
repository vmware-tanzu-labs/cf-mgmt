package testutil

import (
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

var (
	mux           *http.ServeMux
	server        *httptest.Server
	fakeUAAServer *httptest.Server
)

type MockRoute struct {
	Method           string
	Endpoint         string
	Output           []string
	UserAgent        string
	Status           int
	QueryString      string
	PostForm         string
	RedirectLocation string
}

func SetupFakeAPIServer() string {
	if fakeUAAServer == nil {
		SetupFakeUAAServer(3)
	}
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	return server.URL
}

func SetupFakeUAAServer(expiresIn int) string {
	uaaMux := http.NewServeMux()
	fakeUAAServer = httptest.NewServer(uaaMux)
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	count := 1
	r.Post("/oauth/token", func(r render.Render) {
		r.JSON(200, map[string]interface{}{
			"token_type":    "bearer",
			"access_token":  "foobar" + strconv.Itoa(count),
			"refresh_token": "barfoo",
			"expires_in":    expiresIn,
		})
		count = count + 1
	})
	r.NotFound(func() string { return "" })
	m.Action(r.Handle)
	uaaMux.Handle("/", m)
	return fakeUAAServer.URL
}

func Setup(mock MockRoute, t *testing.T) string {
	return SetupMultiple([]MockRoute{mock}, t)
}

func SetupMultiple(mockEndpoints []MockRoute, t *testing.T) string {
	if server == nil {
		SetupFakeAPIServer()
	}
	m := martini.New()
	m.Use(render.Renderer())
	r := martini.NewRouter()
	for _, mock := range mockEndpoints {
		method := mock.Method
		endpoint := mock.Endpoint
		output := mock.Output
		if len(output) == 0 {
			output = []string{""}
		}
		userAgent := mock.UserAgent
		status := mock.Status
		queryString := mock.QueryString
		postFormBody := mock.PostForm
		redirectLocation := mock.RedirectLocation
		switch method {
		case "GET":
			count := 0
			r.Get(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				singleOutput := output[count]
				count++
				return status, singleOutput
			})
		case "POST":
			r.Post(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "DELETE":
			r.Delete(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PUT":
			r.Put(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PATCH":
			r.Patch(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testQueryString(req.URL.RawQuery, queryString, t)
				testReqBody(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		case "PUT-FILE":
			r.Put(endpoint, func(res http.ResponseWriter, req *http.Request) (int, string) {
				testUserAgent(req.Header.Get("User-Agent"), userAgent, t)
				testBodyContains(req, postFormBody, t)
				if redirectLocation != "" {
					res.Header().Add("Location", redirectLocation)
				}
				return status, output[0]
			})
		}
	}
	r.Get("/", func(r render.Render) {
		r.JSON(200, map[string]any{
			"links": map[string]any{
				"cloud_controller_v2": map[string]any{
					"href": server.URL + "/v2",
					"meta": map[string]any{
						"version": "2.155.0",
					},
				},
				"cloud_controller_v3": map[string]any{
					"href": server.URL + "/v3",
					"meta": map[string]any{
						"version": "3.90.0",
					},
				},
				"network_policy_v0": map[string]any{
					"href": "https://api.example.org/networking/v0/external",
				},
				"network_policy_v1": map[string]any{
					"href": "https://api.example.org/networking/v1/external",
				},
				"uaa": map[string]any{
					"href": fakeUAAServer.URL,
				},
				"login": map[string]any{
					"href": fakeUAAServer.URL,
				},
				"credhub": map[string]any{
					"href": "",
				},
				"routing": map[string]any{
					"href": "https://api.example.org/routing",
				},
				"logging": map[string]any{
					"href": "wss://doppler.example.org:443",
				},
				"log_cache": map[string]any{
					"href": "https://log-cache.example.org",
				},
				"log_stream": map[string]any{
					"href": "https://log-stream.example.org",
				},
				"app_ssh": map[string]any{
					"href": "ssh.example.org:2222",
					"meta": map[string]any{
						"host_key_fingerprint": "Y411oivJwZCUQnXHq83mdM5SKCK4ftyoSXI31RRe4Zs",
						"oauth_client":         "ssh-proxy",
					},
				},
			},
		})
	})

	m.Action(r.Handle)
	mux.Handle("/", m)

	return server.URL
}

func Teardown() {
	if server != nil {
		server.Close()
		server = nil
	}
	if fakeUAAServer != nil {
		fakeUAAServer.Close()
		fakeUAAServer = nil
	}
}

func testQueryString(QueryString string, QueryStringExp string, t *testing.T) {
	t.Helper()
	if QueryStringExp == "" {
		return
	}

	value, _ := url.QueryUnescape(QueryString)
	if QueryStringExp != value {
		t.Errorf("Error: Query string '%s' should be equal to '%s'", QueryStringExp, value)
	}
}

func testUserAgent(UserAgent string, UserAgentExp string, t *testing.T) {
	t.Helper()
	if len(UserAgentExp) < 1 {
		UserAgentExp = config.UserAgent
	}
	if UserAgent != UserAgentExp {
		t.Errorf("Error: Agent %s should be equal to %s", UserAgent, UserAgentExp)
	}
}

func testReqBody(req *http.Request, postFormBody string, t *testing.T) {
	t.Helper()
	if postFormBody != "" {
		if body, err := io.ReadAll(req.Body); err != nil {
			t.Error("No request body but expected one")
		} else {
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)
			require.JSONEq(t, postFormBody, string(body),
				"Expected request body (%s) does not equal request body (%s)", postFormBody, body)
		}
	}
}

func testBodyContains(req *http.Request, expected string, t *testing.T) {
	t.Helper()
	if expected != "" {
		if body, err := io.ReadAll(req.Body); err != nil {
			t.Error("No request body but expected one")
		} else {
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(req.Body)
			if !strings.Contains(string(body), expected) {
				t.Errorf("Expected request body (%s) was not found in actual request body (%s)", expected, body)
			}
		}
	}
}
