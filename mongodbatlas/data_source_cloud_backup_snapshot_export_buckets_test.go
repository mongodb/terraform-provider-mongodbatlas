package mongodbatlas_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccBackupDSCloudBackupSnapshotExportBuckets_basic(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName     = os.Getenv("AWS_S3_BUCKET")
		iamRoleID      = os.Getenv("IAM_ROLE_ID")
		datasourceName = "mongodbatlas_cloud_backup_snapshot_export_buckets"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketsConfig(projectID, bucketName, iamRoleID),
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

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportBucketsConfig(projectID, iamRoleID, bucketName string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id   = "%[1]s"
  iam_role_id = "%[2]s"
  bucket_name = "%[3]s"
  cloud_provider = "AWS"
}

data "mongodbatlas_cloud_backup_snapshot_export_buckets" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
}`, projectID, iamRoleID, bucketName)
}
