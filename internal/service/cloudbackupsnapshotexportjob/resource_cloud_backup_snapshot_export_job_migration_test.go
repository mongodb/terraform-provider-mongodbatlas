package cloudbackupsnapshotexportjob_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupSnapshotExportJob_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.1")
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
