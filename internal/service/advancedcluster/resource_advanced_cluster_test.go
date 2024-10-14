package advancedcluster_test

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	admin20240530 "go.mongodb.org/atlas-sdk/v20240530005/admin"
	"go.mongodb.org/atlas-sdk/v20240805004/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(projectID, clusterName),
				Check:  checkTenant(projectID, clusterName),
			},
			{
				Config: configTenant(projectID, clusterNameUpdated),
				Check:  checkTenant(projectID, clusterNameUpdated),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: acc.ImportStateClusterIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccClusterAdvancedCluster_replicaSetAWSProvider(t *testing.T) {
	resource.ParallelTest(t, replicaSetAWSProviderTestCase(t))
}
func replicaSetAWSProviderTestCase(t *testing.T) resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	return resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicaSetAWSProvider(projectID, clusterName, 60, 3),
				Check:  checkReplicaSetAWSProvider(projectID, clusterName, 60, 3, true, true),
			},
			{
				Config: configReplicaSetAWSProvider(projectID, clusterName, 50, 5),
				Check:  checkReplicaSetAWSProvider(projectID, clusterName, 50, 5, true, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	}
}

func TestAccClusterAdvancedCluster_replicaSetMultiCloud(t *testing.T) {
	resource.ParallelTest(t, replicaSetMultiCloudTestCase(t))
}
func replicaSetMultiCloudTestCase(t *testing.T) resource.TestCase {
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
				Config: configReplicaSetMultiCloud(orgID, projectName, clusterName),
				Check:  checkReplicaSetMultiCloud(clusterName, 3),
			},
			{
				Config: configReplicaSetMultiCloud(orgID, projectName, clusterNameUpdated),
				Check:  checkReplicaSetMultiCloud(clusterNameUpdated, 3),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	}
}

func TestAccClusterAdvancedCluster_singleShardedMultiCloud(t *testing.T) {
	resource.ParallelTest(t, singleShardedMultiCloudTestCase(t))
}

