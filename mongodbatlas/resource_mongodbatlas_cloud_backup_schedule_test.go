package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mwielbut/pointy"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccBackupRSCloudBackupSchedule_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_cloud_backup_schedule.schedule_test"
		dataSourceName = "data.mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acctest.RandomWithPrefix("test-acc")
		clusterName    = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleConfigNoPolicies(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_hour_of_day"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_minute_of_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restore_window_days"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_hourly.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_daily.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_weekly.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policy_item_monthly.#")),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleNewPoliciesConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(0),
					ReferenceMinuteOfHour: pointy.Int64(0),
					RestoreWindowDays:     pointy.Int64(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
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
					resource.TestCheckResourceAttr(dataSourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_hour_of_day"),
					resource.TestCheckResourceAttrSet(dataSourceName, "reference_minute_of_hour"),
					resource.TestCheckResourceAttrSet(dataSourceName, "restore_window_days"),
				),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleAdvancedPoliciesConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(0),
					ReferenceMinuteOfHour: pointy.Int64(0),
					RestoreWindowDays:     pointy.Int64(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "0"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "0"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.#", "2"),
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
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.frequency_interval", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.retention_unit", "weeks"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_weekly.1.retention_value", "5"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.frequency_interval", "6"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.retention_unit", "months"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_monthly.1.retention_value", "4"),
				),
			},
		},
	})
}

