package streamworkspace_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigStreamRSStreamWorkspace_basic(t *testing.T) {
	var (
		resourceName  = "mongodbatlas_stream_workspace.test"
		projectID     = acc.ProjectIDExecution(t)
		workspaceName = acc.RandomName()
		config        = acc.StreamInstanceConfig(projectID, workspaceName, region, cloudProvider)
	)
	mig.SkipIfVersionBelow(t, "1.16.0") // when reached GA

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheckBasic(t) },
		CheckDestroy: acc.CheckDestroyStreamInstance,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             streamWorkspaceAttributeChecks(resourceName, workspaceName, region, cloudProvider),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
