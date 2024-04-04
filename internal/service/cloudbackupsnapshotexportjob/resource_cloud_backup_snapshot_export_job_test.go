package cloudbackupsnapshotexportjob_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

var (
	resourceName         = "mongodbatlas_cloud_backup_snapshot_export_job.test"
	dataSourceName       = "data.mongodbatlas_cloud_backup_snapshot_export_job.test"
	dataSourcePluralName = "data.mongodbatlas_cloud_backup_snapshot_export_jobs.test"
)

func TestAccBackupSnapshotExportJob_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()
	acc.SkipTestForCI(tb) // needs AWS IAM role and S3 bucket

	var (
		projectID  = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName = os.Getenv("AWS_S3_BUCKET")
		iamRoleID  = os.Getenv("IAM_ROLE_ID")
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(tb); acc.PreCheckS3Bucket(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, bucketName, iamRoleID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(resourceName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttrSet(resourceName, "iam_role_id"),

					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(dataSourceName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttrSet(dataSourceName, "iam_role_id"),

					resource.TestCheckResourceAttr(dataSourcePluralName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourcePluralName, "bucket_name", "example-bucket"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "cloud_provider", "AWS"),
					resource.TestCheckResourceAttrSet(dataSourcePluralName, "iam_role_id"),
					resource.TestCheckResourceAttr(dataSourcePluralName, "results.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().CloudProviderSnapshotExportJobs.Get(context.Background(), ids["project_id"], ids["cluster_name"], ids["export_job_id"])
		if err == nil {
			return nil
		}
		return fmt.Errorf("snapshot export job (%s) does not exist", ids["export_job_id"])
	}
}

func checkDestroy(state *terraform.State) error {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_export_job" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		snapshotExportBucket, _, err := acc.Conn().CloudProviderSnapshotExportJobs.Get(context.Background(), ids["project_id"], ids["cluster_name"], ids["export_job_id"])
		if err == nil && snapshotExportBucket != nil {
			return fmt.Errorf("snapshot export job (%s) still exists", ids["export_job_id"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s-%s", ids["project_id"], ids["cluster_name"], ids["export_job_id"]), nil
	}
}

func configBasic(projectID, bucketName, iamRoleID string) string {
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
}

data "mongodbatlas_cloud_backup_snapshot_export_job" "test" {
  project_id = "%[1]s"
  cluster_name = mongodbatlas_cluster.my_cluster.name
  export_job_id = mongodbatlas_cloud_backup_snapshot_export_job.myjob.export_job_id
}

data "mongodbatlas_cloud_backup_snapshot_export_jobs" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot_export_bucket.test.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name
}
`, projectID, bucketName, iamRoleID)
}
