package client

import (
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"github.com/cloudfoundry-community/go-cfclient/v3/testutil"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type RouteTest struct {
	Description string
	Route       testutil.MockRoute
	Expected    string
	Expected2   string
	Expected3   string
	Action      func(c *Client, t *testing.T) (any, error)
	Action2     func(c *Client, t *testing.T) (any, any, error)
	Action3     func(c *Client, t *testing.T) (any, any, any, error)
}

func ExecuteTests(tests []RouteTest, t *testing.T) {
	for _, tt := range tests {
		func() {
			serverURL := testutil.Setup(tt.Route, t)
			defer testutil.Teardown()
			details := fmt.Sprintf("%s %s", tt.Route.Method, tt.Route.Endpoint)
			if tt.Description != "" {
				details = tt.Description + ": " + details
			}

			c, _ := config.NewToken(serverURL, "", "fake-refresh-token")
			cl, err := New(c)
			require.NoError(t, err, details)

			assertEq := func(t *testing.T, expected string, obj any) {
				if isJSON(expected) {
					actualJSON, err := json.Marshal(obj)
					require.NoError(t, err, details)
					require.JSONEq(t, expected, string(actualJSON), details)
				} else {
					if s, ok := obj.(string); ok {
						require.Equal(t, expected, s, details)
					}
				}
			}

			if tt.Action != nil {
				obj1, err := tt.Action(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
			} else if tt.Action2 != nil {
				obj1, obj2, err := tt.Action2(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
				assertEq(t, tt.Expected2, obj2)
			} else if tt.Action3 != nil {
				obj1, obj2, obj3, err := tt.Action3(cl, t)
				require.NoError(t, err, details)
				assertEq(t, tt.Expected, obj1)
				assertEq(t, tt.Expected2, obj2)
				assertEq(t, tt.Expected3, obj3)
			}
		}()
	}
}

func isJSON(obj string) bool {
	return strings.HasPrefix(obj, "{") || strings.HasPrefix(obj, "[")
}
