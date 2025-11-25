package encryptionatrestprivateendpoint_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigEncryptionAtRestPrivateEndpoint_Azure_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.19.0")
	testCase := basicTestCaseAzure(t)
	mig.CreateAndRunTestNonParallel(t, testCase)
}

func TestMigEncryptionAtRestPrivateEndpoint_AWS_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.28.0")
	testCase := basicTestCaseAWS(t)
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, testCase, mig.ExternalProvidersWithAWS(), nil)
}
