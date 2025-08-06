package cloudbackupsnapshotrestorejob_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.22.0") // this is when the new `failed` field was added
	mig.CreateAndRunTest(t, basicTestCase(t))
}
