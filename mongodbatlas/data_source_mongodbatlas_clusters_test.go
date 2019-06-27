package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasClusters_basic(t *testing.T) {
	var cluster matlas.Cluster

	resourceName := "mongodbatlas_cluster.test"
	dataSourceName := "data.mongodbatlas_clusters.test"
	projectID := "5cf5a45a9ccf6400e60981b6" // Modify until project data source is created.
	name := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasClustersConfig(projectID, name, "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasClusterExists(resourceName, &cluster),
					testAccCheckMongoDBAtlasClusterAttributes(&cluster, name),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disk_size_gb", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "mongo_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.#"),
					resource.TestCheckResourceAttrSet(resourceName, "replication_specs.0.regions_config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.replication_specs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "results.0.name"),
				),
			},
		},
	})

}

func testAccDataSourceMongoDBAtlasClustersConfig(projectID, name, backupEnabled string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "test" {
	project_id   = "%s"
	name         = "%s"
	disk_size_gb = 100
	num_shards   = 1
	
	replication_factor           = 3
	backup_enabled               = %s
	auto_scaling_disk_gb_enabled = true
	
	//Provider Settings "block"
	provider_name               = "AWS"
	provider_disk_iops          = 300
	provider_encrypt_ebs_volume = false
	provider_instance_size_name = "M40"
	provider_region_name        = "US_EAST_1"
}

data "mongodbatlas_clusters" "test" {
	project_id = mongodbatlas_cluster.test.project_id
}
`, projectID, name, backupEnabled)
}
