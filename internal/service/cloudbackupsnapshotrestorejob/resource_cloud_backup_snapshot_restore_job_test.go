package cloudbackupsnapshotrestorejob_test

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

const (
	resourceName   = "mongodbatlas_cloud_backup_snapshot_restore_job.test"
	dataSourceName = "data.mongodbatlas_cloud_backup_snapshot_restore_job.test"
)

func clusterRequest() *acc.ClusterRequest {
	return &acc.ClusterRequest{
		CloudBackup: true,
		ReplicationSpecs: []acc.ReplicationSpecRequest{
			{Region: "US_WEST_2"},
		},
	}
}

func TestAccCloudBackupSnapshotRestoreJob_basic(t *testing.T) {
	resource.ParallelTest(t, *basicTestCase(t))
}

func TestAccCloudBackupSnapshotRestoreJob_basicDownload(t *testing.T) {
	var (
		clusterInfo         = acc.GetClusterInfo(t, clusterRequest())
		clusterName         = clusterInfo.Name
		description         = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays     = "1"
		useSnapshotID       = true
		clusterTerraformStr = clusterInfo.TerraformStr
		clusterResourceName = clusterInfo.ResourceName
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 acc.PreCheckBasicSleep(t, &clusterInfo, "", ""),
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configDownload(clusterTerraformStr, clusterResourceName, description, retentionInDays, useSnapshotID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.download", "true"),
				),
			},
			{
				Config:      configDownload(clusterTerraformStr, clusterResourceName, description, retentionInDays, !useSnapshotID),
				ExpectError: regexp.MustCompile("SNAPSHOT_NOT_FOUND"),
			},
		},
	})
}

func basicTestCase(tb testing.TB) *resource.TestCase {
	tb.Helper()

	var (
		snapshotsDataSourceName           = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.test"
		snapshotsDataSourcePaginationName = "data.mongodbatlas_cloud_backup_snapshot_restore_jobs.pagination"
		clusterInfo                       = acc.GetClusterInfo(tb, clusterRequest())
		clusterName                       = clusterInfo.Name
		description                       = fmt.Sprintf("My description in %s", clusterName)
		retentionInDays                   = "1"
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasicSleep(tb, &clusterInfo, "", ""); mig.PreCheckOldPreviewEnv(tb) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasic(clusterInfo.TerraformStr, clusterInfo.ResourceName, description, retentionInDays),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.automated", "true"),
					resource.TestCheckResourceAttr(resourceName, "delivery_type_config.0.target_cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "failed", "false"),
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
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]
		restoreJobID := ids["snapshot_restore_job_id"]
		if _, _, err := acc.ConnV2().CloudBackupsApi.GetBackupRestoreJob(context.Background(), projectID, clusterName, restoreJobID).Execute(); err == nil {
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
		projectID := ids["project_id"]
		clusterName := ids["cluster_name"]
		restoreJobID := ids["snapshot_restore_job_id"]
		res, _, _ := acc.ConnV2().CloudBackupsApi.GetBackupRestoreJob(context.Background(), projectID, clusterName, restoreJobID).Execute()
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

func configBasic(terraformStr, clusterResourceName, description, retentionInDays string) string {
	return fmt.Sprintf(`
		%[1]s
		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			project_id        = %[2]s.project_id
			cluster_name      = %[2]s.name
			description       = %[3]q
			retention_in_days = %[4]q
			depends_on = [%[2]s]
		}

		resource "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
			project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
			cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
			snapshot_id     = mongodbatlas_cloud_backup_snapshot.test.id

			delivery_type_config   {
				automated           = true
				target_project_id   = mongodbatlas_cloud_backup_snapshot.test.project_id
				target_cluster_name = mongodbatlas_cloud_backup_snapshot.test.cluster_name
			}
			depends_on = [mongodbatlas_cloud_backup_snapshot.test]			
		}

		data "mongodbatlas_cloud_backup_snapshot_restore_job" "test" {
			project_id      = mongodbatlas_cloud_backup_snapshot.test.project_id
			cluster_name    = mongodbatlas_cloud_backup_snapshot.test.cluster_name
			snapshot_restore_job_id       = mongodbatlas_cloud_backup_snapshot_restore_job.test.snapshot_restore_job_id
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
	`, terraformStr, clusterResourceName, description, retentionInDays)
}

func configDownload(terraformStr, clusterResourceName, description, retentionInDays string, useSnapshotID bool) string {
	var snapshotIDField string
	if useSnapshotID {
		snapshotIDField = `snapshot_id  = mongodbatlas_cloud_backup_snapshot.test.id`
	}

	return fmt.Sprintf(`
		%[1]s
		resource "mongodbatlas_cloud_backup_snapshot" "test" {
			project_id        = %[2]s.project_id
			cluster_name      = %[2]s.name
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
	`, terraformStr, clusterResourceName, description, retentionInDays, snapshotIDField)
}
