package privatelinkendpoint_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkPrivateLinkEndpoint_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		orgID        = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName  = acc.RandomProjectName()
		region       = "us-east-1"
		providerName = "AWS"
		config       = configBasic(orgID, projectName, providerName, region)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttr(resourceName, "provider_name", providerName),
					resource.TestCheckResourceAttr(resourceName, "region", region),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
