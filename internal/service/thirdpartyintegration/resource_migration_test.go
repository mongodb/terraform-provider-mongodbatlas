package thirdpartyintegration_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigThirdPartyIntegration_basic(t *testing.T) {
	// does not run in parallel to reuse same execution project
	mig.CreateAndRunTestNonParallel(t, basicPagerDutyTest(t))
}
