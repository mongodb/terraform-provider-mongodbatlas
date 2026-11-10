package pushbasedlogexport_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

// Push-based log export migration tests do not run in CI:
//   * PBLE is deprecated in favor of log_integration and will be removed in the next major version.
// 	 * Atlas no longer accepts new PBLE configurations.
// You may run these locally passing a project ID via env vars.

func TestMigPushBasedLogExport_basic(t *testing.T) {
	acc.SkipTestForCI(t)
	mig.SkipIfVersionBelow(t, "1.16.0") // this feature was introduced in provider version 1.16.0
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}

func TestMigPushBasedLogExport_noPrefixPath(t *testing.T) {
	acc.SkipTestForCI(t)
	mig.SkipIfVersionBelow(t, "1.16.0") // this feature was introduced in provider version 1.16.0
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, noPrefixPathTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
