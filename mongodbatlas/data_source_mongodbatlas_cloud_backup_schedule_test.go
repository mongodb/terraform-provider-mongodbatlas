package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccdataSourceMongoDBAtlasCloudBackupSchedule_basic(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName    = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasCloudBackupScheduleConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(datasourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(datasourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_hourly.#", "0"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_monthly.#", "0"),
				),
			},
		},
	})
}

func TestAccdataSourceMongoDBAtlasCloudBackupSchedule_withOnePolicy(t *testing.T) {
	var (
		datasourceName = "data.mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName    = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasCloudBackupScheduleWithPoliciesConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(datasourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(datasourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(datasourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(datasourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(datasourceName, "policy_item_monthly.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceMongoDBAtlasCloudBackupScheduleConfig(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d
		}
	 
		 data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id = mongodbatlas_cloud_backup_schedule.schedule_test.project_id
			cluster_name = mongodbatlas_cloud_backup_schedule.schedule_test.cluster_name
		 }	 

	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccDataSourceMongoDBAtlasCloudBackupScheduleWithPoliciesConfig(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = "%s"
			name         = "%s"
			
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d

			policy_item_hourly {
				frequency_interval = 2
				retention_unit     = "days"
				retention_value    = 1
			}
		}
	 
		 data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id = mongodbatlas_cloud_backup_schedule.schedule_test.project_id
			cluster_name = mongodbatlas_cloud_backup_schedule.schedule_test.cluster_name
		 }	 

	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}
