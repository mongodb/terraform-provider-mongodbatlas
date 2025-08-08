package advancedcluster_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

// last version that did not support new sharding schema or attributes
const versionBeforeISSRelease = "1.17.6"

func TestMigAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	migTest(t, replicaSetAWSProviderTestCase)
}

func TestMigAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	migTest(t, replicaSetMultiCloudTestCase)
}

func TestMigAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	migTest(t, singleShardedMultiCloudTestCase)
}

func TestMigAdvancedCluster_symmetricGeoShardedOldSchema(t *testing.T) {
	migTest(t, symmetricGeoShardedOldSchemaTestCase)
}

func TestMigAdvancedCluster_asymmetricShardedNewSchema(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.23.0") // version where sharded cluster tier auto-scaling was introduced
	migTest(t, asymmetricShardedNewSchemaTestCase)
}

func TestMigAdvancedCluster_shardedMigrationFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, false, false),
				Check:             checkShardedTransitionOldToNewSchema(false, false),
			},
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, true, false),
				Check:             checkShardedTransitionOldToNewSchema(false, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true, false),
				Check:                    checkShardedTransitionOldToNewSchema(true, true),
			},
		},
	})
}

func TestMigAdvancedCluster_geoShardedMigrationFromOldToNewSchema(t *testing.T) {
	projectID, clusterName := acc.ProjectIDExecutionWithCluster(t, 8)
	versionBeforeTPFGARelease := os.Getenv("MONGODB_ATLAS_LAST_1X_VERSION")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t); mig.PreCheckLast1XVersion(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, false, projectID, clusterName, false),
				Check:             checkGeoShardedTransitionOldToNewSchema(false, false),
			},
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeTPFGARelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, false, projectID, clusterName, true),
				Check:             checkGeoShardedTransitionOldToNewSchema(false, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(true, true),
			},
		},
	})
}

// migTest is a helper function to run migration tests using existing test case functions:
// TPF -> TPF: for versions 2.0.0+ (tests same config with older TPF provider vs newer TPF provider)
// SDKv2 -> TPF: when MONGODB_ATLAS_TEST_SDKV2_TO_TPF=true (tests SDKv2 config vs TPF config with MONGODB_ATLAS_LAST_VERSION = MONGODB_ATLAS_LAST_1X_VERSION)
func migTest(t *testing.T, testCaseFunc func(t *testing.T, useSDKv2 ...bool) resource.TestCase) {
	t.Helper()

	if acc.IsTestSDKv2ToTPF() {
		t.Log("Running migration test: SDKv2 to TPF")
		sdkv2TestCase := testCaseFunc(t, true)

		// Step 1: SDKv2 config with external provider (latest v1.x.x version)
		// Step 2: Converted TPF config with current provider (latest v2.x.x version)
		mig.CreateAndRunTest(t, &sdkv2TestCase)
	} else {
		mig.SkipIfVersionBelow(t, "2.0.0")
		t.Log("Running migration test: TPF to TPF")
		testCase := testCaseFunc(t)
		mig.CreateAndRunTest(t, &testCase)
	}
}
