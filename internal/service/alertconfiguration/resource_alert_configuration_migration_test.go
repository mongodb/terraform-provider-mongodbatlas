package alertconfiguration_test

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115007/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationConfigRSAlertConfiguration_withNotificationsMetricThreshold(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
		config       = configBasicRS(orgID, projectName, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_withThreshold(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
		config       = configWithThresholdUpdated(orgID, projectName, true, 1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threshold_config.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_withEmptyOptionalBlocks(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
		config       = configWithEmptyOptionalBlocks(orgID, projectName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "threshold_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "metric_threshold_config.#", "0"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_withMultipleMatchers(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
		config       = configWithMatchers(orgID, projectName, true, false, true,
			map[string]interface{}{
				"fieldName": "TYPE_NAME",
				"operator":  "EQUALS",
				"value":     "SECONDARY",
			},
			map[string]interface{}{
				"fieldName": "TYPE_NAME",
				"operator":  "CONTAINS",
				"value":     "MONGOS",
			})
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_withEmptyOptionalAttributes(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		alert        = &admin.GroupAlertsConfig{}
		config       = configWithEmptyOptionalAttributes(orgID, projectName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

// configWithEmptyOptionalAttributes does not define notification.delay_min, notification.sms_enabled, and metric_threshold_config.threshold.
func configWithEmptyOptionalAttributes(orgID, projectName string) string {
	return fmt.Sprintf(`
		resource "mongodbatlas_project" "test" {
			name   = %[2]q
			org_id = %[1]q
		}
		resource "mongodbatlas_alert_configuration" "test" {
			project_id = mongodbatlas_project.test.id
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
	`, orgID, projectName)
}

func configWithEmptyOptionalBlocks(orgID, projectName string) string {
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
	`, orgID, projectName)
}
