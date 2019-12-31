package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	matlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAccResourceMongoDBAtlasAlertConfiguration_basic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var alert = &matlas.AlertConfiguration{}

	resourceName := "mongodbatlas_alert_configuration.test"

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
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var alert = &matlas.AlertConfiguration{}

	resourceName := "mongodbatlas_alert_configuration.test"

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
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var alert = &matlas.AlertConfiguration{}

	resourceName := "mongodbatlas_alert_configuration.test"

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
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	var alert = &matlas.AlertConfiguration{}

	resourceName := "mongodbatlas_alert_configuration.test"

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
func TestAccResourceMongoDBAtlasAlertConfiguration_importBasic(t *testing.T) {
	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")

	resourceName := "mongodbatlas_alert_configuration.test"

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

		alertResp, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), ids["project_id"], ids["id"])
		if err != nil {
			return fmt.Errorf("Alert Configuration(%s) does not exist", ids["id"])
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

		alert, _, err := conn.AlertConfigurations.GetAnAlert(context.Background(), ids["project_id"], ids["id"])
		if alert != nil {
			return fmt.Errorf("project Alert Configuration(%s) still exists %s", ids["id"], err)
		}
	}
	return nil
}

func testAccCheckMongoDBAtlasAlertConfigurationImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s-%s", ids["project_id"], ids["id"]), nil
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
