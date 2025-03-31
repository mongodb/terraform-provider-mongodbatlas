package alertconfiguration_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/service/alertconfiguration"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName         = "mongodbatlas_alert_configuration.test"
	dataSourceName       = "data.mongodbatlas_alert_configuration.test"
	dataSourcePluralName = "data.mongodbatlas_alert_configurations.test"
)

func TestAccConfigRSAlertConfiguration_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configBasicRS(projectID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			{
				Config: configBasicRS(projectID, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withEmptyMetricThresholdConfig(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyMetricThresholdConfig(projectID, true),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withEmptyMatcherMetricThresholdConfig(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyMatcherMetricThresholdConfig(projectID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
		},
	})
}
func TestAccConfigRSAlertConfiguration_withNotifications(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithNotifications(projectID, true, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configWithNotifications(projectID, false, false, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withMatchers(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithMatchers(projectID, true, false, true,
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "EQUALS",
						"value":     "SECONDARY",
					},
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "CONTAINS",
						"value":     "MONGOS",
					}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configWithMatchers(projectID, false, true, false,
					map[string]interface{}{
						"fieldName": "TYPE_NAME",
						"operator":  "NOT_EQUALS",
						"value":     "SECONDARY",
					},
					map[string]interface{}{
						"fieldName": "HOSTNAME",
						"operator":  "EQUALS",
						"value":     "PRIMARY",
					}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withMetricUpdated(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithMetricUpdated(projectID, true, 99.0),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configWithMetricUpdated(projectID, false, 89.7),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withThresholdUpdated(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithThresholdUpdated(projectID, true, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configWithThresholdUpdated(projectID, false, 3),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id", "updated", "matcher.0.field_name"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withoutRoles(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithoutRoles(projectID, true, 99.0),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withoutOptionalAttributes(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithEmptyOptionalAttributes(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_importIncorrectId(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configBasicRS(projectID, true),
			},
			{
				ResourceName:  resourceName,
				ImportState:   true,
				ImportStateId: "incorrect_id_without_project_id_and_dash",
				ExpectError:   regexp.MustCompile("import format error"),
			},
		},
	})
}

// dummy keys used for credential values in third party notifications
const dummy32CharKey = "11111111111111111111111111111111"
const dummy36CharKey = "11111111-1111-1111-1111-111111111111"

func TestAccConfigRSAlertConfiguration_updatePagerDutyWithNotifierId(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		serviceKey = dummy32CharKey
		notifierID = "651dd9336afac13e1c112222"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDutyNotifierID(projectID, notifierID, 10, &serviceKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.0.delay_min", "10"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.service_key", serviceKey),
				),
			},
			{
				Config: configWithPagerDutyNotifierID(projectID, notifierID, 15, nil),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.0.delay_min", "15"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withDataDog(t *testing.T) {
	resource.Test(t, *datadogTestCase(t)) // not run in parallel so acc and mig tests don't interfere
}

func datadogTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	var (
		projectID = acc.ProjectIDExecution(t)
		ddAPIKey  = dummy32CharKey
		ddRegion  = "US"

		groupNotificationMap = map[string]string{
			"type_name":    "GROUP",
			"interval_min": "5",
			"delay_min":    "0",
		}

		ddNotificationMap = map[string]string{
			"type_name":       "DATADOG",
			"interval_min":    "5",
			"delay_min":       "0",
			"datadog_api_key": ddAPIKey,
			"datadog_region":  ddRegion,
		}

		ddNotificationUpdatedMap = map[string]string{
			"type_name":       "DATADOG",
			"interval_min":    "6",
			"delay_min":       "0",
			"datadog_api_key": ddAPIKey,
			"datadog_region":  ddRegion,
		}
	)

	return &resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithDataDog(projectID, ddAPIKey, ddRegion, true, ddNotificationMap, groupNotificationMap),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", groupNotificationMap),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", ddNotificationMap),
				),
			},
			{
				Config: configWithDataDog(projectID, ddAPIKey, ddRegion, true, ddNotificationUpdatedMap, groupNotificationMap),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", groupNotificationMap),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", ddNotificationUpdatedMap),
				),
			},
		},
	}
}
func TestAccConfigRSAlertConfiguration_withPagerDuty(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		serviceKey = dummy32CharKey
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDuty(projectID, serviceKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateProjectIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// service key is not returned by api in import operation
				// integration_id is not returned during Create
				ImportStateVerifyIgnore: []string{"updated", "notification.0.service_key", "notification.0.integration_id"},
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withEmailToPagerDuty(t *testing.T) {
	var (
		projectID  = acc.ProjectIDExecution(t)
		serviceKey = dummy32CharKey
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithEmail(projectID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				Config: configWithPagerDuty(projectID, serviceKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateProjectIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// service_key is not returned by api in import operation
				// integration_id is not returned during Create
				ImportStateVerifyIgnore: []string{"updated", "notification.0.service_key", "notification.0.integration_id"},
			},
		},
	})
}

func TestAccConfigAlertConfiguration_PagerDutyUsingIntegrationID(t *testing.T) {
	// create a new project as it need to ensure no third party integration has already been created
	var (
		orgID       = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName = acc.RandomProjectName()
		serviceKey  = dummy32CharKey
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDutyIntegrationID(orgID, projectName, serviceKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "notification.0.integration_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "notification.0.integration_id"),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withOpsGenie(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		apiKey    = dummy36CharKey
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithOpsGenie(projectID, apiKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func TestAccConfigRSAlertConfiguration_withVictorOps(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		apiKey    = dummy36CharKey
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithVictorOps(projectID, apiKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
				),
			},
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.ConnV2().AlertConfigurationsApi.GetAlertConfiguration(context.Background(), ids[alertconfiguration.EncodedIDKeyProjectID], ids[alertconfiguration.EncodedIDKeyAlertID]).Execute()
		if err != nil {
			return fmt.Errorf("the Alert Configuration(%s) does not exist", ids[alertconfiguration.EncodedIDKeyAlertID])
		}
		return nil
	}
}

func checkDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_alert_configuration" {
				continue
			}
			ids := conversion.DecodeStateID(rs.Primary.ID)
			alert, _, err := acc.ConnV2().AlertConfigurationsApi.GetAlertConfiguration(context.Background(), ids[alertconfiguration.EncodedIDKeyProjectID], ids[alertconfiguration.EncodedIDKeyAlertID]).Execute()
			if alert != nil {
				return fmt.Errorf("the Project Alert Configuration(%s) still exists %s", ids[alertconfiguration.EncodedIDKeyAlertID], err)
			}
		}
		return nil
	}
}

func importStateProjectIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s-%s", rs.Primary.Attributes["project_id"], rs.Primary.Attributes["alert_configuration_id"]), nil
	}
}

func configBasicRS(projectID string, enabled bool) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_alert_configuration" "test" {
		project_id = %[1]q
		enabled    = %[2]t
		event_type = "OUTSIDE_METRIC_THRESHOLD"

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
	`, projectID, enabled)
}

func configWithNotifications(projectID string, enabled, smsEnabled, emailEnabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			event_type = "NO_PRIMARY"
			enabled    = %[2]t

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

func configWithMatchers(projectID string, enabled, smsEnabled, emailEnabled bool, m1, m2 map[string]interface{}) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "HOST_DOWN"

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[3]t
				email_enabled = %[4]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}

			matcher {
				field_name = %[5]q
				operator   = %[6]q
				value      = %[7]q
			}
			matcher {
				field_name = %[8]q
				operator   = %[9]q
				value      = %[10]q
			}
		}
	`, projectID, enabled, smsEnabled, emailEnabled,
		m1["fieldName"], m1["operator"], m1["value"],
		m2["fieldName"], m2["operator"], m2["value"])
}

func configWithMetricUpdated(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "OUTSIDE_METRIC_THRESHOLD"

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
				threshold   = %[3]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled, threshold)
}

func configWithoutRoles(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "OUTSIDE_METRIC_THRESHOLD"

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
				threshold   = %[3]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled, threshold)
}

func configWithThresholdUpdated(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"

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
				threshold   = %[3]f
			}
		}
	`, projectID, enabled, threshold)
}

func configWithDataDog(projectID, dataDogAPIKey, dataDogRegion string, enabled bool, ddNotificationMap, groupNotificationMap map[string]string) string {
	ddNotificationBlock := fmt.Sprintf(`
	notification {
		type_name = %[1]q
		datadog_api_key = mongodbatlas_third_party_integration.atlas_datadog.api_key
		datadog_region = mongodbatlas_third_party_integration.atlas_datadog.region
		interval_min  = %[2]v
		delay_min     = %[3]v
	}
	`, ddNotificationMap["type_name"], ddNotificationMap["interval_min"], ddNotificationMap["delay_min"])

	groupNotificationBlock := fmt.Sprintf(`
	notification {
		type_name     = %[1]q
		interval_min  = %[2]v
		delay_min     = %[3]v
		sms_enabled   = false
		email_enabled = true
		roles         = ["GROUP_OWNER"]
	}
	`, groupNotificationMap["type_name"], groupNotificationMap["interval_min"], groupNotificationMap["delay_min"])

	return fmt.Sprintf(`
		resource "mongodbatlas_third_party_integration" "atlas_datadog" {
			project_id = %[1]q
			api_key    = %[2]q
			region     = %[3]q
			type = "DATADOG"
		}

		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_third_party_integration.atlas_datadog.project_id
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
			enabled    = %[4]t

			%[5]s

			%[6]s

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
	`, projectID, dataDogAPIKey, dataDogRegion, enabled, groupNotificationBlock, ddNotificationBlock)
}

func configWithPagerDuty(projectID, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[3]t
			event_type = "NO_PRIMARY"

			notification {
				type_name    = "PAGER_DUTY"
				service_key  = %[2]q
				delay_min    = 0
			}
		}
	`, projectID, serviceKey, enabled)
}

