package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshotBackupPolicy_basic(t *testing.T) {
	resourceName := "mongodbatlas_cloud_provider_snapshot_backup_policy.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotBackupPolicyConfig(projectID, clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudProviderSnapshotBackupPolicy_importBasic(t *testing.T) {

	resourceName := "mongodbatlas_cloud_provider_snapshot_backup_policy.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotBackupPolicyConfig(projectID, clusterName),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{""},
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]

		schedule, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
		if err != nil || schedule == nil {
			return fmt.Errorf("Cloud Provider Snapshot Schedule (%s) does not exist: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_snapshot_backup_policy" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]

		schedule, _, err := conn.CloudProviderSnapshotBackupPolicies.Get(context.Background(), projectID, clusterName)
		if schedule != nil || err == nil {
			return fmt.Errorf("Cloud Provider Snapshot Schedule (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		ids := decodeStateID(rs.Primary.ID)
		return fmt.Sprintf("%s-%s", ids["project_id"], ids["cluster_name"]), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotBackupPolicyConfig(projectID, clusterName string) string {
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
	`, projectID, clusterName)
}
