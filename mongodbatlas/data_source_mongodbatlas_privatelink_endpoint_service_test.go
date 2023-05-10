package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkDSPrivateLinkEndpointServiceAWS_basic(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "data.mongodbatlas_privatelink_endpoint_service.test"

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	vpcID := os.Getenv("AWS_VPC_ID")
	subnetID := os.Getenv("AWS_SUBNET_ID")
	securityGroupID := os.Getenv("AWS_SECURITY_GROUP_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateLinkEndpointServiceDataSourceConfig(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateLinkEndpointServiceExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint_service_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateLinkEndpointServiceDataSourceConfig(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "us-east-1"
			access_key = "%s"
			secret_key = "%s"
		}

		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%s"
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%s"]
			security_group_ids = ["%s"]
		}

		resource "mongodbatlas_privatelink_endpoint_service" "test" {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id       =  aws_vpc_endpoint.ptfe_service.id
			private_link_id = mongodbatlas_privatelink_endpoint.test.private_link_id
			provider_name = "%[4]s"
		}

		data "mongodbatlas_privatelink_endpoint_service" "test" {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			private_link_id       =  mongodbatlas_privatelink_endpoint_service.test.private_link_id
			endpoint_service_id = mongodbatlas_privatelink_endpoint_service.test.endpoint_service_id
			provider_name = "%[4]s"
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID)
}
