package cloudbackupsnapshotrestorejob_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.29.0") // version when advanced cluster TPF was introduced
	mig.CreateAndRunTest(t, basicTestCase(t))
}
