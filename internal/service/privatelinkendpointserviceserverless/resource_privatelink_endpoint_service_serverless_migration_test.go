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
		projectID      = acc.ProjectIDExecution(t)
		instanceName   = acc.RandomClusterName()
		commentOrigin  = "this is a comment for serverless private link endpoint"
		config         = configBasic(projectID, instanceName, commentOrigin)
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
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
					resource.TestCheckResourceAttr(resourceName, "comment", commentOrigin),
					resource.TestCheckResourceAttr(datasourceName, "comment", commentOrigin),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}

func TestMigServerlessPrivateLinkEndpointService_AWSVPC(t *testing.T) {
	mig.SkipIfVersionBelow(t, "1.16.0") // bug fix included for https://github.com/mongodb/terraform-provider-mongodbatlas/issues/2011
	var (
		resourceName = "mongodbatlas_privatelink_endpoint_service_serverless.test"
		projectID    = acc.ProjectIDExecution(t)
		instanceName = acc.RandomClusterName()
		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
		config       = configAWSEndpoint(projectID, instanceName, awsAccessKey, awsSecretKey, true, "test comment")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { mig.PreCheckBasic(t) },
		CheckDestroy: checkDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProvidersWithAWS(),
				Config:            config,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "provider_name", "AWS"),
				),
			},
			mig.TestStepCheckEmptyPlan(config),
		},
	})
}
