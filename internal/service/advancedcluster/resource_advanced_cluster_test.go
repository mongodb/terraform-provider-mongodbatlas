package advancedcluster_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"
	"time"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20241113003/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/advancedcluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_advanced_cluster.test"
	dataSourceName       = "data.mongodbatlas_advanced_cluster.test"
	dataSourcePluralName = "data.mongodbatlas_advanced_clusters.test"
)

var (
	configServerManagementModeFixedToDedicated = "FIXED_TO_DEDICATED"
	configServerManagementModeAtlasManaged     = "ATLAS_MANAGED"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(t, true, projectID, clusterName),
				Check:  checkTenant(true, projectID, clusterName),
			},
			{
				Config: configTenant(t, true, projectID, clusterNameUpdated),
				Check:  checkTenant(true, projectID, clusterNameUpdated),
			},
			acc.TestStepImportCluster(resourceName),
		},
	})
}

func TestAccClusterAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	resource.ParallelTest(t, replicaSetAWSProviderTestCase(t, true))
}

func replicaSetAWSProviderTestCase(t *testing.T, isAcc bool) resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetAWSProvider(t, isAcc, projectID, clusterName, 60, 3),
				Check:  checkReplicaSetAWSProvider(isAcc, projectID, clusterName, 60, 3, true, true),
			},
			{
				Config: configReplicaSetAWSProvider(t, isAcc, projectID, clusterName, 50, 5),
				Check:  checkReplicaSetAWSProvider(isAcc, projectID, clusterName, 50, 5, true, true),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs", "retain_backups_enabled"),
		},
	}
}

func TestAccClusterAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	resource.ParallelTest(t, replicaSetMultiCloudTestCase(t, true))
}

func replicaSetMultiCloudTestCase(t *testing.T, isAcc bool) resource.TestCase {
	t.Helper()
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetMultiCloud(t, isAcc, orgID, projectName, clusterName),
				Check:  checkReplicaSetMultiCloud(isAcc, clusterName, 3),
			},
			{
				Config: configReplicaSetMultiCloud(t, isAcc, orgID, projectName, clusterNameUpdated),
				Check:  checkReplicaSetMultiCloud(isAcc, clusterNameUpdated, 3),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs", "retain_backups_enabled"),
		},
	}
}

func TestAccClusterAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	resource.ParallelTest(t, singleShardedMultiCloudTestCase(t, true))
}

func singleShardedMultiCloudTestCase(t *testing.T, isAcc bool) resource.TestCase {
	t.Helper()
	// TODO: Already prepared for TPF but getting this error:
	// resource_advanced_cluster_test.go:119: Step 1/3 error: Check failed: Check 9/12 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.0.electable_specs.0.disk_iops' expected to be set
	// Check 10/12 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.0.analytics_specs.0.disk_iops' expected to be set
	// Check 11/12 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.1.electable_specs.0.disk_iops' expected to be set
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaMultiCloud(t, isAcc, orgID, projectName, clusterName, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(isAcc, clusterName, 1, "M10", true, nil),
			},
			{
				Config: configShardedOldSchemaMultiCloud(t, isAcc, orgID, projectName, clusterNameUpdated, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(isAcc, clusterNameUpdated, 1, "M10", true, nil),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"),
		},
	}
}

func TestAccClusterAdvancedCluster_unpausedToPaused(t *testing.T) {
	var (
		projectID           = acc.ProjectIDExecution(t)
		clusterName         = acc.RandomClusterName()
		instanceSize        = "M10"
		anotherInstanceSize = "M20"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, false),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, true),
			},
			{
				Config:      configSingleProviderPaused(t, true, projectID, clusterName, true, anotherInstanceSize),
				ExpectError: regexp.MustCompile("CANNOT_UPDATE_PAUSED_CLUSTER"),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"),
		},
	})
}

