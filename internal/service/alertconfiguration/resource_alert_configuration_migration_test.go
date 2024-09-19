package alertconfiguration_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigConfigRSAlertConfiguration_withNotificationsMetricThreshold(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		config    = configBasicRS(projectID, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigRSAlertConfiguration_withThreshold(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		config    = configWithThresholdUpdated(projectID, true, 1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "threshold_config.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigRSAlertConfiguration_withEmptyOptionalBlocks(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		config    = configWithEmptyOptionalBlocks(projectID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
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

func TestMigConfigRSAlertConfiguration_withMultipleMatchers(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		config    = configWithMatchers(projectID, true, false, true,
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
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "matcher.#", "2"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigRSAlertConfiguration_withEmptyOptionalAttributes(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		config    = configWithEmptyOptionalAttributes(projectID)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckResourceAttr(resourceName, "notification.#", "1"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigConfigRSAlertConfiguration_withDataDog(t *testing.T) {
	var (
		projectID = acc.ProjectIDExecution(t)
		ddAPIKey  = dummy32CharKey
		ddRegion  = "US"

		groupNotificationMap = map[string]string{
			"type_name":    "GROUP",
			"interval_min": "5",
			"delay_min":    "0",
		}

		datadogNotificationMap = map[string]string{
			"type_name":       "DATADOG",
			"interval_min":    "5",
			"delay_min":       "0",
			"datadog_api_key": ddAPIKey,
			"datadog_region":  ddRegion,
		}

		config = configWithDataDog(projectID, ddAPIKey, ddRegion, true, datadogNotificationMap, groupNotificationMap)
	)

	// 1.20.0 introduced handling of a breaking change from the API which required notification.#.integration_id to be
	// updated to Option/Computed from an Optional attribute. This impacted only notifications with integrations.
	mig.SkipIfVersionBelow(t, "1.20.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "project_id", projectID),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", groupNotificationMap),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "notification.*", datadogNotificationMap),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
