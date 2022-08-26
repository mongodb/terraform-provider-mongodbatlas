package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudBackupSnapshot_basic(t *testing.T) {
	var (
		cloudProviderSnapshot matlas.CloudProviderSnapshot
		projectID             = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName           = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description           = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays       = "1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotExists("mongodbatlas_cloud_backup_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  
  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_2"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

data "mongodbatlas_cloud_backup_snapshot" "test" {
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
}
`, projectID, clusterName, description, retentionInDays)
}
