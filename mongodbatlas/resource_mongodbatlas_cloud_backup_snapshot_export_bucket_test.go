package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasBackupSnapshotExportBucket_basic(t *testing.T) {
	var (
		snapshotExportBackup matlas.CloudProviderSnapshotExportBucket
		resourceName         = "mongodbatlas_cloud_backup_snapshot_export_bucket.test"
		projectID            = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupSnapshotExportBucketExists(resourceName, &snapshotExportBackup),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasBackupSnapshotExportBucket_importBasic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_snapshot_export_bucket.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID),
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

func testAccCheckMongoDBAtlasBackupSnapshotExportBucketExists(resourceName string, snapshotExportBackup *matlas.CloudProviderSnapshotExportBucket) resource.TestCheckFunc {
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
			*snapshotExportBackup = *response
			return nil
		}

		return fmt.Errorf("snapshot export backup (%s) does not exist", ids["id"])
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
			return fmt.Errorf("snapshot export backup (%s) still exists", ids["id"])
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

func testAccMongoDBAtlasBackupSnapshotExportBucketConfig(projectID string) string {
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
	`, projectID)
}
