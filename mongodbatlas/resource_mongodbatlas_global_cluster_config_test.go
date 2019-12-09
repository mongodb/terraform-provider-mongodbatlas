package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasGlobalCluster_basic(t *testing.T) {
	var globalConfig matlas.GlobalCluster

	resourceName := "mongodbatlas_global_cluster_config.config"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	name := fmt.Sprintf("test-acc-global-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasGlobalClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasGlobalClusterConfig(projectID, name, "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasGlobalClusterExists(resourceName, &globalConfig),
					resource.TestCheckResourceAttrSet(resourceName, "managed_namespaces.#"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_zone_mappings.#"),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", name),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					testAccCheckMongoDBAtlasGlobalClusterAttributes(&globalConfig, 1),
					// resource.TestCheckResourceAttr(resourceName, "cluster_type", "GEOSHARDED"),
					// resource.TestCheckResourceAttr(resourceName, "replication_specs.#", "2"),
					// resource.TestCheckResourceAttr(resourceName, "replication_specs.0.regions_config.#", "1"),
					// resource.TestCheckResourceAttr(resourceName, "replication_specs.1.regions_config.#", "1"),
				),
			},
		},
	})

}

func testAccCheckMongoDBAtlasGlobalClusterExists(resourceName string, globalConfig *matlas.GlobalCluster) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		if globalConfigResp, _, err := conn.GlobalClusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"]); err == nil {
			*globalConfig = *globalConfigResp

			if len(globalConfig.CustomZoneMapping) > 0 || len(globalConfig.ManagedNamespaces) > 0 {
				return nil
			}
		}

		return fmt.Errorf("global config for cluster(%s) does not exist", rs.Primary.Attributes["cluster_name"])
	}
}

func testAccCheckMongoDBAtlasGlobalClusterAttributes(globalCluster *matlas.GlobalCluster, managedNamespacesCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(globalCluster.ManagedNamespaces) != managedNamespacesCount {
			return fmt.Errorf("bad managed namespaces: %s", globalCluster.ManagedNamespaces)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasGlobalClusterDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_global_cluster_config" {
			continue
		}

		// Try to find the cluster
		globalConfig, _, err := conn.GlobalClusters.Get(context.Background(), rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"])
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("No cluster named %s exists in group %s", rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["project_id"])) {
				return nil
			}
			return err
		}

		if len(globalConfig.CustomZoneMapping) > 0 || len(globalConfig.ManagedNamespaces) > 0 {
			return fmt.Errorf("global cluster configuration for cluster(%s) still exists", rs.Primary.Attributes["cluster_name"])
		}
	}

	return nil
}

func testAccMongoDBAtlasGlobalClusterConfig(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cluster" "test" {
		project_id              = "%s"
		name                    = "%s"
		disk_size_gb            = 80
		num_shards              = 1
		backup_enabled          = %s
		provider_backup_enabled = true
		cluster_type            = "GEOSHARDED"
		
		//Provider Settings "block"
		provider_name               = "AWS"
		provider_disk_iops          = 240
		provider_instance_size_name = "M30"
		
		replication_specs {
			zone_name  = "Zone 1"
			num_shards = 2
			regions_config {
			region_name     = "EU_CENTRAL_1"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}
		
		replication_specs { 
			zone_name  = "Zone 2"
			num_shards = 2
			regions_config {
			region_name     = "US_EAST_2"
			electable_nodes = 3
			priority        = 7
			read_only_nodes = 0
			}
		}
	}
	
	resource "mongodbatlas_global_cluster_config" "config" {
		project_id = mongodbatlas_cluster.test.project_id
		cluster_name = mongodbatlas_cluster.test.name
	
		managed_namespaces {
			db 				 = "mydata"
			collection 		 = "publishers"
			custom_shard_key = "city"
		}
	
		custom_zone_mappings {
			location ="CA"
			zone =  "Zone 1"
		}
	}
	`, projectID, name, backupEnabled)
}
