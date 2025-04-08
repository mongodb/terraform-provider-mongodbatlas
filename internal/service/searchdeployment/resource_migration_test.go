package searchdeployment_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigSearchDeployment_basic(t *testing.T) {
	var (
		resourceName    = "mongodbatlas_search_deployment.test"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_EAR_PE_AWS_ID") // to use RequirePrivateNetworking, Atlas Project is required to have FF enabled
		clusterName     = acc.RandomClusterName()
		instanceSize    = "S30_HIGHCPU_NVME"
		searchNodeCount = 3
		config          = configBasic(projectID, clusterName, instanceSize, searchNodeCount)
	)
	mig.SkipIfVersionBelow(t, "1.32.0") // enabled_for_search_nodes introduced in this version
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check:             resource.ComposeAggregateTestCheckFunc(searchNodeChecks(resourceName, clusterName, instanceSize, searchNodeCount)...),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
