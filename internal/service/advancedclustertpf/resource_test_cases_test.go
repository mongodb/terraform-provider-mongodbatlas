package advancedclustertpf_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
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

func SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(projectID, clusterName, 50),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(50),
			},
			{
				Config: configShardedOldSchemaDiskSizeGBElectableLevel(projectID, clusterName, 55),
				Check:  checkShardedOldSchemaDiskSizeGBElectableLevel(55),
			},
		},
	}
}

func configShardedOldSchemaDiskSizeGBElectableLevel(projectID, name string, diskSizeGB int) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id = %[1]q
		name = %[2]q
		backup_enabled = false
		mongo_db_major_version = "7.0"
		cluster_type   = "SHARDED"

		replication_specs = [{
			num_shards = 2

			region_configs = [{
			electable_specs = {
				instance_size = "M10"
				node_count    = 3
				disk_size_gb  = %[3]d
			}
			analytics_specs = {
				instance_size = "M10"
				node_count    = 0
				disk_size_gb  = %[3]d
			}
			provider_name = "AWS"
			priority      = 7
			region_name   = "US_EAST_1"
			},
			]
		}]
	}
	`, projectID, name, diskSizeGB)
}

func checkShardedOldSchemaDiskSizeGBElectableLevel(diskSizeGB int) resource.TestCheckFunc {
	return checkAggr(
		[]string{},
		map[string]string{
			"replication_specs.0.num_shards": "2",
			"disk_size_gb":                   fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.electable_specs.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
			"replication_specs.0.region_configs.0.analytics_specs.disk_size_gb": fmt.Sprintf("%d", diskSizeGB),
		})
}

func checkAggr(attrsSet []string, attrsMap map[string]string, extra ...resource.TestCheckFunc) resource.TestCheckFunc {
	checks := []resource.TestCheckFunc{}
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	checks = acc.AddAttrSetChecks(resourceName, checks, attrsSet...)
	checks = append(checks, extra...)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func SymmetricShardedOldSchema(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configShardedOldSchemaMultiCloud(projectID, clusterName, 2, "M10", &configServerManagementModeFixedToDedicated),
				Check:  checkShardedOldSchemaMultiCloud(clusterName, 2, "M10", false, &configServerManagementModeFixedToDedicated),
			},
			{
				Config: configShardedOldSchemaMultiCloud(projectID, clusterName, 2, "M20", &configServerManagementModeAtlasManaged),
				Check:  checkShardedOldSchemaMultiCloud(clusterName, 2, "M20", false, &configServerManagementModeAtlasManaged),
			},
		},
	}
}

func configShardedOldSchemaMultiCloud(projectID, name string, numShards int, analyticsSize string, configServerManagementMode *string) string {
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
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "SHARDED"
		%[5]s

		replication_specs = [{
			num_shards = %[3]d
			region_configs = [{
				electable_specs = {
					instance_size = "M10"
					node_count    = 3
				}
				analytics_specs = {
					instance_size = %[4]q
					node_count    = 1
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}, {
				electable_specs = {
					instance_size = "M10"
					node_count    = 2
				}
				provider_name = "AZURE"
				priority      = 6
				region_name   = "US_EAST_2"
			},]
		},]
	}
	`, projectID, name, numShards, analyticsSize, rootConfig)
}

func checkShardedOldSchemaMultiCloud(name string, numShards int, analyticsSize string, verifyExternalID bool, configServerManagementMode *string) resource.TestCheckFunc {
	additionalChecks := []resource.TestCheckFunc{
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.electable_specs.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.0.analytics_specs.disk_iops", acc.IntGreatThan(0)),
		resource.TestCheckResourceAttrWith(resourceName, "replication_specs.0.region_configs.1.electable_specs.disk_iops", acc.IntGreatThan(0)),
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
		)
	}

	return checkAggr(
		[]string{"project_id", "replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"},
		map[string]string{
			"name":                           name,
			"replication_specs.0.num_shards": strconv.Itoa(numShards),
			"replication_specs.0.region_configs.0.analytics_specs.instance_size": analyticsSize,
		},
		additionalChecks...)
}

func BasicTenantTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID          = acc.ProjectIDExecution(t)
		clusterName        = acc.RandomClusterName()
		clusterNameUpdated = acc.RandomClusterName()
	)
	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		Steps: []resource.TestStep{
			{
				Config: configTenant(projectID, clusterName),
				Check:  checkTenant(projectID, clusterName),
			},
			{
				Config: configTenant(projectID, clusterNameUpdated),
				Check:  checkTenant(projectID, clusterNameUpdated),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func configTenant(projectID, name string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"

			replication_specs = [{
				region_configs = [{
					electable_specs = {
						instance_size = "M5"
					}
					provider_name         = "TENANT"
					backing_provider_name = "AWS"
					region_name           = "US_EAST_1"
					priority              = 7
				}]
			}]
		}
	`, projectID, name)
}

func checkTenant(projectID, name string) resource.TestCheckFunc {
	attrsSet := []string{"replication_specs.#", "replication_specs.0.id", "replication_specs.0.region_configs.#"}
	attrsMap := map[string]string{
		"project_id":                           projectID,
		"name":                                 name,
		"termination_protection_enabled":       "false",
		"global_cluster_self_managed_sharding": "false",
		"labels.#":                             "0",
	}
	checks := acc.AddAttrSetChecks(resourceName, nil, attrsSet...)
	checks = acc.AddAttrChecks(resourceName, checks, attrsMap)
	return resource.ComposeAggregateTestCheckFunc(checks...)
}

func TenantUpgrade(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	return &resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, nil, projectID, clusterName),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: acc.ConvertAdvancedClusterToTPF(t, configTenant(projectID, clusterName)),
				Check:  checkTenant(projectID, clusterName),
			},
			{
				Config: acc.ConvertAdvancedClusterToTPF(t, configTenantUpgraded(projectID, clusterName)),
				Check:  checksTenantUpgraded(projectID, clusterName),
			},
		},
	}
}

func configTenantUpgraded(projectID, name string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "REPLICASET"
		
		replication_specs {
			region_configs {
				priority        = 7
				provider_name = "AWS"
				region_name     = "US_EAST_1"
				electable_specs {
					node_count = 3
					instance_size = "M10"
				}
			}
		}
	}
	`, projectID, name)
}

func enableChecksLatestTpf(checkMap map[string]string) map[string]string {
	newMap := map[string]string{}
	for k, v := range checkMap {
		modifiedKey := strings.ReplaceAll(k, "electable_specs.0", "electable_specs")
		newMap[modifiedKey] = v
	}
	return newMap
}

func checksTenantUpgraded(projectID, name string) resource.TestCheckFunc {
	originalChecks := checkTenant(projectID, name)
	checks := []resource.TestCheckFunc{originalChecks}
	checkMap := map[string]string{
		"replication_specs.0.region_configs.0.electable_specs.0.node_count":    "3",
		"replication_specs.0.region_configs.0.electable_specs.0.instance_size": "M10",
		"replication_specs.0.region_configs.0.provider_name":                   "AWS",
	}
	if config.AdvancedClusterV2Schema() {
		checkMap = enableChecksLatestTpf(checkMap)
	}
	checks = acc.AddAttrChecks(resourceName, checks, checkMap)
	return resource.ComposeAggregateTestCheckFunc(originalChecks, resource.ComposeAggregateTestCheckFunc(checks...))
}

