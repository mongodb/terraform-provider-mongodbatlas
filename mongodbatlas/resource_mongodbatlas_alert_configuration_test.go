package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	matlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAccConfigRSAlertConfiguration_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(orgID, projectName, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(orgID, projectName, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_Notifications(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(orgID, projectName, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(orgID, projectName, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_WithMatchers(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(orgID, projectName, true, false, true,
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
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(orgID, projectName, false, true, false,
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

func TestAccConfigRSAlertConfiguration_whitMetricUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(orgID, projectName, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(orgID, projectName, false, 89.7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_whitThresholdUpdated(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(orgID, projectName, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(orgID, projectName, false, 3),
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

func TestAccConfigRSAlertConfiguration_whitoutRoles(t *testing.T) {
	var (
		alert        = &matlas.AlertConfiguration{}
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigWithoutRoles(orgID, projectName, true, 99.0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_importBasic(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfig(orgID, projectName, true),
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

func TestAccConfigRSAlertConfiguration_importConfigNotifications(t *testing.T) {
	var (
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		resourceName = "mongodbatlas_alert_configuration.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckBasic(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationConfigNotifications(orgID, projectName, true, true, false),
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

func TestAccConfigRSAlertConfiguration_DataDog(t *testing.T) {
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
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
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

func TestAccConfigRSAlertConfiguration_PagerDuty(t *testing.T) {
	SkipTestExtCred(t) // Will skip because requires external credentials aka api key
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		serviceKey   = os.Getenv("PAGER_DUTY_SERVICE_KEY")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationPagerDutyConfig(projectID, serviceKey, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_OpsGenie(t *testing.T) {
	SkipTestExtCred(t) // Will skip because requires external credentials aka api key
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey       = os.Getenv("OPS_GENIE_API_KEY")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationOpsGenieConfig(projectID, apiKey, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_VictorOps(t *testing.T) {
	SkipTestExtCred(t) // Will skip because requires external credentials aka api key
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		projectID    = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		apiKey       = os.Getenv("VICTOR_OPS_API_KEY")
		alert        = &matlas.AlertConfiguration{}
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasAlertConfigurationVictorOpsConfig(projectID, apiKey, true),
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
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

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

func testAccMongoDBAtlasAlertConfigurationConfig(orgID, projectName string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}

resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "OUTSIDE_METRIC_THRESHOLD"
  enabled    = "%[3]t"

  notification {
    type_name     = "GROUP"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = false
    email_enabled = true
    roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"]
  }

  notification {
    type_name     = "ORG"
    interval_min  = 5
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
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
	`, orgID, projectName, enabled)
}

func testAccMongoDBAtlasAlertConfigurationConfigNotifications(orgID, projectName string, enabled, smsEnabled, emailEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "NO_PRIMARY"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[4]t
				email_enabled = %[5]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}

			notification {
				type_name     = "ORG"
				interval_min  = 5
				delay_min     = 1
				sms_enabled   = %[4]t
				email_enabled = %[5]t
			}
		}
	`, orgID, projectName, enabled, smsEnabled, emailEnabled)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(orgID, projectName string, enabled, smsEnabled, emailEnabled bool, m1, m2 matlas.Matcher) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "HOST_DOWN"
			enabled    = "%[3]t"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[4]t
				email_enabled = %[5]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = %[6]q
				operator   = %[7]q
				value      = %[8]q
			}
			matcher {
				field_name = %[9]q
				operator   = %[10]q
				value      = %[11]q
			}
		}
	`, orgID, projectName, enabled, smsEnabled, emailEnabled,
		m1.FieldName, m1.Operator, m1.Value,
		m2.FieldName, m2.Operator, m2.Value)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithMetrictUpdated(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%[3]t"

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

			metric_threshold_config {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %[4]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, orgID, projectName, enabled, threshold)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithoutRoles(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "OUTSIDE_METRIC_THRESHOLD"
			enabled    = "%[3]t"

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

			metric_threshold_config {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %[4]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, orgID, projectName, enabled, threshold)
}

func testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(orgID, projectName string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
			enabled    = "%[3]t"

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
				threshold   = %[4]f
			}
		}
	`, orgID, projectName, enabled, threshold)
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
    field_name = "REPLICA_SET_NAME"
    operator   = "EQUALS"
    value      = "SECONDARY"
  }

  threshold_config {
    operator    = "LESS_THAN"
    threshold   = 72
    units       = "HOURS"
  }
}
	`, projectID, enabled, dataDogAPIKey, dataDogRegion)
}

func testAccMongoDBAtlasAlertConfigurationPagerDutyConfig(projectID, serviceKey string, enabled bool) string {
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
	`, projectID, serviceKey, enabled)
}

func testAccMongoDBAtlasAlertConfigurationOpsGenieConfig(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_alert_configuration" "test" {
  project_id = %[1]q
  event_type = "NO_PRIMARY"
  enabled    = "%[3]t"

  notification {
    type_name          = "OPS_GENIE"
    ops_genie_api_key  = %[2]q
    ops_genie_region   = "US"
    delay_min          = 0
  }
}
	`, projectID, apiKey, enabled)
}

func testAccMongoDBAtlasAlertConfigurationVictorOpsConfig(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_alert_configuration" "test" {
  project_id = %[1]q
  event_type = "NO_PRIMARY"
  enabled    = "%[3]t"

  notification {
    type_name              = "VICTOR_OPS"
    victor_ops_api_key     = %[2]q
    victor_ops_routing_key = "testing"
    delay_min              = 0
  }
}
	`, projectID, apiKey, enabled)
}
