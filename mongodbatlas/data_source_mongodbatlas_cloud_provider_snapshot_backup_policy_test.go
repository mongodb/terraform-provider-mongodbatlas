package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceMongoDBAtlasCloudProviderSnapshotBackupPolicy_basic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasDataSourceCloudProviderSnapshotBackupPolicyConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyExists("mongodbatlas_cloud_provider_snapshot_backup_policy.test"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot_backup_policy.test", "project_id"),
					resource.TestCheckResourceAttrSet("mongodbatlas_cloud_provider_snapshot_backup_policy.test", "cluster_name"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasDataSourceCloudProviderSnapshotBackupPolicyConfig(projectID, clusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5

			//Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true //enable cloud provider snapshots
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot_backup_policy" "test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = 3
			reference_minute_of_hour = 45
			restore_window_days      = 4


			policies {
				id = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.id

				policy_item {
					id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.0.id
					frequency_interval = 1
					frequency_type     = "hourly"
					retention_unit     = "days"
					retention_value    = 1
				}
				policy_item {
					id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.1.id
					frequency_interval = 1
					frequency_type     = "daily"
					retention_unit     = "days"
					retention_value    = 2
				}
				policy_item {
					id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.2.id
					frequency_interval = 4
					frequency_type     = "weekly"
					retention_unit     = "weeks"
					retention_value    = 3
				}
				policy_item {
					id                 = mongodbatlas_cluster.my_cluster.snapshot_backup_policy.0.policies.0.policy_item.3.id
					frequency_interval = 5
					frequency_type     = "monthly"
					retention_unit     = "months"
					retention_value    = 4
				}
			}
		}

		data "mongodbatlas_cloud_provider_snapshot_backup_policy" "test" {
			project_id   = mongodbatlas_cloud_provider_snapshot_backup_policy.test.project_id
			cluster_name = mongodbatlas_cloud_provider_snapshot_backup_policy.test.cluster_name
		}
`, projectID, clusterName)
}
