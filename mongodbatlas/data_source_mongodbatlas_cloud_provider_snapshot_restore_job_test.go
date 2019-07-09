package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJob_basic(t *testing.T) {
	var cloudProviderSnapshotRestoreJob matlas.CloudProviderSnapshotRestoreJob
	var cloudProviderSnapshot matlas.CloudProviderSnapshot

	resourceName := "data.mongodbatlas_cloud_provider_snapshot_restore_job.test"

	projectID := "5cf5a45a9ccf6400e60981b6"
	clusterName := "MyCluster"
	description := "myDescription"
	retentionInDays := "1"
	// targetClusterName := "MyCluster"
	// targetGroupID := "5cf5a45a9ccf6400e60981b6"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists("mongodbatlas_cloud_provider_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigWithDS(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id          = "%s"
			cluster_name      = "%s"
			description       = "%s"
			retention_in_days = %s
		}
		
		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id      = "${mongodbatlas_cloud_provider_snapshot.test.project_id}"
			cluster_name  = "${mongodbatlas_cloud_provider_snapshot.test.cluster_name}"
			snapshot_id   = "${mongodbatlas_cloud_provider_snapshot.test.id}"
			delivery_type = {
				download = true
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigWithDS(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		%s		
		
		data "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id     = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.project_id}"
			cluster_name = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.cluster_name}"
			job_id       = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.id}"
		}
		`, testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(projectID, clusterName, description, retentionInDays))
}
