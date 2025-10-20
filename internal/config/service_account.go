package config

import (
	"context"
	"fmt"
	"log"
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
	log.Println("Service Account: Token source created")
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

func CloseTokenSource() {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()
	if saInfo.tokenSource == nil {
		return
	}
	conf := getConfig(saInfo.clientID, saInfo.clientSecret, saInfo.baseURL)
	token, err := saInfo.tokenSource.Token()
	saInfo.closed = true
	saInfo.tokenSource = nil
	saInfo.clientID = ""
	saInfo.clientSecret = ""
	saInfo.baseURL = ""
	if err != nil {
		log.Printf("Service Account: Error retrieving token to revoke: %v\n", err)
		return
	}
	if err := conf.RevokeToken(context.Background(), token); err != nil {
		log.Printf("Service Account: Error revoking token: %v\n", err)
		return
	}
	log.Println("Service Account: Token revoked successfully")
}
