package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccOutageSimulationCluster_SingleRegion_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		dataSourceName = "mongodbatlas_cluster_outage_simulation.test_outage"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-project")
		clusterName    = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterOutageSimulationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasClusterOutageSimulationConfigSingleRegion(projectName, orgID, clusterName),
				Check: resource.ComposeTestCheckFunc(
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
	SkipTestExtCred(t)
	var (
		dataSourceName = "mongodbatlas_cluster_outage_simulation.test_outage"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc-project")
		clusterName    = acctest.RandomWithPrefix("test-acc-cluster")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasClusterOutageSimulationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasClusterOutageSimulationConfigMultiRegion(projectName, orgID, clusterName),
				Check: resource.ComposeTestCheckFunc(
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

func testAccDataSourceMongoDBAtlasClusterOutageSimulationConfigSingleRegion(projectName, orgID, clusterName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "outage_project" {
		name   = "%s"
		org_id = "%s"
	}

	resource "mongodbatlas_cluster" "atlas_cluster" {
		project_id                  = mongodbatlas_project.outage_project.id
   		provider_name               = "AWS"
   		name                        = "%s"
   		backing_provider_name       = "AWS"
   		provider_region_name        = "US_EAST_1"
   		provider_instance_size_name = "M10"
	  }

	  resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
		project_id = mongodbatlas_project.outage_project.id
		cluster_name = mongodbatlas_cluster.atlas_cluster.name
		 outage_filters {
		  cloud_provider = "AWS"
		  region_name    = "US_EAST_1"
		}
	}
	`, projectName, orgID, clusterName)
}

func testAccDataSourceMongoDBAtlasClusterOutageSimulationConfigMultiRegion(projectName, orgID, clusterName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_project" "outage_project" {
		name   = "%s"
		org_id = "%s"
	}

	resource "mongodbatlas_cluster" "atlas_cluster" {
		project_id   = mongodbatlas_project.outage_project.id
		name         = "%s"
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
			region_name     = "US_WEST_1"
			electable_nodes = 2
			priority        = 5
			read_only_nodes = 2
		  }
		}
	  }

	  resource "mongodbatlas_cluster_outage_simulation" "test_outage" {
		project_id = mongodbatlas_project.outage_project.id
		cluster_name = mongodbatlas_cluster.atlas_cluster.name
		 outage_filters {
		  cloud_provider = "AWS"
		  region_name    = "US_EAST_1"
		}
		outage_filters {
			   cloud_provider = "AWS"
			   region_name    = "US_EAST_2"
		}
	}
	`, projectName, orgID, clusterName)
}

func testAccCheckMongoDBAtlasClusterOutageSimulationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cluster_outage_simulation" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.ClusterOutageSimulation.GetOutageSimulation(context.Background(), ids["project_id"], ids["cluster_name"])
		if err == nil {
			return fmt.Errorf("cluster outage simulation for project (%s) and cluster (%s) still exists", ids["project_id"], ids["cluster_name"])
		}
	}

	return nil
}
