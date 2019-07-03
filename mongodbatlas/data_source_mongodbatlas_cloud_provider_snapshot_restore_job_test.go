package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	matlas "github.com/mongodb-partners/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotRestoreJob_basic(t *testing.T) {
	var cloudProviderSnapshotRestoreJob matlas.CloudProviderSnapshotRestoreJob
	var cloudProviderSnapshot matlas.CloudProviderSnapshot

	resourceName := "data.mongodbatlas_cloud_provider_snapshot_restore_job.test"

	groupID := "5cf5a45a9ccf6400e60981b6"
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
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(groupID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotExists("mongodbatlas_cloud_provider_snapshot.test", &cloudProviderSnapshot),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "group_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigWithDS(groupID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					resource.TestCheckResourceAttrSet(resourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(resourceName, "group_id"),
					resource.TestCheckResourceAttrSet(resourceName, "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(groupID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			group_id          = "%s"
			cluster_name      = "%s"
			description       = "%s"
			retention_in_days = %s
		}
		
		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			group_id      = "${mongodbatlas_cloud_provider_snapshot.test.group_id}"
			cluster_name  = "${mongodbatlas_cloud_provider_snapshot.test.cluster_name}"
			snapshot_id   = "${mongodbatlas_cloud_provider_snapshot.test.id}"
			delivery_type = {
				download = true
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, groupID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigWithDS(groupID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		%s		
		
		data "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			group_id     = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.group_id}"
			cluster_name = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.cluster_name}"
			job_id       = "${mongodbatlas_cloud_provider_snapshot_restore_job.test.id}"
		}
		`, testAccMongoDBAtlasDataSourceCloudProviderSnapshotRestoreJobConfig(groupID, clusterName, description, retentionInDays))
}
