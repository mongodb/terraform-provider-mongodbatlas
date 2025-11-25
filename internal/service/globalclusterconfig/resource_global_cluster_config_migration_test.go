package globalclusterconfig_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGlobalClusterConfig_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0")
	mig.CreateAndRunTest(t, basicTestCase(t, false))
}
