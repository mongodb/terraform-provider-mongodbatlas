package config_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDigestClient_ResendsAfterConnectionClose covers a server that closes a keep-alive
// connection after the digest challenge, without answering the authenticated request sent on it.
// The client must resend the request on a new connection instead of failing with EOF.
func TestDigestClient_ResendsAfterConnectionClose(t *testing.T) {
	const projectID = "111111111111111111111111"
	var dropped atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.Header().Set("WWW-Authenticate", `Digest realm="MMS Public API", nonce="abc", qop="auth", algorithm=MD5`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if dropped.CompareAndSwap(false, true) {
			conn, _, err := w.(http.Hijacker).Hijack()
			if assert.NoError(t, err) {
				assert.NoError(t, conn.Close())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"id":"` + projectID + `"}`))
		assert.NoError(t, err)
	}))
	defer srv.Close()

	client, err := config.NewClient(&config.Credentials{PublicKey: "public", PrivateKey: "private", BaseURL: srv.URL}, "test")
	require.NoError(t, err)

	project, _, err := client.AtlasV2.ProjectsAPI.GetGroup(t.Context(), projectID).Execute()
	require.NoError(t, err)
	assert.True(t, dropped.Load(), "server didn't close the connection")
	assert.Equal(t, projectID, project.GetId())
}
