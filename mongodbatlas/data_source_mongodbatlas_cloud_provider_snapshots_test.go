package mongodbatlas

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshots_basic(t *testing.T) {
	projectID := "5cf5a45a9ccf6400e60981b6"
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	description := fmt.Sprintf("My description in %s", clusterName)
	retentionInDays := "1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "cluster_name"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "description"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot.test", "retention_in_days"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotsConfig(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5
			
			//Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
		}
		
		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}
		
		data "mongodbatlas_cloud_provider_snapshots" "test" {
			project_id   = mongodbatlas_cloud_provider_snapshot.test.project_id
			cluster_name = mongodbatlas_cloud_provider_snapshot.test.cluster_name
		}
	`, projectID, clusterName, description, retentionInDays)
}
