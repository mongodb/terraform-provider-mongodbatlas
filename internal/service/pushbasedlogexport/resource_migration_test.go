package pushbasedlogexport_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigPushBasedLogExport_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // this feature was introduced in provider version 1.16.0
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}

func TestMigPushBasedLogExport_noPrefixPath(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // this feature was introduced in provider version 1.16.0
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