func singleShardedMultiCloudTestCase(t *testing.T) resource.TestCase {
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
				Config: configShardedOldSchemaMultiCloud(orgID, projectName, clusterName, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(clusterName, 1, "M10", true, nil),
			},
			{
				Config: configShardedOldSchemaMultiCloud(orgID, projectName, clusterNameUpdated, 1, "M10", nil),
				Check:  checkShardedOldSchemaMultiCloud(clusterNameUpdated, 1, "M10", true, nil),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(clusterName, false),
			},
			{
				Config: configSingleProviderPaused(projectID, clusterName, true, instanceSize),
				Check:  checkSingleProviderPaused(clusterName, true),
			},
			{
				Config:      configSingleProviderPaused(projectID, clusterName, true, anotherInstanceSize),
				ExpectError: regexp.MustCompile("CANNOT_UPDATE_PAUSED_CLUSTER"),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(projectID, clusterName, true, instanceSize),
				Check:  checkSingleProviderPaused(clusterName, true),
			},
			{
				Config: configSingleProviderPaused(projectID, clusterName, false, instanceSize),
				Check:  checkSingleProviderPaused(clusterName, false),
			},
			{
				Config:      configSingleProviderPaused(projectID, clusterName, true, instanceSize),
				ExpectError: regexp.MustCompile("CANNOT_PAUSE_RECENTLY_RESUMED_CLUSTER"),
			},
			{
				Config: configSingleProviderPaused(projectID, clusterName, false, instanceSize),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
		processArgs        = &admin20240530.ClusterDescriptionProcessArgs{
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
		processArgsUpdated = &admin20240530.ClusterDescriptionProcessArgs{
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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvanced(projectID, clusterName, processArgs, nil, nil),
				Check:  checkAdvanced(clusterName, "TLS1_1", "-1"),
			},
			{
				Config: configAdvanced(projectID, clusterNameUpdated, processArgsUpdated, conversion.IntPtr(100), nil),
				Check:  checkAdvanced(clusterNameUpdated, "TLS1_2", "100"),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_advancedConfig_MongoDBVersion5(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		processArgs = &admin20240530.ClusterDescriptionProcessArgs{
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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvanced(projectID, clusterName, processArgs, nil, conversion.StringPtr("5")),
				Check:  checkAdvanced(clusterName, "TLS1_1", "-1"),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_defaultWrite(t *testing.T) {
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedDefaultWrite(projectID, clusterName, processArgs),
				Check:  checkAdvancedDefaultWrite(clusterName, "1", "TLS1_1"),
			},
			{
				Config: configAdvancedDefaultWrite(projectID, clusterNameUpdated, processArgsUpdated),
				Check:  checkAdvancedDefaultWrite(clusterNameUpdated, "majority", "TLS1_2"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAutoScaling(t *testing.T) {
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAutoScaling(projectID, clusterName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_min_retention_hours", "5.5"),
				),
			},
			{
				Config: configReplicationSpecsAutoScaling(projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAnalyticsAutoScaling(t *testing.T) {
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(projectID, clusterName, autoScaling),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
				),
			},
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_singleShardedTransitionToOldSchemaExpectsError(t *testing.T) {
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
				Config: configGeoShardedOldSchema(orgID, projectName, clusterName, 1, 1, false),
				Check:  checkGeoShardedOldSchema(clusterName, 1, 1, true, true),
			},
			{
				Config:      configGeoShardedOldSchema(orgID, projectName, clusterName, 1, 2, false),
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
				Config: configWithTags(orgID, projectName, clusterName),
				Check:  checkTags(clusterName),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, acc.ClusterTagsMap1, acc.ClusterTagsMap2),
				Check:  checkTags(clusterName, acc.ClusterTagsMap1, acc.ClusterTagsMap2),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, acc.ClusterTagsMap3),
				Check:  checkTags(clusterName, acc.ClusterTagsMap3),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_selfManagedSharding(t *testing.T) {
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
				Config: configGeoShardedOldSchema(orgID, projectName, clusterName, 1, 1, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "global_cluster_self_managed_sharding", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "global_cluster_self_managed_sharding", "true"),
				),
			},
			{
				Config:      configGeoShardedOldSchema(orgID, projectName, clusterName, 1, 1, false),
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
		PreCheck:                 acc.PreCheckBasicSleep(t),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configIncorrectTypeGobalClusterSelfManagedSharding(projectID, clusterName),
				ExpectError: regexp.MustCompile("CANNOT_SET_SELF_MANAGED_SHARDING_FOR_NON_GLOBAL_CLUSTER"),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchema(t *testing.T) {
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
				Config: configShardedOldSchemaMultiCloud(orgID, projectName, clusterName, 2, "M10", &configServerManagementModeFixedToDedicated),
				Check:  checkShardedOldSchemaMultiCloud(clusterName, 2, "M10", false, &configServerManagementModeFixedToDedicated),
			},
			{
				Config: configShardedOldSchemaMultiCloud(orgID, projectName, clusterName, 2, "M20", &configServerManagementModeAtlasManaged),
				Check:  checkShardedOldSchemaMultiCloud(clusterName, 2, "M20", false, &configServerManagementModeAtlasManaged),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricGeoShardedOldSchema(t *testing.T) {
	resource.ParallelTest(t, symmetricGeoShardedOldSchemaTestCase(t))
}

func symmetricGeoShardedOldSchemaTestCase(t *testing.T) resource.TestCase {
	t.Helper()
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
				Config: configGeoShardedOldSchema(orgID, projectName, clusterName, 2, 2, false),
				Check:  checkGeoShardedOldSchema(clusterName, 2, 2, true, false),
			},
			{
				Config: configGeoShardedOldSchema(orgID, projectName, clusterName, 3, 3, false),
				Check:  checkGeoShardedOldSchema(clusterName, 3, 3, true, false),
			},
		},
	}
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) {
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
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, clusterName, 50),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(50),
			},
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, clusterName, 55),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(55),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_symmetricShardedNewSchemaToAsymmetricAddingRemovingShard(t *testing.T) {
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
				Config: configShardedNewSchema(orgID, projectName, clusterName, 50, "M10", "M10", nil, nil, false),
				Check:  checkShardedNewSchema(50, "M10", "M10", nil, nil, false, false),
			},
			{
				Config: configShardedNewSchema(orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, true), // add middle replication spec and transition to asymmetric
				Check:  checkShardedNewSchema(55, "M10", "M20", nil, nil, true, true),
			},
			{
				Config: configShardedNewSchema(orgID, projectName, clusterName, 55, "M10", "M20", nil, nil, false), // removes middle replication spec
				Check:  checkShardedNewSchema(55, "M10", "M20", nil, nil, true, false),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_asymmetricShardedNewSchema(t *testing.T) {
	resource.ParallelTest(t, asymmetricShardedNewSchemaTestCase(t))
}

func asymmetricShardedNewSchemaTestCase(t *testing.T) resource.TestCase {
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
				Config: configShardedNewSchema(orgID, projectName, clusterName, 50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), false),
				Check:  checkShardedNewSchema(50, "M30", "M40", admin.PtrInt(2000), admin.PtrInt(2500), true, false),
			},
		},
	}
}

func TestAccClusterAdvancedClusterConfig_asymmetricGeoShardedNewSchemaAddingRemovingShard(t *testing.T) {
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
				Config: configGeoShardedNewSchema(orgID, projectName, clusterName, false),
				Check:  checkGeoShardedNewSchema(false),
			},
			{
				Config: configGeoShardedNewSchema(orgID, projectName, clusterName, true),
				Check:  checkGeoShardedNewSchema(true),
			},
			{
				Config: configGeoShardedNewSchema(orgID, projectName, clusterName, false),
				Check:  checkGeoShardedNewSchema(false),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_shardedTransitionFromOldToNewSchema(t *testing.T) {
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
				Config: configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, false),
				Check:  checkShardedTransitionOldToNewSchema(false),
			},
			{
				Config: configShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				Check:  checkShardedTransitionOldToNewSchema(true),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_geoShardedTransitionFromOldToNewSchema(t *testing.T) {
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
				Config: configGeoShardedTransitionOldToNewSchema(orgID, projectName, clusterName, false),
				Check:  checkGeoShardedTransitionOldToNewSchema(false),
			},
			{
				Config: configGeoShardedTransitionOldToNewSchema(orgID, projectName, clusterName, true),
				Check:  checkGeoShardedTransitionOldToNewSchema(true),
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
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(orgID, projectName, clusterName, "WORKLOAD_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("WORKLOAD_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(orgID, projectName, clusterName, "SEQUENTIAL", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("SEQUENTIAL", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(orgID, projectName, clusterName, "NODE_TYPE", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("NODE_TYPE", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogData(orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("NODE_TYPE", false),
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
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(orgID, projectName, clusterName, "WORKLOAD_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("WORKLOAD_TYPE", false),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(orgID, projectName, clusterName, "SEQUENTIAL", true),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("SEQUENTIAL", true),
			},
			{
				Config: configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(orgID, projectName, clusterName, "NODE_TYPE", false),
				Check:  checkReplicaSetScalingStrategyAndRedactClientLogData("NODE_TYPE", false),
			},
		},
	})
}

// TestAccClusterAdvancedCluster_priorityOldSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityOldSchema(t *testing.T) {
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
				Config:      configPriority(orgID, projectName, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(orgID, projectName, clusterName, true, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(orgID, projectName, clusterName, true, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
		},
	})
}

// TestAccClusterAdvancedCluster_priorityNewSchema will be able to be simplied or deleted in CLOUDP-275825
func TestAccClusterAdvancedCluster_priorityNewSchema(t *testing.T) {
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
				Config:      configPriority(orgID, projectName, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
			{
				Config: configPriority(orgID, projectName, clusterName, false, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.#", "2"),
			},
			{
				Config:      configPriority(orgID, projectName, clusterName, false, true),
				ExpectError: regexp.MustCompile("priority values in region_configs must be in descending order"),
			},
		},
	})
}

func checkAggr(attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{checkExists(resourceName)}
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	checks = acc.AddAttrChecks(dataSourceName, checks, attrsMap)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	checks = acc.AddAttrSetChecks(dataSourceName, checks, attrsSet...)
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		err := acc.CheckClusterExistsHandlingRetry(ids["project_id"], ids["cluster_name"])
		if err == nil {
			return nil
		}
		return fmt.Errorf("cluster(%s:%s) does not exist: %w", rs.Primary.Attributes["project_id"], rs.Primary.ID, err)
	}
}

func configTenant(projectID, name string) string {
	return fmt.Sprintf(`
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
	`, projectID, name)
}

func checkTenant(projectID, name string) resource.TestCheckFunc {
	pluralChecks := acc.AddAttrSetChecks(dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)
	return checkAggr(
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":                           projectID,
			"name":                                 name,
			"termination_protection_enabled":       "false",
			"global_cluster_self_managed_sharding": "false",
			"labels.#":                             "0"},
		pluralChecks...)
}

func configWithTags(orgID, projectName, name string, tags ...map[string]string) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, label["key"], label["value"])
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
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
	`, orgID, projectName, name, tagsConf)
}

func checkTags(name string, tags ...map[string]string) resource.TestCheckFunc {
	lenStr := strconv.Itoa(len(tags))
	tagChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "tags.#", lenStr),
		resource.TestCheckResourceAttr(dataSourceName, "tags.#", lenStr),
		resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", lenStr),
	}
	for _, tag := range tags {
		tagChecks = append(tagChecks,
			resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", tag),
			resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", tag),
			resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", tag))
	}
	return checkAggr(
		[]string{"project_id"},
		map[string]string{
			"name": name,
		},
		tagChecks...)
}

func configReplicaSetAWSProvider(projectID, name string, diskSizeGB, nodeCountElectable int) string {
	return fmt.Sprintf(`
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
	`, projectID, name, diskSizeGB, nodeCountElectable)
}

func checkReplicaSetAWSProvider(projectID, name string, diskSizeGB, nodeCountElectable int, checkDiskSizeGBInnerLevel, checkExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
	}
	if checkDiskSizeGBInnerLevel {
		additionalChecks = append(additionalChecks,
			checkAggr([]string{}, map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			}),
		)
	}

	if checkExternalID {
		additionalChecks = append(additionalChecks, resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.external_id"))
	}

	return checkAggr(
		[]string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"project_id":   projectID,
			"disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.node_count": fmt.Sprintf("%d", nodeCountElectable),
			"name": name},
		additionalChecks...,
	)
}

func configIncorrectTypeGobalClusterSelfManagedSharding(projectID, name string) string {
	return fmt.Sprintf(`
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
	`, projectID, name)
}

func configReplicaSetMultiCloud(orgID, projectName, name string) string {
	return fmt.Sprintf(`
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

				dynamic "region_configs" {
					for_each = [
						"US_EAST_4",
						"NORTH_AMERICA_NORTHEAST_1"
					]

					content {
						provider_name = "GCP"
						priority      = 0
						region_name   = region_configs.value

						read_only_specs {
							instance_size = "M10"
							node_count    = 2
						}
					}
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
	`, orgID, projectName, name)
}

func checkReplicaSetMultiCloud(name string, regionConfigs int) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.0.replication_specs.0.region_configs.#", acc.JSONEquals(strconv.Itoa(regionConfigs))),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
		resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.external_id"),
	}
	return checkAggr(
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id"},
		map[string]string{
			"name": name},
		additionalChecks...,
	)
}

func configShardedOldSchemaMultiCloud(orgID, projectName, name string, numShards int, analyticsSize string, configServerManagementMode *string) string {
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
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, numShards, analyticsSize, rootConfig)
}

func checkShardedOldSchemaMultiCloud(name string, numShards int, analyticsSize string, verifyExternalID bool, configServerManagementMode *string) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.analytics_specs.0.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.1.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
	}

	if verifyExternalID {
		additionalChecks = append(
			additionalChecks,
			resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.external_id"))
	}
	if configServerManagementMode != nil {
		additionalChecks = append(
			additionalChecks,
			resource.TestCheckResourceAttr(resourceName, "config_server_management_mode", *configServerManagementMode),
			resource.TestCheckResourceAttrSet(resourceName, "config_server_type"),
			resource.TestCheckResourceAttr(dataSourceName, "config_server_management_mode", *configServerManagementMode),
			resource.TestCheckResourceAttrSet(dataSourceName, "config_server_type"),
			resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.config_server_management_mode", *configServerManagementMode),
			resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.config_server_type"),
		)
	}

	return checkAggr(
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":                           name,
			"replication_specs.0.num_shards": strconv.Itoa(numShards),
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": analyticsSize,
		},
		additionalChecks...)
}

func configSingleProviderPaused(projectID, clusterName string, paused bool, instanceSize string) string {
	return fmt.Sprintf(`
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
	`, projectID, clusterName, paused, instanceSize)
}

func checkSingleProviderPaused(name string, paused bool) resource.TestCheckFunc {
	return checkAggr(
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":   name,
			"paused": strconv.FormatBool(paused)})
}

func configAdvanced(projectID, clusterName string, p *admin20240530.ClusterDescriptionProcessArgs, changeStreamOptions *int, mongodbMajorVersion *string) string {
	changeStreamOptionsString := ""
	mongodbMajorVersionString := ""
	if changeStreamOptions != nil {
		changeStreamOptionsString = fmt.Sprintf(`change_stream_options_pre_and_post_images_expire_after_seconds = %[1]d`, *changeStreamOptions)
	}
	if mongodbMajorVersion != nil {
		mongodbMajorVersionString = fmt.Sprintf(`mongo_db_major_version = %[1]q`, *mongodbMajorVersion)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id             = %[1]q
			name                   = %[2]q
			cluster_type           = "REPLICASET"
			%[12]s

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
		p.GetFailIndexKeyTooLong(), p.GetJavascriptEnabled(), p.GetMinimumEnabledTlsProtocol(), p.GetNoTableScan(),
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetTransactionLifetimeLimitSeconds(), changeStreamOptionsString, mongodbMajorVersionString)
}

func checkAdvanced(name, tls, changeStreamOptions string) resource.TestCheckFunc {
	return checkAggr(
		[]string{"project_id", "replication_specs.#", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name": name,
			"advanced_configuration.0.minimum_enabled_tls_protocol":                                   tls,
			"advanced_configuration.0.fail_index_key_too_long":                                        "false",
			"advanced_configuration.0.javascript_enabled":                                             "true",
			"advanced_configuration.0.no_table_scan":                                                  "false",
			"advanced_configuration.0.oplog_size_mb":                                                  "1000",
			"advanced_configuration.0.sample_refresh_interval_bi_connector":                           "310",
			"advanced_configuration.0.sample_size_bi_connector":                                       "110",
			"advanced_configuration.0.transaction_lifetime_limit_seconds":                             "300",
			"advanced_configuration.0.change_stream_options_pre_and_post_images_expire_after_seconds": changeStreamOptions},
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"))
}

func configAdvancedDefaultWrite(projectID, clusterName string, p *admin20240530.ClusterDescriptionProcessArgs) string {
	return fmt.Sprintf(`
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
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetDefaultReadConcern(), p.GetDefaultWriteConcern())
}

