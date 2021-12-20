package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasCloudBackupSchedule_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleConfigNoPolicies(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
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
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "0"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleNewPoliciesConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(0),
					ReferenceMinuteOfHour: pointy.Int64(0),
					RestoreWindowDays:     pointy.Int64(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "3"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudBackupSchedule_onepolicy(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
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
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "4"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleOnePolicyConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(0),
					ReferenceMinuteOfHour: pointy.Int64(0),
					RestoreWindowDays:     pointy.Int64(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudBackupScheduleImport_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(projectID, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
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
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.frequency_interval", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.0.retention_value", "3"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.0.retention_value", "4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceMongoDBAtlasCloudBackupSchedule_azure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleBasicConfig(projectID, clusterName, &matlas.PolicyItem{
					FrequencyInterval: 1,
					RetentionUnit:     "days",
					RetentionValue:    1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleBasicConfig(projectID, clusterName, &matlas.PolicyItem{
					FrequencyInterval: 2,
					RetentionUnit:     "days",
					RetentionValue:    3,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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
			return fmt.Errorf("cloud Provider Snapshot Schedule (%s) does not exist: %s", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccCheckMongoDBAtlasCloudBackupScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_cloud_backup_schedule" {
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
			return fmt.Errorf("cloud Provider Snapshot Schedule (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccMongoDBAtlasCloudBackupScheduleConfigNoPolicies(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d
		}
	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 2
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 3
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 4
			}
		}
	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleOnePolicyConfig(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %d
			reference_minute_of_hour = %d
			restore_window_days      = %d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
		}
	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleNewPoliciesConfig(projectID, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			policy_item_daily {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 4
			}
			policy_item_weekly {
				frequency_interval = 4
				retention_unit     = "weeks"
				retention_value    = 2
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 3
			}
		}
	`, projectID, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleBasicConfig(projectID, clusterName string, policy *matlas.PolicyItem) string {
	return fmt.Sprintf(`
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = %[1]q
  name         = %[2]q

  // Provider Settings "block"
  provider_name               = "AZURE"
  provider_region_name        = "US_EAST_2"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud provider snapshots
}

resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
  project_id   = mongodbatlas_cluster.my_cluster.project_id
  cluster_name = mongodbatlas_cluster.my_cluster.name

  policy_item_hourly {
    frequency_interval = %[3]d
    retention_unit     = %[4]q
    retention_value    = %[5]d
  }
}
	`, projectID, clusterName, policy.FrequencyInterval, policy.RetentionUnit, policy.RetentionValue)
}
