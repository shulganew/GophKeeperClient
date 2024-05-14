package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shulganew/GophKeeperClient/internal/app/logging"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/stretchr/testify/require"
)

func TestSiteAdd(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		method string
	}{
		{
			name:   "Check site list.",
			path:   "/user/site",
			status: http.StatusOK,
			method: http.MethodPost,
		},
	}
	logging.InitLog()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				t.Log("Path: ", r.URL.Path)
				t.Log("Method: ", r.Method)

				// Test path.
				require.Equal(t, r.URL.String(), tt.path)

				// Test method.

				require.Equal(t, r.Method, tt.method)

				// Set status.
				w.WriteHeader(tt.status)

			}))

			// Close the server when test finishes
			defer server.Close()

			// Create client.
			c, err := oapi.NewClient(server.URL, oapi.WithHTTPClient(server.Client()))
			require.NoError(t, err)

			// Use Client & URL from our local test server.
			_, status, err := SiteAdd(c, "jwt", "Test site methods", "www.ru", "user", "123")
			require.NoError(t, err)
			require.Equal(t, status, tt.status)
		})
	}
}

func TestSiteList(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		method string
	}{
		{
			name:   "Check site list.",
			path:   "/user/site",
			status: http.StatusOK,
			method: http.MethodGet,
		},
	}
	logging.InitLog()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				t.Log("Path: ", r.URL.Path)
				t.Log("Method: ", r.Method)

				// Test path.
				require.Equal(t, r.URL.String(), tt.path)

				// Test method.

				require.Equal(t, r.Method, tt.method)

				// Set status.
				w.WriteHeader(tt.status)

			}))

			// Close the server when test finishes
			defer server.Close()

			// Create client.
			c, err := oapi.NewClient(server.URL, oapi.WithHTTPClient(server.Client()))
			require.NoError(t, err)

			// Use Client & URL from our local test server.
			_, status, err := SiteList(c, "")
			require.NoError(t, err)
			require.Equal(t, status, tt.status)
		})
	}
}

func TestSiteUpdate(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		method string
	}{
		{
			name:   "Check site update.",
			path:   "/user/site",
			status: http.StatusOK,
			method: http.MethodPut,
		},
	}
	logging.InitLog()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Start a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

				t.Log("Path: ", r.URL.Path)
				t.Log("Method: ", r.Method)

				// Test path.
				require.Equal(t, r.URL.String(), tt.path)

				// Test method.
				require.Equal(t, r.Method, tt.method)

				// Set status.
				w.WriteHeader(tt.status)

			}))

			// Close the server when test finishes
			defer server.Close()

			// Create client.
			c, err := oapi.NewClient(server.URL, oapi.WithHTTPClient(server.Client()))
			require.NoError(t, err)

			// Use Client & URL from our local test server. string, siteID, def, siteURL, slogin, spw string
			status, err := SiteUpdate(c, "jwt", "siteID", "Test site methods", "www.ru", "user", "123")
			require.NoError(t, err)
			require.Equal(t, status, tt.status)
		})
	}
}
