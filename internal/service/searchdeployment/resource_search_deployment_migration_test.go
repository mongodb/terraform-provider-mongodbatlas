package searchdeployment_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationSearchDeployment_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_search_deployment.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acctest.RandomWithPrefix("test-acc-search-dep")
		clusterName  = acctest.RandomWithPrefix("test-acc-search-dep")
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            configBasic(orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
				Check:             resource.ComposeTestCheckFunc(searchNodeChecks(resourceName, clusterName, "S30_HIGHCPU_NVME", 4)...),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config:                   configBasic(orgID, projectName, clusterName, "S30_HIGHCPU_NVME", 4),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
