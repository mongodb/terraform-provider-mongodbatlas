package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs_basic(t *testing.T) {
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
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists("mongodbatlas_cloud_provider_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJobs_withPagination(t *testing.T) {
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
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithPagination(projectID, clusterName, description, retentionInDays, 1, 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists("mongodbatlas_cloud_provider_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}

		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id    = mongodbatlas_cloud_provider_snapshot.test.project_id
			cluster_name  = mongodbatlas_cloud_provider_snapshot.test.cluster_name
			snapshot_id   = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
			delivery_type_config {
				download = true
			}
		}

		data "mongodbatlas_cloud_provider_snapshot_restore_jobs" "test" {
			project_id   = mongodbatlas_cloud_provider_snapshot_restore_job.test.project_id
			cluster_name = mongodbatlas_cloud_provider_snapshot_restore_job.test.cluster_name
		}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobsConfigWithPagination(projectID, clusterName, description, retentionInDays string, pageNum, itemPage int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}

		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id    = mongodbatlas_cloud_provider_snapshot.test.project_id
			cluster_name  = mongodbatlas_cloud_provider_snapshot.test.cluster_name
			snapshot_id   = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
			delivery_type = {
				download = true
			}
		}

		data "mongodbatlas_cloud_provider_snapshot_restore_jobs" "test" {
			project_id   = mongodbatlas_cloud_provider_snapshot_restore_job.test.project_id
			cluster_name = mongodbatlas_cloud_provider_snapshot_restore_job.test.cluster_name
			page_num = %d
			items_per_page = %d
		}
	`, projectID, clusterName, description, retentionInDays, pageNum, itemPage)
}