func configWithEmail(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "NO_PRIMARY"

			notification {
				type_name     = "EMAIL"
				interval_min  = 60
				email_address = "test@mongodbtest.com"
			}
		}
	`, projectID, enabled)
}

func configWithPagerDutyIntegrationID(orgID, projectName, serviceKey string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			org_id = %[1]q
			name   = %[2]q
		}
	
		resource "mongodbatlas_third_party_integration" "test" {
			project_id = mongodbatlas_project.test.id
			type = "PAGER_DUTY"
			service_key = %[3]q
		}

		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
			enabled    = true
			event_type = "USERS_WITHOUT_MULTI_FACTOR_AUTH"
		  
			notification {
				type_name     = "PAGER_DUTY"
				integration_id = mongodbatlas_third_party_integration.test.id
			}
		}

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = mongodbatlas_alert_configuration.test.project_id
			alert_configuration_id = mongodbatlas_alert_configuration.test.id
		}
	`, orgID, projectName, serviceKey)
}

func configWithPagerDutyNotifierID(projectID, notifierID string, delayMin int, serviceKey *string) string {
	var serviceKeyString string
	if serviceKey != nil {
		serviceKeyString = fmt.Sprintf(`service_key = %q`, *serviceKey)
	}
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = true
			event_type = "NO_PRIMARY"

			notification {
				type_name    = "PAGER_DUTY"
				notifier_id  = %[2]q
				%[3]s
				delay_min    = %[4]d
			}
		}
	`, projectID, notifierID, serviceKeyString, delayMin)
}

func configWithOpsGenie(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[3]t
			event_type = "NO_PRIMARY"

			notification {
				type_name          = "OPS_GENIE"
				ops_genie_api_key  = %[2]q
				ops_genie_region   = "US"
				delay_min          = 0
			}
		}
	`, projectID, apiKey, enabled)
}

func configWithVictorOps(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[3]t
			event_type = "NO_PRIMARY"

			notification {
				type_name              = "VICTOR_OPS"
				victor_ops_api_key     = %[2]q
				victor_ops_routing_key = "testing"
				delay_min              = 0
			}
		}
	`, projectID, apiKey, enabled)
}

func configWithEmptyMetricThresholdConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = %[2]t
			event_type = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"

			notification {
				type_name     = "GROUP"
				interval_min  = 60
				delay_min     = 0
				sms_enabled   = true
				email_enabled = false
			roles         = ["GROUP_OWNER"]
			}

			threshold_config {
				operator    = "LESS_THAN"
				threshold   = 72
				units       = "HOURS"
			}
		}
	`, projectID, enabled)
}

func configWithEmptyMatcherMetricThresholdConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_alert_configuration" "test" {
  project_id = %[1]q
  enabled    = %[2]t
  event_type = "CLUSTER_MONGOS_IS_MISSING"

  notification {
    type_name     = "GROUP"
    interval_min  = 60
    delay_min     = 0
    sms_enabled   = true
    email_enabled = false
	roles         = ["GROUP_OWNER"]
  }
}
	`, projectID, enabled)
}

// configWithEmptyOptionalAttributes does not define notification.delay_min, notification.sms_enabled, and metric_threshold_config.threshold.
func configWithEmptyOptionalAttributes(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			event_type = "OUTSIDE_METRIC_THRESHOLD"

			notification {
			  type_name     = "ORG"
			  interval_min  = 5
			  email_enabled   = true
			}

			metric_threshold_config {
			  metric_name = "ASSERT_REGULAR"
			  operator    = "LESS_THAN"
			  units       = "RAW"
			  mode        = "AVERAGE"
			}
		  }
	`, projectID)
}

func configWithEmptyOptionalBlocks(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			event_type = "NO_PRIMARY"
			enabled    = true

			notification {
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = true
				email_enabled = false
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}
		}
	`, projectID)
}
