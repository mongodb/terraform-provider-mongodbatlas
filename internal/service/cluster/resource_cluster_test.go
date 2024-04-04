package cluster_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	clustersvc "github.com/mongodb/terraform-provider-mongodbatlas/internal/service/cluster"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	resourceName   = "mongodbatlas_cluster.test"
	dataSourceName = "data.mongodbatlas_cluster.test"
)

func TestAccClusterRSCluster_basicAWS_simple(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWS(projectID, clusterName, true, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_release_system", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "accept_data_risks_and_force_replica_set_reconfig", ""),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_disk_gb_enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_backup_policy.#"),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_backup_policy.0.policies.#"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_strings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_strings.0.private_endpoint.#"),
				),
			},
			{
				Config: configAWS(projectID, clusterName, false, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "pit_enabled", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "retain_backups_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "version_release_system", "LTS"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_disk_gb_enabled", "false"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_backup", "retain_backups_enabled"},
			},
		},
	})
}

func TestAccClusterRSCluster_basicAWS_instanceScale(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSNVMEInstance(projectID, clusterName, "M40_NVME"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_instance_size_name", "M40_NVME"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
			{
				Config: configAWSNVMEInstance(projectID, clusterName, "M50_NVME"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_instance_size_name", "M50_NVME"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basic_Partial_AdvancedConf(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedConf(projectID, clusterName, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              conversion.Pointer(false),
					JavascriptEnabled:                conversion.Pointer(true),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      conversion.Pointer(false),
					OplogSizeMB:                      conversion.Pointer[int64](1000),
					SampleRefreshIntervalBIConnector: conversion.Pointer[int64](310),
					SampleSizeBIConnector:            conversion.Pointer[int64](110),
					TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterName),
					resource.TestCheckResourceAttr(dataSourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(dataSourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.#"),
					resource.TestCheckResourceAttr(dataSourceName, "version_release_system", "LTS"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(dataSourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
				),
			},
			{
				Config: configAdvancedConfPartial(projectID, clusterName, "false", &matlas.ProcessArgs{
					MinimumEnabledTLSProtocol: "TLS1_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
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

func TestAccClusterRSCluster_basic_DefaultWriteRead_AdvancedConf(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedConfDefaultWriteRead(projectID, clusterName, "false", &matlas.ProcessArgs{
					DefaultReadConcern:               "available",
					DefaultWriteConcern:              "1",
					FailIndexKeyTooLong:              conversion.Pointer(false),
					JavascriptEnabled:                conversion.Pointer(true),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      conversion.Pointer(false),
					OplogSizeMB:                      conversion.Pointer[int64](1000),
					SampleRefreshIntervalBIConnector: conversion.Pointer[int64](310),
					SampleSizeBIConnector:            conversion.Pointer[int64](110),
					TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
				),
			},
			{
				Config: configAdvancedConfPartialDefault(projectID, clusterName, "false", &matlas.ProcessArgs{
					MinimumEnabledTLSProtocol: "TLS1_2",
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_read_concern", "available"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.default_write_concern", "1"),
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

func TestAccClusterRSCluster_emptyAdvancedConf(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedConfPartial(projectID, clusterName, "false", &matlas.ProcessArgs{
					MinimumEnabledTLSProtocol: "TLS1_2",
				}),
			},
			{
				Config: configAdvancedConf(projectID, clusterName, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              conversion.Pointer(false),
					JavascriptEnabled:                conversion.Pointer(true),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      conversion.Pointer(false),
					OplogSizeMB:                      conversion.Pointer[int64](1000),
					SampleRefreshIntervalBIConnector: conversion.Pointer[int64](310),
					SampleSizeBIConnector:            conversion.Pointer[int64](110),
					TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicAdvancedConf(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAdvancedConf(projectID, clusterName, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              conversion.Pointer(false),
					JavascriptEnabled:                conversion.Pointer(true),
					MinimumEnabledTLSProtocol:        "TLS1_2",
					NoTableScan:                      conversion.Pointer(true),
					OplogSizeMB:                      conversion.Pointer[int64](1000),
					SampleRefreshIntervalBIConnector: conversion.Pointer[int64](310),
					SampleSizeBIConnector:            conversion.Pointer[int64](110),
					TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](300),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_2"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "true"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "1000"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "110"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "310"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "300"),
				),
			},
			{
				Config: configAdvancedConf(projectID, clusterName, "false", &matlas.ProcessArgs{
					FailIndexKeyTooLong:              conversion.Pointer(false),
					JavascriptEnabled:                conversion.Pointer(false),
					MinimumEnabledTLSProtocol:        "TLS1_1",
					NoTableScan:                      conversion.Pointer(false),
					OplogSizeMB:                      conversion.Pointer[int64](990),
					SampleRefreshIntervalBIConnector: conversion.Pointer[int64](0),
					SampleSizeBIConnector:            conversion.Pointer[int64](0),
					TransactionLifetimeLimitSeconds:  conversion.Pointer[int64](60),
				}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.fail_index_key_too_long", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.javascript_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.minimum_enabled_tls_protocol", "TLS1_1"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.no_table_scan", "false"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.oplog_size_mb", "990"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_size_bi_connector", "0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.sample_refresh_interval_bi_connector", "0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.0.transaction_lifetime_limit_seconds", "60"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicAzure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.basic_azure"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		clusterName  = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAzure(orgID, projectName, clusterName, "true", "M30", true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: configAzure(orgID, projectName, clusterName, "false", "M30", true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_AzureUpdateToNVME(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.basic_azure"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		clusterName  = acc.RandomClusterName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAzure(orgID, projectName, clusterName, "true", "M60", true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_instance_size_name", "M60"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: configAzure(orgID, projectName, clusterName, "true", "M60_NVME", false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_instance_size_name", "M60_NVME"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicGCP(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.basic_gcp"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		clusterName  = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCP(orgID, projectName, clusterName, "true"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
			{
				Config: configGCP(orgID, projectName, clusterName, "false"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_WithBiConnectorGCP(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.basic_gcp"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		clusterName  = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCPWithBiConnector(orgID, projectName, clusterName, "true", false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "false"),
				),
			},
			{
				Config: configGCPWithBiConnector(orgID, projectName, clusterName, "false", true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "40"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "bi_connector_config.0.enabled", "true"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_MultiRegion(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.multi_region"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName  = acc.RandomClusterName()
	)

	createRegionsConfig := `regions_config {
					region_name     = "US_EAST_1"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}`

	updatedRegionsConfig := `regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 6
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_WEST_1"
					electable_nodes = 1
					priority        = 5
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_EAST_1"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configMultiRegion(orgID, projectName, clusterName, "true", createRegionsConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
				),
			},
			{
				Config: configMultiRegion(orgID, projectName, clusterName, "false", updatedRegionsConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_ProviderRegionName(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.multi_region"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName  = acc.RandomClusterName()
	)

	updatedRegionsConfig := `regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 6
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_WEST_1"
					electable_nodes = 1
					priority        = 5
					read_only_nodes = 0
				}
				regions_config {
					region_name     = "US_EAST_1"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config:      configMultiRegionWithProviderRegionNameInvalid(orgID, projectName, clusterName, "false", updatedRegionsConfig),
				ExpectError: regexp.MustCompile("attribute must be set ONLY for single-region clusters"),
			},
			{
				Config: configSingleRegionWithProviderRegionName(orgID, projectName, clusterName, "false"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
				),
			},
			{
				Config: configMultiRegion(orgID, projectName, clusterName, "false", updatedRegionsConfig),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "REPLICASET"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "3"),
				),
			},
			{
				Config: configMultiRegion(orgID, projectName, clusterName, "false", updatedRegionsConfig),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestAccClusterRSCluster_Global(t *testing.T) {
	var (
		resourceSuffix = "global_cluster"
		resourceName   = fmt.Sprintf("mongodbatlas_cluster.%s", resourceSuffix)
		projectID      = acc.ProjectIDExecution(t)
		clusterName    = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: acc.ConfigClusterGlobal(resourceSuffix, projectID, clusterName, "false"),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.1.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "80"),
					resource.TestCheckResourceAttr(resourceName, "cluster_type", "GEOSHARDED"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.1.regions_config.#", "1"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_AWSWithLabels(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.aws_with_labels"
		projectID    = acc.ProjectIDExecution(t)
		clusterName  = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, clusterName, "false", "M10", "US_WEST_2", []matlas.Label{}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, clusterName, "false", "M10", "US_WEST_2",
					[]matlas.Label{
						{
							Key:   "key 4",
							Value: "value 4",
						},
						{
							Key:   "key 3",
							Value: "value 3",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "3"),
				),
			},
			{
				Config: testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, clusterName, "false", "M10", "US_WEST_2",
					[]matlas.Label{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 5",
							Value: "value 5",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "2"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_WithTags(t *testing.T) {
	var (
		dataSourceClustersName = "data.mongodbatlas_clusters.test"
		orgID                  = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName            = acc.RandomProjectName() // No ProjectIDExecution because this test has plural datasource
		clusterName            = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithTags(orgID, projectName, clusterName, "false", "M10", "US_WEST_2", []matlas.Tag{}),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "0"),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "0"),
				),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, "false", "M10", "US_WEST_2",
					[]matlas.Tag{
						{
							Key:   "key 1",
							Value: "value 1",
						},
						{
							Key:   "key 2",
							Value: "value 2",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap2),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap1),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap2),
				),
			},
			{
				Config: configWithTags(orgID, projectName, clusterName, "false", "M10", "US_WEST_2",
					[]matlas.Tag{
						{
							Key:   "key 3",
							Value: "value 3",
						},
					},
				),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttr(resourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceName, "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceName, "tags.*", acc.ClusterTagsMap3),
					resource.TestCheckResourceAttr(dataSourceClustersName, "results.0.tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(dataSourceClustersName, "results.0.tags.*", acc.ClusterTagsMap3),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withPrivateEndpointLink(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	var (
		resourceName = "mongodbatlas_cluster.with_endpoint_link"

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		region       = os.Getenv("AWS_REGION")
		providerName = "AWS"

		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
		clusterName     = acc.RandomClusterName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t); acc.PreCheckPeeringEnvAWS(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configWithPrivateEndpointLink(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withAzureNetworkPeering(t *testing.T) {
	acc.SkipTestForCI(t) // needs Azure configuration

	var (
		resourceName = "mongodbatlas_cluster.with_azure_peering"

		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		directoryID       = os.Getenv("AZURE_DIRECTORY_ID")
		subcrptionID      = os.Getenv("AZURE_SUBSCRIPTION_ID")
		resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		vNetName          = os.Getenv("AZURE_VNET_NAME")
		providerName      = "AZURE"
		region            = os.Getenv("AZURE_REGION")

		atlasCidrBlock = "192.168.208.0/21"
		clusterName    = acc.RandomClusterName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAzureWithNetworkPeering(projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withGCPNetworkPeering(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		projectID        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		gcpRegion        = os.Getenv("GCP_REGION_NAME")
		gcpProjectID     = os.Getenv("GCP_PROJECT_ID")
		providerName     = "GCP"
		gcpPeeringName   = acc.RandomName()
		clusterName      = acc.RandomClusterName()
		gcpClusterRegion = os.Getenv("GCP_CLUSTER_REGION_NAME")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPeeringEnvGCP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCPWithNetworkPeering(gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withAzureAndContainerID(t *testing.T) {
	acc.SkipTestForCI(t) // needs Azure configuration

	var (
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName       = acc.RandomClusterName()
		providerName      = "AZURE"
		region            = os.Getenv("AZURE_REGION")
		directoryID       = os.Getenv("AZURE_DIRECTORY_ID")
		subcrptionID      = os.Getenv("AZURE_SUBSCRIPTION_ID")
		resourceGroupName = os.Getenv("AZURE_RESOURCE_GROUP_NAME")
		vNetName          = os.Getenv("AZURE_VNET_NAME")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPeeringEnvAzure(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAzureWithContainerID(projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withAWSAndContainerID(t *testing.T) {
	acc.SkipTestForCI(t) // needs AWS configuration

	var (
		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = acc.RandomClusterName()
		providerName = "AWS"
		awsRegion    = os.Getenv("AWS_REGION")
		vpcCIDRBlock = os.Getenv("AWS_VPC_CIDR_BLOCK")
		awsAccountID = os.Getenv("AWS_ACCOUNT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSWithContainerID(awsAccessKey, awsSecretKey, projectID, clusterName, providerName, awsRegion, vpcCIDRBlock, awsAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "container_id"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withGCPAndContainerID(t *testing.T) {
	acc.SkipTestForCI(t) // needs GCP configuration

	var (
		gcpProjectID     = os.Getenv("GCP_PROJECT_ID")
		gcpRegion        = os.Getenv("GCP_REGION_NAME")
		projectID        = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName      = acc.RandomClusterName()
		providerName     = "GCP"
		gcpClusterRegion = os.Getenv("GCP_CLUSTER_REGION_NAME")
		gcpPeeringName   = acc.RandomName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckPeeringEnvGCP(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCPWithContainerID(gcpProjectID, gcpRegion, projectID, clusterName, providerName, gcpClusterRegion, gcpPeeringName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_withAutoScalingAWS(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()

		instanceSize = "M30"
		minSize      = ""
		maxSize      = "M60"

		instanceSizeUpdated = "M60"
		minSizeUpdated      = "M20"
		maxSizeUpdated      = "M80"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSWithAutoscaling(projectID, clusterName, "true", "false", "true", "false", minSize, maxSize, instanceSize),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "provider_auto_scaling_compute_max_instance_size", maxSize),
					resource.TestCheckResourceAttr(dataSourceName, "name", clusterName),
					resource.TestCheckResourceAttr(dataSourceName, "auto_scaling_compute_enabled", "true"),
					resource.TestCheckResourceAttr(dataSourceName, "provider_auto_scaling_compute_max_instance_size", maxSize),
					resource.TestCheckResourceAttrSet(dataSourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(dataSourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttr(dataSourceName, "version_release_system", "LTS"),
					resource.TestCheckResourceAttr(dataSourceName, "termination_protection_enabled", "false"),
				),
			},
			{
				Config: configAWSWithAutoscaling(projectID, clusterName, "false", "true", "true", "true", minSizeUpdated, maxSizeUpdated, instanceSizeUpdated),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "auto_scaling_compute_scale_down_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "provider_auto_scaling_compute_min_instance_size", minSizeUpdated),
					resource.TestCheckResourceAttr(resourceName, "provider_auto_scaling_compute_max_instance_size", maxSizeUpdated),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_tenant(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cluster.tenant"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because tenant
		clusterName  = acc.RandomClusterName()
	)

	dbMajorVersion := testAccGetMongoDBAtlasMajorVersion()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(orgID, projectName, clusterName, "M2", "2", dbMajorVersion),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
			{
				Config: configTenantUpdated(orgID, projectName, clusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_tenant_m5(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_cluster.tenant"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because tenant
		clusterName    = acc.RandomClusterName()
		dbMajorVersion = testAccGetMongoDBAtlasMajorVersion()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configTenant(orgID, projectName, clusterName, "M5", "5", dbMajorVersion),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "5"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicGCPRegionNameWesternUS(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		clusterName = acc.RandomClusterName()
		regionName  = "WESTERN_US"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCPRegionName(orgID, projectName, clusterName, regionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_region_name", regionName),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicGCPRegionNameUSWest2(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		clusterName = acc.RandomClusterName() // No ProjectIDExecution to avoid cross-region limits because no AWS
		regionName  = "US_WEST_2"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configGCPRegionName(orgID, projectName, clusterName, regionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "provider_region_name", regionName),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_RegionsConfig(t *testing.T) {
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName() // No ProjectIDExecution to avoid cross-region limits because multi-region
		clusterName = acc.RandomClusterName()
	)

	replications := `replication_specs {
		num_shards = 1
		zone_name = "us2"
		regions_config{
			region_name     = "US_WEST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	  }
	 replication_specs {
		num_shards = 1
		zone_name = "us3"
		regions_config{
			region_name     = "US_EAST_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	 }
	 replication_specs {
		num_shards = 1
		zone_name = "us1"
		regions_config{
			region_name     = "US_WEST_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	}`

	replicationsUpdate := `replication_specs {
		num_shards = 1
		zone_name = "us2"
		regions_config{
			region_name     = "US_WEST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	  }

	 replication_specs {
		num_shards = 1
		zone_name = "us1"
		regions_config{
			region_name     = "US_WEST_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	}`

	replicationsShardsUpdate := `replication_specs {
		num_shards = 2
		zone_name = "us2"
		regions_config{
			region_name     = "US_WEST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	  }

	 replication_specs {
		num_shards = 1
		zone_name = "us1"
		regions_config{
			region_name     = "US_WEST_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
		}
	}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configRegions(orgID, projectName, clusterName, replications),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "3"),
				),
			},
			{
				Config: configRegions(orgID, projectName, clusterName, replicationsUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
				),
			},
			{
				Config: configRegions(orgID, projectName, clusterName, replicationsShardsUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
					// Note: replication_specs is a set for the cluster resource, therefore the order will not be consistent
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_specs.*", map[string]string{"num_shards": "1"}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "replication_specs.*", map[string]string{"num_shards": "2"}),
				),
			},
		},
	})
}

func TestAccClusterRSCluster_basicAWS_UnpauseToPaused(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSPaused(projectID, clusterName, true, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
			{
				Config: configAWSPaused(projectID, clusterName, false, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       acc.ImportStateClusterIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cloud_backup", "backup_enabled"},
			},
		},
	})
}

func TestAccClusterRSCluster_basicAWS_PausedToUnpaused(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configAWSPaused(projectID, clusterName, true, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "true"),
				),
			},
			{
				Config: configAWSPaused(projectID, clusterName, false, false),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttr(resourceName, "paused", "false"),
				),
			},
		},
	})
}

func testAccGetMongoDBAtlasMajorVersion() string {
	conn, _ := matlas.New(http.DefaultClient, matlas.SetBaseURL(matlas.CloudURL))
	majorVersion, _, _ := conn.DefaultMongoDBMajorVersion.Get(context.Background())

	return majorVersion
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
		log.Printf("[DEBUG] projectID: %s, name %s", ids["project_id"], ids["cluster_name"])
		if _, _, err := acc.Conn().Clusters.Get(context.Background(), ids["project_id"], ids["cluster_name"]); err == nil {
			return nil
		}
		return fmt.Errorf("cluster(%s:%s) does not exist", rs.Primary.Attributes["project_id"], rs.Primary.ID)
	}
}

func configAWS(projectID, name string, backupEnabled, autoDiskGBEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id                   = %[1]q
			name                         = %[2]q
			disk_size_gb                 = 100
			cluster_type = "REPLICASET"
			replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_WEST_2"
			     electable_nodes = 3
			     priority        = 7
							read_only_nodes = 0
		       }
		    }
			cloud_backup                 = %[3]t
			pit_enabled                  = %[3]t
			retain_backups_enabled       = true
			auto_scaling_disk_gb_enabled = %[4]t
			provider_name               = "AWS"
			provider_instance_size_name = "M30"
		}
	`, projectID, name, backupEnabled, autoDiskGBEnabled)
}

func configAWSNVMEInstance(projectID, name, instanceName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cloud_backup                 = true
			provider_region_name     = "US_WEST_2"
			provider_name               = "AWS"
			provider_instance_size_name = %[3]q
			provider_volume_type        = "PROVISIONED"
		}
	`, projectID, name, instanceName)
}

func configAdvancedConf(projectID, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_WEST_2"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[3]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"

			advanced_configuration  {
				fail_index_key_too_long              = %[4]t
				javascript_enabled                   = %[5]t
				minimum_enabled_tls_protocol         = %[6]q
				no_table_scan                        = %[7]t
				oplog_size_mb                        = %[8]d
				sample_size_bi_connector			 = %[9]d
				sample_refresh_interval_bi_connector = %[10]d
				transaction_lifetime_limit_seconds   = %[11]d
			}
		}

		data "mongodbatlas_cluster" "test" {
			project_id = mongodbatlas_cluster.test.project_id
			name 	     = mongodbatlas_cluster.test.name
		}
	`, projectID, name, autoscalingEnabled,
		*p.FailIndexKeyTooLong, *p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, *p.TransactionLifetimeLimitSeconds)
}

func configAdvancedConfDefaultWriteRead(projectID, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10
			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[3]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"

			advanced_configuration {
				javascript_enabled                   = %[4]t
				minimum_enabled_tls_protocol         = %[5]q
				no_table_scan                        = %[6]t
				oplog_size_mb                        = %[7]d
				sample_size_bi_connector             = %[8]d
				sample_refresh_interval_bi_connector = %[9]d
				default_read_concern                 = %[10]q
				default_write_concern                = %[11]q
			}
		}
	`, projectID, name, autoscalingEnabled,
		*p.JavascriptEnabled, p.MinimumEnabledTLSProtocol, *p.NoTableScan,
		*p.OplogSizeMB, *p.SampleSizeBIConnector, *p.SampleRefreshIntervalBIConnector, p.DefaultReadConcern, p.DefaultWriteConcern)
}

func configAdvancedConfPartial(projectID, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_WEST_2"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[3]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "US_WEST_2"

			advanced_configuration {
				minimum_enabled_tls_protocol         = %[4]q
			}
		}
	`, projectID, name, autoscalingEnabled, p.MinimumEnabledTLSProtocol)
}

func configAdvancedConfPartialDefault(projectID, name, autoscalingEnabled string, p *matlas.ProcessArgs) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10

			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}

			backup_enabled               = false
			auto_scaling_disk_gb_enabled =  %[3]s

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name        = "US_WEST_2"

			advanced_configuration {
				minimum_enabled_tls_protocol = %[4]q
			}
		}
	`, projectID, name, autoscalingEnabled, p.MinimumEnabledTLSProtocol)
}

func configAzure(orgID, projectName, name, backupEnabled, instanceSizeName string, includeDiskType bool) string {
	var diskType string
	if includeDiskType {
		diskType = `provider_disk_type_name     = "P6"`
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "basic_azure" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_EAST_2"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			cloud_backup                 = %[4]q
			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "AZURE"
			%[5]s
			provider_instance_size_name = %[6]q
			provider_region_name        = "US_EAST_2"
		}
	`, orgID, projectName, name, backupEnabled, diskType, instanceSizeName)
}

func configGCP(orgID, projectName, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "basic_gcp" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			disk_size_gb = 40

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_EAST_4"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			cloud_backup                 = %[4]q
			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "GCP"
			provider_instance_size_name = "M30"
		}
	`, orgID, projectName, name, backupEnabled)
}

func configGCPWithBiConnector(orgID, projectName, name, backupEnabled string, biConnectorEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "basic_gcp" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			disk_size_gb = 40

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "US_EAST_4"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			cloud_backup                 = %[4]q
			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "GCP"
			provider_instance_size_name = "M30"
			bi_connector_config {
				enabled = %[5]t
			}
		}
	`, orgID, projectName, name, backupEnabled, biConnectorEnabled)
}

func configMultiRegion(orgID, projectName, name, backupEnabled, regionsConfig string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "multi_region" {
			project_id              = mongodbatlas_project.cluster_project.id
			name                    = %[3]q
			disk_size_gb            = 100
			num_shards              = 1
			cloud_backup            = %[4]s
			cluster_type            = "REPLICASET"

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"

			replication_specs {
				num_shards = 1

				%[5]s
			}
		}
	`, orgID, projectName, name, backupEnabled, regionsConfig)
}

func configMultiRegionWithProviderRegionNameInvalid(orgID, projectName, name, backupEnabled, regionsConfig string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "multi_region" {
			project_id              = mongodbatlas_project.cluster_project.id
			name                    = %[3]q
			disk_size_gb            = 100
			num_shards              = 1
			cloud_backup            = %[4]s
			cluster_type            = "REPLICASET"

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name = "US_WEST_2"

			replication_specs {
				num_shards = 1

				%[5]s
			}
		}
	`, orgID, projectName, name, backupEnabled, regionsConfig)
}

func configSingleRegionWithProviderRegionName(orgID, projectName, name, backupEnabled string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "multi_region" {
			project_id              = mongodbatlas_project.cluster_project.id
			name                    = %[3]q
			disk_size_gb            = 100
			num_shards              = 1
			cloud_backup            = %[4]s
			cluster_type            = "REPLICASET"

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			provider_region_name = "US_WEST_2"

			replication_specs {
				num_shards = 1

				regions_config {
					region_name     = "US_WEST_2"
					electable_nodes = 3
					priority        = 7
					read_only_nodes = 0
				}
			}
		}
	`, orgID, projectName, name, backupEnabled)
}

func configTenant(orgID, projectName, name, instanceSize, diskSize, majorDBVersion string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "cluster_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cluster" "tenant" {
		project_id = mongodbatlas_project.cluster_project.id
		name       = %[3]q

		provider_name         = "TENANT"
		backing_provider_name = "AWS"
		provider_region_name  = "US_EAST_1"
	  	//M2 must be 2, M5 must be 5
	  	disk_size_gb            = %[4]q

		provider_instance_size_name  = %[5]q
		//These must be the following values
 	 	mongo_db_major_version = %[6]q
	  }
	`, orgID, projectName, name, diskSize, instanceSize, majorDBVersion)
}

func configTenantUpdated(orgID, projectName, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "cluster_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cluster" "tenant" {
		project_id = mongodbatlas_project.cluster_project.id
		name       = %[3]q

		provider_name        = "AWS"
		provider_region_name = "EU_CENTRAL_1"

		provider_instance_size_name  = "M10"
		disk_size_gb                 = 10
		auto_scaling_disk_gb_enabled = true
	  }
	`, orgID, projectName, name)
}

func testAccMongoDBAtlasClusterAWSConfigdWithLabels(projectID, name, backupEnabled, tier, region string, labels []matlas.Label) string {
	var labelsConf string
	for _, label := range labels {
		labelsConf += fmt.Sprintf(`
			labels {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "aws_with_labels" {
			project_id   = %[1]q
			name         = %[2]q
			disk_size_gb = 10
  
			backup_enabled               = %[3]s
			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = %[4]q
			cluster_type = "REPLICASET"
			  replication_specs {
				num_shards = 1
				regions_config {
				  region_name     = %[5]q
				  electable_nodes = 3
				  priority        = 7
				  read_only_nodes = 0
				}
		  	}
			%[6]s
		}
	`, projectID, name, backupEnabled, tier, region, labelsConf)
}

func configWithTags(orgID, projectName, name, backupEnabled, tier, region string, tags []matlas.Tag) string {
	var tagsConf string
	for _, label := range tags {
		tagsConf += fmt.Sprintf(`
			tags {
				key   = "%s"
				value = "%s"
			}
		`, label.Key, label.Value)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_project" "cluster_project" {
			org_id = %[1]q
			name   = %[2]q
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = mongodbatlas_project.cluster_project.id
			name         = %[3]q
			disk_size_gb = 10
  
			backup_enabled               = %[4]s
			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_instance_size_name = %[5]q
			cluster_type = "REPLICASET"
			replication_specs {
			num_shards = 1
			regions_config {
				region_name     = %[6]q
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
			}
		  	}
			%[7]s
		}

		data "mongodbatlas_cluster" "test" {
			project_id = mongodbatlas_cluster.test.project_id
			name 	     = mongodbatlas_cluster.test.name
		}
	
		data "mongodbatlas_clusters" "test" {
			project_id = mongodbatlas_cluster.test.project_id
		}

	`, orgID, projectName, name, backupEnabled, tier, region, tagsConf)
}

func configWithPrivateEndpointLink(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "${lower(replace("%[5]s", "_", "-"))}"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%[6]s"
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%[7]s"]
			security_group_ids = ["%[8]s"]
		}

		resource "mongodbatlas_privatelink_endpoint_service" "test" {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id       = mongodbatlas_privatelink_endpoint.test.private_link_id
			endpoint_service_id = aws_vpc_endpoint.ptfe_service.id
			provider_name = "%[4]s"
		}

		resource "mongodbatlas_cluster" "with_endpoint_link" {
		  project_id             = "%[3]s"
		  name                   = "%[9]s"
		  disk_size_gb           = 5

		  // Provider Settings "block"
		  provider_name               = "AWS"
		  provider_region_name        = "${upper(replace("%[5]s", "-", "_"))}"
		  provider_instance_size_name = "M10"
		  cloud_backup                = true // enable cloud provider snapshots
		  depends_on                  = ["mongodbatlas_privatelink_endpoint_service.test"]
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, clusterName)
}

func configAzureWithNetworkPeering(projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_network_container" "test" {
			project_id       = "%[1]s"
			atlas_cidr_block = "%[8]s"
			provider_name    = "%[2]s"
			region           = "%[9]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id            = "%[1]s"
			atlas_cidr_block      = "192.168.0.0/21"
			container_id          = mongodbatlas_network_container.test.container_id
			provider_name         = "%[2]s"
			azure_directory_id    = "%[3]s"
			azure_subscription_id = "%[4]s"
			resource_group_name   = "%[5]s"
			vnet_name             = "%[6]s"
		}

		resource "mongodbatlas_cluster" "with_azure_peering" {
			project_id   = "%[1]s"
			name         = "%[7]s"

			cluster_type = "REPLICASET"
			  replication_specs {
				num_shards = 1
				regions_config {
				  region_name     = "%[9]s"
				  electable_nodes = 3
				  priority        = 7
				  read_only_nodes = 0
				}
		  	}

			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "%[2]s"
			provider_disk_type_name     = "P6"
			provider_instance_size_name = "M10"

			depends_on = ["mongodbatlas_network_peering.test"]
		}
	`, projectID, providerName, directoryID, subcrptionID, resourceGroupName, vNetName, clusterName, atlasCidrBlock, region)
}

func configGCPWithNetworkPeering(gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion string) string {
	return fmt.Sprintf(`
		provider "google" {
			project     = "%[1]s"
			region      = "%[2]s"
		}

		resource "mongodbatlas_network_container" "test" {
			project_id       = "%[3]s"
			atlas_cidr_block = "192.168.192.0/18"
			provider_name    = "%[4]s"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id     = "%[3]s"
			container_id   = mongodbatlas_network_container.test.container_id
			provider_name  = "%[4]s"
			gcp_project_id = "%[1]s"
			network_name   = "default"
		}

		data "google_compute_network" "default" {
			name = "default"
		}

		resource "google_compute_network_peering" "gcp_peering" {
			name         = "%[5]s"
			network      = data.google_compute_network.default.self_link
			peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[3]s"
			name         = "%[6]s"
			
            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "%[7]s"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "%[4]s"
			provider_instance_size_name = "M10"

			depends_on = ["google_compute_network_peering.gcp_peering"]
		}
	`, gcpProjectID, gcpRegion, projectID, providerName, gcpPeeringName, clusterName, gcpClusterRegion)
}

func configAzureWithContainerID(projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[1]s"
			name         = "%[2]s"

            cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "%[4]s"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "%[3]s"
			provider_instance_size_name = "M10"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id            = "%[1]s"
			atlas_cidr_block      = "192.168.0.0/21"
			container_id          = mongodbatlas_cluster.test.container_id
			provider_name         = "%[3]s"
			azure_directory_id    = "%[5]s"
			azure_subscription_id = "%[6]s"
			resource_group_name   = "%[7]s"
			vnet_name             = "%[8]s"
		}
	`, projectID, clusterName, providerName, region, directoryID, subcrptionID, resourceGroupName, vNetName)
}

func configAWSWithContainerID(awsAccessKey, awsSecretKey, projectID, clusterName, providerName, region, vpcCIDRBlock, awsAccountID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = lower(replace("%[6]s", "_", "-"))
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[3]s"
			name         = "%[4]s"
			
			cluster_type = "REPLICASET"
		    replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "%[6]s"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			auto_scaling_disk_gb_enabled = false

			// Provider Settings "block"
			provider_name               = "%[5]s"
			provider_instance_size_name = "M10"
		}

		resource "aws_default_vpc" "default" {
			tags = {
				Name = "Default VPC"
			}
		}

		resource "mongodbatlas_network_peering" "mongo_peer" {
			accepter_region_name   = lower(replace("%[6]s", "_", "-"))
			project_id             = "%[3]s"
			container_id           = mongodbatlas_cluster.test.container_id
			provider_name          = "%[5]s"
			route_table_cidr_block = "%[7]s"
			vpc_id                 = aws_default_vpc.default.id
			aws_account_id         = "%[8]s"
		}

		resource "aws_vpc_peering_connection_accepter" "aws_peer" {
			vpc_peering_connection_id = mongodbatlas_network_peering.mongo_peer.connection_id
			auto_accept               = true

			tags = {
				Side = "Accepter"
			}
		}
	`, awsAccessKey, awsSecretKey, projectID, clusterName, providerName, region, vpcCIDRBlock, awsAccountID)
}

func configGCPWithContainerID(gcpProjectID, gcpRegion, projectID, clusterName, providerName, gcpClusterRegion, gcpPeeringName string) string {
	return fmt.Sprintf(`
		provider "google" {
			project     = "%[1]s"
			region      = "%[2]s"
		}

		resource "mongodbatlas_cluster" "test" {
			project_id   = "%[3]s"
			name         = "%[4]s"
			
            cluster_type = "REPLICASET"
			replication_specs {
			  num_shards = 1
			  regions_config {
			     region_name     = "%[6]s"
			     electable_nodes = 3
			     priority        = 7
                 read_only_nodes = 0
		       }
		    }

			auto_scaling_disk_gb_enabled = true

			// Provider Settings "block"
			provider_name               = "%[5]s"
			provider_instance_size_name = "M10"
		}

		resource "mongodbatlas_network_peering" "test" {
			project_id     = "%[3]s"
			container_id   = mongodbatlas_cluster.test.container_id
			provider_name  = "%[5]s"
			gcp_project_id = "%[1]s"
			network_name   = "default"
		}

		data "google_compute_network" "default" {
			name = "default"
		}

		resource "google_compute_network_peering" "gcp_peering" {
			name         = "%[7]s"
			network      = data.google_compute_network.default.self_link
			peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
		}
	`, gcpProjectID, gcpRegion, projectID, clusterName, providerName, gcpClusterRegion, gcpPeeringName)
}

func configAWSWithAutoscaling(projectID, name, backupEnabled, autoDiskEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName, instanceSizeName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id                              = %[1]q
			name                                    = %[2]q
			disk_size_gb                            = 100

			cluster_type = "REPLICASET"
			replication_specs {
				num_shards = 1
				regions_config {
				region_name     = "US_WEST_2"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
			}
			cloud_backup                            = %[3]s
			auto_scaling_disk_gb_enabled            = %[4]s
			auto_scaling_compute_enabled            = %[5]s
			auto_scaling_compute_scale_down_enabled = %[6]s

			//Provider Settings "block"
			provider_name                                   = "AWS"
			provider_auto_scaling_compute_min_instance_size = %[7]q
			provider_auto_scaling_compute_max_instance_size = %[8]q
			provider_instance_size_name                     = %[9]q

			lifecycle { // To simulate if there a new instance size name to avoid scale cluster down to original value
				ignore_changes = [provider_instance_size_name]
			}
		}

		data "mongodbatlas_cluster" "test" {
			project_id = mongodbatlas_cluster.test.project_id
			name 	     = mongodbatlas_cluster.test.name
		}
	`, projectID, name, backupEnabled, autoDiskEnabled, autoScalingEnabled, scaleDownEnabled, minSizeName, maxSizeName, instanceSizeName)
}

func configGCPRegionName(
	orgID, projectName, name, regionName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "cluster_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cluster" "test" {
  project_id                   = mongodbatlas_project.cluster_project.id
  name                         = %[3]q
  auto_scaling_disk_gb_enabled = true
  provider_name                = "GCP"
  disk_size_gb                 = 10
  provider_instance_size_name  = "M10"
  num_shards                   = 1
  provider_region_name         = %[4]q
}
	`, orgID, projectName, name, regionName)
}

func configRegions(
	orgID, projectName, name, replications string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "cluster_project" {
		name   = %[2]q
		org_id = %[1]q
	}
	resource "mongodbatlas_cluster" "test" {
	  project_id              = mongodbatlas_project.cluster_project.id
	  name                    = "%[3]s"
	  disk_size_gb            = 400
	  num_shards              = 3
	  cloud_backup            = true
	  cluster_type            = "GEOSHARDED"
	  // Provider Settings "block"
	  provider_name               = "AWS"
	  provider_disk_iops          = 1200
	  provider_instance_size_name = "M30"
	  %[4]s

		lifecycle {
		# avoid cluster has been auto-scaled to different instance size
		ignore_changes = [provider_instance_size_name, disk_size_gb]
	  }
	}
	`, orgID, projectName, name, replications)
}

func configAWSPaused(projectID, name string, backupEnabled, paused bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "test" {
  project_id                   = %[1]q
  name                         = %[2]q
  disk_size_gb                 = 100
  cluster_type = "REPLICASET"
  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_WEST_2"
      electable_nodes = 3
      priority        = 7
      read_only_nodes = 0
    }
  }
  cloud_backup                 = %[3]t
  paused                       = %[4]t
  // Provider Settings "block"

  provider_name               = "AWS"
  provider_instance_size_name = "M30"
}
	`, projectID, name, backupEnabled, paused)
}

func TestIsMultiRegionCluster(t *testing.T) {
	tests := []struct {
		name     string
		repSpecs []matlas.ReplicationSpec
		want     bool
	}{
		{
			name:     "No ReplicationSpecs",
			repSpecs: []matlas.ReplicationSpec{},
			want:     false,
		},
		{
			name: "Single ReplicationSpec Single Region",
			repSpecs: []matlas.ReplicationSpec{
				{
					RegionsConfig: map[string]matlas.RegionsConfig{
						"region1": {},
					},
				},
			},
			want: false,
		},
		{
			name: "Single ReplicationSpec Multiple Regions",
			repSpecs: []matlas.ReplicationSpec{
				{
					RegionsConfig: map[string]matlas.RegionsConfig{
						"region1": {},
						"region2": {},
					},
				},
			},
			want: true,
		},
		{
			name: "Multiple ReplicationSpecs",
			repSpecs: []matlas.ReplicationSpec{
				{
					RegionsConfig: map[string]matlas.RegionsConfig{
						"region1": {},
					},
				},
				{
					RegionsConfig: map[string]matlas.RegionsConfig{
						"region2": {},
					},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := clustersvc.IsMultiRegionCluster(tt.repSpecs); got != tt.want {
				t.Errorf("isMultiRegionCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateProviderRegionName(t *testing.T) {
	tests := []struct {
		name               string
		clusterType        string
		providerRegionName string
		repSpecs           []matlas.ReplicationSpec
		wantErr            bool
	}{
		{
			name:               "Single Region with Provider Name",
			clusterType:        "REPLICASET",
			providerRegionName: "us-east-1",
			repSpecs:           []matlas.ReplicationSpec{{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}}}},
			wantErr:            false,
		},
		{
			name:               "Single Region without Provider Name",
			clusterType:        "REPLICASET",
			providerRegionName: "",
			repSpecs:           []matlas.ReplicationSpec{{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}}}},
			wantErr:            false,
		},
		{
			name:               "Multi Region with Provider Name",
			clusterType:        "REPLICASET",
			providerRegionName: "us-east-1",
			repSpecs: []matlas.ReplicationSpec{
				{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}, "region2": {}}},
			},
			wantErr: true,
		},
		{
			name:               "Multi Region without Provider Name",
			clusterType:        "REPLICASET",
			providerRegionName: "",
			repSpecs: []matlas.ReplicationSpec{
				{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}, "region2": {}}},
			},
			wantErr: false,
		},
		{
			name:               "Geosharded with Provider Name",
			clusterType:        "GEOSHARDED",
			providerRegionName: "us-east-1",
			repSpecs:           []matlas.ReplicationSpec{{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}}}},
			wantErr:            true,
		},
		{
			name:               "Geosharded without Provider Name",
			clusterType:        "GEOSHARDED",
			providerRegionName: "",
			repSpecs:           []matlas.ReplicationSpec{{RegionsConfig: map[string]matlas.RegionsConfig{"region1": {}}}},
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := clustersvc.ValidateProviderRegionName(tt.clusterType, tt.providerRegionName, tt.repSpecs)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateProviderRegionName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
