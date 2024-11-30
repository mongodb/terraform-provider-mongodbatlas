package tc

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func SymmetricShardedOldSchemaDiskSizeGBAtElectableLevel(t *testing.T, orgID, projectName, clusterName string) *resource.TestCase {
	t.Helper()
	return &resource.TestCase{
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
	}
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

			replication_specs = [{
				num_shards = 2

				region_configs = [{
				electable_specs = {
					instance_size = "M10"
					node_count    = 3
					disk_size_gb  = %[4]d
				}
				analytics_specs = {
					instance_size = "M10"
					node_count    = 0
					disk_size_gb  = %[4]d
				}
				provider_name = "AWS"
				priority      = 7
				region_name   = "US_EAST_1"
				},
				]
			}]
		}
	`, orgID, projectName, name, diskSizeGB)
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

func SymmetricShardedOldSchema(t *testing.T, orgID, projectName, clusterName string) *resource.TestCase {
	t.Helper()
	return &resource.TestCase{
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
	}
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

			replication_specs = [{
				num_shards = %[4]d
				region_configs = [{
					electable_specs = {
						instance_size = "M10"
						node_count    = 3
					}
					analytics_specs = {
						instance_size = %[5]q
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
	`, orgID, projectName, name, numShards, analyticsSize, rootConfig)
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
