package advancedcluster_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115008/admin"
)

const (
	resourceName         = "mongodbatlas_advanced_cluster.test"
	dataSourceName       = "data.mongodbatlas_advanced_cluster.test"
	dataSourcePluralName = "data.mongodbatlas_advanced_clusters.test"
)

func TestAccClusterAdvancedCluster_basicTenant(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "termination_protection_enabled"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterName),
					resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.termination_protection_enabled"),
				),
			},
			{
				Config: configTenant(projectID, clusterNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.termination_protection_enabled"),
				),
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

func TestAccClusterAdvancedCluster_singleProvider(t *testing.T) {
	var (
		dataSourceName = fmt.Sprintf("data.%s", resourceName)
		projectID      = acc.ProjectIDExecution(t)
		clusterName    = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProvider(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrWith(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.0.disk_iops", acc.IntGreatThan(0)),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_multicloud(t *testing.T) {
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMultiCloud(orgID, projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterName),
				),
			},
			{
				Config: configMultiCloud(orgID, projectName, clusterNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterNameUpdated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"replication_specs", "retain_backups_enabled"},
			},
		},
	})
}

func TestAccClusterAdvancedCluster_multicloudSharded(t *testing.T) {
	var (
		orgID              = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName        = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMultiCloudSharded(orgID, projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
			},
			{
				Config: configMultiCloudSharded(orgID, projectName, clusterNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
				),
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

func TestAccClusterAdvancedCluster_unpausedToPaused(t *testing.T) {
	var (
		projectID           = acc.ProjectIDExecution(t)
		clusterName         = acc.RandomClusterName()
		instanceSize        = "M10"
		anotherInstanceSize = "M20"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(projectID, clusterName, false, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
			{
				Config: configSingleProviderPaused(projectID, clusterName, true, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configSingleProviderPaused(projectID, clusterName, true, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				Config: configSingleProviderPaused(projectID, clusterName, false, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
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
		processArgs        = &admin.ClusterDescriptionProcessArgs{
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
		processArgsUpdated = &admin.ClusterDescriptionProcessArgs{
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvanced(projectID, clusterName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
				),
			},
			{
				Config: configAdvanced(projectID, clusterNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.name"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedCluster_defaultWrite(t *testing.T) {
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
		processArgs        = &admin.ClusterDescriptionProcessArgs{
			DefaultReadConcern:               conversion.StringPtr("available"),
			DefaultWriteConcern:              conversion.StringPtr("1"),
			JavascriptEnabled:                conversion.Pointer(true),
			MinimumEnabledTlsProtocol:        conversion.StringPtr("TLS1_1"),
			NoTableScan:                      conversion.Pointer(false),
			OplogSizeMB:                      conversion.Pointer(1000),
			SampleRefreshIntervalBIConnector: conversion.Pointer(310),
			SampleSizeBIConnector:            conversion.Pointer(110),
		}
		processArgsUpdated = &admin.ClusterDescriptionProcessArgs{
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedDefaultWrite(projectID, clusterName, processArgs),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
			{
				Config: configAdvancedDefaultWrite(projectID, clusterNameUpdated, processArgsUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "majority"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAutoScaling(projectID, clusterName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.0.compute_enabled", "false"),
				),
			},
			{
				Config: configReplicationSpecsAutoScaling(projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
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
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(projectID, clusterName, autoScaling),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "false"),
				),
			},
			{
				Config: configReplicationSpecsAnalyticsAutoScaling(projectID, clusterNameUpdated, autoScalingUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_auto_scaling.0.compute_enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterAdvancedClusterConfig_replicationSpecsAndShardUpdating(t *testing.T) {
	var (
		orgID            = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName      = acc.RandomProjectName()
		clusterName      = acc.RandomClusterName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		numShards        = "1"
		numShardsUpdated = "2"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMultiZoneWithShards(orgID, projectName, clusterName, numShards, numShards),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.num_shards", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.1.num_shards", "1"),
				),
			},
			{
				Config: configMultiZoneWithShards(orgID, projectName, clusterName, numShardsUpdated, numShards),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.num_shards", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.1.num_shards", "1"),
				),
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
				Config: configWithTags(orgID, projectName, clusterName, nil),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, []admin.ResourceTag{
					{
						Key:   "key 1",
						Value: "value 1",
					},
					{
						Key:   "key 2",
						Value: "value 2",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, []admin.ResourceTag{
					{
						Key:   "key 3",
						Value: "value 3",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.0.tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourcePluralName, "results.0.tags.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
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
		if _, _, err := acc.ConnV2().ClustersApi.GetCluster(context.Background(), ids["project_id"], ids["cluster_name"]).Execute(); err == nil {
			return nil
		}
		return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
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

func configWithTags(orgID, projectName, name string, tags []admin.ResourceTag) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, label.GetKey(), label.GetValue())
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

func configSingleProvider(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			retain_backups_enabled = "true"

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

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name 	     = mongodbatlas_advanced_cluster.test.name
		}
	`, projectID, name)
}

func configMultiCloud(orgID, projectName, name string) string {
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
				region_configs {
					electable_specs {
						instance_size = "M10"
						node_count    = 2
					}
					provider_name = "GCP"
					priority      = 6
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
	`, orgID, projectName, name)
}

func configMultiCloudSharded(orgID, projectName, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}	

		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			cluster_type = "SHARDED"

			replication_specs {
				num_shards = 1
				region_configs {
					electable_specs {
						instance_size = "M30"
						node_count    = 3
					}
					analytics_specs {
						instance_size = "M30"
						node_count    = 1
					}
					provider_name = "AWS"
					priority      = 7
					region_name   = "EU_WEST_1"
				}
				region_configs {
					electable_specs {
						instance_size = "M30"
						node_count    = 2
					}
					provider_name = "AZURE"
					priority      = 6
					region_name   = "US_EAST_2"
				}
			}
		}
	`, orgID, projectName, name)
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
	`, projectID, clusterName, paused, instanceSize)
}

func configAdvanced(projectID, clusterName string, p *admin.ClusterDescriptionProcessArgs) string {
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
				fail_index_key_too_long              = %[3]t
				javascript_enabled                   = %[4]t
				minimum_enabled_tls_protocol         = %[5]q
				no_table_scan                        = %[6]t
				oplog_size_mb                        = %[7]d
				sample_size_bi_connector			 = %[8]d
				sample_refresh_interval_bi_connector = %[9]d
			transaction_lifetime_limit_seconds   = %[10]d
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
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetTransactionLifetimeLimitSeconds())
}

func configAdvancedDefaultWrite(projectID, clusterName string, p *admin.ClusterDescriptionProcessArgs) string {
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
	`, projectID, clusterName, p.GetJavascriptEnabled(), p.GetMinimumEnabledTlsProtocol(), p.GetNoTableScan(),
		p.GetOplogSizeMB(), p.GetSampleSizeBIConnector(), p.GetSampleRefreshIntervalBIConnector(), p.GetDefaultReadConcern(), p.GetDefaultWriteConcern())
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

func configMultiZoneWithShards(orgID, projectName, name, numShardsFirstZone, numShardsSecondZone string) string {
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
				num_shards = %[4]q

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
				num_shards = %[5]q

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
	`, orgID, projectName, name, numShardsFirstZone, numShardsSecondZone)
}
