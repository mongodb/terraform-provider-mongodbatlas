package alertconfigurationapi_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

const (
	resourceName = "mongodbatlas_alert_configuration_api.test"
)

func TestAccAlertConfigurationAPI_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configBasic(projectID, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "severity_override"),
					resource.TestCheckResourceAttrSet(resourceName, "notifications.0.notifier_id"),
					resource.TestCheckResourceAttr(resourceName, "matchers.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "metric_threshold.metric_name"),
				),
			},
			{
				Config: configBasic(projectID, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.#", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "severity_override"),
					resource.TestCheckResourceAttrSet(resourceName, "notifications.0.notifier_id"),
					resource.TestCheckResourceAttr(resourceName, "matchers.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "metric_threshold.metric_name"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       importStateProjectIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated"},
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withEmptyMetricThresholdConfig(t *testing.T) {
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.#", "1"),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withEmptyMatcherMetricThresholdConfig(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.#", "1"),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withMatchers(t *testing.T) {
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
					map[string]any{
						"fieldName": "TYPE_NAME",
						"operator":  "EQUALS",
						"value":     "SECONDARY",
					},
					map[string]any{
						"fieldName": "TYPE_NAME",
						"operator":  "CONTAINS",
						"value":     "MONGOS",
					}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
			{
				Config: configWithMatchers(projectID, false, true, false,
					map[string]any{
						"fieldName": "TYPE_NAME",
						"operator":  "NOT_EQUALS",
						"value":     "SECONDARY",
					},
					map[string]any{
						"fieldName": "HOSTNAME",
						"operator":  "EQUALS",
						"value":     "PRIMARY",
					}),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withMetricUpdated(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
			{
				Config: configWithMetricUpdated(projectID, true, 89.7),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
		},
	})
}

// dummy keys used for credential values in third party notifications
const dummy32CharKey = "11111111111111111111111111111111"
const dummy36CharKey = "11111111-1111-1111-1111-111111111111"

func TestAccAlertConfigurationAPI_updatePagerDutyWithNotifierId(t *testing.T) {
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
				Config: configWithPagerDutyNotifierID(projectID, notifierID, serviceKey, 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.delay_min", "10"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.service_key", serviceKey),
				),
			},
			{
				Config: configWithPagerDutyNotifierID(projectID, notifierID, serviceKey, 15),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.delay_min", "15"),
					resource.TestCheckResourceAttr(resourceName, "notifications.0.service_key", serviceKey),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withDataDog(t *testing.T) {
	resource.Test(t, *datadogTestCase(t)) // not run in parallel so acc and mig tests don't interfere
}

func datadogTestCase(t *testing.T) *resource.TestCase {
	t.Helper()

	var (
		projectID = acc.ProjectIDExecution(t)
		ddAPIKey  = dummy32CharKey
		ddRegion  = "US"

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
				Config: configWithDataDog(projectID, ddAPIKey, ddRegion, true, ddNotificationMap),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notifications.*", ddNotificationMap),
				),
			},
			{
				Config: configWithDataDog(projectID, ddAPIKey, ddRegion, true, ddNotificationUpdatedMap),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notifications.*", ddNotificationUpdatedMap),
				),
			},
		},
	}
}

func TestAccAlertConfigurationAPI_withPagerDuty(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateProjectIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// service key is not returned by api in import operation
				// integration_id is not returned during Create
				ImportStateVerifyIgnore: []string{"updated", "notifications.0.service_key", "notifications.0.integration_id"},
			},
		},
	})
}

func TestAccAlertConfiguration_withEmailToPagerDuty(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
			{
				Config: configWithPagerDuty(projectID, serviceKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: importStateProjectIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// service_key is not returned by api in import operation
				// integration_id is not returned during Create
				ImportStateVerifyIgnore: []string{"updated", "notifications.0.service_key", "notifications.0.integration_id"},
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withOpsGenie(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withVictorOps(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "group_id", projectID),
				),
			},
		},
	})
}

