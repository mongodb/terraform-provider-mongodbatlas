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

type saCacheKey struct {
	clientID         string
	clientSecret     string
	baseURL          string
	terraformVersion string
}

type saCacheEntry struct {
	tokenSource auth.TokenSource
}

// saTokenCache holds OAuth token sources keyed by service account credentials.
//
// The provider process may call NewClient with different (clientID, clientSecret) pairs
// during one Terraform run—for example organization3 secret rotation (same client ID,
// new secret) or multiple SA-backed resources. A single global token source cannot
// represent more than one credential set; the previous implementation rejected credential
// changes with "service account credentials changed" and broke Read after rotation.
//
// Each key gets its own ReuseTokenSourceWithExpiry so tokens are reused per credential
// tuple without cross-talk. CloseTokenSource revokes every cached token on plugin exit.
type saTokenCacheState struct {
	entries map[saCacheKey]saCacheEntry
	mu      sync.Mutex
	closed  bool
}

var saTokenCache = saTokenCacheState{
	entries: make(map[saCacheKey]saCacheEntry),
}

// getTokenSource returns a cached TokenSource for the given SA credentials, creating
// and caching one on first use. baseURL is normalized before lookup.
func getTokenSource(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	saTokenCache.mu.Lock()
	defer saTokenCache.mu.Unlock()

	if saTokenCache.closed {
		return nil, fmt.Errorf("service account token source already closed")
	}

	key := saCacheKey{
		clientID:         clientID,
		clientSecret:     clientSecret,
		baseURL:          NormalizeBaseURL(baseURL),
		terraformVersion: terraformVersion,
	}
	if entry, ok := saTokenCache.entries[key]; ok {
		return entry.tokenSource, nil
	}

	// Use a new context to avoid "context canceled" errors as the token source is reused and can outlast the callee context.
	ctx := context.WithValue(context.Background(), auth.HTTPClient, NewOAuthHTTPClient(terraformVersion))
	conf := GetServiceAccountConfig(clientID, clientSecret, key.baseURL)
	tokenSource := oauth2.ReuseTokenSourceWithExpiry(nil, conf.TokenSource(ctx), saTokenExpiryBuffer)
	if _, err := tokenSource.Token(); err != nil { // Retrieve token to fail-fast if credentials are invalid.
		return nil, err
	}
	saTokenCache.entries[key] = saCacheEntry{tokenSource: tokenSource}
	return tokenSource, nil
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

// CloseTokenSource is called just before the provider finishes. It sets saTokenCache.closed
// to avoid future calls to getTokenSource, which should not happen as the provider is exiting.
// It best-effort revokes every cached OAuth token, not only the last credential set.
func CloseTokenSource() {
	saTokenCache.mu.Lock()
	defer saTokenCache.mu.Unlock()
	if saTokenCache.closed {
		return
	}
	saTokenCache.closed = true
	for key, entry := range saTokenCache.entries {
		if token, err := entry.tokenSource.Token(); err == nil {
			conf := GetServiceAccountConfig(key.clientID, key.clientSecret, key.baseURL)
			ctx := context.WithValue(context.Background(), auth.HTTPClient, NewOAuthHTTPClient(key.terraformVersion))
			_ = conf.RevokeToken(ctx, token) // Best-effort, no need to do anything if it fails.
		}
	}
	clear(saTokenCache.entries)
}
