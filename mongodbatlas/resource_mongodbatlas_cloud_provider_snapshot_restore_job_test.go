package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJob_basic(t *testing.T) {
	var cloudProviderSnapshotRestoreJob = matlas.CloudProviderSnapshotRestoreJob{}

	resourceName := "mongodbatlas_cloud_provider_snapshot_restore_job.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	description := fmt.Sprintf("My description in %s", clusterName)
	retentionInDays := "1"
	targetClusterName := clusterName
	targetGroupID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobAttributes(&cloudProviderSnapshotRestoreJob, "automated"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type.target_cluster_name", targetClusterName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type.target_project_id", targetGroupID),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobAttributes(&cloudProviderSnapshotRestoreJob, "download"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type.download", "true"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJob_importBasic(t *testing.T) {

	resourceName := "mongodbatlas_cloud_provider_snapshot_restore_job.test"

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	clusterName := fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	description := fmt.Sprintf("My description in %s", clusterName)
	retentionInDays := "1"
	targetClusterName := clusterName
	targetGroupID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName string, cloudProviderSnapshotRestoreJob *matlas.CloudProviderSnapshotRestoreJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.Attributes["snapshot_restore_job_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshotRestoreJob ID: %s", rs.Primary.Attributes["snapshot_restore_job_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			JobID:       rs.Primary.Attributes["snapshot_restore_job_id"],
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		if snapshotRes, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters); err == nil {
			*cloudProviderSnapshotRestoreJob = *snapshotRes
			return nil
		}
		return fmt.Errorf("cloudProviderSnapshotRestoreJob (%s) does not exist", rs.Primary.Attributes["snapshot_restore_job_id"])
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobAttributes(cloudProviderSnapshotRestoreJob *matlas.CloudProviderSnapshotRestoreJob, deliveryType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudProviderSnapshotRestoreJob.DeliveryType != deliveryType {
			return fmt.Errorf("bad cloudProviderSnapshotRestoreJob deliveryType: %s", cloudProviderSnapshotRestoreJob.DeliveryType)
		}
		return nil
	}
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_snapshot_restore_job" {
			continue
		}

		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
			JobID:       rs.Primary.Attributes["snapshot_restore_job_id"],
		}

		res, _, _ := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters)

		if res != nil {
			return fmt.Errorf("cloudProviderSnapshotRestoreJob (%s) still exists", rs.Primary.Attributes["snapshot_restore_job_id"])
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_restore_job_id"]), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5

		//Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_WEST_2"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true   // enable cloud provider snapshots
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}

		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
			cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
			snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
			delivery_type   = {
				automated           = true
				target_cluster_name = "%s"
				target_project_id   = "%s"
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			disk_size_gb = 5
	
		//Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_WEST_2"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true   // enable cloud provider snapshots
			provider_disk_iops          = 100
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cloud_provider_snapshot" "test" {
			project_id        = mongodbatlas_cluster.my_cluster.project_id
			cluster_name      = mongodbatlas_cluster.my_cluster.name
			description       = "%s"
			retention_in_days = %s
		}
		
		resource "mongodbatlas_cloud_provider_snapshot_restore_job" "test" {
			project_id      = mongodbatlas_cloud_provider_snapshot.test.project_id
			cluster_name    = mongodbatlas_cloud_provider_snapshot.test.cluster_name
			snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.snapshot_id
			delivery_type   = {
				download = true
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays)
}