func TestAccClusterAdvancedCluster_pausedToUnpaused(t *testing.T) {
	var (
		projectID    = acc.ProjectIDExecution(t)
		clusterName  = acc.RandomClusterName()
		instanceSize = "M10"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, true),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(true, clusterName, false),
			},
			{
				Config:      configSingleProviderPaused(t, true, projectID, clusterName, true, instanceSize),
				ExpectError: regexp.MustCompile("CANNOT_PAUSE_RECENTLY_RESUMED_CLUSTER"),
			},
			{
				Config: configSingleProviderPaused(t, true, projectID, clusterName, false, instanceSize),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"),
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig_oldMongoDBVersion(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// unexpected new value: .advanced_configuration.fail_index_key_too_long: was cty.False, but now null
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()

		processArgs20240530 = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_1"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgs = &admin.ClusterDescriptionProcessArgs20240805{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(-1), // this will not be set in the TF configuration
			DefaultMaxTimeMS: conversion.IntPtr(65),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configAdvanced(t, true, projectID, clusterName, "6.0", processArgs20240530, processArgs),
				ExpectError: regexp.MustCompile(advancedcluster.ErrorDefaultMaxTimeMinVersion),
			},
			{
				Config: configAdvanced(t, true, projectID, clusterName, "6.0", processArgs20240530, &admin.ClusterDescriptionProcessArgs20240805{}),
				Check:  checkAdvanced(true, clusterName, "TLS1_1", &admin.ClusterDescriptionProcessArgs20240805{}),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig(t *testing.T) {
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		projectID           = acc.ProjectIDExecution(t)
		clusterName         = acc.RandomClusterName()
		clusterNameUpdated  = acc.RandomClusterName()
		processArgs20240530 = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_1"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgs = &admin.ClusterDescriptionProcessArgs20240805{
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(-1), // this will not be set in the TF configuration
		}

		processArgs20240530Updated = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("0"),
			FailIndexKeyTooLong:              conversion.Pointer(false),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
		processArgsUpdated = &admin.ClusterDescriptionProcessArgs20240805{
			DefaultMaxTimeMS: conversion.IntPtr(65),
			ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds: conversion.IntPtr(100),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvanced(t, true, projectID, clusterName, "", processArgs20240530, processArgs),
				Check:  checkAdvanced(true, clusterName, "TLS1_1", processArgs),
			},
			{
				Config: configAdvanced(t, true, projectID, clusterNameUpdated, "", processArgs20240530Updated, processArgsUpdated),
				Check:  checkAdvanced(true, clusterNameUpdated, "TLS1_2", processArgsUpdated),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_defaultWrite(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// Check failed: Check 8/14 error: mongodbatlas_advanced_cluster.test: Attribute 'advanced_configuration.fail_index_key_too_long' not found
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
		processArgs        = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_1"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
		}
		processArgsUpdated = &admin20240530.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("majority"),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_2"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
			TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedDefaultWrite(t, true, projectID, clusterName, processArgs),
				Check:  checkAdvancedDefaultWrite(true, clusterName, "1", "TLS1_1"),
			},
			{
				Config: configAdvancedDefaultWrite(t, true, projectID, clusterNameUpdated, processArgsUpdated),
				Check:  checkAdvancedDefaultWrite(true, clusterNameUpdated, "majority", "TLS1_2"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAutoScaling(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "INVALID_ENUM_VALUE") Detail: An invalid enumeration value  was specified. Reason: Bad Request. Params: [],
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
		autoScaling        = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(false), MaxInstanceSize: conversion.StringPtr("")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
		autoScalingUpdated = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(true), MaxInstanceSize: conversion.StringPtr("M20")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAutoScaling(t, true, projectID, clusterName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetSchemaV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "advanced_configuration.0.oplog_min_retention_hours", "5.5"),
				),
			},
			{
				Config: configReplicationSpecsAutoScaling(t, true, projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "name", clusterNameUpdated),
					acc.TestCheckResourceAttrSetSchemaV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAnalyticsAutoScaling(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "INVALID_ENUM_VALUE") Detail: An invalid enumeration value  was specified. Reason: Bad Request. Params: [],
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
		autoScaling        = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(false), MaxInstanceSize: conversion.StringPtr("")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
		autoScalingUpdated = &admin.AdvancedAutoScalingSettings{
			Compute: &admin.AdvancedComputeAutoScaling{Enabled: conversion.Pointer(true), MaxInstanceSize: conversion.StringPtr("M20")},
			DiskGB:  &admin.DiskGBAutoScaling{Enabled: conversion.Pointer(true)},
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "name", clusterName),
					acc.TestCheckResourceAttrSetSchemaV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
				),
			},
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(t, true, projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					acc.CheckExistsCluster(resourceName),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "name", clusterNameUpdated),
					acc.TestCheckResourceAttrSetSchemaV2(true, resourceName, "replication_specs.0.region_configs.#"),
					acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_singleShardedTransitionToOldSchemaExpectsError(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "ASYMMETRIC_REGION_TOPOLOGY_IN_ZONE"). Detail: All shards in the same zone must have the same region topology.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, true, orgID, projectName, clusterName, 1, 1, false),
				Check:  checkGeoShardedOldSchema(true, clusterName, 1, 1, true, true),
			},
			{
				Config:      configGeoShardedOldSchema(t, true, orgID, projectName, clusterName, 1, 2, false),
				ExpectError: regexp.MustCompile(advancedcluster.ErrorOperationNotPermitted),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_withTags(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to check correctly plural data source in the different test steps
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags"),
				Check:  checkKeyValueBlocks(true, clusterName, "tags"),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags", acc.ClusterTagsMap1, acc.ClusterTagsMap2),
				Check:  checkKeyValueBlocks(true, clusterName, "tags", acc.ClusterTagsMap1, acc.ClusterTagsMap2),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "tags", acc.ClusterTagsMap3),
				Check:  checkKeyValueBlocks(true, clusterName, "tags", acc.ClusterTagsMap3),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_withLabels(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to check correctly plural data source in the different test steps
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels"),
				Check:  checkKeyValueBlocks(true, clusterName, "labels"),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels", acc.ClusterLabelsMap1, acc.ClusterLabelsMap2),
				Check:  checkKeyValueBlocks(true, clusterName, "labels", acc.ClusterLabelsMap1, acc.ClusterLabelsMap2),
			},
			{
				Config: configWithKeyValueBlocks(t, true, orgID, projectName, clusterName, "labels", acc.ClusterLabelsMap3),
				Check:  checkKeyValueBlocks(true, clusterName, "labels", acc.ClusterLabelsMap3),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_selfManagedSharding(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "ASYMMETRIC_REGION_TOPOLOGY_IN_ZONE"). Detail: All shards in the same zone must have the same region topology.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
		checks      = []resource.TestCheckFunc{
			acc.CheckExistsCluster(resourceName),
			acc.TestCheckResourceAttrSchemaV2(true, resourceName, "global_cluster_self_managed_sharding", "true"),
			acc.TestCheckResourceAttrSchemaV2(true, dataSourceName, "global_cluster_self_managed_sharding", "true"),
		}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, true, orgID, projectName, clusterName, 1, 1, true),
				Check: resource.ComposeAggregateTestCheckFunc(checks...,
				),
			},
			{
				Config:      configGeoShardedOldSchema(t, true, orgID, projectName, clusterName, 1, 1, false),
				ExpectError: regexp.MustCompile("CANNOT_MODIFY_GLOBAL_CLUSTER_MANAGEMENT_SETTING"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_selfManagedShardingIncorrectType(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configIncorrectTypeGobalClusterSelfManagedSharding(t, true, projectID, clusterName),
				ExpectError: regexp.MustCompile("CANNOT_SET_SELF_MANAGED_SHARDING_FOR_NON_GLOBAL_CLUSTER"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// resource_advanced_cluster_test.go:545: Step 1/2 error: Check failed: Check 3/13 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.num_shards' expected "2", got "1"
	// Check 9/13 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.0.electable_specs.0.disk_iops' expected to be set
	// Check 10/13 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.0.analytics_specs.0.disk_iops' expected to be set
	// Check 11/13 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.region_configs.1.electable_specs.0.disk_iops' expected to be set
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaMultiCloud(t, true, orgID, projectName, clusterName, 2, "M10", &configServerManagementModeFixedToDedicated),
				Check:  checkShardedOldSchemaMultiCloud(true, clusterName, 2, "M10", false, &configServerManagementModeFixedToDedicated),
			},
			{
				Config: configShardedOldSchemaMultiCloud(t, true, orgID, projectName, clusterName, 2, "M20", &configServerManagementModeAtlasManaged),
				Check:  checkShardedOldSchemaMultiCloud(true, clusterName, 2, "M20", false, &configServerManagementModeAtlasManaged),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricGeoShardedOldSchema(t *testing.T) {
	resource.ParallelTest(t, symmetricGeoShardedOldSchemaTestCase(t, true))
}

func symmetricGeoShardedOldSchemaTestCase(t *testing.T, isAcc bool) resource.TestCase {
	t.Helper()
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "INVALID_ENUM_VALUE") Detail: An invalid enumeration value  was specified. Reason: Bad Request. Params: [],
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedOldSchema(t, isAcc, orgID, projectName, clusterName, 2, 2, false),
				Check:  checkGeoShardedOldSchema(isAcc, clusterName, 2, 2, true, false),
			},
			{
				Config: configGeoShardedOldSchema(t, isAcc, orgID, projectName, clusterName, 3, 3, false),
				Check:  checkGeoShardedOldSchema(isAcc, clusterName, 3, 3, true, false),
			},
		},
	}
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// Check failed: Check 2/5 error: mongodbatlas_advanced_cluster.test: Attribute 'replication_specs.0.num_shards' expected \"2\", got \"1\"
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(t, true, orgID, projectName, clusterName, 50),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(true, 50),
			},
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(t, true, orgID, projectName, clusterName, 55),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(true, 55),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedNewSchemaToAsymmetricAddingRemovingShard(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// PATCH: HTTP 400 Bad Request (Error code: \"AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG\") Detail: If any regionConfigs specify an autoScaling object, all regionConfigs must also specify an autoScaling object.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 50, "M10", "M10", nil, nil, false),
				Check:  checkShardedNewSchema(true, 50, "M10", "M10", nil, nil, false, false),
			},
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, true), // add middle replication spec and transition to asymmetric
				Check:  checkShardedNewSchema(true, 55, "M10", "M20", nil, nil, true, true),
			},
			{
				Config: configShardedNewSchema(t, true, orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, false), // removes middle replication spec
				Check:  checkShardedNewSchema(true, 55, "M10", "M20", nil, nil, true, false),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_asymmetricShardedNewSchema(t *testing.T) {
	resource.ParallelTest(t, asymmetricShardedNewSchemaTestCase(t, true))
}

func asymmetricShardedNewSchemaTestCase(t *testing.T, isAcc bool) resource.TestCase {
	t.Helper()
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedNewSchema(t, isAcc, orgID, projectName, clusterName, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), false),
				Check:  checkShardedNewSchema(isAcc, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), true, false),
			},
		},
	}
}

