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

type saCacheEntry struct {
	tokenSource      auth.TokenSource
	clientID         string
	clientSecret     string
	baseURL          string
	terraformVersion string
}

// saInfo caches token sources per service account so a single provider process can
// authenticate as more than one SA (e.g. org-creator SA, then the SA created with a new org).
var saInfo = struct {
	sources map[string]*saCacheEntry
	mu      sync.Mutex
	closed  bool
}{}

// createTokenSourceFn is the OAuth token-source factory; overridden in unit tests.
var createTokenSourceFn = defaultCreateTokenSource

func saCacheKey(clientID, baseURL string) string {
	return clientID + "\x00" + NormalizeBaseURL(baseURL)
}

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

func getTokenSource(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()

	if saInfo.closed {
		return nil, fmt.Errorf("service account token source already closed")
	}

	baseURL = NormalizeBaseURL(baseURL)
	key := saCacheKey(clientID, baseURL)
	if saInfo.sources == nil {
		saInfo.sources = make(map[string]*saCacheEntry)
	}
	if entry, ok := saInfo.sources[key]; ok {
		if entry.clientSecret != clientSecret {
			return nil, fmt.Errorf("service account credentials changed")
		}
		return entry.tokenSource, nil
	}

	tokenSource, err := createTokenSourceFn(clientID, clientSecret, baseURL, terraformVersion)
	if err != nil {
		return nil, err
	}
	saInfo.sources[key] = &saCacheEntry{
		tokenSource:      tokenSource,
		clientID:         clientID,
		clientSecret:     clientSecret,
		baseURL:          baseURL,
		terraformVersion: terraformVersion,
	}
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

// CloseTokenSource is called just before the provider finishes, it does a best-effort try to revoke all cached Service Account tokens.
// It sets saInfo.closed = true to avoid future calls to getTokenSource, that should't happen as the provider is exiting.
func CloseTokenSource() {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()
	if saInfo.closed {
		return
	}
	saInfo.closed = true
	for _, entry := range saInfo.sources {
		if token, err := entry.tokenSource.Token(); err == nil {
			conf := GetServiceAccountConfig(entry.clientID, entry.clientSecret, entry.baseURL)
			ctx := context.WithValue(context.Background(), auth.HTTPClient, NewOAuthHTTPClient(entry.terraformVersion))
			_ = conf.RevokeToken(ctx, token) // Best-effort, no need to do anything if it fails.
		}
	}
}
