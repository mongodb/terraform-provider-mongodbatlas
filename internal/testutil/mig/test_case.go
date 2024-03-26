package mig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

// ConvertToMigrationTest returns an updated TestCase that reuses step 1 and adds a TestStepCheckEmptyPlan
// Requires: `MONGODB_ATLAS_LAST_VERSION` to be present
func ConvertToMigrationTest(t *testing.T, test *resource.TestCase) resource.TestCase {
	t.Helper()
	validateReusableCase(t, test)
	firstStep := test.Steps[0]
	steps := []resource.TestStep{
		useExternalProvider(&firstStep, ExternalProviders()),
		TestStepCheckEmptyPlan(firstStep.Config),
	}
	return reuseCase(test, steps)
}

// ConvertToMigrationTestUseExternalProvider returns an updated TestCase that reuses step 1 and adds a TestStepCheckEmptyPlan with the additionalProviders
// Requires: `MONGODB_ATLAS_LAST_VERSION` to be present
// externalProviders: e.g., ExternalProvidersWithAWS() or ExternalProviders("specific_sem_ver")
// additionalProviders: e.g., acc.ExternalProvidersOnlyAWS(), can also be nil
func ConvertToMigrationTestUseExternalProvider(t *testing.T, test *resource.TestCase, externalProviders, additionalProviders map[string]resource.ExternalProvider) resource.TestCase {
	t.Helper()
	validateReusableCase(t, test)
	firstStep := test.Steps[0]
	require.NotContains(t, additionalProviders, "mongodbatlas", "Will use the local provider, cannot specify mongodbatlas provider")
	emptyPlanStep := TestStepCheckEmptyPlan(firstStep.Config)
	steps := []resource.TestStep{
		useExternalProvider(&firstStep, externalProviders),
		useExternalProvider(&emptyPlanStep, additionalProviders),
	}
	return reuseCase(test, steps)
}

func validateReusableCase(t *testing.T, test *resource.TestCase) {
	t.Helper()
	checkLastVersion(t)
	require.GreaterOrEqual(t, len(test.Steps), 1, "Must have at least 1 test step.")
	require.NotEmpty(t, test.Steps[0].Config, "First step of migration test must use Config")
}

func useExternalProvider(step *resource.TestStep, provider map[string]resource.ExternalProvider) resource.TestStep {
	step.ExternalProviders = provider
	return *step
}

// Note how we don't set ProtoV6ProviderFactories and instead specify providers on each step
func reuseCase(test *resource.TestCase, steps []resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:     test.PreCheck,
		CheckDestroy: test.CheckDestroy,
		ErrorCheck:   test.ErrorCheck,
		Steps:        steps,
	}
}
