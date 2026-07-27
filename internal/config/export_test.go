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

func ResetSAInfoForTest() {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()
	saInfo.sources = nil
	saInfo.closed = false
}

func SetSAInfoClosedForTest(closed bool) {
	saInfo.mu.Lock()
	defer saInfo.mu.Unlock()
	saInfo.closed = closed
}