func TestAccClusterAdvancedClusterConfig_asymmetricGeoShardedNewSchemaAddingRemovingShard(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	//  POST: HTTP 400 Bad Request (Error code: "ASYMMETRIC_REGION_TOPOLOGY_IN_ZONE"). Detail: All shards in the same zone must have the same region topology.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedNewSchema(t, true, orgID, projectName, clusterName, false),
				Check:  checkGeoShardedNewSchema(true, false),
			},
			{
				Config: configGeoShardedNewSchema(t, true, orgID, projectName, clusterName, true),
				Check:  checkGeoShardedNewSchema(true, true),
			},
			{
				Config: configGeoShardedNewSchema(t, true, orgID, projectName, clusterName, false),
				Check:  checkGeoShardedNewSchema(true, false),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_shardedTransitionFromOldToNewSchema(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// PATCH: HTTP 400 Bad Request (Error code: "AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG") Detail: If any regionConfigs specify an autoScaling object, all regionConfigs must also specify an autoScaling object.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedTransitionOldToNewSchema(t, true, orgID, projectName, clusterName, false),
				Check:  checkShardedTransitionOldToNewSchema(true, false),
			},
			{
				Config: configShardedTransitionOldToNewSchema(t, true, orgID, projectName, clusterName, true),
				Check:  checkShardedTransitionOldToNewSchema(true, true),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_geoShardedTransitionFromOldToNewSchema(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// POST: HTTP 400 Bad Request (Error code: "ASYMMETRIC_REGION_TOPOLOGY_IN_ZONE"). Detail: All shards in the same zone must have the same region topology.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGeoShardedTransitionOldToNewSchema(t, true, orgID, projectName, clusterName, false),
				Check:  checkGeoShardedTransitionOldToNewSchema(true, false),
			},
			{
				Config: configGeoShardedTransitionOldToNewSchema(t, true, orgID, projectName, clusterName, true),
				Check:  checkGeoShardedTransitionOldToNewSchema(true, true),
			},
		},
	})
}

func TestAccAdvancedCluster_replicaSetScalingStrategyAndRedactClientLogData(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "WORKLOAD_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "WORKLOAD_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "SEQUENTIAL", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "SEQUENTIAL", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "NODE_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(t, true, orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", false),
			},
		},
	})
}

func TestAccAdvancedCluster_replicaSetScalingStrategyAndRedactClientLogDataOldSchema(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "WORKLOAD_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "WORKLOAD_TYPE", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "SEQUENTIAL", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "SEQUENTIAL", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t, true, orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData(true, "NODE_TYPE", false),
			},
		},
	})
}

// TestAccClusterAdvancedCluster_priorityOldSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityOldSchema(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// .replication_specs[0].region_configs[0].electable_specs.node_count: was cty.NumberIntVal(1), but now cty.NumberIntVal(2)
	// .replication_specs[0].region_configs[0].priority: was cty.NumberIntVal(6), but now cty.NumberIntVal(7).
	//  .replication_specs[0].region_configs[0].region_name: was cty.StringVal("US_WEST_2"), but now cty.StringVal("US_EAST_1").
	// .replication_specs[0].region_configs[1].electable_specs.node_count: was cty.NumberIntVal(2), but now cty.NumberIntVal(1).
	// .replication_specs[0].region_configs[1].priority: was cty.NumberIntVal(7), but now cty.NumberIntVal(6).
	// .replication_specs[0].region_configs[1].region_name: was cty.StringVal("US_EAST_1"), but now cty.StringVal("US_WEST_2").
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configPriority(t, true, orgID, projectName, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(t, true, orgID, projectName, clusterName, true, false),
				Check:  acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(t, true, orgID, projectName, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
		},
	})
}

