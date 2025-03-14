package mig

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestStepCheckEmptyPlan(config string) resource.TestStep {
	testStep := acc.TestStepCheckEmptyPlan(config)
	// migration tests need provider to be defined in each step
	testStep.ProtoV6ProviderFactories = acc.TestAccProviderV6Factories
	return testStep
}
