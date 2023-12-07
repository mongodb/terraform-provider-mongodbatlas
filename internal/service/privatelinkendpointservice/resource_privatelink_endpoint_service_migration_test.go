package privatelinkendpointservice_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/mig"
)

func TestAccMigrationNetworkRSPrivateLinkEndpointService_Complete(t *testing.T) {
	acc.SkipTestExtCred(t)
	var (
		resourceSuffix = "test"
		resourceName   = fmt.Sprintf("mongodbatlas_privatelink_endpoint_service.%s", resourceSuffix)

		awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
		awsSecretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		providerName    = "AWS"
		projectID       = os.Getenv("MONGODB_ATLAS_PROJECT_ID")
		region          = os.Getenv("AWS_REGION")
		vpcID           = os.Getenv("AWS_VPC_ID")
		subnetID        = os.Getenv("AWS_SUBNET_ID")
		securityGroupID = os.Getenv("AWS_SECURITY_GROUP_ID")
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t) },
		CheckDestroy: testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy,
		Steps: []resource.TestStep{
			{
				ExternalProviders: mig.ExternalProvidersWithAWS(t, "5.1.0"),
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceConfigCompleteAWS(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, resourceSuffix,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_service_id"),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"aws": {
						VersionConstraint: "5.1.0",
						Source:            "hashicorp/aws",
					},
				},
				ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceConfigCompleteAWS(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, resourceSuffix,
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPreRefresh: []plancheck.PlanCheck{
						acc.DebugPlan(),
					},
				},
				PlanOnly: true,
			},
		},
	})
}
