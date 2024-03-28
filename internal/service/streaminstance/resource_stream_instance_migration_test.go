package streaminstance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamInstance_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_instance.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
	)
	mig.SkipIfVersionBelow(t, "1.14.0")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckPreviewFlag(t); acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider),
				Check:             streamInstanceAttributeChecks(resourceName, instanceName, region, cloudProvider),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						acc.DebugPlan(),
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
