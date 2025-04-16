package resourcepolicy_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigResourcePolicy_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.33.0") // this feature was GA (no need of MONGODB_ATLAS_ENABLE_PREVIEW env variable) in 1.33.0
	mig.CreateAndRunTestNonParallel(t, basicTestCase(t))
}