func checkAdvancedDefaultWrite(name, writeConcern, tls string) resource.TestCheckFunc {
	return checkAggr(
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
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
		resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"))
}

func configReplicationSpecsAutoScaling(projectID, clusterName string, p *admin.AdvancedAutoScalingSettings) string {
	return fmt.Sprintf(`
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
	`, projectID, clusterName, p.Compute.GetEnabled(), p.DiskGB.GetEnabled(), p.Compute.GetMaxInstanceSize())
}

func configReplicationSpecsAnalyticsAutoScaling(projectID, clusterName string, p *admin.AdvancedAutoScalingSettings) string {
	return fmt.Sprintf(`
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
	`, projectID, clusterName, p.Compute.GetEnabled(), p.DiskGB.GetEnabled(), p.Compute.GetMaxInstanceSize())
}

func configGeoShardedOldSchema(orgID, projectName, name string, numShardsFirstZone, numShardsSecondZone int, selfManagedSharding bool) string {
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, numShardsFirstZone, numShardsSecondZone, selfManagedSharding)
}

func checkGeoShardedOldSchema(name string, numShardsFirstZone, numShardsSecondZone int, isLatestProviderVersion, verifyExternalID bool) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{}

	if verifyExternalID {
		additionalChecks = append(additionalChecks, resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.external_id"))
	}

	if isLatestProviderVersion { // checks that will not apply if doing migration test with older version
		additionalChecks = append(additionalChecks, checkAggr(
			[]string{"replication_specs.0.zone_id", "replication_specs.0.zone_id"},
			map[string]string{
				"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": "60",
				"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": "60",
			}))
	}

	return checkAggr(
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

func configShardedOldSchemaDiskSizeGBElectableLevel(orgID, projectName, name string, diskSizeGB int) string {
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, diskSizeGB)
}

func checkShardedOldSchemaDiskSizeGBElectableLevel(diskSizeGB int) resource.TestCheckFunc {
	return checkAggr(
		[]string{},
		map[string]string{
			"replication_specs.0.num_shards": "2",
			"disk_size_gb":                   fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.analytics_specs.0.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
		})
}

func configShardedNewSchema(orgID, projectName, name string, diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIOPS, lastDiskIOPS *int, includeMiddleSpec bool) string {
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
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, firstInstanceSize, lastInstanceSize, firstDiskIOPSAttrs, lastDiskIOPSAttrs, thirdReplicationSpec, diskSizeGB)
}