// TestAccClusterAdvancedCluster_priorityNewSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityNewSchema(t *testing.T) {
	// TODO: Already prepared for TPF but getting this error:
	// Error: errorUpdateLegacy. PATCH: HTTP 400 Bad Request (Error code: "AUTO_SCALINGS_MUST_BE_IN_EVERY_REGION_CONFIG") Detail: If any regionConfigs specify an autoScaling object, all regionConfigs must also specify an autoScaling object.
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configPriority(t, true, orgID, projectName, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(t, true, orgID, projectName, clusterName, false, false),
				Check:  acc.TestCheckResourceAttrSchemaV2(true, resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(t, true, orgID, projectName, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_biConnectorConfig(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configBiConnectorConfig(t, true, projectID, clusterName, false),
				Check:  checkTenantBiConnectorConfig(true, projectID, clusterName, false),
			},
			{
				Config: configBiConnectorConfig(t, true, projectID, clusterName, true),
				Check:  checkTenantBiConnectorConfig(true, projectID, clusterName, true),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_pinnedFCVWithVersionUpgradeAndDowngrade(t *testing.T) {
	acc.SkipIfAdvancedClusterV2Schema(t)
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // Using single project to assert plural data source
		clusterName = acc.RandomClusterName()
	)

	now := time.Now()
	// Time 7 days from now, truncated to the beginning of the day
	sevenDaysFromNow := now.AddDate(0, 0, 7).Truncate(24 * time.Hour)
	firstExpirationDate := conversion.TimeToString(sevenDaysFromNow)
	// Time 8 days from now
	eightDaysFromNow := sevenDaysFromNow.AddDate(0, 0, 1)
	updatedExpirationDate := conversion.TimeToString(eightDaysFromNow)
	invalidDateFormat := "invalid"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configFCVPinning(orgID, projectName, clusterName, nil, "7.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 7, nil, nil),
			},
			{ // pins fcv
				Config: configFCVPinning(orgID, projectName, clusterName, &firstExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(firstExpirationDate), admin.PtrInt(7)),
			},
			{ // using incorrect format
				Config:      configFCVPinning(orgID, projectName, clusterName, &invalidDateFormat, "7.0"),
				ExpectError: regexp.MustCompile("expiration_date format is incorrect: " + invalidDateFormat),
			},
			{ // updates expiration date of fcv
				Config: configFCVPinning(orgID, projectName, clusterName, &updatedExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // upgrade mongodb version with fcv pinned
				Config: configFCVPinning(orgID, projectName, clusterName, &updatedExpirationDate, "8.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 8, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // downgrade mongodb version with fcv pinned
				Config: configFCVPinning(orgID, projectName, clusterName, &updatedExpirationDate, "7.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 7, admin.PtrString(updatedExpirationDate), admin.PtrInt(7)),
			},
			{ // unpins fcv
				Config: configFCVPinning(orgID, projectName, clusterName, nil, "7.0"),
				Check:  acc.CheckFCVPinningConfig(resourceName, dataSourceName, dataSourcePluralName, 7, nil, nil),
			},
		},
	})
}

func checkAggr(isAcc bool, attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{acc.CheckExistsCluster(resourceName)}
	checks = acc.AddAttrChecksSchemaV2(isAcc, resourceName, checks, attrsMap)
	checks = acc.AddAttrSetChecksSchemaV2(isAcc, resourceName, checks, attrsSet...)
	checks = acc.AddAttrChecksSchemaV2(isAcc, dataSourceName, checks, attrsMap)
	checks = acc.AddAttrSetChecksSchemaV2(isAcc, dataSourceName, checks, attrsSet...)
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func configTenant(t *testing.T, isAcc bool, projectID, name string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M5"
					}
					provider_name         = "TENANT"
					backing_provider_name = "AWS"
					region_name           = "US_EAST_1"
					priority              = 7
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, projectID, name))
}

func checkTenant(isAcc bool, projectID, name string) resource.TestCheckFunc {
	pluralChecks := acc.AddAttrSetChecksSchemaV2(isAcc, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)
	return checkAggr(isAcc,
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":                           projectID,
			"name":                                 name,
			"termination_protection_enabled":       "false",
			"global_cluster_self_managed_sharding": "false"},
		pluralChecks...)
}

func configWithKeyValueBlocks(t *testing.T, isAcc bool, orgID, projectName, clusterName, blockName string, blocks ...map[string]string) string {
	t.Helper()
	var extraConfig string
	for _, block := range blocks {
		extraConfig += fmt.Sprintf(`
			%[1]s {
				key   = %[2]q
				value = %[3]q
			}
		`, blockName, block["key"], block["value"])
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}
		
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_EAST_1"
				}
			}

			%[4]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, orgID, projectName, clusterName, extraConfig))
}

func checkKeyValueBlocks(isAcc bool, clusterName, blockName string, blocks ...map[string]string) resource.TestCheckFunc {
	const pluralPrefix = "results.0."
	lenStr := strconv.Itoa(len(blocks))
	keyHash := fmt.Sprintf("%s.#", blockName)
	keyStar := fmt.Sprintf("%s.*", blockName)
	checks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSchemaV2(isAcc, resourceName, keyHash, lenStr),
		acc.TestCheckResourceAttrSchemaV2(isAcc, dataSourceName, keyHash, lenStr),
		acc.TestCheckResourceAttrSchemaV2(isAcc, dataSourcePluralName, pluralPrefix+keyHash, lenStr),
	}
	for _, block := range blocks {
		checks = append(checks,
			acc.TestCheckTypeSetElemNestedAttrsSchemaV2(isAcc, resourceName, keyStar, block),
			acc.TestCheckTypeSetElemNestedAttrsSchemaV2(isAcc, dataSourceName, keyStar, block),
			acc.TestCheckTypeSetElemNestedAttrsSchemaV2(isAcc, dataSourcePluralName, pluralPrefix+keyStar, block))
	}
	return checkAggr(isAcc,
		[]string{"project_id"},
		map[string]string{
			"name": clusterName,
		},
		checks...)
}

func configReplicaSetAWSProvider(t *testing.T, isAcc bool, projectID, name string, diskSizeGB, nodeCountElectable int) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			retain_backups_enabled = "true"
			disk_size_gb = %[3]d

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = %[4]d
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}
	`, projectID, name, diskSizeGB, nodeCountElectable))
}

func checkReplicaSetAWSProvider(isAcc bool, projectID, name string, diskSizeGB, nodeCountElectable int, checkDiskSizeGBInnerLevel, checkExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSchemaV2(isAcc, resourceName, "retain_backups_enabled", "true"),
	}
	additionalChecks = append(additionalChecks,
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)))

	if checkDiskSizeGBInnerLevel {
		additionalChecks = append(additionalChecks,
			checkAggr(isAcc, []string{}, map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			}),
		)
	}

	if checkExternalID {
		additionalChecks = append(additionalChecks, acc.TestCheckResourceAttrSetSchemaV2(isAcc, resourceName, "replication_specs.0.external_id"))
	}

	return checkAggr(isAcc,
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":   projectID,
			"disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.node_count": fmt.Sprintf("%d", nodeCountElectable),
			"name": name},
		additionalChecks...,
	)
}

func configIncorrectTypeGobalClusterSelfManagedSharding(t *testing.T, isAcc bool, projectID, name string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q

			cluster_type = "REPLICASET"
			global_cluster_self_managed_sharding = true # invalid, can only by used with GEOSHARDED clusters

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
	`, projectID, name))
}

