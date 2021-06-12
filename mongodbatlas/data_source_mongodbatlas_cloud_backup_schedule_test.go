package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccdataSourceMongoDBAtlasCloudBackupSchedule_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMongoDBAtlasCloudBackupScheduleConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
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
			disk_size_gb = 5

			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
			provider_disk_iops          = 1000
		}

		resource "mongodbatlas_cloud_backup_schedule" "test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d
		}
	 
		 data "mongodbatlas_cloud_backup_schedule" "test" {
			project_id = mongodbatlas_cloud_backup_schedule.test.project_id
			cluster_name = mongodbatlas_cloud_backup_schedule.test.cluster_name
		 }	 

	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}