func TestAccBackupRSCloudBackupSchedule_export(t *testing.T) {
	SkipTestExtCred(t)
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
		policyName   = acctest.RandomWithPrefix("test-acc")
		roleName     = acctest.RandomWithPrefix("test-acc")
		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		region       = os.Getenv("AWS_REGION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleExportPoliciesConfig(orgID, projectName, clusterName, policyName, roleName, awsAccessKey, awsSecretKey, region),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "auto_export_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "20"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "5"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "4"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_daily.0.retention_value", "4"),
				),
			},
		},
	})
}
func TestAccBackupRSCloudBackupSchedule_onepolicy(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
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
				Config: testAccMongoDBAtlasCloudBackupScheduleOnePolicyConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(0),
					ReferenceMinuteOfHour: pointy.Int64(0),
					RestoreWindowDays:     pointy.Int64(7),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
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

func TestAccBackupRSCloudBackupSchedule_copySettings(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleCopySettingsConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(1),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "reference_hour_of_day", "3"),
					resource.TestCheckResourceAttr(resourceName, "reference_minute_of_hour", "45"),
					resource.TestCheckResourceAttr(resourceName, "restore_window_days", "1"),
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
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.cloud_provider", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.region_name", "US_EAST_1"),
					resource.TestCheckResourceAttr(resourceName, "copy_settings.0.should_copy_oplogs", "true"),
				),
			},
		},
	})
}
func TestAccBackupRSCloudBackupScheduleImport_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(orgID, projectName, clusterName, &matlas.CloudProviderSnapshotBackupPolicy{
					ReferenceHourOfDay:    pointy.Int64(3),
					ReferenceMinuteOfHour: pointy.Int64(45),
					RestoreWindowDays:     pointy.Int64(4),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
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

func TestAccBackupRSCloudBackupSchedule_azure(t *testing.T) {
	var (
		resourceName = "mongodbatlas_cloud_backup_schedule.schedule_test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		clusterName  = fmt.Sprintf("test-acc-%s", acctest.RandString(10))
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasCloudBackupScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleBasicConfig(orgID, projectName, clusterName, &matlas.PolicyItem{
					FrequencyInterval: 1,
					RetentionUnit:     "days",
					RetentionValue:    1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "cluster_name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.frequency_interval", "1"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_unit", "days"),
					resource.TestCheckResourceAttr(resourceName, "policy_item_hourly.0.retention_value", "1")),
			},
			{
				Config: testAccMongoDBAtlasCloudBackupScheduleBasicConfig(orgID, projectName, clusterName, &matlas.PolicyItem{
					FrequencyInterval: 2,
					RetentionUnit:     "days",
					RetentionValue:    3,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasCloudBackupScheduleExists(resourceName),
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

func testAccMongoDBAtlasCloudBackupScheduleConfigNoPolicies(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d
		}

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name
		 }	
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleDefaultConfig(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d

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

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name
		 }	
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleCopySettingsConfig(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "backup_project" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_cluster" "my_cluster" {
			project_id   = mongodbatlas_project.backup_project.id
			name         = %[3]q
			
			cluster_type = "REPLICASET"
            replication_specs {
            num_shards = 1
            regions_config {
              region_name     = "US_EAST_2"
              electable_nodes = 3
              priority        = 7
              read_only_nodes = 0
              }
            }
			// Provider Settings "block"
			provider_name               = "AWS"
			provider_region_name        = "US_EAST_2"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
			pit_enabled = true // enable point in time restore. you cannot copy oplogs when pit is not enabled.
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d

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
			copy_settings {
				cloud_provider = "AWS"
				frequencies = ["HOURLY",
							"DAILY",
							"WEEKLY",
							"MONTHLY",
							"ON_DEMAND"]
				region_name = "US_EAST_1"
				replication_spec_id = mongodbatlas_cluster.my_cluster.replication_specs.*.id[0]
				should_copy_oplogs = true
			  }
		}
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleOnePolicyConfig(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d

			policy_item_hourly {
				frequency_interval = 1
				retention_unit     = "days"
				retention_value    = 1
			}
		}
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleNewPoliciesConfig(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name

			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d

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

		data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name
		 }	
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleBasicConfig(orgID, projectName, clusterName string, policy *matlas.PolicyItem) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "backup_project" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q

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
    frequency_interval = %[4]d
    retention_unit     = %[5]q
    retention_value    = %[6]d
  }
}

data "mongodbatlas_cloud_backup_schedule" "schedule_test" {
	project_id   = mongodbatlas_cluster.my_cluster.project_id
	cluster_name = mongodbatlas_cluster.my_cluster.name
}	
	`, orgID, projectName, clusterName, policy.FrequencyInterval, policy.RetentionUnit, policy.RetentionValue)
}

func testAccMongoDBAtlasCloudBackupScheduleAdvancedPoliciesConfig(orgID, projectName, clusterName string, p *matlas.CloudProviderSnapshotBackupPolicy) string {
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
			provider_region_name        = "EU_CENTRAL_1"
			provider_instance_size_name = "M10"
			cloud_backup     = true //enable cloud provider snapshots
		}

		resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
			project_id   = mongodbatlas_cluster.my_cluster.project_id
			cluster_name = mongodbatlas_cluster.my_cluster.name
			auto_export_enabled = false
			reference_hour_of_day    = %[4]d
			reference_minute_of_hour = %[5]d
			restore_window_days      = %[6]d

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
			policy_item_weekly {
				frequency_interval = 5
				retention_unit     = "weeks"
				retention_value    = 5
			}
			policy_item_monthly {
				frequency_interval = 5
				retention_unit     = "months"
				retention_value    = 3
			}
			policy_item_monthly {
				frequency_interval = 6
				retention_unit     = "months"
				retention_value    = 4
			}
		}
	`, orgID, projectName, clusterName, *p.ReferenceHourOfDay, *p.ReferenceMinuteOfHour, *p.RestoreWindowDays)
}

func testAccMongoDBAtlasCloudBackupScheduleExportPoliciesConfig(orgID, projectName, clusterName, policyName, roleName, awsAccessKey, awsSecretKey, region string) string {
	return fmt.Sprintf(`

resource "mongodbatlas_project" "backup_project" {
	name   = %[2]q
	org_id = %[1]q
}
locals {
	mongodbatlas_project_id = mongodbatlas_project.backup_project.id
}

provider "aws" {
	region     = %[8]q
	access_key = %[6]q
	secret_key = %[7]q
}

resource "mongodbatlas_cluster" "my_cluster" {
  project_id   = mongodbatlas_project.backup_project.id
  name         = %[3]q

  // Provider Settings "block"
  provider_name               = "AWS"
  provider_region_name        = "US_WEST_2"
  provider_instance_size_name = "M10"
  cloud_backup                = true //enable cloud provider snapshots
  depends_on = ["mongodbatlas_cloud_backup_snapshot_export_bucket.test"]
}

resource "mongodbatlas_cloud_backup_schedule" "schedule_test" {
  project_id               = mongodbatlas_cluster.my_cluster.project_id
  cluster_name             = mongodbatlas_cluster.my_cluster.name
  auto_export_enabled      = true
  reference_hour_of_day    = 20
  reference_minute_of_hour = "05"
  restore_window_days      = 4

  policy_item_daily {
	frequency_interval = 1
	retention_unit     = "days"
	retention_value    = 4
  }
  export {
		export_bucket_id = mongodbatlas_cloud_backup_snapshot_export_bucket.test.export_bucket_id
		frequency_type   = "daily"
  }
}

resource "aws_s3_bucket" "backup" {
	bucket = "${local.mongodbatlas_project_id}-s3-mongodb-backups"
	force_destroy = true
    object_lock_configuration {
      object_lock_enabled = "Enabled"
    }
}

resource "mongodbatlas_cloud_provider_access_setup" "setup_only" {
  project_id    = mongodbatlas_project.backup_project.id
  provider_name = "AWS"
}

resource "mongodbatlas_cloud_provider_access_authorization" "auth_role" {
  project_id = mongodbatlas_project.backup_project.id
  role_id    = mongodbatlas_cloud_provider_access_setup.setup_only.role_id

	aws {
	  iam_assumed_role_arn = aws_iam_role.test_role.arn
	}
}

resource "mongodbatlas_cloud_backup_snapshot_export_bucket" "test" {
  project_id = mongodbatlas_project.backup_project.id

  iam_role_id    = mongodbatlas_cloud_provider_access_authorization.auth_role.role_id
  bucket_name    = aws_s3_bucket.backup.bucket
  cloud_provider = "AWS"
}

resource "aws_iam_role_policy" "test_policy" {
	name = mongodbatlas_project.backup_project.id
	role = aws_iam_role.test_role.id

	policy = <<-EOF
	{
	  "Version": "2012-10-17",
	  "Statement": [
		{
		  "Effect": "Allow",
		  "Action": "*",
		  "Resource": "*"
		}
	  ]
	}
	EOF
}

resource "aws_iam_role" "test_role" {
  name = %[5]q

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_aws_account_arn}"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${mongodbatlas_cloud_provider_access_setup.setup_only.aws_config.0.atlas_assumed_role_external_id}"
        }
      }
    }
  ]
}
EOF

}
	`, orgID, projectName, clusterName, policyName, roleName, awsAccessKey, awsSecretKey, region)
}

func testAccCheckMongoDBAtlasCloudProviderSnapshotBackupPolicyImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["cluster_name"]), nil
	}
}
