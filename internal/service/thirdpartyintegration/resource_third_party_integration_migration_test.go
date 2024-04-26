package thirdpartyintegration_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigThirdPartyIntegration_basic(t *testing.T) {
	mig.CreateAndRunTest(t, basicOpsGenie(t))
}
