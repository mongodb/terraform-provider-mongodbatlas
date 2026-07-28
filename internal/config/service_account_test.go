package config_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/mongodb/atlas-sdk-go/auth"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

type staticTokenSource struct {
	token *oauth2.Token
	err   error
}

func (s staticTokenSource) Token() (*oauth2.Token, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.token, nil
}

func resetSAInfo(t *testing.T) {
	t.Helper()
	config.ResetSAInfoForTest()
	t.Cleanup(func() {
		config.ResetSAInfoForTest()
		config.ResetCreateTokenSourceForTest()
		config.ResetRevokeTokenForTest()
	})
}

func TestGetTokenSource_AllowsDifferentClientIDs(t *testing.T) {
	resetSAInfo(t)

	var mu sync.Mutex
	created := make([]string, 0, 2)
	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		mu.Lock()
		created = append(created, clientID)
		mu.Unlock()
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok-" + clientID, TokenType: "Bearer"}}, nil
	})

	ts1, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	ts2, err := config.GetTokenSourceForTest("client-b", "secret-b", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)

	tok1, err := ts1.Token()
	require.NoError(t, err)
	tok2, err := ts2.Token()
	require.NoError(t, err)
	assert.Equal(t, "tok-client-a", tok1.AccessToken)
	assert.Equal(t, "tok-client-b", tok2.AccessToken)
	assert.Equal(t, []string{"client-a", "client-b"}, created)
}

func TestGetTokenSource_ReusesCacheForSameClientID(t *testing.T) {
	resetSAInfo(t)

	calls := 0
	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		calls++
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok", TokenType: "Bearer"}}, nil
	})

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestGetTokenSource_RejectsSecretChangeForSameClientID(t *testing.T) {
	resetSAInfo(t)

	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok", TokenType: "Bearer"}}, nil
	})

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-a", "secret-b", "https://cloud-qa.mongodb.com", "1.0.0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service account credentials changed")
}

func TestGetTokenSource_RejectsBaseURLChangeForSameClientID(t *testing.T) {
	resetSAInfo(t)

	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok", TokenType: "Bearer"}}, nil
	})

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud.mongodb.com", "1.0.0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "service account credentials changed")
}

func TestGetTokenSource_ReusesCacheWhenBaseURLDiffersOnlyByTrailingSlash(t *testing.T) {
	resetSAInfo(t)

	calls := 0
	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		calls++
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok", TokenType: "Bearer"}}, nil
	})

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com/", "1.0.0")
	require.NoError(t, err)
	assert.Equal(t, 1, calls)
}

func TestCloseTokenSource_RevokesAllCachedSources(t *testing.T) {
	resetSAInfo(t)

	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		return staticTokenSource{token: &oauth2.Token{AccessToken: "tok-" + clientID, TokenType: "Bearer"}}, nil
	})

	var mu sync.Mutex
	revoked := make([]string, 0, 2)
	config.SetRevokeTokenForTest(func(clientID string) {
		mu.Lock()
		revoked = append(revoked, clientID)
		mu.Unlock()
	})

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)
	_, err = config.GetTokenSourceForTest("client-b", "secret-b", "https://cloud-qa.mongodb.com", "1.0.0")
	require.NoError(t, err)

	config.CloseTokenSource()
	assert.ElementsMatch(t, []string{"client-a", "client-b"}, revoked)
}

func TestGetTokenSource_ErrorsWhenClosed(t *testing.T) {
	resetSAInfo(t)
	config.SetCreateTokenSourceForTest(func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
		return nil, errors.New("should not be called")
	})
	config.SetSAInfoClosedForTest(true)

	_, err := config.GetTokenSourceForTest("client-a", "secret-a", "", "1.0.0")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already closed")
}
