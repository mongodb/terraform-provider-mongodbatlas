package config_test

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockOAuthTokenServer(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var tokenPosts atomic.Int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case clientcredentials.TokenAPIPath:
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			tokenPosts.Add(1)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"access_token":"test-access-token","token_type":"Bearer","expires_in":3600}`))
		case clientcredentials.RevokeAPIPath:
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	})
	return httptest.NewServer(handler), &tokenPosts
}

func TestGetTokenSource_cacheHit(t *testing.T) {
	config.ResetSATokenCacheForTest()
	server, tokenPosts := mockOAuthTokenServer(t)
	defer server.Close()

	const (
		clientID         = "client-a"
		clientSecret     = "secret-a"
		terraformVersion = "1.9.0"
	)
	baseURL := server.URL

	first, err := config.GetTokenSourceForTest(clientID, clientSecret, baseURL, terraformVersion)
	require.NoError(t, err)
	require.Equal(t, int32(1), tokenPosts.Load())

	second, err := config.GetTokenSourceForTest(clientID, clientSecret, baseURL, terraformVersion)
	require.NoError(t, err)
	assert.Same(t, first, second)
	assert.Equal(t, int32(1), tokenPosts.Load(), "cache hit should not request another token")
}

func TestGetTokenSource_differentSecretNewSource(t *testing.T) {
	config.ResetSATokenCacheForTest()
	server, tokenPosts := mockOAuthTokenServer(t)
	defer server.Close()

	const (
		clientID         = "client-rotate"
		terraformVersion = "1.9.0"
	)
	baseURL := server.URL

	first, err := config.GetTokenSourceForTest(clientID, "secret-old", baseURL, terraformVersion)
	require.NoError(t, err)

	second, err := config.GetTokenSourceForTest(clientID, "secret-new", baseURL, terraformVersion)
	require.NoError(t, err)
	assert.NotSame(t, first, second)
	assert.Equal(t, int32(2), tokenPosts.Load())
}

func TestCloseTokenSource_clearsCache(t *testing.T) {
	config.ResetSATokenCacheForTest()
	server, _ := mockOAuthTokenServer(t)
	defer server.Close()

	const terraformVersion = "1.9.0"
	baseURL := server.URL

	_, err := config.GetTokenSourceForTest("client-1", "secret-1", baseURL, terraformVersion)
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-2", "secret-2", baseURL, terraformVersion)
	require.NoError(t, err)
	require.Equal(t, 2, config.SATokenCacheLenForTest())

	config.CloseTokenSource()
	require.Equal(t, 0, config.SATokenCacheLenForTest())

	_, err = config.GetTokenSourceForTest("client-1", "secret-1", baseURL, terraformVersion)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already closed")
}

func TestGetTokenSource_normalizesBaseURL(t *testing.T) {
	config.ResetSATokenCacheForTest()
	server, tokenPosts := mockOAuthTokenServer(t)
	defer server.Close()

	const (
		clientID         = "client-trim"
		clientSecret     = "secret-trim"
		terraformVersion = "1.9.0"
	)
	baseURL := server.URL + "/"

	first, err := config.GetTokenSourceForTest(clientID, clientSecret, baseURL, terraformVersion)
	require.NoError(t, err)

	second, err := config.GetTokenSourceForTest(clientID, clientSecret, server.URL, terraformVersion)
	require.NoError(t, err)
	assert.Same(t, first, second)
	assert.Equal(t, int32(1), tokenPosts.Load())
}

func TestCloseTokenSource_idempotent(t *testing.T) {
	config.ResetSATokenCacheForTest()
	config.CloseTokenSource()
	require.NotPanics(t, config.CloseTokenSource)
}
