package advancedclustertpf_test

import (
	"testing"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0")
	mig.CreateAndRunTest(t, replicaSetAWSProviderTestCase(t))
}

func TestMigAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0")
	mig.CreateAndRunTest(t, replicaSetMultiCloudTestCase(t))
}

func TestMigAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	mig.SkipIfVersionBelow(t, "2.0.0")
	mig.CreateAndRunTest(t, singleShardedMultiCloudTestCase(t))
}
