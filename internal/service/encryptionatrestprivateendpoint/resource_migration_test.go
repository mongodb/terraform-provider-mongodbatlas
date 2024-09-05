package encryptionatrestprivateendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigEncryptionAtRestPrivateEndpoint_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.19.0")
	testCase := basicTestCase(t)
	mig.CreateAndRunTestNonParallel(t, testCase)
}
