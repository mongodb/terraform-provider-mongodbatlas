package unit

import (
	"sync"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

var (
	shortRefresh                 = 100 * time.Millisecond
	MockConfigAdvancedClusterTPF = MockHTTPDataConfig{AllowMissingRequests: true, SideEffect: shortenClusterTPFRetries, IsDiffMustSubstrings: []string{"/clusters"}, QueryVars: []string{"providerName"}}
	onceShortenClusterTPFRetries sync.Once
)

func shortenClusterTPFRetries() error {
	onceShortenClusterTPFRetries.Do(func() {
		advancedclustertpf.RetryMinTimeout = shortRefresh
		advancedclustertpf.RetryDelay = shortRefresh
		advancedclustertpf.RetryPollInterval = shortRefresh
	})
	return nil
}
