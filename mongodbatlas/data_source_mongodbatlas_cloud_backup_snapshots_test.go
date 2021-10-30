package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudBackupSnapshots_basic(t *testing.T) {
	var (
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName     = acctest.RandomWithPrefix("test-acc")
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotsConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasCloudBackupSnapshots_withPagination(t *testing.T) {
	var (
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName     = acctest.RandomWithPrefix("test-acc")
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotsConfigWithPagination(projectID, clusterName, description, retentionInDays, 1, 5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_backup_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotsConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  disk_size_gb = 5

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

data "mongodbatlas_cloud_backup_snapshots" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotsConfigWithPagination(projectID, clusterName, description, retentionInDays string, pageNum, itemPage int) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  disk_size_gb = 5

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "EU_CENTRAL_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

data "mongodbatlas_cloud_backup_snapshots" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  page_num = %[5]d
  items_per_page = %[6]d
}
	`, projectID, clusterName, description, retentionInDays, pageNum, itemPage)
}
