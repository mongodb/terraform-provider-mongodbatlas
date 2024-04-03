package backupcompliancepolicy_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupCompliancePolicy_basic(t *testing.T) {
	useYearly := mig.IsProviderVersionAtLeast("1.16.0") // attribute introduced in this version
	mig.CreateAndRunTest(t, basicTestCase(t, useYearly))
}
