package mig_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
	"github.com/stretchr/testify/assert"
)

func TestConvertToMigration(t *testing.T) {
	var (
		preCheckCalled     = false
		checkDestroyCalled = false
		config             = "someTerraformConfig"
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
		newTest := mig.ConvertToMigrationTest(t, &test)
		newTest.PreCheck()
		if newTest.CheckDestroy != nil {
			asserter.NoError(newTest.CheckDestroy(nil))
		}
		return newTest
	}
	defaultAssertions := func(test resource.TestCase) {
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

	t.Run("check destroy is nil has no panic", func(t *testing.T) {
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

	t.Run("explicit ExternalProvider version", func(t *testing.T) {
		test := mig.ConvertToMigrationTest(t, &resource.TestCase{
			PreCheck: preCheck,
			Steps:    []resource.TestStep{firstStep},
		}, acc.ExternalProviders("1.2.3"))
		defaultAssertions(test)
		asserter.Equal("1.2.3", test.Steps[0].ExternalProviders["mongodbatlas"].VersionConstraint)
	})

	t.Run("multiple ExternalProviders", func(t *testing.T) {
		test := mig.ConvertToMigrationTest(t, &resource.TestCase{
			PreCheck: preCheck,
			Steps:    []resource.TestStep{firstStep},
		}, acc.ExternalProviders("1.2.3"), acc.ExternalProvidersWithAWS("4.5.6"))
		asserter.Len(test.Steps, 3)
		asserter.Equal("1.2.3", test.Steps[0].ExternalProviders["mongodbatlas"].VersionConstraint)
		asserter.Equal("4.5.6", test.Steps[1].ExternalProviders["mongodbatlas"].VersionConstraint)
		// must be upgraded when the aws provider version is changed
		asserter.Equal("5.1.0", test.Steps[1].ExternalProviders["aws"].VersionConstraint)
		planStep := test.Steps[2]
		asserter.Equal(mig.TestStepCheckEmptyPlan(config), planStep)
	})
}