func TestAccAlertConfigurationAPI_withSeverityOverride(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithSeverityOverride(projectID, conversion.StringPtr("ERROR")),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "severity_override", "ERROR"),
				),
			},
			// TODO: Should check for no attr once CLOUDP-353933 is fixed.
			// {
			// 	Config: configWithSeverityOverride(projectID, nil),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		checkExists(resourceName),
			// 		resource.TestCheckNoResourceAttr(resourceName, "severity_override"),
			// 	),
			// },
		},
	})
}

func checkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		id := rs.Primary.Attributes["id"]
		projectID := rs.Primary.Attributes["group_id"]
		if id == "" || projectID == "" {
			return fmt.Errorf("checkExists, attributes not found for: %s", resourceName)
		}
		if _, _, err := acc.ConnV2().AlertConfigurationsApi.GetAlertConfig(context.Background(), projectID, id).Execute(); err != nil {
			return fmt.Errorf("the Alert Configuration(%s) does not exist", id)
		}
		return nil
	}
}

func checkDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "mongodbatlas_alert_configuration_api" {
				continue
			}
			id := rs.Primary.Attributes["id"]
			projectID := rs.Primary.Attributes["group_id"]
			if id == "" || projectID == "" {
				return fmt.Errorf("checkDestroy, attributes not found for: %s", resourceName)
			}
			alert, _, err := acc.ConnV2().AlertConfigurationsApi.GetAlertConfig(context.Background(), projectID, id).Execute()
			if alert != nil {
				return fmt.Errorf("the Project Alert Configuration(%s) still exists %s", id, err)
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
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["group_id"], rs.Primary.ID), nil
	}
}

func configBasic(projectID string, enabled bool) string {
	return fmt.Sprintf(`
	resource "mongodbatlas_alert_configuration_api" "test" {
		group_id = %[1]q
		enabled  = %[2]t
		event_type_name = "OUTSIDE_METRIC_THRESHOLD"

		notifications = [{
			type_name     = "GROUP"
			interval_min  = 5
			delay_min     = 0
			sms_enabled   = false
			email_enabled = true
			roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER", "GROUP_DATA_ACCESS_ADMIN"]
		}, {
			type_name     = "ORG"
			interval_min  = 5
			delay_min     = 0
			sms_enabled   = true
			email_enabled = false
		}]

		matchers = [{
			field_name = "HOSTNAME_AND_PORT"
			operator   = "EQUALS"
			value      = "SECONDARY"
		}]

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

func configWithMatchers(projectID string, enabled, smsEnabled, emailEnabled bool, m1, m2 map[string]any) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = %[1]q
			enabled    = %[2]t
			event_type_name = "HOST_DOWN"

			notifications = [{
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = %[3]t
				email_enabled = %[4]t
				roles = ["GROUP_DATA_ACCESS_READ_ONLY", "GROUP_CLUSTER_MANAGER"]
			}]

			matchers = [{
				field_name = %[5]q
				operator   = %[6]q
				value      = %[7]q
			}, {
				field_name = %[8]q
				operator   = %[9]q
				value      = %[10]q
			}]
		}
	`, projectID, enabled, smsEnabled, emailEnabled,
		m1["fieldName"], m1["operator"], m1["value"],
		m2["fieldName"], m2["operator"], m2["value"])
}

func configWithMetricUpdated(projectID string, enabled bool, threshold float64) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = %[1]q
			enabled    = %[2]t
			event_type_name = "OUTSIDE_METRIC_THRESHOLD"

			notifications = [{
				type_name     = "GROUP"
				interval_min  = 5
				delay_min     = 0
				sms_enabled   = false
				email_enabled = true
				roles = ["GROUP_DATA_ACCESS_READ_ONLY"]
			}]

			matchers = [{
				field_name = "HOSTNAME_AND_PORT"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}]

			metric_threshold = {
				metric_name = "ASSERT_REGULAR"
				operator    = "LESS_THAN"
				threshold   = %[3]f
				units       = "RAW"
				mode        = "AVERAGE"
			}
		}
	`, projectID, enabled, threshold)
}

func configWithDataDog(projectID, dataDogAPIKey, dataDogRegion string, enabled bool, ddNotificationMap map[string]string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_third_party_integration" "atlas_datadog" {
			project_id = %[1]q
			api_key  = %[2]q
			region   = %[3]q
			type     = "DATADOG"
		}

		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = mongodbatlas_third_party_integration.atlas_datadog.project_id
			event_type_name = "REPLICATION_OPLOG_WINDOW_RUNNING_OUT"
			enabled    = %[4]t

			notifications = [{
				type_name       = %[5]q
				datadog_api_key = mongodbatlas_third_party_integration.atlas_datadog.api_key
				datadog_region  = mongodbatlas_third_party_integration.atlas_datadog.region
				interval_min    = %[6]v
				delay_min       = %[7]v
			}]

			matchers = [{
				field_name = "REPLICA_SET_NAME"
				operator   = "EQUALS"
				value      = "SECONDARY"
			}]

			threshold = {
				operator    = "LESS_THAN"
				threshold   = 72
				units       = "HOURS"
			}
		}
	`, projectID, dataDogAPIKey, dataDogRegion, enabled, ddNotificationMap["type_name"], ddNotificationMap["interval_min"], ddNotificationMap["delay_min"])
}

