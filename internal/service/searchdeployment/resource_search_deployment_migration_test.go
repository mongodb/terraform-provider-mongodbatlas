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
		orgID           = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName     = acc.RandomProjectName()
		clusterName     = acc.RandomClusterName()
		instanceSize    = "S30_HIGHCPU_NVME"
		searchNodeCount = 3
		config          = configBasic(orgID, projectName, clusterName, instanceSize, searchNodeCount)
	)
	mig.SkipIfVersionBelow(t, "1.13.0")
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
