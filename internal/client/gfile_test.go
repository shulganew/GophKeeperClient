package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shulganew/GophKeeperClient/internal/app/logging"
	"github.com/shulganew/GophKeeperClient/internal/client/oapi"
	"github.com/stretchr/testify/require"
)

func TestGfileAdd(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		method string
	}{
		{
			name:   "Check file list.",
			path:   "/user/file",
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
			// File not existed.
			_, _, err = FileAdd(c, "jwt", "Test file methods", "readme.md")
			require.Error(t, err)

		})
	}
}

func TestGfileList(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		status int
		method string
	}{
		{
			name:   "Check file list.",
			path:   "/user/file",
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
			_, status, err := GfileList(c, "")
			require.NoError(t, err)
			require.Equal(t, status, tt.status)
		})
	}
}
