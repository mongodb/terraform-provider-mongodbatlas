package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mongodb/atlas-sdk-go/auth"
	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
	"golang.org/x/oauth2"
)

const saTokenExpiryBuffer = 10 * time.Minute

var saInfo = struct {
	tokenSource  auth.TokenSource
	clientID     string
	clientSecret string
	baseURL      string
	mu           sync.Mutex
}{}

func getTokenSource(c *Config, tokenRenewalBase http.RoundTripper) (auth.TokenSource, error) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()

	if saInfo.tokenSource != nil { // Token source in cache.
		if saInfo.clientID != c.ClientID || saInfo.clientSecret != c.ClientSecret || saInfo.baseURL != c.BaseURL {
			return nil, fmt.Errorf("service account credentials changed")
		}
		return saInfo.tokenSource, nil
	}

	conf := clientcredentials.NewConfig(c.ClientID, c.ClientSecret)
	if c.BaseURL != "" {
		baseURL := strings.TrimRight(c.BaseURL, "/")
		conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
		conf.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	}
	// Use a new context to avoid "context canceled" errors as the token source is reused and can outlast the callee context.
	ctx := context.WithValue(context.Background(), auth.HTTPClient, &http.Client{Transport: tokenRenewalBase})
	tokenSource := oauth2.ReuseTokenSourceWithExpiry(nil, conf.TokenSource(ctx), saTokenExpiryBuffer)
	if _, err := tokenSource.Token(); err != nil { // Retrieve token to fail-fast if credentials are invalid.
		return nil, err
	}
	saInfo.clientID = c.ClientID
	saInfo.clientSecret = c.ClientSecret
	saInfo.baseURL = c.BaseURL
	saInfo.tokenSource = tokenSource
	return saInfo.tokenSource, nil
}