func checkShardedNewSchema(diskSizeGB int, firstInstanceSize, lastInstanceSize string, firstDiskIops, lastDiskIops *int, isAsymmetricCluster, includeMiddleSpec bool) resource.TestCheckFunc {
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
	additionalChecks := acc.AddAttrSetChecks(dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replication_specs.#", "results.0.replication_specs.0.region_configs.#", "results.0.name", "results.0.termination_protection_enabled", "results.0.global_cluster_self_managed_sharding"}...)
	additionalChecks = acc.AddAttrChecksPrefix(dataSourcePluralName, additionalChecks, clusterChecks, "results.0")

	// expected id attribute only if cluster is symmetric
	if isAsymmetricCluster {
		additionalChecks = append(additionalChecks, checkAggr([]string{}, map[string]string{
			"replication_specs.0.id": "",
			"replication_specs.1.id": "",
		}))
		additionalChecks = acc.AddAttrChecks(dataSourcePluralName, additionalChecks, map[string]string{
			"results.0.replication_specs.0.id": "",
			"results.0.replication_specs.1.id": "",
		})
	} else {
		additionalChecks = append(additionalChecks, checkAggr([]string{"replication_specs.0.id", "replication_specs.1.id"}, map[string]string{}))
		additionalChecks = acc.AddAttrSetChecks(dataSourcePluralName, additionalChecks, "results.0.replication_specs.0.id", "results.0.replication_specs.1.id")
	}

	return checkAggr(
		[]string{"replication_specs.0.external_id", "replication_specs.0.zone_id", "replication_specs.1.external_id", "replication_specs.1.zone_id"},
		clusterChecks,
		additionalChecks...,
	)
}

func configGeoShardedNewSchema(orgID, projectName, name string, includeThirdShardInFirstZone bool) string {
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
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, thirdReplicationSpec)
}