func configWithPagerDuty(projectID, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = %[1]q
			enabled    = %[3]t
			event_type_name = "NO_PRIMARY"

			notifications = [{
				type_name   = "PAGER_DUTY"
				service_key = %[2]q
				delay_min   = 0
			}]
		}
	`, projectID, serviceKey, enabled)
}

func configWithEmail(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = %[1]q
			enabled    = %[2]t
			event_type_name = "NO_PRIMARY"

			notifications = [{
				type_name     = "EMAIL"
				interval_min  = 60
				email_address = "test@mongodbtest.com"
			}]
		}
	`, projectID, enabled)
}

func configWithPagerDutyNotifierID(projectID, notifierID, serviceKey string, delayMin int) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id = %[1]q
			enabled  = true
			event_type_name = "NO_PRIMARY"

			notifications = [{
				type_name   = "PAGER_DUTY"
				notifier_id = %[2]q
				service_key = %[3]q
				delay_min   = %[4]d
			}]
		}
	`, projectID, notifierID, serviceKey, delayMin)
}

func configWithOpsGenie(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id = %[1]q
			enabled    = %[3]t
			event_type_name = "NO_PRIMARY"

			notifications = [{
				type_name          = "OPS_GENIE"
				ops_genie_api_key  = %[2]q
				ops_genie_region   = "US"
				delay_min          = 0
			}]
		}
	`, projectID, apiKey, enabled)
}

func configWithVictorOps(projectID, apiKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id = %[1]q
			enabled    = %[3]t
			event_type_name = "NO_PRIMARY"

			notifications = [{
				type_name              = "VICTOR_OPS"
				victor_ops_api_key     = %[2]q
				victor_ops_routing_key = "testing"
				delay_min              = 0
			}]
		}
	`, projectID, apiKey, enabled)
}

func configWithEmptyMetricThresholdConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id   = %[1]q
			enabled    = %[2]t
			event_type_name = "CLUSTER_MONGOS_IS_MISSING"

			notifications =[{
				type_name     = "GROUP"
				interval_min  = 60
				delay_min     = 0
				sms_enabled   = true
				email_enabled = false
				roles         = ["GROUP_OWNER"]
			}]
		}
	`, projectID, enabled)
}

func configWithEmptyMatcherMetricThresholdConfig(projectID string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
		  group_id   = %[1]q
		  enabled    = %[2]t
		  event_type_name = "CLUSTER_MONGOS_IS_MISSING"

		  notifications = [{
			type_name     = "GROUP"
			interval_min  = 60
			delay_min     = 0
			sms_enabled   = true
			email_enabled = false
			roles         = ["GROUP_OWNER"]
		  }]
		}
	`, projectID, enabled)
}

func configWithSeverityOverride(projectID string, severity *string) string {
	severityOverride := ""
	if severity != nil {
		severityOverride = fmt.Sprintf("severity_override = %[1]q", *severity)
	}

	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration_api" "test" {
			group_id        = %[1]q
			enabled         = true
			event_type_name = "NO_PRIMARY"
			%[2]s

			notifications = [{
				type_name     = "EMAIL"
				interval_min  = 60
				email_address = "test@mongodbtest.com"
			}]
		}
		`, projectID, severityOverride)
}
