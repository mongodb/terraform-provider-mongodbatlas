package cloudbackupsnapshotrestorejob_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	resourceName   = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
	dataSourceName = "data.mongodbatlas_cloud_backup_snapshot_restore_job.test"
)

func TestAccCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccCloudBackupSnapshotRestoreJob_basicDownload(t *testing.T) {
	var (
		projectID       = acc.ProjectIDExecution(t)
		clusterName     = acc.RandomClusterName()
		description     = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays = "1"
		useSnapshotID   = true
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDownload(projectID, clusterName, description, retentionInDays, useSnapshotID),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.download", "true"),
				),
			},
			{
				Config:      configDownload(projectID, clusterName, description, retentionInDays, !useSnapshotID),
				ExpectError: regexp.MustCompile("SNAPSHOT_NOT_FOUND"),
			},
		},
	})
}

func TestAccCloudBackupSnapshotRestoreJobWithPointTime_basic(t *testing.T) {
	var (
		orgID             = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectID         = acc.ProjectIDExecution(t)
		targetProjectName = acc.RandomProjectName()
		clusterName       = acc.RandomClusterName()
		description       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays   = "1"
		timeUtc           = int64(1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configPointInTime(orgID, projectID, clusterName, description, retentionInDays, targetProjectName, timeUtc),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		snapshotsDataSourceName           = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.test"
		snapshotsDataSourcePaginationName = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.pagination"
		projectID                         = acc.ProjectIDExecution(tb)
		clusterName                       = acc.RandomClusterName()
		targetClusterName                 = acc.RandomClusterName()
		description                       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays                   = "1"
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, clusterName, description, retentionInDays, targetClusterName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.automated", "true"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_cluster_name", targetClusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "cluster_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "snapshot_id"),
					resource.TestCheckResourceAttrSet(snapshotsDataSourceName, "results.#"),
					resource.TestCheckResourceAttrSet(snapshotsDataSourcePaginationName, "results.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"retention_in_days", "snapshot_id"},
			},
		},
	}
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		if ids["snapshot_restore_job_id"] == "" {
			return fmt.Errorf("no ID is set")
		}
		log.Printf("[DEBUG] cloudBackupSnapshotRestoreJob ID: %s", rs.Primary.Attributes["snapshot_restore_job_id"])
		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
			JobID:       ids["snapshot_restore_job_id"],
		}
		if _, _, err := acc.Conn().CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters); err == nil {
			return nil
		}
		return fmt.Errorf("cloudBackupSnapshotRestoreJob (%s) does not exist", rs.Primary.Attributes["snapshot_restore_job_id"])
	}
}

func checkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_snapshot_restore_job" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		requestParameters := &matlas.SnapshotReqPathParameters{
			GroupID:     ids["project_id"],
			ClusterName: ids["cluster_name"],
			JobID:       ids["snapshot_restore_job_id"],
		}
		res, _, _ := acc.Conn().CloudProviderSnapshotRestoreJobs.Get(context.Background(), requestParameters)
		if res != nil {
			return fmt.Errorf("cloudBackupSnapshotRestoreJob (%s) still exists", rs.Primary.Attributes["snapshot_restore_job_id"])
		}
	}
	return nil
}

func importStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found:: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["cluster_name"], rs.Primary.Attributes["snapshot_restore_job_id"]), nil
	}
}

func configBasic(projectID, clusterName, description, retentionInDays, targetClusterName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = %[1]q
			name         = %[2]q
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_1"
			provider_instance_size_name = "M10"
			cloud_backup                = true
		}

		resource "mongodbatlas_cluster" "targer_cluster" {
			project_id   = %[1]q
			name         = %[5]q
			
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
	`, projectID, clusterName, description, retentionInDays, targetClusterName)
}

func configDownload(projectID, clusterName, description, retentionInDays string, useSnapshotID bool) string {
	var snapshotIDField string
	if useSnapshotID {
		snapshotIDField = `snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id`
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = %[1]q
			name         = %[2]q
			
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
			%[5]s

			delivery_type_config {
				download = true
			}
		}
	`, projectID, clusterName, description, retentionInDays, snapshotIDField)
}

func configPointInTime(orgID, projectID, clusterName, description, retentionInDays, targetProjectName string, pointTimeUTC int64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "target_project" {
			org_id = %[1]q
			name   = %[6]q
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
			project_id   = %[2]q
			name         = %[3]q
			disk_size_gb = 10

			// Provider Settings "block"
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

			delivery_type_config {
				point_in_time       = true
				target_cluster_name = mongodbatlas_cluster.target_cluster.name
				target_project_id   = mongodbatlas_cluster.target_cluster.project_id
				oplog_ts            = %[7]d
				oplog_inc           = 300
			}
		}
	`, orgID, projectID, clusterName, description, retentionInDays, targetProjectName, pointTimeUTC)
}
