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
		config    = configBasic(projectID, true)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy(),
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
		config    = configWithThreshold(projectID, true, 1)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy(),
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
		CheckDestroy: checkDestroy(),
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
		CheckDestroy: checkDestroy(),
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
		CheckDestroy: checkDestroy(),
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
	mig.SkipIfVersionBelow(t, "1.20.0")
	mig.CreateAndRunTestNonParallel(t, datadogTestCase(t)) // not run in parallel so acc and mig tests don't interfere
}
