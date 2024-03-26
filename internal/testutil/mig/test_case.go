package mig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func ConvertToMigrationTest(t *testing.T, test *resource.TestCase, externalProviders ...map[string]resource.ExternalProvider) resource.TestCase {
	t.Helper()
	checkLastVersion(t)
	require.GreaterOrEqual(t, len(test.Steps), 1, "Must have at least 1 test step.")
	firstStep := test.Steps[0]

	if len(externalProviders) == 0 {
		externalProviders = append(externalProviders, ExternalProviders())
	}
	steps := []resource.TestStep{}

	for _, provider := range externalProviders {
		steps = append(steps, *useExternalProvider(firstStep, provider))
	}
	steps = append(steps, TestStepCheckEmptyPlan(firstStep.Config))
	return resource.TestCase{
		PreCheck:     test.PreCheck,
		CheckDestroy: test.CheckDestroy,
		Steps:        steps,
	}
}

//nolint:gocritic
func useExternalProvider(step resource.TestStep, provider map[string]resource.ExternalProvider) *resource.TestStep {
	step.ExternalProviders = provider
	return &step
}
