package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasAdvancedCluster_basic(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = "data.mongodbatlas_advanced_cluster.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasAdvancedClusterConfig(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasAdvancedCluster_multicloud(t *testing.T) {
	var (
		cluster        matlas.AdvancedCluster
		resourceName   = "mongodbatlas_advanced_cluster.test"
		dataSourceName = "data.mongodbatlas_advanced_cluster.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		name           = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAdvancedClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasAdvancedClusterMultiCloudConfig(projectID, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAdvancedClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasAdvancedClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.region_configs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasAdvancedClusterConfig(projectID, name string) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_advanced_cluster" "test" {
  project_id = mongodbatlas_advanced_cluster.test.project_id
  name 	     = mongodbatlas_advanced_cluster.test.name
}
	`, testAccMongoDBAtlasAdvancedClusterConfigTenant(projectID, name))
}

func testAccDataSourceMongoDBAtlasAdvancedClusterMultiCloudConfig(projectID, name string) string {
	return fmt.Sprintf(`
%s

data "mongodbatlas_advanced_cluster" "test" {
  project_id = mongodbatlas_advanced_cluster.test.project_id
  name 	     = mongodbatlas_advanced_cluster.test.name
}
	`, testAccMongoDBAtlasAdvancedClusterConfigMultiCloud(projectID, name))
}
