package clusteroutagesimulation_test

import (
	"context"
	"fmt"
	"regexp"
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
	resource.ParallelTest(t, *singleRegionTestCase(t))
}

func singleRegionTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		singleRegionRequest = acc.ClusterRequest{
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{Region: "US_WEST_2", InstanceSize: "M10"},
			},
		}
		clusterInfo = acc.GetClusterInfo(t, &singleRegionRequest)
		clusterName = clusterInfo.Name
	)
	return &resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configSingleRegion(&clusterInfo),
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
	}
}

func TestAccOutageSimulationCluster_MultiRegion_basic(t *testing.T) {
	resource.ParallelTest(t, *multiRegionTestCase(t))
}

func multiRegionTestCase(t *testing.T) *resource.TestCase {
	t.Helper()
	var (
		multiRegionRequest = acc.ClusterRequest{ReplicationSpecs: []acc.ReplicationSpecRequest{
			{
				Region:    "US_EAST_1",
				NodeCount: 3,
				ExtraRegionConfigs: []acc.ReplicationSpecRequest{
					{Region: "US_EAST_2", NodeCount: 2, Priority: 6},
					{Region: "US_WEST_2", NodeCount: 2, Priority: 5, NodeCountReadOnly: 2},
				},
			},
		}}
		clusterInfo = acc.GetClusterInfo(t, &multiRegionRequest)
		clusterName = clusterInfo.Name
	)

	return &resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configMultiRegion(&clusterInfo),
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
	}
}

func configSingleRegion(info *acc.ClusterInfo) string {
	return fmt.Sprintf(`
			%[1]s
			resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
				project_id = %[2]q
				cluster_name = %[3]q
				outage_filters {
					cloud_provider = "AWS"
					region_name    = "US_WEST_2"
				}
				depends_on = [%[4]s]
			}

			data "mongodbatlas_cluster_outage_simulation" "test" {
				project_id = %[2]q
				cluster_name = %[3]q
				depends_on = [mongodbatlas_cluster_outage_simulation.test_outage]
			}		
	`, info.TerraformStr, info.ProjectID, info.Name, info.ResourceName)
}

func configMultiRegion(info *acc.ClusterInfo) string {
	return fmt.Sprintf(`
		%[1]s
		resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
			project_id   = %[2]q
			cluster_name = %[3]q

			outage_filters {
				cloud_provider = "AWS"
				region_name    = "US_EAST_1"
			}
			outage_filters {
					cloud_provider = "AWS"
					region_name    = "US_EAST_2"
			}
			depends_on = [%[4]s]
		}

		data "mongodbatlas_cluster_outage_simulation" "test" {
			project_id = %[2]q
			cluster_name = %[3]q
			depends_on = [mongodbatlas_cluster_outage_simulation.test_outage]
		}		
	`, info.TerraformStr, info.ProjectID, info.Name, info.ResourceName)
}

func TestAccClusterOutageSimulation_deleteOnCreateTimeout(t *testing.T) {
	var (
		singleRegionRequest = acc.ClusterRequest{
			ReplicationSpecs: []acc.ReplicationSpecRequest{
				{Region: "US_WEST_2", InstanceSize: "M10"},
			},
		}
		clusterInfo = acc.GetClusterInfo(t, &singleRegionRequest)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config:      configDeleteOnCreateTimeout(&clusterInfo, "1s", true),
				ExpectError: regexp.MustCompile("will run cleanup because delete_on_create_timeout is true"),
			},
		},
	})
}

func configDeleteOnCreateTimeout(info *acc.ClusterInfo, timeout string, deleteOnTimeout bool) string {
	return fmt.Sprintf(`
		%[1]s
		resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
			project_id = %[2]q
			cluster_name = %[3]q
			delete_on_create_timeout = %[5]t
			
			timeouts {
				create = %[4]q
			}
			
			outage_filters {
				cloud_provider = "AWS"
				region_name    = "US_WEST_2"
			}
			
			depends_on = [%[6]s]
		}
	`, info.TerraformStr, info.ProjectID, info.Name, timeout, deleteOnTimeout, info.ResourceName)
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
