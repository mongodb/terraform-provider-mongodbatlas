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

// Renew token if it expires within 10 minutes to avoid authentication errors during Atlas API calls.
const saTokenExpiryBuffer = 10 * time.Minute

var saInfo = struct {
	tokenSource  auth.TokenSource
	clientID     string
	clientSecret string
	baseURL      string
	mu           sync.Mutex
	closed       bool
}{}

func getTokenSource(clientID, clientSecret, baseURL string, tokenRenewalBase http.RoundTripper) (auth.TokenSource, error) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()

	if saInfo.closed {
		return nil, fmt.Errorf("service account token source already closed")
	}

	baseURL = NormalizeBaseURL(baseURL)
	if saInfo.tokenSource != nil { // Token source in cache.
		if saInfo.clientID != clientID || saInfo.clientSecret != clientSecret || saInfo.baseURL != baseURL {
			return nil, fmt.Errorf("service account credentials changed")
		}
		return saInfo.tokenSource, nil
	}

	// Use a new context to avoid "context canceled" errors as the token source is reused and can outlast the callee context.
	ctx := context.WithValue(context.Background(), auth.HTTPClient, &http.Client{Transport: tokenRenewalBase})
	conf := getConfig(clientID, clientSecret, baseURL)
	tokenSource := oauth2.ReuseTokenSourceWithExpiry(nil, conf.TokenSource(ctx), saTokenExpiryBuffer)
	if _, err := tokenSource.Token(); err != nil { // Retrieve token to fail-fast if credentials are invalid.
		return nil, err
	}
	saInfo.clientID = clientID
	saInfo.clientSecret = clientSecret
	saInfo.baseURL = baseURL
	saInfo.tokenSource = tokenSource
	return saInfo.tokenSource, nil
}

func NormalizeBaseURL(baseURL string) string {
	return strings.TrimRight(baseURL, "/")
}

func getConfig(clientID, clientSecret, baseURL string) *clientcredentials.Config {
	config := clientcredentials.NewConfig(clientID, clientSecret)
	if baseURL != "" {
		config.TokenURL = baseURL + clientcredentials.TokenAPIPath
		config.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	}
	return config
}

// CloseTokenSource is called just before the provider finishes, it does a best-effort try to revoke the Service Access token.
// It sets saInfo.closed = true to avoid future calls to getTokenSource, that should't happen as the provider is exiting.
func CloseTokenSource() {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()
	if saInfo.closed {
		return
	}
	saInfo.closed = true
	if saInfo.tokenSource == nil { // No need to do anything if SA was not initialized.
		return
	}
	if token, err := saInfo.tokenSource.Token(); err == nil {
		conf := getConfig(saInfo.clientID, saInfo.clientSecret, saInfo.baseURL)
		_ = conf.RevokeToken(context.Background(), token) // Best-effort, no need to do anything if it fails.
	}
}
