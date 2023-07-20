package mongodbatlas

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkRSPrivateLinkEndpointServiceAWS_Complete(t *testing.T) {
	SkipTestExtCred(t)
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
		PreCheck:          func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy,
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
		},
	})
}

func TestAccNetworkRSPrivateLinkEndpointServiceAWS_import(t *testing.T) {
	SkipTestExtCred(t)
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
		PreCheck:          func() { testAccPreCheck(t); testCheckAwsEnv(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy,
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

		ids := decodeStateID(rs.Primary.ID)

		return fmt.Sprintf("%s--%s--%s--%s", ids["project_id"], ids["private_link_id"], ids["endpoint_service_id"], ids["provider_name"]), nil
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*MongoDBClient).Atlas

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		ids := decodeStateID(rs.Primary.ID)

		_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["endpoint_service_id"])
		if err == nil {
			return nil
		}

		return fmt.Errorf("the MongoDB Interface Endpoint(%s) for the project(%s) does not exist", rs.Primary.Attributes["endpoint_service_id"], rs.Primary.Attributes["project_id"])
	}
}

func testAccCheckMongoDBAtlasPrivateLinkEndpointServiceDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*MongoDBClient).Atlas

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "mongodbatlas_privatelink_endpoint_service" {
			continue
		}

		ids := decodeStateID(rs.Primary.ID)
		_, _, err := conn.PrivateEndpoints.GetOnePrivateEndpoint(context.Background(), ids["project_id"], ids["provider_name"], ids["private_link_id"], ids["endpoint_service_id"])
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

func testAccMongoDBAtlasPrivateLinkEndpointServiceConfigUnmanagedAWS(awsAccessKey, awsSecretKey, projectID, providerName, region, serviceResourceName string) string {
	return fmt.Sprintf(`
		provider "aws" {
			region     = "%[5]s"
			access_key = "%[1]s"
			secret_key = "%[2]s"
		}
		resource "mongodbatlas_privatelink_endpoint" "test" {
			project_id    = "%[3]s"
			provider_name = "%[4]s"
			region        = "%[5]s"
		}
		resource "aws_vpc_endpoint" "ptfe_service" {
			vpc_id             = aws_vpc.primary.id
			service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
			vpc_endpoint_type  = "Interface"
			subnet_ids         = [aws_subnet.primary-az1.id]
			security_group_ids = [aws_security_group.primary_default.id]
			
		}
		resource "mongodbatlas_privatelink_endpoint_service" %[6]q {
			project_id            = mongodbatlas_privatelink_endpoint.test.project_id
			endpoint_service_id   =  aws_vpc_endpoint.ptfe_service.id
			private_link_id       = mongodbatlas_privatelink_endpoint.test.id
			provider_name         = %[4]q
		}
		resource "aws_vpc" "primary" {
			cidr_block           = "10.0.0.0/16"
			enable_dns_hostnames = true
			enable_dns_support   = true
		}

		resource "aws_internet_gateway" "primary" {
			vpc_id = aws_vpc.primary.id
		}

		resource "aws_route" "primary-internet_access" {
			route_table_id         = aws_vpc.primary.main_route_table_id
			destination_cidr_block = "0.0.0.0/0"
			gateway_id             = aws_internet_gateway.primary.id
		}
		  
		  # Subnet-A
		  resource "aws_subnet" "primary-az1" {
			vpc_id                  = aws_vpc.primary.id
			cidr_block              = "10.0.1.0/24"
			map_public_ip_on_launch = true
			availability_zone       = "%[5]sa"
		  }
		  
		  # Subnet-B
		  resource "aws_subnet" "primary-az2" {
			vpc_id                  = aws_vpc.primary.id
			cidr_block              = "10.0.2.0/24"
			map_public_ip_on_launch = false
			availability_zone       = "%[5]sb"
		  }
		  
		  resource "aws_security_group" "primary_default" {
			name_prefix = "default-"
			description = "Default security group for all instances in vpc"
			vpc_id      = aws_vpc.primary.id
			ingress {
			  from_port = 0
			  to_port   = 0
			  protocol  = "tcp"
			  cidr_blocks = [
				"0.0.0.0/0",
			  ]
			}
			egress {
			  from_port   = 0
			  to_port     = 0
			  protocol    = "-1"
			  cidr_blocks = ["0.0.0.0/0"]
			}
		  }
	`, awsAccessKey, awsSecretKey, projectID, providerName, region, serviceResourceName)
}
