package alertconfiguration_test

import (
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/atlas-sdk/v20231115001/admin"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccMigrationConfigRSAlertConfiguration_NotificationsWithMetricThreshold(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_alert_configuration.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		alert                 = &admin.GroupAlertsConfig{}
		config                = testAccMongoDBAtlasAlertConfigurationConfig(orgID, projectName, true)
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_WithThreshold(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_alert_configuration.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		alert                 = &admin.GroupAlertsConfig{}
		config                = testAccMongoDBAtlasAlertConfigurationConfigWithThresholdUpdated(orgID, projectName, true, 1)
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threshold_config.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_EmptyOptionalBlocks(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_alert_configuration.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		alert                 = &admin.GroupAlertsConfig{}
		config                = testAccMongoDBAtlasAlertConfigurationConfigEmptyOptionalBlocks(orgID, projectName)
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "threshold_config.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "metric_threshold_config.#", "0"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_MultipleMatchers(t *testing.T) {
	var (
		resourceName = "mongodbatlas_alert_configuration.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc")
		alert        = &admin.GroupAlertsConfig{}
		config       = testAccMongoDBAtlasAlertConfigurationConfigWithMatchers(orgID, projectName, true, false, true,
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
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "2"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

func TestAccMigrationConfigRSAlertConfiguration_EmptyOptionalAttributes(t *testing.T) {
	var (
		resourceName          = "mongodbatlas_alert_configuration.test"
		orgID                 = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName           = acctest.RandomWithPrefix("test-acc")
		alert                 = &admin.GroupAlertsConfig{}
		config                = testAccMongoDBAtlasAlertConfigurationConfigWithEmptyOptionalAttributes(orgID, projectName)
		lastVersionConstraint = os.Getenv("MONGODB_ATLAS_LAST_VERSION")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasicMigration(t) },
		CheckDestroy: testAccCheckMongoDBAtlasAlertConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"mongodbatlas": {
						VersionConstraint: lastVersionConstraint,
						Source:            "mongodb/mongodbatlas",
					},
				},
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasAlertConfigurationExists(resourceName, alert),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}

// does not define notification.delay_min, notification.sms_enabled, and metric_threshold_config.threshold
func testAccMongoDBAtlasAlertConfigurationConfigWithEmptyOptionalAttributes(orgID, projectName string) string {
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

func testAccMongoDBAtlasAlertConfigurationConfigEmptyOptionalBlocks(orgID, projectName string) string {
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
