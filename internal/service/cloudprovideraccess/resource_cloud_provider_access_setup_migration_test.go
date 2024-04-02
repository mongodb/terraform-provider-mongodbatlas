package cloudprovideraccess_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudProviderAccessSetupAWS_basic(t *testing.T) {
	mig.CreateAndRunTest(t, basicSetupTestCase(t))
}