func checkGeoShardedNewSchema(includeThirdShardInFirstZone bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if includeThirdShardInFirstZone {
		amtOfReplicationSpecs = 3
	} else {
		amtOfReplicationSpecs = 2
	}
	clusterChecks := map[string]string{
		"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
	}

	return checkAggr(
		[]string{},
		clusterChecks,
	)
}

func configShardedTransitionOldToNewSchema(orgID, projectName, name string, useNewSchema bool) string {
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

	return fmt.Sprintf(`
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
	`, orgID, projectName, name, replicationSpecs, dataSourceFlag)
}

func checkShardedTransitionOldToNewSchema(useNewSchema bool) resource.TestCheckFunc {
	var amtOfReplicationSpecs int
	if useNewSchema {
		amtOfReplicationSpecs = 2
	} else {
		amtOfReplicationSpecs = 1
	}
	var checksForNewSchema []resource.TestCheckFunc
	if useNewSchema {
		checksForNewSchema = []resource.TestCheckFunc{
			checkAggr([]string{"replication_specs.1.id", "replication_specs.0.external_id", "replication_specs.1.external_id"},
				map[string]string{
					"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
					"replication_specs.1.region_configs.0.electable_specs.0.instance_size": "M10",
					"replication_specs.1.region_configs.0.analytics_specs.0.instance_size": "M10",
				}),
		}
	}

	return checkAggr(
		[]string{"replication_specs.0.id"},
		map[string]string{
			"replication_specs.#": fmt.Sprintf("%d", amtOfReplicationSpecs),
			"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
			"replication_specs.0.region_configs.0.analytics_specs.0.instance_size": "M10",
		},
		checksForNewSchema...,
	)
}

