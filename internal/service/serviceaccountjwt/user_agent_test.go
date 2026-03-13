package serviceaccountjwt_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/serviceaccountjwt"
)

func TestUserAgent_SentOnTokenRequest(t *testing.T) {
	var captured capturedHeaders
	srv := newFakeOAuthServer(t, &captured)
	defer srv.Close()

	r := &serviceaccountjwt.ES{
		ESCommon: config.ESCommon{
			EphemeralResourceData: &config.EphemeralResourceData{TerraformVersion: "1.11.0"},
			ResourceName:          serviceaccountjwt.ResourceTypeName,
		},
	}
	ctx := config.AddUserAgentExtra(context.Background(), config.UserAgentExtra{
		Name:      "service_account_jwt",
		Operation: config.UserAgentOperationValueOpen,
	})

	token, err := r.GenerateToken(ctx, "test-id", "test-secret", srv.URL)
	require.NoError(t, err)
	require.NotEmpty(t, token.AccessToken)

	ua := captured.get("token")
	require.NotEmpty(t, ua, "token request must carry a User-Agent header")
	assert.Contains(t, ua, "terraform-provider-mongodbatlas/")
	assert.Contains(t, ua, "Terraform/1.11.0")
	assert.Contains(t, ua, "Name/service_account_jwt")
	assert.Contains(t, ua, "Operation/open")
}

func TestUserAgent_SentOnRevokeRequest(t *testing.T) {
	var captured capturedHeaders
	srv := newFakeOAuthServer(t, &captured)
	defer srv.Close()

	r := &serviceaccountjwt.ES{
		ESCommon: config.ESCommon{
			EphemeralResourceData: &config.EphemeralResourceData{TerraformVersion: "1.11.0"},
			ResourceName:          serviceaccountjwt.ResourceTypeName,
		},
	}
	ctx := config.AddUserAgentExtra(context.Background(), config.UserAgentExtra{
		Name:      "service_account_jwt",
		Operation: config.UserAgentOperationValueClose,
	})

	conf := config.GetServiceAccountConfig("test-id", "test-secret", srv.URL)
	err := conf.RevokeToken(r.WithUserAgentClient(ctx), &oauth2.Token{AccessToken: "tok-abc"})
	require.NoError(t, err)

	ua := captured.get("revoke")
	require.NotEmpty(t, ua, "revoke request must carry a User-Agent header")
	assert.Contains(t, ua, "terraform-provider-mongodbatlas/")
	assert.Contains(t, ua, "Terraform/1.11.0")
	assert.Contains(t, ua, "Name/service_account_jwt")
	assert.Contains(t, ua, "Operation/close")
}

// --- helpers ---

type capturedHeaders struct {
	token  string
	revoke string
	mu     sync.Mutex
}

func (c *capturedHeaders) set(kind, ua string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch kind {
	case "token":
		c.token = ua
	case "revoke":
		c.revoke = ua
	}
}

func (c *capturedHeaders) get(kind string) string {
	c.mu.Lock()
	defer c.mu.Unlock()
	switch kind {
	case "token":
		return c.token
	case "revoke":
		return c.revoke
	}
	return ""
}

func newFakeOAuthServer(t *testing.T, captured *capturedHeaders) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/oauth/token", func(w http.ResponseWriter, r *http.Request) {
		captured.set("token", r.Header.Get("User-Agent"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "fake-token-12345",
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	})
	mux.HandleFunc("/api/oauth/revoke", func(w http.ResponseWriter, r *http.Request) {
		captured.set("revoke", r.Header.Get("User-Agent"))
		w.WriteHeader(http.StatusOK)
	})
	return httptest.NewServer(mux)
}
