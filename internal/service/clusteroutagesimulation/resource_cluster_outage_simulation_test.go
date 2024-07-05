package clusteroutagesimulation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName   = "mongodbatlas_cluster_outage_simulation.test_outage"
	dataSourceName = "data.mongodbatlas_cluster_outage_simulation.test"
)

func TestAccOutageSimulationCluster_SingleRegion_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSingleRegion(projectID, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(resourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
				),
			},
		},
	})
}

func TestAccOutageSimulationCluster_MultiRegion_basic(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		clusterName = acc.RandomClusterName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configMultiRegion(projectID, clusterName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(resourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(resourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),

					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "outage_filters.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "start_request_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "simulation_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "state"),
				),
			},
		},
	})
}

func configSingleRegion(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
				project_id                  = %[1]q
				name                        = %[2]q
				provider_name               = "AWS"
				provider_region_name        = "US_WEST_2"
				provider_instance_size_name = "M10"
			}

			resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
				project_id = %[1]q
				cluster_name = %[2]q
				outage_filters {
					cloud_provider = "AWS"
					region_name    = "US_WEST_2"
				}
				depends_on = ["mongodbatlas_cluster.test"]
			}

			data "mongodbatlas_cluster_outage_simulation" "test" {
				project_id = %[1]q
				cluster_name = %[2]q
				depends_on = [mongodbatlas_cluster_outage_simulation.test_outage]
			}		
	`, projectID, clusterName)
}

func configMultiRegion(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "test" {
			project_id   = %[1]q
			name         = %[2]q
			cluster_type = "REPLICASET"
			
			provider_name               = "AWS"
			provider_instance_size_name = "M10"
			
			replication_specs {
				num_shards = 1
				regions_config {
				region_name     = "US_EAST_1"
				electable_nodes = 3
				priority        = 7
				read_only_nodes = 0
				}
				regions_config {
				region_name     = "US_EAST_2"
				electable_nodes = 2
				priority        = 6
				read_only_nodes = 0
				}
				regions_config {
				region_name     = "US_WEST_2"
				electable_nodes = 2
				priority        = 5
				read_only_nodes = 2
				}
			}
		}

		resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
			project_id   = %[1]q
			cluster_name = %[2]q

			outage_filters {
				cloud_provider = "AWS"
				region_name    = "US_EAST_1"
			}
			outage_filters {
					cloud_provider = "AWS"
					region_name    = "US_EAST_2"
			}
			depends_on = ["mongodbatlas_cluster.test"]
		}

		data "mongodbatlas_cluster_outage_simulation" "test" {
			project_id = %[1]q
			cluster_name = %[2]q
			depends_on = [mongodbatlas_cluster_outage_simulation.test_outage]
		}		
	`, projectID, clusterName)
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_outage_simulation" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().ClusterOutageSimulationApi.GetOutageSimulation(context.Background(), ids["project_id"], ids["cluster_name"]).Execute()
		if err == nil {
			return fmt.Errorf("cluster outage simulation for project (%s) and cluster (%s) still exists", ids["project_id"], ids["cluster_name"])
		}
	}
	return nil
}
