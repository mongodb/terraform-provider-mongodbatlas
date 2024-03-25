package mig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func ConvertToMigrationTest(t *testing.T, test *resource.TestCase) resource.TestCase {
	t.Helper()
	checkLastVersion(t)

	firstStep := test.Steps[0]
	firstStep.ExternalProviders = ExternalProviders()

	return resource.TestCase{
		PreCheck:     test.PreCheck,
		CheckDestroy: test.CheckDestroy,
		Steps: []resource.TestStep{
			firstStep,
			TestStepCheckEmptyPlan(firstStep.Config),
		},
	}
}
