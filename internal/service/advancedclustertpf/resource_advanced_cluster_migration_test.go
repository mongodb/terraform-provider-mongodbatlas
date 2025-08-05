package advancedclustertpf_test

import (
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, false, false),
				Check:             checkShardedTransitionOldToNewSchema(false, false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configShardedTransitionOldToNewSchema(t, false, projectID, clusterName, true, false),
				Check:                    checkShardedTransitionOldToNewSchema(false, true),
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviders(versionBeforeISSRelease),
				Config:            configGeoShardedTransitionOldToNewSchema(t, false, projectID, clusterName, false),
				Check:             checkGeoShardedTransitionOldToNewSchema(false, false),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, false, projectID, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(false, true),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configGeoShardedTransitionOldToNewSchema(t, true, projectID, clusterName, true),
				Check:                    checkGeoShardedTransitionOldToNewSchema(true, true),
			},
		},
	})
}

// migTest is a helper function to run migration tests:
// - TPF -> TPF: for versions 2.0.0+ (tests same config with older TPF provider vs newer TPF provider)
// - SDKv2 -> TPF: when MONGODB_ATLAS_TEST_SDKV2_TO_TPF=true (tests SDKv2 config vs TPF config with MONGODB_ATLAS_LAST_VERSION=1.39.0)
func migTest(t *testing.T, testCaseFunc func(t *testing.T, usePreviewProvider bool) resource.TestCase) {
	t.Helper()

	if acc.IsTestSDKv2ToTPF() {
		// SDKv2 to TPF migration: first step uses SDKv2, second step uses TPF
		t.Log("Running migration test: SDKv2 to TPF")
		testCase := testCaseFunc(t, false) // Get SDKv2 configuration

		migrationTestCase := resource.TestCase{
			PreCheck:     testCase.PreCheck,
			CheckDestroy: testCase.CheckDestroy,
			ErrorCheck:   testCase.ErrorCheck,
			Steps: []resource.TestStep{
				{
					ExternalProviders: mig.ExternalProviders(),
					Config:            testCase.Steps[0].Config, // SDKv2 config
					Check:             testCase.Steps[0].Check,
				},
				{
					ProtoV6ProviderFactories: testCase.ProtoV6ProviderFactories,
					Config:                   getTPFConfig(t, testCaseFunc),
					Check:                    testCase.Steps[0].Check,
				},
			},
		}
		mig.CreateAndRunTestNonParallel(t, &migrationTestCase)
	} else {
		mig.SkipIfVersionBelow(t, "2.0.0")
		t.Log("Running migration test: TPF to TPF")
		testCase := testCaseFunc(t, true)
		mig.CreateAndRunTest(t, &testCase)
	}
}

func getTPFConfig(t *testing.T, testCaseFunc func(t *testing.T, usePreviewProvider bool) resource.TestCase) string {
	t.Helper()
	tpfTestCase := testCaseFunc(t, true)
	return tpfTestCase.Steps[0].Config
}

// func IsTestSDKv2ToTPF() bool {
// 	env, _ := strconv.ParseBool(os.Getenv("MONGODB_ATLAS_TEST_SDKV2_TO_TPF"))
// 	return env
// }