func configReplicaSetMultiCloud(t *testing.T, isAcc bool, orgID, projectName, name string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "REPLICASET"
			retain_backups_enabled = false

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}

				region_configs {
					read_only_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "GCP"
					priority      = 0
					region_name   = "US_EAST_4"
				}

				region_configs {
					read_only_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "GCP"
					priority      = 0
					region_name   = "NORTH_AMERICA_NORTHEAST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, orgID, projectName, name))
}

func checkReplicaSetMultiCloud(isAcc bool, name string, regionConfigs int) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSchemaV2(isAcc, resourceName, "retain_backups_enabled", "false"),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, resourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, resourceName, "replication_specs.0.external_id"),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourcePluralName, "results.0.replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.name"),
	}
	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id"},
		map[string]string{
			"name": name},
		additionalChecks...,
	)
}

func configShardedOldSchemaMultiCloud(t *testing.T, isAcc bool, orgID, projectName, name string, numShards int, analyticsSize string, configServerManagementMode *string) string {
	t.Helper()
	var rootConfig string
	if configServerManagementMode != nil {
		// valid values: FIXED_TO_DEDICATED or ATLAS_MANAGED (default)
		// only valid for Major version 8 and later
		// cluster must be SHARDED
		rootConfig = fmt.Sprintf(`
		  mongo_db_major_version = "8"
		  config_server_management_mode = %[1]q
		`, *configServerManagementMode)
	}
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}	

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "SHARDED"
			%[6]s

			replication_specs {
				num_shards = %[4]d
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = %[5]q
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "AZURE"
					priority      = 6
					region_name   = "US_EAST_2"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, orgID, projectName, name, numShards, analyticsSize, rootConfig))
}

func checkShardedOldSchemaMultiCloud(isAcc bool, name string, numShards int, analyticsSize string, verifyExternalID bool, configServerManagementMode *string) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		acc.TestCheckResourceAttrWithSchemaV2(isAcc, dataSourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
	}
	if verifyExternalID {
		additionalChecks = append(
			additionalChecks,
			acc.TestCheckResourceAttrSetSchemaV2(isAcc, resourceName, "replication_specs.0.external_id"))
	}
	if configServerManagementMode != nil {
		additionalChecks = append(additionalChecks,
			acc.TestCheckResourceAttrSchemaV2(isAcc, resourceName, "config_server_management_mode", *configServerManagementMode),
			acc.TestCheckResourceAttrSetSchemaV2(isAcc, resourceName, "config_server_type"),
			acc.TestCheckResourceAttrSchemaV2(isAcc, dataSourceName, "config_server_management_mode", *configServerManagementMode),
			acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourceName, "config_server_type"),
			acc.TestCheckResourceAttrSchemaV2(isAcc, dataSourcePluralName, "results.0.config_server_management_mode", *configServerManagementMode),
			acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.config_server_type"))
	}

	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":                           name,
			"replication_specs.0.num_shards": strconv.Itoa(numShards),
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": analyticsSize,
		},
		additionalChecks...)
}

func configSingleProviderPaused(t *testing.T, isAcc bool, projectID, clusterName string, paused bool, instanceSize string) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			paused       = %[3]t
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[4]q
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}
	`, projectID, clusterName, paused, instanceSize))
}

func checkSingleProviderPaused(isAcc bool, name string, paused bool) resource.TestCheckFunc {
	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":   name,
			"paused": strconv.FormatBool(paused)})
}

func configAdvanced(t *testing.T, isAcc bool, projectID, clusterName, mongoDBMajorVersion string, p20240530 *admin20240530.ClusterDescriptionProcessArgs, p *admin.ClusterDescriptionProcessArgs20240805) string {
	t.Helper()
	changeStreamOptionsString := ""
	defaultMaxTimeString := ""
	mongoDBMajorVersionString := ""

	if p != nil {
		if p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != nil && p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != conversion.IntPtr(-1) {
			changeStreamOptionsString = fmt.Sprintf(`change_stream_options_pre_and_post_images_expire_after_seconds = %[1]d`, *p.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)
		}
		if p.DefaultMaxTimeMS != nil {
			defaultMaxTimeString = fmt.Sprintf(`default_max_time_ms = %[1]d`, *p.DefaultMaxTimeMS)
		}
	}
	if mongoDBMajorVersion != "" {
		mongoDBMajorVersionString = fmt.Sprintf(`mongo_db_major_version = %[1]q`, mongoDBMajorVersion)
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"
			%[13]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			advanced_configuration  {
				fail_index_key_too_long              = %[3]t
				javascript_enabled                   = %[4]t
				minimum_enabled_tls_protocol         = %[5]q
				no_table_scan                        = %[6]t
				oplog_size_mb                        = %[7]d
				sample_size_bi_connector			 = %[8]d
				sample_refresh_interval_bi_connector = %[9]d
			    transaction_lifetime_limit_seconds   = %[10]d
			    %[11]s
				%[12]s
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, projectID, clusterName,
		p20240530.GetFailIndexKeyTooLong(), p20240530.GetJavascriptEnabled(), p20240530.GetMinimumEnabledTlsProtocol(), p20240530.GetNoTableScan(),
		p20240530.GetOplogSizeMB(), p20240530.GetSampleSizeBIConnector(), p20240530.GetSampleRefreshIntervalBIConnector(), p20240530.GetTransactionLifetimeLimitSeconds(),
		changeStreamOptionsString, defaultMaxTimeString, mongoDBMajorVersionString))
}

func checkAdvanced(isAcc bool, name, tls string, processArgs *admin.ClusterDescriptionProcessArgs20240805) resource.TestCheckFunc {
	advancedConfig := map[string]string{
		"name": name,
		"advanced_configuration.0.minimum_enabled_tls_protocol":         tls,
		"advanced_configuration.0.fail_index_key_too_long":              "false",
		"advanced_configuration.0.javascript_enabled":                   "true",
		"advanced_configuration.0.no_table_scan":                        "false",
		"advanced_configuration.0.oplog_size_mb":                        "1000",
		"advanced_configuration.0.sample_refresh_interval_bi_connector": "310",
		"advanced_configuration.0.sample_size_bi_connector":             "110",
		"advanced_configuration.0.transaction_lifetime_limit_seconds":   "300",
	}

	if processArgs.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds != nil {
		advancedConfig["advanced_configuration.0.change_stream_options_pre_and_post_images_expire_after_seconds"] = strconv.Itoa(*processArgs.ChangeStreamOptionsPreAndPostImagesExpireAfterSeconds)
	}

	if processArgs.DefaultMaxTimeMS != nil {
		advancedConfig["advanced_configuration.0.default_max_time_ms"] = strconv.Itoa(*processArgs.DefaultMaxTimeMS)
	}

	pluralChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.name"),
	}

	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		advancedConfig,
		pluralChecks...,
	)
}

func configAdvancedDefaultWrite(t *testing.T, isAcc bool, projectID, clusterName string, p *admin20240530.ClusterDescriptionProcessArgs) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			advanced_configuration  {
				javascript_enabled                   = %[3]t
				minimum_enabled_tls_protocol         = %[4]q
				no_table_scan                        = %[5]t
				oplog_size_mb                        = %[6]d
				sample_size_bi_connector			 = %[7]d
				sample_refresh_interval_bi_connector = %[8]d
				default_read_concern                 = %[9]q
				default_write_concern                = %[10]q
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, projectID, clusterName, p.GetJavascriptEnabled(), p.GetMinimumEnabledTlsProtocol(), p.GetNoTableScan(),
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetDefaultReadConcern(), p.GetDefaultWriteConcern()))
}

func checkAdvancedDefaultWrite(isAcc bool, name, writeConcern, tls string) resource.TestCheckFunc {
	pluralChecks := []resource.TestCheckFunc{
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.replication_specs.#"),
		acc.TestCheckResourceAttrSetSchemaV2(isAcc, dataSourcePluralName, "results.0.name"),
	}
	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name": name,
			"advanced_configuration.0.minimum_enabled_tls_protocol":         tls,
			"advanced_configuration.0.default_write_concern":                writeConcern,
			"advanced_configuration.0.default_read_concern":                 "available",
			"advanced_configuration.0.fail_index_key_too_long":              "false",
			"advanced_configuration.0.javascript_enabled":                   "true",
			"advanced_configuration.0.no_table_scan":                        "false",
			"advanced_configuration.0.oplog_size_mb":                        "1000",
			"advanced_configuration.0.sample_refresh_interval_bi_connector": "310",
			"advanced_configuration.0.sample_size_bi_connector":             "110"},
		pluralChecks...)
}

func configReplicationSpecsAutoScaling(t *testing.T, isAcc bool, projectID, clusterName string, p *admin.AdvancedAutoScalingSettings) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
				auto_scaling {
						compute_enabled = %[3]t
						disk_gb_enabled = %[4]t
				compute_max_instance_size = %[5]q
				}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
			advanced_configuration  {
			    oplog_min_retention_hours = 5.5
			}
		}
	`, projectID, clusterName, p.Compute.GetEnabled(), p.DiskGB.GetEnabled(), p.Compute.GetMaxInstanceSize()))
}

