package config

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/mongodb/atlas-sdk-go/auth"
	"github.com/mongodb/atlas-sdk-go/auth/clientcredentials"
	"golang.org/x/oauth2"
)

var saInfo = struct {
	tokenSource  auth.TokenSource
	clientID     string
	clientSecret string
	baseURL      string
	mu           sync.Mutex
}{}

func tokenSource(ctx context.Context, c *Config, base http.RoundTripper) (auth.TokenSource, error) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()

	if saInfo.tokenSource != nil {
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
	ctx = context.WithValue(ctx, auth.HTTPClient, &http.Client{Transport: base})
	token, err := conf.TokenSource(ctx).Token()
	if err != nil {
		return nil, err
	}
	saInfo.clientID = c.ClientID
	saInfo.clientSecret = c.ClientSecret
	saInfo.baseURL = c.BaseURL
	// TODO: token will be refreshed in a follow-up PR
	saInfo.tokenSource = oauth2.StaticTokenSource(token)
	return saInfo.tokenSource, nil
}
