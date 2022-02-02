package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceMongoDBAtlasCloudBackupSnapshotExportBuckets_basic(t *testing.T) {
	var (
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		datasourceName = "mongodbatlas_cloud_backup_snapshot_export_buckets"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketsConfig(projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "iam_role_id"),
					resource.TestCheckResourceAttr(datasourceName, "results.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(datasourceName, "cloud_provider", "AWS"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketsConfig(projectID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_provider_access" "test" {
		project_id = "%[1]s"
		provider_name = "AWS"
	 }

	resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
			project_id   = "%[1]s"
			
    	  	iam_role_id = mongodbatlas_cloud_provider_access.test.role_id
       		bucket_name = "example-bucket"
       		cloud_provider = "AWS"
		}

data "mongodbatlas_cloud_backup_snapshot_export_buckets" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
}
	`, projectID)
}
