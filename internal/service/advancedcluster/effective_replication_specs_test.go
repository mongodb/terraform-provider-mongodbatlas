package advancedcluster_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// TestAccAdvancedCluster_effectiveReplicationSpecs verifies that:
// 1. Resource replication_specs shows user-configured values
// 2. Data source replication_specs shows user-configured values
// 3. Data source effective_replication_specs shows actual running configuration
func TestAccAdvancedCluster_effectiveReplicationSpecs(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
		instanceSize           = "M10"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configEffectiveReplicationSpecs(projectID, clusterName, instanceSize, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Resource has user-configured values
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.instance_size", instanceSize),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.node_count", "3"),

					// Singular data source has both configured and effective values
					resource.TestCheckResourceAttr(dataSourceName, "replication_specs.0.region_configs.0.electable_specs.instance_size", instanceSize),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.instance_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.node_count"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.disk_size_gb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.disk_iops"),

					// Plural data source has both configured and effective values
					resource.TestCheckResourceAttrWith(dataSourcePluralName, "results.#", acc.IntGreatThan(0)),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.effective_replication_specs.0.region_configs.0.electable_specs.instance_size"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "results.0.effective_replication_specs.0.region_configs.0.electable_specs.node_count"),
				),
			},
			acc.TestStepImportCluster(resourceName, "replication_specs"),
		},
	})
}

// TestAccAdvancedCluster_effectiveReplicationSpecsWithAutoScaling verifies effective_replication_specs
// shows actual running values when auto-scaling is enabled and changes instance size
func TestAccAdvancedCluster_effectiveReplicationSpecsWithAutoScaling(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
		instanceSize           = "M10"
		maxInstanceSize        = "M40"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configEffectiveReplicationSpecsWithAutoScaling(projectID, clusterName, instanceSize, maxInstanceSize, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Resource has user-configured instance size
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.electable_specs.instance_size", instanceSize),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.compute_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.auto_scaling.compute_max_instance_size", maxInstanceSize),

					// Data source effective_replication_specs shows actual running values
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.instance_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.disk_size_gb"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.disk_iops"),
				),
			},
		},
	})
}

// TestAccAdvancedCluster_effectiveReplicationSpecsReadOnlyAnalytics verifies effective_replication_specs
// includes read_only_specs and analytics_specs when configured
func TestAccAdvancedCluster_effectiveReplicationSpecsReadOnlyAnalytics(t *testing.T) {
	var (
		projectID, clusterName = acc.ProjectIDExecutionWithCluster(t, 3)
		instanceSize           = "M10"
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             acc.CheckDestroyCluster,
		Steps: []resource.TestStep{
			{
				Config: configEffectiveReplicationSpecsWithReadOnlyAnalytics(projectID, clusterName, instanceSize, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Resource has configured specs
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.read_only_specs.instance_size", instanceSize),
					resource.TestCheckResourceAttr(resourceName, "replication_specs.0.region_configs.0.analytics_specs.instance_size", instanceSize),

					// Data source effective_replication_specs shows actual running values for all spec types
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.electable_specs.instance_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.read_only_specs.instance_size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "effective_replication_specs.0.region_configs.0.analytics_specs.instance_size"),
				),
			},
		},
	})
}

func configEffectiveReplicationSpecs(projectID, clusterName, instanceSize string, nodeCount int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [
				{
					region_configs = [
						{
							priority      = 7
							provider_name = "AWS"
							region_name   = "US_EAST_1"
							electable_specs = {
								node_count    = %[4]d
								instance_size = %[3]q
							}
						}
					]
				}
			]
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name       = mongodbatlas_advanced_cluster.test.name
			depends_on = [mongodbatlas_advanced_cluster.test]
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, projectID, clusterName, instanceSize, nodeCount)
}

func configEffectiveReplicationSpecsWithAutoScaling(projectID, clusterName, instanceSize, maxInstanceSize string, nodeCount int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [
				{
					region_configs = [
						{
							priority      = 7
							provider_name = "AWS"
							region_name   = "US_EAST_1"
							electable_specs = {
								node_count    = %[5]d
								instance_size = %[3]q
							}
							auto_scaling = {
								compute_enabled          = true
								compute_max_instance_size = %[4]q
							}
						}
					]
				}
			]
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name       = mongodbatlas_advanced_cluster.test.name
			depends_on = [mongodbatlas_advanced_cluster.test]
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, projectID, clusterName, instanceSize, maxInstanceSize, nodeCount)
}

func configEffectiveReplicationSpecsWithReadOnlyAnalytics(projectID, clusterName, instanceSize string, nodeCount int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_advanced_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			replication_specs = [
				{
					region_configs = [
						{
							priority      = 7
							provider_name = "AWS"
							region_name   = "US_EAST_1"
							electable_specs = {
								node_count    = %[4]d
								instance_size = %[3]q
							}
							read_only_specs = {
								node_count    = 1
								instance_size = %[3]q
							}
							analytics_specs = {
								node_count    = 1
								instance_size = %[3]q
							}
						}
					]
				}
			]
		}

		data "mongodbatlas_advanced_cluster" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			name       = mongodbatlas_advanced_cluster.test.name
			depends_on = [mongodbatlas_advanced_cluster.test]
		}

		data "mongodbatlas_advanced_clusters" "test" {
			project_id = mongodbatlas_advanced_cluster.test.project_id
			depends_on = [mongodbatlas_advanced_cluster.test]
		}
	`, projectID, clusterName, instanceSize, nodeCount)
}
