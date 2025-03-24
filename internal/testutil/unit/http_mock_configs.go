package unit

import (
	"time"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedclustertpf"
)

const (
	shortRefresh = 100 * time.Millisecond
)

var (
	MockConfigAdvancedClusterTPF = MockHTTPDataConfig{AllowMissingRequests: true, RunBeforeEach: shortenClusterTPFRetries, IsDiffMustSubstrings: []string{"/clusters"}, QueryVars: []string{"providerName"}}
)

func shortenClusterTPFRetries() error {
	advancedclustertpf.RetryMinTimeout = shortRefresh
	advancedclustertpf.RetryDelay = shortRefresh
	advancedclustertpf.RetryPollInterval = shortRefresh
	return nil
}
