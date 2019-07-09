package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJob_basic(t *testing.T) {
	var cloudProviderSnapshotRestoreJob = matlas.CloudProviderSnapshotRestoreJob{}

	resourceName := "mongodbatlas_cloud_provider_snapshot_restore_job.test"

	projectID := "5cf5a45a9ccf6400e60981b6"
	clusterName := "MyCluster"
	description := "myDescription"
	retentionInDays := "1"
	targetClusterName := "MyCluster"
	targetGroupID := "5cf5a45a9ccf6400e60981b6"

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

	projectID := "5cf5a45a9ccf6400e60981b6"
	clusterName := "MyCluster"
	description := "myDescription"
	retentionInDays := "1"
	targetClusterName := "MyCluster"
	targetGroupID := "5cf5a45a9ccf6400e60981b6"

	resource.Test(t, resource.TestCase{
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
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshotRestoreJob ID: %s", rs.Primary.ID)

		requestParameters := &matlas.SnapshotReqPathParameters{
			JobID:       rs.Primary.ID,
			GroupID:     rs.Primary.Attributes["project_id"],
			ClusterName: rs.Primary.Attributes["cluster_name"],
		}

		if snapshotRes, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters); err == nil {
			*cloudProviderSnapshotRestoreJob = *snapshotRes
			return nil
		}
		return fmt.Errorf("cloudProviderSnapshotRestoreJob (%s) does not exist", rs.Primary.ID)
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
			JobID:       rs.Primary.ID,
		}

		snapshotReq, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters)
		if err != nil {
			return fmt.Errorf("error getting cloudProviderSnapshotRestoreJob Information: %s", err)
		}

		if snapshotReq.DeliveryType == "download" {
			_, err := conn.CloudProviderSnapshotRestoreJobs.Delete(context.Background(), requestParameters)
			if err != nil {
				return fmt.Errorf("cloudProviderSnapshotRestoreJob (%s) still exists", rs.Primary.ID)
			}
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
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.ID), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID string) string {
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
				automated           = true
				target_cluster_name = "%s"
				target_project_id     = "%s"
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays string) string {
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