func configGeoShardedTransitionOldToNewSchema(orgID, projectName, name string, useNewSchema bool) string {
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

	return fmt.Sprintf(`
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
	`, orgID, projectName, name, replicationSpecs, dataSourceFlag)
}

func checkGeoShardedTransitionOldToNewSchema(useNewSchema bool) resource.TestCheckFunc {
	if useNewSchema {
		return checkAggr(
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
	return checkAggr(
		[]string{"replication_specs.0.id", "replication_specs.1.id"},
		map[string]string{
			"replication_specs.#":           "2",
			"replication_specs.0.zone_name": "zone 1",
			"replication_specs.1.zone_name": "zone 2",
		},
	)
}

func configReplicaSetScalingStrategyAndRedactClientLogData(orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData)
}

func configReplicaSetScalingStrategyAndRedactClientLogDataOldSchema(orgID, projectName, name, replicaSetScalingStrategy string, redactClientLogData bool) string {
	return fmt.Sprintf(`
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
	`, orgID, projectName, name, replicaSetScalingStrategy, redactClientLogData)
}

func checkReplicaSetScalingStrategyAndRedactClientLogData(replicaSetScalingStrategy string, redactClientLogData bool) resource.TestCheckFunc {
	clusterChecks := map[string]string{
		"replica_set_scaling_strategy": replicaSetScalingStrategy,
		"redact_client_log_data":       strconv.FormatBool(redactClientLogData),
	}

	// plural data source checks
	additionalChecks := acc.AddAttrSetChecks(dataSourcePluralName, nil,
		[]string{"results.#", "results.0.replica_set_scaling_strategy", "results.0.redact_client_log_data"}...)
	return checkAggr(
		[]string{},
		clusterChecks,
		additionalChecks...,
	)
}

func configPriority(orgID, projectName, clusterName string, oldSchema, swapPriorities bool) string {
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

	return fmt.Sprintf(`
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
	`, orgID, projectName, clusterName, strType, strNumShards, strConfigs)
}
