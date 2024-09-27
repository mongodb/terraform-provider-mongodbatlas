package globalclusterconfig_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigGlobalClusterConfig_basic(t *testing.T) {
	checkZoneID := mig.IsProviderVersionAtLeast("1.21.0")
	mig.CreateAndRunTest(t, basicTestCase(t, checkZoneID, false))
}
