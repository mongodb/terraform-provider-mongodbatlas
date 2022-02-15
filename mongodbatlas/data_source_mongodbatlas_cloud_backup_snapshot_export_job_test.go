package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudBackupSnapshotExportJob_basic(t *testing.T) {
	SkipTestExtCred(t)
	var (
		snapshotExportJob matlas.CloudProviderSnapshotExportJob
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		bucketName        = os.Getenv("AWS_S3_BUCKET")
		iamRoleID         = os.Getenv("IAM_ROLE_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportJobConfig(projectID, iamRoleID, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasBackupSnapshotExportJobExists("mongodbatlas_cloud_backup_snapshot_export_job.test", &snapshotExportJob),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_job.test", "iam_role_id"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_job.test", "bucket_name"),
					resource.TestCheckResourceAttrSet("data.mongodbatlas_cloud_backup_snapshot_export_job.test", "cloud_provider"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudBackupSnapshotExportJobConfig(projectID, iamRoleID, bucketName string) string {
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


resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = "%[1]s"
  name         = "MyCluster"
  disk_size_gb = 1

  //Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true // enable cloud backup snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = "%[1]s"
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = "myDescription"
  retention_in_days = 1
}


resource "mongodbatlas_cloud_backup_snapshot_export_job" "myjob" {
  project_id   = "%[1]s"
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
`, projectID, iamRoleID, bucketName)
}
