package unit

import (
	"sync"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
)

const (
	shortRefresh = 100 * time.Millisecond
)

var (
	MockConfigAdvancedCluster = MockHTTPDataConfig{AllowMissingRequests: true, RunBeforeEach: shortenClusterRetries, IsDiffMustSubstrings: []string{"/clusters"}, QueryVars: []string{"providerName"}}
	onceShortenClusterRetries sync.Once
)

// shortenClusterRetries must meet the interface func() error as it is used in RunBeforeEach which runs as part of TestCase.PreCheck()
func shortenClusterRetries() error {
	onceShortenClusterRetries.Do(func() {
		advancedcluster.RetryMinTimeout = shortRefresh
		advancedcluster.RetryDelay = shortRefresh
		advancedcluster.RetryPollInterval = shortRefresh
	})
	return nil
}