func ReplicasetAdvConfigUpdate(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
		fullUpdate  = `
	backup_enabled = true
	bi_connector_config = {
		enabled = true
	}
	# config_server_management_mode = "ATLAS_MANAGED" UNSTABLE: After applying this test step, the non-refresh plan was not empty
	labels = [{
		key   = "env"
		value = "test"
	}]
	tags = [{
		key   = "env"
		value = "test"
	}]
	mongo_db_major_version = "8.0"
	pit_enabled = true
	redact_client_log_data = true
	replica_set_scaling_strategy = "NODE_TYPE"
	# retain_backups_enabled = true # only set on delete
	root_cert_type = "ISRGROOTX1"
	# termination_protection_enabled = true # must be reset to false to enable delete
	version_release_system = "CONTINUOUS"
	
	advanced_configuration = {
		change_stream_options_pre_and_post_images_expire_after_seconds = 100
		default_read_concern                                           = "available"
		default_write_concern                                          = "majority"
		javascript_enabled                                             = true
		minimum_enabled_tls_protocol                                   = "TLS1_0"
		no_table_scan                                                  = true
		sample_refresh_interval_bi_connector                           = 310
		sample_size_bi_connector                                       = 110
		transaction_lifetime_limit_seconds                             = 300
	}
`
	// # oplog_min_retention_hours                                      = 5.5
	// # oplog_size_mb                                                  = 1000
	// # fail_index_key_too_long 								        = true # only valid for MongoDB version 4.4 and earlier
	)
	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, ""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state_name", "IDLE"),
					resource.TestCheckResourceAttr(resourceName, "timeouts.create", "2000s"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.container_id.AWS:US_EAST_1", "67345bd9905b8c30c54fd220"),
				),
			},
			{
				Config: configBasic(projectID, clusterName, fullUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "mongo_db_major_version", "8.0"),
					resource.TestCheckResourceAttr(resourceName, "advanced_configuration.change_stream_options_pre_and_post_images_expire_after_seconds", "100"),
				),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func shardedBasic(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)
	return &resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: configSharded(projectID, clusterName, false),
				Check:  resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
			},
			{
				Config: configSharded(projectID, clusterName, true),
				Check:  resource.TestCheckResourceAttr(resourceName, "name", clusterName),
			},
			acc.TestStepImportCluster(resourceName),
		},
	}
}

func configSharded(projectID, clusterName string, withUpdate bool) string {
	var autoScaling, analyticsSpecs string
	if withUpdate {
		autoScaling = `
			auto_scaling = {
				disk_gb_enabled = true
			}`
		analyticsSpecs = `
			analytics_specs = {
				instance_size   = "M30"
				node_count      = 1
				ebs_volume_type = "PROVISIONED"
				disk_iops       = 2000
			}`
	}
	// SDK v2 Implementation receives many warnings, one of them: `.replication_specs[1].region_configs[0].analytics_specs[0].disk_iops: was cty.NumberIntVal(2000), but now cty.NumberIntVal(1000)`
	// Therefore, in TPF we are forced to set the value that will be returned by the API (1000)
	// The rule is: For any replication spec, the `(analytics|electable|read_only)_spec.disk_iops` must be the same across all region_configs
	// The API raises no errors, but the response reflects this rule
	analyticsSpecsForSpec2 := strings.ReplaceAll(analyticsSpecs, "2000", "1000")
	return fmt.Sprintf(`
	resource "mongodbatlas_advanced_cluster" "test" {
		project_id   = %[1]q
		name         = %[2]q
		cluster_type = "SHARDED"

		replication_specs = [
			{ # shard 1
			region_configs = [{
				electable_specs = {
					instance_size   = "M30"
					disk_iops       = 2000
					node_count      = 3
					ebs_volume_type = "PROVISIONED"
				}
				%[3]s
				%[4]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
				}]
				},
				{ # shard 2
			region_configs = [{
				electable_specs = {
					instance_size   = "M30"
					ebs_volume_type = "PROVISIONED"
					disk_iops       = 1000
					node_count      = 3
				}
				%[3]s
				%[5]s
				provider_name = "AWS"
				priority      = 7
				region_name   = "EU_WEST_1"
			}]
		}]
	}
	`, projectID, clusterName, autoScaling, analyticsSpecs, analyticsSpecsForSpec2)
}
