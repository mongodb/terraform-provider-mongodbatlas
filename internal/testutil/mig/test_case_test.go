package mig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestConvertToMigration(t *testing.T) {
	t.Setenv("MONGODB_ATLAS_LAST_VERSION", "1.2.3")
	var (
		preCheckCalled     = false
		checkDestroyCalled = false
		config             = "resource \"dummy\" \"this\" {}"
	)
	preCheck := func() {
		preCheckCalled = true
	}
	firstStep := resource.TestStep{
		Config: config,
		Check:  resource.TestCheckResourceAttrSet("someTarget", "someAttribute"),
	}

	asserter := assert.New(t)

	convertAndCall := func(test resource.TestCase) resource.TestCase {
		newTest := mig.CreateTest(t, &test)
		newTest.PreCheck()
		if newTest.CheckDestroy != nil {
			asserter.NoError(newTest.CheckDestroy(nil))
		}
		return newTest
	}
	defaultAssertions := func(test resource.TestCase) {
		t.Helper()
		asserter.Len(test.Steps, 2, "Expected 2 steps (one extra test step)")
		newFirstStep := test.Steps[0]
		asserter.Equal(config, newFirstStep.Config)

		planStep := test.Steps[1]
		asserter.Equal(mig.TestStepCheckEmptyPlan(config), planStep)
	}

	t.Run("normal call with check and destroy", func(t *testing.T) {
		checkDestroy := func(*terraform.State) error {
			checkDestroyCalled = true
			return nil
		}
		test := convertAndCall(resource.TestCase{
			PreCheck:     preCheck,
			CheckDestroy: checkDestroy,
			Steps: []resource.TestStep{
				firstStep,
			},
		})
		asserter.True(preCheckCalled)
		asserter.True(checkDestroyCalled)
		defaultAssertions(test)
	})

	t.Run("checkDestroy=nil has no panic", func(t *testing.T) {
		test := convertAndCall(resource.TestCase{
			PreCheck: preCheck,
			Steps: []resource.TestStep{
				firstStep,
			},
		})
		defaultAssertions(test)
	})

	t.Run("more than 1 step uses only 1 step", func(t *testing.T) {
		test := convertAndCall(resource.TestCase{
			PreCheck: preCheck,
			Steps: []resource.TestStep{
				firstStep,
				{
					Config: "differentConfig",
					Check:  resource.TestCheckResourceAttrSet("target", "attribute"),
				},
			},
		})
		defaultAssertions(test)
	})
	// ConvertToMigrationTestUseExternalProvider

	t.Run("explicit ExternalProvider version an no additional providers", func(t *testing.T) {
		test := mig.CreateTestUseExternalProvider(t, &resource.TestCase{
			PreCheck: preCheck,
			Steps:    []resource.TestStep{firstStep},
		}, acc.ExternalProviders("1.2.3"), nil)
		asserter.Len(test.Steps, 2, "Expected 2 steps (one extra test step)")
		newFirstStep := test.Steps[0]
		asserter.Equal(config, newFirstStep.Config)
		asserter.Equal("1.2.3", test.Steps[0].ExternalProviders["mongodbatlas"].VersionConstraint)
	})

	t.Run("explicit ExternalProviders and additional providers", func(t *testing.T) {
		test := mig.CreateTestUseExternalProvider(t, &resource.TestCase{
			PreCheck: preCheck,
			Steps:    []resource.TestStep{firstStep},
		}, acc.ExternalProvidersWithAWS("1.2.3"), acc.ExternalProvidersOnlyAWS())
		asserter.Len(test.Steps, 2, "Expected 2 steps (one extra test step)")
		newFirstStep := test.Steps[0]
		asserter.Equal(config, newFirstStep.Config)

		asserter.Equal("1.2.3", test.Steps[0].ExternalProviders["mongodbatlas"].VersionConstraint)
		asserter.Equal(acc.AwsProviderVersion, test.Steps[0].ExternalProviders["aws"].VersionConstraint)
		asserter.Equal(acc.AwsProviderVersion, test.Steps[1].ExternalProviders["aws"].VersionConstraint)
	})
}