func configReplicationSpecsAnalyticsAutoScaling(t *testing.T, isAcc bool, projectID, clusterName string, p *admin.AdvancedAutoScalingSettings) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
				analytics_auto_scaling {
						compute_enabled = %[3]t
						disk_gb_enabled = %[4]t
				compute_max_instance_size = %[5]q
				}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}
	`, projectID, clusterName, p.Compute.GetEnabled(), p.DiskGB.GetEnabled(), p.Compute.GetMaxInstanceSize()))
}

func configGeoShardedOldSchema(t *testing.T, isAcc bool, orgID, projectName, name string, numShardsFirstZone, numShardsSecondZone int, selfManagedSharding bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "GEOSHARDED"
			global_cluster_self_managed_sharding = %[6]t
			disk_size_gb  = 60

			replication_specs {
				zone_name  = "zone n1"
				num_shards = %[4]d

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}

			replication_specs {
				zone_name  = "zone n2"
				num_shards = %[5]d

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}
	`, orgID, projectName, name, numShardsFirstZone, numShardsSecondZone, selfManagedSharding))
}

func checkGeoShardedOldSchema(isAcc bool, name string, numShardsFirstZone, numShardsSecondZone int, isLatestProviderVersion, verifyExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{}

	if verifyExternalID {
		additionalChecks = append(additionalChecks, acc.TestCheckResourceAttrSetSchemaV2(isAcc, resourceName, "replication_specs.0.external_id"))
	}

	if isLatestProviderVersion { // checks that will not apply if doing migration test with older version
		additionalChecks = append(additionalChecks, checkAggr(isAcc,
			[]string{"replication_specs.0.zone_id", "replication_specs.0.zone_id"},
			map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": "60",
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": "60",
			}))
	}

	return checkAggr(isAcc,
		[]string{"project_id", "replication_specs.0.id", "replication_specs.1.id"},
		map[string]string{
			"name":                           name,
			"disk_size_gb":                   "60",
			"replication_specs.0.num_shards": strconv.Itoa(numShardsFirstZone),
			"replication_specs.1.num_shards": strconv.Itoa(numShardsSecondZone),
		},
		additionalChecks...,
	)
}

func configShardedOldSchemaDiskSizeGBElectableLevel(t *testing.T, isAcc bool, orgID, projectName, name string, diskSizeGB int) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "SHARDED"

			replication_specs {
				num_shards = 2

				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
					disk_size_gb  = %[4]d
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 0
					disk_size_gb  = %[4]d
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}
	`, orgID, projectName, name, diskSizeGB))
}

func checkShardedOldSchemaDiskSizeGBElectableLevel(isAcc bool, diskSizeGB int) resource.TestCheckFunc {
	return checkAggr(isAcc,
		[]string{},
		map[string]string{
			"replication_specs.0.num_shards": "2",
			"disk_size_gb":                   fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
		})
}

func configShardedNewSchema(t *testing.T, isAcc bool, orgID, projectName, name string, diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIOPS, lastDiskIOPS *int, includeMiddleSpec bool) string {
	t.Helper()
	var thirdReplicationSpec string
	if includeMiddleSpec {
		thirdReplicationSpec = fmt.Sprintf(`
			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[1]q
						node_count    = 3
						disk_size_gb  = %[2]d
					}
					analytics_specs {
						instance_size = %[1]q
						node_count    = 1
						disk_size_gb  = %[2]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		`, firstInstanceSize, diskSizeGB)
	}
	var firstDiskIOPSAttrs string
	if firstDiskIOPS != nil {
		firstDiskIOPSAttrs = fmt.Sprintf(`
			disk_iops = %d
			ebs_volume_type = "PROVISIONED"
		`, *firstDiskIOPS)
	}
	var lastDiskIOPSAttrs string
	if lastDiskIOPS != nil {
		lastDiskIOPSAttrs = fmt.Sprintf(`
			disk_iops = %d
			ebs_volume_type = "PROVISIONED"
		`, *lastDiskIOPS)
	}
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[4]q
						node_count    = 3
						disk_size_gb  = %[9]d
						%[6]s
					}
					analytics_specs {
						instance_size = %[4]q
						node_count    = 1
						disk_size_gb  = %[9]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}

			%[8]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = %[5]q
						node_count    = 3
						disk_size_gb  = %[9]d
						%[7]s
					}
					analytics_specs {
						instance_size = %[5]q
						node_count    = 1
						disk_size_gb  = %[9]d
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			use_replication_spec_per_shard = true
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			use_replication_spec_per_shard = true
		}
	`, orgID, projectName, name, firstInstanceSize, lastInstanceSize, firstDiskIOPSAttrs, lastDiskIOPSAttrs, thirdReplicationSpec, diskSizeGB))
}

