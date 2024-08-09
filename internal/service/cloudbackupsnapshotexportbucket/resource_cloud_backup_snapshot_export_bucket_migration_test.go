package cloudbackupsnapshotexportbucket_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigBackupSnapshotExportBucket_basic(t *testing.T) {
	mig.CreateTestAndRunUseExternalProviderNonParallel(t, basicAWSTestCase(t), mig.ExternalProvidersWithAWS(), nil)
}
