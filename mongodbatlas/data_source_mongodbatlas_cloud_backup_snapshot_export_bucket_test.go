package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupDSCloudBackupSnapshotExportBucket_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		snapshotExportBackup matlas.CloudProviderSnapshotExportBucket
		projectID            = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName           = os.Getenv("AWS_S3_BUCKET")
		iamRoleID            = os.Getenv("IAM_ROLE_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketConfig(projectID, iamRoleID, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupSnapshotExportBucketExists("mongodbatlas_cloud_backup_snapshot_export_bucket.test", &snapshotExportBackup),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_bucket.test", "iam_role_id"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_bucket.test", "bucket_name"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_bucket.test", "cloud_provider"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketConfig(projectID, iamRoleID, bucketName string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
			project_id   = "%[1]s"
			
    	  	iam_role_id = "%[2]s"
       		bucket_name = "%[3]s"
       		cloud_provider = "AWS"
		}

data "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
  id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.id
}
`, projectID, iamRoleID, bucketName)
}
