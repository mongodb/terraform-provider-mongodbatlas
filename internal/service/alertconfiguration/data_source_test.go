package alertconfiguration_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccConfigDSAlertConfiguration_basic(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configBasicDS(projectID),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "notification.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "notification.0.notifier_id"),
					resource.TestCheckResourceAttr(dataSourceName, "matcher.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metric_threshold_config.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "threshold_config.#", "0"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfiguration_withThreshold(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithThreshold(projectID, true, 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "matcher.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metric_threshold_config.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "threshold_config.#", "1"),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfiguration_withOutput(t *testing.T) {
	var (
		projectID   = acc.ProjectIDExecution(t)
		outputLabel = "resource_import"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy(),
		Steps: []resource.TestStep{
			{
				Config: configWithOutputs(projectID, outputLabel),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(dataSourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "output.0.label", outputLabel),
					resource.TestCheckResourceAttr(dataSourceName, "output.0.type", "resource_import"),
					resource.TestCheckResourceAttrWith(dataSourceName, "output.0.value", acc.MatchesExpression("terraform import mongodbatlas_alert_configuration.*")),
					resource.TestCheckResourceAttr(dataSourceName, "output.1.label", outputLabel),
					resource.TestCheckResourceAttr(dataSourceName, "output.1.type", "resource_hcl"),
					resource.TestCheckResourceAttrWith(dataSourceName, "output.1.value", acc.MatchesExpression("resource \"mongodbatlas_alert_configuration\".*")),
				),
			},
		},
	})
}

func TestAccConfigDSAlertConfiguration_withPagerDuty(t *testing.T) {
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
				Config: configWithPagerDutyDS(projectID, serviceKey, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "project_id", projectID),
				),
			},
		},
	})
}

func configBasicDS(projectID string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
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
			project_id             = mongodbatlas_alert_configuration.test.project_id
			alert_configuration_id = mongodbatlas_alert_configuration.test.id
		}
	`, projectID)
}

func configWithThreshold(projectID string, enabled bool, threshold float64) string {
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

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = mongodbatlas_alert_configuration.test.project_id
			alert_configuration_id = mongodbatlas_alert_configuration.test.id
		}
	`, projectID, enabled, threshold)
}

func configWithOutputs(projectID, outputLabel string) string {
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

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = mongodbatlas_alert_configuration.test.project_id
			alert_configuration_id = mongodbatlas_alert_configuration.test.id

			output {
				type = "resource_import"
				label = %[2]q
			}
			output {
				type = "resource_hcl"
				label = %[2]q
			}
		}
	`, projectID, outputLabel)
}

func configWithPagerDutyDS(projectID, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = %[1]q
			enabled    = "%[3]t"
			event_type = "NO_PRIMARY"

			notification {
				type_name    = "PAGER_DUTY"
				service_key  = %[2]q
				delay_min    = 0
			}
		}

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = mongodbatlas_alert_configuration.test.project_id
			alert_configuration_id = mongodbatlas_alert_configuration.test.id
		}
	`, projectID, serviceKey, enabled)
}
