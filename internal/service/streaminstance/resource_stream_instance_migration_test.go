package streaminstance_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamInstance_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_stream_instance.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomName()
		config       = acc.StreamInstanceConfig(projectID, instanceName, region, cloudProvider)
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             streamInstanceAttributeChecks(resourceName, instanceName, region, cloudProvider),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
