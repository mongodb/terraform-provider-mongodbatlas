package privatelinkendpointserviceserverless_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestMigServerlessPrivateLinkEndpointService_basic(t *testing.T) {
	var (
		resourceName   = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		datasourceName = "data.mongodbatlas_privatelink_endpoint_service_serverless.test"
		orgID          = os.Getenv("MONGODB_ATLAS_ORG_ID")
		projectName    = acc.RandomProjectName()
		instanceName   = acc.RandomClusterName()
		commentOrigin  = "this is a comment for serverless private link endpoint"
		config         = configBasic(orgID, projectName, instanceName, commentOrigin)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProviders(),
				Config:            config,
				Check: resource.ComposeTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttr(datasourceName, "comment", commentOrigin),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
