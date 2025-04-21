package unit

import (
	"sync"
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

const (
	shortRefresh = 100 * time.Millisecond
)

var (
	MockConfigAdvancedClusterTPF = MockHTTPDataConfig{AllowMissingRequests: true, RunBeforeEach: shortenClusterTPFRetries, IsDiffMustSubstrings: []string{"/clusters"}, QueryVars: []string{"providerName"}}
	onceShortenClusterTPFRetries sync.Once
)

// shortenClusterTPFRetries must meet the interface func() error as it is used in RunBeforeEach which runs as part of TestCase.PreCheck()
func shortenClusterTPFRetries() error {
	onceShortenClusterTPFRetries.Do(func() {
		advancedclustertpf.RetryMinTimeout = shortRefresh
		advancedclustertpf.RetryDelay = shortRefresh
		advancedclustertpf.RetryPollInterval = shortRefresh
	})
	return nil
}
