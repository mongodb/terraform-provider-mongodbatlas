package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupRSBackupSnapshotExportBucket_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		snapshotExportBucket matlas.CloudProviderSnapshotExportBucket
		resourceName         = "mongodbatlas_cloud_backup_snapshot_export_bucket.test"
		projectID            = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName           = os.Getenv("AWS_S3_BUCKET")
		iamRoleID            = os.Getenv("IAM_ROLE_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID, bucketName, iamRoleID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupSnapshotExportBucketExists(resourceName, &snapshotExportBucket),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
				),
			},
		},
	})
}

func TestAccBackupRSBackupSnapshotExportBucket_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_cloud_backup_snapshot_export_bucket.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName   = os.Getenv("AWS_S3_BUCKET")
		iamRoleID    = os.Getenv("IAM_ROLE_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID, bucketName, iamRoleID),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasBackupSnapshotExportBucketImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasBackupSnapshotExportBucketExists(resourceName string, snapshotExportBucket *matlas.CloudProviderSnapshotExportBucket) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		response, _, err := conn.CloudProviderSnapshotExportBuckets.Get(context.Background(), ids["project_id"], ids["id"])
		if err == nil {
			*snapshotExportBucket = *response
			return nil
		}

		return fmt.Errorf("snapshot export bucket (%s) does not exist", ids["id"])
	}
}

func testAccCheckMongoDBAtlasBackupSnapshotExportBucketDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_export_bucket" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		snapshotExportBucket, _, err := conn.CloudProviderSnapshotExportBuckets.Get(context.Background(), ids["project_id"], ids["id"])
		if err == nil && snapshotExportBucket != nil {
			return fmt.Errorf("snapshot export bucket (%s) still exists", ids["id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasBackupSnapshotExportBucketImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["id"]), nil
	}
}

func testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID, bucketName, iamRoleID string) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
			project_id     = "%[1]s"
    	  	iam_role_id    = "%[3]s"
       		bucket_name    = "%[2]s"
       		cloud_provider = "AWS"
    }
	`, projectID, bucketName, iamRoleID)
}