func checkShardedNewSchema(isAcc bool, diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIops, lastDiskIops *int, isAsymmetricCluster, includeMiddleSpec bool) resource.TestCheckFunc {
	amtOfReplicationSpecs := 2
	if includeMiddleSpec {
		amtOfReplicationSpecs = 3
	}

	lastSpecIndex := 1
	if includeMiddleSpec {
		lastSpecIndex = 2
	}

	clusterChecks := map[string]string{
		"disk_size_gb":        fmt.Sprintf("%d", diskSizeGB),
		"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
		"replication_specs.0.region_configs.0.electable_specs.0.instance_size":                              firstInstanceSize,
		fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.instance_size", lastSpecIndex): lastInstanceSize,
		"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb":                               fmt.Sprintf("%d", diskSizeGB),
		fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.disk_size_gb", lastSpecIndex):  fmt.Sprintf("%d", diskSizeGB),
		"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb":                               fmt.Sprintf("%d", diskSizeGB),
		fmt.Sprintf("replication_specs.%d.region_configs.0.analytics_specs.0.disk_size_gb", lastSpecIndex):  fmt.Sprintf("%d", diskSizeGB),
	}
	if firstDiskIops != nil {
		clusterChecks["replication_specs.0.region_configs.0.electable_specs.0.disk_iops"] = fmt.Sprintf("%d", *firstDiskIops)
	}
	if lastDiskIops != nil {
		clusterChecks[fmt.Sprintf("replication_specs.%d.region_configs.0.electable_specs.0.disk_iops", lastSpecIndex)] = fmt.Sprintf("%d", *lastDiskIops)
	}

	// plural data source checks
	pluralChecks := acc.AddAttrSetChecksSchemaV2(isAcc, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.replication_specs.0.region_configs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)

	pluralChecks = acc.AddAttrChecksPrefixSchemaV2(isAcc, dataSourcePluralName, pluralChecks, clusterChecks, "results.0")

	// expected id attribute only if cluster is symmetric
	if isAsymmetricCluster {
		pluralChecks = append(pluralChecks, checkAggr(isAcc, []string{}, map[string]string{
			"replication_specs.0.id": "",
			"replication_specs.1.id": "",
		}))
		pluralChecks = acc.AddAttrChecksSchemaV2(isAcc, dataSourcePluralName, pluralChecks, map[string]string{
			"results.0.replication_specs.0.id": "",
			"results.0.replication_specs.1.id": "",
		})
	} else {
		pluralChecks = append(pluralChecks, checkAggr(isAcc, []string{"replication_specs.0.id", "replication_specs.1.id"}, map[string]string{}))
		pluralChecks = acc.AddAttrSetChecksSchemaV2(isAcc, dataSourcePluralName, pluralChecks, "results.0.replication_specs.0.id", "results.0.replication_specs.1.id")
	}

	return checkAggr(isAcc,
		[]string{"replication_specs.0.external_id", "replication_specs.0.zone_id", "replication_specs.1.external_id", "replication_specs.1.zone_id"},
		clusterChecks,
		pluralChecks...,
	)
}

func configGeoShardedNewSchema(t *testing.T, isAcc bool, orgID, projectName, name string, includeThirdShardInFirstZone bool) string {
	t.Helper()
	var thirdReplicationSpec string
	if includeThirdShardInFirstZone {
		thirdReplicationSpec = `
			replication_specs {
				zone_name  = "zone n1"
				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
		`
	}
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			mongo_db_major_version = "7.0"
			cluster_type   = "GEOSHARDED"
			replication_specs {
				zone_name  = "zone n1"
				region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				}
			}
			%[4]s
			replication_specs {
				zone_name  = "zone n2"
				region_configs {
				electable_specs {
					instance_size = "M20"
					node_count    = 3
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}
			}
    	}
		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			use_replication_spec_per_shard = true
		}
		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			use_replication_spec_per_shard = true
		}
	`, orgID, projectName, name, thirdReplicationSpec))
}

func checkGeoShardedNewSchema(isAcc, includeThirdShardInFirstZone bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if includeThirdShardInFirstZone {
		amtOfReplicationSpecs = 3
	} else {
		amtOfReplicationSpecs = 2
	}
	clusterChecks := map[string]string{
		"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
	}

	return checkAggr(isAcc, []string{}, clusterChecks)
}

func configShardedTransitionOldToNewSchema(t *testing.T, isAcc bool, orgID, projectName, name string, useNewSchema bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
	}
	replicationSpec := fmt.Sprintf(`
		replication_specs {
			%[1]s
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 1
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}
		}
	`, numShardsStr)

	var replicationSpecs string
	if useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[1]s
		`, replicationSpec)
	} else {
		replicationSpecs = replicationSpec
	}

	var dataSourceFlag string
	if useNewSchema {
		dataSourceFlag = `use_replication_spec_per_shard = true`
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"

			%[4]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			%[5]s
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			%[5]s
		}
	`, orgID, projectName, name, replicationSpecs, dataSourceFlag))
}

func checkShardedTransitionOldToNewSchema(isAcc, useNewSchema bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if useNewSchema {
		amtOfReplicationSpecs = 2
	} else {
		amtOfReplicationSpecs = 1
	}
	var checksForNewSchema []resource.TestCheckFunc
	if useNewSchema {
		checksForNewSchema = []resource.TestCheckFunc{
			checkAggr(isAcc, []string{"replication_specs.1.id", "replication_specs.0.external_id", "replication_specs.1.external_id"},
				map[string]string{
					"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
					"replication_specs.1.region_configs.0.electable_specs.0.instance_size": "M10",
					"replication_specs.1.region_configs.0.analytics_specs.0.instance_size": "M10",
				}),
		}
	}

	return checkAggr(isAcc,
		[]string{"replication_specs.0.id"},
		map[string]string{
			"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
			"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": "M10",
		},
		checksForNewSchema...,
	)
}

func configGeoShardedTransitionOldToNewSchema(t *testing.T, isAcc bool, orgID, projectName, name string, useNewSchema bool) string {
	t.Helper()
	var numShardsStr string
	if !useNewSchema {
		numShardsStr = `num_shards = 2`
	}
	replicationSpec := `
		replication_specs {
			%[1]s
			region_configs {
				electable_specs {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs {
					instance_size = "M10"
					node_count    = 1
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = %[2]q
			}
			zone_name = %[3]q
		}
	`

	var replicationSpecs string
	if !useNewSchema {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	} else {
		replicationSpecs = fmt.Sprintf(`
			%[1]s
			%[2]s
			%[3]s
			%[4]s
		`, fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"), fmt.Sprintf(replicationSpec, numShardsStr, "US_EAST_1", "zone 1"),
			fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"), fmt.Sprintf(replicationSpec, numShardsStr, "EU_WEST_1", "zone 2"))
	}

	var dataSourceFlag string
	if useNewSchema {
		dataSourceFlag = `use_replication_spec_per_shard = true`
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "GEOSHARDED"

			%[4]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			%[5]s
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			%[5]s
		}
	`, orgID, projectName, name, replicationSpecs, dataSourceFlag))
}

