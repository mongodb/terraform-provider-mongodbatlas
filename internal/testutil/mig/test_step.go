package mig

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

// TestStep is being deprecated in favor of TestStepCheckEmptyPlan
func TestStep(config string) resource.TestStep {
	return TestStepCheckEmptyPlan(config)
}

func TestStepCheckEmptyPlan(config string) resource.TestStep {
	return resource.TestStep{
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		Config:                   config,
		ConfigPlanChecks: resource.ConfigPlanChecks{
			PreApply: []plancheck.PlanCheck{
				acc.DebugPlan(),
				plancheck.ExpectEmptyPlan(),
			},
		},
	}
}
