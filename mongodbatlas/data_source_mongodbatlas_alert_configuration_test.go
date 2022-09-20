package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccDataSourceMongoDBAtlaAlertConfiguration_basic(t *testing.T) {
	var (
		alert          = &matlas.AlertConfiguration{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAlertConfiguration(projectID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlaAlertConfiguration_withThreshold(t *testing.T) {
	var (
		alert          = &matlas.AlertConfiguration{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAlertConfigurationConfigWithThreshold(projectID, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccDataSourceMongoDBAtlaAlertConfiguration_withPagerDuty(t *testing.T) {
	SkipTestExtCred(t) // Will skip because requires external credentials aka api key
	var (
		alert          = &matlas.AlertConfiguration{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		projectID      = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		serviceKey     = os.Getenv("PAGER_DUTY_SERVICE_KEY")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDSMongoDBAtlasAlertConfigurationConfigWithPagerDuty(projectID, serviceKey, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func testAccDSMongoDBAtlasAlertConfiguration(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = true

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold_config {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = 99.0
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
			alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"
		}
	`, projectID)
}

func testAccDSMongoDBAtlasAlertConfigurationConfigWithThreshold(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
			enabled    = "%t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = "REPLICA_SET_NAME"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			threshold_config {
				operator    = "LESS_THAN"
				units       = "HOURS"
				threshold   = %f
			}
		}

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
			alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"
		}
	`, projectID, enabled, threshold)
}

func testAccDSMongoDBAtlasAlertConfigurationConfigWithPagerDuty(projectID, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_alert_configuration" "test" {
  project_id = %[1]q
  event_type = "NO_PRIMARY"
  enabled    = "%[3]t"

  notification {
    type_name    = "PAGER_DUTY"
    service_key  = %[2]q
    delay_min    = 0
  }
}

data "mongodbatlas_alert_configuration" "test" {
  project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
  alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"
}
	`, projectID, serviceKey, enabled)
}
