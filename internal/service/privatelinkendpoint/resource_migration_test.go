package privatelinkendpoint_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/constant"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigNetworkPrivateLinkEndpoint_basic(t *testing.T) {
	var (
		resourceName = "mongodbatlas_privatelink_endpoint.test"
		projectID    = acc.ProjectIDExecution(t)
		region       = "us-west-2"
		providerName = constant.AWS
		config       = configBasic(projectID, providerName, region)
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
