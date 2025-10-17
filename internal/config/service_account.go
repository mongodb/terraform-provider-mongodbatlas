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
}{}

func getTokenSource(clientID, clientSecret, baseURL string, tokenRenewalBase http.RoundTripper) (auth.TokenSource, error) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()

	baseURL = NormalizeBaseURL(baseURL)
	if saInfo.tokenSource != nil { // Token source in cache.
		if saInfo.clientID != clientID || saInfo.clientSecret != clientSecret || saInfo.baseURL != baseURL {
			return nil, fmt.Errorf("service account credentials changed")
		}
		return saInfo.tokenSource, nil
	}

	conf := clientcredentials.NewConfig(clientID, clientSecret)
	if baseURL != "" {
		conf.TokenURL = baseURL + clientcredentials.TokenAPIPath
		conf.RevokeURL = baseURL + clientcredentials.RevokeAPIPath
	}
	// Use a new context to avoid "context canceled" errors as the token source is reused and can outlast the callee context.
	ctx := context.WithValue(context.Background(), auth.HTTPClient, &http.Client{Transport: tokenRenewalBase})
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
