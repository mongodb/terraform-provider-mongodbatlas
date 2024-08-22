package encryptionatrestprivateendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigEncryptionAtRestPrivateEndpoint_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.18.1")
	testCase := basicTestCase(t)
	mig.CreateAndRunTest(t, testCase)
}
