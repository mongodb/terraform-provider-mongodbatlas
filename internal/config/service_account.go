package config

import (
	"context"
	"net/http"
	"strings"
	"sync"

	"github.com/mongodb/atlas-sdk-go/auth"
	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
)

var mu sync.Mutex
var ts auth.TokenSource

func tokenSource(ctx context.Context, c *Config, base http.RoundTripper) (auth.TokenSource, error) {
	mu.Lock()
	defer mu.Unlock()

	if ts != nil {
		return ts, nil
	}

	conf := clientcredentials.NewConfig(c.ClientID, c.ClientSecret)
	// Override TokenURL and RevokeURL if custom BaseURL is provided
	if c.BaseURL != "" {
		baseURL := strings.TrimRight(c.BaseURL, "/")
		conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
		conf.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	}

	// Create a base HTTP client for token acquisition
	baseHTTPClient := &http.Client{
		Transport: base,
	}

	// Set the HTTP client in context for token acquisition
	ctx = context.WithValue(ctx, auth.HTTPClient, baseHTTPClient)

	tokenSource := conf.TokenSource(ctx)

	// Acquire an initial token upfront for several reasons:
	// 1. OAuth2 token caching: The oauth2 library only caches tokens after successful acquisition
	// 2. Early credential validation: Fail fast during provider init rather than first resource operation
	// 3. Performance: Subsequent requests use cached tokens instead of blocking for token acquisition
	if _, err := tokenSource.Token(); err != nil {
		return nil, err
	}
	ts = tokenSource
	return tokenSource, nil
}