func checkGeoShardedTransitionOldToNewSchema(isAcc, useNewSchema bool) resource.TestCheckFunc {
	if useNewSchema {
		return checkAggr(isAcc,
			[]string{"replication_specs.0.id", "replication_specs.1.id", "replication_specs.2.id", "replication_specs.3.id",
				"replication_specs.0.external_id", "replication_specs.1.external_id", "replication_specs.2.external_id", "replication_specs.3.external_id",
			},
			map[string]string{
				"replication_specs.#":           "4",
				"replication_specs.0.zone_name": "zone 1",
				"replication_specs.1.zone_name": "zone 1",
				"replication_specs.2.zone_name": "zone 2",
				"replication_specs.3.zone_name": "zone 2",
			},
		)
	}
	return checkAggr(isAcc,
		[]string{"replication_specs.0.id", "replication_specs.1.id"},
		map[string]string{
			"replication_specs.#":           "2",
			"replication_specs.0.zone_name": "zone 1",
			"replication_specs.1.zone_name": "zone 2",
		},
	)
}

func configReplicaSetScalingStrategyAndRedactClientLogData(t *testing.T, isAcc bool, orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"
			replica_set_scaling_strategy = %[4]q
			redact_client_log_data = %[5]t

			replication_specs {
				region_configs {
					electable_specs {
						instance_size ="M10"
						node_count    = 3
						disk_size_gb  = 10
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
						disk_size_gb  = 10
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			use_replication_spec_per_shard = true
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			use_replication_spec_per_shard = true
		}
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData))
}

func configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(t *testing.T, isAcc bool, orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	t.Helper()
	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_project.cluster_project.id
			name = %[3]q
			backup_enabled = false
			cluster_type   = "SHARDED"
			replica_set_scaling_strategy = %[4]q
			redact_client_log_data = %[5]t

			replication_specs {
				num_shards = 2
				region_configs {
					electable_specs {
						instance_size ="M10"
						node_count    = 3
						disk_size_gb  = 10
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
						disk_size_gb  = 10
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData))
}

func checkReplicaSetScalingStrategyAndRedactClientLogData(isAcc bool, replicaSetScalingStrategy string, redactClientLogData bool) resource.TestCheckFunc {
	clusterChecks := map[string]string{
		"replica_set_scaling_strategy": replicaSetScalingStrategy,
		"redact_client_log_data":       strconv.FormatBool(redactClientLogData),
	}

	// plural data source checks
	pluralChecks := acc.AddAttrSetChecksSchemaV2(isAcc, dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replica_set_scaling_strategy", "results.0.redact_client_log_data"}...)

	return checkAggr(isAcc,
		[]string{},
		clusterChecks,
		pluralChecks...,
	)
}

func configPriority(t *testing.T, isAcc bool, orgID, projectName, clusterName string, oldSchema, swapPriorities bool) string {
	t.Helper()
	const (
		config7 = `
			region_configs {
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				electable_specs {
					node_count    = 2
					instance_size = "M10"
				}
			}
		`
		config6 = `
			region_configs {
				provider_name = "AWS"
				priority      = 6
				region_name   = "US_WEST_2"
				electable_specs {
					node_count    = 1
					instance_size = "M10"
				}
			}
		`
	)
	strType, strNumShards, strConfigs := "REPLICASET", "", config7+config6
	if oldSchema {
		strType = "SHARDED"
		strNumShards = "num_shards = 2"
	}
	if swapPriorities {
		strConfigs = config6 + config7
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q
			cluster_type   = %[4]q
			backup_enabled = false
			
			replication_specs {
 					%[5]s
 					%[6]s
			}
		}
	`, orgID, projectName, clusterName, strType, strNumShards, strConfigs))
}

func configBiConnectorConfig(t *testing.T, isAcc bool, projectID, name string, enabled bool) string {
	t.Helper()
	additionalConfig := `
		bi_connector_config {
			enabled = false
		}	
	`
	if enabled {
		additionalConfig = `
			bi_connector_config {
				enabled         = true
				read_preference = "secondary"
			}	
		`
	}

	return acc.ConvertAdvancedClusterToSchemaV2(t, isAcc, fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M10"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}

			%[3]s
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
			depends_on = [mongodbatlas_advanced_cluster.test]
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, projectID, name, additionalConfig))
}

func checkTenantBiConnectorConfig(isAcc bool, projectID, name string, enabled bool) resource.TestCheckFunc {
	attrsMap := map[string]string{
		"project_id": projectID,
		"name":       name,
	}
	if enabled {
		attrsMap["bi_connector_config.0.enabled"] = "true"
		attrsMap["bi_connector_config.0.read_preference"] = "secondary"
	} else {
		attrsMap["bi_connector_config.0.enabled"] = "false"
	}
	return checkAggr(isAcc, nil, attrsMap)
}

func configFCVPinning(orgID, projectName, clusterName string, pinningExpirationDate *string, mongoDBMajorVersion string) string {
	var pinnedFCVAttr string
	if pinningExpirationDate != nil {
		pinnedFCVAttr = fmt.Sprintf(`
		pinned_fcv {
    		expiration_date = %q
  		}
		`, *pinningExpirationDate)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
		
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.test.id
			name         = %[3]q

			cluster_type = "REPLICASET"

			mongo_db_major_version = %[4]q

			%[5]s

			replication_specs {
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 3
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "US_WEST_2"
				}
			}
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
		}
	`, orgID, projectName, clusterName, mongoDBMajorVersion, pinnedFCVAttr)
}
