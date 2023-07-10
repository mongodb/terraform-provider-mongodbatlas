package mongodbatlas

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupRSCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	var (
		cloudBackupSnapshotRestoreJob     = matlas.CloudProviderSnapshotRestoreJob{}
		resourceName                      = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
		snapshotsDataSourceName           = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.test"
		snapshotsDataSourcePaginationName = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.pagination"
		dataSourceName                    = "data.mongodbatlas_cloud_backup_snapshot_restore_job.test"
		orgID                             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                       = acctest.RandomWithPrefix("test-snapshot-acc")
		targetProjectName                 = acctest.RandomWithPrefix("test-acc")
		clusterName                       = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description                       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays                   = "1"
		targetClusterName                 = fmt.Sprintf("test-acc-target-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigAutomated(orgID, projectName, clusterName, description, retentionInDays, targetProjectName, targetClusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobExists(resourceName, &cloudBackupSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobAttributes(&cloudBackupSnapshotRestoreJob, "automated"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_cluster_name", targetClusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(snapshotsDataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(snapshotsDataSourcePaginationName, "results.#"),
				),
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

func TestAccBackupRSCloudBackupSnapshotRestoreJob_basicDownload(t *testing.T) {
	var (
		cloudBackupSnapshotRestoreJob = matlas.CloudProviderSnapshotRestoreJob{}
		resourceName                  = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
		orgID                         = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName                   = acctest.RandomWithPrefix("test-acc")
		clusterName                   = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		description                   = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays               = "1"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigDownload(orgID, projectName, clusterName, description, retentionInDays),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobExists(resourceName, &cloudBackupSnapshotRestoreJob),
					testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobAttributes(&cloudBackupSnapshotRestoreJob, "download"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.download", "true"),
				),
			},
		},
	})
}

func TestAccBackupRSCloudBackupSnapshotRestoreJobWithPointTime_basic(t *testing.T) {
	SkipTestForCI(t)
	var (
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName       = acctest.RandomWithPrefix("test-acc")
		targetProjectName = acctest.RandomWithPrefix("test-acc-target")
		clusterName       = acctest.RandomWithPrefix("test-acc")
		description       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays   = "1"
		timeUtc           = int64(1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupSnapshotRestoreJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigPointInTime(orgID, projectName, clusterName, description, retentionInDays, targetProjectName, timeUtc),
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

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigAutomated(orgID, projectName, clusterName, description, retentionInDays, targetProjectName, targetClusterName string) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "backup_project" {
	name   = %[2]q
	org_id = %[1]q
}


resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q
  
  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true
}

resource "mongodbatlas_cluster" "targer_cluster" {
	project_id   = mongodbatlas_project.backup_project.id
	name         = %[7]q
	
	// Provider Settings "block"
	provider_name               = "AWS"
	provider_region_name        = "US_EAST_1"
	provider_instance_size_name = "M10"
	cloud_backup                = true
  }

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[4]q
  retention_in_days = %[5]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id     = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config   {
    automated           = true
    target_cluster_name = mongodbatlas_cluster.targer_cluster.name
    target_project_id   = mongodbatlas_cluster.targer_cluster.project_id
  }
}

data "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
	project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
	cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
	job_id       = mongodbatlas_cloud_backup_snapshot_restore_job.test.id  
}

data "mongodbatlas_cloud_backup_snapshot_restore_jobs" "test" {
	project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
	cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
}

data "mongodbatlas_cloud_backup_snapshot_restore_jobs" "pagination" {
	project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
	cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
	page_num = 1
	items_per_page = 5
}

	`, orgID, projectName, clusterName, description, retentionInDays, targetProjectName, targetClusterName)
}

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigDownload(orgID, projectName, clusterName, description, retentionInDays string) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "backup_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q
  
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true   // enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_snapshot" "test" {
  project_id        = mongodbatlas_cluster.my_cluster.project_id
  cluster_name      = mongodbatlas_cluster.my_cluster.name
  description       = %[4]q
  retention_in_days = %[5]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config {
    download = true
  }
}
	`, orgID, projectName, clusterName, description, retentionInDays)
}

func testAccMongoDBAtlasCloudBackupSnapshotRestoreJobConfigPointInTime(orgID, projectName, clusterName, description, retentionInDays, targetProjectName string, pointTimeUTC int64) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "backup_project" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_project" "target_project" {
	name   = %[6]q
	org_id = %[1]q
}

resource "mongodbatlas_cluster" "target_cluster" {
  project_id   = mongodbatlas_project.target_project.id
  name         = "cluster-target"
  disk_size_gb = 10

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_EAST_1"
  provider_instance_size_name = "M10"
  cloud_backup                = true
}

resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q
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
  description       = %[4]q
  retention_in_days = %[5]q
}

resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
  project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
  cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
  snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id

  delivery_type_config {
    point_in_time       = true
    target_cluster_name = mongodbatlas_cluster.target_cluster.name
    target_project_id   = mongodbatlas_cluster.target_cluster.project_id
    oplog_ts            = %[7]d
    oplog_inc           = 300
  }
}
	`, orgID, projectName, clusterName, description, retentionInDays, targetProjectName, pointTimeUTC)
}
