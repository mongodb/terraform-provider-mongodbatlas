package privatelinkendpointservice_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/common/conversion"
	"github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/acc"
)

func TestAccNetworkRSPrivateLinkEndpointServiceAWS_Complete(t *testing.T) {
	acc.SkipTestForCI(t)
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
		PreCheck:                 func() { acc.PreCheck(t); acc.PreCheckAwsEnv(t) },
		CheckDestroy:             testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy,
		ProtoV6ProviderFactories: acc.TestAccProviderV6Factories,
		ExternalProviders:        acc.ExternalProvidersOnlyAWS(),
		Steps: []resource.TestStep{
			{
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
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccCheckMongoDBAtlasPrivateLinkEndpointServiceImportStateIDFunc(resourceName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_link_id"},
			},
		},
	})
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		ids := conversion.DecodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s--%s", ids["project_id"], ids["private_link_id"], ids["endpoint_service_id"], ids["provider_name"]), nil
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["endpoint_service_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("the MongoDB Interface Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["endpoint_service_id"], rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service" {
			continue
		}
		ids := conversion.DecodeStateID(rs.Primary.ID)
		_, _, err := acc.Conn().PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["endpoint_service_id"])
		if err == nil {
			return fmt.Errorf("the MongoDB Private Endpoint(%s) still exists", ids["endpoint_service_id"])
		}
	}
	return nil
}

func testAccMongoDBAtlasPrivateLinkEndpointServiceConfigCompleteAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, serviceResourceName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region        = "%[5]s"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}

		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
		}

		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = "%[6]s"
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = ["%[7]s"]
			security_group_ids = ["%[8]s"]
			
		}

		resource "mongodbatlas_privatelink_endpoint_service" %[9]q {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id   =  aws_vpc_endpoint.ptfe_service.id
			private_link_id       = mongodbatlas_privatelink_endpoint.test.id
			provider_name         = "%[4]s"
		}
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, vpcID, subnetID, securityGroupID, serviceResourceName)
}
