package alertconfiguration_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"go.mongodb.org/atlas-sdk/v20231115007/admin"
)

func TestAccConfigDSAlertConfiguration_basic(t *testing.T) {
	var (
		alert          = &admin.GroupAlertsConfig{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configBasicDS(orgID, projectName),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
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
		alert          = &admin.GroupAlertsConfig{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithThreshold(orgID, projectName, true, 1),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
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
		alert          = &admin.GroupAlertsConfig{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		outputLabel    = "resource_import"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithOutputs(orgID, projectName, outputLabel),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
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
		alert          = &admin.GroupAlertsConfig{}
		dataSourceName = "data.mongodbatlas_alert_configuration.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		serviceKey     = dummy32CharKey
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheckBasic(t) },
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configWithPagerDutyDS(orgID, projectName, serviceKey, true),
				Check: resource.ComposeTestCheckFunc(
					checkExists(dataSourceName, alert),
					resource.TestCheckResourceAttrSet(dataSourceName, "project_id"),
				),
			},
		},
	})
}

func configBasicDS(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
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
	`, orgID, projectName)
}

func configWithThreshold(orgID, projectName string, enabled bool, threshold float64) string {
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

		data "mongodbatlas_alert_configuration" "test" {
			project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
			alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"
		}
	`, orgID, projectName, enabled, threshold)
}

func configWithOutputs(orgID, projectName, outputLabel string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id

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
			project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
			alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"

			output {
				type = "resource_import"
				label = %[3]q
			}
			output {
				type = "resource_hcl"
				label = %[3]q
			}
		}
	`, orgID, projectName, outputLabel)
}

func configWithPagerDutyDS(orgID, projectName, serviceKey string, enabled bool) string {
	return fmt.Sprintf(`
resource "mongodbatlas_project" "test" {
	name   = %[2]q
	org_id = %[1]q
}
resource "mongodbatlas_alert_configuration" "test" {
  project_id = mongodbatlas_project.test.id
  event_type = "NO_PRIMARY"
  enabled    = "%[4]t"

  notification {
    type_name    = "PAGER_DUTY"
    service_key  = %[3]q
    delay_min    = 0
  }
}

data "mongodbatlas_alert_configuration" "test" {
  project_id             = "${mongodbatlas_alert_configuration.test.project_id}"
  alert_configuration_id = "${mongodbatlas_alert_configuration.test.id}"
}
	`, orgID, projectName, serviceKey, enabled)
}
