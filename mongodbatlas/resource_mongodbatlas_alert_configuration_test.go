package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasAlertConfiguration_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(projectID, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(projectID, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_Notifications(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(projectID, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(projectID, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_WithMatchers(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(projectID, true, false, true,
					matlas.Matcher{
						FieldName: "TYPE_NAME",
						Operator:  "EQUALS",
						Value:     "SECONDARY",
					},
					matlas.Matcher{
						FieldName: "TYPE_NAME",
						Operator:  "CONTAINS",
						Value:     "MONGOS",
					}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(projectID, false, true, false,
					matlas.Matcher{
						FieldName: "TYPE_NAME",
						Operator:  "NOT_EQUALS",
						Value:     "SECONDARY",
					},
					matlas.Matcher{
						FieldName: "HOSTNAME",
						Operator:  "EQUALS",
						Value:     "PRIMARY",
					}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_whitMetricUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(projectID, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(projectID, false, 89.7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_whitThresholdUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(projectID, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(projectID, false, 3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasAlertConfigurationImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "matcher.0.field_name"},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_whitoutRoles(t *testing.T) {
	var (
		alert        = &matlas.AlertConfiguration{}
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithoutRoles(projectID, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_importBasic(t *testing.T) {
	var (
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(projectID, true),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasAlertConfigurationImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_importConfigNotifications(t *testing.T) {
	var (
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(projectID, true, true, false),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasAlertConfigurationImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}

func TestAccResourceMongoDBAtlasAlertConfiguration_DataDog(t *testing.T) {
	SkipTestExtCred(t) // Will skip because requires external credentials aka api key
	SkipTest(t)        // Will force skip if enabled
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		ddAPIKey     = os.Getenv("DD_API_KEY")
		ddRegion     = "US"
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithDataDog(projectID, ddAPIKey, ddRegion, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName string, alert *matlas.AlertConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*matlas.Client)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		alertResp, _, err := conn.AlertConfigurations.GetAnAlertConfig(context.Background(), ids["project_id"], ids["id"])
		if err != nil {
			return fmt.Errorf("the Alert Configuration(%s) does not exist", ids["id"])
		}

		alert = alertResp

		return nil
	}
}

func testAccCheckMongoDBAtlasAlertConfigurationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*matlas.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_alert_configuration" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)

		alert, _, err := conn.AlertConfigurations.GetAnAlertConfig(context.Background(), ids["project_id"], ids["id"])
		if alert != nil {
			return fmt.Errorf("the Project Alert Configuration(%s) still exists %s", ids["id"], err)
		}
	}

	return nil
}

func testAccCheckMongoDBAtlasAlertConfigurationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["alert_configuration_id"]), nil
	}
}

func testAccMongoDBAtlasAlertConfigurationConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"]
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold = {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = 99.0
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled)
}

func testAccMongoDBAtlasAlertConfigurationConfigNotifications(projectID string, enabled, smsEnabled, emailEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%[1]s"
			event_type = "NO_PRIMARY"
			enabled    = "%[2]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[3]t
				email_enabled = %[4]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}

			notification {
				type_name     = "ORG"
				interval_min  = 5
				delay_min     = 1
				sms_enabled   = %[3]t
				email_enabled = %[4]t
			}
		}
	`, projectID, enabled, smsEnabled, emailEnabled)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(projectID string, enabled, smsEnabled, emailEnabled bool, m1, m2 matlas.Matcher) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "NO_PRIMARY"
			enabled    = "%t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %t
				email_enabled = %t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = "%s"
				operator   = "%s"
				value      = "%s"
			}
			matcher {
				field_name = "%s"
				operator   = "%s"
				value      = "%s"
			}
		}
	`, projectID, enabled, smsEnabled, emailEnabled,
		m1.FieldName, m1.Operator, m1.Value,
		m2.FieldName, m2.Operator, m2.Value)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold = {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled, threshold)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithoutRoles(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = "%s"
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%t"

			notification {
				type_name     = "EMAIL"
				email_address = "mongodbatlas.testing@gmail.com"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = false
			}

			matcher {
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			metric_threshold = {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled, threshold)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(projectID string, enabled bool, threshold float64) string {
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
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}

			threshold = {
				operator    = "LESS_THAN"
				units       = "HOURS"
				threshold   = %f
			}
		}
	`, projectID, enabled, threshold)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithDataDog(projectID, dataDogAPIKey, dataDogRegion string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_third_party_integration" "atlas_datadog" {
  project_id = "%[1]s"
  type = "DATADOG"
  api_key = "%[3]s"
  region = "%[4]s"
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = "%[1]s"
  event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
  enabled    = %t

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles         = ["GROUP_OWNER"]
  }

  notification {
    type_name = "DATADOG"
    datadog_api_key = mongodbatlas_third_party_integration.atlas_datadog.api_key
    datadog_region = mongodbatlas_third_party_integration.atlas_datadog.region
    interval_min  = 5
    delay_min     = 0
  }

  matcher {
    field_name = "HOSTNAME_AND_PORT"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  threshold = {
    operator    = "LESS_THAN"
    threshold   = 72
    units       = "HOURS"
  }
}
	`, projectID, enabled, dataDogAPIKey, dataDogRegion)
}
