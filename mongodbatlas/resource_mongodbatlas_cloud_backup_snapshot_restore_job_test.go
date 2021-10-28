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

func TestAccResourceMongoDBAtlasCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	var (
		cloudBackupSnapshotRestoreJob = matlas.CloudProviderSnapshotRestoreJob{}
		resourceName                  = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
		projectID                     = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName                   = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description                   = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays               = "1"
		targetClusterName             = clusterName
		targetGroupID                 = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobExists(resourceName, &cloudBackupSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobAttributes(&cloudBackupSnapshotRestoreJob, "automated"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_cluster_name", targetClusterName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_project_id", targetGroupID),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobExists(resourceName, &cloudBackupSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobAttributes(&cloudBackupSnapshotRestoreJob, "download"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.download", "true"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudBackupSnapshotRestoreJob_importBasic(t *testing.T) {
	var (
		resourceName      = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
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
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days", "snapshot_id"},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudBackupSnapshotRestoreJobWithPointTime_basic(t *testing.T) {
	SkipTest(t)
	var (
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName     = acctest.RandomWithPrefix("test-acc")
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "1"
		targetGroupID   = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		timeUtc         = int64(1)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigPointInTime(projectID, clusterName, description, retentionInDays, targetGroupID, timeUtc),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobExists(resourceName string, cloudBackupSnapshotRestoreJob *matlas.CloudProviderSnapshotRestoreJob) resource.TestCheckFunc {
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

		log.Printf("[DEBUG] cloudBackupSnapshotRestoreJob ID: %s", rs.Primary.Attributes["snapshot_restore_job_id"])

		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
			JobID:       ids["snapshot_restore_job_id"],
		}

		if snapshotRes, _, err := conn.CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters); err == nil {
			*cloudBackupSnapshotRestoreJob = *snapshotRes
			return nil
		}

		return fmt.Errorf("cloudBackupSnapshotRestoreJob (%s) does not exist", rs.Primary.Attributes["snapshot_restore_job_id"])
	}
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobAttributes(cloudBackupSnapshotRestoreJob *matlas.CloudProviderSnapshotRestoreJob, deliveryType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if cloudBackupSnapshotRestoreJob.DeliveryType != deliveryType {
			return fmt.Errorf("bad cloudBackupSnapshotRestoreJob deliveryType: %s", cloudBackupSnapshotRestoreJob.DeliveryType)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_restore_job" {
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
			return fmt.Errorf("cloudBackupSnapshotRestoreJob (%s) still exists", rs.Primary.Attributes["snapshot_restore_job_id"])
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found:: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_restore_job_id"]), nil
	}
}

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigAutomated(projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  disk_size_gb = 5

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id     = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config   {
    automated           = true
    target_cluster_name = %[5]q
    target_project_id   = %[6]q
  }
}
	`, projectID, clusterName, description, retentionInDays, targetClusterName, targetGroupID)
}

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigDownload(projectID, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  disk_size_gb = 5

  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true   // enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config {
    download = true
  }
}
	`, projectID, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigPointInTime(projectID, clusterName, description, retentionInDays, targetGroupID string, pointTimeUTC int64) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "target_cluster" {
  project_id   = %[1]q
  name         = "cluster-target"
  disk_size_gb = 10

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true
}

resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q
  disk_size_gb = 10

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  provider_backup_enabled     = true   // enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[3]q
  retention_in_days = %[4]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config {
    point_in_time       = true
    target_cluster_name = mongodbatlas_cluster.target_cluster.name
    target_project_id   = %[5]q
    oplog_ts            = %[6]d
    oplog_inc           = 300
  }
}
	`, projectID, clusterName, description, retentionInDays, targetGroupID, pointTimeUTC)
}
