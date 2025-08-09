package cloudbackupsnapshotexportjob_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupSnapshotExportJob_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.29.0") // version when advanced cluster TPF was introduced
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
