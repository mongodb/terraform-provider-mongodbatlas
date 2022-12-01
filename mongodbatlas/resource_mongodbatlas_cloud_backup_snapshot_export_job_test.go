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

func TestAccBackupRSBackupSnapshotExportJob_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		snapshotExportJob matlas.CloudProviderSnapshotExportJob
		resourceName      = "mongodbatlas_cloud_backup_snapshot_export_job.test"
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName        = os.Getenv("AWS_S3_BUCKET")
		iamRoleID         = os.Getenv("IAM_ROLE_ID")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportJobConfig(projectID, bucketName, iamRoleID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupSnapshotExportJobExists(resourceName, &snapshotExportJob),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
				),
			},
		},
	})
}

func TestAccBackupRSBackupSnapshotExportJob_importBasic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_cloud_backup_snapshot_export_job.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName   = os.Getenv("AWS_S3_BUCKET")
		iamRoleID    = os.Getenv("IAM_ROLE_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasBackupSnapshotExportJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasBackupSnapshotExportJobConfig(projectID, bucketName, iamRoleID),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasBackupSnapshotExportJobImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasBackupSnapshotExportJobExists(resourceName string, snapshotExportJob *matlas.CloudProviderSnapshotExportJob) resource.TestCheckFunc {
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

		response, _, err := conn.CloudProviderSnapshotExportJobs.Get(context.Background(), ids["project_id"], ids["cluster_name"], ids["export_job_id"])
		if err == nil {
			*snapshotExportJob = *response
			return nil
		}

		return fmt.Errorf("snapshot export job (%s) does not exist", ids["export_job_id"])
	}
}

func testAccCheckMongoDBAtlasBackupSnapshotExportJobDestroy(state *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_export_job" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		snapshotExportBucket, _, err := conn.CloudProviderSnapshotExportJobs.Get(context.Background(), ids["project_id"], ids["cluster_name"], ids["export_job_id"])
		if err == nil && snapshotExportBucket != nil {
			return fmt.Errorf("snapshot export job (%s) still exists", ids["export_job_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasBackupSnapshotExportJobImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["cluster_name"], ids["export_job_id"]), nil
	}
}

func testAccMongoDBAtlasBackupSnapshotExportJobConfig(projectID, bucketName, iamRoleID string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = var.project_id
  name         = "MyCluster"
  disk_size_gb = 1
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true // enable cloud backup snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = var.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = "myDescription"
  retention_in_days = 1
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id     = "%[1]s"
  iam_role_id    = "%[3]s"
  bucket_name    = "%[2]s"
  cloud_provider = "AWS"
}

resource "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id   = var.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name
  snapshot_id = mongodbatlas_cloud_backup_snapshot.test.snapshot_id
  export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id

  custom_data {
    key   = "exported by"
    value = "myName"
  }
}`, projectID, bucketName, iamRoleID)
}
