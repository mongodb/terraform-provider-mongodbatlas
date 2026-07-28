package config

import "github.com/mongodb/atlas-sdk-go/auth"

// Hooks for service_account_test.go (package config_test).

func GetTokenSourceForTest(clientID, clientSecret, baseURL, terraformVersion string) (auth.TokenSource, error) {
	return getTokenSource(clientID, clientSecret, baseURL, terraformVersion)
}

func ResetSATokenCacheForTest() {
	saTokenCache.mu.Lock()
	defer saTokenCache.mu.Unlock()
	saTokenCache.closed = false
	saTokenCache.entries = make(map[saCacheKey]saCacheEntry)
}

func SATokenCacheLenForTest() int {
	saTokenCache.mu.Lock()
	defer saTokenCache.mu.Unlock()
	return len(saTokenCache.entries)
}
