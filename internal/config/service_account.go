package config

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mongodb/atlas-sdk-go/auth"
	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
	"golang.org/x/oauth2"
)

// Renew token if it expires within 10 minutes to avoid authentication errors during Atlas API calls.
const saTokenExpiryBuffer = 10 * time.Minute

type saTokenSourceEntry struct {
	tokenSource      auth.TokenSource
	clientSecret     string
	baseURL          string
	terraformVersion string
}

// saTokenSourceCache caches token sources per service account (keyed by clientID) so a single
// provider process can authenticate as more than one SA (e.g. org-creator SA, then the SA
// created with a new org).
var saTokenSourceCache = struct {
	sources map[string]*saTokenSourceEntry
	mu      sync.Mutex
	closed  bool
}{sources: make(map[string]*saTokenSourceEntry)}

// createTokenSourceFn is the OAuth token-source factory; overridden in unit tests.
var createTokenSourceFn = defaultCreateTokenSource

// revokeTokenFn revokes a single Service Account token; overridden in unit tests.
var revokeTokenFn = defaultRevokeToken

func defaultCreateTokenSource(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	// Use a new context to avoid "context canceled" errors as the token source is reused and can outlast the callee context.
	ctx := context.WithValue(context.Background(), auth.HTTPClient, NewOAuthHTTPClient(terraformVersion))
	conf := GetServiceAccountConfig(clientID, clientSecret, baseURL)
	tokenSource := oauth2.ReuseTokenSourceWithExpiry(nil, conf.TokenSource(ctx), saTokenExpiryBuffer)
	if _, err := tokenSource.Token(); err != nil { // Retrieve token to fail-fast if credentials are invalid.
		return nil, err
	}
	return tokenSource, nil
}

func defaultRevokeToken(clientID string, entry *saTokenSourceEntry) {
	token, err := entry.tokenSource.Token()
	if err != nil {
		return // Best-effort, no need to do anything if the token can't be retrieved.
	}
	conf := GetServiceAccountConfig(clientID, entry.clientSecret, entry.baseURL)
	ctx := context.WithValue(context.Background(), auth.HTTPClient, NewOAuthHTTPClient(entry.terraformVersion))
	_ = conf.RevokeToken(ctx, token) // Best-effort, no need to do anything if it fails.
}

func getTokenSource(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	baseURL = NormalizeBaseURL(baseURL)

	// Read from cache.
	saTokenSourceCache.mu.Lock()
	ts, err := getCachedTokenSource(clientID, clientSecret, baseURL)
	saTokenSourceCache.mu.Unlock()
	if ts != nil || err != nil {
		return ts, err
	}

	// Cache miss, fetch token.
	tokenSource, err := createTokenSourceFn(clientID, clientSecret, baseURL, terraformVersion)
	if err != nil {
		return nil, err
	}

	// Re-acquire lock.
	saTokenSourceCache.mu.Lock()
	defer saTokenSourceCache.mu.Unlock()

	// Re-check cache.
	if ts, err = getCachedTokenSource(clientID, clientSecret, baseURL); ts != nil || err != nil {
		return ts, err
	}

	// Write to cache.
	saTokenSourceCache.sources[clientID] = &saTokenSourceEntry{
		tokenSource:      tokenSource,
		clientSecret:     clientSecret,
		baseURL:          baseURL,
		terraformVersion: terraformVersion,
	}
	return tokenSource, nil
}

// getCachedTokenSource checks the closed flag and reads from the sources map. Callers must hold a lock on saTokenSourceCache.mu.
func getCachedTokenSource(clientID, clientSecret, baseURL string) (auth.TokenSource, error) {
	if saTokenSourceCache.closed {
		return nil, fmt.Errorf("service account token source already closed")
	}
	entry, ok := saTokenSourceCache.sources[clientID]
	if !ok {
		return nil, nil
	}
	if entry.clientSecret != clientSecret || entry.baseURL != baseURL {
		return nil, fmt.Errorf("service account credentials changed")
	}
	return entry.tokenSource, nil
}

func NormalizeBaseURL(baseURL string) string {
	return strings.TrimRight(baseURL, "/")
}

func GetServiceAccountConfig(clientID, clientSecret, baseURL string) *clientcredentials.Config {
	config := clientcredentials.NewConfig(clientID, clientSecret)
	if baseURL != "" {
		config.TokenURL = baseURL + clientcredentials.TokenAPIPath
		config.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	}
	return config
}

// CloseTokenSource is called just before the provider finishes, it does a best-effort try to revoke all cached Service Account tokens.
// It sets saTokenSourceCache.closed = true to avoid future calls to getTokenSource, that shouldn't happen as the provider is exiting.
func CloseTokenSource() {
	saTokenSourceCache.mu.Lock()
	defer saTokenSourceCache.mu.Unlock()
	if saTokenSourceCache.closed {
		return
	}
	saTokenSourceCache.closed = true
	for clientID, entry := range saTokenSourceCache.sources {
		revokeTokenFn(clientID, entry)
	}
}
