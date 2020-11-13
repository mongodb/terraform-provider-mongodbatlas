package mongodbatlas

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceMongoDBAtlasPrivateEndpointLinkAWS_basic(t *testing.T) {
	SkipTestExtCred(t)
	resourceName := "data.mongodbatlas_private_endpoint_service_link.test"

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	projectID := os.Getenv("MONGODB_ATLAS_PROJECT_ID")
	region := os.Getenv("AWS_REGION")
	providerName := "AWS"

	vpcID := os.Getenv("AWS_VPC_ID")
	subnetID := os.Getenv("AWS_SUBNET_ID")
	securityGroupID := os.Getenv("AWS_SECURITY_GROUP_ID")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); checkAwsEnv(t); checkPeeringEnvAWS(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongoDBAtlasPrivateEndpointLinkDataSourceConfig(
					awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongoDBAtlasPrivateEndpointLinkExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "private_link_id"),
					resource.TestCheckResourceAttrSet(resourceName, "interface_endpoint_id"),
				),
			},
		},
	})
}

func testAccMongoDBAtlasPrivateEndpointLinkDataSourceConfig(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "us-east-1"
			access_key = "%s"
			secret_key = "%s"
		}

		resource "mongodbatlas_private_endpoint" "test" {
			project_id    = "%s"
			provider_name = "%s"
			region        = "%s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%s"
			service_name       = "${mongodbatlas_private_endpoint.test.endpoint_service_name}"
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%s"]
			security_group_ids = ["%s"]
		}

		resource "mongodbatlas_private_endpoint_service_link" "test" {
			project_id            = "${mongodbatlas_private_endpoint.test.project_id}"
			private_link_id       = "${mongodbatlas_private_endpoint.test.private_link_id}"
			interface_endpoint_id = "${aws_vpc_endpoint.ptfe_service.id}"
		}

		data "mongodbatlas_private_endpoint_service_link" "test" {
			project_id            = "${mongodbatlas_private_endpoint_service_link.test.project_id}"
			private_link_id       = "${mongodbatlas_private_endpoint_service_link.test.private_link_id}"
			interface_endpoint_id = "${mongodbatlas_private_endpoint_service_link.test.interface_endpoint_id}"
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID)
}
