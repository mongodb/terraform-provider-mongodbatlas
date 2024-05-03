package mig

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func CreateAndRunTest(t *testing.T, test *resource.TestCase) {
	t.Helper()
	acc.SkipInUnitTest(t) // Migration tests create external resources and use MONGODB_ATLAS_LAST_VERSION env-var.
	resource.ParallelTest(t, CreateTest(t, test))
}

// avoids running migration test in parallel
func CreateAndRunTestNonParallel(t *testing.T, test *resource.TestCase) {
	t.Helper()
	acc.SkipInUnitTest(t) // Migration tests create external resources and use MONGODB_ATLAS_LAST_VERSION env-var.
	resource.Test(t, CreateTest(t, test))
}

func CreateTestAndRunUseExternalProvider(t *testing.T, test *resource.TestCase, externalProviders, additionalProviders map[string]resource.ExternalProvider) {
	t.Helper()
	acc.SkipInUnitTest(t) // Migration tests create external resources and use MONGODB_ATLAS_LAST_VERSION env-var.
	resource.ParallelTest(t, CreateTestUseExternalProvider(t, test, externalProviders, additionalProviders))
}

func CreateTestAndRunUseExternalProviderNonParallel(t *testing.T, test *resource.TestCase, externalProviders, additionalProviders map[string]resource.ExternalProvider) {
	t.Helper()
	acc.SkipInUnitTest(t) // Migration tests create external resources and use MONGODB_ATLAS_LAST_VERSION env-var.
	resource.Test(t, CreateTestUseExternalProvider(t, test, externalProviders, additionalProviders))
}

// CreateTest returns a new TestCase that reuses step 1 and adds a TestStepCheckEmptyPlan.
// Requires: `MONGODB_ATLAS_LAST_VERSION` to be present.
func CreateTest(t *testing.T, test *resource.TestCase) resource.TestCase {
	t.Helper()
	validateReusableCase(t, test)
	firstStep := test.Steps[0]
	steps := []resource.TestStep{
		useExternalProvider(&firstStep, ExternalProviders()),
		TestStepCheckEmptyPlan(firstStep.Config),
	}
	newTest := reuseCase(test, steps)
	return newTest
}

// CreateTestUseExternalProvider returns a new TestCase that reuses step 1 and adds a TestStepCheckEmptyPlan with the additionalProviders.
// Requires: `MONGODB_ATLAS_LAST_VERSION` to be present.
// externalProviders: e.g., ExternalProvidersWithAWS() or ExternalProviders("specific_sem_ver").
// additionalProviders: e.g., acc.ExternalProvidersOnlyAWS(), can also be nil.
func CreateTestUseExternalProvider(t *testing.T, test *resource.TestCase, externalProviders, additionalProviders map[string]resource.ExternalProvider) resource.TestCase {
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

func validateReusableCase(tb testing.TB, test *resource.TestCase) {
	tb.Helper()
	checkLastVersion(tb)
	require.GreaterOrEqual(tb, len(test.Steps), 1, "Must have at least 1 test step.")
	require.NotEmpty(tb, test.Steps[0].Config, "First step of migration test must use Config")
}

func useExternalProvider(step *resource.TestStep, provider map[string]resource.ExternalProvider) resource.TestStep {
	step.ExternalProviders = provider
	return *step
}

// Note how we don't set ProtoV6ProviderFactories and instead specify providers on each step.
func reuseCase(test *resource.TestCase, steps []resource.TestStep) resource.TestCase {
	return resource.TestCase{
		PreCheck:     test.PreCheck,
		CheckDestroy: test.CheckDestroy,
		ErrorCheck:   test.ErrorCheck,
		Steps:        steps,
	}
}
