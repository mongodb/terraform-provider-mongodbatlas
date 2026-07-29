package config

import "github.com/mongodb/atlas-sdk-go/auth"

// Test helpers exported only for package config_test (see service_account_test.go).

func GetTokenSourceForTest(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	return getTokenSource(clientID, clientSecret, baseURL, terraformVersion)
}

func SetCreateTokenSourceForTest(fn func(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error)) {
	createTokenSourceFn = fn
}

func ResetCreateTokenSourceForTest() {
	createTokenSourceFn = defaultCreateTokenSource
}

// SetRevokeTokenForTest records which clientIDs get revoked without hitting the network.
func SetRevokeTokenForTest(fn func(clientID string)) {
	revokeTokenFn = func(clientID string, _ *saTokenSourceEntry) {
		fn(clientID)
	}
}

func ResetRevokeTokenForTest() {
	revokeTokenFn = defaultRevokeToken
}

func ResetSATokenSourceCacheForTest() {
	saTokenSourceCache.mu.Lock()
	defer saTokenSourceCache.mu.Unlock()
	saTokenSourceCache.sources = nil
	saTokenSourceCache.baseURL = ""
	saTokenSourceCache.terraformVersion = ""
	saTokenSourceCache.closed = false
}

func SetSATokenSourceCacheClosedForTest(closed bool) {
	saTokenSourceCache.mu.Lock()
	defer saTokenSourceCache.mu.Unlock()
	saTokenSourceCache.closed = closed
}
