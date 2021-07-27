package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJob_basic(t *testing.T) {
	var (
		cloudProviderSnapshotRestoreJob = matlas.CloudProviderSnapshotRestoreJob{}
		resourceName                    = "mongodbatlas_cloud_provider_snapshot_restore_job.test"
		projectID                       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName                     = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description                     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays                 = "1"
		targetClusterName               = clusterName
		targetGroupID                   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobAttributes(&cloudProviderSnapshotRestoreJob, "automated"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_cluster_name", targetClusterName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_project_id", targetGroupID),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName, &cloudProviderSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobAttributes(&cloudProviderSnapshotRestoreJob, "download"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.download", "true"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJob_importBasic(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_cloud_provider_snapshot_restore_job.test"
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName       = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays   = "1"
		targetClusterName = clusterName
		targetGroupID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days", "snapshot_id"},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudProviderSnapshotRestoreJobWithPointTime_basic(t *testing.T) {
	t.Skip()
	var (
		projectID         = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName       = acctest.RandomWithPrefix("test-acc")
		description       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays   = "1"
		targetClusterName = clusterName
		targetGroupID     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		timeUtc           = int64(1)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigPointInTime(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID, timeUtc),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotRestoreJobExists(resourceName string, cloudProviderSnapshotRestoreJob *matlas.CloudProviderSnapshotRestoreJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		if ids["snapshot_restore_job_id"] == "" {
			return fmt.Errorf("no ID is set")
		}

		log.Printf("[DEBUG] cloudProviderSnapshotRestoreJob ID: %s", rs.Primary.Attributes["snapshot_restore_job_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
			JobID:       ids["snapshot_restore_job_id"],
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
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_provider_snapshot_restore_job" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
			JobID:       ids["snapshot_restore_job_id"],
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
			return "", fmt.Errorf("not found:: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_restore_job_id"]), nil
	}
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 5

		// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true
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
			snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.id
			delivery_type_config   {
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

			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true   // enable cloud provider snapshots
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
			snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.id
			delivery_type_config   {
				download = true
			}
			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudProviderSnapshotRestoreJobConfigPointInTime(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID string, pointTimeUTC int64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "target_cluster" {
			project_id   = "%[1]s"
			name         = "cluster-target"
			disk_size_gb = 10

		// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true
			provider_encrypt_ebs_volume = false
		}

		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%[1]s"
			name         = "%[2]s"
			disk_size_gb = 10

		// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			provider_backup_enabled     = true   // enable cloud provider snapshots
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
			snapshot_id     = mongodbatlas_cloud_provider_snapshot.test.id

			delivery_type_config   {
				point_in_time       = true
				target_cluster_name = mongodbatlas_cluster.target_cluster.name
				target_project_id   = "%s"
				oplog_ts = %v
				oplog_inc = 300
			}

			depends_on = ["mongodbatlas_cloud_provider_snapshot.test"]
		}
	`, projectID, clusterName, description, retentionInDays, targetGroupID, pointTimeUTC)
}
